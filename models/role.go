package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
} 