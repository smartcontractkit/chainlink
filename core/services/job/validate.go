package job

import (
	"bytes"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

var (
	ErrNoPipelineSpec       = errors.New("pipeline spec not specified")
	ErrInvalidJobType       = errors.New("invalid job type")
	ErrNoSchemaVersion      = errors.New("schema version not specified")
	ErrInvalidSchemaVersion = errors.New("schema version invalid")
	schemaVersionsSupported = map[Type]map[uint32]struct{}{
		Cron:              {1: {}},
		DirectRequest:     {1: {}},
		FluxMonitor:       {1: {}},
		OffchainReporting: {1: {}},
		Keeper:            {1: {}},
		VRF:               {1: {}},
		Webhook:           {1: {}},
	}
)

// Common spec validation
func ValidateSpec(ts string) (Job, error) {
	var jb Job
	err := toml.NewDecoder(bytes.NewReader([]byte(ts))).Strict(true).Decode(&jb)
	if err != nil {
		return jb, err
	}
	versions, ok := schemaVersionsSupported[jb.Type]
	if !ok {
		return jb, ErrInvalidJobType
	}
	if jb.SchemaVersion == 0 {
		return jb, ErrNoSchemaVersion
	}
	if _, ok := versions[jb.SchemaVersion]; !ok {
		return jb, ErrInvalidSchemaVersion
	}

	if jb.Type.HasPipelineSpec() && (jb.Pipeline.Source == "") {
		return jb, ErrNoPipelineSpec
	}
	return jb, nil
}
