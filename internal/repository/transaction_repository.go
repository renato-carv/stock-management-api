package repository

import (
	"auth-register-sistem/internal/model/transaction"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateTransaction(t transaction.Transaction) (uuid.UUID, error)
	GetAllTransactions() ([]transaction.Transaction, error)
}

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) CreateTransaction(t transaction.Transaction) (uuid.UUID, error) {
	id := uuid.New()
	t.ID = id

	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO transactions (id, name, quantity, type, created_by)
		VALUES ($1, $2, $3, $4, $5)`,
		id, t.Name, t.Quantity, t.Type, t.CreatedBy)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	var currentQty int
	err = tx.QueryRow(
		`SELECT quantity FROM stock WHERE name = $1 FOR UPDATE`,
		t.Name).Scan(&currentQty)

	if err == sql.ErrNoRows {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("stock item not found")
	} else if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to fetch current stock quantity: %w", err)
	}

	newQty := currentQty
	switch t.Type {
	case "ENTRY":
		newQty += t.Quantity
	case "EXIT":
		if t.Quantity > currentQty {
			tx.Rollback()
			return uuid.Nil, fmt.Errorf("insufficient stock for EXIT transaction")
		}
		newQty -= t.Quantity
	default:
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("invalid transaction type: %s", t.Type)
	}

	_, err = tx.Exec(
		`UPDATE stock SET quantity = $1, updated_at = now() WHERE name = $2`,
		newQty, t.Name)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to update stock quantity: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return t.ID, nil
}

func (r *TransactionRepo) GetAllTransactions() ([]transaction.Transaction, error) {
	rows, err := r.db.Query("SELECT id, name, quantity, type, created_by, created_at, updated_at FROM transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	defer rows.Close()

	var transactions []transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		if err := rows.Scan(&t.ID, &t.Name, &t.Quantity, &t.Type, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return transactions, nil
}
