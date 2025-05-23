package usecases

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"
)

type BillRefreshUsecase struct {
	repo         ports.Repository
	providerSvc  ports.ProviderAPIService
	cacheSvc     ports.CacheService
	maxRetries   int
	retryBackoff time.Duration
}

func NewBillRefreshUsecase(repo ports.Repository, providerSvc ports.ProviderAPIService, cacheSvc ports.CacheService) *BillRefreshUsecase {
	return &BillRefreshUsecase{
		repo:         repo,
		providerSvc:  providerSvc,
		cacheSvc:     cacheSvc,
		maxRetries:   3,
		retryBackoff: time.Second * 2,
	}
}

func (u *BillRefreshUsecase) RefreshBills(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user's linked accounts
	accounts, err := u.repo.GetLinkedAccountsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch linked accounts", http.StatusInternalServerError)
		return
	}

	// Process each account
	for _, account := range accounts {
		// Try to get bills from cache first
		cacheKey := "bills:" + account.ID
		if cachedBills, err := u.cacheSvc.Get(r.Context(), cacheKey); err == nil {
			var bills []*domain.Bill
			if err := json.Unmarshal([]byte(cachedBills), &bills); err == nil {
				continue // Skip if we have valid cached data
			}
		}

		// Fetch bills from provider with retry logic
		var bills []*domain.Bill
		var fetchErr error
		for i := 0; i < u.maxRetries; i++ {
			bills, fetchErr = u.providerSvc.FetchBills(r.Context(), *account)
			if fetchErr == nil {
				break
			}
			time.Sleep(u.retryBackoff * time.Duration(i+1))
		}

		if fetchErr != nil {
			// Log error but continue with other accounts
			continue
		}

		// Save bills to database
		for _, bill := range bills {
			// Generate unique ID for the bill if not already set
			if bill.ID == "" {
				bill.ID = uuid.New().String()
			}
			if err := u.repo.CreateBill(r.Context(), bill); err != nil {
				continue
			}
		}

		// Cache the bills
		if billData, err := json.Marshal(bills); err == nil {
			u.cacheSvc.Set(r.Context(), cacheKey, string(billData), time.Hour*24)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bill refresh completed"})
}
