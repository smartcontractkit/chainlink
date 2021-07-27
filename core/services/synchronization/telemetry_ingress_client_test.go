package synchronization_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/services/synchronization/mocks"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/smartcontractkit/wsrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestTelemetryIngressClient_Send_HappyPath(t *testing.T) {

	// Create mocks
	wsrpcDial := new(mocks.WsrpcDial)
	newTelemClient := new(mocks.NewTelemClient)
	telemClient := new(mocks.TelemClient)
	csaKeystore := new(ksmocks.CSAKeystoreInterface)

	// Set mock handlers for keystore
	clientPubKey, _ := crypto.PublicKeyFromHex("1111111111111111111111111111111111111111111111111111111111111111")
	csaKeyList := []csakey.Key{{PublicKey: *clientPubKey}}
	clientPrivKey := []byte("222222222")
	csaKeystore.On("ListCSAKeys").Return(csaKeyList, nil)
	csaKeystore.On("Unsafe_GetUnlockedPrivateKey", *clientPubKey).Return(clientPrivKey, nil)

	// Set mock handlers for wsrpcDial
	conn := wsrpc.ClientConn{}
	wsrpcDial.On("Execute", mock.Anything, mock.Anything).Return(&conn, nil)

	// Set mock handlers on the new telem client
	newTelemClient.On("Execute", mock.Anything).Return(telemClient)

	// Wire up the telem ingress client
	url := &url.URL{}
	serverPubKeyHex := "33333333333"
	telemIngressClient := synchronization.NewTelemetryIngressClient(url, serverPubKeyHex, wsrpcDial.Execute, newTelemClient.Execute, csaKeystore, false)
	require.NoError(t, telemIngressClient.Start())

	// Create the telemetry payload
	telemetry := []byte("101010")
	address := common.HexToAddress("0xa")
	telemPayload := synchronization.TelemPayload{
		Ctx:             context.Background(),
		Telemetry:       telemetry,
		ContractAddress: address,
	}

	// Assert the telemetry payload is correctly sent to wsrpc
	called := abool.New()
	telemClient.On("Telem", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		called.Set()
		telemReq := args.Get(1).(*telemPb.TelemRequest)
		assert.Equal(t, telemPayload.ContractAddress.String(), telemReq.Address)
		assert.Equal(t, telemPayload.Telemetry, telemReq.Telemetry)
	})

	// Give the goroutines time to spin up
	time.Sleep(1 * time.Second)

	// Send telemetry (and do so continuously in case the sleep wasn't long enough)
	for called.IsSet() == false {
		telemIngressClient.Send(telemPayload)
		time.Sleep(10 * time.Millisecond)
	}

	// FIXME: The above `conn` isn't mocked, so calling close panics
	// defer telemIngressClient.Close()
}
