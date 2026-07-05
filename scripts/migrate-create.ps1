# Создание новой SQL-миграции через golang-migrate (Windows).
# Вызывается из Makefile: make migrate-create seq=<имя>

param (
    [string]$seq
)

. "$PSScriptRoot/_common.ps1"

if ([string]::IsNullOrWhiteSpace($seq)) {
    Write-Error "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"
    exit 1
}

$migrations = Join-Path $env:PROJECT_ROOT 'migrations'
if (-not (Test-Path $migrations)) {
    New-Item -ItemType Directory -Path $migrations | Out-Null
}

docker compose run --rm todoapp-postgres-migrate `
    create `
    -ext sql `
    -dir /migrations `
    -seq "$seq"
