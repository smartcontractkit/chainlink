package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// CapabilityConfig is an untyped map representation of the CapabilityConfig proto message
// It provides methods to marshal/unmarshal to/from proto bytes
type CapabilityConfig map[string]any

// UnmarshalWithValidation unmarshals JSON into a proto message, but first validates
// that all fields in the JSON match the proto schema. If there are unknown fields,
// it returns an error listing them
func UnmarshalWithValidation(jsonData []byte, msg proto.Message) error {
	strictOps := protojson.UnmarshalOptions{DiscardUnknown: false}
	if err := strictOps.Unmarshal(jsonData, msg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

// MarshalProto marshals the CapabilityConfig to proto bytes
// If the CapabilityConfig is nil, it returns nil, nil, to support empty configs
func (c CapabilityConfig) MarshalProto() ([]byte, error) {
	lggr := logger.Nop()

	if c == nil {
		return nil, nil
	}
	jsonEncodedCfg, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal config: %w", err)
	}

	pbCfg := &pb.CapabilityConfig{}
	if errValidation := UnmarshalWithValidation(jsonEncodedCfg, pbCfg); errValidation != nil {
		// log the error with the specific field names that don't match
		lggr.Warnf("⚠️  WARNING: Config validation failed: %v", errValidation)
		// Also print to stderr so it's visible no matter what
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: Config validation failed: %v\n", errValidation)

		return nil, fmt.Errorf("failed to protojson unmarshal json encoded config: %w", errValidation)
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
