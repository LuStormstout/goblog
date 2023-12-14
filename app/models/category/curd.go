package category

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// Create 创建分类，通过 category.ID 来判断是否创建成功
func (category *Category) Create() (err error) {
	if err = model.DB.Create(&category).Error; err != nil {
		return err
	}

	return nil
}

// All 获取全部分类
func All() ([]Category, error) {
	var categories []Category
	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

// Get 通过 ID 获取分类
func Get(stringID string) (Category, error) {
	var category Category
	if err := model.DB.First(&category, types.StringToUint64(stringID)).Error; err != nil {
		return category, err
	}

	return category, nil
}
