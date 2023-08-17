package models

import "gorm.io/gorm"

type ProjectEntry struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      uint
}
