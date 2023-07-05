package logprovider

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_log_automation"
)

type LogDataPacker interface {
	PackLogData(log logpoller.Log) ([]byte, error)
}

type logEventsPacker struct {
	abi abi.ABI
}

func NewLogEventsPacker(logDataABI abi.ABI) *logEventsPacker {
	return &logEventsPacker{abi: logDataABI}
}

func (p *logEventsPacker) PackLogData(log logpoller.Log) ([]byte, error) {
	topics := [][32]byte{}
	for _, topic := range log.GetTopics() {
		topics = append(topics, topic)
	}
	b, err := p.abi.Pack("checkLog", &i_log_automation.Log{
		Index:       big.NewInt(log.LogIndex),
		TxIndex:     big.NewInt(0), // TODO
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
	return b[4:], nil
}
