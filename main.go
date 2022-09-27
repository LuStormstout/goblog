package main

import (
	"fmt"
	"net/http"
)

//defaultHandle
func defaultHandle(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if request.URL.Path == "/" {
		_, err := fmt.Fprint(writer, "<h1>Hello, 这里是 GoBlog!</h1>")
		if err != nil {
			return
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprint(writer, "<h1>请求页面未找到 :(</h1><p>如有疑虑，请联系我们。</p>")
		if err != nil {
			return
		}
	}
}

//aboutHandle
func aboutHandle(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if request.URL.Path == "/about" {
		_, err := fmt.Fprint(writer, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
			"<a href=\"mailto:lustormstout@gmail.com\">lustormstout@gmail.com</a>")
		if err != nil {
			return
		}
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", defaultHandle)
	router.HandleFunc("/about", aboutHandle)

	// 文章详情
	router.HandleFunc("/articles", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			_, err := fmt.Fprint(writer, "访问文章列表")
			if err != nil {
				return
			}
		case "POST":
			_, err := fmt.Fprint(writer, "创建新的文章")
			if err != nil {
				return
			}
		}
	})

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		return
	}
}
