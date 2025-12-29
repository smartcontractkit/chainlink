package environment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	cldf_jd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type StartedJD struct {
	JDOutput *jd.Output
	Client   *cldf_jd.JobDistributor
}

// getJDCredentials determines the appropriate gRPC credentials for JD connection
func getJDCredentials(lggr zerolog.Logger, infraInput infra.Provider, jdOutput *jd.Output) credentials.TransportCredentials {
	// Determine if TLS should be used based on configuration or default port-based logic
	creds := insecure.NewCredentials()

	if infraInput.IsKubernetes() {
		// For Kubernetes, check if TLS is explicitly configured or default to TLS for port 443
		useTLS := false
		if infraInput.Kubernetes != nil && infraInput.Kubernetes.UseTLSForJD != nil {
			useTLS = *infraInput.Kubernetes.UseTLSForJD
		} else {
			// Default behavior: use TLS for port 443
			useTLS = strings.Contains(jdOutput.ExternalGRPCUrl, ":443")
		}

		if useTLS {
			// Passing nil uses the system cert pool for TLS verification.
			creds = credentials.NewTLS(nil)
			lggr.Info().Msg("Using TLS credentials for JD GRPC connection")
		} else {
			lggr.Info().Msg("Using insecure credentials for JD GRPC connection (Kubernetes)")
		}
	} else {
		lggr.Info().Msg("Using insecure credentials for JD GRPC connection (non-Kubernetes)")
	}

	return creds
}

func StartJD(ctx context.Context, lggr zerolog.Logger, jdInput jd.Input, infraInput infra.Provider) (*StartedJD, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")

	var jdOutput *jd.Output
	var jdErr error

	if infraInput.Type == infra.CRIB {
		deployCribJdInput := &crib.DeployCribJdInput{
			JDInput:        jdInput,
			CribConfigsDir: infra.CribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		jdInput.Out, jdErr = crib.DeployJd(ctx, deployCribJdInput)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
		}
	} else if infraInput.IsKubernetes() {
		// For Kubernetes, JD is already running in the cluster, generate service URLs
		lggr.Info().Msg("Generating Kubernetes service URLs for Job Distributor (already running in cluster)")
		jdOutput, jdErr = infra.GenerateKubernetesJDOutput(&infraInput, lggr)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to generate Kubernetes JD output")
		}
	}

	// Only start JD container for Docker provider
	if jdOutput == nil {
		jdOutput, jdErr = jd.NewWithContext(ctx, &jdInput)
		if jdErr != nil {
			jdErr = fmt.Errorf("failed to start JD container for image %s: %w", jdInput.Image, jdErr)

			// useful end user messages
			if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
				jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
			}

			infra.PrintFailedContainerLogs(lggr, 30)

			return nil, jdErr
		}
	}

	// Configure gRPC credentials for JD connection
	creds := getJDCredentials(lggr, infraInput, jdOutput)

	jdConfig := cldf_jd.JDConfig{
		GRPC:  jdOutput.ExternalGRPCUrl,
		WSRPC: jdOutput.InternalWSRPCUrl,
		Creds: creds,
	}

	lggr.Info().Msgf("Connecting to JD GRPC at: %s", jdOutput.ExternalGRPCUrl)

	jdClient, jdErr := cldf_jd.NewJDClient(jdConfig)
	if jdErr != nil {
		return nil, pkgerrors.Wrap(jdErr, "failed to create JD client")
	}

	lggr.Info().Msgf("Job Distributor started in %.2f seconds", time.Since(startTime).Seconds())

	return &StartedJD{
		JDOutput: jdOutput,
		Client:   jdClient,
	}, nil
}
