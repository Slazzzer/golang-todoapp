#!/usr/bin/env bash
# Остановка и удаление контейнера todoapp (Linux/macOS).
# Вызывается из Makefile: make todoapp-undeploy

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

docker compose --project-directory "$PROJECT_ROOT" down todoapp
