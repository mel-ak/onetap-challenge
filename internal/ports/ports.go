package ports

import (
	"context"

	"github.com/mel-ak/onetap-challenge/internal/domain"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	SaveUser(ctx context.Context, user domain.User) (string, error)
	GetUserByID(ctx context.Context, userID string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, userID string) (bool, error)
}

// AccountRepository defines the interface for account persistence
type AccountRepository interface {
	SaveAccount(ctx context.Context, account domain.LinkedAccount) (string, error)
	GetAccountsByUserID(ctx context.Context, userID string) ([]domain.LinkedAccount, error)
	DeleteAccount(ctx context.Context, accountID string) (bool, error)
}

// BillRepository defines the interface for bill persistence
type BillRepository interface {
	SaveBill(ctx context.Context, bill domain.Bill) error
}

// ProviderAPIService defines the interface for third-party provider APIs
type ProviderAPIService interface {
	FetchBills(ctx context.Context, account domain.LinkedAccount) ([]domain.Bill, error)
}

// CacheService defines the interface for caching and rate limiting
type CacheService interface {
	RateLimit(ctx context.Context, key string, limit int, window int64) error
	GetBills(ctx context.Context, key string) ([]domain.Bill, error)
	CacheBills(ctx context.Context, key string, bills []domain.Bill, ttl int64) error
}
