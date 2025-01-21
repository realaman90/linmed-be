package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

type UserRole string

const (
	AdminUser    UserRole = "admin"
	CustomerUser UserRole = "user"
)

type User struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username          string     `gorm:"size:50;not null;unique" json:"username" validate:"required,min=3,max=50"`
	Email             string     `gorm:"size:100;not null;unique" json:"email" validate:"required,email"`
	PasswordHash      string     `gorm:"size:255;not null" json:"password_hash" validate:"required"`
	FirstName         string     `gorm:"size:50" json:"first_name" validate:"omitempty,max=50"`
	LastName          string     `gorm:"size:50" json:"last_name" validate:"omitempty,max=50"`
	PhoneNumber       *string    `gorm:"size:15" json:"phone_number" validate:"omitempty,e164"`
	ProfilePictureURL *string    `gorm:"type:text" json:"profile_picture_url" validate:"omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	Role              string     `gorm:"type:enum('admin','user','moderator');default:'user'" json:"role" validate:"required,oneof=admin user moderator"`
	LastLogin         *time.Time `json:"last_login"`
}

// Validate validates the user model
var validate = validator.New()

func (u *User) Validate() error {
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) Activate() {
	u.IsActive = true
}
