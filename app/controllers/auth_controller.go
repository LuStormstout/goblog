package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
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
	// Initialize a user
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// Validate form input
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// If validation fails, show the error message, redirect to the registration page
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		// If validation passes, create user and redirect to home page
		_ = _user.Create()
		if _user.ID > 0 {
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "Failed to create user, please contact administrator")
		}
	}
}

// Login is the login page
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin handles the login logic
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// Initialize form data
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// Verify credentials
	if err := auth.Attempt(email, password); err == nil {
		// If the login is successful, jump to the home page
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// If the login fails, display the error message and jump to the login page
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}
}

func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	http.Redirect(w, r, "/", http.StatusFound)
}
