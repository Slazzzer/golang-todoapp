# Generate docs/docs.go, docs/swagger.json and docs/swagger.yaml via swag (Windows).

. "$PSScriptRoot/_common.ps1"

& "$PSScriptRoot/swagger-pull.ps1"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

docker compose --project-directory $env:PROJECT_ROOT run --rm swagger `
    init `
    -g cmd/todoapp/main.go `
    -o docs `
    --parseInternal `
    --parseDependency
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Swagger docs generated in docs/"
