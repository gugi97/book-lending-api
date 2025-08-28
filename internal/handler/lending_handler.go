package handler

import (
    "book-lending-api/internal/domain"
    "book-lending-api/internal/middleware"
    "book-lending-api/internal/usecase"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

// LendingHandler wires lending use cases to HTTP routes.
type LendingHandler struct {
    lendingUseCase usecase.LendingUseCase
}

// NewLendingHandler constructs a new LendingHandler.
func NewLendingHandler(uc usecase.LendingUseCase) *LendingHandler {
    return &LendingHandler{lendingUseCase: uc}
}

// BorrowBook allows an authenticated user to borrow a book.  It
// returns the created lending record on success.  Validation errors
// result in 400 responses and conflicts (e.g. already borrowed) in 409.
func (h *LendingHandler) BorrowBook(c *gin.Context) {
    userID, exists := middleware.GetUserIDFromContext(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "Unauthorized", Message: "User not found in context"})
        return
    }
    var req domain.BorrowBookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: err.Error()})
        return
    }
    record, err := h.lendingUseCase.BorrowBook(userID, req.BookID)
    if err != nil {
        status := http.StatusInternalServerError
        switch err.Error() {
        case "book not found":
            status = http.StatusNotFound
        case "you have already borrowed this book",
            "borrowing limit exceeded: maximum 5 books per week",
            "book is not available for borrowing":
            status = http.StatusConflict
        }
        c.JSON(status, domain.ErrorResponse{Error: "Failed to borrow book", Message: err.Error()})
        return
    }
    c.JSON(http.StatusCreated, record)
}

// ReturnBook marks an existing lending record as returned.  The
// record ID is provided in the path.  Only the owner of the record
// may return it.  Returns the updated record on success.
func (h *LendingHandler) ReturnBook(c *gin.Context) {
    userID, exists := middleware.GetUserIDFromContext(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "Unauthorized", Message: "User not found in context"})
        return
    }
    idParam := c.Param("id")
    recordID, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Bad Request", Message: "Invalid lending record ID"})
        return
    }
    record, err := h.lendingUseCase.ReturnBook(userID, uint(recordID))
    if err != nil {
        status := http.StatusInternalServerError
        switch err.Error() {
        case "lending record not found":
            status = http.StatusNotFound
        case "unauthorized: this lending record does not belong to you":
            status = http.StatusForbidden
        case "book has already been returned":
            status = http.StatusConflict
        }
        c.JSON(status, domain.ErrorResponse{Error: "Failed to return book", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, record)
}

// GetBorrowingHistory returns a paginated list of a user's past borrowing
// records.  Page and limit parameters are optional and default to
// page=1 limit=10.
func (h *LendingHandler) GetBorrowingHistory(c *gin.Context) {
    userID, exists := middleware.GetUserIDFromContext(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "Unauthorized", Message: "User not found in context"})
        return
    }
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
    result, err := h.lendingUseCase.GetUserBorrowingHistory(userID, pagination.Page, pagination.Limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: "Failed to retrieve borrowing history", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

// GetActiveBorrowings lists all currently active borrowings for the
// authenticated user.  An empty slice is returned when there are
// none.
func (h *LendingHandler) GetActiveBorrowings(c *gin.Context) {
    userID, exists := middleware.GetUserIDFromContext(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "Unauthorized", Message: "User not found in context"})
        return
    }
    records, err := h.lendingUseCase.GetActiveBorrowings(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: "Failed to retrieve active borrowings", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, records)
}