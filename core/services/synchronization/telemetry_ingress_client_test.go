package synchronization_test

import (
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	telemPb "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	telem "github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

func TestTelemetryIngressClient_Send_HappyPath(t *testing.T) {
	// Create mocks
	telemClient := mocks.NewTelemClient(t)

	csaKeystore := new(ksmocks.CSA)

	// Set mock handlers for keystore
	key := cltest.DefaultCSAKey
	keyList := []csakey.KeyV2{key}
	csaKeystore.On("EnsureKey", mock.Anything).Return(nil)
	csaKeystore.On("GetAll").Return(keyList, nil)
	csaKeystore.On("Get", key.ID()).Return(key, nil)

	// Wire up the telem ingress client
	url := &url.URL{}
	serverPubKeyHex := telem.GetDummyKeyString()
	telemIngressClient := synchronization.NewTestTelemetryIngressClient(t, url, serverPubKeyHex, csaKeystore, telemClient)
	servicetest.Run(t, telemIngressClient)

	// Create the telemetry payload
	telemetry := []byte("101010")
	address := common.HexToAddress("0xa")
	telemPayload := synchronization.TelemPayload{
		Telemetry:  telemetry,
		ContractID: address.String(),
		TelemType:  synchronization.OCR,
	}

	// Assert the telemetry payload is correctly sent to wsrpc
	var called atomic.Bool
	telemClient.On("Telem", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		called.Store(true)
		telemReq := args.Get(1).(*telemPb.TelemRequest)
		assert.Equal(t, telemPayload.ContractID, telemReq.Address)
		assert.Equal(t, telemPayload.Telemetry, telemReq.Telemetry)
		assert.Equal(t, string(synchronization.OCR), telemReq.TelemetryType)
		assert.Positive(t, telemReq.SentAt)
	})

	// Send telemetry
	telemIngressClient.Send(testutils.Context(t), telemPayload.Telemetry, telemPayload.ContractID, telemPayload.TelemType)

	// Wait for the telemetry to be handled
	gomega.NewWithT(t).Eventually(called.Load).Should(gomega.BeTrue())
}
