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

func (m *MockRepository) GetLinkedAccountsByUserID(ctx context.Context, userID string) ([]*domain.LinkedAccount, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.LinkedAccount), args.Error(1)
}

func (m *MockRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
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

	service.StartPeriodicUpdates(ctx)
	time.Sleep(50 * time.Millisecond) // Wait for some updates
}
