package database

import (
	"log"

	"github.com/itpourya/Haze/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&entity.Product{}, &entity.User{}, &entity.Wallet{})
	if err != nil {
		panic("Failed : Unable to migrate your sqlite database")
	}

	return db
}

func Close() {
	db := New()
	database, _ := db.DB()
	err := database.Close()
	if err != nil {
		log.Fatal(err)
	}
}
