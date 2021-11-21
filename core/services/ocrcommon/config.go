package ocrcommon

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/libocr/commontypes"
)

type Config interface {
	DefaultChainID() *big.Int
	Dev() bool
	EvmGasLimitDefault() uint64
	JobPipelineResultWriteQueueDepth() uint64
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRDefaultTransactionQueueDepth() uint32
	OCRKeyBundleID() (string, error)
	OCRObservationGracePeriod() time.Duration
	OCRObservationTimeout() time.Duration
	OCRTraceLogging() bool
	OCRTransmitterAddress() (ethkey.EIP55Address, error)
	P2PBootstrapPeers() ([]string, error)
	P2PPeerID() p2pkey.PeerID
	P2PV2Bootstrappers() []commontypes.BootstrapperLocator
	FlagsContractAddress() string
	ChainType() chains.ChainType
}
