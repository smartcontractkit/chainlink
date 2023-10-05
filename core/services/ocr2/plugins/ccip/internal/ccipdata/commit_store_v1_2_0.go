package ccipdata

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Do not change the JSON format of this struct without consulting with
// the RDD people first.
type CommitOffchainConfigV1_2_0 struct {
	SourceFinalityDepth      uint32
	DestFinalityDepth        uint32
	GasPriceHeartBeat        models.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      models.Duration
	TokenPriceDeviationPPB   uint32
	MaxGasPrice              uint64
	InflightCacheExpiry      models.Duration
}

func (c CommitOffchainConfigV1_2_0) Validate() error {
	if c.SourceFinalityDepth == 0 {
		return errors.New("must set SourceFinalityDepth")
	}
	if c.DestFinalityDepth == 0 {
		return errors.New("must set DestFinalityDepth")
	}
	if c.GasPriceHeartBeat.Duration() == 0 {
		return errors.New("must set GasPriceHeartBeat")
	}
	if c.DAGasPriceDeviationPPB == 0 {
		return errors.New("must set DAGasPriceDeviationPPB")
	}
	if c.ExecGasPriceDeviationPPB == 0 {
		return errors.New("must set ExecGasPriceDeviationPPB")
	}
	if c.TokenPriceHeartBeat.Duration() == 0 {
		return errors.New("must set TokenPriceHeartBeat")
	}
	if c.TokenPriceDeviationPPB == 0 {
		return errors.New("must set TokenPriceDeviationPPB")
	}
	if c.MaxGasPrice == 0 {
		return errors.New("must set MaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}

	return nil
}

type CommitStoreV1_2_0 struct {
	*CommitStoreV1_0_0
	// Dynamic config
	configMu          sync.RWMutex
	gasPriceEstimator prices.DAGasPriceEstimator
	offchainConfig    CommitOffchainConfig
}

func (c *CommitStoreV1_2_0) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[CommitOnchainConfig](onchainConfig)
	if err != nil {
		return common.Address{}, err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[CommitOffchainConfigV1_2_0](offchainConfig)
	if err != nil {
		return common.Address{}, err
	}
	c.configMu.Lock()
	c.gasPriceEstimator = prices.NewDAGasPriceEstimator(
		c.estimator,
		big.NewInt(int64(offchainConfigParsed.MaxGasPrice)),
		int64(offchainConfigParsed.ExecGasPriceDeviationPPB),
		int64(offchainConfigParsed.DAGasPriceDeviationPPB),
	)
	c.offchainConfig = CommitOffchainConfig{
		SourceFinalityDepth:    offchainConfigParsed.SourceFinalityDepth,
		GasPriceDeviationPPB:   offchainConfigParsed.ExecGasPriceDeviationPPB,
		TokenPriceDeviationPPB: offchainConfigParsed.TokenPriceDeviationPPB,
		InflightCacheExpiry:    offchainConfigParsed.InflightCacheExpiry.Duration(),
		DestFinalityDepth:      offchainConfigParsed.DestFinalityDepth,
	}
	c.configMu.Unlock()

	c.lggr.Infow("ChangeConfig",
		"offchainConfig", offchainConfigParsed,
		"onchainConfig", onchainConfigParsed,
	)
	return onchainConfigParsed.PriceRegistry, nil
}

func (c *CommitStoreV1_2_0) OffchainConfig() CommitOffchainConfig {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.offchainConfig
}

func (c *CommitStoreV1_2_0) GasPriceEstimator() prices.GasPriceEstimatorCommit {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.gasPriceEstimator
}

func NewCommitStoreV1_2_0(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (*CommitStoreV1_2_0, error) {
	commitStoreV100, err := NewCommitStoreV1_0_0(lggr, addr, ec, lp, estimator)
	if err != nil {
		return nil, err
	}
	return &CommitStoreV1_2_0{CommitStoreV1_0_0: commitStoreV100}, nil
}
