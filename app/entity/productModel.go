package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID          int    `gorm:"primary_key:auto_increment" json:"-"`
	ProductName string `gorm:"type:varchar(100)" json:"productname"`
	Price       string `gorm:"type:varchar(100)" json:"price"`
	Expire      int32
}

type User struct {
	gorm.Model
	ID          int    `gorm:"primary_key:auto_increment" json:"-"`
	UserID      string `gorm:"type:varchar(100)" json:"user_id"`
	UsernameSub string `gorm:"type:varchar(100);unique" json:"username_sub"`
}

type Wallet struct {
	gorm.Model
	ID      int    `gorm:"primary_key:auto_increment" json:"-"`
	UserID  string `gorm:"type:varchar(100);unique" json:"user_id"`
	Balance int64
}
