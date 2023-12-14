package routes

import (
	"goblog/app/controllers"
	"goblog/app/http/middlewares"
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
	router.HandleFunc("/articles/create", middlewares.Auth(articlesController.Create)).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles", middlewares.Auth(articlesController.Store)).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(articlesController.Edit)).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(articlesController.Update)).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(articlesController.Delete)).Methods("POST").Name("articles.delete")

	// 文章分类相关页面
	categoryController := new(controllers.CategoriesController)
	router.HandleFunc("/categories/create", middlewares.Auth(categoryController.Create)).Methods("GET").Name("categories.create")
	router.HandleFunc("/categories", middlewares.Auth(categoryController.Store)).Methods("POST").Name("categories.store")

	// 用户认证
	authController := new(controllers.AuthController)
	router.HandleFunc("/auth/register", middlewares.Guest(authController.Register)).Methods("GET").Name("auth.register")
	router.HandleFunc("/auth/do-register", middlewares.Guest(authController.DoRegister)).Methods("POST").Name("auth.do-register")
	router.HandleFunc("/auth/login", middlewares.Guest(authController.Login)).Methods("GET").Name("auth.login")
	router.HandleFunc("/auth/do-login", middlewares.Guest(authController.DoLogin)).Methods("POST").Name("auth.do-login")
	router.HandleFunc("/auth/logout", middlewares.Auth(authController.Logout)).Methods("POST").Name("auth.logout")

	// 用户相关
	userController := new(controllers.UsersController)
	router.HandleFunc("/users/{id:[0-9]+}", userController.Show).Methods("GET").Name("users.show")

	// 静态资源
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// Middleware: Use the global session middleware
	router.Use(middlewares.StartSession)
}
