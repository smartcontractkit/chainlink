package customendpoint

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var StaticTransmitters = []types.Account{"str"}
var StaticOnChainPublicKeys = []types.OnchainPublicKey{[]byte{'a'}}

// Answer contains the value and other details of a particular transmission
type Answer struct {
	Data      *big.Int
	Timestamp time.Time
	epoch     uint32
	round     uint8
}

type ContractTracker struct {
	digester          OffchainConfigDigester
	bridgeRequestData string
	bridgeInputAtKey  string

	lggr        logger.Logger
	pipelineORM pipeline.ORM
	config      config.GeneralConfig

	transmittersWg sync.WaitGroup
	answer         Answer
	ansLock        *sync.RWMutex

	utils.StartStopOnce
}

func NewTracker(spec OCR2Spec, configDigester OffchainConfigDigester, lggr logger.Logger, pipelineORM pipeline.ORM, config config.GeneralConfig) ContractTracker {
	return ContractTracker{
		digester:          configDigester,
		bridgeRequestData: spec.BridgeRequestData,
		bridgeInputAtKey:  spec.BridgeInputAtKey,
		lggr:              lggr,
		pipelineORM:       pipelineORM,
		config:            config,
		answer: Answer{
			Data:      nil,
			Timestamp: time.Now(),
			epoch:     0,
			round:     0,
		},
	}
}

func (c *ContractTracker) GetLastTransmittedAnswer() Answer {
	c.ansLock.RLock()
	defer c.ansLock.RUnlock()
	return c.answer
}

// Since we are returning a fixed config, so no need to notify changes about this config.
func (c *ContractTracker) Notify() <-chan struct{} {
	return nil
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (c *ContractTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	digest, err := c.digester.configDigest()
	return 1, digest, err
}

// LatestConfig always returns a fixed config, as it doesn't change.
func (c *ContractTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return c.getContractConfig()
}

// LatestBlockHeight isn't used if LatestConfig() always returns a static config that doesn't
// change. So we can return a static value from here, which is a no-op.
func (c *ContractTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 1, nil
}

// Return a fixed config.
// TODO: Check if the Signers and Transmitters are ok to be fixed static values
func (c *ContractTracker) getContractConfig() (types.ContractConfig, error) {
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

func (c *ContractTracker) Start() error {
	return nil
}

func (c *ContractTracker) Close() error {
	c.transmittersWg.Wait()
	return nil
}

func (c *ContractTracker) Ready() error {
	return nil
}

func (c *ContractTracker) Healthy() error {
	return nil
}
