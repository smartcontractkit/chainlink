package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

func TestUpkeepFilter_Select(t *testing.T) {
	var zeroBytes [32]byte
	emptyTopic := common.BytesToHash(zeroBytes[:])
	contractAddress := common.HexToAddress("0xB9F3af0c2CbfE108efd0E23F7b0a151Ea42f764E")
	uid := big.NewInt(123456)
	topic10 := "0x000000000000000000000000000000000000000000000000000000000000007b" // decimal 123 encoded
	topic20 := "0x0000000000000000000000000000000000000000000000000000000000000001" // bool true encoded
	topic30 := "0x00000000000000000000000082b8b466f4be252e56af8a00aa28838866686062" // address encoded
	topic11 := "0x000000000000000000000000000000000000000000000000000000000000007a" // decimal 122 encoded
	topic21 := "0x0000000000000000000000000000000000000000000000000000000000000000" // bool false encoded
	topic31 := "0x000000000000000000000000f91a27d2f37a36f1e6acc681b07b1dd2e288aebc" // address encoded

	log1 := logpoller.Log{
		Topics: [][]byte{
			contractAddress.Bytes(),
			hexutil.MustDecode(topic10),
			hexutil.MustDecode(topic20),
			hexutil.MustDecode(topic30),
		},
	}
	log2 := logpoller.Log{
		Topics: [][]byte{
			contractAddress.Bytes(),
			hexutil.MustDecode(topic11),
			hexutil.MustDecode(topic21),
			hexutil.MustDecode(topic31),
		},
	}
	log3 := logpoller.Log{
		Topics: [][]byte{
			contractAddress.Bytes(),
			hexutil.MustDecode(topic11),
			hexutil.MustDecode(topic20),
			hexutil.MustDecode(topic31),
		},
	}
	log4 := logpoller.Log{
		Topics: [][]byte{
			contractAddress.Bytes(),
			hexutil.MustDecode(topic10),
			hexutil.MustDecode(topic21),
			hexutil.MustDecode(topic31),
		},
	}
	log5 := logpoller.Log{
		Topics: [][]byte{
			contractAddress.Bytes(),
			hexutil.MustDecode(topic10),
			hexutil.MustDecode(topic20),
			hexutil.MustDecode(topic30),
		},
	}

	tests := []struct {
		name         string
		filter       upkeepFilter
		logs         []logpoller.Log
		expectedLogs []logpoller.Log
	}{
		{
			"no selector configured - all logs are returned",
			upkeepFilter{
				selector: 0,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), emptyTopic, emptyTopic, emptyTopic},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
			},
			[]logpoller.Log{
				log1,
				log2,
			},
		},
		{
			"selector is 1 - topics 1 is used to filter logs",
			upkeepFilter{
				selector: 1,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), common.HexToHash(topic10), emptyTopic, emptyTopic},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
			},
			[]logpoller.Log{
				log1,
			},
		},
		{
			"selector is 2 - topic 2 is used to filter logs",
			upkeepFilter{
				selector: 2,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), emptyTopic, common.HexToHash(topic21), emptyTopic},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
			},
			[]logpoller.Log{
				log2,
			},
		},
		{
			"selector is 3 - topics 1 2 are used to filter logs",
			upkeepFilter{
				selector: 3,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), common.HexToHash(topic10), common.HexToHash(topic21), emptyTopic},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
				log3,
				log4,
			},
			[]logpoller.Log{
				log4,
			},
		},
		{
			"selector is 4 - topic 3 is used to filter logs",
			upkeepFilter{
				selector: 4,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), emptyTopic, emptyTopic, common.HexToHash(topic31)},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
			},
			[]logpoller.Log{
				log2,
			},
		},
		{
			"selector is 5 - topics 1 3 are used to filter logs",
			upkeepFilter{
				selector: 5,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), common.HexToHash(topic11), emptyTopic, common.HexToHash(topic31)},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
				log3,
				log4,
			},
			[]logpoller.Log{
				log2,
				log3,
			},
		},
		{
			"selector is 7 - topics 1 2 3 are used to filter logs",
			upkeepFilter{
				selector: 7,
				topics:   []common.Hash{common.BytesToHash(contractAddress.Bytes()), common.HexToHash(topic10), common.HexToHash(topic20), common.HexToHash(topic30)},
				upkeepID: uid,
			},
			[]logpoller.Log{
				log1,
				log2,
				log3,
				log4,
				log5,
			},
			[]logpoller.Log{
				log1,
				log5,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filteredLogs := tc.filter.Select(tc.logs...)
			assert.Equal(t, tc.expectedLogs, filteredLogs)
		})
	}
}
