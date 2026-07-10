package main

// @title           Todoapp API
// @version         1.0
// @description     REST API для todo-приложения (задачи, пользователи, статистика).
// @description     Разграничение доступа без пароля: заголовок X-User-ID указывает, от чьего имени выполняется запрос.
// @description     Публично без заголовка: GET/POST /users (выбор и регистрация). Остальные эндпоинты требуют X-User-ID.
// @host            localhost:5050
// @BasePath        /api/v1

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Slazzzer/golang-todoapp/internal/core/config"
	core_adminauth "github.com/Slazzzer/golang-todoapp/internal/core/adminauth"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
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
	web_fs_repository "github.com/Slazzzer/golang-todoapp/internal/features/web/repository/file_system"
	web_service "github.com/Slazzzer/golang-todoapp/internal/features/web/service"
	web_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/web/transport/http"

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

	logger.Debug("Initializing feature", core_logger.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("Initializing feature", core_logger.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("Initializing feature", core_logger.String("feature", "web"))
	webRepository := web_fs_repository.NewWebFSRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("Initializing HTTP server...")

	httpConfig := core_http_server.NewConfigMust()
	corsConfig := core_http_middleware.NewCORSConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		// Порядок: CORS → RequestID → Logger → LimitBody → Trace → Panic → handler
		core_http_middleware.CORS(corsConfig),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.LimitBody(httpConfig.MaxBodyBytes),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	adminCfg := core_adminauth.NewConfig()

	logger.Debug("Initializing feature", core_logger.String("feature", "auth"))
	authService := auth_service.NewAuthService(adminCfg)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	core_http_probes.NewHandler(pool).Register(httpServer.Mux())

	apiV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersionV1)
	apiV1.Use(core_http_middleware.ActingUser(adminCfg))
	apiV1.RegisterRoutes(authTransportHTTP.Routes()...)
	apiV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiV1)

	httpServer.RegisterWebRouters(webTransportHTTP.Routes()...)
	httpServer.RegisterSwaggerRouter()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP transport run error: ", core_logger.Error(err))
	}
}
