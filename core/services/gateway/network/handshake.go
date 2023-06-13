package network

import (
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
)

// The handshake works as follows:
//
//	  Client (Initiator)                  Server (Acceptor)
//
//	 NewAuthHeader()
//	             -------auth header-------->
//	                                       StartHandshake()
//	             <-------challenge----------
//	ChallengeResponse()
//	             ---------response--------->
//	                                     FinalizeHandshake()
//
//go:generate mockery --quiet --name ConnectionInitiator --output ./mocks/ --case=underscore
type ConnectionInitiator interface {
	// Generate authentication header value specific to node and gateway
	NewAuthHeader(url *url.URL) ([]byte, error)

	// Sign challenge to prove identity.
	ChallengeResponse(challenge []byte) ([]byte, error)
}

//go:generate mockery --quiet --name ConnectionAcceptor --output ./mocks/ --case=underscore
type ConnectionAcceptor interface {
	// Verify auth header, save state of the attempt and generate a challenge for the node.
	StartHandshake(authHeader []byte) (attemptId string, challenge []byte, err error)

	// Verify signed challenge and update connection, if successful.
	FinalizeHandshake(attemptId string, response []byte, conn *websocket.Conn) error

	// Clear attempt's state.
	AbortHandshake(attemptId string)
}

// Components going into the auth header, excluding the signature.
type AuthHeaderElems struct {
	Timestamp  uint32
	DonId      string
	GatewayURL string
}

func Pack(elems *AuthHeaderElems) []byte {
	packed := common.Uint32ToBytes(elems.Timestamp)
	packed = append(packed, common.StringToAlignedBytes(elems.DonId, HandshakeDonIdLen)...)
	packed = append(packed, []byte(elems.GatewayURL)...)
	return packed
}

func Unpack(data []byte) (*AuthHeaderElems, error) {
	unpacked := &AuthHeaderElems{}
	unpacked.Timestamp = common.BytesToUint32(data[0:HandshakeTimestampLen])
	unpacked.DonId = common.AlignedBytesToString(data[HandshakeTimestampLen : HandshakeTimestampLen+HandshakeDonIdLen])
	unpacked.GatewayURL = string(data[HandshakeTimestampLen+HandshakeDonIdLen:])
	return unpacked, nil
}
