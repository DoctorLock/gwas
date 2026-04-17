package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HandlerFunc func(r *PostRequest, w http.ResponseWriter, httpRequest *http.Request) (map[string]string, string, error)

type HTTPAccept string

const (
	HtML HTTPAccept = "HTML"
	JSON HTTPAccept = "JSON"
)

type PostRequest struct {
	RequireAuth bool
	RequestVars map[string]*RequestVar

	Handler HandlerFunc
	Loaded  bool
}

func sanitize(input string) string {
	// remove common injection characters
	replacer := strings.NewReplacer(
		"'", "",
		"\"", "",
		";", "",
		"--", "",
	)
	return replacer.Replace(strings.TrimSpace(input))
}

func (r *PostRequest) Validate(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}
	log.Printf("REQUEST[POST] -- Starting Validation")
	for key, reqVar := range r.RequestVars {
		log.Printf("REQUEST[POST] --- %s", key)

		raw := req.FormValue(key)
		exists := raw != ""

		if reqVar.Required && (!exists || strings.TrimSpace(raw) == "") {
			return fmt.Errorf("missing required field: %s", key)
		}

		if !exists {
			continue
		}

		clean := sanitize(raw)

		switch reqVar.Type {
		case String:
			reqVar.Value = clean

		case Int:
			val, err := strconv.Atoi(clean)
			if err != nil {
				return fmt.Errorf("invalid int for field %s", key)
			}
			reqVar.Value = val

		case Bool:
			val, err := strconv.ParseBool(clean)
			if err != nil {
				return fmt.Errorf("invalid bool for field %s", key)
			}
			reqVar.Value = val

		default:
			return fmt.Errorf("unknown data type for field: %s", key)
		}
	}

	r.Loaded = true
	return nil
}

func (selectedRequest *PostRequest) Execute(r *http.Request, w http.ResponseWriter) (string, map[string]string, error) {
	log.Printf("REQUEST[POST] -- Executing Request")

	// 1. Validate input
	if err := selectedRequest.Validate(r); err != nil {
		return "", nil, err
	}
	log.Printf("REQUEST[POST] -- Request Validated")

	// 2. Ensure handler exists
	if selectedRequest.Handler == nil {
		return "", nil, fmt.Errorf("no handler defined")
	}
	log.Printf("REQUEST[POST] -- Found Handler function")

	// 3. Execute handler
	response, redirect, err := selectedRequest.Handler(selectedRequest, w, r)
	if err != nil {
		return redirect, nil, err
	}
	log.Printf("redirect: %s", redirect)
	// 4. Return results (no HTTP handling here)
	return redirect, response, nil
}
