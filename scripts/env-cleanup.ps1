# Очистка локального окружения (Windows).
# Вызывается из Makefile: make env-cleanup
#
# 1. Спрашивает подтверждение у пользователя.
# 2. Останавливает docker compose проект.
# 3. Удаляет каталог out/pgdata с данными Postgres.

. "$PSScriptRoot/_common.ps1"

$consoleEncoding = [Text.Encoding]::GetEncoding(
    [int][System.Globalization.CultureInfo]::CurrentCulture.TextInfo.OEMCodePage
)
[Console]::InputEncoding = [Console]::OutputEncoding = $consoleEncoding

$pgdata = Join-Path $env:PROJECT_ROOT 'out/pgdata'

$ans = Read-Host 'Очистить все volume-файлы окружения? Опасность потери данных! [y/N]'
if ($ans -match '^[yY]') {
    docker compose --project-directory $env:PROJECT_ROOT down
    if (Test-Path $pgdata) {
        Remove-Item -Recurse -Force $pgdata
    }
    Write-Host 'Файлы окружения успешно очищены!'
} else {
    Write-Host 'Очистка окружения отменена!'
}
