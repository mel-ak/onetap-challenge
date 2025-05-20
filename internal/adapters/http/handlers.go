package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mel-ak/onetap-challenge/internal/ports"
)

// BillHandler handles HTTP requests for bill-related operations
type BillHandler struct {
	billService ports.BillService
}

// NewBillHandler creates a new bill handler
func NewBillHandler(billService ports.BillService) *BillHandler {
	return &BillHandler{
		billService: billService,
	}
}

// RegisterRoutes registers the bill handler routes
func (h *BillHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/bills", h.GetBills).Methods("GET")
	r.HandleFunc("/api/v1/bills/{provider}", h.GetBillsByProvider).Methods("GET")
	r.HandleFunc("/api/v1/bills/refresh", h.RefreshBills).Methods("POST")
}

// GetBills handles GET /api/v1/bills
func (h *BillHandler) GetBills(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, we would get the user ID from the JWT token
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bills, err := h.billService.FetchBills(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

// GetBillsByProvider handles GET /api/v1/bills/{provider}
func (h *BillHandler) GetBillsByProvider(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	providerID := vars["provider"]

	bills, err := h.billService.FetchBillsByProvider(r.Context(), userID, providerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

// RefreshBills handles POST /api/v1/bills/refresh
func (h *BillHandler) RefreshBills(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.billService.RefreshBills(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
