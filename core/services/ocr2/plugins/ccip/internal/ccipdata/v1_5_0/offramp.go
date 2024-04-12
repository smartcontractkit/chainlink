package v1_5_0

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

var (
	abiOffRamp                        = abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
	_          ccipdata.OffRampReader = &OffRamp{}
)

type ExecOnchainConfig evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig

type OffRamp struct {
	*v1_2_0.OffRamp
	offRampV150 evm_2_evm_offramp.EVM2EVMOffRampInterface
}

// GetTokens Returns no data as the offRamps no longer have this information.
func (o *OffRamp) GetTokens(ctx context.Context) (cciptypes.OffRampTokens, error) {
	return cciptypes.OffRampTokens{
		SourceTokens:      []cciptypes.Address{},
		DestinationTokens: []cciptypes.Address{},
		DestinationPool:   make(map[cciptypes.Address]cciptypes.Address),
	}, nil
}

func (o *OffRamp) GetSourceToDestTokensMapping(ctx context.Context) (map[cciptypes.Address]cciptypes.Address, error) {
	return map[cciptypes.Address]cciptypes.Address{}, nil
}

func NewOffRamp(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int) (*OffRamp, error) {
	v120, err := v1_2_0.NewOffRamp(lggr, addr, ec, lp, estimator, destMaxGasPrice)
	if err != nil {
		return nil, err
	}

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(addr, ec)
	if err != nil {
		return nil, err
	}

	v120.ExecutionReportArgs = abihelpers.MustGetMethodInputs("manuallyExecute", abiOffRamp)[:1]

	return &OffRamp{
		OffRamp:     v120,
		offRampV150: offRamp,
	}, nil
}
