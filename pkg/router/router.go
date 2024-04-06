package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type (
	Middleware = func(next http.Handler) http.Handler
	Routers    = map[string]chi.Router
)

func MakeRoutes(basePath string, routers Routers, middlewares []Middleware) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares...)

	for routerPath, router := range routers {
		r.Mount(fmt.Sprintf("%s%s", basePath, routerPath), router)
	}

	return r
}
