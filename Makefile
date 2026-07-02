# =============================================================================
# Makefile — локальное окружение PostgreSQL и SQL-миграции
# =============================================================================
#
# Требования:
#   - Docker + Docker Compose
#   - GNU Make
#   - Windows: встроенный PowerShell
#   - Linux/macOS: bash
#
# Переменные из .env подключаются и экспортируются в дочерние процессы
# (docker compose, скрипты migrate-*).
#
# Примеры:
#   make env-up                      # поднять Postgres
#   make migrate-create seq=init     # создать пару файлов миграции
#   make migrate-up                  # применить миграции
#   make migrate-down                # откатить последнюю миграцию
#   make env-cleanup                 # остановить compose и удалить out/pgdata
#   make env-port-forward            # запустить port-forwarder для доступа к Postgres из контейнера в локальную сеть
#   make env-port-close              # остановить port-forwarder и закрыть доступ к Postgres из контейнера в локальную сеть
#
# На Windows логика в scripts/*.ps1, на Unix — в scripts/*.sh.
# Makefile только выбирает нужный раннер и передаёт параметры.
# =============================================================================

include .env
export

.PHONY: env-up env-down env-cleanup \
	env-port-forwarder env-port-close \
	migrate-create migrate-up migrate-down migrate-action

# --- Кросс-платформенные настройки ---

# PROJECT_ROOT — абсолютный путь к корню репозитория.
# Используется в docker-compose.yaml для bind-mount:
#   out/pgdata   -> данные Postgres
#   migrations/  -> SQL-файлы миграций
ifeq ($(OS),Windows_NT)
export PROJECT_ROOT := $(CURDIR)
RUN_SCRIPT := powershell -NoProfile -ExecutionPolicy Bypass -File
else
export PROJECT_ROOT := $(shell pwd)
RUN_SCRIPT := bash
endif

# --- Окружение ---

# Поднять контейнер Postgres в фоне.
env-up:
	@docker compose up -d todoapp-postgres

# Остановить и удалить контейнеры compose-проекта.
env-down:
	@docker compose down

# Интерактивная очистка: остановка compose + удаление каталога out/pgdata.
env-cleanup:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) scripts/env-cleanup.ps1
else
	@$(RUN_SCRIPT) scripts/env-cleanup.sh
endif

# Проброс Postgres на хост
env-port-forward:
	@docker compose up -d port-forwarder

# Остановить port-forwarder и закрыть доступ к Postgres из контейнера в локальную сеть
env-port-close:
	@docker compose down port-forwarder

# --- Миграции (образ migrate/migrate) ---

# Создать новую sequential-миграцию: migrations/NNNNNN_<seq>.{up,down}.sql
# Параметр seq обязателен: make migrate-create seq=add_users
migrate-create:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) scripts/migrate-create.ps1 -seq "$(seq)"
else
	@$(RUN_SCRIPT) scripts/migrate-create.sh "$(seq)"
endif

# Применить все неприменённые миграции.
migrate-up:
	@$(MAKE) migrate-action action=up

# Откатить одну последнюю миграцию.
migrate-down:
	@$(MAKE) migrate-action action=down

# Низкоуровневый таргет: передать action=up|down|force VERSION и т.д.
# Пример: make migrate-action action=force VERSION=1
migrate-action:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) scripts/migrate-action.ps1 -action "$(action)"
else
	@$(RUN_SCRIPT) scripts/migrate-action.sh "$(action)"
endif
