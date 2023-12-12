package user

import (
	"goblog/pkg/password"
	"gorm.io/gorm"
)

// BeforeSave GORM hook, used to hash the user password before saving to the database.
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if !password.IsHashed(user.Password) {
		user.Password = password.Hash(user.Password)
	}
	return
}
