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
	router.HandleFunc("/about", pagesController.About).Methods("GET").Name("about")

	// 文章相关页面
	articlesController := new(controllers.ArticlesController)
	router.HandleFunc("/", articlesController.Index).Methods("GET").Name("home")
	router.HandleFunc("/articles", articlesController.Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesController.Show).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles/create", articlesController.Create).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles", articlesController.Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesController.Edit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesController.Update).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesController.Delete).Methods("POST").Name("articles.delete")

	// 用户认证
	authController := new(controllers.AuthController)
	router.HandleFunc("/auth/register", authController.Register).Methods("GET").Name("auth.register")
	router.HandleFunc("/auth/do-register", authController.DoRegister).Methods("POST").Name("auth.do-register")

	// 静态资源
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))
}
