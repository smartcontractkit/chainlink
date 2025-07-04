package api

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

// Codec implements (de)serialization of Message objects.
type Codec interface {
	DecodeRawRequest(msgBytes []byte, jwtToken string) (*Message, error)

	DecodeJSONRequest(request jsonrpc2.Request[json.RawMessage]) (*Message, error)

	// EncodeLegacyRequest creates a Json request with a Message object
	// embedded in jsonrpc2.Request.Params as opposed to new requests,
	// which add payload fields directly in jsonrpc2.Request.Params.
	EncodeLegacyRequest(msg *Message) ([]byte, error)

	DecodeLegacyResponse(msgBytes []byte) (*Message, error)

	EncodeLegacyResponse(msg *Message) []byte

	EncodeNewErrorResponse(id string, code int64, message string, data []byte) []byte
}
