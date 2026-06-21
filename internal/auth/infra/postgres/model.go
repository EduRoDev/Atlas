package postgres

import (
	"time"

	"github.com/EduRoDev/Atlas/internal/auth/domain"
	"github.com/google/uuid"
)

type userModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	FullName     string    `gorm:"not null"`
	IsVerified   bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}

func (userModel) TableName() string {
	return "users"
}

func (m *userModel) toDomain() *domain.User {
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		FullName:     m.FullName,
		IsVerified:   m.IsVerified,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func fromDomain(u *domain.User) *userModel {
	return &userModel{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		IsVerified:   u.IsVerified,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
