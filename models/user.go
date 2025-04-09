package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primarykey"`
	Username   string `json:"username" gorm:"unique" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
