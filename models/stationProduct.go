package models

import "time"

type StationProduct struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StationID      uint      `gorm:"index;not null" json:"station_id" validate:"required"`
	ProductID      uint      `gorm:"index;not null" json:"product_id" validate:"required"`
	InstalledDate  time.Time `gorm:"not null" json:"installed_date" validate:"required"`
	ExpiryDate     time.Time `gorm:"not null" json:"expiry_date" validate:"required,gtfield=InstalledDate"`
	InspectionDate time.Time `gorm:"not null" json:"inspection_date" validate:"required"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Station  *Station  `gorm:"foreignKey:StationID" json:"-"`
	Product  *Product  `gorm:"foreignKey:ProductID" json:"-"`
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"-"`
}

func (sp *StationProduct) Validate() error {
	return validate.Struct(sp)
}

type ExpiringProductResponse struct {
	ID             int       `json:"id"`
	StationID      int       `json:"station_id"`
	ProductID      int       `json:"product_id"`
	InstalledDate  time.Time `json:"installed_date"`
	ExpiryDate     time.Time `json:"expiry_date"`
	InspectionDate time.Time `json:"inspection_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CustomerName   string    `json:"customer_name"`
	ProductName    string    `json:"product_name"`
}
