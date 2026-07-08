#!/usr/bin/env bash
# Compose with nginx overlay (Linux/macOS).

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

docker compose \
    --project-directory "$PROJECT_ROOT" \
    -f docker-compose.yaml \
    -f docker-compose.nginx.yaml \
    "$@"
