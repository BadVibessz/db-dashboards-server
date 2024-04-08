package middleware

import (
	"context"
	"db-dashboards/internal/domain/entity"
	handlerutils "db-dashboards/pkg/utils/handler"
	jwtutils "db-dashboards/pkg/utils/jwt"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler = func(http.Handler) http.Handler

type Service interface {
	Login(ctx context.Context, email, password string) (*entity.User, error)
}

func JWTAuthMiddleware(secret string, logger *logrus.Logger) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				msg := "authorization header is empty"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			token := authHeader[len("Bearer "):]

			payload, err := jwtutils.ValidateToken(token, secret)
			if err != nil {
				msg := fmt.Sprintf("error occurred validating token: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			// todo: to private func
			idAny, exists := payload["id"]
			if !exists {
				msg := "invalid payload: not contains id"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			id, ok := idAny.(float64)
			if !ok {
				msg := "cannot parse id from payload to float64"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			emailAny, exists := payload["email"]
			if !exists {
				msg := "invalid payload: not contains email"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			email, ok := emailAny.(string)
			if !ok {
				msg := "cannot parse email from payload to string"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			req.Header.Set("id", strconv.Itoa(int(id)))
			req.Header.Set("email", email)

			next.ServeHTTP(rw, req)
		})
	}
}
