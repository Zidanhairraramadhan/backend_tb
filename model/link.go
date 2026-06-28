package model

import "time"

type Link struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Platform  string    `gorm:"type:varchar(50);not null" json:"platform"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	Active    bool      `gorm:"default:true" json:"active"`
	Clicks    int       `gorm:"default:0" json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLinkRequest struct {
	Platform string `json:"platform" validate:"required"`
	Title    string `json:"title" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
}

type UpdateLinkRequest struct {
	Platform string `json:"platform"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Active   *bool  `json:"active"`
}
