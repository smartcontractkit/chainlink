package llo

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	StreamsCapabilityName    = "streams"
	StreamsCapabilityVersion = "1.0.0"
)

type configProvider struct {
	ConfigPoller
	digester ocr2types.OffchainConfigDigester
}

func NewConfigProvider(lggr logger.Logger, chain legacyevm.Chain, relayCfg types.RelayConfig, opts *types.RelayOpts) (commontypes.ConfigProvider, error) {
	// ContractID should be the address of the capabilities registry contract
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	capabilitiesRegistryAddr := common.HexToAddress(opts.ContractID)
	digester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: capabilitiesRegistryAddr,
	}

	cp, err := NewConfigPoller(
		CPConfig{
			Logger:                      lggr,
			Client:                      chain.Client(),
			CapabilitiesRegistryAddress: capabilitiesRegistryAddr,
			DonID:                       relayCfg.DonID,
			CapabilityName:              StreamsCapabilityName,
			CapabilityVersion:           StreamsCapabilityVersion,
		},
	)
	if err != nil {
		return nil, err
	}

	return &configProvider{cp, digester}, nil
}

func (cp *configProvider) ContractConfigTracker() ocr2types.ContractConfigTracker {
	return cp
}

func (cp *configProvider) OffchainConfigDigester() ocr2types.OffchainConfigDigester {
	return cp.digester
}
