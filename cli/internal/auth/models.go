package auth

import "time"

// Outgoing: Struct to send new user register data
type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Incoming: Struct to receive new user data
type NewUserIncoming struct {
	ID        any       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
