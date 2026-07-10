# Wait for PostgreSQL before migrations/deploy (Windows).
# Called from scripts/migrate-action.ps1

. "$PSScriptRoot/_common.ps1"

$maxAttempts = 60
if ($env:POSTGRES_WAIT_ATTEMPTS) {
    $maxAttempts = [int]$env:POSTGRES_WAIT_ATTEMPTS
}

Write-Host 'Waiting for PostgreSQL...'

for ($attempt = 0; $attempt -lt $maxAttempts; $attempt++) {
    docker compose exec -T todoapp-postgres `
        pg_isready -U $env:POSTGRES_USER -d $env:POSTGRES_DB 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host 'PostgreSQL is ready.'
        exit 0
    }
    Start-Sleep -Seconds 1
}

Write-Error "PostgreSQL did not become ready within $maxAttempts seconds."
exit 1
