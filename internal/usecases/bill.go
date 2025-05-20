package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid or missing user_id", http.StatusBadRequest)
		return
	}

	// Fetch accounts
	accounts, err := u.repo.GetAccountsByUserID(context.Background(), userIDStr)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Fetch bills concurrently
	billsChan := make(chan []domain.Bill, len(accounts))
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
	var allBills []domain.Bill
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

	resp := map[string]interface{}{
		"bills":      allBills,
		"total_due":  totalDue,
		"bill_count": len(allBills),
	}
	json.NewEncoder(w).Encode(resp)
}

func (u *BillUsecase) fetchBillsWithRetry(ctx context.Context, acc domain.LinkedAccount) ([]domain.Bill, error) {
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
	var bills []domain.Bill
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
