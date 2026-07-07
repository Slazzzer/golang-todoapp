#!/usr/bin/env bash
# Полный деплой: Postgres → миграции → todoapp (Linux/macOS).
# Вызывается из Makefile: make deploy-up

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

"${MAKE:-make}" -C "$PROJECT_ROOT" env-up
"${MAKE:-make}" -C "$PROJECT_ROOT" migrate-up
"${MAKE:-make}" -C "$PROJECT_ROOT" todoapp-deploy
