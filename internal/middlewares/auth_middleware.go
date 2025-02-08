package middlewares

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

// Store session globally
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Init registers types for encoding into session
func Init() {
	gob.Register(map[string]interface{}{})
}

// Authenticate ensures a user is logged in before accessing protected routes
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// Check if user is authenticated
		if _, ok := session.Values["user_id"]; !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call next middleware/handler
		next.ServeHTTP(w, r)
	})
}

// SetSession sets user details in session upon login
func SetSession(w http.ResponseWriter, r *http.Request, userID int, nickname string) error {
	session, _ := store.Get(r, "session")
	session.Values["user_id"] = userID
	session.Values["nickname"] = nickname
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
	}
	return session.Save(r, w)
}
func GetUserIDFromSession(r *http.Request) int {
	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}
	return userID
}
// Logout clears user session
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1 // Expire session
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
