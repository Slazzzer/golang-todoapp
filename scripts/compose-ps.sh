#!/usr/bin/env bash
# Статус сервисов docker compose (Linux/macOS).
# Вызывается из Makefile: make ps

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

docker compose --project-directory "$PROJECT_ROOT" ps
