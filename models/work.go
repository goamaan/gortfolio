package models

import "gorm.io/gorm"

type WorkEntry struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Author      uint
}
