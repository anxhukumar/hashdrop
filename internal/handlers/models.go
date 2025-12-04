package handlers

import (
	"time"
)

// User struct to receive from the user
type UserIncoming struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User struct to send a response after creation
type UserOutgoing struct {
	ID        any       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
