package network

const (
	WsServerHandshakeAuthHeaderName      string = "Authorization"
	WsServerHandshakeChallengeHeaderName string = "Challenge"

	HandshakeTimestampLen            int = 4
	HandshakeDonIdLen                int = 64
	HandshakeGatewayURLMinLen        int = 10
	HandshakeSignatureLen            int = 65
	HandshakeAuthHeaderMinLen        int = HandshakeTimestampLen + HandshakeDonIdLen + HandshakeGatewayURLMinLen + HandshakeSignatureLen
	HandshakeEncodedAuthHeaderMaxLen int = 512
)
