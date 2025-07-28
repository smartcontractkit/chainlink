package vault

import (
	"encoding/json"

	"github.com/gibson042/canonicaljson-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ToCanonicalJSON converts a protobuf message to a stable, deterministic
// representation, including consistent sorting of keys and fields, and
// consistent spacing.
func ToCanonicalJSON(msg proto.Message) ([]byte, error) {
	jsonb, err := protojson.MarshalOptions{
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		return nil, err
	}

	jsond := map[string]any{}
	err = json.Unmarshal(jsonb, &jsond)
	if err != nil {
		return nil, err
	}

	return canonicaljson.Marshal(jsond)
}
