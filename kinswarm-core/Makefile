# KinSwarm Core: High-Performance Orchestration
# Author: The Keeper

.PHONY: all build benchmark integration clean verify

all: build

build:
	@echo "--- Building KinSwarm Runtime ---"
	bash build_kinswarm_runtime.sh

benchmark:
	@echo "--- Running Logic Core Benchmark ---"
	bash benchmark_chainlink_core.sh

integration:
	@echo "--- Verifying Integration Adapter ---"
	cd KinSwarm_Integration && DYLD_LIBRARY_PATH="./lib" go run adapter_benchmark.go

verify:
	@echo "--- Auditing System Integrity ---"
	python3 audit_integrity.py

clean:
	@echo "--- Cleaning Build Artifacts ---"
	rm -rf out/
	rm -rf cache/
	rm -f *.raw
	rm -f *.hex
	@echo "Clean complete."
