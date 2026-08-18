package capabilityrunner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// HTTPPortArg is the required command-line argument naming the port of the
// runner binary's HTTP server, which serves the limits reload endpoint.
const HTTPPortArg = "--http.port"

func ValidatedCapabilityRunnerSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		ExternalJobID: uuid.New(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}

	var spec job.CapabilityRunnerSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	jb.CapabilityRunnerSpec = &spec
	if jb.Type != job.CapabilityRunner {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	if jb.CapabilityRunnerSpec.Command == "" {
		return jb, errors.New("command must be set")
	}
	if _, err := HTTPPortFromArgs(jb.CapabilityRunnerSpec.Args); err != nil {
		return jb, err
	}

	return jb, nil
}

// HTTPPortFromArgs extracts the required --http.port argument, accepting both
// "--http.port=8080" and "--http.port 8080" forms.
func HTTPPortFromArgs(args []string) (int, error) {
	for i, arg := range args {
		var v string
		switch {
		case strings.HasPrefix(arg, HTTPPortArg+"="):
			v = strings.TrimPrefix(arg, HTTPPortArg+"=")
		case arg == HTTPPortArg && i+1 < len(args):
			v = args[i+1]
		default:
			continue
		}
		port, err := strconv.Atoi(v)
		if err != nil || port <= 0 || port > 65535 {
			return 0, fmt.Errorf("invalid %s value %q", HTTPPortArg, v)
		}
		return port, nil
	}
	return 0, fmt.Errorf("args must include %s", HTTPPortArg)
}
