#!/usr/bin/env bash
# Удаление лог-файлов приложения (Linux/macOS).
# Вызывается из Makefile: make logs-cleanup

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

logs_dir="${PROJECT_ROOT}/out/logs"

if [[ ! -d "$logs_dir" ]]; then
    echo "Каталог логов не найден: $logs_dir"
    exit 0
fi

shopt -s nullglob
log_files=("$logs_dir"/*.log)
shopt -u nullglob

if [[ ${#log_files[@]} -eq 0 ]]; then
    echo "Лог-файлы не найдены в: $logs_dir"
    exit 0
fi

read -r -p "Удалить ${#log_files[@]} лог-файл(ов) из $logs_dir ? [y/N] " ans
if [[ "$ans" =~ ^[yY]$ ]]; then
    rm -f "${log_files[@]}"
    echo "Удалено лог-файлов: ${#log_files[@]}"
    echo "Каталог: $logs_dir"
else
    echo 'Очистка логов отменена.'
fi
