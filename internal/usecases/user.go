package usecases

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/mel-ak/onetap-challenge/internal/adapters/auth"
	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase handles user-related business logic
type UserUsecase struct {
	repo       ports.UserRepository
	jwtService *auth.JWTService
}

// NewUserUsecase creates a new user use case
func NewUserUsecase(repo ports.UserRepository, jwtService *auth.JWTService) *UserUsecase {
	return &UserUsecase{
		repo:       repo,
		jwtService: jwtService,
	}
}

// CreateUser handles POST /users
func (u *UserUsecase) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Check if email exists
	if user, _ := u.repo.GetUserByEmail(r.Context(), req.Email); user != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := domain.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repo.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"user_id": user.ID,
		"message": "User created successfully",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetUser handles GET /users/{user_id}
func (u *UserUsecase) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := u.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"user_id":    user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}
	json.NewEncoder(w).Encode(resp)
}

// UpdateUser handles PUT /users/{user_id}
func (u *UserUsecase) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := u.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Email != "" {
		if !isValidEmail(req.Email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}
		if existing, err := u.repo.GetUserByEmail(r.Context(), req.Email); err == nil {
			if existing.ID != userID {
				http.Error(w, "Email already exists", http.StatusConflict)
				return
			}
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		if len(req.Password) < 8 {
			http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hash)
	}

	user.UpdatedAt = time.Now()

	if err := u.repo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "User updated successfully"}
	json.NewEncoder(w).Encode(resp)
}

// DeleteUser handles DELETE /users/{user_id}
func (u *UserUsecase) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ok, err := u.repo.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	resp := map[string]string{"message": "User deleted successfully"}
	json.NewEncoder(w).Encode(resp)
}

// ListUsers handles GET /users
func (u *UserUsecase) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.repo.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Convert users to response format (excluding sensitive data)
	var response []map[string]interface{}
	for _, user := range users {
		response = append(response, map[string]interface{}{
			"user_id":    user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (u *UserUsecase) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := u.repo.GetUserByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password hashes using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"user_id": user.ID,
	})
}
