// Unit tests for BookHandler endpoints.
package handler

import (
	"book-lending-api/internal/domain"
	"book-lending-api/internal/usecase"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/gin-gonic/gin"
)

// mockBookUseCase implements usecase.BookUseCase for handler tests.
type mockBookUseCase struct{}

const notImpl = "not implemented"

func (m *mockBookUseCase) CreateBook(req domain.CreateBookRequest) (*domain.Book, error) {
	return nil, errors.New(notImpl)
}
func (m *mockBookUseCase) GetBookByID(id uint) (*domain.Book, error) {
	if id == 1 {
		return &domain.Book{ID: 1, Title: "Dune", Author: "Frank Herbert", ISBN: "9780441172719", Quantity: 3, Category: "Sci-Fi"}, nil
	}
	return nil, errors.New("book not found")
}
func (m *mockBookUseCase) UpdateBook(id uint, req domain.UpdateBookRequest) (*domain.Book, error) {
	return nil, errors.New(notImpl)
}
func (m *mockBookUseCase) DeleteBook(id uint) error { return errors.New(notImpl) }
func (m *mockBookUseCase) ListBooks(page, limit int) (*domain.PaginatedResponse, error) {
	return nil, errors.New(notImpl)
}

// Ensure mock matches interface
var _ usecase.BookUseCase = (*mockBookUseCase)(nil)

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

func TestBookHandlerGetBookSuccess(t *testing.T) {
	r := setupGin()
	h := NewBookHandler(&mockBookUseCase{})
	r.GET("/books/:id", h.GetBook)

	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !containsAll(body, []string{"Dune", "Frank Herbert", "9780441172719"}) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestBookHandlerGetBookInvalidID(t *testing.T) {
	r := setupGin()
	h := NewBookHandler(&mockBookUseCase{})
	r.GET("/books/:id", h.GetBook)

	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

// containsAll is a tiny helper to assert substrings in the response.
func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
