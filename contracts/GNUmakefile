# To make a snapshot for a specific product, either set the `FOUNDRY_PROFILE` env var or call the target with `FOUNDRY_PROFILE=product`
# env var example
#		export FOUNDRY_PROFILE=llo
#		make snapshot
# make call example
# 		make FOUNDRY_PROFILE=llo snapshot
.PHONY: snapshot
snapshot:
	export FOUNDRY_PROFILE=$(FOUNDRY_PROFILE) && forge snapshot --snap gas-snapshots/$(FOUNDRY_PROFILE).gas-snapshot

.PHONY: snapshot-all
snapshot-all:
	for foundry_profile in vrf automation llo functions automation-dev ; do \
		make snapshot FOUNDRY_PROFILE=$$foundry_profile ; \
	done

help:
	@echo ""
	@echo "         .__           .__       .__  .__        __"
	@echo "    ____ |  |__ _____  |__| ____ |  | |__| ____ |  | __"
	@echo "  _/ ___\|  |  \\\\\\__  \ |  |/    \|  | |  |/    \|  |/ /"
	@echo "  \  \___|   Y  \/ __ \|  |   |  \  |_|  |   |  \    <"
	@echo "   \___  >___|  (____  /__|___|  /____/__|___|  /__|_ \\"
	@echo "       \/     \/     \/        \/             \/     \/"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
