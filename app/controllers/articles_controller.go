package controllers

import (
	"errors"
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. 获取结果集
	articles, err := article.GetAll()

	// 2. 如果出现错误
	if err != nil {
		// 2.1 数据库错误，跳转到 500 错误页面
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, "500 服务器内部错误")
		logger.LogError(err)
		return
	}

	// 设置模板相对路径
	viewDir := "resources/views"

	// 所有布局模板文件 Slice
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logger.LogError(err)

	// 在 Slice 里新增我们的目标文件
	newFiles := append(files, viewDir+"/articles/index.gohtml")

	tpl, err := template.ParseFiles(newFiles...)
	logger.LogError(err)

	// 4. 渲染模板，将所有文章的数据传输进去
	err = tpl.ExecuteTemplate(w, "app", articles)
	logger.LogError(err)
}

// Show 文章详情页
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	articleInfo, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到，执行 404 处理
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, "404 文章未找到")
			if err != nil {
				logger.LogError(err)
			}
		} else {
			// 3.2 数据库错误，执行 500 处理
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			if err != nil {
				logger.LogError(err)
			}
		}
	} else {
		// 4. 读取成功，显示文章
		tpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"Name2URL":       route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		err = tpl.Execute(w, articleInfo)
		if err != nil {
			logger.LogError(err)
		}
	}
}

// ArticlesFormData 用于存储表单数据
type ArticlesFormData struct {
	Title, Body, URL string
	Errors           map[string]string
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		URL:    storeURL,
		Body:   "",
		Title:  "",
		Errors: nil,
	}
	tpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
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
		storeURL := route.Name2URL("articles.store")
		data := map[string]interface{}{
			"URL":   storeURL,
			"Title": title,
			"Body":  body,
			"Error": validateFormDataErrors,
		}
		tpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			logger.LogError(err)
		}

		err = tpl.Execute(w, data)
		if err != nil {
			logger.LogError(err)
		}
	}
}

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到，执行 404 处理
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, "404 文章未找到")
			if err != nil {
				logger.LogError(err)
			}
		} else {
			// 3.2 数据库错误，执行 500 处理
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			if err != nil {
				logger.LogError(err)
			}
		}
	} else {
		// 4. 读取成功，显示文章
		updateURL := route.Name2URL("articles.update", "id", id)
		data := ArticlesFormData{
			URL:    updateURL,
			Title:  _article.Title,
			Body:   _article.Body,
			Errors: nil,
		}
		tpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		logger.LogError(err)
		err = tpl.Execute(w, data)
		if err != nil {
			logger.LogError(err)
		}
	}
}

// Update 文章更新页面
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误
		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")
		validateFromDataErrors := validateArticleFormData(title, body)

		if len(validateFromDataErrors) == 0 {
			// 4.2 表单验证通过，更新数据
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()

			// 数据库错误，更新失败
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprint(w, "500 服务器内部错误")
			}

			// √ 更新成功，跳转到文章详情页
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				_, _ = fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {
			// 4.3 表单验证不通过，显示理由
			updateURL := route.Name2URL("articles.update", "id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: validateFromDataErrors,
			}
			tpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			err = tpl.Execute(w, data)
			logger.LogError(err)
		}
	}
}

// Delete 文章删除页面
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := _article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 SQL 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
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
