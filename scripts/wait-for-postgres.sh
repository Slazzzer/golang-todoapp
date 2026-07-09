#!/usr/bin/env bash
# Ожидание готовности PostgreSQL перед миграциями и деплоем.
# Вызывается из scripts/migrate-action.sh

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

max_attempts="${POSTGRES_WAIT_ATTEMPTS:-60}"
attempt=0

echo "Ожидание готовности PostgreSQL..."

until docker compose exec -T todoapp-postgres \
    pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [[ $attempt -ge $max_attempts ]]; then
        echo "PostgreSQL не ответил за ${max_attempts} с." >&2
        exit 1
    fi
    sleep 1
done

echo "PostgreSQL готов."
