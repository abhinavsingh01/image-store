package repository

type AuthRepo interface {
	CheckImageAuth(userId int, imageId string) (bool, error)
	CheckAlbumAuth(userId int, albumId string) (bool, error)
}
