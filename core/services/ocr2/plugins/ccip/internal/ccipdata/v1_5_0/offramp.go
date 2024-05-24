package v1_5_0

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

var (
	abiOffRamp                                        = abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
	_                          ccipdata.OffRampReader = &OffRamp{}
	RateLimitTokenAddedEvent                          = abihelpers.MustGetEventID("TokenAggregateRateLimitAdded", abiOffRamp)
	RateLimitTokenRemovedEvent                        = abihelpers.MustGetEventID("TokenAggregateRateLimitRemoved", abiOffRamp)
)

type ExecOnchainConfig evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig

type OffRamp struct {
	*v1_2_0.OffRamp
	offRampV150           evm_2_evm_offramp.EVM2EVMOffRampInterface
	cachedRateLimitTokens cache.AutoSync[cciptypes.OffRampTokens]
}

// GetTokens Returns no data as the offRamps no longer have this information.
func (o *OffRamp) GetTokens(ctx context.Context) (cciptypes.OffRampTokens, error) {
	sourceTokens, destTokens, err := o.GetSourceAndDestRateLimitTokens(ctx)
	if err != nil {
		return cciptypes.OffRampTokens{}, err
	}
	return cciptypes.OffRampTokens{
		SourceTokens:      sourceTokens,
		DestinationTokens: destTokens,
	}, nil
}

func (o *OffRamp) GetSourceAndDestRateLimitTokens(ctx context.Context) (sourceTokens []cciptypes.Address, destTokens []cciptypes.Address, err error) {
	cachedTokens, err := o.cachedRateLimitTokens.Get(ctx, func(ctx context.Context) (cciptypes.OffRampTokens, error) {
		tokens, err2 := o.offRampV150.GetAllRateLimitTokens(&bind.CallOpts{Context: ctx})
		if err2 != nil {
			return cciptypes.OffRampTokens{}, err2
		}

		if len(tokens.SourceTokens) != len(tokens.DestTokens) {
			return cciptypes.OffRampTokens{}, errors.New("source and destination tokens are not the same length")
		}

		return cciptypes.OffRampTokens{
			DestinationTokens: ccipcalc.EvmAddrsToGeneric(tokens.DestTokens...),
			SourceTokens:      ccipcalc.EvmAddrsToGeneric(tokens.SourceTokens...),
		}, nil
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get rate limit tokens, if token set is large (~400k) batching may be needed")
	}
	return cachedTokens.SourceTokens, cachedTokens.DestinationTokens, nil
}

func (o *OffRamp) GetSourceToDestTokensMapping(ctx context.Context) (map[cciptypes.Address]cciptypes.Address, error) {
	sourceTokens, destTokens, err := o.GetSourceAndDestRateLimitTokens(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rate limit tokens, if token set is large (~400k) batching may be needed")
	}

	if sourceTokens == nil || destTokens == nil {
		return nil, errors.New("source or destination tokens are nil")
	}

	mapping := make(map[cciptypes.Address]cciptypes.Address)
	for i, sourceToken := range sourceTokens {
		mapping[sourceToken] = destTokens[i]
	}
	return mapping, nil
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
		cachedRateLimitTokens: cache.NewLogpollerEventsBased[cciptypes.OffRampTokens](
			lp,
			[]common.Hash{RateLimitTokenAddedEvent, RateLimitTokenRemovedEvent},
			offRamp.Address(),
		),
	}, nil
}
