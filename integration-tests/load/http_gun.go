package load

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

var ReportTypes = GetReportTypes()

func mustNewType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}

func GetReportTypes() abi.Arguments {
	return []abi.Argument{
		{Name: "feedId", Type: mustNewType("bytes32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "observationsBlocknumber", Type: mustNewType("uint64")},
		{Name: "median", Type: mustNewType("int192")},
	}
}

type MercuryHTTPGun struct {
	BaseURL   string
	client    *client.MercuryServer
	netclient blockchain.EVMClient
	feedID    string
	bn        atomic.Uint64
}

func NewHTTPGun(baseURL string, client *client.MercuryServer, feedID string, bn uint64) *MercuryHTTPGun {
	g := &MercuryHTTPGun{
		BaseURL: baseURL,
		client:  client,
		feedID:  feedID,
	}
	g.bn.Store(bn)
	return g
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *MercuryHTTPGun) Call(l *loadgen.Generator) loadgen.CallResult {
	tn := time.Now()
	answer, res, err := m.client.GetReports(m.feedID, m.bn.Load())
	if err != nil {
		return loadgen.CallResult{Error: "connection error"}
	}
	if res.Status != "200 OK" {
		return loadgen.CallResult{Error: "not 200"}
	}
	log.Info().Dur("Elapsed", time.Since(tn)).Send()
	reportElements := map[string]interface{}{}
	if err = ReportTypes.UnpackIntoMap(reportElements, []byte(answer.ChainlinkBlob)); err != nil {
		return loadgen.CallResult{Error: "blob unpacking error"}
	}
	if err := testsetups.ValidateReport(reportElements); err != nil {
		return loadgen.CallResult{Error: "report validation error"}
	}
	return loadgen.CallResult{}
}
