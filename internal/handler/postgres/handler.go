package postgres

import (
	"context"
	"database/sql"
	"db-dashboards/internal/domain/entity/postgres"
	"db-dashboards/internal/handler/mapper"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"

	postgresrepo "db-dashboards/internal/repository/postgres"
	handlerutils "db-dashboards/pkg/utils/handler"

	sliceutils "db-dashboards/pkg/utils/slice"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service interface {
	GetAllTables(ctx context.Context, repo *postgresrepo.Repo) ([]*postgres.Table, error)
	GetColumnsFromTable(ctx context.Context, repo *postgresrepo.Repo, tableName string) ([]*postgres.Column, error)
	GetAllRowsFromTable(ctx context.Context, repo *postgresrepo.Repo, tableName string) ([]*postgres.Row, error)
}

type Middleware = func(http.Handler) http.Handler

type Handler struct {
	Service     Service
	Middlewares []Middleware

	logger    *logrus.Logger
	validator *validator.Validate
}

func New(service Service,
	logger *logrus.Logger,
	validator *validator.Validate,
	middlewares ...Middleware,
) *Handler {
	return &Handler{
		Service:     service,
		Middlewares: middlewares,
		logger:      logger,
		validator:   validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(h.Middlewares...)

		r.Get("/tables", h.GetAllTables)
		r.Get("/columns", h.GetColumnsFromTable)
		r.Get("/data", h.GetAllRowsFromTable)
	})

	return router
}

// GetAllTables godoc
//
//		@Summary		Get all tables from db
//		@Description	Get all tables from db
//		@Security		JWT
//		@Tags			Postgres
//	 	@Param 			connection-string 	header 	string true "connection string"
//		@Produce		json
//		@Success		200	{object}	[]response.GetTableResponse
//		@Failure		401	{string}	Unauthorized
//		@Router			/db-dashboards/api/v1/postgres/tables [get]
func (h *Handler) GetAllTables(rw http.ResponseWriter, req *http.Request) {
	connStr, err := handlerutils.GetStringHeaderByKey(req, "connection-string")
	if err != nil {
		msg := "no connection string header provided"

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		msg := fmt.Sprintf("cannot connect to db with conn str: %v", connStr)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	// TODO: ping db first
	repo := postgresrepo.New(sqlx.NewDb(conn, "postgres"))

	tables, err := h.Service.GetAllTables(req.Context(), repo)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch tables from db")

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return

	}

	render.JSON(rw, req, sliceutils.Map(tables, mapper.MapTableToTableResponse))
	rw.WriteHeader(http.StatusOK)
}

// GetColumnsFromTable godoc
//
//		@Summary		Get all columns from table
//		@Description	Get all columns from table
//		@Security		JWT
//		@Tags			Postgres
//	 	@Param 			connection-string 	header 	string true "connection string"
//	 	@Param 			table-name 	header 	string true "name of the table"
//		@Produce		json
//		@Success		200	{object}	[]response.GetColumnsResponse
//		@Failure		401	{string}	Unauthorized
//		@Router			/db-dashboards/api/v1/postgres/columns [get]
func (h *Handler) GetColumnsFromTable(rw http.ResponseWriter, req *http.Request) {
	connStr, err := handlerutils.GetStringHeaderByKey(req, "connection-string")
	if err != nil {
		msg := "no connection string header provided"

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		msg := fmt.Sprintf("cannot connect to db with conn str: %v", connStr)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	repo := postgresrepo.New(sqlx.NewDb(conn, "postgres"))

	tableName, err := handlerutils.GetStringHeaderByKey(req, "table-name")
	if err != nil {
		msg := "no connection string header provided"

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	columns, err := h.Service.GetColumnsFromTable(req.Context(), repo, tableName)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch columns from db")

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	render.JSON(rw, req, sliceutils.Map(columns, mapper.MapColumnToColumnResponse))
	rw.WriteHeader(http.StatusOK)
}

// GetAllRowsFromTable godoc
//
//		@Summary		Get all data from table
//		@Description	Get all data from table
//		@Security		JWT
//		@Tags			Postgres
//	 	@Param 			connection-string 	header 	string true "connection string"
//	 	@Param 			table-name 	header 	string true "name of the table"
//		@Produce		json
//		@Success		200	{object}	[]response.GetColumnsResponse
//		@Failure		401	{string}	Unauthorized
//		@Router			/db-dashboards/api/v1/postgres/data [get]
func (h *Handler) GetAllRowsFromTable(rw http.ResponseWriter, req *http.Request) {
	connStr, err := handlerutils.GetStringHeaderByKey(req, "connection-string")
	if err != nil {
		msg := "no connection string header provided"

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		msg := fmt.Sprintf("cannot connect to db with conn str: %v", connStr)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	repo := postgresrepo.New(sqlx.NewDb(conn, "postgres"))

	tableName, err := handlerutils.GetStringHeaderByKey(req, "table-name")
	if err != nil {
		msg := "no connection string header provided"

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	rows, err := h.Service.GetAllRowsFromTable(req.Context(), repo, tableName)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch columns from db")

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	bytes, err := json.Marshal(rows)
	if err != nil {
		msg := fmt.Sprintf("cannot marshall rows to json")

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if _, err = rw.Write(bytes); err != nil {
		return
	}

	rw.WriteHeader(http.StatusOK)
}
