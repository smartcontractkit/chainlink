.PHONY: run clean fresh logs transmitter-logs test test-quick help

# The environment directory (parent of llo_e2e)
ENV_DIR := $(shell cd .. && pwd)

# Run the full LLO E2E test via CRE CLI (interactive mode)
run:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Starting LLO Streams Trigger E2E Test                       ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@cd $(ENV_DIR) && CRE_ENV_DIR=$(ENV_DIR) go run . llo-e2e

# Run the full E2E test using Go test framework with assertions
test:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Running LLO E2E Test (Go Test Framework)                    ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@cd $(ENV_DIR) && go test -v -timeout 30m ./llo_e2e -run TestLLOStreamsE2E

# Quick verification test (assumes environment is already running)
test-quick:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Quick E2E Verification (environment must be running)        ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@cd $(ENV_DIR) && go test -v -timeout 5m ./llo_e2e -run TestLLOStreamsE2E_QuickVerify

# Clean up all containers and start fresh
clean:
	@echo "Cleaning up Docker containers..."
	@docker ps -aq | xargs docker rm -f 2>/dev/null || true
	@docker ps -a --format "{{.Names}}" | xargs -I {} docker rm -f {} 2>/dev/null || true
	@docker rmi -f mock-ea-e2e mock-ea mock-ea:latest 2>/dev/null || true
	@docker network prune -f 2>/dev/null || true
	@docker volume prune -f 2>/dev/null || true
	@# Also clean up any CRE state files that might cause issues
	@rm -f $(ENV_DIR)/state/*.state.yaml 2>/dev/null || true
	@sleep 3
	@echo "Environment cleared"

# Clean and run interactive mode
fresh: clean run

# Clean and run test mode
fresh-test: clean test

# Show logs from workflow nodes (proof of E2E)
logs:
	@echo "=== WORKFLOW REPORT LOGS (proof of E2E) ==="
	@docker logs workflow-node0 2>&1 | grep "WORKFLOW_REPORT" | tail -10 || echo "No workflow logs yet"

# Show CRE Transmitter activity
transmitter-logs:
	@echo "=== CRE TRANSMITTER LOGS ==="
	@docker logs streams-node0 2>&1 | grep "ProcessReport" | tail -10 || echo "No transmitter logs yet"

help:
	@echo "LLO E2E Test Commands:"
	@echo ""
	@echo "  Interactive Mode:"
	@echo "    make run              - Run the E2E test (interactive output)"
	@echo "    make fresh            - Clean and run from scratch"
	@echo ""
	@echo "  Test Framework Mode (with assertions):"
	@echo "    make test             - Run full E2E test with Go test framework"
	@echo "    make test-quick       - Quick verify (env must be running)"
	@echo "    make fresh-test       - Clean and run test from scratch"
	@echo ""
	@echo "  Utilities:"
	@echo "    make clean            - Clean up all Docker containers"
	@echo "    make logs             - Show workflow report logs (proof)"
	@echo "    make transmitter-logs - Show CRE transmitter activity"
