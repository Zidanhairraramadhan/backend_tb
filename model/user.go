package model

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	Password      string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	Name          string    `gorm:"type:varchar(100)" json:"name"`
	Bio           string    `gorm:"type:text" json:"bio"`
	Genre         string    `gorm:"type:varchar(50)" json:"genre"`
	Country       string    `gorm:"type:varchar(50)" json:"country"`
	AvatarInitial string    `gorm:"type:varchar(5)" json:"avatar_initial"`
	Verified      bool      `gorm:"default:false" json:"verified"`
	CreatedAt     time.Time `json:"created_at"`
	Links         []Link    `gorm:"foreignKey:UserID" json:"links"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ProfileRequest struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Genre   string `json:"genre"`
	Country string `json:"country"`
}
