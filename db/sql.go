package db

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"gwas/handlers/random"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var (
	conn        *sql.DB
	initialized bool
)

const timeout = 5 * time.Second

// --------------------------------------------------
// Init
// --------------------------------------------------
func Init(dbPath string, schemaFilePath string) error {
	if initialized {
		return errors.New("database already initialized")
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Open (creates if not exists)
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return err
	}
	log.Println("Datbase connection successful")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	// Enable required SQLite settings
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
	}

	for _, p := range pragmas {
		if _, err := db.ExecContext(ctx, p); err != nil {
			return err
		}
	}

	// Check if schema already exists
	exists, err := tableExists(ctx, db, "lists")
	if err != nil {
		return err
	}
	conn = db
	initialized = true

	// Only run schema if tables are missing
	if !exists {
		schemaPath := filepath.Join(dir, schemaFilePath)

		schemaBytes, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file at %s: %w", schemaPath, err)
		}

		if _, err := db.ExecContext(ctx, string(schemaBytes)); err != nil {
			return err
		}
		log.Println("INFO - Succesfully created database schema!")
		var password = random.String(12)
		err = CreateUser("admin", password, true)
		if err != nil {
			return err
		}
		log.Println("*****************************************************")
		log.Println("#####################################################")
		log.Println("                DEFAULT ADMIN ACCOUNT                ")
		log.Println("                |USERNAME 'admin'                  ")
		log.Printf("                |PASSWORD '%s'        ", password)
		log.Println("#####################################################")
		log.Println("*****************************************************")

	} else {
		log.Println("INFO - Datbase schema already exists!")

	}

	return nil
}

func CreateUser(username, password string, force_password_reset bool) error {
	if err := checkInit(); err != nil {
		return err
	}

	// Validate inputs
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert user
	_, err = Insert("users", map[string]interface{}{
		"username":               username,
		"password_hash":          string(hash),
		"require_password_reset": force_password_reset,
	})
	return err
}

// --------------------------------------------------
func AuthenticateUser(username, password string) (int, error) {
	if err := checkInit(); err != nil {
		return 0, err
	}

	if username == "" || password == "" {
		return 0, errors.New("username and password cannot be empty")
	}

	query := "SELECT id, password_hash FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := conn.QueryRowContext(ctx, query, username)
	var userId int
	var passwordHash string

	err := row.Scan(&userId, &passwordHash)
	if err == sql.ErrNoRows {
		return 0, errors.New("Username or Password is incorrect")
	} else if err != nil {
		return 0, err
	}

	// Compare hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		// password incorrect
		return 0, errors.New("Username or Password is incorrect")
	}

	return userId, nil
}

// --------------------------------------------------
// Check if table exists
// --------------------------------------------------

func tableExists(ctx context.Context, db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT name 
		FROM sqlite_master 
		WHERE type='table' AND name=?;
	`

	row := db.QueryRowContext(ctx, query, tableName)

	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// --------------------------------------------------
// Close
// --------------------------------------------------

func Close() error {
	if !initialized {
		return errors.New("database not initialized")
	}
	initialized = false
	return conn.Close()
}

// --------------------------------------------------
// Internal check
// --------------------------------------------------

func checkInit() error {
	if !initialized || conn == nil {
		return errors.New("database not initialized")
	}
	return nil
}

// -------------------
// INSERT
// -------------------
func Insert(table string, data map[string]interface{}) (sql.Result, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("insert data cannot be empty")
	}

	columns := []string{}
	values := []interface{}{}
	placeholders := []string{}

	for k, v := range data {
		if strings.TrimSpace(k) == "" {
			return nil, errors.New("invalid column name")
		}
		columns = append(columns, k)
		values = append(values, v)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return conn.ExecContext(ctx, query, values...)
}

// -------------------
// UPDATE
// -------------------
func Update(table string, data map[string]interface{}, where string, whereArgs ...interface{}) (sql.Result, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("update data cannot be empty")
	}

	setParts := []string{}
	values := []interface{}{}

	for k, v := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(setParts, ", "),
		where,
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return conn.ExecContext(ctx, query, append(values, whereArgs...)...)
}

// -------------------
// DELETE
// -------------------
func Delete(table string, where string, whereArgs ...interface{}) (sql.Result, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return conn.ExecContext(ctx, query, whereArgs...)
}

// -------------------
// SELECT
// -------------------
func Select(query string, args ...interface{}) (*sql.Rows, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return conn.QueryContext(ctx, query, args...)
}
