package route

import (
	"goblog/pkg/logger"

	"net/http"

	"github.com/gorilla/mux"
)

// RouteName2URL 根据路由名称获取 URL
func RouteName2URL(routeName string, pairs ...string) string {
	var route *mux.Router
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
	}

	return url.String()
}

// GetRouteVariable 获取路由变量
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
