#!/bin/sh
export NUM_OF_PAGES=all
export ENVIRONMENT=integration
export DRY_RUN=false
export REPOSITORY=smartcontractkit/chainlink
export REF=fix/golint
export GITHUB_ACTION=true

pnpm start
