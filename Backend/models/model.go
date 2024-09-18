package models

import "gorm.io/gorm"

type Book struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	Title     string         `json:"title"`
	Author    string         `json:"author"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
