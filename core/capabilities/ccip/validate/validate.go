package validate

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

// ValidatedCCIPSpec validates the given toml string as a CCIP spec.
func ValidatedCCIPSpec(tomlString string) (jb job.Job, err error) {
	var spec job.CCIPSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml error on load: %w", err)
	}
	// Note this validates all the fields which implement an UnmarshalText
	err = tree.Unmarshal(&spec)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml unmarshal error on job: %w", err)
	}
	jb.CCIPSpec = &spec

	if jb.Type != job.CCIP {
		return job.Job{}, fmt.Errorf("the only supported type is currently 'ccip', got %s", jb.Type)
	}
	if jb.CCIPSpec.CapabilityLabelledName == "" {
		return job.Job{}, fmt.Errorf("capabilityLabelledName must be set")
	}
	if jb.CCIPSpec.CapabilityVersion == "" {
		return job.Job{}, fmt.Errorf("capabilityVersion must be set")
	}
	if jb.CCIPSpec.P2PKeyID == "" {
		return job.Job{}, fmt.Errorf("p2pKeyID must be set")
	}

	// ensure that the P2PV2Bootstrappers is in the right format.
	for _, bootstrapperLocator := range jb.CCIPSpec.P2PV2Bootstrappers {
		// needs to be of the form <peer_id>@<ip-address>:<port>
		_, err := ocrcommon.ParseBootstrapPeers([]string{bootstrapperLocator})
		if err != nil {
			return job.Job{}, fmt.Errorf("p2p v2 bootstrapper locator %s is not in the correct format: %w", bootstrapperLocator, err)
		}
	}

	return jb, nil
}

type SpecArgs struct {
	P2PV2Bootstrappers     []string          `toml:"p2pV2Bootstrappers"`
	CapabilityVersion      string            `toml:"capabilityVersion"`
	CapabilityLabelledName string            `toml:"capabilityLabelledName"`
	OCRKeyBundleIDs        map[string]string `toml:"ocrKeyBundleIDs"`
	P2PKeyID               string            `toml:"p2pKeyID"`
	RelayConfigs           map[string]any    `toml:"relayConfigs"`
	PluginConfig           map[string]any    `toml:"pluginConfig"`
}

// NewCCIPSpecToml creates a new CCIP spec in toml format from the given spec args.
func NewCCIPSpecToml(spec SpecArgs) (string, error) {
	type fullSpec struct {
		SpecArgs
		Type          string `toml:"type"`
		SchemaVersion uint64 `toml:"schemaVersion"`
		Name          string `toml:"name"`
		ExternalJobID string `toml:"externalJobID"`
	}
	extJobID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate external job id: %w", err)
	}
	marshaled, err := toml.Marshal(fullSpec{
		SpecArgs:      spec,
		Type:          "ccip",
		SchemaVersion: 1,
		Name:          fmt.Sprintf("%s-%s", "ccip", extJobID.String()),
		ExternalJobID: extJobID.String(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal spec into toml: %w", err)
	}

	return string(marshaled), nil
}
