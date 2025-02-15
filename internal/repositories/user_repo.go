package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/internal/models"
	"strings"
)

// UserRepository handles database operations for users
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	if db == nil {
		log.Fatal("❌ UserRepository received nil database connection")
	}
	return &UserRepository{DB: db}
}

// CreateUser inserts a new user into the database
func (repo *UserRepository) CreateUser(user *models.User) error {
	if repo.DB == nil {
		log.Println("❌ Database connection is nil in CreateUser")
		return sql.ErrConnDone
	}

	_, err := repo.DB.Exec(`
		INSERT INTO users (nickname, email, password, first_name, last_name) 
		VALUES (?, ?, ?, ?, ?)`,
		user.Nickname, user.Email, user.Password, user.FirstName, user.LastName)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("email or nickname already in use")
		}
		log.Printf("❌ Failed to insert user into database: %v", err)
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func (repo *UserRepository) GetUserByEmailOrNickname(identifier string) (*models.User, string, error) {
	var user models.User
	var storedPassword string
	var age sql.NullInt64     // ✅ Handle NULL values for int
	var gender sql.NullString // ✅ Handle NULL values for string

	err := repo.DB.QueryRow(`
		SELECT id, nickname, email, password, age, gender, first_name, last_name 
		FROM users WHERE email = ? OR nickname = ?`, identifier, identifier).
		Scan(&user.ID, &user.Nickname, &user.Email, &storedPassword,
			&age, &gender, // ✅ Store as sql.NullInt64 and sql.NullString
			&user.FirstName, &user.LastName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		log.Println("❌ Error querying user:", err)
		return nil, "", err
	}

	// ✅ Convert `sql.NullInt64` to `*int`
	if age.Valid {
		ageInt := int(age.Int64)
		user.Age = &ageInt
	} else {
		user.Age = nil
	}

	// ✅ Convert `sql.NullString` to `string`
	if gender.Valid {
		user.Gender = gender.String
	} else {
		user.Gender = "" // Default to empty string if NULL
	}

	return &user, storedPassword, nil
}
