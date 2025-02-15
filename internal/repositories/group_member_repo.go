package repositories

import (
	"database/sql"
)

// GroupMemberRepository handles user-group relationships
type GroupMemberRepository struct {
	DB *sql.DB
}

// NewGroupMemberRepository creates a new instance of GroupMemberRepository
func NewGroupMemberRepository(db *sql.DB) *GroupMemberRepository {
	return &GroupMemberRepository{DB: db}
}

// AddUserToGroup adds a user to a group
func (repo *GroupMemberRepository) AddUserToGroup(groupID, userID int, role string) error {
	_, err := repo.DB.Exec(`
		INSERT INTO group_members (group_id, user_id, role) 
		VALUES (?, ?, ?)`, groupID, userID, role)
	return err
}

// RemoveUserFromGroup removes a user from a group
func (repo *GroupMemberRepository) RemoveUserFromGroup(groupID, userID int) error {
	_, err := repo.DB.Exec(`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, userID)
	return err
}

// GetUserGroups retrieves all groups a user belongs to
func (repo *GroupMemberRepository) GetUserGroups(userID int) ([]int, error) {
	rows, err := repo.DB.Query(`SELECT group_id FROM group_members WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupIDs []int
	for rows.Next() {
		var groupID int
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)
	}

	return groupIDs, nil
}
