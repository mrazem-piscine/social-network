package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	"social-network/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type SessionData struct {
	UserID int `json:"user_id"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("❌ Error decoding request body:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Error hashing password:", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Get database connection
	db := config.GetDB()
	userRepo := repositories.NewUserRepository(db)

	// Insert user into database
	err = userRepo.CreateUser(&user)
	if err != nil {
		log.Println("❌ Error inserting user into database:", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	userRepo := repositories.NewUserRepository(db)

	// ✅ Retrieve user by email or nickname
	user, hashedPassword, err := userRepo.GetUserByEmailOrNickname(credentials.Email)
	if err != nil || user == nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// ✅ Verify password
	if !middlewares.CheckPasswordHash(credentials.Password, hashedPassword) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// ✅ Encode session JSON
	sessionData := SessionData{UserID: user.ID}
	sessionJSON, _ := json.Marshal(sessionData)

	// ✅ Fix the session format: `userID|Base64`
	encodedSession := fmt.Sprintf("%d|%s", user.ID, base64.StdEncoding.EncodeToString(sessionJSON))

	// ✅ Store it in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    encodedSession,
		HttpOnly: true,
		Path:     "/",
	})

	// ✅ Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Login successful",
		"nickname": user.Nickname,
		"user_id":  user.ID,
	})
}

// LogoutUser handles session logout
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	middlewares.Logout(w, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
