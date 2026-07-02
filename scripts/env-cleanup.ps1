# Очистка локального окружения (Windows).
# Вызывается из Makefile: make env-cleanup
#
# 1. Спрашивает подтверждение у пользователя.
# 2. Останавливает docker compose проект.
# 3. Удаляет каталог out/pgdata с данными Postgres.

# Кодировка консоли: OEM (на RU Windows — CP866), чтобы кириллица в Read-Host
# отображалась корректно. Файл сохранён в UTF-8 BOM.
$consoleEncoding = [Text.Encoding]::GetEncoding(
    [int][System.Globalization.CultureInfo]::CurrentCulture.TextInfo.OEMCodePage
)
[Console]::InputEncoding = [Console]::OutputEncoding = $consoleEncoding

$ans = Read-Host 'Очистить все volume-файлы окружения? Опасность потери данных! [y/N]'
if ($ans -match '^[yY]') {
    docker compose down
    if (Test-Path 'out/pgdata') {
        Remove-Item -Recurse -Force 'out/pgdata'
    }
    Write-Host 'Файлы окружения успешно очищены!'
} else {
    Write-Host 'Очистка окружения отменена!'
}
