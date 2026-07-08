#!/bin/sh
set -eu

: "${NGINX_DOMAIN:?NGINX_DOMAIN is required}"

CERT_PATH="/etc/letsencrypt/live/${NGINX_DOMAIN}/fullchain.pem"
CONF_DIR="/etc/nginx/conf.d"

rm -f "${CONF_DIR}"/*.conf

export NGINX_DOMAIN

if [ -f "$CERT_PATH" ]; then
    envsubst '${NGINX_DOMAIN}' < /etc/nginx/templates/https.conf.template > "${CONF_DIR}/default.conf"
else
    envsubst '${NGINX_DOMAIN}' < /etc/nginx/templates/http.conf.template > "${CONF_DIR}/default.conf"
fi
