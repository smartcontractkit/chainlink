package v1_0_0

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func CreateExecutionStateChangeEventLog(t *testing.T, seqNr uint64, blockNumber int64, messageID common.Hash) logpoller.Log {
	tAbi, err := evm_2_evm_offramp.EVM2EVMOffRampMetaData.GetAbi()
	require.NoError(t, err)
	eseEvent, ok := tAbi.Events["ExecutionStateChanged"]
	require.True(t, ok)

	logData, err := eseEvent.Inputs.NonIndexed().Pack(uint8(1), []byte("some return data"))
	require.NoError(t, err)
	seqNrBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqNrBytes, seqNr)
	seqNrTopic := common.BytesToHash(seqNrBytes)
	topic0 := evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{}.Topic()

	return logpoller.Log{
		Topics: [][]byte{
			topic0[:],
			seqNrTopic[:],
			messageID[:],
		},
		Data:        logData,
		LogIndex:    1,
		BlockHash:   utils.RandomBytes32(),
		BlockNumber: blockNumber,
		EventSig:    topic0,
		Address:     testutils.NewAddress(),
		TxHash:      utils.RandomBytes32(),
	}
}
