-include .env

.PHONY: help migrate-up migrate-down migrate-status migrate-create

.PHONY: help
help: ## Показать эту справку [default]
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: migrate-up
migrate-up:  ## Применить все миграции
	docker-compose exec backend goose -dir ./migrations postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:  ## Откатить последнюю миграцию
	docker-compose exec backend goose -dir ./migrations postgres "$(DB_URL)" down

.PHONY: migrate-status
migrate-status:  ## Показать статус миграций
	docker-compose exec backend goose -dir ./migrations postgres "$(DB_URL)" status

.PHONY: migrate-create
migrate-create:  ## Создать новую миграцию (make migrate-create name=my_migration)
	docker-compose exec backend goose -dir ./migrations create "$(name)" sql