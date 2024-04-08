package postgres

import (
	"context"
	"database/sql"
	"db-dashboards/internal/domain/entity/postgres"
	"db-dashboards/internal/handler/mapper"
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
