package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/user"
)

// ValidateRegistrationForm used to validate the registration form
func ValidateRegistrationForm(data user.User) map[string][]string {
	// Custom rules
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20"},
		"email":            []string{"required", "email", "min:4", "max:30"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	// Custom error messages
	messages := govalidator.MapData{
		"name": []string{
			"required:Username is required",
			"alpha_num:Invalid format, only alphanumeric characters are allowed",
			"between:Username length must be between 3 and 20",
		},
		"email": []string{
			"required:Email is required",
			"min:Email length must be greater than 4",
			"max:Email length must be less than 30",
			"email:Invalid email format",
		},
		"password": []string{
			"required:Password is required",
			"min:Password length must be greater than 6",
		},
		"password_confirm": []string{
			"required:Password confirmation is required",
		},
	}

	// Validate form input
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", // This is the tag identifier used by govalidator
		Messages:      messages,
	}

	// Start validation
	errs := govalidator.New(opts).ValidateStruct()

	// Custom password confirmation validation rule, which is not included in the govalidator
	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "The two passwords do not match")
	}

	return errs
}
