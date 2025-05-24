package providers

import (
	"context"

	"github.com/mel-ak/onetap-challenge/internal/domain"
)

// Provider defines the interface that all provider adapters must implement
type Provider interface {
	// FetchBills retrieves bills for a given account
	FetchBills(ctx context.Context, accountID string) ([]*domain.Bill, error)

	// ValidateCredentials validates the provided credentials
	ValidateCredentials(ctx context.Context, credentials string) error

	// GetProviderInfo returns information about the provider
	GetProviderInfo() *domain.Provider
}
