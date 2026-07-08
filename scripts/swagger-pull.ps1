# Pull golang base image and build local swagger generator image (Windows).

. "$PSScriptRoot/_common.ps1"

docker pull golang:1.26.3-bookworm
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

docker compose --project-directory $env:PROJECT_ROOT build swagger
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
