package types

import (
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

// CCIPOracle represents either a CCIP commit or exec oracle or a bootstrap node.
type CCIPOracle interface {
	Close() error
	Start() error
}

// OracleCreator is an interface for creating CCIP oracles.
// Whether the oracle uses a LOOPP or not is an implementation detail.
type OracleCreator interface {
	// CreatePlugin creates a new oracle that will run either the commit or exec ccip plugin.
	// The oracle must be returned unstarted.
	CreatePluginOracle(pluginType PluginType, config OCR3ConfigWithMeta) (CCIPOracle, error)

	// CreateBootstrapOracle creates a new bootstrap node with the given OCR config.
	// The oracle must be returned unstarted.
	CreateBootstrapOracle(config OCR3ConfigWithMeta) (CCIPOracle, error)
}
