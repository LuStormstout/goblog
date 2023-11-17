package routes

import (
	"goblog/app/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(router *mux.Router) {
	// 静态页面
	pagesController := new(controllers.PagesController)
	router.NotFoundHandler = http.HandlerFunc(pagesController.NotFound)
	router.HandleFunc("/", pagesController.Home).Methods("GET").Name("home")
	router.HandleFunc("/about", pagesController.About).Methods("GET").Name("about")
}
