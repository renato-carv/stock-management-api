package stock

import (
	"github.com/google/uuid"
	"time"
)

type Stock struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
}
