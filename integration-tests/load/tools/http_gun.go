package tools

import (
	"encoding/hex"
	"math/rand"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	mercurysetup "github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

type MercuryHTTPGun struct {
	BaseURL   string
	client    *client.MercuryServer
	netclient blockchain.EVMClient
	feeds     [][32]byte
	Bn        atomic.Uint64
}

func NewHTTPGun(baseURL string, client *client.MercuryServer, feeds [][32]byte, bn uint64) *MercuryHTTPGun {
	g := &MercuryHTTPGun{
		BaseURL: baseURL,
		client:  client,
		feeds:   feeds,
	}
	g.Bn.Store(bn)
	return g
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *MercuryHTTPGun) Call(l *loadgen.Generator) loadgen.CallResult {
	randFeedIdStr := mercurysetup.Byte32ToString(m.feeds[rand.Intn(len(m.feeds))])
	answer, res, err := m.client.GetReportsByFeedIdStr(randFeedIdStr, m.Bn.Load())
	if err != nil {
		return loadgen.CallResult{Error: "connection error", Failed: true}
	}
	if res.Status != "200 OK" {
		return loadgen.CallResult{Error: "not 200", Failed: true}
	}
	reportBytes, err := hex.DecodeString(answer.ChainlinkBlob[2:])
	if err != nil {
		return loadgen.CallResult{Error: "report validation error", Failed: true}
	}
	report, err := mercury.DecodeReport(reportBytes)
	_ = report
	if err != nil {
		return loadgen.CallResult{Error: "report validation error", Failed: true}
	}
	return loadgen.CallResult{}
}
