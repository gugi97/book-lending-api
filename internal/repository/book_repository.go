package repository

import (
	"book-lending-api/internal/domain"

	"gorm.io/gorm"
)

// BookRepository abstracts persistence operations for books.
type BookRepository interface {
	Create(book *domain.Book) error
	GetByID(id uint) (*domain.Book, error)
	GetByISBN(isbn string) (*domain.Book, error)
	Update(book *domain.Book) error
	Delete(id uint) error
	List(offset, limit int) ([]domain.Book, int64, error)
	GetAvailableQuantity(bookID uint) (int, error)
	UpdateQuantity(bookID uint, quantity int) error
}

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository returns a new BookRepository using the provided
// gorm DB.
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *domain.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) GetByID(id uint) (*domain.Book, error) {
	var book domain.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetByISBN(isbn string) (*domain.Book, error) {
	var book domain.Book
	if err := r.db.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Book{}, id).Error
}

// List returns a slice of books along with the total count.  Offset
// and limit control pagination.
func (r *bookRepository) List(offset, limit int) ([]domain.Book, int64, error) {
	var books []domain.Book
	var total int64
	if err := r.db.Model(&domain.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, err
	}
	return books, total, nil
}

// GetAvailableQuantity calculates the number of books available for
// borrowing by subtracting the number of active lending records from
// the total quantity.
func (r *bookRepository) GetAvailableQuantity(bookID uint) (int, error) {
	var book domain.Book
	if err := r.db.Select("quantity").First(&book, bookID).Error; err != nil {
		return 0, err
	}
	var borrowedCount int64
	if err := r.db.Model(&domain.LendingRecord{}).
		Where("book_id = ? AND return_date IS NULL", bookID).
		Count(&borrowedCount).Error; err != nil {
		return 0, err
	}
	return book.Quantity - int(borrowedCount), nil
}

// UpdateQuantity allows adjusting the total quantity for a book.
func (r *bookRepository) UpdateQuantity(bookID uint, quantity int) error {
	return r.db.Model(&domain.Book{}).Where("id = ?", bookID).
		Update("quantity", quantity).Error
}
