package handler

import (
	"auth-register-sistem/internal/middleware"
	"auth-register-sistem/internal/model/transaction"
	"auth-register-sistem/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type TransactionHandler struct {
	Repo repository.TransactionRepository
}

func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{Repo: repo}
}

// CreateTransaction handles the creation of a new transaction
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(middleware.UserIDKey)
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user ID from context
	userID, ok := userIDVal.(string)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		Name	 string `json:"name"`
		Quantity int    `json:"quantity"`
		Type     string `json:"type"`
	}

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate transaction type
	if req.Type != "ENTRY" && req.Type != "EXIT" {
		http.Error(w, "Invalid transaction type", http.StatusBadRequest)
		return
	}

	//validate quantity
	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	//validate name
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Create transaction model
	transactionData := transaction.Transaction{
		Name:      req.Name,
		Quantity:  req.Quantity,
		Type:      transaction.TransactionType(req.Type),
		CreatedBy: uuid.MustParse(userID),
	}

	// Call repository to create transaction
	id, err := h.Repo.CreateTransaction(transactionData)
	if err != nil {
		http.Error(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Transaction created successfully",
	})
}

// GetAllTransactions retrieves all transactions
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.Repo.GetAllTransactions()
	if err != nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}
}