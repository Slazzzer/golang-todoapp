# Запуск Go-приложения локально (Windows).
# Вызывается из Makefile: make todoapp-run
#
# LOGGER_FOLDER — путь зависит от PROJECT_ROOT (экспортируется Make).
# POSTGRES_HOST — хост БД при локальном запуске (приложение с хоста, Postgres в Docker).
# Остальные переменные (LOGGER_LEVEL, HTTP_*, POSTGRES_USER/...) — из .env.

$env:LOGGER_FOLDER = Join-Path $env:PROJECT_ROOT 'out/logs'
$env:POSTGRES_HOST = 'localhost'
if ($env:POSTGRES_HOST_PORT) {
    $env:POSTGRES_PORT = $env:POSTGRES_HOST_PORT
}

go mod tidy
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go run cmd/todoapp/main.go
