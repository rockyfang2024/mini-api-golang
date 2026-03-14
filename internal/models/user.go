package models

import "time"

// User represents the user model in the application.
type User struct {
	ID        int       `json:"id"`        // unique identifier for the user
	Username  string    `json:"username"`  // username of the user
	Email     string    `json:"email"`     // email address of the user
	Password  string    `json:"password"`  // password of the user
	CreatedAt time.Time `json:"created_at"` // timestamp when the user was created
}