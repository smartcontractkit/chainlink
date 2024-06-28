package evm

import (
	"context"
	"encoding/json"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewWriteTarget(ctx context.Context, relayer *Relayer, chain legacyevm.Chain, lggr logger.Logger) (*targets.WriteTarget, error) {
	// EVM-specific init
	config := chain.Config().EVM().Workflow()
	fromAddress := config.FromAddress()
	forwarderAddress := config.ForwarderAddress().String()
	l := lggr.Named("NewEVMWriteTarget").With("fromAddress", fromAddress.String(), "forwarderAddress", forwarderAddress)

	// generate ID based on chain selector
	id := fmt.Sprintf("write_%v@1.0.0", chain.ID())
	chainName, err := chainselectors.NameFromChainId(chain.ID().Uint64())
	if err == nil {
		id = fmt.Sprintf("write_%v@1.0.0", chainName)
		l = l.With("capabilityID", id)
	} else {
		l.Warnw("failed to get chain name from chain ID", "chainID", chain.ID().Uint64())
		l = l.With("capabilityID", id)
	}

	// Initialize a reader to check whether a value was already transmitted on chain
	contractReaderConfigEncoded, err := json.Marshal(relayevmtypes.ChainReaderConfig{
		Contracts: map[string]relayevmtypes.ChainContractReader{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainReaderDefinition{
					"getTransmitter": {
						ChainSpecificName: "getTransmitter",
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract reader config %v", err)
	}
	cr, err := relayer.NewContractReader(contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}
	err = cr.Bind(ctx, []commontypes.BoundContract{{
		Address: forwarderAddress,
		Name:    "forwarder",
	}})
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
						Checker:           "simulate",
						FromAddress:       fromAddress.Address(),
						GasLimit:          200_000,
					},
				},
			},
		},
	}

	chainWriterConfig.MaxGasPrice = chain.Config().EVM().GasEstimator().PriceMax()
	cw, err := NewChainWriterService(l.Named("ChainWriter"), chain.Client(), chain.TxManager(), chain.GasEstimator(), chainWriterConfig)
	if err != nil {
		return nil, err
	}

	return targets.NewWriteTarget(l, id, cr, cw, forwarderAddress), nil
}
