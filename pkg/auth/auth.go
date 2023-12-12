package auth

import (
	"errors"
	"goblog/app/models/user"
	"goblog/pkg/session"
	"gorm.io/gorm"
)

func getUID() string {
	if uid, ok := session.Get("uid").(string); ok && uid != "" {
		return uid
	}
	return ""
}

// User gets the currently logged-in user information
func User() user.User {
	uid := getUID()
	if len(uid) > 0 {
		_user, err := user.Get(uid)
		if err == nil {
			return _user
		}
	}
	return user.User{}
}

// Attempt to authenticate the user with email and password.
func Attempt(email string, password string) error {
	// Get the user information by email
	_user, err := user.GetByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wrong email or password")
		} else {
			return errors.New("internal error")
		}
	}

	// Verify password
	if !_user.ComparePassword(password) {
		return errors.New("wrong email or password")
	}

	// Login user via session
	session.Put("uid", _user.GetStringID())
	return nil
}

// Login the user
func Login(_user user.User) {
	session.Put("uid", _user.GetStringID())
}

// Logout the user
func Logout() {
	session.Forget("uid")
}

// Check if the user is logged in
func Check() bool {
	return len(getUID()) > 0
}
