package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

// ExtraDataDecoded contains a generic representation of chain specific message parameters. A
// map from string to any is used to account for different parameters required for sending messages
// to different destinations.
type ExtraDataDecoded struct {
	// ExtraArgsDecoded contain message specific extra args.
	ExtraArgsDecoded map[string]any
	// DestExecDataDecoded contain token transfer specific extra args.
	DestExecDataDecoded []map[string]any
}

// HashedCapabilityID returns the hashed capability id in a manner equivalent to the capability registry.
func HashedCapabilityID(capabilityLabelledName, capabilityVersion string) (r [32]byte, err error) {
	// TODO: investigate how to avoid parsing the ABI everytime.
	tabi := `[{"type": "string"}, {"type": "string"}]`
	abiEncoded, err := utils.ABIEncode(tabi, capabilityLabelledName, capabilityVersion)
	if err != nil {
		return r, fmt.Errorf("failed to ABI encode capability version and labelled name: %w", err)
	}

	h := crypto.Keccak256(abiEncoded)
	copy(r[:], h)
	return r, nil
}

type OffChainConfig struct {
	Commit  *pluginconfig.CommitOffchainConfig
	Execute *pluginconfig.ExecuteOffchainConfig
}

func (ofc OffChainConfig) CommitEmpty() bool {
	return ofc.Commit == nil
}

func (ofc OffChainConfig) ExecEmpty() bool {
	return ofc.Execute == nil
}

// IsValid Exactly one of both plugins should be empty at any given time.
func (ofc OffChainConfig) IsValid() bool {
	return (ofc.CommitEmpty() && !ofc.ExecEmpty()) || (!ofc.CommitEmpty() && ofc.ExecEmpty())
}
