package types

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// Generate types from third-party repos:
//
//go:generate mockery --name Logger --structname MockLogger --srcpkg "github.com/smartcontractkit/libocr/commontypes" --outpkg types --output . --case=underscore --filename logger.generated.go

// HeadSubscriber represents head subscriber behaviour; used for evm chains;
//
//go:generate mockery --name HeadSubscriber --inpackage --output . --case=underscore --filename head_subscribed.generated.go
type HeadSubscriber interface {
	HeadTicker() chan BlockKey
}

// EVMClient represents evm client's behavior
//
//go:generate mockery --name EVMClient --inpackage --output . --case=underscore --filename evm_client.generated.go
type EVMClient interface {
	HeadSubscriber
	bind.ContractCaller
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

// Registry represents keeper registry behaviour
//
//go:generate mockery --name Registry --inpackage --output . --case=underscore --filename registry.generated.go
type Registry interface {
	GetActiveUpkeepIDs(context.Context) ([]UpkeepIdentifier, error)
	CheckUpkeep(context.Context, ...UpkeepKey) (UpkeepResults, error)
}

// ReportEncoder represents the report encoder behaviour
//
//go:generate mockery --name ReportEncoder --inpackage --output . --case=underscore --filename report_encoder.generated.go
type ReportEncoder interface {
	EncodeReport([]UpkeepResult) ([]byte, error)
	DecodeReport([]byte) ([]UpkeepResult, error)
}

// PerformLogProvider represents the perform log provider
//
//go:generate mockery --name PerformLogProvider --inpackage --output . --case=underscore --filename perform_log_provider.generated.go
type PerformLogProvider interface {
	PerformLogs(context.Context) ([]PerformLog, error)
	StaleReportLogs(context.Context) ([]StaleReportLog, error)
}

type PerformLog struct {
	Key             UpkeepKey
	TransmitBlock   BlockKey
	Confirmations   int64
	TransactionHash string
}

type StaleReportLog struct {
	Key             UpkeepKey
	TransmitBlock   BlockKey
	Confirmations   int64
	TransactionHash string
}

type BlockKey interface {
	After(BlockKey) (bool, error)
	Next() (BlockKey, error)
	BigInt() (*big.Int, bool)
	fmt.Stringer
}

type Address []byte

// UpkeepKey is an identifier of an upkeep at a moment in time, typically an
// upkeep at a block number
type UpkeepKey interface {
	BlockKeyAndUpkeepID() (BlockKey, UpkeepIdentifier, error)
	fmt.Stringer
}

// UpkeepIdentifier is an identifier for an active upkeep, typically a big int
type UpkeepIdentifier []byte

// BigInt creates and returns big int from the given upkeep identifier
func (ui UpkeepIdentifier) BigInt() (*big.Int, bool) {
	return big.NewInt(0).SetString(string(ui), 10)
}

type UpkeepResults []UpkeepResult

type UpkeepResult struct {
	Key              UpkeepKey
	State            UpkeepState
	FailureReason    uint8
	GasUsed          *big.Int
	PerformData      []byte
	FastGasWei       *big.Int
	LinkNative       *big.Int
	CheckBlockNumber uint32
	CheckBlockHash   [32]byte
	ExecuteGas       uint32
}

type UpkeepState uint

const (
	NotEligible UpkeepState = iota
	Eligible
)

type OffchainConfig struct {
	// PerformLockoutWindow is the window in which a single upkeep cannot be
	// performed again while waiting for a confirmation. Standard setting is
	// 100 blocks * average block time. Units are in milliseconds
	PerformLockoutWindow int64 `json:"performLockoutWindow"`

	// TargetProbability is the probability that all upkeeps will be checked
	// within the provided number rounds
	TargetProbability string `json:"targetProbability"`

	// TargetInRounds is the number of rounds for the above probability to be
	// calculated
	TargetInRounds int `json:"targetInRounds"`

	// SamplingJobDuration is the time allowed for a sampling run to complete
	// before forcing a new job on the latest block. Units are in milliseconds.
	SamplingJobDuration int64 `json:"samplingJobDuration"`

	// MinConfirmations limits registered log events to only those that have
	// the provided number of confirmations.
	MinConfirmations int `json:"minConfirmations"`

	// GasLimitPerReport is the max gas that could be spent per one report.
	// This is needed for calculation of how many upkeeps could be within report.
	GasLimitPerReport uint32 `json:"gasLimitPerReport"`

	// GasOverheadPerUpkeep is gas overhead per upkeep taken place in the report.
	GasOverheadPerUpkeep uint32 `json:"gasOverheadPerUpkeep"`

	// MaxUpkeepBatchSize is the max upkeep batch size of the OCR2 report.
	MaxUpkeepBatchSize int `json:"maxUpkeepBatchSize"`

	// ReportBlockLag is the number to subtract from median block number during report phase.
	ReportBlockLag int `json:"reportBlockLag"`
}

func DecodeOffchainConfig(b []byte) (OffchainConfig, error) {
	var config OffchainConfig
	var err error

	if len(b) > 0 {
		err = json.Unmarshal(b, &config)
	}

	if config.PerformLockoutWindow <= 0 {
		config.PerformLockoutWindow = 20 * 60 * 1000 // default of 20 minutes (100 blocks on eth)
	}

	if len(config.TargetProbability) == 0 {
		config.TargetProbability = "0.99999"
	}

	if config.TargetInRounds <= 0 {
		config.TargetInRounds = 1
	}

	if config.SamplingJobDuration <= 0 {
		config.SamplingJobDuration = 3000 // default of 3 seconds if not set
	}

	if config.MinConfirmations <= 0 {
		config.MinConfirmations = 0 // default of 0
	}

	if config.GasLimitPerReport == 0 { // defined as uint so cannot be < 0
		config.GasLimitPerReport = 5_300_000
	}

	if config.GasOverheadPerUpkeep == 0 { // defined as uint so cannot be < 0
		config.GasOverheadPerUpkeep = 300_000
	}

	if config.MaxUpkeepBatchSize <= 0 {
		config.MaxUpkeepBatchSize = 1
	}

	if config.ReportBlockLag < 0 {
		config.ReportBlockLag = 0
	}

	return config, err
}

func (c OffchainConfig) Encode() []byte {
	b, err := json.Marshal(&c)
	if err != nil {
		panic(fmt.Sprintf("unexpected error json encoding OffChainConfig: %s", err))
	}

	return b
}
