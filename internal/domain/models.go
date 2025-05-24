package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is not exposed in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Provider represents a utility provider
type Provider struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	APIEndpoint string    `json:"api_endpoint"`
	AuthType    string    `json:"auth_type"` // e.g., "oauth2", "api_key", "basic"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LinkedAccount represents a user's linked utility account
type LinkedAccount struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProviderID  string    `json:"provider_id"`
	AccountID   string    `json:"account_id"` // Provider's account ID
	Credentials string    `json:"-"`          // Encrypted credentials
	Status      string    `json:"status"`     // active, inactive, error
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Bill represents a utility bill
type Bill struct {
	ID              string    `json:"id"`
	LinkedAccountID string    `json:"linked_account_id"`
	ProviderID      string    `json:"provider_id"`
	Amount          float64   `json:"amount"`
	DueDate         time.Time `json:"due_date"`
	Status          string    `json:"status"` // paid, unpaid, overdue
	BillDate        time.Time `json:"bill_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BillSummary represents aggregated bill information
type BillSummary struct {
	BillCount int     `json:"bill_count"`
	Bills     []Bill  `json:"bills"`
	TotalDue  float64 `json:"total_due"`
}
