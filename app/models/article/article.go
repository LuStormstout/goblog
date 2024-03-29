package article

import (
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/route"
	"strconv"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title  string `gorm:"column:title;type:varchar(255);not null;" valid:"title"`
	Body   string `gorm:"column:body;type:longtext;not null;" valid:"body"`
	UserID uint64 `gorm:"not null;index"`
	User   user.User

	CategoryID uint64 `gorm:"not null;default:1;index"`
}

// Link 方法用来生成文章链接
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}

// CreatedAtDate 日期格式化
func (a Article) CreatedAtDate() string {
	return a.CreatedAt.Format("2006-01-02")
}
