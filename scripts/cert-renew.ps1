# Manual certificate renewal + nginx reload (Windows).

. "$PSScriptRoot/_common.ps1"

& "$PSScriptRoot/compose-nginx.ps1" --profile certbot run --rm --entrypoint certbot certbot renew --webroot -w /var/www/certbot
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

& "$PSScriptRoot/compose-nginx.ps1" restart nginx
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Certificate renewal check completed."
