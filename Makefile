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
#   make env-down                    # остановить Postgres
#   make migrate-create seq=init     # создать пару файлов миграции
#   make migrate-up                  # применить миграции
#   make migrate-down                # откатить последнюю миграцию
#   make migrate-action action=up VERSION=1 # применить конкретную миграцию
#   make env-cleanup                 # остановить compose и удалить out/pgdata
#   make env-port-forward            # запустить port-forwarder для доступа к Postgres из контейнера в локальную сеть
#   make env-port-close              # остановить port-forwarder и закрыть доступ к Postgres из контейнера в локальную сеть
#   make todoapp-run                 # запустить Go-приложение локально
#   make logs-cleanup                # удалить все лог-файлы из out/logs
#   make todoapp-deploy              # собрать образ и поднять контейнер todoapp
#   make todoapp-undeploy            # остановить и удалить контейнер todoapp
#   make ps                          # статус контейнеров compose-проекта
#
# На Windows логика в scripts/*.ps1, на Unix — в scripts/*.sh.
# PROJECT_ROOT экспортируется Make и обязателен для всех скриптов.
# Makefile только выбирает нужный раннер и передаёт параметры.
# =============================================================================

include .env
export

.PHONY: env-up env-down env-cleanup \
	env-port-forward env-port-close \
	migrate-create migrate-up migrate-down migrate-action \
	todoapp-run logs-cleanup \
	todoapp-deploy todoapp-undeploy ps

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
	@docker compose --project-directory "$(PROJECT_ROOT)" up -d todoapp-postgres

# Остановить и удалить контейнеры compose-проекта.
env-down:
	@docker compose --project-directory "$(PROJECT_ROOT)" down

# Интерактивная очистка: остановка compose + удаление каталога out/pgdata.
env-cleanup:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/env-cleanup.ps1"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/env-cleanup.sh"
endif

# Проброс Postgres на хост
env-port-forward:
	@docker compose --project-directory "$(PROJECT_ROOT)" up -d port-forwarder

# Остановить только port-forwarder (Postgres и сеть compose не трогаем).
env-port-close:
	@docker compose --project-directory "$(PROJECT_ROOT)" rm -sf port-forwarder

# --- Миграции (образ migrate/migrate) ---

# Создать новую sequential-миграцию: migrations/NNNNNN_<seq>.{up,down}.sql
# Параметр seq обязателен: make migrate-create seq=add_users
migrate-create:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/migrate-create.ps1" -seq "$(seq)"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/migrate-create.sh" "$(seq)"
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
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/migrate-action.ps1" -action "$(action)"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/migrate-action.sh" "$(action)"
endif

# Запуск Go-приложения локально.
# POSTGRES_HOST=localhost задаётся в scripts/todoapp-run.{ps1,sh} — приложение
# подключается к Postgres на хосте (нужны make env-up + make env-port-forward).
todoapp-run:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/todoapp-run.ps1"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/todoapp-run.sh"
endif

# Удалить все *.log из out/logs (с подтверждением).
logs-cleanup:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/logs-cleanup.ps1"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/logs-cleanup.sh"
endif

# Собрать образ и поднять контейнер todoapp в Docker.
todoapp-deploy:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/todoapp-deploy.ps1"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/todoapp-deploy.sh"
endif

# Остановить и удалить только контейнер todoapp.
todoapp-undeploy:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/todoapp-undeploy.ps1"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/todoapp-undeploy.sh"
endif

# Статус контейнеров compose-проекта.
ps:
ifeq ($(OS),Windows_NT)
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/compose-ps.ps1"
else
	@$(RUN_SCRIPT) "$(PROJECT_ROOT)/scripts/compose-ps.sh"
endif