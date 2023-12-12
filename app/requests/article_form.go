package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/article"
)

// ValidateArticleForm used to validate the article form
func ValidateArticleForm(data article.Article) map[string][]string {
	// Custom rules
	rules := govalidator.MapData{
		"title": []string{"required", "min:3", "max:40"},
		"body":  []string{"required", "min:10"},
	}

	// Custom error messages
	messages := govalidator.MapData{
		"title": []string{
			"required:Title is required",
			"min:Title length must be greater than 3",
			"max:Title length must be less than 40",
		},
		"body": []string{
			"required:Body is required",
			"min:Body length must be greater than 10",
		},
	}

	// Set the configuration
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 struct 标签标识符
		Messages:      messages,
	}

	// Start validation
	return govalidator.New(opts).ValidateStruct()
}
