package database

import (
	"github.com/charmbracelet/log"
	"github.com/itpourya/Haze/app/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	connStr := "host=redzone-database user=admin password=redzoneadmin dbname=redzone_db port=5432 sslmode=disable TimeZone=Asia/Tehran"
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
