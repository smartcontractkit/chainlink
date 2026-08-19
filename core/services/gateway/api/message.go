package api

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"

	gw_common "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	MessageSignatureLen           = 65
	MessageSignatureHexEncodedLen = 2 + 2*MessageSignatureLen
	MessageIDMaxLen               = 128
	MessageMethodMaxLen           = 64
	MessageDonIDMaxLen            = 64
	MessageReceiverLen            = 2 + 2*20
	NullChar                      = "\x00"
	MethodInternalError           = "internal_error"
)

/*
 * Top-level Message structure containing:
 *   - universal fields identifying the request, the sender and the target DON/service
 *   - product-specific payload
 *
 * Signature, Receiver and Sender are hex-encoded with a "0x" prefix.
 */
type Message struct {
	Signature string      `json:"signature"`
	Body      MessageBody `json:"body"`
}

type MessageBody struct {
	MessageID string `json:"message_id"`
	Method    string `json:"method"`
	DonID     string `json:"don_id"`
	Receiver  string `json:"receiver"`
	// Service-specific payload, decoded inside the Handler.
	Payload json.RawMessage `json:"payload,omitempty"`

	// Fields only used locally for convenience. Not serialized.
	Sender string `json:"-"`
}

func (m *Message) Validate() error {
	if m == nil {
		return errors.New("nil message")
	}
	if len(m.Signature) != MessageSignatureHexEncodedLen {
		return errors.New("invalid hex-encoded signature length")
	}
	if len(m.Body.MessageID) == 0 || len(m.Body.MessageID) > MessageIDMaxLen {
		return errors.New("invalid message ID length")
	}
	if strings.HasSuffix(m.Body.MessageID, NullChar) {
		return errors.New("message ID ending with null bytes")
	}
	if len(m.Body.Method) == 0 || len(m.Body.Method) > MessageMethodMaxLen {
		return errors.New("invalid method name length")
	}
	if strings.HasSuffix(m.Body.Method, NullChar) {
		return errors.New("method name ending with null bytes")
	}
	if len(m.Body.DonID) == 0 || len(m.Body.DonID) > MessageDonIDMaxLen {
		return errors.New("invalid DON ID length")
	}
	if strings.HasSuffix(m.Body.DonID, NullChar) {
		return errors.New("DON ID ending with null bytes")
	}
	if len(m.Body.Receiver) != 0 && len(m.Body.Receiver) != MessageReceiverLen {
		return errors.New("invalid Receiver length")
	}
	signerBytes, err := m.ExtractSigner()
	if err != nil {
		return err
	}
	m.Body.Sender = utils.StringToHex(string(signerBytes))
	return nil
}

// Message signatures are over the following data:
//  1. MessageID aligned to 128 bytes
//  2. Method aligned to 64 bytes
//  3. DonID aligned to 64 bytes
//  4. Receiver (in hex) aligned to 42 bytes
//  5. Payload (raw bytes before parsing)
func (m *Message) Sign(privateKey *ecdsa.PrivateKey) error {
	if m == nil {
		return errors.New("nil message")
	}
	rawData := GetRawMessageBody(&m.Body)
	signature, err := gw_common.SignData(privateKey, rawData...)
	if err != nil {
		return err
	}
	m.Signature = utils.StringToHex(string(signature))
	m.Body.Sender = strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	return nil
}

func (m *Message) SignKS(ctx context.Context, ks keys.MessageSigner, signer common.Address) error {
	if m == nil {
		return errors.New("nil message")
	}
	rawData := GetRawMessageBody(&m.Body)
	signature, err := ks.SignMessage(ctx, signer, gw_common.Flatten(rawData...))
	if err != nil {
		return err
	}
	m.Signature = utils.StringToHex(string(signature))
	m.Body.Sender = strings.ToLower(signer.Hex())
	return nil
}

func (m *Message) ExtractSigner() (signerAddress []byte, err error) {
	if m == nil {
		return nil, errors.New("nil message")
	}
	rawData := GetRawMessageBody(&m.Body)
	signatureBytes, err := hex.DecodeString(m.Signature)
	if err != nil {
		return nil, err
	}
	return gw_common.ExtractSigner(signatureBytes, rawData...)
}

func GetRawMessageBody(msgBody *MessageBody) [][]byte {
	alignedMessageID := make([]byte, MessageIDMaxLen)
	copy(alignedMessageID, msgBody.MessageID)
	alignedMethod := make([]byte, MessageMethodMaxLen)
	copy(alignedMethod, msgBody.Method)
	alignedDonID := make([]byte, MessageDonIDMaxLen)
	copy(alignedDonID, msgBody.DonID)
	alignedReceiver := make([]byte, MessageReceiverLen)
	copy(alignedReceiver, msgBody.Receiver)
	return [][]byte{alignedMessageID, alignedMethod, alignedDonID, alignedReceiver, msgBody.Payload}
}
