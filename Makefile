include .env
export

service-run:
	@go run main.go

# Docker Compose
service-deploy: ## Задеплоить сервис в Docker
	@docker compose up -d app

service-undeploy: ## Удалить сервис из Docker
	@docker compose down app

run-db: ## Запустить БД
	@docker compose up -d postgres

stop-db: ## Остановить БД
	@docker compose stop postgres

# Логи
logs: ## Показать все логи
	@docker compose logs -f

logs-app: ## Логи приложения
	@docker compose logs -f app

logs-db: ## Логи БД
	@docker compose logs -f postgres

# Сборка
build: ## Собрать образы
	@docker compose build

clean: ## Очистить всё
	@docker compose down -v --remove-orphans
	@rm -rf ./out/pgdata

# База данных
db-reset: ## Полный сброс БД
	@docker compose down -v
	@rm -rf ./out/pgdata
	@mkdir -p ./out/pgdata
	@echo "Database reset complete. Run 'make build' to recreate."

# Утилиты
ps: ## Показать контейнеры
	@docker compose ps

shell: ## Shell в контейнер app
	@docker compose exec app sh

db-shell: ## Shell в контейнер postgres
	@docker compose exec postgres sh

# Основные команды
up: ## Запустить всё (все сервисы)
	@docker compose up --build

down: ## Остановить всё
	@docker compose down


	