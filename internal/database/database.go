package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/itpourya/Haze/config"
	entity2 "github.com/itpourya/Haze/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tehran", config.DB_HOST, config.DB_USERNAME, config.DB_PASSWORD, config.DATABASE)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("NewDB: ", err)
	}

	log.Info("Connected to Database.")

	err = db.AutoMigrate(&entity2.Product{}, &entity2.User{}, &entity2.Wallet{}, &entity2.Manager{})
	if err != nil {
		panic("Failed : Unable to migrate your sqlite database")
	}

	log.Info("Migrate database models.")

	return db
}

func Close() {
	db := New()
	database, _ := db.DB()
	err := database.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Database closed.")
}
