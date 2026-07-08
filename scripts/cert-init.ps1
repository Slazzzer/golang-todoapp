# Obtain initial Let's Encrypt certificate (webroot via nginx) (Windows).

. "$PSScriptRoot/_common.ps1"

if ([string]::IsNullOrWhiteSpace($env:NGINX_DOMAIN)) {
    Write-Error "Set NGINX_DOMAIN in .env"
    exit 1
}
if ([string]::IsNullOrWhiteSpace($env:CERTBOT_EMAIL)) {
    Write-Error "Set CERTBOT_EMAIL in .env"
    exit 1
}

Write-Host "Ensuring nginx is running on port 80 for ACME challenge..."
& "$PSScriptRoot/compose-nginx.ps1" up -d nginx todoapp
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Requesting certificate for $($env:NGINX_DOMAIN)..."
& "$PSScriptRoot/compose-nginx.ps1" --profile certbot run --rm --entrypoint certbot certbot certonly `
    --webroot `
    -w /var/www/certbot `
    -d $env:NGINX_DOMAIN `
    --email $env:CERTBOT_EMAIL `
    --agree-tos `
    --no-eff-email
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Reloading nginx with HTTPS configuration..."
& "$PSScriptRoot/compose-nginx.ps1" restart nginx
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host ""
Write-Host "Certificate issued. Test:"
Write-Host "  curl -s https://$($env:NGINX_DOMAIN)/health"
