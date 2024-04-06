package user

import (
	"context"
	"db-dashboards/internal/domain/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
)

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var query string

	if limit == math.MaxInt64 {
		query = fmt.Sprintf("SELECT * FROM users ORDER BY created_at OFFSET %v", offset)
	} else {
		query = fmt.Sprintf("SELECT * FROM users ORDER BY created_at LIMIT %v OFFSET %v", limit, offset)
	}

	rows, err := r.DB.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []*entity.User

	for rows.Next() {
		var user entity.User

		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *Repo) getUserByArg(ctx context.Context, argName string, arg any) (*entity.User, error) {
	var query string

	switch arg.(type) {
	case string:
		query = fmt.Sprintf("SELECT * FROM users WHERE %v = '%v'", argName, arg)

	case int, float64:
		query = fmt.Sprintf("SELECT * FROM users WHERE %v = %v", argName, arg)
	}

	row := r.DB.QueryRowxContext(ctx, query)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user entity.User

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return r.getUserByArg(ctx, "id", id)
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.getUserByArg(ctx, "email", email)
}

func (r *Repo) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	result, err := r.DB.NamedQueryContext(ctx,
		`INSERT INTO users (email, hashed_password, created_at, updated_at) 
VALUES (:email, :hashed_password, :created_at, :updated_at) 
RETURNING id, email, hashed_password, created_at, updated_at`,
		&user)
	if err != nil {
		return nil, err
	}

	var usr entity.User

	if result.Next() {
		if err = result.StructScan(&usr); err != nil {
			return nil, err
		}
	}

	return &usr, nil
}

func (r *Repo) DeleteUser(ctx context.Context, id int) (*entity.User, error) {
	row := r.DB.QueryRowxContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user entity.User

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) CheckUniqueConstraints(ctx context.Context, email string) error {
	got, err := r.GetUserByEmail(ctx, email)
	if got != nil || err == nil {
		return ErrEmailExists
	}

	return nil
}
