# Запуск Go-приложения локально (Windows).
# Вызывается из Makefile: make todoapp-run

. "$PSScriptRoot/_common.ps1"

$env:LOGGER_FOLDER = Join-Path $env:PROJECT_ROOT 'out/logs'
$env:POSTGRES_HOST = 'localhost'
if ($env:POSTGRES_HOST_PORT) {
    $env:POSTGRES_PORT = $env:POSTGRES_HOST_PORT
}

$main = Join-Path $env:PROJECT_ROOT 'cmd/todoapp/main.go'

go mod tidy
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go run $main
