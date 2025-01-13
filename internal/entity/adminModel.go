package entity

import "gorm.io/gorm"

type Manager struct {
	gorm.Model
	ID     int    `gorm:"primary_key:auto_increment" json:"-"`
	UserID string `gorm:"type:varchar(100)" json:"user_id"`
	Dept   int64  `json:"dept"`
}
