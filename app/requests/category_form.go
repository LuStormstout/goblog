package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/category"
)

// ValidateCategoryForm 验证表单，返回 errs 长度等于零即通过
func ValidateCategoryForm(data category.Category) map[string][]string {
	// 定制认证规则
	rules := govalidator.MapData{
		"name": []string{"required", "min:3", "max:32", "not_exists:categories,name"},
	}

	// 定制错误消息
	messages := govalidator.MapData{
		"name": []string{
			"required:分类名称为必填项",
			"min_cn:分类名称长度需大于 3 个字",
			"max_cn:分类名称长度需小于 8 个字",
		},
	}

	// 配置初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	// 开始认证
	return govalidator.New(opts).ValidateStruct()
}
