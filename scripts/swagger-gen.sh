#!/usr/bin/env bash
# Generate docs/docs.go, docs/swagger.json and docs/swagger.yaml via swag.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

bash "$(dirname "${BASH_SOURCE[0]}")/swagger-pull.sh"

docker compose --project-directory "$PROJECT_ROOT" run --rm swagger \
    init \
    -g cmd/todoapp/main.go \
    -o docs \
    --parseInternal \
    --parseDependency

echo "Swagger docs generated in docs/"
