package usecases

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase handles user-related business logic
type UserUsecase struct {
	repo ports.UserRepository
}

// NewUserUsecase creates a new user use case
func NewUserUsecase(repo ports.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
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
	if _, err := u.repo.GetUserByEmail(r.Context(), req.Email); err == nil {
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
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	userID, err := u.repo.SaveUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"user_id": userID,
		"message": "User created successfully",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetUser handles GET /users/{user_id}
func (u *UserUsecase) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
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
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
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
		if existing, err := u.repo.GetUserByEmail(r.Context(), req.Email); err == nil && existing.ID != userID {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
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
		user.PasswordHash = string(hash)
	}

	if err := u.repo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "User updated successfully"}
	json.NewEncoder(w).Encode(resp)
}

// DeleteUser handles DELETE /users/{user_id}
func (u *UserUsecase) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
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

// isValidEmail validates email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
