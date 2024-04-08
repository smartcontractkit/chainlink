package automationv2_1

import (
	"math/big"
	"sync"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/smartcontractkit/wasp"

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
	data             [][]byte
	addresses        []string
	multiCallAddress string
	client           *seth.Client
	numberOfClients  int
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
	client *seth.Client,
	numberOfClients int,
	multicallAddress string,
) *LogTriggerGun {
	var data [][]byte
	var addresses []string

	for _, c := range TriggerConfigs {
		if c.NumberOfEvents > 0 {
			d := generateCallData(1, 1, c.NumberOfEvents)
			data = append(data, d)
			addresses = append(addresses, c.Address)
		}
		if c.NumberOfSpamMatchingEvents > 0 {
			d := generateCallData(1, 2, c.NumberOfSpamMatchingEvents)
			data = append(data, d)
			addresses = append(addresses, c.Address)
		}
		if c.NumberOfSpamNonMatchingEvents > 0 {
			d := generateCallData(2, 2, c.NumberOfSpamNonMatchingEvents)
			data = append(data, d)
			addresses = append(addresses, c.Address)
		}
	}

	return &LogTriggerGun{
		addresses:        addresses,
		data:             data,
		logger:           logger,
		multiCallAddress: multicallAddress,
		client:           client,
		numberOfClients:  numberOfClients,
	}
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	var wg sync.WaitGroup
	var dividedData [][][]byte
	d := m.data
	chunkSize := 100
	for i := 0; i < len(d); i += chunkSize {
		end := i + chunkSize
		if end > len(d) {
			end = len(d)
		}
		dividedData = append(dividedData, d[i:end])
	}

	// semaphoreCh := make(chan struct{}, m.numberOfClients)

	for i, a := range dividedData {
		wg.Add(1)
		// semaphoreCh <- struct{}{}
		go func(a [][]byte, m *LogTriggerGun, i int) *wasp.Response {
			defer wg.Done()
			// defer func() { <-semaphoreCh }()

			clientIndex := i + 1
			_, err := contracts.MultiCallLogTriggerLoadGen(m.client, clientIndex, m.multiCallAddress, m.addresses, a)
			if err != nil {
				m.logger.Error().Err(err).Msg("Error calling MultiCallLogTriggerLoadGen")
				return &wasp.Response{Error: err.Error(), Failed: true}
			}
			return &wasp.Response{}
		}(a, m, i)
	}
	wg.Wait()
	return &wasp.Response{}
}
