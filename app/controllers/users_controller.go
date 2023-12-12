package controllers

import (
	"errors"
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"gorm.io/gorm"
	"net/http"
)

type UsersController struct {
}

func (*UsersController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 获取用户信息
	_user, err := user.Get(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据未找到，执行 404 处理
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, "404 用户未找到")
			if err != nil {
				logger.LogError(err)
			}
		} else {
			// 数据库错误，跳转到 500 错误页面
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			logger.LogError(err)
		}
	} else {
		// 读取成功，显示用户文章列表
		articles, err := article.GetByUserID(_user.GetStringID())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			logger.LogError(err)
			return
		} else {
			view.Render(w, view.D{
				"Articles": articles,
			}, "articles.index", "articles._article_meta")
		}
	}
}
