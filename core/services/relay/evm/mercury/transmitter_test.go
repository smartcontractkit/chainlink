package mercury

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/report"
)

type MockWSRPCClient struct {
	transmit func(ctx context.Context, in *report.ReportRequest) (*report.ReportResponse, error)
}

func (m MockWSRPCClient) Name() string                   { return "" }
func (m MockWSRPCClient) Start(context.Context) error    { return nil }
func (m MockWSRPCClient) Close() error                   { return nil }
func (m MockWSRPCClient) Healthy() error                 { return nil }
func (m MockWSRPCClient) HealthReport() map[string]error { return map[string]error{} }
func (m MockWSRPCClient) Ready() error                   { return nil }
func (m MockWSRPCClient) Transmit(ctx context.Context, in *report.ReportRequest) (*report.ReportResponse, error) {
	return m.transmit(ctx, in)
}

var _ wsrpc.Client = &MockWSRPCClient{}

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	reportURL := "http://report.test/foo"
	username := "my username"
	password := "my password"

	t.Run("successful transmit", func(t *testing.T) {
		c := MockWSRPCClient{
			transmit: func(ctx context.Context, in *report.ReportRequest) (out *report.ReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, samplePayloadHex, hexutil.Encode(in.Payload))
				out = new(report.ReportResponse)
				out.Code = 42
				out.Error = ""
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, c, sampleFromAccount, reportURL, username, password)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.NoError(t, err)
	})

	t.Run("failing transmit", func(t *testing.T) {
		c := MockWSRPCClient{
			transmit: func(ctx context.Context, in *report.ReportRequest) (out *report.ReportResponse, err error) {
				return nil, errors.New("foo error")
			},
		}
		mt := NewTransmitter(lggr, c, sampleFromAccount, reportURL, username, password)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "foo error")
	})
}
