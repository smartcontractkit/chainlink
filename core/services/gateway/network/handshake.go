package network

import (
	"net/url"

	"github.com/gorilla/websocket"
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
	NewAuthHeader(url *url.URL) []byte
	ChallengeResponse(challenge []byte) ([]byte, error)
}

type ConnectionAcceptor interface {
	StartHandshake(authHeader []byte) (attemptId string, challenge []byte, err error)
	FinalizeHandshake(attemptId string, response []byte, conn *websocket.Conn) error
	AbortHandshake(attemptId string)
}
