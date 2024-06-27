package logprovider

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

type LogDataPacker interface {
	PackLogData(log logpoller.Log) ([]byte, error)
}

type logEventsPacker struct {
	abi abi.ABI
}

func NewLogEventsPacker() *logEventsPacker {
	return &logEventsPacker{abi: core.UtilsABI}
}

func (p *logEventsPacker) PackLogData(log logpoller.Log) ([]byte, error) {
	var topics [][32]byte
	for _, topic := range log.GetTopics() {
		topics = append(topics, topic)
	}
	b, err := p.abi.Methods["_log"].Inputs.Pack(&automation_utils_2_1.Log{
		Index:       big.NewInt(log.LogIndex),
		Timestamp:   big.NewInt(log.BlockTimestamp.Unix()),
		TxHash:      log.TxHash,
		BlockNumber: big.NewInt(log.BlockNumber),
		BlockHash:   log.BlockHash,
		Source:      log.Address,
		Topics:      topics,
		Data:        log.Data,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}
