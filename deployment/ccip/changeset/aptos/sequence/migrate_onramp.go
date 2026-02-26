package sequence

import (
	"fmt"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
)

var MigrateOnRampDestChainConfigsToV2Sequence = operations.NewSequence(
	"migrate-onramp-dest-chain-configs-to-v2-sequence",
	operation.Version1_0_0,
	"Migrates OnRamp destination chain configs from V1 to V2",
	migrateOnRampDestChainConfigsToV2Sequence,
)

func migrateOnRampDestChainConfigsToV2Sequence(b operations.Bundle, deps dependency.AptosDeps, in operation.MigrateOnRampDestChainConfigsToV2Input) (mcmstypes.BatchOperation, error) {
	report, err := operations.ExecuteOperation(b, operation.MigrateOnRampDestChainConfigsToV2Op, deps, in)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to migrate OnRamp dest chain configs to V2: %w", err)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  report.Output,
	}, nil
}
