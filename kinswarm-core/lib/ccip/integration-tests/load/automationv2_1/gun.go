package automationv2_1

import (
	"math/big"
	"sync"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

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

	resultCh := make(chan *wasp.Response, len(dividedData))

	for _, a := range dividedData {
		wg.Add(1)
		go func(a [][]byte, m *LogTriggerGun) {
			defer wg.Done()

			_, err := contracts.MultiCallLogTriggerLoadGen(m.client, m.multiCallAddress, m.addresses, a)
			if err != nil {
				m.logger.Error().Err(err).Msg("Error calling MultiCallLogTriggerLoadGen")
				resultCh <- &wasp.Response{Error: err.Error(), Failed: true}
				return
			}
			resultCh <- &wasp.Response{}
		}(a, m)
	}
	wg.Wait()
	close(resultCh)

	r := &wasp.Response{}
	for result := range resultCh {
		if result.Failed {
			r.Failed = true
			if r.Error != "" {
				r.Error += "; " + result.Error
			} else {
				r.Error = result.Error
			}
		}
	}

	return r
}
