package auth

import "time"

// Outgoing: Struct to send new user register data
type NewUserOutgoing struct {
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

// Outgoing: Login struct to send while login
type UserLoginOutgoing struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Incoming: struct to receive tokens after login
type UserLoginIncoming struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Outgoing: Refresh token struct
type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
