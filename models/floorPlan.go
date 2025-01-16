package models

import "time"

type FloorPlan struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Layout     string    `gorm:"type:text" json:"layout"` // JSON or other format for layout data
	CustomerID uint      `gorm:"index" json:"customer_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Stations   []Station `gorm:"foreignKey:FloorPlanID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE" json:"stations"`
}

func (f *FloorPlan) Validate() error {
	return validate.Struct(f)
}
