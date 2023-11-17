package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router 路由实例
var Router *mux.Router

// Initialize 初始化路由实例
func Initialize() {
	Router = mux.NewRouter()
}

// RouteName2URL 根据路由名称获取 URL
func RouteName2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		panic(err)
	}

	return url.String()
}

// GetRouteVariable 获取路由变量
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
