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

	// 3. 加载模板
	tpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)

	// 4. 渲染模板，将所有文章的数据传输进去
	err = tpl.Execute(w, articles)
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

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("文章创建页面"))
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("文章创建页面"))
}

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("文章更新页面"))
}

// Update 文章更新页面
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("文章更新页面"))
}

// Delete 文章删除页面
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("文章删除页面"))
}

// ArticlesFormData 用于存储表单数据
type ArticlesFormData struct {
	Title  string
	Body   string
	Errors map[string]string
}

// Validate 表单验证
func (a ArticlesFormData) Validate() bool {
	a.Errors = make(map[string]string)

	if a.Title == "" {
		a.Errors["Title"] = "标题不能为空"
	}

	if utf8.RuneCountInString(a.Title) < 3 || utf8.RuneCountInString(a.Title) > 40 {
		a.Errors["Title"] = "标题长度需介于 3-40"
	}

	if a.Body == "" {
		a.Errors["Body"] = "内容不能为空"
	}

	if utf8.RuneCountInString(a.Body) < 10 {
		a.Errors["Body"] = "内容长度需大于或等于 10 个字节"
	}

	return len(a.Errors) == 0
}
