package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")

	if request.URL.Path == "/" {
		_, err := fmt.Fprint(response, "<h1>Hello, 这里是 GoBlog!</h1>")
		if err != nil {
			return
		}
	} else if request.URL.Path == "/about" {
		_, err := fmt.Fprint(response, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
			"<a href=\"mailto:lustormstout@gmail.com\">lustormstout@gmail.com</a>")
		if err != nil {
			return
		}
	} else {
		response.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprint(response, "<h1>请求页面未找到 :(</h1><p>如有疑虑，请联系我们。</p>")
		if err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		return
	}
}
