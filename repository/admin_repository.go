package repository

import (
	"gorm.io/gorm"
	"musiclink-backend/model"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

type GlobalStats struct {
	TotalUsers  int64
	TotalLinks  int64
	TotalClicks int64
	RecentUsers []model.User
}

func (r *AdminRepository) GetGlobalStats() (*GlobalStats, error) {
	var stats GlobalStats

	// Count total users
	if err := r.db.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Count total links
	if err := r.db.Model(&model.Link{}).Count(&stats.TotalLinks).Error; err != nil {
		return nil, err
	}

	// Sum total clicks
	type clickResult struct {
		Total int64
	}
	var cResult clickResult
	if err := r.db.Model(&model.Link{}).Select("COALESCE(SUM(clicks), 0) as total").Scan(&cResult).Error; err != nil {
		return nil, err
	}
	stats.TotalClicks = cResult.Total

	// Get 5 recent users
	if err := r.db.Model(&model.User{}).Select("id, username, role, created_at").Order("created_at desc").Limit(5).Find(&stats.RecentUsers).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
