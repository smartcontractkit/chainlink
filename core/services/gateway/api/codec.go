package api

// Codec implements (de)serialization of Message objects.
type Codec interface {
	DecodeRequest(msgBytes []byte) (*Message, error)

	EncodeRequest(msg *Message) ([]byte, error)

	DecodeResponse(msgBytes []byte) (*Message, error)

	EncodeResponse(msg *Message) ([]byte, error)

	EncodeNewErrorResponse(id string, code int, message string, data []byte) ([]byte, error)
}
