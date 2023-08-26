package logprovider

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestLogEventProvider_LifeCycle(t *testing.T) {
	tests := []struct {
		name           string
		errored        bool
		upkeepID       *big.Int
		upkeepCfg      LogTriggerConfig
		cfgUpdateBlock uint64
		mockPoller     bool
		unregister     bool
	}{
		{
			"new upkeep",
			false,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			uint64(1),
			true,
			false,
		},
		{
			"empty config",
			true,
			big.NewInt(111),
			LogTriggerConfig{},
			uint64(0),
			false,
			false,
		},
		{
			"invalid config",
			true,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{}, 32)),
			},
			uint64(2),
			false,
			false,
		},
		{
			"existing config",
			true,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			uint64(0),
			true,
			false,
		},
		{
			"existing config with newer block",
			false,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			uint64(2),
			true,
			true,
		},
	}

	mp := new(mocks.LogPoller)
	mp.On("RegisterFilter", mock.Anything).Return(nil)
	mp.On("UnregisterFilter", mock.Anything).Return(nil)
	mp.On("ReplayAsync", mock.Anything).Return(nil)
	p := NewLogProvider(logger.TestLogger(t), mp, &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := p.RegisterFilter(FilterOptions{
				UpkeepID:      tc.upkeepID,
				TriggerConfig: tc.upkeepCfg,
				UpdateBlock:   tc.cfgUpdateBlock,
			})
			if tc.errored {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.unregister {
					require.NoError(t, p.UnregisterFilter(tc.upkeepID))
				}
			}
		})
	}
}

func TestEventLogProvider_RefreshActiveUpkeeps(t *testing.T) {
	mp := new(mocks.LogPoller)
	mp.On("RegisterFilter", mock.Anything).Return(nil)
	mp.On("UnregisterFilter", mock.Anything).Return(nil)
	mp.On("ReplayAsync", mock.Anything).Return(nil)

	p := NewLogProvider(logger.TestLogger(t), mp, &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200))

	require.NoError(t, p.RegisterFilter(FilterOptions{
		UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1111").BigInt(),
		TriggerConfig: LogTriggerConfig{
			ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
			Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
		},
		UpdateBlock: uint64(0),
	}))
	require.NoError(t, p.RegisterFilter(FilterOptions{
		UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2222").BigInt(),
		TriggerConfig: LogTriggerConfig{
			ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
			Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
		},
		UpdateBlock: uint64(0),
	}))
	require.Equal(t, 2, p.filterStore.Size())

	newIds, err := p.RefreshActiveUpkeeps()
	require.NoError(t, err)
	require.Len(t, newIds, 0)
	newIds, err = p.RefreshActiveUpkeeps(
		core.GenUpkeepID(ocr2keepers.LogTrigger, "2222").BigInt(),
		core.GenUpkeepID(ocr2keepers.LogTrigger, "1234").BigInt(),
		core.GenUpkeepID(ocr2keepers.LogTrigger, "123").BigInt())
	require.NoError(t, err)
	require.Len(t, newIds, 2)
	require.Equal(t, 1, p.filterStore.Size())
}

func TestLogEventProvider_ValidateLogTriggerConfig(t *testing.T) {
	contractAddress := common.HexToAddress("0xB9F3af0c2CbfE108efd0E23F7b0a151Ea42f764E")
	eventSig := common.HexToHash("0x3bdab8bffae631cfee411525ebae27f3fb61b10c662c09ec2a7dbb5854c87e8c")
	tests := []struct {
		name        string
		cfg         LogTriggerConfig
		expectedErr error
	}{
		{
			"success",
			LogTriggerConfig{
				ContractAddress: contractAddress,
				FilterSelector:  0,
				Topic0:          eventSig,
			},
			nil,
		},
		{
			"invalid contract address",
			LogTriggerConfig{
				ContractAddress: common.Address{},
				FilterSelector:  0,
				Topic0:          eventSig,
			},
			fmt.Errorf("invalid contract address: zeroed"),
		},
		{
			"invalid topic0",
			LogTriggerConfig{
				ContractAddress: contractAddress,
				FilterSelector:  0,
			},
			fmt.Errorf("invalid topic0: zeroed"),
		},
		{
			"success",
			LogTriggerConfig{
				ContractAddress: contractAddress,
				FilterSelector:  8,
				Topic0:          eventSig,
			},
			fmt.Errorf("invalid filter selector: larger or equal to 8"),
		},
	}

	p := NewLogProvider(logger.TestLogger(t), nil, &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200))
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := p.validateLogTriggerConfig(tc.cfg)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestLogEventProvider_FilterLogsByContent(t *testing.T) {
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
				topics:   []common.Hash{contractAddress.Hash(), emptyTopic, emptyTopic, emptyTopic},
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
			"selector is 1 - only topics[1] is used to filter logs",
			upkeepFilter{
				selector: 1,
				topics:   []common.Hash{contractAddress.Hash(), common.HexToHash(topic10), emptyTopic, emptyTopic},
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
			"selector is 2 - only topics[2] is used to filter logs",
			upkeepFilter{
				selector: 2,
				topics:   []common.Hash{contractAddress.Hash(), emptyTopic, common.HexToHash(topic21), emptyTopic},
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
			"selector is 3 - topics[1] and [2] are used to filter logs",
			upkeepFilter{
				selector: 3,
				topics:   []common.Hash{contractAddress.Hash(), common.HexToHash(topic10), common.HexToHash(topic21), emptyTopic},
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
			"selector is 4 - only topics[3] is used to filter logs",
			upkeepFilter{
				selector: 4,
				topics:   []common.Hash{contractAddress.Hash(), emptyTopic, emptyTopic, common.HexToHash(topic31)},
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
			"selector is 5 - both topics[1] and [3] are used to filter logs",
			upkeepFilter{
				selector: 5,
				topics:   []common.Hash{contractAddress.Hash(), common.HexToHash(topic11), emptyTopic, common.HexToHash(topic31)},
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
			"selector is 7 - topics[1], [2] and [3] are used to filter logs",
			upkeepFilter{
				selector: 7,
				topics:   []common.Hash{contractAddress.Hash(), common.HexToHash(topic10), common.HexToHash(topic20), common.HexToHash(topic30)},
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

	p := NewLogProvider(logger.TestLogger(t), nil, &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200))
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filteredLogs := p.filterLogsByContent(tc.filter, tc.logs)
			assert.Equal(t, tc.expectedLogs, filteredLogs)
		})
	}
}
