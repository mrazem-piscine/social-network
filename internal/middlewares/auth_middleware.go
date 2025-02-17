package middlewares

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Store session globally
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// SessionData represents stored user session data
type SessionData struct {
	UserID int `json:"user_id"`
}

// Init registers types for encoding into session
func Init() {
	gob.Register(map[string]interface{}{})
}

// Authenticate ensures a user is logged in before accessing protected routes
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromSession(r)
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SetSession sets user details in session upon login
func SetSession(w http.ResponseWriter, r *http.Request, userID int) error {
	// ✅ Encode session data into JSON
	sessionData := SessionData{UserID: userID}
	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		log.Println("❌ Failed to encode session data:", err)
		return err
	}

	// ✅ Convert JSON to Base64
	sessionBase64 := base64.StdEncoding.EncodeToString(sessionJSON)

	// ✅ Store session as "userID|Base64Data"
	sessionToken := fmt.Sprintf("%d|%s", userID, sessionBase64)

	// ✅ Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400, // 1 day
	})

	log.Printf("✅ Session set for UserID: %d", userID)
	return nil
}

// GetUserIDFromSession retrieves the user ID from the session cookie
func GetUserIDFromSession(r *http.Request) int {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("❌ No session cookie found")
		return 0
	}

	// ✅ Extract "userID|Base64Data" format
	parts := strings.SplitN(cookie.Value, "|", 2)
	if len(parts) != 2 {
		log.Println("❌ Invalid session token format: missing `|` separator")
		return 0
	}

	// ✅ Extract userID as integer
	var userID int
	_, err = fmt.Sscanf(parts[0], "%d", &userID)
	if err != nil {
		log.Println("❌ Failed to extract userID:", err)
		return 0
	}

	// ✅ Decode Base64 session data
	decodedData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println("❌ Base64 decoding failed:", err)
		return 0
	}

	// ✅ Parse JSON session data
	var sessionData SessionData
	if err := json.Unmarshal(decodedData, &sessionData); err != nil {
		log.Println("❌ Failed to parse session JSON:", err)
		return 0
	}

	log.Printf("🔍 Extracted userID from session: %d", sessionData.UserID)
	return sessionData.UserID
}

// Logout clears user session
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire session immediately
	})
	log.Println("✅ User logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HashPassword hashes a password for secure storage
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies if a given password matches a stored hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("❌ Password mismatch:", err)
		return false
	}
	return true
}
