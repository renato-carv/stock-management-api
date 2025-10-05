package transaction

import (
	"github.com/google/uuid"
	"time"
)

type TransactionType string

const (
	TypeIn  TransactionType = "ENTRY"
	TypeOut TransactionType = "EXIT"
)

type Transaction struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Quantity  int             `json:"quantity"`
	Type      TransactionType `json:"type"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	CreatedBy uuid.UUID       `json:"created_by"`
}
