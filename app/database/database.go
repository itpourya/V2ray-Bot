package database

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/itpourya/Haze/app/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Can't load database environment file", err)
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Tehran", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DATABASE"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("NewDB: ", err)
	}

	log.Info("Connected to Database.")

	err = db.AutoMigrate(&entity.Product{}, &entity.User{}, &entity.Wallet{}, &entity.Manager{})
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
