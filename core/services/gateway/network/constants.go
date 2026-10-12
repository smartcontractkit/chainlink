package network

const (
	WsServerHandshakeAuthHeaderName      string = "Authorization"
	WsServerHandshakeChallengeHeaderName string = "Challenge"

	HandshakeTimestampLen            int = 4
	HandshakeDonIDLen                int = 64
	HandshakeGatewayURLLen           int = 128
	HandshakeSignatureLen            int = 65
	HandshakeAuthHeaderLen           int = HandshakeTimestampLen + HandshakeDonIDLen + HandshakeGatewayURLLen + HandshakeSignatureLen
	HandshakeEncodedAuthHeaderMaxLen int = 512
	HandshakeChallengeMinLen         int = HandshakeTimestampLen + HandshakeGatewayURLLen + 1
)
