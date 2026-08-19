package network

import (
	"context"
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
	NewAuthHeader(ctx context.Context, url *url.URL) ([]byte, error)

	// Sign challenge to prove identity.
	ChallengeResponse(ctx context.Context, url *url.URL, challenge []byte) ([]byte, error)
}

type ConnectionAcceptor interface {
	// Verify auth header, save state of the attempt and generate a challenge for the node.
	StartHandshake(authHeader []byte) (attemptID string, challenge []byte, err error)

	// Verify signed challenge and update connection, if successful.
	FinalizeHandshake(attemptID string, response []byte, conn *websocket.Conn) error

	// Clear attempt's state.
	AbortHandshake(attemptID string)
}

// Components going into the auth header, excluding the signature.
type AuthHeaderElems struct {
	Timestamp uint32
	DonID     string
	GatewayID string
}

type ChallengeElems struct {
	Timestamp      uint32
	GatewayID      string
	ChallengeBytes []byte
}

var (
	ErrAuthHeaderParse           = errors.New("unable to parse auth header")
	ErrAuthInvalidDonID          = errors.New("invalid DON ID")
	ErrAuthInvalidNode           = errors.New("unexpected node address")
	ErrAuthInvalidGateway        = errors.New("invalid gateway ID")
	ErrAuthInvalidTimestamp      = errors.New("timestamp outside of tolerance range")
	ErrChallengeTooShort         = errors.New("challenge too short")
	ErrChallengeAttemptNotFound  = errors.New("attempt not found")
	ErrChallengeInvalidSignature = errors.New("invalid challenge signature")
)

func PackAuthHeader(elems *AuthHeaderElems) []byte {
	packed := common.Uint32ToBytes(elems.Timestamp)
	packed = append(packed, common.StringToAlignedBytes(elems.DonID, HandshakeDonIDLen)...)
	packed = append(packed, common.StringToAlignedBytes(elems.GatewayID, HandshakeGatewayURLLen)...)
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
	elems.DonID = common.AlignedBytesToString(data[offset : offset+HandshakeDonIDLen])
	offset += HandshakeDonIDLen
	elems.GatewayID = common.AlignedBytesToString(data[offset : offset+HandshakeGatewayURLLen])
	offset += HandshakeGatewayURLLen
	signature := data[offset:]
	signer, err = common.ExtractSigner(signature, data[:len(data)-HandshakeSignatureLen])
	return
}

func PackChallenge(elems *ChallengeElems) []byte {
	packed := common.Uint32ToBytes(elems.Timestamp)
	packed = append(packed, common.StringToAlignedBytes(elems.GatewayID, HandshakeGatewayURLLen)...)
	packed = append(packed, elems.ChallengeBytes...)
	return packed
}

func UnpackChallenge(data []byte) (*ChallengeElems, error) {
	if len(data) < HandshakeChallengeMinLen {
		return nil, fmt.Errorf("challenge length is too small (expected at least: %d, got: %d)", HandshakeChallengeMinLen, len(data))
	}
	unpacked := &ChallengeElems{}
	unpacked.Timestamp = common.BytesToUint32(data[0:HandshakeTimestampLen])
	unpacked.GatewayID = common.AlignedBytesToString(data[HandshakeTimestampLen : HandshakeTimestampLen+HandshakeGatewayURLLen])
	unpacked.ChallengeBytes = data[HandshakeTimestampLen+HandshakeGatewayURLLen:]
	return unpacked, nil
}
