package domain

// The DTOs below define the payloads used by HTTP handlers and
// responses.  Having a dedicated set of request/response types keeps
// the API surface explicit and separate from the underlying models.

// LoginRequest represents the login request payload.
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents the user registration request payload.
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse represents the authentication response returned after
// successful login or registration.  It contains a JWT token and the
// authenticated user's information.
type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

// CreateBookRequest represents the request body used when creating a
// new book.  All fields are required.
type CreateBookRequest struct {
    Title    string `json:"title" binding:"required"`
    Author   string `json:"author" binding:"required"`
    ISBN     string `json:"isbn" binding:"required"`
    Quantity int    `json:"quantity" binding:"required,min=1"`
    Category string `json:"category" binding:"required"`
}

// UpdateBookRequest represents a request to update one or more fields
// on an existing book.  Pointers are used so that zero values can be
// distinguished from fields that were not provided.
type UpdateBookRequest struct {
    Title    *string `json:"title"`
    Author   *string `json:"author"`
    ISBN     *string `json:"isbn"`
    Quantity *int    `json:"quantity"`
    Category *string `json:"category"`
}

// BorrowBookRequest represents the request to borrow a book by id.
type BorrowBookRequest struct {
    BookID uint `json:"book_id" binding:"required"`
}

// PaginationRequest contains pagination parameters used when listing
// resources.  Sensible defaults are supplied by the handler when
// fields are omitted.
type PaginationRequest struct {
    Page  int `form:"page,default=1" binding:"min=1"`
    Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

// PaginatedResponse wraps a slice of data together with pagination
// metadata.  It is returned by list endpoints to make it easy for
// clients to navigate through large result sets.
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    Total      int64       `json:"total"`
    TotalPages int         `json:"total_pages"`
}

// ErrorResponse describes an error returned by the API.  It includes
// both a short error code and a human‑readable message.
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
}

// SuccessResponse describes a simple success message returned when
// deleting a resource.
type SuccessResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}