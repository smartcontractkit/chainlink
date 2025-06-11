package read

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/onramp"
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DecodeHardcodedType(t *testing.T) {
	t.Parallel()

	t.Run("decode hardcoded type offramp CoommitReportAccess success", func(t *testing.T) {
		fixtLog := getFixtureOfframpLog()

		log, err := generateOfframpLog(fixtLog)
		require.NoError(t, err)

		var out reader.CommitReportAcceptedEvent
		err = decodeHardcodedType(&out, log)
		require.NoError(t, err)

		require.Equal(t, true, bytes.Equal(fixtLog.BlessedMerkleRoots[0].MerkleRoot[:], out.BlessedMerkleRoots[0].MerkleRoot[:]))
		require.Equal(t, true, bytes.Equal(fixtLog.UnblessedMerkleRoots[0].MerkleRoot[:], out.UnblessedMerkleRoots[0].MerkleRoot[:]))
	})

	t.Run("decode hardcoded type onramp SendRequested success", func(t *testing.T) {
		fixtLog := getFixtureOnrampMessageLog()
		log, err := generateOnRampLog(fixtLog)
		require.NoError(t, err)

		var out chainaccessor.SendRequestedEvent
		err = decodeHardcodedType(&out, log)
		require.NoError(t, err)

		require.Equal(t, ccipocr3.SeqNum(fixtLog.SequenceNumber), out.SequenceNumber)
		require.Equal(t, ccipocr3.ChainSelector(fixtLog.DestChainSelector), out.DestChainSelector)
		require.Equal(t, true, bytes.Equal(fixtLog.Message.Data, out.Message.Data))
		require.Equal(t, true, bytes.Equal(fixtLog.Message.FeeToken.Bytes(), out.Message.FeeToken[:]))
		require.Equal(t, len(fixtLog.Message.TokenAmounts), len(out.Message.TokenAmounts))
	})
	t.Run("decode hardcoded tupe offramp ExecutionStateChange success", func(t *testing.T) {
		fixtLog := getFixtureExecStateChangedLog()
		log, err := generateOffRampStateChangeLog(fixtLog)
		require.NoError(t, err)

		var out reader.ExecutionStateChangedEvent
		err = decodeHardcodedType(&out, log)
		require.NoError(t, err)

		assert.Equal(t, fixtLog.GasUsed, &out.GasUsed)
		assert.Equal(t, ccipocr3.SeqNum(fixtLog.SequenceNumber), out.SequenceNumber)
		assert.Equal(t, true, bytes.Equal(fixtLog.MessageHash[:], out.MessageHash[:]))
		assert.Equal(t, true, bytes.Equal(fixtLog.ReturnData[:], out.ReturnData[:]))
	})
}
func generateOffRampStateChangeLog(log offramp.OffRampExecutionStateChanged) (*logpoller.Log, error) {
	event, ok := offrampABI.Events[executionStateChangedEvent]
	if !ok {
		return nil, fmt.Errorf("event not found %s", commitReportAcceptedEvent)
	}
	data, err := event.Inputs.NonIndexed().Pack(log.MessageHash, log.State, log.ReturnData, log.GasUsed)
	if err != nil {
		return nil, err
	}

	res := &logpoller.Log{}
	res.Data = data
	res.Topics = append(res.Topics, event.ID.Bytes())
	var buf [32]byte
	binary.BigEndian.PutUint64(buf[24:], log.SourceChainSelector) // put into last 8 bytes
	res.Topics = append(res.Topics, common.BytesToHash(buf[24:]).Bytes())
	binary.BigEndian.PutUint64(buf[24:], log.SequenceNumber) // put into last 8 bytes
	res.Topics = append(res.Topics, common.BytesToHash(buf[24:]).Bytes())
	res.Topics = append(res.Topics, log.MessageHash[:])

	return res, nil
}

func generateOfframpLog(log offramp.OffRampCommitReportAccepted) (*logpoller.Log, error) {
	event, ok := offrampABI.Events[commitReportAcceptedEvent]
	if !ok {
		return nil, fmt.Errorf("event not found %s", commitReportAcceptedEvent)
	}
	data, err := event.Inputs.NonIndexed().Pack(log.BlessedMerkleRoots, log.UnblessedMerkleRoots, log.PriceUpdates)
	if err != nil {
		return nil, err
	}

	res := &logpoller.Log{}
	res.Data = data
	res.Topics = append(res.Topics, event.ID.Bytes())

	return res, nil
}

func generateOnRampLog(log onramp.OnRampCCIPMessageSent) (*logpoller.Log, error) {
	event, ok := onrampABI.Events[ccipMessageSentEvent]
	if !ok {
		return nil, fmt.Errorf("event not found %s", ccipMessageSentEvent)
	}

	data, err := event.Inputs.NonIndexed().Pack(log.Message)
	if err != nil {
		return nil, err
	}

	res := &logpoller.Log{}
	res.Data = data
	res.Topics = append(res.Topics, event.ID.Bytes())
	var buf [32]byte
	binary.BigEndian.PutUint64(buf[24:], log.DestChainSelector) // put into last 8 bytes
	res.Topics = append(res.Topics, common.BytesToHash(buf[24:]).Bytes())
	binary.BigEndian.PutUint64(buf[24:], log.SequenceNumber) // put into last 8 bytes
	res.Topics = append(res.Topics, common.BytesToHash(buf[24:]).Bytes())

	return res, nil
}

func getFixtureOnrampMessageLog() onramp.OnRampCCIPMessageSent {
	return onramp.OnRampCCIPMessageSent{
		DestChainSelector: 1001,
		SequenceNumber:    1,
		Message: onramp.InternalEVM2AnyRampMessage{
			Header: onramp.InternalRampMessageHeader{
				MessageId:           [32]byte{1, 3},
				SourceChainSelector: 100,
				DestChainSelector:   200,
				SequenceNumber:      1,
				Nonce:               42,
			},
			Sender:         common.HexToAddress("0x1111111111111111111111111111111111111111"),
			Data:           []byte("mock data payload"),
			Receiver:       []byte{0xde, 0xad, 0xbe, 0xef},
			ExtraArgs:      []byte("extra args"),
			FeeToken:       common.HexToAddress("0x2222222222222222222222222222222222222222"),
			FeeTokenAmount: big.NewInt(1234560000000000000), // 1.23 ETH
			FeeValueJuels:  big.NewInt(9876543210),
			TokenAmounts: []onramp.InternalEVM2AnyTokenTransfer{
				{
					SourcePoolAddress: common.HexToAddress("0x3333333333333333333333333333333333333333"),
					DestTokenAddress:  []byte{0xca, 0xfe},
					ExtraData:         []byte("token extra"),
					Amount:            big.NewInt(1000000000000000000), // 1.0
					DestExecData:      []byte("exec call data"),
				},
			},
		},
	}
}

func getFixtureExecStateChangedLog() offramp.OffRampExecutionStateChanged {
	return offramp.OffRampExecutionStateChanged{
		SourceChainSelector: 10,
		SequenceNumber:      10,
		MessageId:           [32]byte{1, 2, 3},
		MessageHash:         [32]byte{5, 6, 7},
		State:               3,
		ReturnData:          []byte{1, 3},
		GasUsed:             big.NewInt(10),
	}
}

func getFixtureOfframpLog() offramp.OffRampCommitReportAccepted {
	var res offramp.OffRampCommitReportAccepted
	res.BlessedMerkleRoots = []offramp.InternalMerkleRoot{
		{
			SourceChainSelector: 1234,
			OnRampAddress:       bytes.Repeat([]byte{0x11}, 20),
			MinSeqNr:            1,
			MaxSeqNr:            10,
			MerkleRoot:          [32]byte{0xaa},
		},
	}
	res.UnblessedMerkleRoots = []offramp.InternalMerkleRoot{
		{
			SourceChainSelector: 1234,
			OnRampAddress:       bytes.Repeat([]byte{0x11}, 20),
			MinSeqNr:            1,
			MaxSeqNr:            10,
			MerkleRoot:          [32]byte{0xab},
		},
	}
	res.PriceUpdates.TokenPriceUpdates = []offramp.InternalTokenPriceUpdate{

		{
			SourceToken: common.HexToAddress("0x2222222222222222222222222222222222222222"),
			UsdPerToken: big.NewInt(1e18),
		},
	}
	res.PriceUpdates.GasPriceUpdates = []offramp.InternalGasPriceUpdate{
		{
			DestChainSelector: 5678,
			UsdPerUnitGas:     big.NewInt(2e18),
		},
	}

	return res
}
