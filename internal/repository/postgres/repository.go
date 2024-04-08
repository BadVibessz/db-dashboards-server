package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"db-dashboards/internal/domain/entity/postgres"
)

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) GetAllTables(ctx context.Context) ([]*postgres.Table, error) {
	rows, err := r.DB.QueryxContext(ctx, "SELECT * FROM information_schema.tables WHERE table_schema = 'public'") // todo: return only public tables?
	if err != nil {
		return nil, err
	}

	columnNames, _ := rows.Columns()

	var tables []*postgres.Table

	for rows.Next() {
		columns := make([]string, len(columnNames))
		columnPointers := make([]any, len(columnNames))

		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		var table postgres.Table

		for i, columnName := range columnNames {
			if columnName == "table_name" {
				table.Name = columns[i]
				break
			}
		}

		tables = append(tables, &table)
	}

	return tables, nil
}

func (r *Repo) GetColumnsFromTable(ctx context.Context, tableName string) ([]*postgres.Column, error) {
	rows, err := r.DB.QueryxContext(ctx,
		fmt.Sprintf("SELECT * FROM information_schema.columns WHERE table_name = '%v' order by ordinal_position", tableName))

	if err != nil {
		return nil, err
	}

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var columns []*postgres.Column

	for rows.Next() {
		row := make(postgres.Row)

		if err = rows.MapScan(row); err != nil {
			return nil, err
		}

		cols := make([]string, len(columnNames))
		columnPointers := make([]any, len(columnNames))

		for i, _ := range cols {
			columnPointers[i] = &cols[i]
		}

		rows.Scan(columnPointers...)

		var column postgres.Column

		for i, columnName := range columnNames {
			if columnName == "column_name" {
				column.Name = cols[i]
			}

			if columnName == "data_type" {
				column.Type = strings.ToLower(columnTypes[i].DatabaseTypeName())
			}

			if column.Name != "" && column.Type != "" {
				break
			}
		}

		columns = append(columns, &column)
	}

	return columns, nil
}

func (r *Repo) GetAllRowsFromTable(ctx context.Context, tableName string) ([]*postgres.Row, error) {

	dbRows, err := r.DB.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %v", tableName))
	if err != nil {
		return nil, err
	}

	var rows []*postgres.Row

	for dbRows.Next() {
		row := make(postgres.Row)

		if err = dbRows.MapScan(row); err != nil {
			return nil, err
		}

		rows = append(rows, &row)
	}

	return rows, nil
}
