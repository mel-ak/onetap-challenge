package usecases

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"

	"github.com/gorilla/mux"
)

// ProviderUsecase handles provider-related business logic
type ProviderUsecase struct {
	repo ports.Repository
}

// NewProviderUsecase creates a new provider use case
func NewProviderUsecase(repo ports.Repository) *ProviderUsecase {
	return &ProviderUsecase{repo: repo}
}

// CreateProvider handles POST /providers
func (u *ProviderUsecase) CreateProvider(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		APIEndpoint string `json:"api_endpoint"`
		AuthType    string `json:"auth_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if req.Name == "" || req.APIEndpoint == "" || req.AuthType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	provider := &domain.Provider{
		ID:          uuid.New().String(),
		Name:        req.Name,
		APIEndpoint: req.APIEndpoint,
		AuthType:    req.AuthType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.repo.CreateProvider(r.Context(), provider); err != nil {
		http.Error(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"provider_id": provider.ID,
		"message":     "Provider created successfully",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// ListProviders handles GET /providers
func (u *ProviderUsecase) ListProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := u.repo.ListProviders(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch providers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// GetProvider handles GET /providers/{provider_id}
func (u *ProviderUsecase) GetProvider(w http.ResponseWriter, r *http.Request) {
	providerID := mux.Vars(r)["provider_id"]
	if providerID == "" {
		http.Error(w, "Invalid provider ID", http.StatusBadRequest)
		return
	}

	provider, err := u.repo.GetProviderByID(r.Context(), providerID)
	if err != nil {
		http.Error(w, "Failed to fetch provider", http.StatusInternalServerError)
		return
	}
	if provider == nil {
		http.Error(w, "Provider not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}
