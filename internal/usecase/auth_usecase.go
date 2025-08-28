package usecase

import (
	"book-lending-api/internal/domain"
	"book-lending-api/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase defines the operations available for authentication.
type AuthUseCase interface {
	Register(req domain.RegisterRequest) (*domain.User, error)
	Login(req domain.LoginRequest) (*domain.User, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
}

// NewAuthUseCase constructs a new authentication use case.
func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepo: userRepo}
}

// Register registers a new user.  It hashes the password using bcrypt
// and returns an error if the email is already taken.
func (uc *authUseCase) Register(req domain.RegisterRequest) (*domain.User, error) {
	if existing, _ := uc.userRepo.GetByEmail(req.Email); existing != nil {
		return nil, errors.New("user with this email already exists")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashed),
	}
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login authenticates a user by checking the provided credentials.
func (uc *authUseCase) Login(req domain.LoginRequest) (*domain.User, error) {
	user, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
