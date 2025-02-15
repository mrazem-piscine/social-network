package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	"strconv"
)

// RSVPToEventHandler handles user RSVPs to an event
func RSVPToEventHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := strconv.Atoi(r.URL.Query().Get("event_id"))
	if err != nil || eventID == 0 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Status string `json:"status"` // "going" or "not going"
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil || (requestBody.Status != "going" && requestBody.Status != "not going") {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewEventRSVPRepository(db)

	err = repo.RSVPToEvent(eventID, userID, requestBody.Status)
	if err != nil {
		log.Println("Error updating RSVP:", err)
		http.Error(w, "Failed to RSVP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "RSVP updated successfully"})
}

// GetRSVPCountHandler returns the number of users attending an event
func GetRSVPCountHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(r.URL.Query().Get("event_id"))
	if err != nil || eventID == 0 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewEventRSVPRepository(db)

	count, err := repo.GetRSVPCount(eventID)
	if err != nil {
		log.Println("Error retrieving RSVP count:", err)
		http.Error(w, "Failed to retrieve RSVP count", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"going_count": count})
}
