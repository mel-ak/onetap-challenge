package usecases

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/adapters/providers"
	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"
)

type billService struct {
	repo        ports.Repository
	rateLimiter *RateLimiter
	providers   map[string]providers.Provider
}

// NewBillService creates a new instance of the bill service
func NewBillService(repo ports.Repository) ports.BillService {
	// Initialize providers map with mock provider
	providersMap := make(map[string]providers.Provider)
	providersMap["mock-provider"] = providers.NewMockProviderAdapter("http://localhost:8083")

	return &billService{
		repo:        repo,
		rateLimiter: NewRateLimiter(100, time.Minute), // 100 requests per minute
		providers:   providersMap,
	}
}

// FetchBills retrieves all bills for a user
func (s *billService) FetchBills(ctx context.Context, userID string) (*domain.BillSummary, error) {
	// Get all linked accounts for the user
	accounts, err := s.repo.GetLinkedAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return &domain.BillSummary{
			BillCount: 0,
			Bills:     nil,
			TotalDue:  0,
		}, nil
	}

	var wg sync.WaitGroup
	billChan := make(chan []*domain.Bill, len(accounts))
	errChan := make(chan error, len(accounts))

	// Fetch bills for each account concurrently
	for _, account := range accounts {
		wg.Add(1)
		go func(acc *domain.LinkedAccount) {
			defer wg.Done()

			// Apply rate limiting
			if err := s.rateLimiter.Wait(ctx); err != nil {
				errChan <- err
				return
			}

			bills, err := s.fetchBillsForAccount(ctx, acc)
			if err != nil {
				errChan <- err
				return
			}
			billChan <- bills
		}(account)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(billChan)
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	// Collect all bills
	var allBills []*domain.Bill
	for bills := range billChan {
		allBills = append(allBills, bills...)
	}

	// Calculate summary
	summary := &domain.BillSummary{
		BillCount: len(allBills),
		Bills:     make([]domain.Bill, len(allBills)),
	}

	var totalDue float64
	for i, bill := range allBills {
		summary.Bills[i] = *bill
		if bill.Status == "unpaid" || bill.Status == "overdue" {
			totalDue += bill.Amount
		}
	}
	summary.TotalDue = totalDue

	return summary, nil
}

// FetchBillsByProvider retrieves bills for a specific provider
func (s *billService) FetchBillsByProvider(ctx context.Context, userID string, providerID string) ([]*domain.Bill, error) {
	accounts, err := s.repo.GetLinkedAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var providerBills []*domain.Bill
	for _, account := range accounts {
		if account.ProviderID == providerID {
			bills, err := s.fetchBillsForAccount(ctx, account)
			if err != nil {
				return nil, err
			}
			providerBills = append(providerBills, bills...)
		}
	}

	return providerBills, nil
}

// RefreshBills triggers a refresh of all bills for a user
func (s *billService) RefreshBills(ctx context.Context, userID string) error {
	accounts, err := s.repo.GetLinkedAccountsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(accounts))

	for _, account := range accounts {
		wg.Add(1)
		go func(acc *domain.LinkedAccount) {
			defer wg.Done()

			if err := s.rateLimiter.Wait(ctx); err != nil {
				errChan <- err
				return
			}

			bills, err := s.fetchBillsForAccount(ctx, acc)
			if err != nil {
				errChan <- err
				return
			}

			// Update bills in database
			for _, bill := range bills {
				if err := s.repo.UpdateBill(ctx, bill); err != nil {
					errChan <- err
					return
				}
			}
		}(account)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

// GetBillSummary retrieves a summary of bills for a user
func (s *billService) GetBillSummary(ctx context.Context, userID string) (*domain.BillSummary, error) {
	return s.repo.GetBillSummaryByUserID(ctx, userID)
}

// fetchBillsForAccount is a helper method to fetch bills for a specific account
func (s *billService) fetchBillsForAccount(ctx context.Context, account *domain.LinkedAccount) ([]*domain.Bill, error) {
	provider, exists := s.providers[account.ProviderID]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", account.ProviderID)
	}

	return provider.FetchBills(ctx, account.ID)
}

// StartPeriodicUpdates starts the background job for periodic bill updates
func (s *billService) StartPeriodicUpdates(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Update daily
	go func() {
		for {
			select {
			case <-ticker.C:
				users, err := s.repo.ListUsers(ctx)
				if err != nil {
					log.Printf("Error fetching users for periodic update: %v", err)
					continue
				}

				for _, user := range users {
					if err := s.RefreshBills(ctx, user.ID); err != nil {
						log.Printf("Error refreshing bills for user %s: %v", user.ID, err)
					}
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
