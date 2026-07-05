#!/usr/bin/env bash
# Очистка локального окружения (Linux/macOS).
# Вызывается из Makefile: make env-cleanup
#
# 1. Спрашивает подтверждение у пользователя.
# 2. Останавливает docker compose проект.
# 3. Удаляет каталог out/pgdata с данными Postgres.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

pgdata="${PROJECT_ROOT}/out/pgdata"

read -r -p 'Очистить все volume-файлы окружения? Опасность потери данных! [y/N] ' ans
if [[ "$ans" =~ ^[yY]$ ]]; then
    docker compose down
    rm -rf "$pgdata"
    echo 'Файлы окружения успешно очищены!'
else
    echo 'Очистка окружения отменена!'
fi
