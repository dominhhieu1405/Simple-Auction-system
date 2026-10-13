package repositories

import (
	"errors"
	"simple-auction-system/internal/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

// UserRepository định nghĩa các thao tác với database
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository khởi tạo repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicate
	}
	return err
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
