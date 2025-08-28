package repository

import (
	"book-lending-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

// LendingRepository provides persistence methods for lending records.
// Implementations are responsible for handling associations.
type LendingRepository interface {
	Create(record *domain.LendingRecord) error
	GetByID(id uint) (*domain.LendingRecord, error)
	GetActiveByUserAndBook(userID, bookID uint) (*domain.LendingRecord, error)
	Update(record *domain.LendingRecord) error
	GetUserBorrowingHistory(userID uint, offset, limit int) ([]domain.LendingRecord, int64, error)
	GetActiveBorrowingsByUser(userID uint) ([]domain.LendingRecord, error)
	CountUserBorrowsInPeriod(userID uint, since time.Time) (int64, error)
}

type lendingRepository struct {
	db *gorm.DB
}

// NewLendingRepository returns a new LendingRepository using the
// provided gorm DB.
func NewLendingRepository(db *gorm.DB) LendingRepository {
	return &lendingRepository{db: db}
}

func (r *lendingRepository) Create(record *domain.LendingRecord) error {
	return r.db.Create(record).Error
}

func (r *lendingRepository) GetByID(id uint) (*domain.LendingRecord, error) {
	var record domain.LendingRecord
	if err := r.db.Preload("Book").Preload("User").First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *lendingRepository) GetActiveByUserAndBook(userID, bookID uint) (*domain.LendingRecord, error) {
	var record domain.LendingRecord
	if err := r.db.Where("user_id = ? AND book_id = ? AND return_date IS NULL", userID, bookID).
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *lendingRepository) Update(record *domain.LendingRecord) error {
	return r.db.Save(record).Error
}

func (r *lendingRepository) GetUserBorrowingHistory(userID uint, offset, limit int) ([]domain.LendingRecord, int64, error) {
	var records []domain.LendingRecord
	var total int64
	if err := r.db.Model(&domain.LendingRecord{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Preload("Book").
		Where("user_id = ?", userID).
		Order("borrow_date DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (r *lendingRepository) GetActiveBorrowingsByUser(userID uint) ([]domain.LendingRecord, error) {
	var records []domain.LendingRecord
	if err := r.db.Preload("Book").Where("user_id = ? AND return_date IS NULL", userID).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *lendingRepository) CountUserBorrowsInPeriod(userID uint, since time.Time) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.LendingRecord{}).
		Where("user_id = ? AND borrow_date >= ?", userID, since).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
