package web2

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

func ValidatedSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		ExternalJobID: uuid.NewV4(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}

	if jb.Type != job.VRFWeb2 {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	var spec job.VRFWeb2Spec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	if spec.LotteryConsumerAddress.String() == "" {
		return jb, errors.New("lotteryConsumerAddress must be set")
	}

	if len(spec.FromAddresses) == 0 {
		return jb, errors.New("fromAddresses must be set")
	}

	jb.VRFWeb2Spec = &spec
	return jb, nil
}
