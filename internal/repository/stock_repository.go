package repository

import (
	"auth-register-sistem/internal/model/stock"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type StockRepository interface {
	CreateProduct(s stock.Stock) (uuid.UUID, error)
	GetAllProducts() ([]stock.Stock, error)
	UpdateProductById(s stock.Stock) (uuid.UUID, error)
	DeleteProductById(id string) error
}

type stockRepo struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) StockRepository {
	return &stockRepo{db: db}
}

func (r *stockRepo) CreateProduct(s stock.Stock) (uuid.UUID, error) {
	id := uuid.New()
	s.ID = id
	_, err := r.db.Exec(
		`INSERT INTO stock (id, name, quantity, created_by)
		VALUES ($1, $2, $3, $4)`,
		id, s.Name, s.Quantity, s.CreatedBy)
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, fmt.Errorf("failed to create stock: %w", err)
	}
	return id, nil
}

func (r *stockRepo) GetAllProducts() ([]stock.Stock, error) {
	rows, err := r.db.Query("SELECT id, name, quantity, created_by, created_at, updated_at FROM stock")
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var stocks []stock.Stock
	for rows.Next() {
		var s stock.Stock
		if err := rows.Scan(&s.ID, &s.Name, &s.Quantity, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		stocks = append(stocks, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return stocks, nil
}

func (r *stockRepo) UpdateProductById(s stock.Stock) (uuid.UUID, error) {
	_, err := r.db.Exec(
		`UPDATE stock SET name = $1, quantity = $2, updated_at = $3 WHERE id = $4`,
		s.Name, s.Quantity, time.Now(), s.ID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to update stock: %w", err)
	}
	return s.ID, nil
}

func (r *stockRepo) DeleteProductById(id string) error {
	_, err := r.db.Exec("DELETE FROM stock WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete stock: %w", err)
	}
	return nil
}
