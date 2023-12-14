package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
)

// Get 通过 ID 获取文章
func Get(idStr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idStr)
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// GetAll 获取全部文章
func GetAll(r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 1. 初始化分页实例
	db := model.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("home"), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []Article
	_ = _pager.Results(&articles)

	return articles, viewData, nil
}

// Create 创建文章
func (a *Article) Create() (err error) {
	if err = model.DB.Create(&a).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Update 更新文章
func (a *Article) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&a)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// Delete 删除文章
func (a *Article) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&a)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// GetByUserID 获取某个用户的全部文章
func GetByUserID(uid string) (articles []Article, err error) {
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

// GetByCategoryID 获取某个分类下的全部文章
func GetByCategoryID(categoryId string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 初始化分页实例
	db := model.DB.Model(Article{}).Where("category_id = ?", categoryId).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("categories.show"), perPage)

	// 获取视图数据
	viewData := _pager.Paging()

	// 获取数据
	var articles []Article
	_ = _pager.Results(&articles)

	return articles, viewData, nil
}
