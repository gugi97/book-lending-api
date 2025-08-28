// Unit tests for BookUseCase
package usecase

import (
	"book-lending-api/internal/domain"
	"book-lending-api/internal/repository"
	"errors"
	"testing"
)

// mockBookRepo for usecase tests
type mockBookRepo struct{ existingByISBN map[string]*domain.Book }

func (m *mockBookRepo) Create(book *domain.Book) error        { return nil }
func (m *mockBookRepo) GetByID(id uint) (*domain.Book, error) { return &domain.Book{ID: id}, nil }
func (m *mockBookRepo) GetByISBN(isbn string) (*domain.Book, error) {
	if b := m.existingByISBN[isbn]; b != nil {
		return b, nil
	}
	return nil, errors.New("not found")
}
func (m *mockBookRepo) Update(book *domain.Book) error                       { return nil }
func (m *mockBookRepo) Delete(id uint) error                                 { return nil }
func (m *mockBookRepo) List(offset, limit int) ([]domain.Book, int64, error) { return nil, 0, nil }
func (m *mockBookRepo) GetAvailableQuantity(bookID uint) (int, error)        { return 1, nil }
func (m *mockBookRepo) UpdateQuantity(bookID uint, quantity int) error       { return nil }

var _ repository.BookRepository = (*mockBookRepo)(nil)

func TestBookUseCaseCreateBookDuplicateISBN(t *testing.T) {
	repo := &mockBookRepo{existingByISBN: map[string]*domain.Book{"123": {ID: 1, ISBN: "123"}}}
	uc := NewBookUseCase(repo)

	_, err := uc.CreateBook(domain.CreateBookRequest{Title: "T", Author: "A", ISBN: "123", Quantity: 1, Category: "C"})
	if err == nil || err.Error() != "book with this ISBN already exists" {
		t.Fatalf("expected duplicate ISBN error, got %v", err)
	}
}
