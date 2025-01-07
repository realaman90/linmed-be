package models

import "time"

type FloorPlan struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Layout     string    `gorm:"type:text" json:"layout"` // JSON or other format for layout data
	CustomerID uint      `gorm:"index" json:"customer_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Stations   []Station `gorm:"foreignKey:FloorPlanID" json:"stations"`
}

// Station Model
type Station struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CustomerID  uint      `gorm:"index" json:"customer_id"`
	FloorPlanID *uint     `gorm:"index" json:"floor_plan_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Products    []Product `gorm:"foreignKey:StationID" json:"products"`
}
