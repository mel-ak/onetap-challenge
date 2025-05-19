package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
)

// HTTPProvider simulates a third-party provider API
type HTTPProvider struct{}

// NewHTTPProvider creates a new provider
func NewHTTPProvider() *HTTPProvider {
	return &HTTPProvider{}
}

// FetchBills fetches bills from a provider
func (p *HTTPProvider) FetchBills(ctx context.Context, account domain.Account) ([]domain.Bill, error) {
	// Simulate slow or failing API
	if rand.Float32() < 0.2 {
		return nil, fmt.Errorf("provider API timeout")
	}

	return []domain.Bill{
		{
			AccountID: account.ID,
			Provider:  account.Provider,
			Amount:    rand.Float64() * 100,
			DueDate:   time.Now().AddDate(0, 0, rand.Intn(30)),
			Status:    []string{"paid", "unpaid", "overdue"}[rand.Intn(3)],
		},
	}, nil
}
