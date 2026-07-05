# Общие проверки для PowerShell-скриптов Makefile.
# Подключается через: . "$PSScriptRoot/_common.ps1"

if ([string]::IsNullOrWhiteSpace($env:PROJECT_ROOT)) {
    Write-Error 'PROJECT_ROOT is not set (expected from Makefile)'
    exit 1
}

Set-Location $env:PROJECT_ROOT
