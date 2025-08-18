package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone string `gorm:"size:13;unique;not null" json:"phone"`
}
