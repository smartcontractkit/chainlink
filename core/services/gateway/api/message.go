package api

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"

	gw_common "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	Payload json.RawMessage `json:"payload,omitempty"`
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
	signerBytes, err := m.ValidateSignature()
	if err != nil {
		return err
	}
	hexSigner := utils.StringToHex(string(signerBytes))
	if m.Body.Sender != "" && m.Body.Sender != hexSigner {
		return errors.New("sender doesn't match signer")
	}
	m.Body.Sender = hexSigner
	return nil
}

// Message signatures are over the following data:
//  1. MessageId aligned to 128 bytes
//  2. Method aligned to 64 bytes
//  3. DonId aligned to 64 bytes
//  4. Payload (before parsing)
func (m *Message) Sign(privateKey *ecdsa.PrivateKey) error {
	rawData, err := getRawMessageBody(&m.Body)
	if err != nil {
		return err
	}
	signature, err := gw_common.SignData(privateKey, rawData...)
	if err != nil {
		return err
	}
	m.Signature = utils.StringToHex(string(signature))
	return nil
}

func (m *Message) ValidateSignature() (signerAddress []byte, err error) {
	rawData, err := getRawMessageBody(&m.Body)
	if err != nil {
		return
	}
	signatureBytes, err := utils.TryParseHex(m.Signature)
	if err != nil {
		return
	}
	return gw_common.ValidateSignature(signatureBytes, rawData...)
}

func getRawMessageBody(msgBody *MessageBody) ([][]byte, error) {
	if msgBody == nil {
		return nil, errors.New("nil message")
	}
	alignedMessageId := make([]byte, MessageIdMaxLen)
	copy(alignedMessageId, msgBody.MessageId)
	alignedMethod := make([]byte, MessageMethodMaxLen)
	copy(alignedMethod, msgBody.Method)
	alignedDonId := make([]byte, MessageDonIdMaxLen)
	copy(alignedDonId, msgBody.DonId)
	return [][]byte{alignedMessageId, alignedMethod, alignedDonId, msgBody.Payload}, nil
}
