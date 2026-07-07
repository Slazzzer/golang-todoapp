# Full deploy: Postgres -> migrations -> todoapp (Windows).
# Called from Makefile: make deploy-up

. "$PSScriptRoot/_common.ps1"

& make -C $env:PROJECT_ROOT env-up
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

& make -C $env:PROJECT_ROOT migrate-up
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

& make -C $env:PROJECT_ROOT todoapp-deploy
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
