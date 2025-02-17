package repositories

import (
	"database/sql"
	"errors"
	"log"
	"social-network/internal/models"
)

// GroupRepository handles database operations related to groups
type GroupRepository struct {
	DB *sql.DB
}

// NewGroupRepository creates a new instance of GroupRepository
func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

// CreateGroup inserts a new group into the database
func (repo *GroupRepository) CreateGroup(group *models.Group) error {
	_, err := repo.DB.Exec(`
        INSERT INTO groups (name, description, creator_id, created_at) 
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		group.Name, group.Description, group.CreatorID)
	return err
}

// GroupExists checks if a group with the given name exists
func (repo *GroupRepository) GroupExists(groupName string) (bool, error) {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM groups WHERE name = ?)", groupName).Scan(&exists)
	return exists, err
}

// GroupIDExists checks if a group with the given ID exists
func (repo *GroupRepository) GroupIDExists(groupID int) (bool, error) {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM groups WHERE id = ?)", groupID).Scan(&exists)
	return exists, err
}

// IsUserInGroup checks if a user is already in a group
func (repo *GroupRepository) IsUserInGroup(userID, groupID int) (bool, error) {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)", userID, groupID).Scan(&exists)
	return exists, err
}

// ApproveMembership approves a user's membership request
func (repo *GroupRepository) ApproveMembership(groupID, userID, adminID int) error {
	isAdmin, err := repo.IsUserGroupAdmin(adminID, groupID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("only the group admin can approve members")
	}

	_, err = repo.DB.Exec("UPDATE group_members SET role = 'member' WHERE group_id = ? AND user_id = ? AND role = 'pending'", groupID, userID)
	return err
}

// RejectMembership removes the user from pending requests
func (repo *GroupRepository) RejectMembership(groupID, userID, adminID int) error {
	isAdmin, err := repo.IsUserGroupAdmin(adminID, groupID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("only the group admin can reject members")
	}

	_, err = repo.DB.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ? AND role = 'pending'", groupID, userID)
	return err
}

// GetGroupMembers retrieves all approved members of a group
func (repo *GroupRepository) GetGroupMembers(groupID int) ([]models.GroupMember, error) {
	rows, err := repo.DB.Query(`
        SELECT gm.user_id, u.nickname, gm.status, gm.joined_at 
        FROM group_members gm
        JOIN users u ON gm.user_id = u.id
        WHERE gm.group_id = ? AND (gm.status = 'member' OR gm.status = 'admin')`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.GroupMember
	for rows.Next() {
		var member models.GroupMember
		err := rows.Scan(&member.UserID, &member.Nickname, &member.Status, &member.JoinedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

// IsUserGroupAdmin checks if a user is the admin of a group
func (repo *GroupRepository) IsUserGroupAdmin(userID, groupID int) (bool, error) {
	var isAdmin bool
	err := repo.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM groups WHERE id = ? AND creator_id = ?)", groupID, userID).Scan(&isAdmin)
	return isAdmin, err
}

// RequestToJoinGroup allows a user to request membership in a group
func (repo *GroupRepository) RequestToJoinGroup(groupID, userID int) error {
	log.Println("Inserting into group_members:", groupID, userID) // Debug Log

	_, err := repo.DB.Exec(`
        INSERT INTO group_members (group_id, user_id, status) 
        VALUES (?, ?, 'pending')`, groupID, userID)

	if err != nil {
		log.Println("❌ Failed to request to join group:", err) // Show exact error
	}
	return err
}
