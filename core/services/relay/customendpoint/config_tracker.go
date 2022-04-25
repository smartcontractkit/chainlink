package customendpoint

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var StaticTransmitters = []types.Account{"str"}
var StaticOnChainPublicKeys = []types.OnchainPublicKey{[]byte{'a'}}

// answer contains the value and other details of a particular transmission
type answer struct {
	Data      *big.Int
	Timestamp time.Time
	epoch     uint32
	round     uint8
}

type contractTracker struct {
	digester          OffchainConfigDigester
	bridgeRequestData string
	bridgeInputAtKey  string
	multiplierUsed    int32

	lggr        logger.Logger
	pipelineORM pipeline.ORM
	config      config.GeneralConfig
	reportCodec evmreportcodec.ReportCodec
	clock       utils.Nower

	transmittersWg sync.WaitGroup
	storedAnswer   answer
	ansLock        sync.RWMutex

	utils.StartStopOnce
}

func NewTracker(spec OCR2Spec,
	configDigester OffchainConfigDigester,
	lggr logger.Logger,
	pipelineORM pipeline.ORM,
	config config.GeneralConfig,
	reportCodec evmreportcodec.ReportCodec,
	clock utils.Nower) contractTracker {
	return contractTracker{
		digester:          configDigester,
		bridgeRequestData: spec.BridgeRequestData,
		bridgeInputAtKey:  spec.BridgeInputAtKey,
		multiplierUsed:    spec.MultiplierUsed,
		lggr:              lggr,
		pipelineORM:       pipelineORM,
		config:            config,
		reportCodec:       reportCodec,
		clock:             clock,
		storedAnswer: answer{
			Data:      big.NewInt(0),
			Timestamp: clock.Now(),
			epoch:     0,
			round:     0,
		},
	}
}

func (c *contractTracker) getLastTransmittedAnswer() answer {
	c.ansLock.RLock()
	defer c.ansLock.RUnlock()
	return c.storedAnswer
}

// Since we are returning a fixed config, so no need to notify changes about this config.
func (c *contractTracker) Notify() <-chan struct{} {
	return nil
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (c *contractTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	digest, err := c.digester.configDigest()
	return 1, digest, err
}

// LatestConfig always returns a fixed config, as it doesn't change.
func (c *contractTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return c.getContractConfig()
}

// LatestBlockHeight isn't used if LatestConfig() always returns a static config that doesn't
// change. So we can return a static value from here, which is a no-op.
func (c *contractTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 1, nil
}

// Return a fixed config.
// TODO: Figure out where to get config from. Job Spec, or API endpoint, or some onchain.
func (c *contractTracker) getContractConfig() (types.ContractConfig, error) {
	digest, err := c.digester.configDigest()

	return types.ContractConfig{
		ConfigDigest:          digest,
		ConfigCount:           uint64(1),
		Signers:               StaticOnChainPublicKeys,
		Transmitters:          StaticTransmitters,
		F:                     uint8(1),
		OnchainConfig:         []byte{'a'},
		OffchainConfigVersion: uint64(1),
		OffchainConfig:        []byte{'a'},
	}, err
}

func (c *contractTracker) Start() error {
	return nil
}

func (c *contractTracker) Close() error {
	c.transmittersWg.Wait()
	return nil
}

func (c *contractTracker) Ready() error {
	return nil
}

func (c *contractTracker) Healthy() error {
	return nil
}

// Waits till all the transmission threads are done. Used for testing.
func (c *contractTracker) WaitForTransmissions() {
	c.transmittersWg.Wait()
}
