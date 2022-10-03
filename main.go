package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

var router = mux.NewRouter()

// ArticleFormData 创建文章的表单数据
type ArticleFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

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

func articlesShowHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	_, err := fmt.Fprint(writer, "文章 ID："+id)
	if err != nil {
		return
	}
}

func articlesIndexHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "访问文章列表")
	if err != nil {
		return
	}
}

func articlesStoreHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.PostFormValue("title")
	body := request.PostFormValue("body")

	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3 - 40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需要大于或等于 10 个字符"
	}

	// 检查是否有错误
	if len(errors) == 0 {
		fmt.Fprint(writer, "通过验证！<br>")
		fmt.Fprintf(writer, "title 的值为：%v <br>", title)
		fmt.Fprintf(writer, "title 的长度为：%v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(writer, "body 的值为：%v <br>", body)
		fmt.Fprintf(writer, "body 的长度为：%v <br>", utf8.RuneCountInString(body))
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticleFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(writer, data)
		if err != nil {
			panic(err)
		}
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

func articlesCreateHandler(writer http.ResponseWriter, request *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticleFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
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
