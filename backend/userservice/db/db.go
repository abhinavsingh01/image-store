package db

import (
	"log"
	config "userservice/configs"
	"userservice/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {
	appConfig := config.GetConfig()
	dbURL := appConfig.DBUrl
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	db.AutoMigrate(&models.User{})
	return db
}
