package main

// @title           Todoapp API
// @version         1.0
// @description     REST API для todo-приложения с JWT-авторизацией.
// @description     Публичный маршрут: POST /auth/register. Все остальные эндпоинты /api/v1 требуют заголовок Authorization: Bearer &lt;token&gt;.
// @description     Пользователь видит и изменяет только свои данные (профиль, задачи, статистику).
// @host            localhost:5050
// @BasePath        /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT-токен. Формат: Bearer &lt;token&gt;. Токен выдаётся при регистрации (POST /auth/register).

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
	core_config "github.com/Slazzzer/golang-todoapp/internal/core/config"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_ratelimit "github.com/Slazzzer/golang-todoapp/internal/core/ratelimit"
	core_postgres_pool_pgx "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/middleware"
	core_http_probes "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/probes"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
	auth_service "github.com/Slazzzer/golang-todoapp/internal/features/auth/service"
	auth_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/auth/transport/http"
	statistics_postgres_repository "github.com/Slazzzer/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Slazzzer/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/Slazzzer/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Slazzzer/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Slazzzer/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Slazzzer/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/users/transport/http"

	_ "github.com/Slazzzer/golang-todoapp/docs"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.Timezone

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

	logger.Debug("application timezone", core_logger.String("timezone", cfg.Timezone.String()))

	logger.Debug("Initializing postgres connection pool...")
	pool, err := core_postgres_pool_pgx.NewConnectionPool(
		ctx,
		core_postgres_pool_pgx.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool: ", core_logger.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing feature", core_logger.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	tokenManager := core_auth.NewTokenManager(core_auth.NewConfigMust())
	rateLimitConfig := core_ratelimit.NewConfigMust()
	registerRateLimiter := core_ratelimit.NewLimiter(
		rateLimitConfig.RegisterMax,
		rateLimitConfig.RegisterWindow,
	)
	authService := auth_service.NewAuthService(usersService, tokenManager)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(
		authService,
		core_http_middleware.RateLimit(registerRateLimiter),
	)

	logger.Debug("Initializing feature", core_logger.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("Initializing feature", core_logger.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("Initializing HTTP server...")

	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		// Порядок: CORS → RequestID → Logger → LimitBody → Trace → Panic → handler
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.LimitBody(httpConfig.MaxBodyBytes),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	core_http_probes.NewHandler(pool).Register(httpServer.Mux())

	apiV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersionV1)
	apiV1.Use(core_http_middleware.JWTAuth(tokenManager))
	apiV1.RegisterRoutes(authTransportHTTP.Routes()...)
	apiV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiV1)

	httpServer.RegisterSwaggerRouter()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP transport run error: ", core_logger.Error(err))
	}
}
