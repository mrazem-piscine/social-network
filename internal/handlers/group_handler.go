package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	"social-network/internal/repositories"
	"strconv"
)

// CreateGroupHandler allows users to create a new group
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if group.Name == "" || group.Description == "" {
		http.Error(w, "Group name and description are required", http.StatusBadRequest)
		return
	}

	group.CreatorID = userID

	db := config.GetDB()
	repo := repositories.NewGroupRepository(db)

	err := repo.CreateGroup(&group)
	if err != nil {
		log.Println("Error creating group:", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group created successfully"})
}

// RequestToJoinGroupHandler allows a user to request joining a group
func RequestToJoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil || groupID == 0 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewGroupRepository(db)

	err = repo.RequestToJoinGroup(groupID, userID)
	if err != nil {
		http.Error(w, "Failed to request to join group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request sent"})
}

// ApproveMembershipHandler allows group admins to approve members
func ApproveMembershipHandler(w http.ResponseWriter, r *http.Request) {
	adminID := middlewares.GetUserIDFromSession(r)
	if adminID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	userID, err2 := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || err2 != nil || groupID == 0 || userID == 0 {
		http.Error(w, "Invalid group or user ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewGroupRepository(db)

	err = repo.ApproveMembership(groupID, userID, adminID)
	if err != nil {
		log.Println("Error approving member:", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Membership approved"})
}
// GetGroupMembersHandler retrieves all approved members of a group
func GetGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil || groupID == 0 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewGroupRepository(db)

	members, err := repo.GetGroupMembers(groupID)
	if err != nil {
		log.Println("Error getting group members:", err)
		http.Error(w, "Failed to fetch members", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(members)
}

// RejectMembershipHandler allows group admins to reject members
func RejectMembershipHandler(w http.ResponseWriter, r *http.Request) {
	adminID := middlewares.GetUserIDFromSession(r)
	if adminID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	userID, err2 := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || err2 != nil || groupID == 0 || userID == 0 {
		http.Error(w, "Invalid group or user ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewGroupRepository(db)

	err = repo.RejectMembership(groupID, userID, adminID)
	if err != nil {
		log.Println("Error rejecting member:", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Membership rejected"})
}
