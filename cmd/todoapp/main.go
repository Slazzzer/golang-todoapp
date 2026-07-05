package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
	users_postgres_repository "github.com/Slazzzer/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Slazzzer/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/users/transport/http"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Debug("Initializing postgres connection pool...")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool: ", core_logger.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing feature", core_logger.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing HTTP server...")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		// Порядок: RequestID → Logger → Dummy → Trace → Panic → handler
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Dummy(),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	usersRoutes := usersTransportHTTP.Routes()

	for _, version := range []core_http_server.APIVersion{
		core_http_server.APIVersionV1,
		core_http_server.APIVersionV2,
		core_http_server.APIVersionV3,
	} {
		router := core_http_server.NewAPIVersionRouter(version)
		router.RegisterRoutes(usersRoutes...)
		httpServer.RegisterAPIRouters(router)
	}

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP transport run error: ", core_logger.Error(err))
	}
}
