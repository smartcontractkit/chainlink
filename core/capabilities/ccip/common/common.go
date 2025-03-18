package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
)

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

// IsValid Exactly one of both plugins should be empty at any given time.
func (ofc OffChainConfig) IsValid() bool {
	return (ofc.CommitEmpty() && !ofc.ExecEmpty()) || (!ofc.CommitEmpty() && ofc.ExecEmpty())
}
