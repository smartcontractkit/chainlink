package ccipcommon

import cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

// TokenID is a struct that contains the token address and chain selector.
type TokenID struct {
	TokenAddress  cciptypes.Address
	ChainSelector uint64
}
