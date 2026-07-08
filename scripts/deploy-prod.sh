#!/usr/bin/env bash
# Production deploy: Postgres → migrations → todoapp + nginx + certbot.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

"${MAKE:-make}" -C "$PROJECT_ROOT" env-up
"${MAKE:-make}" -C "$PROJECT_ROOT" migrate-up

"$(dirname "${BASH_SOURCE[0]}")/compose-nginx.sh" up -d --build todoapp nginx

echo ""
echo "Production stack is up."
if [ -z "${NGINX_DOMAIN:-}" ]; then
    echo "Set NGINX_DOMAIN and CERTBOT_EMAIL in .env, then run: make cert-init"
else
    echo "  HTTP:  http://${NGINX_DOMAIN}"
    echo "  Run 'make cert-init' to enable HTTPS (Let's Encrypt)."
fi
