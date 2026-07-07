# Применение или откат миграций через golang-migrate (Windows).
# Вызывается из Makefile: make migrate-action action=<команда>

param (
    [string]$action
)

. "$PSScriptRoot/_common.ps1"

if ([string]::IsNullOrWhiteSpace($action)) {
    Write-Error "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"
    exit 1
}

$sslMode = $env:POSTGRES_SSLMODE
if ([string]::IsNullOrWhiteSpace($sslMode)) {
    $sslMode = 'disable'
}

$database = "postgres://$($env:POSTGRES_USER):$($env:POSTGRES_PASSWORD)@todoapp-postgres:5432/$($env:POSTGRES_DB)?sslmode=$sslMode"

docker compose run --rm todoapp-postgres-migrate `
    -path /migrations `
    -database $database `
    $action
