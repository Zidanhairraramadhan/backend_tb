package repository

import (
	"gorm.io/gorm"
	"musiclink-backend/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID mencari user berdasarkan UUID (string)
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// GetAllWithLinks returns all users with their associated links preloaded
func (r *UserRepository) GetAllWithLinks() ([]model.User, error) {
	var users []model.User
	err := r.db.Preload("Links").Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByUsernameWithLinks returns a user with their associated links preloaded
func (r *UserRepository) GetByUsernameWithLinks(username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Links").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetPublicProfileByUsername fetches user by username with active links ordered by created_at desc
func (r *UserRepository) GetPublicProfileByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Links", func(db *gorm.DB) *gorm.DB {
		return db.Where("active = ?", true).Order("created_at desc")
	}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// IncrementLinkClick increments the click count of a specific link
func (r *UserRepository) IncrementLinkClick(linkID string) error {
	return r.db.Model(&model.Link{}).Where("id = ?", linkID).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error
}
