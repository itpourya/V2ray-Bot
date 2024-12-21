package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID          int    `gorm:"primary_key:auto_increment" json:"-"`
	ProductName string `gorm:"type:varchar(100)" json:"productname"`
	Price       string `gorm:"type:varchar(100)" json:"price"`
	Expire      int32
}
