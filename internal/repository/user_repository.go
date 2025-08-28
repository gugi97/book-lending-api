package repository

import (
	"book-lending-api/internal/domain"

	"gorm.io/gorm"
)

// UserRepository defines the persistence contract for users.  It is
// implemented against GORM for MySQL but could be swapped out for
// another backend if required.
type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns an implementation of UserRepository
// backed by a gorm.DB instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
