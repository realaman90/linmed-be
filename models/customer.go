package models

import "time"

// Customer Model
// New model for customers who own floor plans and stations.
type Customer struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string      `gorm:"size:100;not null" json:"name"`
	Email      string      `gorm:"size:100;not null;unique" json:"email"`
	Phone      string      `gorm:"size:15" json:"phone"`
	Address    string      `gorm:"size:255" json:"address"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	FloorPlans []FloorPlan `gorm:"foreignKey:CustomerID" json:"floor_plans"`
	Stations   []Station   `gorm:"foreignKey:CustomerID" json:"stations"`
}
