package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"db-dashboards/internal/config"
	"db-dashboards/pkg/router"

	authhandler "db-dashboards/internal/handler/auth"
	postgreshandler "db-dashboards/internal/handler/postgres"
	userhandler "db-dashboards/internal/handler/user"

	userrepo "db-dashboards/internal/repository/user"

	authservice "db-dashboards/internal/service/auth"
	postgreservice "db-dashboards/internal/service/postgres"
	userservice "db-dashboards/internal/service/user"

	middlewares "db-dashboards/internal/handler/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "db-dashboards/docs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Hasher struct{}

func (h *Hasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (h *Hasher) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

const (
	configPath = "config/"
)

func initConfig() (*config.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	// env variables
	if err := godotenv.Load(configPath + "/.env"); err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("db_dashboards")
	viper.AutomaticEnv()

	// validate todo: VALIDATOR!

	conf.Jwt.Secret = viper.GetString("JWT_SECRET")
	if conf.Jwt.Secret == "" {
		return nil, errors.New("CHAT_JWT_SECRET env variable not set")
	}

	return &conf, nil
}

func main() {
	logger := logrus.New()
	ctx, cancel := context.WithCancel(context.Background())

	valid := validator.New(validator.WithRequiredStructEnabled())

	//tempConnStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	//
	//conn, err := sql.Open("pgx", tempConnStr)
	//if err != nil {
	//	logger.Fatalf("cannot open database connection with connection string: %v, err: %v", tempConnStr, err)
	//}
	//
	//db := sqlx.NewDb(conn, "postgres")
	//
	//repo := postgres.New(db)
	//
	//tables, err := repo.GetAllTables(ctx)
	//if err != nil {
	//	logger.Fatalf(err.Error())
	//}
	//
	//for _, table := range tables {
	//	logger.Infof("table: %v", table.Name)
	//
	//	columns, err := repo.GetColumnsFromTable(ctx, table.Name)
	//	if err != nil {
	//		logger.Fatalf(err.Error())
	//	}
	//
	//	for _, col := range columns {
	//		logger.Infof("column %v of type %v", col.Name, col.Type)
	//	}
	//
	//	rows, err := repo.GetAllRowsFromTable(ctx, table.Name)
	//	if err != nil {
	//		logger.Fatalf(err.Error())
	//	}
	//
	//	for _, row := range rows {
	//		logger.Infof("row: %+v", row)
	//	}
	//}

	conf, err := initConfig()
	if err != nil {
		logger.Fatalf("error occurred reading config file: %v", err)
	}

	connStr := conf.Postgres.ConnectionURL()

	conn, err := sql.Open("pgx", conf.Postgres.ConnectionURL())
	if err != nil {
		logger.Fatalf("cannot open database connection with connection string: %v, err: %v", connStr, err)
	}

	db := sqlx.NewDb(conn, "postgres")

	// try to connect to db
	for i := 0; i < conf.Postgres.Retries; i++ {
		conn, err = sql.Open("pgx", conf.Postgres.ConnectionURL())
		if err != nil {
			logger.Fatalf("cannot open database connection with connection string: %v, err: %v", conf.Postgres.ConnectionURL(), err)
		} else {
			db = sqlx.NewDb(conn, "postgres")

			if err = db.Ping(); err != nil {
				logger.Errorf("can't ping database: %v\nconnection string: %v", err, conf.Postgres.ConnectionURL())
				logger.Infof("retrying in %v sec...", conf.Postgres.Interval)
				logger.Infof("retry %v of %v", i+1, conf.Postgres.Retries)

				time.Sleep(time.Duration(conf.Postgres.Interval) * time.Second)
			} else {
				err = nil
				break
			}
		}
	}

	userRepo := userrepo.New(db)

	userService := userservice.New(userRepo, &Hasher{})
	authService := authservice.New(userRepo, &Hasher{})
	postgresService := postgreservice.New()

	authMiddleware := middlewares.JWTAuthMiddleware(conf.Jwt.Secret, logger)

	authHandler := authhandler.New(userService, authService, conf.Jwt, logger, valid)
	userHandler := userhandler.New(userService, logger, valid, authMiddleware)
	postgresHandler := postgreshandler.New(postgresService, logger, valid)

	routers := make(map[string]chi.Router)

	routers["/auth"] = authHandler.Routes()
	routers["/users"] = userHandler.Routes()
	routers["/postgres"] = postgresHandler.Routes()

	middlewars := []router.Middleware{
		middleware.Recoverer,
		middleware.Logger,
	}

	r := router.MakeRoutes("/db-dashboards/api/v1", routers, middlewars)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Server.Port),
		Handler: r,
	}

	// add swagger middleware
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", conf.Server.Port)), // The url pointing to API definition
	))

	logger.Infof("server started at port %v", server.Addr)

	go func() {
		if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			logger.WithError(serverErr).Fatalf("server can't listen requests")
		}
	}()

	logger.Infof("documentation available on: http://localhost:%v/swagger/index.html", conf.Server.Port)

	interrupt := make(chan os.Signal, 1)

	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(interrupt, syscall.SIGINT)

	go func() {
		<-interrupt

		logger.Info("interrupt signal caught")
		logger.Info("server shutting down")

		if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
			logger.WithError(shutdownErr).Fatalf("can't close server listening on '%s'", server.Addr)
		}

		cancel()
	}()

	<-ctx.Done()
}
