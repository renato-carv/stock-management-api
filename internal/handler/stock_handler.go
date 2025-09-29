package handler

import (
	"auth-register-sistem/internal/middleware"
	"auth-register-sistem/internal/model/stock"
	"auth-register-sistem/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type StockHandler struct {
	Repo repository.StockRepository
}

func NewStockHandler(repo repository.StockRepository) *StockHandler {
	return &StockHandler{Repo: repo}
}

func (h *StockHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req stock.Stock
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	userIDVal := r.Context().Value(middleware.UserIDKey)
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusUnauthorized)
		return
	}

	createdByUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	req.CreatedBy = createdByUUID

	id, err := h.Repo.CreateProduct(req)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Product created successfully",
	})
}

func (h *StockHandler) GetAllProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := h.Repo.GetAllProducts()
	if err != nil {
		http.Error(writer, "Failed to get products", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(products)
}

func (h *StockHandler) UpdateProductById(writer http.ResponseWriter, request *http.Request) {
	idStr := request.URL.Query().Get("id")
	if idStr == "" {
		http.Error(writer, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid id format", http.StatusBadRequest)
		return
	}

	var req stock.Stock
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.ID = id

	updatedId, err := h.Repo.UpdateProductById(req)
	if err != nil {
		http.Error(writer, "Failed to update product", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"id":      updatedId,
		"message": "Product updated successfully",
	})
}

func (h *StockHandler) DeleteProductById(writer http.ResponseWriter, request *http.Request) {
	idStr := request.URL.Query().Get("id")
	if idStr == "" {
		http.Error(writer, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid id format", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteProductById(id.String())
	if err != nil {
		http.Error(writer, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message": "Product deleted successfully",
	})
}
