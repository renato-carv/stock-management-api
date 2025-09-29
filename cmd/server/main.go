package main

import (
	"auth-register-sistem/internal/config"
	"auth-register-sistem/internal/handler"
	"auth-register-sistem/internal/repository"
	"auth-register-sistem/internal/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbCfg := config.NewDBConfig()
	dbConn, err := config.SetupDb(dbCfg)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}
	defer dbConn.Close()
	log.Println("Connected to database")
	userRepo := repository.NewUserRepository(dbConn)
	stockRepo := repository.NewStockRepository(dbConn)
	userHandler := handler.NewUserHandler(userRepo)
	stockHandler := handler.NewStockHandler(stockRepo)

	mux := routes.SetupRoutes(userHandler, stockHandler)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}