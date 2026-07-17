package routers

import (
	"fmt"
	"gwas/handlers"
	"gwas/sessions"
	"log"
	"net/http"
	"strconv"
)

type GETRouter struct {
	PageList map[string]handlers.Page
}

func (router *GETRouter) Route(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST[GET] -- PATH: |%s|", r.URL.Path)
	var userId int
	page, ok := router.PageList[r.URL.Path]
	if ok {
		if page.RequireAuth {
			var loggedIn bool
			userIdValue, loggedIn := sessions.GetValue(r, "userId")
			if loggedIn {
				userId = userIdValue.(int)
				if page.PageData == nil {
					page.PageData = make(map[string]string)
				}
				page.PageData["userId"] = strconv.Itoa(userId)
				if usernameVal, _ := sessions.GetValue(r, "username"); usernameVal != nil {
					if usernameStr, ok := usernameVal.(string); ok {
						page.PageData["username"] = usernameStr
					} else {
						page.PageData["username"] = ""
					}
				} else {
					page.PageData["username"] = ""
				}
				log.Printf("REQUEST[GET] -- User is logged in with userId: %d and username: %s", userIdValue, page.PageData["username"])

			} else {
				fmt.Fprintf(w, "<h1>401 Unauthorised</h1><div>You are not allowed to view this page <a href='/'>Go Back</div>")
				return
			}

		}
		page.LoadGetRequest(w, r)
	} else {
		fmt.Fprintf(w, "<h1>404 Not Found</h1><div>This is an error</div>")
	}
}
