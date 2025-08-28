package domain

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	Author   string `json:"author" binding:"required"`
	ISBN     string `json:"isbn" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Category string `json:"category" binding:"required"`
}

type UpdateBookRequest struct {
	Title    *string `json:"title"`
	Author   *string `json:"author"`
	ISBN     *string `json:"isbn"`
	Quantity *int    `json:"quantity"`
	Category *string `json:"category"`
}

type BorrowBookRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

type PaginationRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
