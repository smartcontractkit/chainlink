package environment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
)

func StartJD(lggr zerolog.Logger, nixShell *nix.Shell, jdInput jd.Input, infraInput infra.Input) (*jd.Output, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")

	var jdOutput *jd.Output
	if infraInput.Type == infra.CRIB {
		deployCribJdInput := &cre.DeployCribJdInput{
			JDInput:        &jdInput,
			NixShell:       nixShell,
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
	jdOutput, jdErr = CreateJobDistributor(&jdInput)
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

func SetupJobs(lggr zerolog.Logger, jdInput jd.Input, nixShell *nix.Shell, registryChainBlockchainOutput *blockchain.Output, topology *cre.Topology, infraInput infra.Input, capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet) (*jd.Output, []*cre.WrappedNodeOutput, error) {
	var jdOutput *jd.Output
	jdAndDonsErrGroup := &errgroup.Group{}

	jdAndDonsErrGroup.Go(func() error {
		var startJDErr error
		jdOutput, startJDErr = StartJD(lggr, nixShell, jdInput, infraInput)
		if startJDErr != nil {
			return pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}

		return nil
	})

	nodeSetOutput := make([]*cre.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))

	jdAndDonsErrGroup.Go(func() error {
		var startDonsErr error
		nodeSetOutput, startDonsErr = StartDONs(lggr, nixShell, topology, infraInput, registryChainBlockchainOutput, capabilitiesAwareNodeSets)
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
