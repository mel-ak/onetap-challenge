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

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
	)
	return err
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByEmail retrieves a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// UpdateUser updates a user
func (r *PostgresRepository) UpdateUser(ctx context.Context, user *domain.User) error {
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

// ListUsers retrieves all users
func (r *PostgresRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// SaveAccount saves an account
func (r *PostgresRepository) SaveAccount(ctx context.Context, account domain.LinkedAccount) (string, error) {
	query := `INSERT INTO linked_accounts (id, user_id, provider_id, account_id, credentials, status, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query, account.ID, account.UserID, account.ProviderID, account.AccountID, account.Credentials, account.Status, time.Now(), time.Now()).Scan(&id)
	return id, err
}

// CreateLinkedAccount creates a new linked account
func (r *PostgresRepository) CreateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error {
	query := `INSERT INTO linked_accounts (id, user_id, provider_id, account_id, credentials, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
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

// DeleteLinkedAccount deletes a linked account
func (r *PostgresRepository) DeleteLinkedAccount(ctx context.Context, id string) error {
	query := `DELETE FROM linked_accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CreateBill creates a new bill
func (r *PostgresRepository) CreateBill(ctx context.Context, bill *domain.Bill) error {
	query := `INSERT INTO bills (id, linked_account_id, provider_id, amount, due_date, status, bill_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
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

// SaveBill saves a bill
func (r *PostgresRepository) SaveBill(ctx context.Context, bill domain.Bill) error {
	query := `INSERT INTO bills (id, linked_account_id, provider_id, amount, due_date, status, bill_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, bill.ID, bill.LinkedAccountID, bill.ProviderID, bill.Amount, bill.DueDate, bill.Status, bill.BillDate, time.Now(), time.Now())
	return err
}

// DeleteBill deletes a bill
func (r *PostgresRepository) DeleteBill(ctx context.Context, id string) error {
	query := `DELETE FROM bills WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CreateProvider creates a new provider
func (r *PostgresRepository) CreateProvider(ctx context.Context, provider *domain.Provider) error {
	query := `INSERT INTO providers (id, name, api_endpoint, auth_type, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6)`
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

// DeleteProvider deletes a provider
func (r *PostgresRepository) DeleteProvider(ctx context.Context, id string) error {
	query := `DELETE FROM providers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRepository) GetBillByID(ctx context.Context, id string) (*domain.Bill, error) {
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

func (r *PostgresRepository) GetBillSummaryByUserID(ctx context.Context, userID string) (*domain.BillSummary, error) {
	query := `
		SELECT 
			COUNT(*) as bill_count,
			COALESCE(SUM(CASE WHEN b.status IN ('unpaid', 'overdue') THEN b.amount ELSE 0 END), 0) as total_due
		FROM bills b
		JOIN linked_accounts la ON b.linked_account_id = la.id
		WHERE la.user_id = $1
	`
	summary := &domain.BillSummary{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&summary.BillCount,
		&summary.TotalDue,
	)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (r *PostgresRepository) GetBillsByLinkedAccountID(ctx context.Context, linkedAccountID string) ([]*domain.Bill, error) {
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

func (r *PostgresRepository) GetBillsByUserID(ctx context.Context, userID string) ([]*domain.Bill, error) {
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

func (r *PostgresRepository) GetLinkedAccountByID(ctx context.Context, id string) (*domain.LinkedAccount, error) {
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

func (r *PostgresRepository) GetLinkedAccountsByProviderID(ctx context.Context, providerID string) ([]*domain.LinkedAccount, error) {
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

func (r *PostgresRepository) GetLinkedAccountsByUserID(ctx context.Context, userID string) ([]*domain.LinkedAccount, error) {
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

func (r *PostgresRepository) GetProviderByID(ctx context.Context, id string) (*domain.Provider, error) {
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

func (r *PostgresRepository) GetProviderByName(ctx context.Context, name string) (*domain.Provider, error) {
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

func (r *PostgresRepository) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	query := `SELECT id, name, api_endpoint, auth_type, created_at, updated_at FROM providers ORDER BY name`
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

func (r *PostgresRepository) UpdateBill(ctx context.Context, bill *domain.Bill) error {
	query := `UPDATE bills SET amount = $1, due_date = $2, status = $3, bill_date = $4, updated_at = $5 WHERE id = $6`
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

func (r *PostgresRepository) UpdateLinkedAccount(ctx context.Context, account *domain.LinkedAccount) error {
	query := `UPDATE linked_accounts SET credentials = $1, status = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query,
		account.Credentials,
		account.Status,
		time.Now(),
		account.ID,
	)
	return err
}

func (r *PostgresRepository) UpdateProvider(ctx context.Context, provider *domain.Provider) error {
	query := `UPDATE providers SET name = $1, api_endpoint = $2, auth_type = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.APIEndpoint,
		provider.AuthType,
		time.Now(),
		provider.ID,
	)
	return err
}
