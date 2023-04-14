package gateway

import "encoding/json"

/*
 * Top-level Message structure with
 *   - universal fields identifying the sender and the target DON
 *   - product-specific RequestPayload
 */
type Message struct {
	DonId         string                 `json:"don_id"`
	ServiceId     string                 `json:"service_id"`
	MessageId     string                 `json:"message_id"`
	Method        string                 `json:"method"`
	SenderAddress string                 `json:"sender_address"`
	Signature     string                 `json:"signature"`
	Payload       map[string]interface{} `json:"payload"`
}

func Decode(msgBytes []byte) (msg *Message, err error) {
	err = json.Unmarshal(msgBytes, &msg)
	return msg, err
}

func Encode(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}
