.PHONY: help build up setup down logs clean restart local-dev-up local-dev-setup local-dev-node

COMPOSE := docker compose
COMPOSE_LOCAL_DEV := docker compose -f docker-compose.yml -f docker-compose.local-dev.yml
REPO_ROOT := $(shell cd ../../.. && pwd)
LOCALDON_DIR := $(REPO_ROOT)/tools/docker/localdon

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the chainlink Docker image from local source (optional)
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
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon

setup-full: up _create-keys _restart-nodes setup ## Full setup: start, create keys, restart nodes, deploy+configure
	@echo "Full setup complete."

_create-keys:
	@echo "Creating secondary ETH keys on each node..."
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon --keys-only

_restart-nodes:
	@echo "Restarting nodes so TXMv2 picks up secondary keys..."
	$(COMPOSE) restart node-1 node-2 node-3 node-4
	@echo "Waiting for nodes to become healthy after restart..."
	@for i in 1 2 3 4; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo "All nodes healthy."

down: ## Stop all services
	$(COMPOSE) down

logs: ## Follow logs from all services
	$(COMPOSE) logs -f

clean: ## Stop all services and remove volumes
	$(COMPOSE) down -v --remove-orphans

restart: clean setup-full ## Full restart: clean + full setup (with key creation + node restart)

# --- Local-dev mode: run node-4 on the host ---

local-dev-up: ## Start Docker services (postgres + anvil + 3 nodes; node-4 disabled)
	$(COMPOSE_LOCAL_DEV) up -d
	@echo "Waiting for Docker nodes to become healthy..."
	@for i in 1 2 3; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo "Docker nodes healthy. Start node-4 with: make local-dev-node"

local-dev-setup: ## Run the setup script in local-dev mode
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon --local-dev

local-dev-node: ## Run node-4 locally from source (logs to local-dev/node-4.log)
	@echo "Node-4 logs: $(LOCALDON_DIR)/local-dev/node-4.log"
	cd $(REPO_ROOT) && go run -mod=mod . \
		-c $(LOCALDON_DIR)/local-dev/node-config.toml \
		-s $(LOCALDON_DIR)/local-dev/secrets.toml \
		node start -d \
		-p $(LOCALDON_DIR)/secrets/password.txt \
		-a $(LOCALDON_DIR)/secrets/apicredentials \
		2>&1 | tee $(LOCALDON_DIR)/local-dev/node-4.log
