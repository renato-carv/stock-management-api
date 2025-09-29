package repository

import (
	"auth-register-sistem/internal/model/user"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user user.User) (uuid.UUID, error)
	FindByEmail(email string) (*user.User, error)
	FindByUsername(username string) (*user.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u user.User) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.db.Exec(
		"INSERT INTO users (id, name, username, email, password) VALUES ($1, $2, $3, $4, $5)",
		id, u.Name, u.Username, u.Email, u.Password,
	)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *userRepo) FindByEmail(email string) (*user.User, error) {
	row := r.db.QueryRow("SELECT id, name, username, email, password FROM users WHERE email = $1", email)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // user not found
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepo) FindByUsername(username string) (*user.User, error) {
	row := r.db.QueryRow("SELECT id, name, username, email, password FROM users WHERE username = $1", username)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // user not found
		}
		return nil, err
	}
	return u, nil
}
