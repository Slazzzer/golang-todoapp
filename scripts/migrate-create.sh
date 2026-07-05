#!/usr/bin/env bash
# Создание новой SQL-миграции через golang-migrate (Linux/macOS).
# Вызывается из Makefile: make migrate-create seq=<имя>

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

seq="${1:-}"
if [[ -z "$seq" ]]; then
    echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init" >&2
    exit 1
fi

mkdir -p "${PROJECT_ROOT}/migrations"

docker compose run --rm todoapp-postgres-migrate \
    create \
    -ext sql \
    -dir /migrations \
    -seq "$seq"
