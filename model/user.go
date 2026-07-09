package model

import "time"

// User merepresentasikan tabel 'users' di database.
// Kolom: id (UUID PK), username, password, role, name, bio, genre, country, avatar_initial, verified, created_at
type User struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	Password      string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	Name          string    `gorm:"type:varchar(100)" json:"name"`
	Bio           string    `gorm:"type:text" json:"bio"`
	Genre         string    `gorm:"type:varchar(50)" json:"genre"`
	Country       string    `gorm:"type:varchar(50)" json:"country"`
	AvatarInitial string    `gorm:"type:varchar(5);column:avatar_initial" json:"avatar_initial"`
	AvatarURL     string    `gorm:"type:text;column:avatar_url" json:"avatar_url"`
	CoverURL      string    `gorm:"type:text;column:cover_url" json:"cover_url"`
	Verified      bool      `gorm:"default:false" json:"verified"`
	CreatedAt     time.Time `json:"created_at"`
	Links         []Link    `gorm:"foreignKey:UserID" json:"links,omitempty"`
}

// RegisterRequest adalah body untuk endpoint POST /register
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest adalah body untuk endpoint POST /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangePasswordRequest adalah body untuk endpoint PUT /api/change-password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ProfileRequest adalah body untuk endpoint PUT /api/profile
type ProfileRequest struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Genre     string `json:"genre"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatar_url"`
	CoverURL  string `json:"cover_url"`
}
