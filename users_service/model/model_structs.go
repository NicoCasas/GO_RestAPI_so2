package model

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token"`
}

type OSUser struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}
