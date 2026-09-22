package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}