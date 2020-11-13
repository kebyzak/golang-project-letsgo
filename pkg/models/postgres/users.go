package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"kazybek.net/finalproject/pkg/models"
	"time"
)

const (
	insertUser = "INSERT INTO users (name, email, hashed_password, created) VALUES($1,$2,$3,$4) RETURNING id"
	selectUser = "SELECT id, hashed_password FROM users WHERE email = $1 AND active = TRUE"
)

type UserModel struct {
	Pool *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	var id uint64
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	row := m.Pool.QueryRow(context.Background(), insertUser, name, email, string(hashedPassword), time.Now())
	err = row.Scan(&id)
	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			return models.ErrDuplicateEmail
		}
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := m.Pool.QueryRow(context.Background(), selectUser, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
