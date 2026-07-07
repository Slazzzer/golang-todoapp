#!/usr/bin/env bash
# Сборка и запуск контейнера todoapp (Linux/macOS).
# Вызывается из Makefile: make todoapp-deploy

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

docker compose --project-directory "$PROJECT_ROOT" up -d --build todoapp
