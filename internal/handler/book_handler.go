package handler

import (
	"book-lending-api/internal/domain"
	"book-lending-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookHandler exposes book use cases over HTTP.
type BookHandler struct {
	bookUseCase usecase.BookUseCase
}

// NewBookHandler constructs a new BookHandler.
func NewBookHandler(uc usecase.BookUseCase) *BookHandler {
	return &BookHandler{bookUseCase: uc}
}

// CreateBook handles creating a new book.  The endpoint is
// authenticated via middleware upstream.  Duplicate ISBNs return 409.
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req domain.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		return
	}
	book, err := h.bookUseCase.CreateBook(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "book with this ISBN already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, domain.ErrorResponse{Error: "Failed to create book", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// GetBook retrieves a single book by id.  If the id is invalid or the
// book is not found appropriate HTTP statuses are returned.
func (h *BookHandler) GetBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: "Invalid book ID"})
		return
	}
	book, err := h.bookUseCase.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Error: "Not Found", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// UpdateBook updates an existing book by id.  Conflicts and not found
// cases return 409 and 404 respectively.
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: "Invalid book ID"})
		return
	}
	var req domain.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		return
	}
	book, err := h.bookUseCase.UpdateBook(uint(id), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "book not found":
			status = http.StatusNotFound
		case "book with this ISBN already exists":
			status = http.StatusConflict
		}
		c.JSON(status, domain.ErrorResponse{Error: "Failed to update book", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book.  Not found errors return 404.
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: "Invalid book ID"})
		return
	}
	if err := h.bookUseCase.DeleteBook(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "book not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, domain.ErrorResponse{Error: "Failed to delete book", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Book deleted successfully"})
}

// ListBooks lists books with pagination.  Defaults to page=1 and
// limit=10 when parameters are omitted.  Invalid parameters return a
// 400 response.
func (h *BookHandler) ListBooks(c *gin.Context) {
	var pagination domain.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		return
	}
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}
	result, err := h.bookUseCase.ListBooks(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: "Failed to retrieve books", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
