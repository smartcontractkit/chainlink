.PHONY: help build up setup setup-full down logs clean restart \
        _create-keys _restart-nodes _verify \
        local-dev-up local-dev-setup local-dev-setup-full local-dev-node

COMPOSE           := docker compose
COMPOSE_LOCAL_DEV := docker compose -f docker-compose.yml -f docker-compose.local-dev.yml
REPO_ROOT         := $(shell cd ../../.. && pwd)
LOCALDON_DIR      := $(REPO_ROOT)/tools/docker/localdon

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

build: ## Build the chainlink Docker image from local source (optional)
	$(COMPOSE) build node-1

# ── All-Docker (5 nodes) ────────────────────────────────────────────────────

up: ## Start all services (postgres + anvil + flashbots-mock + 5 nodes)
	$(COMPOSE) up -d
	@echo "Waiting for nodes to become healthy..."
	@for i in 1 2 3 4 5; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo "All nodes healthy."

setup: ## Deploy contracts, fund nodes, create jobs (all 5 Docker nodes)
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon 2>&1 | tee /tmp/localdon-setup.log

setup-full: up _create-keys _restart-nodes setup _verify ## Full setup: start, create keys, restart nodes, deploy+configure
	@echo "Full setup complete."

_create-keys:
	@echo "Creating secondary ETH keys on each node..."
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon --keys-only
	@echo "Verifying all nodes have 2 ETH keys..."
	@all_ok=true; \
	for port in 6688 6689 6690 6691 6692; do \
		curl -sf -c /tmp/cjar_$$port \
			-X POST http://localhost:$$port/sessions \
			-H "Content-Type: application/json" \
			-d '{"email":"local-test@chainlink.test","password":"localdon-testing-only-not-a-real-password"}' > /dev/null 2>&1; \
		count=$$(curl -sf -b /tmp/cjar_$$port http://localhost:$$port/v2/keys/eth 2>/dev/null \
			| /usr/bin/ruby -e "require 'json'; d=JSON.parse(STDIN.read); puts d['data'].length" 2>/dev/null); \
		if [ "$$count" = "2" ]; then \
			echo "  ✅ port $$port: 2 ETH keys"; \
		else \
			echo "  ❌ port $$port: expected 2 keys, got '$$count'"; \
			all_ok=false; \
		fi; \
	done; \
	if [ "$$all_ok" = "false" ]; then echo "Key creation failed — aborting"; exit 1; fi

_restart-nodes:
	@echo "Restarting nodes so TXMv2 picks up secondary keys..."
	$(COMPOSE) restart node-1 node-2 node-3 node-4 node-5
	@echo "Waiting for nodes to become healthy after restart..."
	@for i in 1 2 3 4 5; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo "All nodes healthy."

_verify: ## Wait for OCR rounds and verify primary + secondary transmissions
	@echo ""
	@echo "Verifying DON is producing rounds..."
	@echo "  Waiting for latestAnswer != 0 (up to 120s)..."
	@OCR_ADDR=$$(grep "OCR2Aggregator address:" /tmp/localdon-setup.log 2>/dev/null | awk '{print $$NF}'); \
	if [ -z "$$OCR_ADDR" ]; then \
		echo "  [warn] OCR address not found in /tmp/localdon-setup.log, skipping on-chain check"; \
	else \
		for i in $$(seq 1 60); do \
			val=$$(cast call $$OCR_ADDR 'latestAnswer()(int256)' --rpc-url http://localhost:8545 2>/dev/null | tr -d '[:space:]' | sed 's/\[.*//'); \
			if [ "$$val" != "0" ] && [ -n "$$val" ]; then \
				echo "  ✅ Primary: latestAnswer=$$val (rounds advancing)"; \
				break; \
			fi; \
			printf "."; sleep 2; \
		done; \
		if [ "$$val" = "0" ] || [ -z "$$val" ]; then echo "  ❌ Primary: latestAnswer still 0 after 120s — check node logs"; exit 1; fi; \
	fi
	@echo "  Checking flashbots mock secondary tx count..."
	@count1=$$(curl -s http://localhost:9090/health | sed 's/.*secondaryTxCount":\([0-9]*\).*/\1/'); \
	sleep 15; \
	count2=$$(curl -s http://localhost:9090/health | sed 's/.*secondaryTxCount":\([0-9]*\).*/\1/'); \
	if [ "$$count2" -gt "$$count1" ]; then \
		echo "  ✅ Secondary: secondaryTxCount $$count1 → $$count2 (increasing)"; \
	else \
		echo "  ❌ Secondary: secondaryTxCount not increasing ($$count1 → $$count2) — check node logs"; exit 1; \
	fi
	@echo ""

down: ## Stop all services
	$(COMPOSE) down

logs: ## Follow logs from all services
	$(COMPOSE) logs -f

clean: ## Stop all services and remove volumes
	$(COMPOSE) down -v --remove-orphans

restart: clean setup-full ## Full restart: clean + full setup (with key creation + node restart)

# ── Local-dev mode: node-5 on host (nodes 1–4 in Docker) ───────────────────
# Use this for Cases 2 and 3 in the INC-2328 test plan.
# Node-5 config lives in local-dev-node5/node-config.toml — edit DualBroadcast
# or TXMv2 Enabled there, then Ctrl-C + restart make local-dev-node.

local-dev-up: ## Start Docker infra + nodes 1–4 (node-5 disabled; P2P port exposed for host node-5)
	$(COMPOSE_LOCAL_DEV) up -d
	@echo "Waiting for Docker nodes to become healthy..."
	@for i in 1 2 3 4; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo "Docker nodes healthy. Start node-5 with: make local-dev-node"

local-dev-setup: ## Run the setup script in local-dev mode (node-5 on host)
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon --local-dev 2>&1 | tee /tmp/localdon-setup.log

local-dev-setup-full: local-dev-up ## Interactive full setup: Docker nodes 1–4, then host node-5
	@echo ""
	@echo ">>> Start node-5 locally in another terminal: make local-dev-node"
	@echo ">>> Waiting for node-5 to become healthy on localhost:6692..."
	@until curl -sf http://localhost:6692/readyz > /dev/null 2>&1; do \
		printf "."; sleep 2; \
	done
	@echo " node-5 ready"
	@echo ""
	@echo "Creating secondary ETH keys on all 5 nodes..."
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon --keys-only
	@echo ""
	@echo "Restarting Docker nodes 1–4 so TXMv2 picks up secondary keys..."
	$(COMPOSE_LOCAL_DEV) restart node-1 node-2 node-3 node-4
	@echo "Waiting for Docker nodes to become healthy after restart..."
	@for i in 1 2 3 4; do \
		printf "  node-$$i: "; \
		until curl -sf http://localhost:$$(( 6687 + $$i ))/readyz > /dev/null 2>&1; do \
			printf "."; sleep 2; \
		done; \
		echo " ready"; \
	done
	@echo ""
	@echo ">>> Now restart node-5: Ctrl-C the local-dev-node process, then re-run: make local-dev-node"
	@echo ">>> Waiting for node-5 to become healthy on localhost:6692..."
	@until curl -sf http://localhost:6692/readyz > /dev/null 2>&1; do \
		printf "."; sleep 2; \
	done
	@echo " node-5 ready"
	@echo ""
	cd $(REPO_ROOT)/core/scripts && go run ./setup-localdon --local-dev 2>&1 | tee /tmp/localdon-setup.log
	$(MAKE) _verify
	@echo "Local-dev full setup complete."

local-dev-node: ## Run node-5 locally from source (edit local-dev-node5/node-config.toml between test cases)
	cd $(REPO_ROOT) && go run -mod=mod . \
		-c $(LOCALDON_DIR)/local-dev-node5/node-config.toml \
		-s $(LOCALDON_DIR)/local-dev-node5/secrets.toml \
		node start -d \
		-p $(LOCALDON_DIR)/secrets/password.txt \
		-a $(LOCALDON_DIR)/secrets/apicredentials \
		2>&1 | tee $(LOCALDON_DIR)/local-dev-node5/node-5.log
