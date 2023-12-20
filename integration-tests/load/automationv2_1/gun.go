package automationv2_1

import (
	"math/big"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type LogTriggerConfig struct {
	Address                       string
	NumberOfEvents                int64
	NumberOfSpamMatchingEvents    int64
	NumberOfSpamNonMatchingEvents int64
}

type LogTriggerGun struct {
	data             [][][]byte
	addresses        []string
	multiCallAddress string
	evmClient        blockchain.EVMClient
	logger           zerolog.Logger
}

func generateCallData(int1 int64, int2 int64, count int64) []byte {
	abi, err := log_emitter.LogEmitterMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	data, err := abi.Pack("EmitLog4", big.NewInt(int1), big.NewInt(int2), big.NewInt(count))
	if err != nil {
		panic(err)
	}
	return data
}

func NewLogTriggerUser(
	logger zerolog.Logger,
	TriggerConfigs []LogTriggerConfig,
	evmClient blockchain.EVMClient,
	multicallAddress string,
) *LogTriggerGun {

	var data1, data2, data3 [][]byte
	var addresses []string

	for _, c := range TriggerConfigs {
		if c.NumberOfEvents > 0 {
			d := generateCallData(1, 1, c.NumberOfEvents)
			data1 = append(data1, d)
			addresses = append(addresses, c.Address)
		}
		if c.NumberOfSpamMatchingEvents > 0 {
			d := generateCallData(1, 2, c.NumberOfSpamMatchingEvents)
			data2 = append(data2, d)
			addresses = append(addresses, c.Address)
		}
		if c.NumberOfSpamNonMatchingEvents > 0 {
			d := generateCallData(2, 2, c.NumberOfSpamNonMatchingEvents)
			data3 = append(data3, d)
			addresses = append(addresses, c.Address)
		}
	}

	var data [][][]byte
	data = append(data, data1)
	data = append(data, data2)
	data = append(data, data3)

	return &LogTriggerGun{
		addresses:        addresses,
		data:             data,
		logger:           logger,
		multiCallAddress: multicallAddress,
		evmClient:        evmClient,
	}
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {

	for _, d := range m.data {
		var dividedData [][][]byte
		chuckSize := 100
		for i := 0; i < len(d); i += chuckSize {
			end := i + chuckSize
			if end > len(d) {
				end = len(d)
			}
			dividedData = append(dividedData, d[i:end])

		}

		for _, a := range dividedData {
			_, err := contracts.MultiCallLogTriggerLoadGen(m.evmClient, m.multiCallAddress, m.addresses, a)
			if err != nil {
				return &wasp.Response{Error: err.Error(), Failed: true}
			}
		}
	}

	return &wasp.Response{}
}
