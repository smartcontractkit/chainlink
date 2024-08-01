package wsrpc

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// MarshalProtoMessage returns the protobuf message wire format of v.
func MarshalProtoMessage(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}

	return proto.Marshal(vv)
}

// Unmarshal parses the protobuf wire format into v.
func UnmarshalProtoMessage(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}

	return proto.Unmarshal(data, vv)
}
