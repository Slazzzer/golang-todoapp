#!/usr/bin/env bash
# Создание новой SQL-миграции через golang-migrate (Linux/macOS).
# Вызывается из Makefile: make migrate-create seq=<имя>
#
# Параметр seq — суффикс файлов, например init → 000001_init.up.sql
# Образ todoapp-postgres-migrate монтирует ./migrations в /migrations.

set -euo pipefail

seq="${1:-}"
if [[ -z "$seq" ]]; then
    echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init" >&2
    exit 1
fi

mkdir -p migrations

docker compose run --rm todoapp-postgres-migrate \
    create \
    -ext sql \
    -dir /migrations \
    -seq "$seq"
