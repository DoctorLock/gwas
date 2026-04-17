package sessions

import (
	"errors"
	"gwas/db"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	cookieStore *sessions.CookieStore
	sessionName = "SwipeDecide-session"
)

func initCookieStore() {
	if cookieStore == nil {
		// 32-byte key for signing (required)
		key := []byte("super-secret-key-32-bytes-oooo!!")

		cookieStore = sessions.NewCookieStore(key)
		cookieStore.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   36000,
			HttpOnly: true,
			Secure:   false,
		}
	}
}

func SetValue(w http.ResponseWriter, r *http.Request, key string, value any) error {
	initCookieStore()
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[key] = value
	return session.Save(r, w)
}

func GetValue(r *http.Request, key string) (any, bool) {
	initCookieStore()
	session, err := cookieStore.Get(r, sessionName)
	log.Printf("REQUEST[POST] -- ")

	if err != nil {
		return nil, false
	}

	val, ok := session.Values[key]
	return val, ok
}

func Login(w http.ResponseWriter, r *http.Request, username string, password string) (bool, error) {
	log.Printf("REQUEST[POST] -- Starting session")
	_, exists := GetValue(r, "userId")

	if exists {
		return false, errors.New("Already logged in")
	}
	log.Printf("REQUEST[POST] -- Starting checking if authed")
	userId, err := db.AuthenticateUser(username, password)
	if err != nil {
		return false, err
	}
	log.Printf("REQUEST[POST] -- User has succesfully logged in ID: |%s|%s|", userId, err)
	SetValue(w, r, "userId", userId)
	return true, nil
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, sessionName)

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
