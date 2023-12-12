package user

import (
	"goblog/app/models"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/route"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email    string `gorm:"type:varchar(255);unique" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`

	// gorm "-": This field is ignored, it will not be created in the database, nor will it bind parameters, it is only used for validation
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

// ComparePassword compares the user's password with the incoming password
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}

// Link to generate user link
func (user User) Link() string {
	return route.Name2URL("users.show", "id", user.GetStringID())
}

// All gets all users
func All() ([]User, error) {
	var users []User
	if err := model.DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
