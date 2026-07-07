# Build and start todoapp container (Windows).
# Called from Makefile: make todoapp-deploy

. "$PSScriptRoot/_common.ps1"

docker compose --project-directory $env:PROJECT_ROOT up -d --build todoapp
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
