package domain

import "time"

// User represents a registered user
type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// Account represents a linked utility account
type Account struct {
	ID          int
	UserID      int
	Provider    string
	Credentials string
	CreatedAt   time.Time
}

// Bill represents a utility bill
type Bill struct {
	ID        int
	AccountID int
	Provider  string
	Amount    float64
	DueDate   time.Time
	Status    string
	CreatedAt time.Time
}
