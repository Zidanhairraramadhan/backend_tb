package model

import "time"

// Link merepresentasikan tabel 'links' di database.
// Kolom: id (UUID PK), user_id (UUID FK), platform, url, title, image_url, active, clicks, created_at
type Link struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;index;not null;column:user_id" json:"user_id"`
	Platform  string    `gorm:"type:varchar(50);not null" json:"platform"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	URL       string    `gorm:"type:text;not null;column:url" json:"url"`
	ImageURL  string    `gorm:"type:text;column:image_url" json:"image_url"`
	Active    bool      `gorm:"default:true" json:"active"`
	Clicks    int       `gorm:"default:0" json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateLinkRequest adalah body untuk endpoint POST /api/links
type CreateLinkRequest struct {
	Platform string `json:"platform"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
}

// UpdateLinkRequest adalah body untuk endpoint PUT /api/links/:id
type UpdateLinkRequest struct {
	Platform string `json:"platform"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
	Active   *bool  `json:"active"`
}
