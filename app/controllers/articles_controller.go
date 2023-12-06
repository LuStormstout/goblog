package controllers

import (
	"errors"
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"unicode/utf8"
)

// ArticlesController is a struct that groups all the methods related to articles.
// These methods handle HTTP requests related to articles, such as creating, reading, updating, and deleting articles.
type ArticlesController struct {
}

// Index 文章列表页
func (*ArticlesController) Index(
	w http.ResponseWriter,
	_ *http.Request,
) {
	// 获取结果集
	articles, err := article.GetAll()

	// 如果出现错误
	if err != nil {
		// 数据库错误，跳转到 500 错误页面
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, "500 服务器内部错误")
		logger.LogError(err)
		return
	}

	// 加载模板
	view.Render(w, articles, "articles.index")
}

// Show 文章详情页
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	articleInfo, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据未找到，执行 404 处理
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, "404 文章未找到")
			if err != nil {
				logger.LogError(err)
			}
		} else {
			// 数据库错误，执行 500 处理
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			if err != nil {
				logger.LogError(err)
			}
		}
	} else {
		// 读取成功，显示文章
		view.Render(w, articleInfo, "articles.show")
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, _ *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Store 文章添加到数据库中
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	validateFormDataErrors := validateArticleFormData(title, body)

	if len(validateFormDataErrors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		err := _article.Create()
		if err == nil {
			_, err := fmt.Fprint(w, "插入成功, ID 为"+strconv.FormatUint(_article.ID, 10))
			if err != nil {
				logger.LogError(err)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "服务器内部错误")
			if err != nil {
				logger.LogError(err)
			}
		}
	} else {
		view.Render(w, view.D{
			"Title":  title,
			"Body":   body,
			"Errors": validateFormDataErrors,
		}, "articles.create", "articles._form_field")
	}
}

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	_article, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据未找到，执行 404 处理
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误，执行 500 处理
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 读取成功，显示编辑文章表单
		view.Render(w, view.D{
			"Title":   _article.Title,
			"Body":    _article.Body,
			"Article": _article,
			"Errors":  nil,
		}, "articles.edit", "articles._form_field")
	}
}

// Update 文章更新页面
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	_article, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据未找到
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 未出现错误
		// 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")
		validateFromDataErrors := validateArticleFormData(title, body)

		if len(validateFromDataErrors) == 0 {
			// 表单验证通过，更新数据
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()

			// 数据库错误，更新失败
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprint(w, "500 服务器内部错误")
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
				"Title":   title,
				"Body":    body,
				"Article": _article,
				"Errors":  validateFromDataErrors,
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Delete 文章删除页面
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应的文章数据
	_article, err := article.Get(id)

	// 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据未找到
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "500 服务器内部错误")
		}
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

func validateArticleFormData(title, body string) map[string]string {
	validateFromDataErrors := make(map[string]string)
	// 验证标题
	if title == "" {
		validateFromDataErrors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		validateFromDataErrors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		validateFromDataErrors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		validateFromDataErrors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return validateFromDataErrors
}
