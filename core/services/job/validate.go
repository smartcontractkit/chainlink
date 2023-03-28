package job

import (
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

var (
	ErrNoPipelineSpec       = errors.New("pipeline spec not specified")
	ErrInvalidJobType       = errors.New("invalid job type")
	ErrInvalidSchemaVersion = errors.New("invalid schema version")
	jobTypes                = map[Type]struct{}{
		Cron:               {},
		DirectRequest:      {},
		FluxMonitor:        {},
		OffchainReporting:  {},
		OffchainReporting2: {},
		Keeper:             {},
		VRF:                {},
		Webhook:            {},
		BlockhashStore:     {},
		Bootstrap:          {},
	}
)

// ValidateSpec is the common spec validation
func ValidateSpec(ts string) (Type, error) {
	var jb Job
	// Note we can't use:
	//   toml.NewDecoder(bytes.NewReader([]byte(ts))).Strict(true).Decode(&jb)
	// to error in the case of unrecognized keys because all the keys in the toml are at
	// the top level and so decoding for the job will have undecodable keys meant for the job
	// type specific struct and vice versa. Should we upgrade the schema,
	// we put the type specific config in its own subtree e.g.
	// 	schemaVersion=1
	//  name="test"
	//  [vrf_spec]
	//  publicKey="0x..."
	// and then we could use it.
	tree, err := toml.Load(ts)
	if err != nil {
		return "", err
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return "", err
	}
	if _, ok := jobTypes[jb.Type]; !ok {
		return "", ErrInvalidJobType
	}
	if jb.Type.SchemaVersion() != jb.SchemaVersion {
		return "", ErrInvalidSchemaVersion
	}
	if jb.Type.RequiresPipelineSpec() && (jb.Pipeline.Source == "") {
		return "", ErrNoPipelineSpec
	}
	if jb.Pipeline.RequiresPreInsert() && !jb.Type.SupportsAsync() {
		return "", errors.Errorf("async=true tasks are not supported for %v", jb.Type)
	}

	if strings.Contains(ts, "<{}>") {
		return "", errors.Errorf("'<{}>' syntax is not supported. Please use \"{}\" instead")
	}

	return jb.Type, nil
}
