package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	"social-network/internal/websocket"
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
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	eventRepo := repositories.NewGroupEventRepository(db) // ✅ Correct repository
	notifRepo := repositories.NewNotificationRepository(db)

	// ✅ Store RSVP in the database
	err = eventRepo.RSVPToEvent(eventID, userID, requestBody.Status)
	if err != nil {
		log.Println("❌ Error updating RSVP:", err)
		http.Error(w, "Failed to RSVP", http.StatusInternalServerError)
		return
	}

	// ✅ Get event details
	event, err := eventRepo.GetEventByID(eventID)
	if err != nil {
		log.Println("❌ Error retrieving event:", err)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	// ✅ Notify the event creator
	message := fmt.Sprintf("User %d has RSVP'd as %s to your event: %s", userID, requestBody.Status, event.Title)
	err = notifRepo.CreateNotification(event.CreatorID, message)
	if err != nil {
		log.Println("❌ Error creating notification:", err)
	}

	// ✅ Send real-time WebSocket notification
	go websocket.SendNotification(event.CreatorID, message)

	// ✅ Response
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
