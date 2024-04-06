package user

import (
	"context"

	"db-dashboards/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type Repo interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) ([]*entity.User, error)
	DeleteUser(ctx context.Context, id int) (*entity.User, error)
	CheckUniqueConstraints(ctx context.Context, email string) error
}

type Hasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type Service struct {
	Repo   Repo
	Hasher Hasher
}

func New(repo Repo, hasher Hasher) *Service {
	return &Service{
		Repo:   repo,
		Hasher: hasher,
	}
}

func (s *Service) GetAllUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	return s.Repo.GetAllUsers(ctx, offset, limit)
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.Repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) RegisterUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// ensure that user with this email does not exist
	err := s.Repo.CheckUniqueConstraints(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	// user model sent with plain password
	hash, err := s.Hasher.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = string(hash)

	created, err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int) (*entity.User, error) { // todo: authorize admin rights
	deleted, err := s.Repo.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}
