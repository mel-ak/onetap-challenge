package ports

import (
	"context"

	"github.com/mel-ak/onetap-challenge/internal/domain"
)

// Repository defines the interface for all repository operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) (bool, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)

	// Provider operations
	CreateProvider(ctx context.Context, provider *domain.Provider) error
	GetProviderByID(ctx context.Context, id string) (*domain.Provider, error)
	GetProviderByName(ctx context.Context, name string) (*domain.Provider, error)
	ListProviders(ctx context.Context) ([]*domain.Provider, error)
	UpdateProvider(ctx context.Context, provider *domain.Provider) error
	DeleteProvider(ctx context.Context, id string) error

	// LinkedAccount operations
	CreateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error
	GetLinkedAccountByID(ctx context.Context, id string) (*domain.LinkedAccount, error)
	GetLinkedAccountsByUserID(ctx context.Context, userID string) ([]*domain.LinkedAccount, error)
	GetLinkedAccountsByProviderID(ctx context.Context, providerID string) ([]*domain.LinkedAccount, error)
	UpdateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error
	DeleteLinkedAccount(ctx context.Context, id string) error

	// Bill operations
	CreateBill(ctx context.Context, bill *domain.Bill) error
	GetBillByID(ctx context.Context, id string) (*domain.Bill, error)
	GetBillsByLinkedAccountID(ctx context.Context, linkedAccountID string) ([]*domain.Bill, error)
	GetBillsByUserID(ctx context.Context, userID string) ([]*domain.Bill, error)
	GetBillSummaryByUserID(ctx context.Context, userID string) (*domain.BillSummary, error)
	UpdateBill(ctx context.Context, bill *domain.Bill) error
	DeleteBill(ctx context.Context, id string) error
}

// AccountRepository defines the interface for account-related database operations
// ... existing code ...
