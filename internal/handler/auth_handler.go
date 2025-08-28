package handler

import (
    "book-lending-api/internal/domain"
    "book-lending-api/internal/usecase"
    "book-lending-api/pkg"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

// AuthHandler wires authentication use cases to HTTP requests.
type AuthHandler struct {
    authUseCase usecase.AuthUseCase
    jwtUtil     *pkg.JWTUtil
}

// NewAuthHandler constructs a new AuthHandler.
func NewAuthHandler(authUseCase usecase.AuthUseCase, jwtUtil *pkg.JWTUtil) *AuthHandler {
    return &AuthHandler{authUseCase: authUseCase, jwtUtil: jwtUtil}
}

// Register handles user registration.  On success it returns a
// newly minted JWT token and user object.  Validation errors yield
// 400 responses and conflicts yield 409.
func (h *AuthHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.ErrorResponse{
            Error:   "Bad Request",
            Message: err.Error(),
        })
        return
    }
    user, err := h.authUseCase.Register(req)
    if err != nil {
        status := http.StatusInternalServerError
        if err.Error() == "user with this email already exists" {
            status = http.StatusConflict
        }
        c.JSON(status, domain.ErrorResponse{
            Error:   "Registration failed",
            Message: err.Error(),
        })
        return
    }
    token, err := h.jwtUtil.GenerateToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
            Error:   "Token generation failed",
            Message: err.Error(),
        })
        return
    }
    c.JSON(http.StatusCreated, domain.AuthResponse{Token: token, User: *user})
}

// Login handles user authentication.  On success it returns a new
// token and user object.  Invalid credentials return 401.
func (h *AuthHandler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.ErrorResponse{
            Error:   "Bad Request",
            Message: err.Error(),
        })
        return
    }
    user, err := h.authUseCase.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
            Error:   "Login failed",
            Message: err.Error(),
        })
        return
    }
    token, err := h.jwtUtil.GenerateToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
            Error:   "Token generation failed",
            Message: err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, domain.AuthResponse{Token: token, User: *user})
}