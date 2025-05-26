package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"
)

// BillUsecase handles bill-related business logic
type BillUsecase struct {
	repo     ports.AccountRepository
	provider ports.ProviderAPIService
	cache    ports.CacheService
}

// NewBillUsecase creates a new bill use case
func NewBillUsecase(repo ports.AccountRepository, provider ports.ProviderAPIService, cache ports.CacheService) *BillUsecase {
	return &BillUsecase{repo: repo, provider: provider, cache: cache}
}

// FetchBills handles GET /bills
func (u *BillUsecase) FetchBills(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	// Fetch accounts
	accounts, err := u.repo.GetAccountsByUserID(context.Background(), userIDStr)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Fetch bills concurrently
	billsChan := make(chan []*domain.Bill, len(accounts))
	var wg sync.WaitGroup

	for _, acc := range accounts {
		wg.Add(1)
		go func(acc domain.LinkedAccount) {
			defer wg.Done()
			bills, err := u.fetchBillsWithRetry(context.Background(), acc)
			if err != nil {
				return
			}
			billsChan <- bills
		}(acc)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(billsChan)
	}()

	// Collect bills
	var allBills []*domain.Bill
	for bills := range billsChan {
		allBills = append(allBills, bills...)
	}

	// Calculate total amount due
	var totalDue float64
	for _, bill := range allBills {
		// if bill.Status == "unpaid" || bill.Status == "overdue" {
		totalDue += bill.Amount
		// }
	}

	resp := map[string]interface{}{
		"bills":      allBills,
		"total_due":  totalDue,
		"bill_count": len(allBills),
	}
	json.NewEncoder(w).Encode(resp)
}

func (u *BillUsecase) fetchBillsWithRetry(ctx context.Context, acc domain.LinkedAccount) ([]*domain.Bill, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", acc.ProviderID, acc.ID)
	cacheKey := fmt.Sprintf("bills:%s", acc.ID)

	// Rate limiting
	if err := u.cache.RateLimit(ctx, key, 5, int64(time.Minute.Seconds())); err != nil {
		return nil, err
	}

	// Check cache
	if bills, err := u.cache.GetBills(ctx, cacheKey); err == nil && len(bills) > 0 {
		return bills, nil
	}

	// Fetch from provider with retries
	var bills []*domain.Bill
	var err error
	for i := 0; i < 3; i++ {
		bills, err = u.provider.FetchBills(ctx, acc)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(1<<i))
	}
	if err != nil {
		return nil, err
	}

	// Cache and save bills
	u.cache.CacheBills(ctx, cacheKey, bills, int64(time.Hour.Seconds()))
	return bills, nil
}

// FetchBillsByProvider handles GET /providers/{provider_id}/bills
func (u *BillUsecase) FetchBillsByProvider(w http.ResponseWriter, r *http.Request) {
	// Get provider ID from URL parameters
	vars := mux.Vars(r)
	providerID := vars["provider_id"]

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all accounts for the user
	accounts, err := u.repo.GetAccountsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Filter accounts for the specific provider
	var providerAccounts []domain.LinkedAccount
	for _, acc := range accounts {
		if acc.ProviderID == providerID {
			providerAccounts = append(providerAccounts, acc)
		}
	}

	if len(providerAccounts) == 0 {
		// Return empty response if no accounts found for the provider
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"bills":      []interface{}{},
			"total_due":  0,
			"bill_count": 0,
		})
		return
	}

	// Fetch bills concurrently for all accounts of this provider
	billsChan := make(chan []*domain.Bill, len(providerAccounts))
	var wg sync.WaitGroup

	for _, acc := range providerAccounts {
		wg.Add(1)
		go func(acc domain.LinkedAccount) {
			defer wg.Done()
			bills, err := u.fetchBillsWithRetry(r.Context(), acc)
			if err != nil {
				return
			}
			billsChan <- bills
		}(acc)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(billsChan)
	}()

	// Collect bills
	var allBills []*domain.Bill
	for bills := range billsChan {
		allBills = append(allBills, bills...)
	}

	// Calculate total amount due
	var totalDue float64
	for _, bill := range allBills {
		if bill.Status == "unpaid" || bill.Status == "overdue" {
			totalDue += bill.Amount
		}
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bills":      allBills,
		"total_due":  totalDue,
		"bill_count": len(allBills),
	})
}
