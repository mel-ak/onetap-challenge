package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
)

type MockProviderAdapter struct {
	baseURL    string
	httpClient *http.Client
}

func NewMockProviderAdapter(baseURL string) *MockProviderAdapter {
	return &MockProviderAdapter{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *MockProviderAdapter) FetchBills(ctx context.Context, accountID string) ([]*domain.Bill, error) {
	url := fmt.Sprintf("%s/bills", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bills: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var mockBills []struct {
		ID          string    `json:"id"`
		Provider    string    `json:"provider"`
		Amount      float64   `json:"amount"`
		DueDate     time.Time `json:"due_date"`
		Status      string    `json:"status"`
		Description string    `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mockBills); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	bills := make([]*domain.Bill, len(mockBills))
	for i, mockBill := range mockBills {
		bills[i] = &domain.Bill{
			ID:              mockBill.ID,
			LinkedAccountID: accountID,
			ProviderID:      mockBill.Provider,
			Amount:          mockBill.Amount,
			DueDate:         mockBill.DueDate,
			Status:          mockBill.Status,
			BillDate:        time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
	}

	return bills, nil
}

func (a *MockProviderAdapter) ValidateCredentials(ctx context.Context, credentials string) error {
	// For mock provider, we'll always return success
	return nil
}

func (a *MockProviderAdapter) GetProviderInfo() *domain.Provider {
	return &domain.Provider{
		ID:          "mock-provider",
		Name:        "Mock Provider",
		APIEndpoint: a.baseURL,
		AuthType:    "none",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
