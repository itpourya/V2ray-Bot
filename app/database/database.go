package database

import (
	"fmt"
	"log"

	"github.com/itpourya/Haze/app/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	database = "postgres"
	password = "docker"
	username = "postgres"
	port     = "5432"
	host     = "localhost"
)

func New() *gorm.DB {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Tehran", host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("NewDB: ", err)
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
