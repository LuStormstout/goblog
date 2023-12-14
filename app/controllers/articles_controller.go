package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// ArticlesController is a struct that groups all the methods related to articles.
// These methods handle HTTP requests related to articles, such as creating, reading, updating, and deleting articles.
type ArticlesController struct {
	BaseController
}

// Index 文章列表页
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. 获取结果集
	articles, pagerData, err := article.GetAll(r, config.GetInt("pagination.per_page"))

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {

		// ---  2. 加载模板 ---
		view.Render(w, view.D{
			"Articles":  articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

// Show 文章详情页
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	articleInfo, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 读取成功，显示文章
		view.Render(w, view.D{
			"Article":          articleInfo,
			"CanModifyArticle": policies.CanModifyArticle(articleInfo),
		}, "articles.show", "articles._article_meta")
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, _ *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Store 文章添加到数据库中
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	// 初始化数据
	currentUser := auth.User()
	_article := article.Article{
		Title:  r.PostFormValue("title"),
		Body:   r.PostFormValue("body"),
		UserID: currentUser.ID,
	}

	validateFormDataErrors := requests.ValidateArticleForm(_article)

	if len(validateFormDataErrors) == 0 {
		_ = _article.Create()
		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "创建文章失败，请联系管理员")
			if err != nil {
				logger.LogError(err)
			}
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  validateFormDataErrors,
		}, "articles.create", "articles._form_field")
	}
}

// Edit 文章更新页面
func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	_article, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 读取成功，显示编辑文章表单
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  view.D{},
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Update 文章更新页面
func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	_article, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 未出现错误

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 表单验证
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")
			validateFromDataErrors := requests.ValidateArticleForm(_article)

			if len(validateFromDataErrors) == 0 {
				// 表单验证通过，更新数据
				rowsAffected, err := _article.Update()

				// 数据库错误，更新失败
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprint(w, "500 服务器内部错误")
					return
				}

				// 更新成功，跳转到文章详情页
				if rowsAffected > 0 {
					showURL := route.Name2URL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					_, _ = fmt.Fprint(w, "您没有做任何更改！")
				}
			} else {
				// 表单验证不通过，显示理由
				view.Render(w, view.D{
					"Article": _article,
					"Errors":  validateFromDataErrors,
				}, "articles.edit", "articles._form_field")
			}
		}
	}
}

// Delete 文章删除页面
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	_article, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 未出现错误，执行删除操作
			rowsAffected, err := _article.Delete()

			// 发生错误
			if err != nil {
				// 应该是 SQL 报错了
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprint(w, "500 服务器内部错误")
			} else {
				// 未发生错误
				if rowsAffected > 0 {
					// 重定向到文章列表页
					indexURL := route.Name2URL("articles.index")
					http.Redirect(w, r, indexURL, http.StatusFound)
				} else {
					// Edge case
					w.WriteHeader(http.StatusNotFound)
					_, _ = fmt.Fprint(w, "404 文章未找到")
				}
			}
		}
	}
}
