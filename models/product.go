package models

import (
	"time"
)

// Product Model

func (c *Category) Validate() error {
	return validate.Struct(c)
}

type Product struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Description    string    `gorm:"type:text" json:"description" validate:"required"`
	ImageURL       string    `gorm:"type:text" json:"image_url" validate:"omitempty"`
	ParentID       *uint     `gorm:"index" json:"parent_id" validate:"omitempty"`
	Price          float64   `gorm:"not null" json:"price" validate:"required"`
	CoverageAmount *float64  `json:"coverage_amount" validate:"omitempty"`
	AgeLimit       *int      `json:"age_limit" validate:"omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Children       []Product `gorm:"foreignKey:ParentID" json:"children" validate:"omitempty"`
	CategoryID     *uint     `gorm:"index" json:"category_id" validate:"required"`
	CategoryName   string    `gorm:"-" json:"category_name"`
}

func (p *Product) Validate() error {
	return validate.Struct(p)
}

type CategoryProducts struct {
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Products     []Product `json:"products"`
}

type CategoryProductsResponse struct {
	Categories []CategoryProducts `json:"categories"`
	TotalCount int                `json:"total_count"`
}
