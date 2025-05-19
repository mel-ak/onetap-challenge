package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLRepository implements AccountRepository and BillRepository
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new repository
func NewMySQLRepository(conn string) (*MySQLRepository, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	return &MySQLRepository{db: db}, nil
}

// Implement SaveAccount, GetAccountsByUserID, DeleteAccount, SaveBill similarly to PostgresRepository
