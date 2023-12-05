package ccipdata

import (
	"context"
	"math/big"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const (
	ManuallyExecute = "manuallyExecute"
)

// Do not change the JSON format of this struct without consulting with
// the RDD people first.
type ExecOffchainConfig struct {
	SourceFinalityDepth         uint32
	DestOptimisticConfirmations uint32
	DestFinalityDepth           uint32
	BatchGasLimit               uint32
	RelativeBoostPerWaitHour    float64
	MaxGasPrice                 uint64
	InflightCacheExpiry         models.Duration
	RootSnoozeTime              models.Duration
}

func (c ExecOffchainConfig) Validate() error {
	if c.SourceFinalityDepth == 0 {
		return errors.New("must set SourceFinalityDepth")
	}
	if c.DestFinalityDepth == 0 {
		return errors.New("must set DestFinalityDepth")
	}
	if c.DestOptimisticConfirmations == 0 {
		return errors.New("must set DestOptimisticConfirmations")
	}
	if c.BatchGasLimit == 0 {
		return errors.New("must set BatchGasLimit")
	}
	if c.RelativeBoostPerWaitHour == 0 {
		return errors.New("must set RelativeBoostPerWaitHour")
	}
	if c.MaxGasPrice == 0 {
		return errors.New("must set MaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	if c.RootSnoozeTime.Duration() == 0 {
		return errors.New("must set RootSnoozeTime")
	}

	return nil
}

type ExecOnchainConfig struct {
	PermissionLessExecutionThresholdSeconds time.Duration
}

func (c ExecOnchainConfig) Validate() error {
	if c.PermissionLessExecutionThresholdSeconds == 0 {
		return errors.New("must set PermissionLessExecutionThresholdSeconds")
	}

	return nil
}

type ExecutionStateChanged struct {
	SequenceNumber uint64
	Finalized      bool
}

type ExecReport struct {
	Messages          []internal.EVM2EVMMessage
	OffchainTokenData [][][]byte
	Proofs            [][32]byte
	ProofFlagBits     *big.Int
}

type OffRampStaticConfig struct {
	CommitStore         common.Address
	ChainSelector       uint64
	SourceChainSelector uint64
	OnRamp              common.Address
	PrevOffRamp         common.Address
	ArmProxy            common.Address
}

type TokenBucketRateLimit struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

//go:generate mockery --quiet --name OffRampReader --filename offramp_reader_mock.go --case=underscore
type OffRampReader interface {
	Closer
	RegisterFilters(qopts ...pg.QOpt) error
	// Will error if messages are not a compatible version.
	EncodeExecutionReport(report ExecReport) ([]byte, error)
	DecodeExecutionReport(report []byte) (ExecReport, error)
	// GetExecutionStateChangesBetweenSeqNums returns all the execution state change events for the provided message sequence numbers (inclusive).
	GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]Event[ExecutionStateChanged], error)
	GetDestinationTokens(ctx context.Context) ([]common.Address, error)
	GetTokenPoolsRateLimits(ctx context.Context, poolAddresses []common.Address) ([]TokenBucketRateLimit, error)
	Address() common.Address
	// Notifies the reader that the config has changed onchain
	ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, common.Address, error)
	OffchainConfig() ExecOffchainConfig
	OnchainConfig() ExecOnchainConfig
	GasPriceEstimator() prices.GasPriceEstimatorExec
	GetSenderNonce(ctx context.Context, sender common.Address) (uint64, error)
	CurrentRateLimiterState(ctx context.Context) (evm_2_evm_offramp.RateLimiterTokenBucket, error)
	GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error)
	GetStaticConfig(ctx context.Context) (OffRampStaticConfig, error)
	GetSourceToDestTokensMapping(ctx context.Context) (map[common.Address]common.Address, error)
	GetDestinationTokenPools(ctx context.Context) (map[common.Address]common.Address, error)
}

// MessageExecutionState defines the execution states of CCIP messages.
type MessageExecutionState uint8

const (
	ExecutionStateUntouched MessageExecutionState = iota
	ExecutionStateInProgress
	ExecutionStateSuccess
	ExecutionStateFailure
)

func NewOffRampReader(lggr logger.Logger, addr common.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (OffRampReader, error) {
	_, version, err := ccipconfig.TypeAndVersion(addr, destClient)
	if err != nil {
		return nil, err
	}
	switch version.String() {
	case V1_0_0, V1_1_0:
		return NewOffRampV1_0_0(lggr, addr, destClient, lp, estimator)
	case V1_2_0, V1_3_0:
		return NewOffRampV1_2_0(lggr, addr, destClient, lp, estimator)
	default:
		return nil, errors.Errorf("unsupported offramp version %v", version.String())
	}
	// TODO can validate it pointing to the correct version
}

func ExecReportToEthTxMeta(typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	if typ != ccipconfig.EVM2EVMOffRamp {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOffRamp, typ)
	}
	switch ver.String() {
	case V1_0_0, V1_1_0:
		offRampABI := abihelpers.MustParseABI(evm_2_evm_offramp_1_0_0.EVM2EVMOffRampABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			execReport, err := decodeExecReportV1_0_0(abihelpers.MustGetMethodInputs(ManuallyExecute, offRampABI)[:1], report)
			if err != nil {
				return nil, err
			}
			return execReportToEthTxMeta(execReport)
		}, nil
	case V1_2_0, V1_3_0:
		offRampABI := abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			execReport, err := decodeExecReportV1_2_0(abihelpers.MustGetMethodInputs(ManuallyExecute, offRampABI)[:1], report)
			if err != nil {
				return nil, err
			}
			return execReportToEthTxMeta(execReport)
		}, nil
	default:
		return nil, errors.Errorf("got unexpected version %v", ver.String())
	}
}

// EncodeExecutionReport is only used in tests
// TODO should remove it and update tests to use Reader interface.
func EncodeExecutionReport(report ExecReport) ([]byte, error) {
	offRampABI := abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
	return encodeExecutionReportV1_2_0(abihelpers.MustGetMethodInputs(ManuallyExecute, offRampABI)[:1], report)
}

func execReportToEthTxMeta(execReport ExecReport) (*txmgr.TxMeta, error) {
	msgIDs := make([]string, len(execReport.Messages))
	for i, msg := range execReport.Messages {
		msgIDs[i] = hexutil.Encode(msg.MessageId[:])
	}

	return &txmgr.TxMeta{
		MessageIDs: msgIDs,
	}, nil
}
