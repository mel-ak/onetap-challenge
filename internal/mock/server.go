package mock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Bill struct {
	ID          string    `json:"id"`
	Provider    string    `json:"provider"`
	Amount      float64   `json:"amount"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}

type MockServer struct {
	port      int
	providers []string
}

func NewMockServer(port int) *MockServer {
	return &MockServer{
		port: port,
		providers: []string{
			"Electricity Co",
			"Water Works",
			"Gas Supply",
			"Internet Provider",
			"Phone Company",
		},
	}
}

func (s *MockServer) Start() error {
	http.HandleFunc("/bills", s.handleBills)
	http.HandleFunc("/bills/", s.handleBill)
	http.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Mock server starting on port %d\n", s.port)
	return http.ListenAndServe(addr, nil)
}

func (s *MockServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *MockServer) handleBills(w http.ResponseWriter, r *http.Request) {
	// Simulate random delay
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Simulate random errors
	if rand.Float32() < 0.1 { // 10% chance of error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	bills := s.generateRandomBills(5)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

func (s *MockServer) handleBill(w http.ResponseWriter, r *http.Request) {
	// Simulate random delay
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Simulate random errors
	if rand.Float32() < 0.1 { // 10% chance of error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	bill := s.generateRandomBill()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

func (s *MockServer) generateRandomBills(count int) []Bill {
	bills := make([]Bill, count)
	for i := 0; i < count; i++ {
		bills[i] = s.generateRandomBill()
	}
	return bills
}

func (s *MockServer) generateRandomBill() Bill {
	statuses := []string{"paid", "unpaid", "overdue"}
	status := statuses[rand.Intn(len(statuses))]

	dueDate := time.Now().AddDate(0, 0, rand.Intn(30))

	return Bill{
		ID:          fmt.Sprintf("BILL-%d", rand.Intn(10000)),
		Provider:    s.providers[rand.Intn(len(s.providers))],
		Amount:      float64(rand.Intn(1000)) + rand.Float64(),
		DueDate:     dueDate,
		Status:      status,
		Description: fmt.Sprintf("Bill for %s services", s.providers[rand.Intn(len(s.providers))]),
	}
}
