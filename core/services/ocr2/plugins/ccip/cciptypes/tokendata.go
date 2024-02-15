package cciptypes

import "context"

type TokenDataReader interface {
	// ReadTokenData returns the attestation bytes if ready, and throws an error if not ready.
	// It supports messages with a single token transfer, the returned []byte has the token data for the first token of the msg.
	ReadTokenData(ctx context.Context, msg EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) (tokenData []byte, err error)
}
