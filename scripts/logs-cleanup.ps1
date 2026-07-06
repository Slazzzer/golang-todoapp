# Удаление лог-файлов приложения (Windows).
# Вызывается из Makefile: make logs-cleanup
#
# 1. Спрашивает подтверждение у пользователя.
# 2. Удаляет все *.log из out/logs.
# 3. Показывает количество удалённых файлов и каталог.

. "$PSScriptRoot/_common.ps1"

$consoleEncoding = [Text.Encoding]::GetEncoding(
    [int][System.Globalization.CultureInfo]::CurrentCulture.TextInfo.OEMCodePage
)
[Console]::InputEncoding = [Console]::OutputEncoding = $consoleEncoding

$logsDir = Join-Path $env:PROJECT_ROOT 'out/logs'

if (-not (Test-Path $logsDir)) {
    Write-Host 'Каталог логов не найден, нечего удалять.'
    exit 0
}

$logFiles = @(Get-ChildItem -Path $logsDir -Filter '*.log' -File -ErrorAction SilentlyContinue)

if ($logFiles.Count -eq 0) {
    Write-Host 'Лог-файлы не найдены, нечего удалять.'
    exit 0
}

$ans = Read-Host 'Удалить все лог-файлы из out/logs? [y/N]'
if ($ans -match '^[yY]') {
    $logFiles | Remove-Item -Force
    Write-Host 'Файлы логов успешно очищены!'
    Write-Host 'Удалено лог-файлов:' $logFiles.Count
    Write-Host 'Каталог:' $logsDir
} else {
    Write-Host 'Очистка логов отменена!'
}
