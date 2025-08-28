package usecase

import (
	"book-lending-api/internal/domain"
	"book-lending-api/internal/repository"
	"errors"
	"math"
	"time"
)

// LendingUseCase defines the operations for borrowing and returning books.
type LendingUseCase interface {
	BorrowBook(userID, bookID uint) (*domain.LendingRecord, error)
	ReturnBook(userID, recordID uint) (*domain.LendingRecord, error)
	GetUserBorrowingHistory(userID uint, page, limit int) (*domain.PaginatedResponse, error)
	GetActiveBorrowings(userID uint) ([]domain.LendingRecord, error)
}

type lendingUseCase struct {
	lendingRepo repository.LendingRepository
	bookRepo    repository.BookRepository
}

// NewLendingUseCase constructs a new lending use case.
func NewLendingUseCase(lendingRepo repository.LendingRepository, bookRepo repository.BookRepository) LendingUseCase {
	return &lendingUseCase{
		lendingRepo: lendingRepo,
		bookRepo:    bookRepo,
	}
}

func (uc *lendingUseCase) BorrowBook(userID, bookID uint) (*domain.LendingRecord, error) {
	// verify book exists
	if _, err := uc.bookRepo.GetByID(bookID); err != nil {
		return nil, errors.New("book not found")
	}
	// ensure user hasn't borrowed this book already
	if rec, _ := uc.lendingRepo.GetActiveByUserAndBook(userID, bookID); rec != nil {
		return nil, errors.New("you have already borrowed this book")
	}
	// enforce weekly borrow limit
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	count, err := uc.lendingRepo.CountUserBorrowsInPeriod(userID, sevenDaysAgo)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, errors.New("borrowing limit exceeded: maximum 5 books per week")
	}
	// ensure availability
	available, err := uc.bookRepo.GetAvailableQuantity(bookID)
	if err != nil {
		return nil, err
	}
	if available <= 0 {
		return nil, errors.New("book is not available for borrowing")
	}
	record := &domain.LendingRecord{
		BookID:     bookID,
		UserID:     userID,
		BorrowDate: time.Now(),
	}
	if err := uc.lendingRepo.Create(record); err != nil {
		return nil, err
	}
	return uc.lendingRepo.GetByID(record.ID)
}

func (uc *lendingUseCase) ReturnBook(userID, recordID uint) (*domain.LendingRecord, error) {
	record, err := uc.lendingRepo.GetByID(recordID)
	if err != nil {
		return nil, errors.New("lending record not found")
	}
	if record.UserID != userID {
		return nil, errors.New("unauthorized: this lending record does not belong to you")
	}
	if record.ReturnDate != nil {
		return nil, errors.New("book has already been returned")
	}
	now := time.Now()
	record.ReturnDate = &now
	if err := uc.lendingRepo.Update(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (uc *lendingUseCase) GetUserBorrowingHistory(userID uint, page, limit int) (*domain.PaginatedResponse, error) {
	offset := (page - 1) * limit
	records, total, err := uc.lendingRepo.GetUserBorrowingHistory(userID, offset, limit)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &domain.PaginatedResponse{
		Data:       records,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (uc *lendingUseCase) GetActiveBorrowings(userID uint) ([]domain.LendingRecord, error) {
	return uc.lendingRepo.GetActiveBorrowingsByUser(userID)
}
