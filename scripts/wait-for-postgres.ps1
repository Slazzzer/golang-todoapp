# Ожидание готовности PostgreSQL перед миграциями и деплоем (Windows).
# Вызывается из scripts/migrate-action.ps1

. "$PSScriptRoot/_common.ps1"

$maxAttempts = 60
if ($env:POSTGRES_WAIT_ATTEMPTS) {
    $maxAttempts = [int]$env:POSTGRES_WAIT_ATTEMPTS
}

Write-Host 'Ожидание готовности PostgreSQL...'

for ($attempt = 0; $attempt -lt $maxAttempts; $attempt++) {
    docker compose exec -T todoapp-postgres `
        pg_isready -U $env:POSTGRES_USER -d $env:POSTGRES_DB 2>$null | Out-Null
    if ($LASTEXITCODE -eq 0) {
        Write-Host 'PostgreSQL готов.'
        exit 0
    }
    Start-Sleep -Seconds 1
}

Write-Error "PostgreSQL не ответил за $maxAttempts с."
exit 1
