# Show compose services status (Windows).
# Called from Makefile: make ps

. "$PSScriptRoot/_common.ps1"

docker compose --project-directory $env:PROJECT_ROOT ps
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
