package synchronization_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

func TestTelemetryIngressClient_Send_HappyPath(t *testing.T) {

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
	telemIngressClient := synchronization.NewTestTelemetryIngressClient(t, url, serverPubKeyHex, csaKeystore, false, telemClient)
	require.NoError(t, telemIngressClient.Start(testutils.Context(t)))
	defer telemIngressClient.Close()

	// Create the telemetry payload
	telemetry := []byte("101010")
	address := common.HexToAddress("0xa")
	telemPayload := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  telemetry,
		ContractID: address.String(),
	}

	// Assert the telemetry payload is correctly sent to wsrpc
	var called atomic.Bool
	telemClient.On("Telem", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		called.Store(true)
		telemReq := args.Get(1).(*telemPb.TelemRequest)
		assert.Equal(t, telemPayload.ContractID, telemReq.Address)
		assert.Equal(t, telemPayload.Telemetry, telemReq.Telemetry)
	})

	// Send telemetry
	telemIngressClient.Send(telemPayload)

	// Wait for the telemetry to be handled
	gomega.NewWithT(t).Eventually(called.Load).Should(gomega.BeTrue())
}
