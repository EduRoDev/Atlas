package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EduRoDev/Atlas/internal/auth/domain"
	"github.com/google/uuid"
)

type Service struct {
	users  domain.UserRepository
	hasher domain.PasswordHasher
}

func NewService(users domain.UserRepository, hasher domain.PasswordHasher) *Service {
	return &Service{
		users:  users,
		hasher: hasher,
	}
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
	FullName string
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (*domain.User, error) {
	ifExist, err := s.users.ExistsByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("verificando email: %w", err)
	}

	if ifExist {
		return nil, domain.ErrUserAlreadyExists
	}

	ifExist, err = s.users.ExistsByUsername(ctx, in.Username)
	if err != nil {
		return nil, fmt.Errorf("verificando username: %w", err)
	}

	if ifExist {
		return nil, domain.ErrUserAlreadyExists
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hasheando contraseña: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        in.Email,
		Username:     in.Username,
		PasswordHash: hash,
		FullName:     in.FullName,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.users.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("guardando usuario: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, emailOrUsername, password string) (*domain.User, error) {
	user, err := s.users.FindByEmail(ctx, emailOrUsername)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("buscando usuario por email: %w", err)
	}
	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}
