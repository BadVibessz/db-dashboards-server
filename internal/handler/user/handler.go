package user

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"db-dashboards/internal/domain/entity"
)

type Service interface {
	RegisterUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) ([]*entity.User, error)
	DeleteUser(ctx context.Context, id int) (*entity.User, error)
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
		// r.Get("/all", h.GetAll)
	})

	return router
}

// todo:
//func (h *Handler) GetAll(rw http.ResponseWriter, req *http.Request) {
//	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)
//
//	if err := paginationOpts.Validate(h.validator); err != nil {
//		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())
//
//		return
//	}
//
//	users := h.UserService.GetAllUsers(req.Context(), paginationOpts.Offset, paginationOpts.Limit)
//
//	render.JSON(rw, req, sliceutils.Map(users, mapper.MapUserToUserResponse))
//	rw.WriteHeader(http.StatusOK)
//}
