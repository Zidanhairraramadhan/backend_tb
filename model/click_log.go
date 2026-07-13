package model

import "time"

// ClickLog merepresentasikan tabel 'click_logs' di database.
// Menyimpan log detail setiap kali pengunjung mengklik link di halaman profil publik.
type ClickLog struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LinkID    string    `gorm:"type:uuid;index;not null;column:link_id" json:"link_id"`
	UserID    string    `gorm:"type:uuid;index;not null;column:user_id" json:"user_id"` // Pemilik link (untuk query analitik per user)
	Referrer  string    `gorm:"type:varchar(100)" json:"referrer"`                      // Sumber trafik: Instagram, Twitter, Direct, dll
	Device    string    `gorm:"type:varchar(20)" json:"device"`                         // Mobile atau Desktop
	CreatedAt time.Time `json:"created_at"`
}

// ClickTrackRequest adalah body opsional untuk endpoint POST /api/public/link/:id/click
type ClickTrackRequest struct {
	Referrer string `json:"referrer"`
}
