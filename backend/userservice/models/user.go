package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        int            `json:"Id" gorm:"primary_key"`
	Username  string         `json:"username" binding:"required"`
	Name      string         `json:"name" binding:"required"`
	Password  string         `json:"password" binding:"required"`
	Email     string         `json:"email" binding:"required"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserView struct {
	Id       int    `json:"Id" gorm:"primary_key"`
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
}

type UserLoginResponse struct {
	Id       int    `json:"Id" gorm:"primary_key"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
