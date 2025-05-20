package usecases

import (
	"context"
	"sync"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"
)

type billService struct {
	repo        ports.Repository
	rateLimiter *RateLimiter
}

// NewBillService creates a new instance of the bill service
func NewBillService(repo ports.Repository) ports.BillService {
	return &billService{
		repo:        repo,
		rateLimiter: NewRateLimiter(100, time.Minute), // 100 requests per minute
	}
}

// FetchBills retrieves all bills for a user
func (s *billService) FetchBills(ctx context.Context, userID string) (*domain.BillSummary, error) {
	// Get all linked accounts for the user
	accounts, err := s.repo.GetLinkedAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
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
		Bills: make([]domain.Bill, len(allBills)),
	}

	for i, bill := range allBills {
		summary.Bills[i] = *bill
		summary.TotalAmount += bill.Amount
		if bill.Status == "unpaid" {
			summary.DueBills++
		} else if bill.Status == "overdue" {
			summary.OverdueBills++
		}
	}

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
	// In a real implementation, this would call the provider's API
	// For now, we'll return mock data
	return []*domain.Bill{
		{
			ID:              "1",
			LinkedAccountID: account.ID,
			ProviderID:      account.ProviderID,
			Amount:          100.00,
			DueDate:         time.Now().AddDate(0, 0, 7),
			Status:          "unpaid",
			BillDate:        time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}, nil
}
