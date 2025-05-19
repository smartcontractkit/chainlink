package evm

import (
	"context"
	"encoding/json"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewWriteTarget(ctx context.Context, relayer *Relayer, chain legacyevm.Chain, gasLimitDefault uint64, lggr logger.Logger) (*targets.WriteTarget, error) {
	// generate ID based on chain selector
	id := GenerateWriteTargetName(chain.ID().Uint64())

	// EVM-specific init
	config := chain.Config().EVM().Workflow()

	// Initialize a reader to check whether a value was already transmitted on chain
	contractReaderConfigEncoded, err := json.Marshal(relayevmtypes.ChainReaderConfig{
		Contracts: map[string]relayevmtypes.ChainContractReader{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainReaderDefinition{
					"getTransmissionInfo": {
						ChainSpecificName: "getTransmissionInfo",
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract reader config %v", err)
	}
	cr, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}

	chainWriterConfig := relayevmtypes.ChainWriterConfig{
		Contracts: map[string]*relayevmtypes.ContractConfig{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainWriterDefinition{
					"report": {
						ChainSpecificName: "report",
						FromAddress:       config.FromAddress().Address(),
						GasLimit:          gasLimitDefault,
					},
				},
			},
		},
	}
	chainWriterConfig.MaxGasPrice = chain.Config().EVM().GasEstimator().PriceMax()

	encodedWriterConfig, err := json.Marshal(chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chainwriter config: %w", err)
	}

	cw, err := relayer.NewContractWriter(ctx, encodedWriterConfig)
	if err != nil {
		return nil, err
	}

	return targets.NewWriteTarget(logger.Named(lggr, "WriteTarget"), id, cr, cw, config.ForwarderAddress().String(), gasLimitDefault), nil
}

func GenerateWriteTargetName(chainID uint64) string {
	id := fmt.Sprintf("write_%v@1.0.0", chainID)
	chainName, err := chainselectors.NameFromChainId(chainID)
	if err == nil {
		id = fmt.Sprintf("write_%v@1.0.0", chainName)
	}

	return id
}
