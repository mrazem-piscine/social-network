package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	"strconv"
)

// JoinGroupHandler allows a user to join a group
func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewGroupMemberRepository(db)

	err = repo.AddUserToGroup(groupID, userID, "member")
	if err != nil {
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Joined group successfully"})
}

// LeaveGroupHandler allows a user to leave a group
func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewGroupMemberRepository(db)

	err = repo.RemoveUserFromGroup(groupID, userID)
	if err != nil {
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
}
