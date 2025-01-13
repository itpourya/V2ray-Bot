package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID          int    `gorm:"primary_key:auto_increment" json:"-"`
	UserID      string `gorm:"type:varchar(100)" json:"user_id"`
	UsernameSub string `gorm:"type:varchar(100);unique" json:"username_sub"`
	OwnerName   string
}
