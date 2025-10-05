package routes

import (
	"auth-register-sistem/internal/handler"
	"auth-register-sistem/internal/middleware"
	"net/http"
)

func SetupRoutes(userHandler *handler.UserHandler, stockHandler *handler.StockHandler, transactionHandler *handler.TransactionHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	// Stock routes
	mux.HandleFunc("/stock", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			stockHandler.GetAllProducts(w, r)
		case http.MethodPost:
			stockHandler.CreateProduct(w, r)
		case http.MethodPut:
			stockHandler.UpdateProductById(w, r)
		case http.MethodDelete:
			stockHandler.DeleteProductById(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Transaction routes
	mux.HandleFunc("/transaction", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			transactionHandler.CreateTransaction(w, r)
		case http.MethodGet:
			transactionHandler.GetAllTransactions(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	return mux
}
