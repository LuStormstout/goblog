package session

import (
	"github.com/gorilla/sessions"
	"goblog/pkg/logger"
	"net/http"
)

// Store the cookie store
var Store = sessions.NewCookieStore([]byte("AHO99RL1ApfPUQXIgUo3jPWZMpfL5y4o"))

// Session Get a session instance
var Session *sessions.Session

// Request the request instance
var Request *http.Request

// Response the response instance
var Response http.ResponseWriter

// StartSession initialize a session instance, must be called before using Session
func StartSession(w http.ResponseWriter, r *http.Request) {
	var err error

	// Get a session. Get() always returns a session, even if empty.
	Session, err = Store.Get(r, "goblog-session")
	if err != nil {
		panic(err)
	}

	Request = r
	Response = w
}

// Put a key / value pair to the session data
func Put(key string, value interface{}) {
	Session.Values[key] = value
	Save()
}

// Get a key from the session data
func Get(key string) interface{} {
	return Session.Values[key]
}

// Forget Delete a key from the session data
func Forget(key string) {
	delete(Session.Values, key)
	Save()
}

// Flush Delete all data from the session
func Flush() {
	Session.Options.MaxAge = -1
	Save()
}

// Save the session data to the storage
func Save() {
	// Not HTTPS connection cannot use Secure and HttpOnly, The browser will report an error
	// Session.Options.Secure = true
	// Session.Options.HttpOnly = true
	err := Session.Save(Request, Response)
	logger.LogError(err)
}
