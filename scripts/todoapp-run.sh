#!/usr/bin/env bash
# Запуск Go-приложения локально (Linux/macOS).
# Вызывается из Makefile: make todoapp-run

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs"
export POSTGRES_HOST=localhost
if [[ -n "${POSTGRES_HOST_PORT:-}" ]]; then
    export POSTGRES_PORT="${POSTGRES_HOST_PORT}"
fi

go mod tidy
go run "${PROJECT_ROOT}/cmd/todoapp/main.go"
