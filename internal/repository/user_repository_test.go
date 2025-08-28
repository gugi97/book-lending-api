// Unit tests for UserRepository using sqlite in-memory
package repository

import (
	"book-lending-api/internal/domain"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB prepares an in-memory sqlite database and migrates schema
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUserRepositoryCRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Email: "alice@example.com", PasswordHash: "hash"}
	if err := repo.Create(user); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if user.ID == 0 {
		t.Fatalf("expected ID to be set after create")
	}

	gotByEmail, err := repo.GetByEmail("alice@example.com")
	if err != nil || gotByEmail == nil || gotByEmail.Email != user.Email {
		t.Fatalf("get by email failed: user=%v err=%v", gotByEmail, err)
	}

	gotByID, err := repo.GetByID(user.ID)
	if err != nil || gotByID == nil || gotByID.Email != user.Email {
		t.Fatalf("get by id failed: user=%v err=%v", gotByID, err)
	}
}
