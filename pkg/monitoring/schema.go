package monitoring

import (
	"encoding/binary"
	"fmt"

	"github.com/riferrei/srclient"
)

type Schema interface {
	ID() int
	Version() int
	Subject() string
	Encode(interface{}) ([]byte, error)
	Decode([]byte) (interface{}, error)
}

type wrapSchema struct {
	subject string
	*srclient.Schema
}

func (w wrapSchema) ID() int {
	return w.Schema.ID()
}
func (w wrapSchema) Version() int {
	return w.Schema.Version()
}

func (w wrapSchema) Subject() string {
	return w.subject
}

func (w wrapSchema) Encode(value interface{}) ([]byte, error) {
	payload, err := w.Schema.Codec().BinaryFromNative(nil, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode value in avro: %w", err)
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(w.Schema.ID()))

	// Magic 0 byte + 4 bytes of schema ID + the data bytes
	bytes := []byte{0}
	bytes = append(bytes, schemaIDBytes...)
	bytes = append(bytes, payload...)
	return bytes, nil
}

func (w wrapSchema) Decode(buf []byte) (interface{}, error) {
	if buf[0] != 0 {
		return nil, fmt.Errorf("magic byte not 0, instead is %d", buf[0])
	}
	schemaID := int(binary.BigEndian.Uint32(buf[1:5]))
	if schemaID != w.ID() {
		return nil, fmt.Errorf("decoding message for a different schema, found schema id is %d but expected %d", schemaID, w.ID())
	}
	value, _, err := w.Schema.Codec().NativeFromBinary(buf[5:])
	return value, err
}

func (w wrapSchema) String() string {
	return fmt.Sprintf("schema(subject=%s,id=%d,version=%d)", w.subject, w.Schema.ID(), w.Schema.Version())
}

// SubjectFromTopic is a utility to the associated schema subject from a kafka topic name.
func SubjectFromTopic(topic string) string {
	return fmt.Sprintf("%s-value", topic)
}
