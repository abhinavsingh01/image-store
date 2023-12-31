package models

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Id        int       `json:"Id" gorm:"primary_key"`
	AlbumId   int       `json:"album_id"`
	ImageId   string    `json:"image_id"`
	ImageName string    `json:"image_name"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type ImageView struct {
	ImageId   string `json:"image_id"`
	ImageName string `json:"image_name"`
	ImageUrl  string `json:"image_url"`
}
