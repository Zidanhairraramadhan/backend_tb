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

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
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
	// Assuming you want only active links, we can preload with condition if needed,
	// but to follow the instruction we'll preload all. If only active is needed, 
	// we would do .Preload("Links", "active = ?", true)
	err := r.db.Preload("Links").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
