package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	"social-network/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input format.", http.StatusBadRequest)
		return
	}

	if user.Nickname == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Nickname, email, and password are required.", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password.", http.StatusInternalServerError)
		return
	}

	// Get repository and store user
	db := config.GetDB()
	userRepo := repositories.NewUserRepository(db)

	err = userRepo.CreateUser(&user, hashedPassword)
	if err != nil {
		http.Error(w, "Registration failed.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	userRepo := repositories.NewUserRepository(db)

	storedUser, storedPassword, err := userRepo.GetUserByEmailOrNickname(user.Email)
	if err != nil || storedUser == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set session
	err = middlewares.SetSession(w, r, storedUser.ID, storedUser.Nickname)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Login successful",
		"user_id":  storedUser.ID,
		"nickname": storedUser.Nickname,
	})
}

// LogoutUser handles session logout
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	middlewares.Logout(w, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
