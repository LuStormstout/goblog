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

	// 文章相关页面
	articlesController := new(controllers.ArticlesController)
	router.HandleFunc("/articles", articlesController.Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesController.Show).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles/create", articlesController.Create).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles", articlesController.Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesController.Edit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesController.Update).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesController.Delete).Methods("POST").Name("articles.delete")
}