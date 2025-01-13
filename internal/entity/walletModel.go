package entity

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	ID      int    `gorm:"primary_key:auto_increment" json:"-"`
	UserID  string `gorm:"type:varchar(100);unique" json:"user_id"`
	Balance int64
}
