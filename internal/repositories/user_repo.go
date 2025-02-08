package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser inserts a new user into the database
func (repo *UserRepository) CreateUser(user *models.User, hashedPassword []byte) error {
	_, err := repo.DB.Exec(`
		INSERT INTO users (nickname, email, password, age, gender, first_name, last_name) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Nickname, user.Email, string(hashedPassword), user.Age, user.Gender, user.FirstName, user.LastName,
	)
	return err
}

// GetUserByEmailOrNickname retrieves a user by email or nickname
func (repo *UserRepository) GetUserByEmailOrNickname(identifier string) (*models.User, string, error) {
	var user models.User
	var storedPassword string

	err := repo.DB.QueryRow(`
		SELECT id, nickname, email, password, age, gender, first_name, last_name 
		FROM users WHERE email = ? OR nickname = ?`, identifier, identifier).
		Scan(&user.ID, &user.Nickname, &user.Email, &storedPassword, &user.Age, &user.Gender, &user.FirstName, &user.LastName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		log.Println("Error querying user:", err)
		return nil, "", err
	}

	return &user, storedPassword, nil
}
