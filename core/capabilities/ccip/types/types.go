package types

import (
	"context"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// OCR3ConfigWithMeta is a type alias in order to generate correct mocks for the OracleCreator interface.
type OCR3ConfigWithMeta ccipreaderpkg.OCR3ConfigWithMeta

// PluginType represents the type of CCIP plugin.
// It mirrors the OCRPluginType in Internal.sol.
type PluginType uint8

const (
	PluginTypeCCIPCommit PluginType = 0
	PluginTypeCCIPExec   PluginType = 1
)

func (pt PluginType) String() string {
	switch pt {
	case PluginTypeCCIPCommit:
		return "CCIPCommit"
	case PluginTypeCCIPExec:
		return "CCIPExec"
	default:
		return "Unknown"
	}
}

type OracleType uint8

const (
	OracleTypePlugin    OracleType = 0
	OracleTypeBootstrap OracleType = 1
)

// CCIPOracle represents either a CCIP commit or exec oracle or a bootstrap node.
type CCIPOracle interface {
	Close() error
	Start() error
}

// OracleCreator is an interface for creating CCIP oracles.
// Whether the oracle uses a LOOPP or not is an implementation detail.
type OracleCreator interface {
	// Create creates a new oracle that will run either the commit or exec ccip plugin,
	// if its a plugin oracle, or a bootstrap oracle if its a bootstrap oracle.
	// The oracle must be returned unstarted.
	Create(ctx context.Context, donID uint32, config OCR3ConfigWithMeta) (CCIPOracle, error)

	// Type returns the type of oracle that this creator creates.
	// The only valid values are OracleTypePlugin and OracleTypeBootstrap.
	Type() OracleType
}

type OffChainConfig struct {
	CommitOffchainConfig *pluginconfig.CommitOffchainConfig
	ExecOffchainConfig   *pluginconfig.ExecuteOffchainConfig
}

func (ofc OffChainConfig) CommitEmpty() bool {
	return ofc.CommitOffchainConfig == nil
}

func (ofc OffChainConfig) ExecEmpty() bool {
	return ofc.ExecOffchainConfig == nil
}

func (ofc OffChainConfig) Commit() *pluginconfig.CommitOffchainConfig {
	return ofc.CommitOffchainConfig
}

func (ofc OffChainConfig) Exec() *pluginconfig.ExecuteOffchainConfig {
	return ofc.ExecOffchainConfig
}

// Exactly one of both plugins should be empty at any given time.
func (ofc OffChainConfig) IsValid() bool {
	return (ofc.CommitEmpty() && !ofc.ExecEmpty()) || (!ofc.CommitEmpty() && ofc.ExecEmpty())
}

// PluginConfig is a struct that holds all the necessary information for a CCIP plugin to function.
type PluginConfig struct {
	CommitPluginCodec   cciptypes.CommitPluginCodec
	ExecutePluginCodec  cciptypes.ExecutePluginCodec
	MessageHasher       func(lggr logger.Logger) cciptypes.MessageHasher
	TokenDataEncoder    cciptypes.TokenDataEncoder
	GasEstimateProvider cciptypes.EstimateProvider
	RMNCrypto           func(lggr logger.Logger) cciptypes.RMNCrypto
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn    string
	GetChainReaderConfig func(lggr logger.Logger,
		chainID string,
		destChainID string,
		homeChainID string,
		ofc OffChainConfig,
		chainSelector cciptypes.ChainSelector,
	) ([]byte, error)
	GetChainWriter func(
		ctx context.Context,
		chainID string,
		relayer loop.Relayer,
		transmitters map[types.RelayID][]string,
		execBatchGasLimit uint64,
		chainFamily string,
		offrampProgramAddress []byte,
	) (types.ContractWriter, error)
}
