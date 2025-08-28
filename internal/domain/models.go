package domain

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

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

func (LendingRecord) TableName() string { return "lending_records" }
