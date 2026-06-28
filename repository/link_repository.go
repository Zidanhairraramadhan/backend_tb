package repository

import (
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

func (r *LinkRepository) GetAllByUserID(userID uint) ([]model.Link, error) {
	var links []model.Link
	err := r.db.Where("user_id = ?", userID).Order("id desc").Find(&links).Error
	return links, err
}

func (r *LinkRepository) GetByID(id uint) (*model.Link, error) {
	var link model.Link
	err := r.db.First(&link, id).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) Update(link *model.Link) error {
	return r.db.Save(link).Error
}

func (r *LinkRepository) Delete(id uint) error {
	return r.db.Delete(&model.Link{}, id).Error
}

func (r *LinkRepository) GetActiveByUserID(userID uint) ([]model.Link, error) {
	var links []model.Link
	err := r.db.Where("user_id = ? AND active = ?", userID, true).Order("id desc").Find(&links).Error
	return links, err
}

func (r *LinkRepository) IncrementClicks(id uint) error {
	return r.db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn("clicks", gorm.Expr("clicks + 1")).Error
}
