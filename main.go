package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var router = mux.NewRouter()

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
	if err := request.ParseForm(); err != nil {
		// 解析错误，这里应该有错误处理
		fmt.Fprint(writer, "请提供正确的数据！")
		return
	}

	title := request.PostForm.Get("title")

	fmt.Fprintf(writer, "POST PostForm: %v <br>", request.PostForm)
	fmt.Fprintf(writer, "POST Form: %v <br>", request.Form)
	fmt.Fprintf(writer, "title 的值为：%v", title)

	fmt.Fprintf(writer, "request.Form 中的 title 值为：%v <br>", request.FormValue("title"))
	fmt.Fprintf(writer, "request.Form 中的 test 值为：%v <br>", request.FormValue("test"))
	fmt.Fprintf(writer, "request.PostForm 中的 test 值为：%v <br>", request.PostFormValue("test"))
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

func articlesCreateHandler(writer http.ResponseWriter, request *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>创建文章 - 我的技术博客</title>
</head>
<body>
	<form action="%s?test=test-get-data" method="post">
		<p><input type="text" name="title"></p>
		<p><textarea name="body" cols="30" rows="10"></textarea></p>
		<p><button type="submit">提交</button></p>
	</form>
</body>
</html>
`
	storeURL, _ := router.Get("articles.store").URL()
	_, err := fmt.Fprintf(writer, html, storeURL)
	if err != nil {
		return
	}
}

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articleShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articleIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articleStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	err := http.ListenAndServe(":3000", removeTrailingSlash(router))
	if err != nil {
		return
	}
}
