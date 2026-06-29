package changeset

import (
	"encoding/hex"
	"fmt"
	"strings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var SetCallAllowedChangeSet cldf.ChangeSetV2[vaulttypes.SetCallAllowedInput] = setCallAllowed{}

type setCallAllowed struct{}

func (s setCallAllowed) VerifyPreconditions(env cldf.Environment, config vaulttypes.SetCallAllowedInput) error {
	if len(config.Chains) == 0 {
		return fmt.Errorf("chains must not be empty")
	}
	for chainSelector, chainCfg := range config.Chains {
		if err := validateChainSelector(chainSelector, env); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if err := validateEthAddress("automationReceiverAddress", chainCfg.AutomationReceiverAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if err := validateEthAddress("targetAddress", chainCfg.TargetAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if _, err := parseSelectorHex(chainCfg.Selector); err != nil {
			return fmt.Errorf("chain %d: invalid selector %q: %w", chainSelector, chainCfg.Selector, err)
		}
	}
	return nil
}

func (s setCallAllowed) Apply(e cldf.Environment, config vaulttypes.SetCallAllowedInput) (cldf.ChangesetOutput, error) {
	evmChains := e.BlockChains.EVMChains()

	var primaryChainSelector uint64
	for sel := range config.Chains {
		if primaryChainSelector == 0 || sel < primaryChainSelector {
			primaryChainSelector = sel
		}
	}

	primaryChain := evmChains[primaryChainSelector]
	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}

	for chainSelector, chainCfg := range config.Chains {
		selector, err := parseSelectorHex(chainCfg.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: invalid selector: %w", chainSelector, err)
		}

		_, err = operations.ExecuteOperation(
			e.OperationsBundle,
			SetCallAllowedOperation,
			deps,
			SetCallAllowedOperationInput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: chainCfg.AutomationReceiverAddress,
				TargetAddress:             chainCfg.TargetAddress,
				Selector:                  selector,
				Allowed:                   chainCfg.Allowed,
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: setCallAllowed failed: %w", chainSelector, err)
		}
	}

	return cldf.ChangesetOutput{}, nil
}

// parseSelectorHex parses a hex string like "0x4b9f5c20" into a [4]byte selector.
func parseSelectorHex(s string) ([4]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return [4]byte{}, fmt.Errorf("not valid hex: %w", err)
	}
	if len(b) != 4 {
		return [4]byte{}, fmt.Errorf("selector must be exactly 4 bytes, got %d", len(b))
	}
	return [4]byte{b[0], b[1], b[2], b[3]}, nil
}
