package domain

import "time"

// User represents a system user.  Email addresses are unique and
// passwords are stored as a bcrypt hash in PasswordHash.  The
// password hash is omitted from JSON responses to avoid leaking
// sensitive information.
type User struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    Email        string    `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash string    `json:"-" gorm:"not null"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Book represents a book in the library.  ISBN numbers are unique.
type Book struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"not null"`
    Author    string    `json:"author" gorm:"not null"`
    ISBN      string    `json:"isbn" gorm:"uniqueIndex;not null"`
    Quantity  int       `json:"quantity" gorm:"not null;default:1"`
    Category  string    `json:"category" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// LendingRecord captures the act of borrowing a book.  A record is
// active until ReturnDate is set.  Books and users are referenced via
// foreign keys.
type LendingRecord struct {
    ID         uint       `json:"id" gorm:"primaryKey"`
    BookID     uint       `json:"book_id" gorm:"not null"`
    UserID     uint       `json:"user_id" gorm:"not null"`
    BorrowDate time.Time  `json:"borrow_date" gorm:"not null"`
    ReturnDate *time.Time `json:"return_date"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    Book       Book       `json:"book" gorm:"foreignKey:BookID"`
    User       User       `json:"user" gorm:"foreignKey:UserID"`
}

// TableName sets a custom table name for LendingRecord so it uses
// "lending_records" instead of the default which would be
// "lending_records" anyway.  Explicitly defining it improves clarity
// when reading migrations.
func (LendingRecord) TableName() string {
    return "lending_records"
}