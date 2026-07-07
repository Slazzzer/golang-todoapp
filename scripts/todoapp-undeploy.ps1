# Stop and remove todoapp container (Windows).
# Called from Makefile: make todoapp-undeploy

. "$PSScriptRoot/_common.ps1"

docker compose --project-directory $env:PROJECT_ROOT down todoapp
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
