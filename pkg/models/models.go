package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Jazba struct {
	ID int
	Takyryp string
	Mazmuny string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID int
	Name string
	Email string
	Password string
	Created time.Time
	Active bool
}