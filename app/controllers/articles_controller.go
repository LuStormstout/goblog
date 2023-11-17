package controllers

import (
	"net/http"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// // 1. 查询文章数据
	// rows, err := db.Query("SELECT * FROM articles")
	// logger.LogError(err)
	// defer rows.Close()

	// var articles []Article
	// // 2. 遍历查询结果
	// for rows.Next() {
	// 	article := Article{}
	// 	// 2.1 将每一行的结果都赋值到一个 Article 对象中
	// 	err := rows.Scan(&article.ID, &article.Title, &article.Body)
	// 	logger.LogError(err)
	// 	// 2.2 将 Article 对象追加到 articles 的这个数组中
	// 	articles = append(articles, article)
	// }

	// // 3. 检测遍历时是否发生错误
	// err = rows.Err()
	// logger.LogError(err)

	// // 4. 加载模板
	// template, err := template.ParseFiles("resources/views/articles/index.gohtml")
	// logger.LogError(err)

	// // 5. 渲染模板，将所有文章的数据传输进去
	// err = template.Execute(w, articles)
	// logger.LogError(err)
}

// Show 文章详情页
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// // 1. 获取 URL 参数
	// id := route.GetRouteVariable("id", r)

	// // 2. 读取对应的文章数据
	// article, err := getArticleByID(id)

	// // 3. 如果出现错误
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		// 3.1 数据未找到，执行 404 处理
	// 		w.WriteHeader(http.StatusNotFound)
	// 		fmt.Fprint(w, "404 文章未找到")
	// 	} else {
	// 		// 3.2 数据库错误，执行 500 处理
	// 		logger.LogError(err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprint(w, "500 服务器内部错误")
	// 	}
	// } else {
	// 	// 4. 读取成功，显示文章
	// 	template, err := template.New("show.gohtml").
	// 		Funcs(template.FuncMap{
	// 			"RouteName2URL": route.RouteName2URL,
	// 			"Int64ToString": types.Int64ToString,
	// 		}).
	// 		ParseFiles("resources/views/articles/show.gohtml")
	// 	logger.LogError(err)
	// 	template.Execute(w, article)
	// }
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
