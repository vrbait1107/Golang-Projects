package database

import (
	"log"
	"os"

	"example.com/package/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {

	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to Connect Database\n", err.Error())
		os.Exit(2)
	}

	log.Println("Database Connected Successfully")

	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Database Migrations")

	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	Database = DbInstance{Db: db}
}
