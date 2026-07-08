#!/usr/bin/env bash
# Obtain initial Let's Encrypt certificate (webroot via nginx).

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

: "${NGINX_DOMAIN:?Set NGINX_DOMAIN in .env}"
: "${CERTBOT_EMAIL:?Set CERTBOT_EMAIL in .env}"

COMPOSE_NGINX="$(dirname "${BASH_SOURCE[0]}")/compose-nginx.sh"

echo "Ensuring nginx is running on port 80 for ACME challenge..."
"$COMPOSE_NGINX" up -d nginx todoapp

echo "Requesting certificate for ${NGINX_DOMAIN}..."
"$COMPOSE_NGINX" --profile certbot run --rm --entrypoint certbot certbot certonly \
    --webroot \
    -w /var/www/certbot \
    -d "$NGINX_DOMAIN" \
    --email "$CERTBOT_EMAIL" \
    --agree-tos \
    --no-eff-email

echo "Reloading nginx with HTTPS configuration..."
"$COMPOSE_NGINX" restart nginx

echo ""
echo "Certificate issued. Test:"
echo "  curl -s https://${NGINX_DOMAIN}/health"
echo "  curl -s https://${NGINX_DOMAIN}/ready"
