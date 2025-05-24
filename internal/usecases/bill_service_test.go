package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// User operations
func (m *MockRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

// Provider operations
func (m *MockRepository) CreateProvider(ctx context.Context, provider *domain.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockRepository) GetProviderByID(ctx context.Context, id string) (*domain.Provider, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Provider), args.Error(1)
}

func (m *MockRepository) GetProviderByName(ctx context.Context, name string) (*domain.Provider, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*domain.Provider), args.Error(1)
}

func (m *MockRepository) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Provider), args.Error(1)
}

func (m *MockRepository) UpdateProvider(ctx context.Context, provider *domain.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockRepository) DeleteProvider(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// LinkedAccount operations
func (m *MockRepository) CreateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockRepository) GetLinkedAccountByID(ctx context.Context, id string) (*domain.LinkedAccount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.LinkedAccount), args.Error(1)
}

func (m *MockRepository) GetLinkedAccountsByUserID(ctx context.Context, userID string) ([]*domain.LinkedAccount, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.LinkedAccount), args.Error(1)
}

func (m *MockRepository) GetLinkedAccountsByProviderID(ctx context.Context, providerID string) ([]*domain.LinkedAccount, error) {
	args := m.Called(ctx, providerID)
	return args.Get(0).([]*domain.LinkedAccount), args.Error(1)
}

func (m *MockRepository) UpdateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockRepository) DeleteLinkedAccount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Bill operations
func (m *MockRepository) CreateBill(ctx context.Context, bill *domain.Bill) error {
	args := m.Called(ctx, bill)
	return args.Error(0)
}

func (m *MockRepository) GetBillByID(ctx context.Context, id string) (*domain.Bill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Bill), args.Error(1)
}

func (m *MockRepository) GetBillsByLinkedAccountID(ctx context.Context, linkedAccountID string) ([]*domain.Bill, error) {
	args := m.Called(ctx, linkedAccountID)
	return args.Get(0).([]*domain.Bill), args.Error(1)
}

func (m *MockRepository) GetBillsByUserID(ctx context.Context, userID string) ([]*domain.Bill, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Bill), args.Error(1)
}

func (m *MockRepository) GetBillSummaryByUserID(ctx context.Context, userID string) (*domain.BillSummary, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.BillSummary), args.Error(1)
}

func (m *MockRepository) UpdateBill(ctx context.Context, bill *domain.Bill) error {
	args := m.Called(ctx, bill)
	return args.Error(0)
}

func (m *MockRepository) DeleteBill(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestFetchBills(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewBillService(mockRepo)

	// Test case: No linked accounts
	mockRepo.On("GetLinkedAccountsByUserID", mock.Anything, "user1").Return([]*domain.LinkedAccount{}, nil)
	summary, err := service.FetchBills(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Equal(t, 0, summary.BillCount)
	assert.Equal(t, 0.0, summary.TotalDue)
	assert.Empty(t, summary.Bills)

	// Test case: With linked accounts
	accounts := []*domain.LinkedAccount{
		{
			ID:         "acc1",
			UserID:     "user1",
			ProviderID: "mock-provider",
		},
	}
	mockRepo.On("GetLinkedAccountsByUserID", mock.Anything, "user2").Return(accounts, nil)
	summary, err = service.FetchBills(context.Background(), "user2")
	assert.NoError(t, err)
	assert.NotNil(t, summary)
}

func TestRefreshBills(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewBillService(mockRepo)

	accounts := []*domain.LinkedAccount{
		{
			ID:         "acc1",
			UserID:     "user1",
			ProviderID: "mock-provider",
		},
	}
	mockRepo.On("GetLinkedAccountsByUserID", mock.Anything, "user1").Return(accounts, nil)

	// Add expectation for UpdateBill
	mockRepo.On("UpdateBill", mock.Anything, mock.AnythingOfType("*domain.Bill")).Return(nil)

	err := service.RefreshBills(context.Background(), "user1")
	assert.NoError(t, err)
}

func TestPeriodicUpdates(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewBillService(mockRepo)

	users := []*domain.User{
		{ID: "user1"},
		{ID: "user2"},
	}
	mockRepo.On("ListUsers", mock.Anything).Return(users, nil)
	mockRepo.On("GetLinkedAccountsByUserID", mock.Anything, mock.Anything).Return([]*domain.LinkedAccount{}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start periodic updates in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Refresh bills for each user
				for _, user := range users {
					_ = service.RefreshBills(ctx, user.ID)
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	time.Sleep(50 * time.Millisecond) // Wait for some updates
}
