.PHONY: help build up setup down logs clean

COMPOSE := docker compose -f docker-compose.yml

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the chainlink Docker image from local source (optional, default uses published image)
	$(COMPOSE) build node-1

up: ## Start all services (postgres + anvil + 4 nodes)
	$(COMPOSE) up -d
	@echo "Waiting for nodes to become healthy..."
	@for i in 1 2 3 4; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo "All nodes healthy."

setup: ## Run the setup script (deploy contracts, create jobs)
	cd ../../../core/scripts && go run ./localdon

down: ## Stop all services
	$(COMPOSE) down

logs: ## Follow logs from all services
	$(COMPOSE) logs -f

clean: ## Stop all services and remove volumes
	$(COMPOSE) down -v --remove-orphans
