package usecases

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"

	"github.com/gorilla/mux"
)

// AccountUsecase handles account-related business logic
type AccountUsecase struct {
	repo  ports.AccountRepository
	cache ports.CacheService
}

// NewAccountUsecase creates a new account use case
func NewAccountUsecase(repo ports.AccountRepository, cache ports.CacheService) *AccountUsecase {
	return &AccountUsecase{repo: repo, cache: cache}
}

// LinkAccount handles POST /accounts/link
func (u *AccountUsecase) LinkAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      int    `json:"user_id"`
		Provider    string `json:"provider"`
		Credentials string `json:"credentials"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if req.UserID <= 0 || req.Provider == "" || req.Credentials == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	account := domain.LinkedAccount{
		UserID:      strconv.Itoa(req.UserID),
		ProviderID:  req.Provider,
		AccountID:   req.Credentials,
		Credentials: req.Credentials,
		CreatedAt:   time.Now(),
	}

	accountID, err := u.repo.SaveAccount(context.Background(), account)
	if err != nil {
		log.Printf("Failed to link account: %v", err)
		http.Error(w, "Failed to link account", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"account_id": accountID,
		"message":    "Account linked successfully",
	}
	json.NewEncoder(w).Encode(resp)
}

// DeleteAccount handles DELETE /accounts/{account_id}
func (u *AccountUsecase) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	accountIDStr := mux.Vars(r)["account_id"]
	// accountID, err := strconv.Atoi(accountIDStr)
	// if err != nil {
	// 	http.Error(w, "Invalid account ID", http.StatusBadRequest)
	// 	return
	// }

	ok, err := u.repo.DeleteAccount(context.Background(), accountIDStr)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	resp := map[string]string{"message": "Account deleted successfully"}
	json.NewEncoder(w).Encode(resp)
}
