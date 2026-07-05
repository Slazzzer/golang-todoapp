#!/usr/bin/env bash
# Запуск Go-приложения локально (Linux/macOS).
# Вызывается из Makefile: make todoapp-run
#
# LOGGER_FOLDER — путь зависит от PROJECT_ROOT (экспортируется Make).
# POSTGRES_HOST — хост БД при локальном запуске (приложение с хоста, Postgres в Docker).
# Остальные переменные (LOGGER_LEVEL, HTTP_*, POSTGRES_USER/...) — из .env.

set -euo pipefail

export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs"
export POSTGRES_HOST=localhost
if [[ -n "${POSTGRES_HOST_PORT:-}" ]]; then
    export POSTGRES_PORT="${POSTGRES_HOST_PORT}"
fi

go mod tidy
go run cmd/todoapp/main.go
