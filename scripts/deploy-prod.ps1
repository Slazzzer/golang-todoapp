# Production deploy: Postgres -> migrations -> todoapp + nginx + certbot (Windows).

. "$PSScriptRoot/_common.ps1"

& make -C $env:PROJECT_ROOT env-up
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

& make -C $env:PROJECT_ROOT migrate-up
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

& "$PSScriptRoot/compose-nginx.ps1" up -d --build todoapp nginx
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host ""
Write-Host "Production stack is up."
if ([string]::IsNullOrWhiteSpace($env:NGINX_DOMAIN)) {
    Write-Host "Set NGINX_DOMAIN and CERTBOT_EMAIL in .env, then run: make cert-init"
} else {
    Write-Host "  HTTP:  http://$($env:NGINX_DOMAIN)"
    Write-Host "HTTPS: run 'make cert-init' if certificate is not issued yet."
}
