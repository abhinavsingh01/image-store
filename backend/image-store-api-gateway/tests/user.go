package tests

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDetails struct {
	Id       int    `json:"Id" gorm:"primary_key"`
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
}
