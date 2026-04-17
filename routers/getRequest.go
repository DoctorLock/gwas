package routers

import (
	"fmt"
	"gwas/handlers"
	"gwas/sessions"
	"log"
	"net/http"
)

type GETRouter struct {
	PageList map[string]handlers.Page
}

func (router *GETRouter) Route(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST[GET] -- PATH: |%s|", r.URL.Path)
	// var pages = map[string]handlers.Page{
	// 	"/":      {Title: "Swipy", Html: "index", RequireAuth: false},
	// 	"/login": {Title: "Swipy - Login", Html: "login", RequireAuth: false},
	// 	"/home":  {Title: "Swipy - Home", Html: "home", RequireAuth: true},
	// }
	page, ok := router.PageList[r.URL.Path]
	if ok {
		if page.RequireAuth {
			_, loggedIn := sessions.GetValue(r, "userId")
			if !loggedIn {
				fmt.Fprintf(w, "<h1>401 Unauthorised</h1><div>You are not allowed to view this page <a href='/'>Go Back</div>")
			}

		}
		handlers.LoadGetRequest(w, &page)
	} else {
		fmt.Fprintf(w, "<h1>404 Not Found</h1><div>This is an error</div>")
	}
}
