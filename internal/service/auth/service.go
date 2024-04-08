package auth

import (
	"context"
	"db-dashboards/internal/domain/entity"
)

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type Hasher interface {
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type Service struct {
	UserRepo UserRepo
	Hasher   Hasher
}

func New(userRepo UserRepo, hasher Hasher) *Service {
	return &Service{
		UserRepo: userRepo,
		Hasher:   hasher,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*entity.User, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = s.Hasher.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
