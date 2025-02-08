package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
)

// GetNotificationsHandler retrieves a user's notifications
func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db := config.GetDB()
	repo := repositories.NewNotificationRepository(db)
	notifications, err := repo.GetNotifications(userID)
	if err != nil {
		log.Println("Error retrieving notifications:", err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notifications)
}
