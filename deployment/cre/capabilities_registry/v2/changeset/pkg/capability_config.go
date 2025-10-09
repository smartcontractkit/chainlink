package pkg

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// CapabilityConfig is an untyped map representation of the CapabilityConfig proto message
// It provides methods to marshal/unmarshal to/from proto bytes
type CapabilityConfig map[string]any

// MarshalProto marshals the CapabilityConfig to proto bytes
// If the CapabilityConfig is nil, it returns nil, nil, to support empty configs
func (c CapabilityConfig) MarshalProto() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	jsonEncodedCfg, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal config: %w", err)
	}

	fmt.Println("JSON Encoded Config:", string(jsonEncodedCfg))

	pbCfg := &pb.CapabilityConfig{}
	ops := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = ops.Unmarshal(jsonEncodedCfg, pbCfg); err != nil {
		return nil, fmt.Errorf("failed to protojson unmarshal json encoded config %w", err)
	}

	protoEncodedCfg, err := proto.Marshal(pbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to proto marshal %T: %w", pbCfg, err)
	}

	return protoEncodedCfg, nil
}

// UnmarshalProto unmarshals proto bytes into the CapabilityConfig
func (c *CapabilityConfig) UnmarshalProto(data []byte) error {
	pbCfg := &pb.CapabilityConfig{}
	if err := proto.Unmarshal(data, pbCfg); err != nil {
		return fmt.Errorf("failed to proto unmarshal data into %T: %w", pbCfg, err)
	}

	jsonEncodedCfg, err := protojson.Marshal(pbCfg)
	if err != nil {
		return fmt.Errorf("failed to protojson marshal %T: %w", pbCfg, err)
	}
	if err := json.Unmarshal(jsonEncodedCfg, &c); err != nil {
		return fmt.Errorf("failed to json unmarshal into CapabilityConfig: %w", err)
	}

	return nil
}

func (c CapabilityConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(c)) // avoid infinite recursion by casting to underlying type
}

func (c *CapabilityConfig) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("cannot unmarshal into nil CapabilityConfig")
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to json unmarshal into map: %w", err)
	}
	*c = m
	return nil
}
