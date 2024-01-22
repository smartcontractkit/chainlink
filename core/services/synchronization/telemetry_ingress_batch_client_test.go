package synchronization_test

import (
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	telemPb "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

func TestTelemetryIngressBatchClient_HappyPath(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create mocks
	telemClient := mocks.NewTelemClient(t)
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
	servicetest.Run(t, telemIngressClient)

	// Create telemetry payloads for different contracts
	telemPayload1 := synchronization.TelemPayload{
		Telemetry:  []byte("Mock telem 1"),
		ContractID: "0x1",
		TelemType:  synchronization.OCR,
	}
	telemPayload2 := synchronization.TelemPayload{
		Telemetry:  []byte("Mock telem 2"),
		ContractID: "0x2",
		TelemType:  synchronization.OCR2VRF,
	}
	telemPayload3 := synchronization.TelemPayload{
		Telemetry:  []byte("Mock telem 3"),
		ContractID: "0x3",
		TelemType:  synchronization.OCR2Functions,
	}

	// Assert telemetry payloads for each contract are correctly sent to wsrpc
	var contractCounter1 atomic.Uint32
	var contractCounter2 atomic.Uint32
	var contractCounter3 atomic.Uint32
	telemClient.On("TelemBatch", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		telemBatchReq := args.Get(1).(*telemPb.TelemBatchRequest)

		if telemBatchReq.ContractId == "0x1" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter1.Add(1)
				assert.Equal(t, telemPayload1.Telemetry, telem)
				assert.Equal(t, synchronization.OCR, telemPayload1.TelemType)
			}
		}
		if telemBatchReq.ContractId == "0x2" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter2.Add(1)
				assert.Equal(t, telemPayload2.Telemetry, telem)
				assert.Equal(t, synchronization.OCR2VRF, telemPayload2.TelemType)
			}
		}
		if telemBatchReq.ContractId == "0x3" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter3.Add(1)
				assert.Equal(t, telemPayload3.Telemetry, telem)
				assert.Equal(t, synchronization.OCR2Functions, telemPayload3.TelemType)
			}
		}
	})

	// Send telemetry
	testCtx := testutils.Context(t)
	telemIngressClient.Send(testCtx, telemPayload1.Telemetry, telemPayload1.ContractID, telemPayload1.TelemType)
	telemIngressClient.Send(testCtx, telemPayload2.Telemetry, telemPayload2.ContractID, telemPayload2.TelemType)
	telemIngressClient.Send(testCtx, telemPayload3.Telemetry, telemPayload3.ContractID, telemPayload3.TelemType)
	time.Sleep(sendInterval * 2)
	telemIngressClient.Send(testCtx, telemPayload1.Telemetry, telemPayload1.ContractID, telemPayload1.TelemType)
	telemIngressClient.Send(testCtx, telemPayload1.Telemetry, telemPayload1.ContractID, telemPayload1.TelemType)
	telemIngressClient.Send(testCtx, telemPayload2.Telemetry, telemPayload2.ContractID, telemPayload2.TelemType)

	// Wait for the telemetry to be handled
	g.Eventually(func() []uint32 {
		return []uint32{contractCounter1.Load(), contractCounter2.Load(), contractCounter3.Load()}
	}).Should(gomega.Equal([]uint32{3, 2, 1}))
}
