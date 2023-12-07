package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/pkg/view"
	"net/http"
)

// AuthController is a struct that groups all the methods related to authentication.
// These methods handle HTTP requests related to user registration and login.
type AuthController struct {
}

// Register is the registration page
func (*AuthController) Register(w http.ResponseWriter, _ *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister handles the registration logic
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// Form validation
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// If validation passes, create user and redirect to homepage
	_user := user.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	_ = _user.Create()

	if _user.ID > 0 {
		_, _ = fmt.Fprint(w, "User created successfully. ID: "+_user.GetStringID())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "Failed to create user.")
	}

	// If validation does not pass, display reason and re-display form
}
