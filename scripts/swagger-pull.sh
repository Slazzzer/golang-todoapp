#!/usr/bin/env bash
# Pull golang base image and build local swagger generator image.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

docker pull golang:1.26.3-bookworm
docker compose --project-directory "$PROJECT_ROOT" build swagger
