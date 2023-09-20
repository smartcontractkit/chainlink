package contractutil

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// IsCommitStoreDownNow Checks whether the commit store is down by doing an onchain check for Paused and ARM status
func IsCommitStoreDownNow(ctx context.Context, lggr logger.Logger, commitStore commit_store.CommitStoreInterface) bool {
	unPausedAndHealthy, err := commitStore.IsUnpausedAndARMHealthy(&bind.CallOpts{Context: ctx})
	if err != nil {
		// If we cannot read the state, assume the worst
		lggr.Errorw("Unable to read CommitStore IsUnpausedAndARMHealthy", "err", err)
		return true
	}
	return !unPausedAndHealthy
}

func GetMessageIDsAsHexString(messages []evm_2_evm_offramp.InternalEVM2EVMMessage) []string {
	messageIDs := make([]string, 0, len(messages))
	for _, m := range messages {
		messageIDs = append(messageIDs, "0x"+hex.EncodeToString(m.MessageId[:]))
	}
	return messageIDs
}
