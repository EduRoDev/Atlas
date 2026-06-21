package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	FullName     string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
