package models

import "time"

// Customer Model
// New model for customers who own floor plans and stations.
type Customer struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string      `gorm:"size:100;not null" json:"name" validate:"required"`
	Email      string      `gorm:"size:100;not null;unique" json:"email" validate:"required,email"`
	Phone      string      `gorm:"size:15" json:"phone" validate:"omitempty,e164"`
	Address    string      `gorm:"size:255" json:"address" validate:"omitempty"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	FloorPlans []FloorPlan `gorm:"foreignKey:CustomerID" json:"floor_plans"`
	Stations   []Station   `gorm:"foreignKey:CustomerID" json:"stations"`
}

func (c *Customer) Validate() error {
	return validate.Struct(c)
}
