package usecase

import (
	"book-lending-api/internal/domain"
	"book-lending-api/internal/repository"
	"errors"
	"math"
)

// BookUseCase defines business logic operations for books.
type BookUseCase interface {
	CreateBook(req domain.CreateBookRequest) (*domain.Book, error)
	GetBookByID(id uint) (*domain.Book, error)
	UpdateBook(id uint, req domain.UpdateBookRequest) (*domain.Book, error)
	DeleteBook(id uint) error
	ListBooks(page, limit int) (*domain.PaginatedResponse, error)
}

type bookUseCase struct {
	bookRepo repository.BookRepository
}

// NewBookUseCase constructs a new book use case.
func NewBookUseCase(bookRepo repository.BookRepository) BookUseCase {
	return &bookUseCase{bookRepo: bookRepo}
}

func (uc *bookUseCase) CreateBook(req domain.CreateBookRequest) (*domain.Book, error) {
	// check for duplicate ISBN
	if existing, _ := uc.bookRepo.GetByISBN(req.ISBN); existing != nil {
		return nil, errors.New("book with this ISBN already exists")
	}
	book := &domain.Book{
		Title:    req.Title,
		Author:   req.Author,
		ISBN:     req.ISBN,
		Quantity: req.Quantity,
		Category: req.Category,
	}
	if err := uc.bookRepo.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (uc *bookUseCase) GetBookByID(id uint) (*domain.Book, error) {
	book, err := uc.bookRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (uc *bookUseCase) UpdateBook(id uint, req domain.UpdateBookRequest) (*domain.Book, error) {
	book, err := uc.bookRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("book not found")
	}
	// handle ISBN change
	if req.ISBN != nil && *req.ISBN != book.ISBN {
		if existing, _ := uc.bookRepo.GetByISBN(*req.ISBN); existing != nil {
			return nil, errors.New("book with this ISBN already exists")
		}
		book.ISBN = *req.ISBN
	}
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Quantity != nil {
		book.Quantity = *req.Quantity
	}
	if req.Category != nil {
		book.Category = *req.Category
	}
	if err := uc.bookRepo.Update(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (uc *bookUseCase) DeleteBook(id uint) error {
	if _, err := uc.bookRepo.GetByID(id); err != nil {
		return errors.New("book not found")
	}
	return uc.bookRepo.Delete(id)
}

func (uc *bookUseCase) ListBooks(page, limit int) (*domain.PaginatedResponse, error) {
	offset := (page - 1) * limit
	books, total, err := uc.bookRepo.List(offset, limit)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &domain.PaginatedResponse{
		Data:       books,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
