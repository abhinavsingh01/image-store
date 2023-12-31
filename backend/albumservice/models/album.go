package models

import (
	"time"

	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Id               int       `json:"Id" gorm:"primary_key"`
	AlbumId          string    `json:"album_id"`
	UserId           int       `json:"user_id"`
	AlbumName        string    `json:"album_name"`
	AlbumDescription string    `json:"album_description"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type AlbumView struct {
	AlbumId          string `json:"album_id"`
	AlbumName        string `json:"album_name"`
	AlbumDescription string `json:"album_description"`
}
