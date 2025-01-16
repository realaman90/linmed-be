package models

import "time"

// Station Model
type Station struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Description string    `gorm:"type:text" json:"description" validate:"required"`
	CustomerID  uint      `gorm:"index;not null" json:"customer_id"`   // Ensure CustomerID is not null
	FloorPlanID *uint     `gorm:"index;not null" json:"floor_plan_id"` // Ensure FloorPlanID is not null, but it can be empty if nullable
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Products    []Product `gorm:"foreignKey:StationID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE" json:"products"`
}

func (s *Station) Validate() error {
	return validate.Struct(s)
}
