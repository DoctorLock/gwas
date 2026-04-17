package gwas

import (
	"gwas/db"
	"gwas/handlers"
	"gwas/routers"
	"log"
	"net/http"
)

var gwasInited bool = false
var appName string
var databaseLocation string = "/data"
var getRouter routers.GETRouter
var postRouter routers.POSTRouter
var pageList map[string]handlers.Page

func StartGwas(name string, postRequests map[string]handlers.PostRequest, pages map[string]handlers.Page, schemaFileLocation string) {
	appName = name
	pageList = pages
	err := db.Init("./"+databaseLocation+appName+".db", schemaFileLocation)
	if err != nil {
		log.Fatalf("FATAL: Error occured whilst setup of database: %s", err)
	} else {
		log.Println("INFO - Database connection initialisation successful")
	}

	getRouter = routers.GETRouter{PageList: pages}
	postRouter = routers.POSTRouter{RequestList: postRequests}

	http.HandleFunc("/", route)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func route(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRouter.Route(w, r)
	case http.MethodPost:
		postRouter.Route(w, r)
	}
}
