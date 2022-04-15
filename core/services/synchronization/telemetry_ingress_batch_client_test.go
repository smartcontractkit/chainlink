package synchronization_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/services/synchronization/mocks"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
)

func TestTelemetryIngressBatchClient_HappyPath(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create mocks
	telemClient := new(mocks.TelemClient)
	csaKeystore := new(ksmocks.CSA)

	// Set mock handlers for keystore
	key := cltest.DefaultCSAKey
	keyList := []csakey.KeyV2{key}
	csaKeystore.On("GetAll").Return(keyList, nil)

	// Wire up the telem ingress client
	url := &url.URL{}
	serverPubKeyHex := "33333333333"
	sendInterval := time.Millisecond * 5
	telemIngressClient := synchronization.NewTestTelemetryIngressBatchClient(t, url, serverPubKeyHex, csaKeystore, false, telemClient, sendInterval, false)
	require.NoError(t, telemIngressClient.Start(testutils.Context(t)))

	// Create telemetry payloads for different contracts
	telemPayload1 := synchronization.TelemPayload{
		Ctx:        testutils.Context(t),
		Telemetry:  []byte("Mock telem 1"),
		ContractID: "0x1",
	}
	telemPayload2 := synchronization.TelemPayload{
		Ctx:        testutils.Context(t),
		Telemetry:  []byte("Mock telem 2"),
		ContractID: "0x2",
	}
	telemPayload3 := synchronization.TelemPayload{
		Ctx:        testutils.Context(t),
		Telemetry:  []byte("Mock telem 3"),
		ContractID: "0x3",
	}

	// Assert telemetry payloads for each contract are correctly sent to wsrpc
	var contractCounter1 atomic.Uint32
	var contractCounter2 atomic.Uint32
	var contractCounter3 atomic.Uint32
	telemClient.On("TelemBatch", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		telemBatchReq := args.Get(1).(*telemPb.TelemBatchRequest)

		if telemBatchReq.ContractId == "0x1" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter1.Inc()
				assert.Equal(t, telemPayload1.Telemetry, telem)
			}
		}
		if telemBatchReq.ContractId == "0x2" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter2.Inc()
				assert.Equal(t, telemPayload2.Telemetry, telem)
			}
		}
		if telemBatchReq.ContractId == "0x3" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter3.Inc()
				assert.Equal(t, telemPayload3.Telemetry, telem)
			}
		}
	})

	// Send telemetry
	telemIngressClient.Send(telemPayload1)
	telemIngressClient.Send(telemPayload2)
	telemIngressClient.Send(telemPayload3)
	time.Sleep(sendInterval * 2)
	telemIngressClient.Send(telemPayload1)
	telemIngressClient.Send(telemPayload1)
	telemIngressClient.Send(telemPayload2)

	// Wait for the telemetry to be handled
	g.Eventually(func() []uint32 {
		return []uint32{contractCounter1.Load(), contractCounter2.Load(), contractCounter3.Load()}
	}).Should(gomega.Equal([]uint32{3, 2, 1}))

	// Client should shut down
	telemIngressClient.Close()
}
