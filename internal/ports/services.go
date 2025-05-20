package ports

import (
	"context"

	"github.com/mel-ak/onetap-challenge/internal/domain"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT token
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
}

// ProviderService defines the interface for provider-related business logic
type ProviderService interface {
	RegisterProvider(ctx context.Context, provider *domain.Provider) error
	GetProvider(ctx context.Context, id string) (*domain.Provider, error)
	ListProviders(ctx context.Context) ([]*domain.Provider, error)
	UpdateProvider(ctx context.Context, provider *domain.Provider) error
}

// AccountService defines the interface for account linking business logic
type AccountService interface {
	LinkAccount(ctx context.Context, userID string, providerID string, credentials string) (*domain.LinkedAccount, error)
	GetLinkedAccounts(ctx context.Context, userID string) ([]*domain.LinkedAccount, error)
	UnlinkAccount(ctx context.Context, accountID string) error
	RefreshAccountStatus(ctx context.Context, accountID string) error
}

// BillService defines the interface for bill-related business logic
type BillService interface {
	FetchBills(ctx context.Context, userID string) (*domain.BillSummary, error)
	FetchBillsByProvider(ctx context.Context, userID string, providerID string) ([]*domain.Bill, error)
	RefreshBills(ctx context.Context, userID string) error
	GetBillSummary(ctx context.Context, userID string) (*domain.BillSummary, error)
}
