package routers

import (
	"log"
	"net/http"

	"gwas/handlers"
)

type POSTRouter struct {
	RequestList map[string]handlers.PostRequest
}

func (router *POSTRouter) Route(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST[POST] -- PATH: |%s|", r.URL.Path)
	requestHandler, ok := router.RequestList[r.URL.Path]
	if ok {
		redirect, _, _ := requestHandler.Execute(r, w)
		// if err == nil {
		http.Redirect(w, r, redirect, http.StatusSeeOther)
		// } else {

		// 	log.Fatal("Could not redirect to '" + redirect + "'")
		// }
	} else {

	}
}
