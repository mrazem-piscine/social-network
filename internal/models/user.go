package models

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"` // Exclude from JSON output
	Age       *int   `json:"age,omitempty"`      // ✅ Change to a pointer to handle NULL values
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
