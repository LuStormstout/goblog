package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "<h1>Hello, 欢迎来到 go blog!</h1>")
	if err != nil {
		return
	}
}

func aboutHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:lustormstout@gmail.com\">lustormstout@gmail.com</a>")
	if err != nil {
		return
	}
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprint(writer, "<h1>请求页面未找到 :(</h1><p>如有疑虑，请联系我们。</p>")
	if err != nil {
		return
	}
}

func articleShowHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	_, err := fmt.Fprint(writer, "文章 ID："+id)
	if err != nil {
		return
	}
}

func articleIndexHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "访问文章列表")
	if err != nil {
		return
	}
}

func articleStoreHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "创建新的文章")
	if err != nil {
		return
	}
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 1.设置标头
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2.继续处理请求
		next.ServeHTTP(writer, request)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 1.除首页以外，移除所有请求路径后面的斜杠
		if request.URL.Path != "/" {
			request.URL.Path = strings.TrimSuffix(request.URL.Path, "/")
		}
		// 2.将请求传递下去
		next.ServeHTTP(writer, request)
	})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articleShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articleIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articleStoreHandler).Methods("POST").Name("article.store")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL:", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL:", articleURL)

	err := http.ListenAndServe(":3000", removeTrailingSlash(router))
	if err != nil {
		return
	}
}
