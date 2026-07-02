# Создание новой SQL-миграции через golang-migrate (Windows).
# Вызывается из Makefile: make migrate-create seq=<имя>
#
# Параметр seq — суффикс файлов, например init → 000001_init.up.sql
# Образ todoapp-postgres-migrate монтирует ./migrations в /migrations.

param (
    [string]$seq
)

if ([string]::IsNullOrWhiteSpace($seq)) {
    Write-Error "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"
    exit 1
}

if (-not (Test-Path 'migrations')) {
    New-Item -ItemType Directory -Path 'migrations' | Out-Null
}

docker compose run --rm todoapp-postgres-migrate `
    create `
    -ext sql `
    -dir /migrations `
    -seq "$seq"
