package leafer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"

	evm_2_evm_onramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
)

// LeafHasher converts a CCIPSendRequested event into something that can be hashed and hashes it.
type LeafHasher interface {
	HashLeaf(log types.Log) ([32]byte, error)
}

// Version is the contract to use.
type Version string

const (
	V1_2_0 Version = "v1_2_0"
	V1_5_0 Version = "v1_5_0"
)

// MakeLeafHasher is a factory function to construct the onramp implementing the HashLeaf function for a given version.
func MakeLeafHasher(ver Version, cl bind.ContractBackend, sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashutil.Hasher[[32]byte]) (LeafHasher, error) {
	switch ver {
	case V1_2_0:
		or, err := evm_2_evm_onramp_1_2_0.NewEVM2EVMOnRamp(onRampId, cl)
		if err != nil {
			return nil, err
		}
		h := v1_2_0.NewLeafHasher(sourceChainSelector, destChainSelector, onRampId, ctx, or)
		return h, nil
	case V1_5_0:
		or, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampId, cl)
		if err != nil {
			return nil, err
		}
		h := v1_5_0.NewLeafHasher(sourceChainSelector, destChainSelector, onRampId, ctx, or)
		return h, nil
	default:
		return nil, fmt.Errorf("unknown version %q", ver)
	}
}
