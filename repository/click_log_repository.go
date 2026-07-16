package repository

import (
	"time"

	"gorm.io/gorm"
	"musiclink-backend/model"
)

type ClickLogRepository struct {
	db *gorm.DB
}

func NewClickLogRepository(db *gorm.DB) *ClickLogRepository {
	return &ClickLogRepository{db: db}
}

// Create menyimpan entri log klik baru ke database
func (r *ClickLogRepository) Create(log *model.ClickLog) error {
	return r.db.Create(log).Error
}

// DailyClickStats merepresentasikan jumlah klik per hari
type DailyClickStats struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

// MonthlyClickStats merepresentasikan jumlah klik per bulan
type MonthlyClickStats struct {
	Month  string `json:"month"`
	Clicks int    `json:"clicks"`
}

// SourceStats merepresentasikan jumlah klik per sumber trafik
type SourceStats struct {
	Source     string  `json:"source"`
	Clicks    int     `json:"clicks"`
	Percentage float64 `json:"percentage"`
}

// GetDailyClicks mengambil jumlah klik per hari dalam 7 hari terakhir untuk user tertentu
func (r *ClickLogRepository) GetDailyClicks(userID string, days int) ([]DailyClickStats, error) {
	var results []DailyClickStats
	since := time.Now().AddDate(0, 0, -days)

	if r.db.Dialector.Name() == "postgres" {
		// Syntax untuk PostgreSQL
		err := r.db.Model(&model.ClickLog{}).
			Select("created_at::date as date, COUNT(*) as clicks").
			Where("user_id = ? AND created_at >= ?", userID, since).
			Group("created_at::date").
			Order("date ASC").
			Find(&results).Error
		return results, err
	}

	// Syntax default untuk SQLite
	err := r.db.Model(&model.ClickLog{}).
		Select("DATE(created_at) as date, COUNT(*) as clicks").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&results).Error

	return results, err
}

// GetMonthlyClicks mengambil jumlah klik per bulan dalam 12 bulan terakhir untuk user tertentu
func (r *ClickLogRepository) GetMonthlyClicks(userID string) ([]MonthlyClickStats, error) {
	var results []MonthlyClickStats
	since := time.Now().AddDate(-1, 0, 0) // 12 bulan terakhir

	if r.db.Dialector.Name() == "postgres" {
		// Syntax untuk PostgreSQL (TO_CHAR)
		err := r.db.Model(&model.ClickLog{}).
			Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as clicks").
			Where("user_id = ? AND created_at >= ?", userID, since).
			Group("TO_CHAR(created_at, 'YYYY-MM')").
			Order("month ASC").
			Find(&results).Error
		return results, err
	}

	// Syntax default untuk SQLite (strftime)
	err := r.db.Model(&model.ClickLog{}).
		Select("strftime('%Y-%m', created_at) as month, COUNT(*) as clicks").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Group("strftime('%Y-%m', created_at)").
		Order("month ASC").
		Find(&results).Error

	return results, err
}

// GetMonthlyClicksPostgres (Deprecated - dipertahankan untuk backward compatibility)
func (r *ClickLogRepository) GetMonthlyClicksPostgres(userID string) ([]MonthlyClickStats, error) {
	return r.GetMonthlyClicks(userID)
}

// GetTrafficSources mengambil pengelompokan klik berdasarkan sumber trafik untuk user tertentu
func (r *ClickLogRepository) GetTrafficSources(userID string) ([]SourceStats, error) {
	var results []SourceStats

	// Pertama ambil total klik
	var totalClicks int64
	r.db.Model(&model.ClickLog{}).Where("user_id = ?", userID).Count(&totalClicks)

	if totalClicks == 0 {
		return results, nil
	}

	err := r.db.Model(&model.ClickLog{}).
		Select("referrer as source, COUNT(*) as clicks").
		Where("user_id = ?", userID).
		Group("referrer").
		Order("clicks DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// Hitung persentase
	for i := range results {
		results[i].Percentage = float64(results[i].Clicks) / float64(totalClicks) * 100
	}

	return results, nil
}

// GetTotalClicksByUser mengambil total klik untuk user tertentu
func (r *ClickLogRepository) GetTotalClicksByUser(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.ClickLog{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetClicksThisWeek mengambil total klik dalam 7 hari terakhir
func (r *ClickLogRepository) GetClicksThisWeek(userID string) (int64, error) {
	var count int64
	since := time.Now().AddDate(0, 0, -7)
	err := r.db.Model(&model.ClickLog{}).Where("user_id = ? AND created_at >= ?", userID, since).Count(&count).Error
	return count, err
}

// GetClicksLastWeek mengambil total klik dalam 7-14 hari yang lalu (untuk perhitungan growth)
func (r *ClickLogRepository) GetClicksLastWeek(userID string) (int64, error) {
	var count int64
	since := time.Now().AddDate(0, 0, -14)
	until := time.Now().AddDate(0, 0, -7)
	err := r.db.Model(&model.ClickLog{}).Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, since, until).Count(&count).Error
	return count, err
}
