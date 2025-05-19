package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"

	_ "github.com/lib/pq"
)

// PostgresRepository implements UserRepository, AccountRepository, and BillRepository
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new repository
func NewPostgresRepository(conn string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

// SaveUser saves a user
func (r *PostgresRepository) SaveUser(ctx context.Context, user domain.User) (int, error) {
	query := `INSERT INTO users (email, password_hash, created_at) 
              VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, time.Now()).Scan(&id)
	return id, err
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

// GetUserByEmail retrieves a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

// UpdateUser updates a user
func (r *PostgresRepository) UpdateUser(ctx context.Context, user domain.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.ID)
	return err
}

// DeleteUser deletes a user
func (r *PostgresRepository) DeleteUser(ctx context.Context, userID int) (bool, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

// SaveAccount saves an account
func (r *PostgresRepository) SaveAccount(ctx context.Context, account domain.Account) (int, error) {
	query := `INSERT INTO accounts (user_id, provider, credentials, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, account.UserID, account.Provider, account.Credentials, time.Now()).Scan(&id)
	return id, err
}

// GetAccountsByUserID retrieves accounts for a user
func (r *PostgresRepository) GetAccountsByUserID(ctx context.Context, userID int) ([]domain.Account, error) {
	query := `SELECT id, user_id, provider, credentials, created_at FROM accounts WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var acc domain.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Provider, &acc.Credentials, &acc.CreatedAt); err != nil {
			continue
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// DeleteAccount deletes an account
func (r *PostgresRepository) DeleteAccount(ctx context.Context, accountID int) (bool, error) {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, accountID)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

// SaveBill saves a bill
func (r *PostgresRepository) SaveBill(ctx context.Context, bill domain.Bill) error {
	query := `INSERT INTO bills (account_id, provider, amount, due_date, status, created_at)
              VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, bill.AccountID, bill.Provider, bill.Amount, bill.DueDate, bill.Status, time.Now())
	return err
}
