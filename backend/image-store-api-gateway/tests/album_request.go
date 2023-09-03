package tests

type AlbumRequest struct {
	AlbumName string `json:"album_name" binding:"required"`
	AlbumDesc string `json:"album_description"`
}
