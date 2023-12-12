package controllers

import (
	"errors"
	"fmt"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"gorm.io/gorm"
	"net/http"
)

// BaseController is a struct that groups all the methods related to articles.
type BaseController struct {
}

// ResponseForSQLError handles the error of SQL query
func (bc BaseController) ResponseForSQLError(w http.ResponseWriter, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据未找到，执行 404 处理
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "404 Article Not Found")
	} else {
		// 数据库错误，跳转到 500 错误页面
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "500 Internal Server Error")
	}
}

// ResponseForUnauthorized handles the error of unauthorized operation
func (bc BaseController) ResponseForUnauthorized(w http.ResponseWriter, r *http.Request) {
	flash.Warning("Authorization required")
	http.Redirect(w, r, "/", http.StatusFound)
}
