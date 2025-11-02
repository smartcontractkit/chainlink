package vaultutils

import (
	"encoding/json"
	"fmt"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
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

	JSONBytes, err := jsonv2.Marshal(jsond, jsonv2.Deterministic(true))
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	canonicalJSONBytes := jsontext.Value(JSONBytes)
	err = canonicalJSONBytes.Canonicalize()
	if err != nil {
		return nil, fmt.Errorf("error canonicalizing JSON: %w", err)
	}

	return canonicalJSONBytes, nil
}
