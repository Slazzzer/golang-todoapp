#!/usr/bin/env bash
# Manual certificate renewal + nginx reload.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

COMPOSE_NGINX="$(dirname "${BASH_SOURCE[0]}")/compose-nginx.sh"

"$COMPOSE_NGINX" --profile certbot run --rm --entrypoint certbot certbot renew --webroot -w /var/www/certbot
"$COMPOSE_NGINX" restart nginx

echo "Certificate renewal check completed."
