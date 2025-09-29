package handler

import (
	"auth-register-sistem/internal/model/user"
	"auth-register-sistem/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// Register new users
func (h *UserHandler) Register(writer http.ResponseWriter, request *http.Request) {
	var req user.User
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(req.Password)), bcrypt.DefaultCost)
	if err != nil {
		http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPassword)

	id, err := h.Repo.Create(req)
	if err != nil {
		http.Error(writer, "Failed to create user", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"id":      id,
		"message": "User created successfully",
	})
}

// Login existing users
func (h *UserHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	userData, err := h.Repo.FindByUsername(req.Username)
	if err != nil {
		http.Error(writer, "Failed to find user", http.StatusNotFound)
		return
	}

	if userData == nil {
		http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(strings.TrimSpace(req.Password))); err != nil {
		http.Error(writer, "Invalid password", http.StatusUnauthorized)
		log.Println(err)
		return
	} else {
		log.Println("Login successful")
	}

	//Generate JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(writer, "JWT secret not found", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userData.ID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(writer, "Failed to sign token", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{
		"token":   tokenStr,
		"message": "Login successful",
	})
}
