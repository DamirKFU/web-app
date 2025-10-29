package auth

type RegisterRequest struct {
	Username string `json:"username" binding:"username"`
	Password string `json:"password" binding:"password"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
