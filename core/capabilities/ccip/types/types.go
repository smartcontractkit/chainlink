package types

import (
	"context"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
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
