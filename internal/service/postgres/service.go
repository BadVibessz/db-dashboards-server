package postgres

import (
	"context"

	"db-dashboards/internal/domain/entity/postgres"

	postgresrepo "db-dashboards/internal/repository/postgres"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) GetAllTables(ctx context.Context, repo *postgresrepo.Repo) ([]*postgres.Table, error) {
	return repo.GetAllTables(ctx)
}

func (s *Service) GetColumnsFromTable(ctx context.Context, repo *postgresrepo.Repo, tableName string) ([]*postgres.Column, error) {
	return repo.GetColumnsFromTable(ctx, tableName)
}

func (s *Service) GetAllRowsFromTable(ctx context.Context, repo *postgresrepo.Repo, tableName string) ([]*postgres.Row, error) {
	return repo.GetAllRowsFromTable(ctx, tableName)
}
