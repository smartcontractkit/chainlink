package environment

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func StartJD(lggr zerolog.Logger, jdInput jd.Input, infraInput infra.Input) (*jd.Output, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")

	var jdOutput *jd.Output
	if infraInput.Type == infra.CRIB {
		deployCribJdInput := &cre.DeployCribJdInput{
			JDInput:        jdInput,
			CribConfigsDir: cribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		var jdErr error
		jdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
		}
	}

	var jdErr error
	jdOutput, jdErr = CreateJobDistributor(jdInput)
	if jdErr != nil {
		jdErr = fmt.Errorf("failed to start JD container for image %s: %w", jdInput.Image, jdErr)

		// useful end user messages
		if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
			jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
		}
		return nil, jdErr
	}

	lggr.Info().Msgf("Job Distributor started in %.2f seconds", time.Since(startTime).Seconds())

	return jdOutput, nil
}

func CreateJobDistributor(input jd.Input) (*jd.Output, error) {
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(E2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(E2eJobDistributorVersionEnvVarName)
		input.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}

	jdOutput, err := jd.NewJD(&input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create new job distributor")
	}

	return jdOutput, nil
}

func StartDONsAndJD(lggr zerolog.Logger, jdInput *jd.Input, registryChainBlockchainOutput *blockchain.Output, topology *cre.Topology, infraInput infra.Input, capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet) (*jd.Output, []*cre.WrappedNodeOutput, error) {
	if jdInput == nil {
		return nil, nil, errors.New("jd input is nil")
	}
	if registryChainBlockchainOutput == nil {
		return nil, nil, errors.New("registry chain blockchain output is nil")
	}
	if topology == nil {
		return nil, nil, errors.New("topology is nil")
	}
	var jdOutput *jd.Output
	jdAndDonsErrGroup := &errgroup.Group{}

	jdAndDonsErrGroup.Go(func() error {
		var startJDErr error
		jdOutput, startJDErr = StartJD(lggr, *jdInput, infraInput)
		if startJDErr != nil {
			return pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}

		return nil
	})

	nodeSetOutput := make([]*cre.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))

	jdAndDonsErrGroup.Go(func() error {
		var startDonsErr error
		nodeSetOutput, startDonsErr = StartDONs(lggr, topology, infraInput, registryChainBlockchainOutput, capabilitiesAwareNodeSets)
		if startDonsErr != nil {
			return pkgerrors.Wrap(startDonsErr, "failed to start DONs")
		}

		return nil
	})

	if jdAndDonErr := jdAndDonsErrGroup.Wait(); jdAndDonErr != nil {
		return nil, nil, pkgerrors.Wrap(jdAndDonErr, "failed to start Job Distributor or DONs")
	}

	return jdOutput, nodeSetOutput, nil
}
