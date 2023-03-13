package tools

import (
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

type MercuryHTTPGun struct {
	BaseURL   string
	client    *client.MercuryServer
	netclient blockchain.EVMClient
	feedID    string
	Bn        atomic.Uint64
}

func NewHTTPGun(baseURL string, client *client.MercuryServer, feedID string, bn uint64) *MercuryHTTPGun {
	g := &MercuryHTTPGun{
		BaseURL: baseURL,
		client:  client,
		feedID:  feedID,
	}
	g.Bn.Store(bn)
	return g
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *MercuryHTTPGun) Call(l *loadgen.Generator) loadgen.CallResult {
	answer, res, err := m.client.GetReports(m.feedID, m.Bn.Load())
	if err != nil {
		return loadgen.CallResult{Error: "connection error"}
	}
	if res.Status != "200 OK" {
		return loadgen.CallResult{Error: "not 200"}
	}
	if err := mercury.ValidateReport([]byte(answer.ChainlinkBlob)); err != nil {
		return loadgen.CallResult{Error: "report validation error"}
	}
	return loadgen.CallResult{}
}
