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
func (r *PostgresRepository) SaveUser(ctx context.Context, user domain.User) (string, error) {
	query := `INSERT INTO users (id, email, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.Password, time.Now(), time.Now()).Scan(&id)
	return id, err
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

// GetUserByEmail retrieves a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

// UpdateUser updates a user
func (r *PostgresRepository) UpdateUser(ctx context.Context, user domain.User) error {
	query := `UPDATE users SET email = $1, password = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Password, time.Now(), user.ID)
	return err
}

// DeleteUser deletes a user
func (r *PostgresRepository) DeleteUser(ctx context.Context, userID string) (bool, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

// SaveAccount saves an account
func (r *PostgresRepository) SaveAccount(ctx context.Context, account domain.LinkedAccount) (string, error) {
	query := `INSERT INTO linked_accounts (id, user_id, provider_id, account_id, credentials, status, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query, account.ID, account.UserID, account.ProviderID, account.AccountID, account.Credentials, account.Status, time.Now(), time.Now()).Scan(&id)
	return id, err
}

// GetAccountsByUserID retrieves accounts for a user
func (r *PostgresRepository) GetAccountsByUserID(ctx context.Context, userID string) ([]domain.LinkedAccount, error) {
	query := `SELECT id, user_id, provider_id, account_id, credentials, status, created_at, updated_at FROM linked_accounts WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.LinkedAccount
	for rows.Next() {
		var acc domain.LinkedAccount
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.ProviderID, &acc.AccountID, &acc.Credentials, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
			continue
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// DeleteAccount deletes an account
func (r *PostgresRepository) DeleteAccount(ctx context.Context, accountID string) (bool, error) {
	query := `DELETE FROM linked_accounts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, accountID)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

// SaveBill saves a bill
func (r *PostgresRepository) SaveBill(ctx context.Context, bill domain.Bill) error {
	query := `INSERT INTO bills (id, linked_account_id, provider_id, amount, due_date, status, bill_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, bill.ID, bill.LinkedAccountID, bill.ProviderID, bill.Amount, bill.DueDate, bill.Status, bill.BillDate, time.Now(), time.Now())
	return err
}
