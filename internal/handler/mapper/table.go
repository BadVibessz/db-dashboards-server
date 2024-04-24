package mapper

import (
	"db-dashboards/internal/domain/entity/postgres"
	"db-dashboards/internal/handler/response"
)

func MapTableToTableResponse(table *postgres.Table) response.GetTableResponse {
	return response.GetTableResponse{
		Name: table.Name,
	}
}

func MapColumnToColumnResponse(column *postgres.Column) response.GetColumnsResponse {
	return response.GetColumnsResponse{
		Name: column.Name,
		Type: column.Type,
	}
}
