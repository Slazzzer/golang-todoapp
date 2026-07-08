# Compose with nginx overlay (Windows).
# Called from Makefile targets that need -f docker-compose.nginx.yaml

. "$PSScriptRoot/_common.ps1"

$argsList = @(
    "compose",
    "--project-directory", $env:PROJECT_ROOT,
    "-f", "docker-compose.yaml",
    "-f", "docker-compose.nginx.yaml"
) + $args

& docker @argsList
exit $LASTEXITCODE
