package handlers

import (
	"time"
)

// Incoming: ser struct to receive from the user
type UserIncoming struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Outgoing: User struct to send a response after creation
type UserOutgoing struct {
	ID        any       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// Incoming: Login struct to receive from the user
type UserLoginIncoming struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Outgoing: struct to send the user once they are logged in
type UserLoginOutgoing struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Incoming: refresh token
type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// Outgoing: new access token
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
