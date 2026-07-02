#!/usr/bin/env bash
# Применение или откат миграций через golang-migrate (Linux/macOS).
# Вызывается из Makefile: make migrate-action action=<команда>
#
# Типичные значения action: up, down
# Переменные POSTGRES_* берутся из .env (Make экспортирует их в окружение).
# Хост todoapp-postgres — имя сервиса в docker-compose сети.

set -euo pipefail

action="${1:-}"
if [[ -z "$action" ]]; then
    echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up" >&2
    exit 1
fi

database="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable"

docker compose run --rm todoapp-postgres-migrate \
    -path /migrations \
    -database "$database" \
    "$action"
