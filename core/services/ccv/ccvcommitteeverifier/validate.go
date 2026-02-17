package ccvcommitteeverifier

import (
	"errors"
	"fmt"

	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-ccv/verifier/commit"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func ValidatedCCVCommitteeVerifierSpec(tomlString string) (jb job.Job, err error) {
	var spec job.CCVCommitteeVerifierSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&spec)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return job.Job{}, err
	}
	jb.CCVCommitteeVerifierSpec = &spec

	if jb.Type != job.CCVCommitteeVerifier {
		return job.Job{}, fmt.Errorf("the only supported type is currently 'ccvcommitteeverifier', got %s", jb.Type)
	}
	if jb.CCVCommitteeVerifierSpec.CommitteeVerifierConfig == "" {
		return job.Job{}, errors.New("committeeVerifierConfig must be set")
	}

	var cfg commit.Config
	err = toml.Unmarshal([]byte(jb.CCVCommitteeVerifierSpec.CommitteeVerifierConfig), &cfg)
	if err != nil {
		return job.Job{}, fmt.Errorf("failed to unmarshal committeeVerifierConfig into the verifier config struct: %w", err)
	}

	// validation functions are called in ServicesForSpec
	// so that this method can be used in tests w/out crafting
	// a valid configuration.

	return jb, nil
}
