.PHONY: run clean fresh proof streams-logs workflow-logs ea-logs status test test-quick help

# The environment directory (parent of llo_e2e)
ENV_DIR := $(shell cd .. && pwd)

# Run the full LLO E2E test via CRE CLI (interactive mode)
run:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Starting LLO Streams Trigger E2E Test                       ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@cd $(ENV_DIR) && CRE_ENV_DIR=$(ENV_DIR) go run . llo-e2e

# Run the full E2E test using Go test framework with assertions
# -count=1 disables Go test caching (required since we use Docker containers)
test:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Running LLO E2E Test (Go Test Framework)                    ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@cd $(ENV_DIR) && go test -v -count=1 -timeout 30m ./llo_e2e -run TestLLOStreamsE2E

# Quick verification test (assumes environment is already running)
# -count=1 disables Go test caching
test-quick:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Quick E2E Verification (environment must be running)        ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@cd $(ENV_DIR) && go test -v -count=1 -timeout 5m ./llo_e2e -run TestLLOStreamsE2E_QuickVerify

# Clean up all containers and start fresh
clean:
	@echo "Cleaning up Docker containers..."
	@docker ps -aq | xargs docker rm -f 2>/dev/null || true
	@docker ps -a --format "{{.Names}}" | xargs -I {} docker rm -f {} 2>/dev/null || true
	@docker rmi -f mock-ea-e2e mock-ea mock-ea:latest 2>/dev/null || true
	@docker network prune -f 2>/dev/null || true
	@docker volume prune -f 2>/dev/null || true
	@rm -f $(ENV_DIR)/state/*.state.yaml 2>/dev/null || true
	@sleep 3
	@echo "Environment cleared"

# Clean and run interactive mode
fresh: clean run

# Clean and run test mode
fresh-test: clean test

# Show E2E proof - workflow received MAGIC_NUMBER from Streams DON
proof:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  E2E PROOF: Workflow received value from Streams DON         ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Looking for LLO_E2E_VALUE logs (MAGIC_NUMBER=424242)..."
	@docker logs workflow-node0 2>&1 | grep "LLO_E2E_VALUE" | tail -5 || echo "❌ No LLO_E2E_VALUE logs found"

# Show Streams DON CRE Transmitter logs
streams-logs:
	@echo "=== STREAMS DON: CRE Transmitter Activity ==="
	@docker logs streams-node0 2>&1 | grep "ProcessReport pushing event" | tail -10 || echo "No transmitter activity yet"
	@echo ""
	@echo "=== STREAMS DON: Subscriber Registrations ==="
	@docker logs streams-node0 2>&1 | grep "RegisterTrigger" | tail -5 || echo "No subscribers registered yet"

# Show Workflow DON logs
workflow-logs:
	@echo "=== WORKFLOW DON: Trigger Events Received ==="
	@docker logs workflow-node0 2>&1 | grep -E "LLO_E2E_VALUE|LLO_CONSUMER|LLO_WORKFLOW" | tail -10 || echo "No workflow trigger logs yet"

# Show Mock EA logs
ea-logs:
	@echo "=== MOCK EA: Price Requests (MAGIC_NUMBER=424242) ==="
	@docker logs mock-ea 2>&1 | grep -E "Received|424242|TEST/USD" | tail -10 || echo "No EA requests yet"

# Show status of all three stages
status:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  E2E DATA FLOW STATUS                                        ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "STAGE 1: Mock EA → Stream Jobs"
	@docker logs mock-ea 2>&1 | grep -q "424242" && echo "  ✓ Mock EA returning MAGIC_NUMBER" || echo "  ✗ Mock EA not receiving requests"
	@echo ""
	@echo "STAGE 2: LLO Plugin → CRE Transmitter"
	@docker logs streams-node0 2>&1 | grep -q "ProcessReport pushing event" && echo "  ✓ CRE Transmitter emitting events" || echo "  ✗ No CRE Transmitter activity"
	@echo ""
	@echo "STAGE 3: Workflow DON → Consumer Workflow"
	@docker logs workflow-node0 2>&1 | grep -q "LLO_E2E_VALUE" && echo "  ✓ Workflow decoded MAGIC_NUMBER" || echo "  ✗ No workflow output yet"

help:
	@echo "LLO E2E Test Commands:"
	@echo ""
	@echo "  Run Tests:"
	@echo "    make run              - Run E2E test (interactive output)"
	@echo "    make test             - Run E2E test with Go test framework"
	@echo "    make test-quick       - Quick verify (env must be running)"
	@echo "    make fresh            - Clean and run from scratch"
	@echo "    make fresh-test       - Clean and run test from scratch"
	@echo ""
	@echo "  View Logs:"
	@echo "    make status           - Show status of all 3 stages"
	@echo "    make proof            - Show E2E proof (workflow output)"
	@echo "    make ea-logs          - Show Mock EA requests"
	@echo "    make streams-logs     - Show Streams DON transmitter logs"
	@echo "    make workflow-logs    - Show Workflow DON trigger logs"
	@echo ""
	@echo "  Cleanup:"
	@echo "    make clean            - Remove all Docker containers"
