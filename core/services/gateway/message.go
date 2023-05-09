package gateway

import "encoding/json"

/*
 * Top-level Message structure containing:
 *   - universal fields identifying the request, the sender and the target DON/service
 *   - product-specific payload
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
