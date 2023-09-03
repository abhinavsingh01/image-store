package repository

import (
	"errors"
	"imageservice/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ImageRepoImpl struct {
	db *gorm.DB
}

func NewImageRepoImpl(db *gorm.DB) ImageRepo {
	return &ImageRepoImpl{
		db: db,
	}
}

func (repo *ImageRepoImpl) FindAllImages(albumId string) (*[]models.ImageView, error) {
	var imageView []models.ImageView
	// var image models.Image
	// result := repo.db.Model(&image).Select("Id", "ImageName", "ImageUrl").Where("album_id = ?", albumId).Find(&imageView)
	// return &imageView, result.Error

	err := repo.db.Table("images").
		Select("ImageId", "ImageName", "ImageUrl").
		Joins("INNER JOIN albums on albums.id = images.album_id").
		Where("albums.album_id = ?", albumId).
		Where("images.deleted_at IS NULL").
		Find(&imageView).Error

	return &imageView, err
}

func (repo *ImageRepoImpl) GetImageMetadata(imageId string) (*models.ImageView, error) {
	var imageView models.ImageView
	var image models.Image
	result := repo.db.Model(&image).
		Select("ImageId", "ImageName", "ImageUrl").
		Where("image_id = ?", imageId).
		First(&imageView)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("Image not found")
	}
	return &imageView, result.Error
}

func (repo *ImageRepoImpl) Create(albumId string, image *models.Image) error {
	query := `
    INSERT INTO images (image_id, image_url, image_name, album_id) 
    SELECT ?,?,?, albums.id FROM albums WHERE albums.album_id = ?
    `
	result := repo.db.Exec(query, image.ImageId, image.ImageUrl, image.ImageName, albumId)
	return result.Error
}

func (repo *ImageRepoImpl) Delete(imageId string) (*models.Image, error) {
	var image *models.Image
	result := repo.db.Clauses(clause.Returning{}).
		Where("image_id = ?", imageId).
		Delete(&image)
	if result.RowsAffected == 0 {
		return nil, errors.New("Image not found")
	}
	return image, result.Error
}
