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

type StationProduct struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StationID      uint      `gorm:"index;not null" json:"station_id" validate:"required"`
	ProductID      uint      `gorm:"index;not null" json:"product_id" validate:"required"`
	InstalledDate  time.Time `gorm:"not null" json:"installed_date"`
	ExpiryDate     time.Time `gorm:"not null" json:"expiry_date"`
	InspectionDate time.Time `gorm:"not null" json:"inspection_date"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Station Station `gorm:"foreignKey:StationID" json:"station"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
