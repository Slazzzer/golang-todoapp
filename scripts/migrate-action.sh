#!/usr/bin/env bash
# Применение или откат миграций через golang-migrate (Linux/macOS).
# Вызывается из Makefile: make migrate-action action=<команда>

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

action="${1:-}"
if [[ -z "$action" ]]; then
    echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up" >&2
    exit 1
fi

ssl_mode="${POSTGRES_SSLMODE:-disable}"

database="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=${ssl_mode}"

bash "$(dirname "${BASH_SOURCE[0]}")/wait-for-postgres.sh"

docker compose run --rm todoapp-postgres-migrate \
    -path /migrations \
    -database "$database" \
    "$action"
