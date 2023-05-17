package gateway

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

const (
	MessageSignatureLen           = 65
	MessageSignatureHexEncodedLen = 2 + 2*MessageSignatureLen
	MessageIdMaxLen               = 128
	MessageMethodMaxLen           = 64
	MessageDonIdMaxLen            = 64
	MessageSenderLen              = 20
	MessageSenderHexEncodedLen    = 2 + 2*MessageSenderLen
)

/*
 * Top-level Message structure containing:
 *   - universal fields identifying the request, the sender and the target DON/service
 *   - product-specific payload
 *
 * Signature and Sender are hex-encoded with a "0x" prefix.
 */
type Message struct {
	Signature string      `json:"signature"`
	Body      MessageBody `json:"body"`
}

type MessageBody struct {
	MessageId string `json:"message_id"`
	Method    string `json:"method"`
	DonId     string `json:"don_id"`
	Sender    string `json:"sender"`

	// Service-specific payload, decoded inside the Handler.
	Payload json.RawMessage `json:"payload"`
}

func (m *Message) Validate() error {
	if m == nil {
		return errors.New("nil message")
	}
	if len(m.Signature) != MessageSignatureHexEncodedLen {
		return errors.New("invalid hex-encoded signature length")
	}
	if len(m.Body.MessageId) == 0 || len(m.Body.MessageId) > MessageIdMaxLen {
		return errors.New("invalid message ID length")
	}
	if len(m.Body.Method) == 0 || len(m.Body.Method) > MessageMethodMaxLen {
		return errors.New("invalid method name length")
	}
	if len(m.Body.DonId) == 0 || len(m.Body.DonId) > MessageDonIdMaxLen {
		return errors.New("invalid DON ID length")
	}
	if len(m.Body.Sender) != MessageSenderHexEncodedLen || !common.IsHexAddress(m.Body.Sender) {
		return errors.New("invalid hex-encoded sender address")
	}
	return nil
}
