package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Posts          []Post         `json:"posts" gorm:"foreignKey:Author"`
	WorkEntries    []WorkEntry    `json:"work_entries" gorm:"foreignKey:Author"`
	ProjectEntries []ProjectEntry `json:"project_entries" gorm:"foreignKey:Author"`
}
