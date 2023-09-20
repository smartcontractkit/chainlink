package contractutil

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
)

func TestGetMessageIDsAsHexString(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		hashes := make([]common.Hash, 10)
		for i := range hashes {
			hashes[i] = common.HexToHash(strconv.Itoa(rand.Intn(100000)))
		}

		msgs := make([]evm_2_evm_offramp.InternalEVM2EVMMessage, len(hashes))
		for i := range msgs {
			msgs[i] = evm_2_evm_offramp.InternalEVM2EVMMessage{MessageId: hashes[i]}
		}

		messageIDs := GetMessageIDsAsHexString(msgs)
		for i := range messageIDs {
			assert.Equal(t, hashes[i].String(), messageIDs[i])
		}
	})

	t.Run("empty", func(t *testing.T) {
		messageIDs := GetMessageIDsAsHexString(nil)
		assert.Empty(t, messageIDs)
	})
}
