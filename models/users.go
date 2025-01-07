package models

import (
	"time"
)

// User Model

type User struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username          string     `gorm:"size:50;not null;unique" json:"username"`
	Email             string     `gorm:"size:100;not null;unique" json:"email"`
	PasswordHash      string     `gorm:"size:255;not null" json:"-"`
	FirstName         string     `gorm:"size:50" json:"first_name"`
	LastName          string     `gorm:"size:50" json:"last_name"`
	PhoneNumber       *string    `gorm:"size:15" json:"phone_number"`
	ProfilePictureURL *string    `gorm:"type:text" json:"profile_picture_url"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	Role              string     `gorm:"type:enum('admin','user','moderator');default:'user'" json:"role"`
	LastLogin         *time.Time `json:"last_login"`
}
