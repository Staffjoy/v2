.PHONY: env-init env-check compose-config up down

env-init:
	@if [ -f .env ]; then \
		echo ".env already exists"; \
	else \
		cp .env.example .env && echo "Created .env from .env.example"; \
	fi

env-check:
	@./scripts/check-env.sh .env

compose-config:
	@docker compose --env-file .env config > /dev/null
	@echo "docker compose configuration is valid"

up: env-check
	@docker compose --env-file .env up -d

down:
	@docker compose --env-file .env down
