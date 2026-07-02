# Применение или откат миграций через golang-migrate (Windows).
# Вызывается из Makefile: make migrate-action action=<команда>
#
# Типичные значения action: up, down
# Переменные POSTGRES_* берутся из .env (Make экспортирует их в окружение).
# Хост todoapp-postgres — имя сервиса в docker-compose сети.

param (
    [string]$action
)

if ([string]::IsNullOrWhiteSpace($action)) {
    Write-Error "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"
    exit 1
}

$database = "postgres://$($env:POSTGRES_USER):$($env:POSTGRES_PASSWORD)@todoapp-postgres:5432/$($env:POSTGRES_DB)?sslmode=disable"

docker compose run --rm todoapp-postgres-migrate `
    -path /migrations `
    -database $database `
    $action
