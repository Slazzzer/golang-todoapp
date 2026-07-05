#!/usr/bin/env bash
# Общие проверки для bash-скриптов Makefile.
# Подключается через: source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

: "${PROJECT_ROOT:?PROJECT_ROOT is not set (expected from Makefile)}"

cd "$PROJECT_ROOT"
