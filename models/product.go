package models

import (
	"time"
)

// Product Model

type Product struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	ParentID       *uint     `gorm:"index" json:"parent_id"`
	Parent         *Product  `gorm:"foreignKey:ParentID" json:"parent"`
	Price          float64   `gorm:"not null" json:"price"`
	CoverageAmount *float64  `json:"coverage_amount"`
	AgeLimit       *int      `json:"age_limit"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Children       []Product `gorm:"foreignKey:ParentID" json:"children"`
	CategoryID     *uint     `gorm:"index" json:"category_id"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`
}
