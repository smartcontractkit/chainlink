package network

import (
	"errors"
	"fmt"
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
type ConnectionInitiator interface {
	// Generate authentication header value specific to node and gateway
	NewAuthHeader(url *url.URL) ([]byte, error)

	// Sign challenge to prove identity.
	ChallengeResponse(url *url.URL, challenge []byte) ([]byte, error)
}

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
	Timestamp uint32
	DonId     string
	GatewayId string
}

type ChallengeElems struct {
	Timestamp      uint32
	GatewayId      string
	ChallengeBytes []byte
}

var (
	ErrAuthHeaderParse           = errors.New("unable to parse auth header")
	ErrAuthInvalidDonId          = errors.New("invalid DON ID")
	ErrAuthInvalidNode           = errors.New("unexpected node address")
	ErrAuthInvalidGateway        = errors.New("invalid gateway ID")
	ErrAuthInvalidTimestamp      = errors.New("timestamp outside of tolerance range")
	ErrChallengeTooShort         = errors.New("challenge too short")
	ErrChallengeAttemptNotFound  = errors.New("attempt not found")
	ErrChallengeInvalidSignature = errors.New("invalid challenge signature")
)

func PackAuthHeader(elems *AuthHeaderElems) []byte {
	packed := common.Uint32ToBytes(elems.Timestamp)
	packed = append(packed, common.StringToAlignedBytes(elems.DonId, HandshakeDonIdLen)...)
	packed = append(packed, common.StringToAlignedBytes(elems.GatewayId, HandshakeGatewayURLLen)...)
	return packed
}

func UnpackSignedAuthHeader(data []byte) (elems *AuthHeaderElems, signer []byte, err error) {
	if len(data) != HandshakeAuthHeaderLen {
		return nil, nil, fmt.Errorf("auth header length is invalid (expected: %d, got: %d)", HandshakeAuthHeaderLen, len(data))
	}
	elems = &AuthHeaderElems{}
	offset := 0
	elems.Timestamp = common.BytesToUint32(data[offset : offset+HandshakeTimestampLen])
	offset += HandshakeTimestampLen
	elems.DonId = common.AlignedBytesToString(data[offset : offset+HandshakeDonIdLen])
	offset += HandshakeDonIdLen
	elems.GatewayId = common.AlignedBytesToString(data[offset : offset+HandshakeGatewayURLLen])
	offset += HandshakeGatewayURLLen
	signature := data[offset:]
	signer, err = common.ExtractSigner(signature, data[:len(data)-HandshakeSignatureLen])
	return
}

func PackChallenge(elems *ChallengeElems) []byte {
	packed := common.Uint32ToBytes(elems.Timestamp)
	packed = append(packed, common.StringToAlignedBytes(elems.GatewayId, HandshakeGatewayURLLen)...)
	packed = append(packed, elems.ChallengeBytes...)
	return packed
}

func UnpackChallenge(data []byte) (*ChallengeElems, error) {
	if len(data) < HandshakeChallengeMinLen {
		return nil, fmt.Errorf("challenge length is too small (expected at least: %d, got: %d)", HandshakeChallengeMinLen, len(data))
	}
	unpacked := &ChallengeElems{}
	unpacked.Timestamp = common.BytesToUint32(data[0:HandshakeTimestampLen])
	unpacked.GatewayId = common.AlignedBytesToString(data[HandshakeTimestampLen : HandshakeTimestampLen+HandshakeGatewayURLLen])
	unpacked.ChallengeBytes = data[HandshakeTimestampLen+HandshakeGatewayURLLen:]
	return unpacked, nil
}
