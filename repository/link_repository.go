package repository

import (
	"strings"

	"gorm.io/gorm"
	"musiclink-backend/model"
)

type LinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(link *model.Link) error {
	return r.db.Create(link).Error
}

// GetAllByUserID mengambil semua link milik user berdasarkan UUID user
func (r *LinkRepository) GetAllByUserID(userID string) ([]model.Link, error) {
	var links []model.Link
	err := r.db.Where("user_id = ?", userID).Order("position ASC, created_at DESC").Find(&links).Error
	return links, err
}

// LinkQueryParams berisi parameter untuk pagination, search, dan filter
type LinkQueryParams struct {
	Page     int
	Limit    int
	Search   string
	Platform string
}

// GetAllByUserIDPaginated mengambil link dengan pagination, search, dan filter platform
func (r *LinkRepository) GetAllByUserIDPaginated(userID string, params LinkQueryParams) ([]model.Link, int64, error) {
	var links []model.Link
	var total int64

	query := r.db.Model(&model.Link{}).Where("user_id = ?", userID)

	// Filter by platform
	if params.Platform != "" && params.Platform != "all" {
		query = query.Where("LOWER(platform) = ?", strings.ToLower(params.Platform))
	}

	// Search by title or URL
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(url) LIKE ?", searchTerm, searchTerm)
	}

	// Count total records (before pagination)
	query.Count(&total)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	err := query.Order("position ASC, created_at DESC").Offset(offset).Limit(params.Limit).Find(&links).Error

	return links, total, err
}

// CountByUserID mengembalikan jumlah total link milik user
func (r *LinkRepository) CountByUserID(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Link{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetByID mengambil link berdasarkan UUID link
func (r *LinkRepository) GetByID(id string) (*model.Link, error) {
	var link model.Link
	err := r.db.Where("id = ?", id).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) Update(link *model.Link) error {
	return r.db.Save(link).Error
}

// Delete menghapus link berdasarkan UUID link
func (r *LinkRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Link{}).Error
}

// GetActiveByUserID mengambil semua link aktif milik user berdasarkan UUID user
func (r *LinkRepository) GetActiveByUserID(userID string) ([]model.Link, error) {
	var links []model.Link
	err := r.db.Where("user_id = ? AND active = ?", userID, true).Order("position ASC, created_at DESC").Find(&links).Error
	return links, err
}

// IncrementClicks menambah hitungan klik pada sebuah link berdasarkan UUID link
func (r *LinkRepository) IncrementClicks(id string) error {
	return r.db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn("clicks", gorm.Expr("clicks + 1")).Error
}

// ReorderItem merepresentasikan satu item dalam request reorder
type ReorderItem struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// ReorderLinks memperbarui posisi beberapa link sekaligus secara batch
func (r *LinkRepository) ReorderLinks(userID string, items []ReorderItem) error {
	tx := r.db.Begin()
	for _, item := range items {
		if err := tx.Model(&model.Link{}).Where("id = ? AND user_id = ?", item.ID, userID).UpdateColumn("position", item.Position).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

