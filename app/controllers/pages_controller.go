package controllers

import "net/http"

// PagesController 处理静态页面
type PagesController struct {
}

// Home 首页
func (*PagesController) Home(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("首页"))
}

// About 关于我们页面
func (*PagesController) About(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("关于我们"))
}

// NotFound 404 页面
func (*PagesController) NotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("404 Not Found"))
}
