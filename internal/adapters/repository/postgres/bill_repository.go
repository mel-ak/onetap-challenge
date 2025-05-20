package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"
	"github.com/mel-ak/onetap-challenge/internal/ports"
)

type repository struct {
	db *sql.DB
}

// NewRepository creates a new PostgreSQL repository
func NewRepository(db *sql.DB) ports.Repository {
	return &repository{db: db}
}

// Bill operations
func (r *repository) CreateBill(ctx context.Context, bill *domain.Bill) error {
	query := `
		INSERT INTO bills (
			id, linked_account_id, provider_id, amount, due_date, 
			status, bill_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		bill.ID,
		bill.LinkedAccountID,
		bill.ProviderID,
		bill.Amount,
		bill.DueDate,
		bill.Status,
		bill.BillDate,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *repository) GetBillByID(ctx context.Context, id string) (*domain.Bill, error) {
	query := `
		SELECT id, linked_account_id, provider_id, amount, due_date,
			status, bill_date, created_at, updated_at
		FROM bills
		WHERE id = $1
	`
	bill := &domain.Bill{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bill.ID,
		&bill.LinkedAccountID,
		&bill.ProviderID,
		&bill.Amount,
		&bill.DueDate,
		&bill.Status,
		&bill.BillDate,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return bill, err
}

func (r *repository) GetBillsByLinkedAccountID(ctx context.Context, linkedAccountID string) ([]*domain.Bill, error) {
	query := `
		SELECT id, linked_account_id, provider_id, amount, due_date,
			status, bill_date, created_at, updated_at
		FROM bills
		WHERE linked_account_id = $1
		ORDER BY due_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, linkedAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*domain.Bill
	for rows.Next() {
		bill := &domain.Bill{}
		err := rows.Scan(
			&bill.ID,
			&bill.LinkedAccountID,
			&bill.ProviderID,
			&bill.Amount,
			&bill.DueDate,
			&bill.Status,
			&bill.BillDate,
			&bill.CreatedAt,
			&bill.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}
	return bills, rows.Err()
}

func (r *repository) GetBillsByUserID(ctx context.Context, userID string) ([]*domain.Bill, error) {
	query := `
		SELECT b.id, b.linked_account_id, b.provider_id, b.amount, b.due_date,
			b.status, b.bill_date, b.created_at, b.updated_at
		FROM bills b
		JOIN linked_accounts la ON b.linked_account_id = la.id
		WHERE la.user_id = $1
		ORDER BY b.due_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*domain.Bill
	for rows.Next() {
		bill := &domain.Bill{}
		err := rows.Scan(
			&bill.ID,
			&bill.LinkedAccountID,
			&bill.ProviderID,
			&bill.Amount,
			&bill.DueDate,
			&bill.Status,
			&bill.BillDate,
			&bill.CreatedAt,
			&bill.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}
	return bills, rows.Err()
}

func (r *repository) GetBillSummaryByUserID(ctx context.Context, userID string) (*domain.BillSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(b.amount), 0) as total_amount,
			COUNT(CASE WHEN b.status = 'unpaid' THEN 1 END) as due_bills,
			COUNT(CASE WHEN b.status = 'overdue' THEN 1 END) as overdue_bills
		FROM bills b
		JOIN linked_accounts la ON b.linked_account_id = la.id
		WHERE la.user_id = $1
	`
	summary := &domain.BillSummary{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&summary.TotalAmount,
		&summary.DueBills,
		&summary.OverdueBills,
	)
	if err != nil {
		return nil, err
	}

	// Get the actual bills
	bills, err := r.GetBillsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary.Bills = make([]domain.Bill, len(bills))
	for i, bill := range bills {
		summary.Bills[i] = *bill
	}

	return summary, nil
}

func (r *repository) UpdateBill(ctx context.Context, bill *domain.Bill) error {
	query := `
		UPDATE bills
		SET amount = $1, due_date = $2, status = $3, bill_date = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		bill.Amount,
		bill.DueDate,
		bill.Status,
		bill.BillDate,
		time.Now(),
		bill.ID,
	)
	return err
}

func (r *repository) DeleteBill(ctx context.Context, id string) error {
	query := `DELETE FROM bills WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// User operations
func (r *repository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *repository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Password,
		time.Now(),
		user.ID,
	)
	return err
}

func (r *repository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Provider operations
func (r *repository) CreateProvider(ctx context.Context, provider *domain.Provider) error {
	query := `
		INSERT INTO providers (id, name, api_endpoint, auth_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		provider.ID,
		provider.Name,
		provider.APIEndpoint,
		provider.AuthType,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *repository) GetProviderByID(ctx context.Context, id string) (*domain.Provider, error) {
	query := `
		SELECT id, name, api_endpoint, auth_type, created_at, updated_at
		FROM providers
		WHERE id = $1
	`
	provider := &domain.Provider{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&provider.ID,
		&provider.Name,
		&provider.APIEndpoint,
		&provider.AuthType,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return provider, err
}

func (r *repository) GetProviderByName(ctx context.Context, name string) (*domain.Provider, error) {
	query := `
		SELECT id, name, api_endpoint, auth_type, created_at, updated_at
		FROM providers
		WHERE name = $1
	`
	provider := &domain.Provider{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&provider.ID,
		&provider.Name,
		&provider.APIEndpoint,
		&provider.AuthType,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return provider, err
}

func (r *repository) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	query := `
		SELECT id, name, api_endpoint, auth_type, created_at, updated_at
		FROM providers
		ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*domain.Provider
	for rows.Next() {
		provider := &domain.Provider{}
		err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&provider.APIEndpoint,
			&provider.AuthType,
			&provider.CreatedAt,
			&provider.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, rows.Err()
}

func (r *repository) UpdateProvider(ctx context.Context, provider *domain.Provider) error {
	query := `
		UPDATE providers
		SET name = $1, api_endpoint = $2, auth_type = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.APIEndpoint,
		provider.AuthType,
		time.Now(),
		provider.ID,
	)
	return err
}

func (r *repository) DeleteProvider(ctx context.Context, id string) error {
	query := `DELETE FROM providers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// LinkedAccount operations
func (r *repository) CreateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error {
	query := `
		INSERT INTO linked_accounts (
			id, user_id, provider_id, account_id, credentials,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.ProviderID,
		account.AccountID,
		account.Credentials,
		account.Status,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *repository) GetLinkedAccountByID(ctx context.Context, id string) (*domain.LinkedAccount, error) {
	query := `
		SELECT id, user_id, provider_id, account_id, credentials,
			status, created_at, updated_at
		FROM linked_accounts
		WHERE id = $1
	`
	account := &domain.LinkedAccount{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.ProviderID,
		&account.AccountID,
		&account.Credentials,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return account, err
}

func (r *repository) GetLinkedAccountsByUserID(ctx context.Context, userID string) ([]*domain.LinkedAccount, error) {
	query := `
		SELECT id, user_id, provider_id, account_id, credentials,
			status, created_at, updated_at
		FROM linked_accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.LinkedAccount
	for rows.Next() {
		account := &domain.LinkedAccount{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.ProviderID,
			&account.AccountID,
			&account.Credentials,
			&account.Status,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, rows.Err()
}

func (r *repository) GetLinkedAccountsByProviderID(ctx context.Context, providerID string) ([]*domain.LinkedAccount, error) {
	query := `
		SELECT id, user_id, provider_id, account_id, credentials,
			status, created_at, updated_at
		FROM linked_accounts
		WHERE provider_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.LinkedAccount
	for rows.Next() {
		account := &domain.LinkedAccount{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.ProviderID,
			&account.AccountID,
			&account.Credentials,
			&account.Status,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, rows.Err()
}

func (r *repository) UpdateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error {
	query := `
		UPDATE linked_accounts
		SET credentials = $1, status = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		account.Credentials,
		account.Status,
		time.Now(),
		account.ID,
	)
	return err
}

func (r *repository) DeleteLinkedAccount(ctx context.Context, id string) error {
	query := `DELETE FROM linked_accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
