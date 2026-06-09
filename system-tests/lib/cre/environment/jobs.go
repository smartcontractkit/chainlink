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

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
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

func StartJD(
	ctx context.Context,
	lggr zerolog.Logger,
	jdConfig *config.JobDistributor,
	infraInput infra.Provider,
	remoteRuntime *remoteclient.Runtime,
) (*StartedJD, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")
	if jdConfig == nil {
		return nil, errors.New("jd configuration is nil")
	}

	var jdOutput *jd.Output
	var jdErr error

	switch {
	case jdConfig.Placement == config.PlacementRemote:
		jdOutput, jdErr = remoteclient.StartWithRuntimeDescriptor(
			ctx,
			lggr,
			remoteRuntime,
			remoteclient.StartDescriptor[jd.Output]{
				ComponentType: remoteclient.ComponentTypeJD,
				BuildPayload: func() (agent.StartComponentPayload, error) {
					return agent.StartComponentPayload{
						ComponentType: remoteclient.ComponentTypeJD,
						JD:            jdConfig.InputRef(),
						ReusePolicy:   string(jdConfig.RemoteStartPolicy),
					}, nil
				},
				Rewrite: rewriteJDForDirectAccess,
			},
		)
		if jdErr != nil {
			return nil, jdErr
		}
	case infraInput.IsKubernetes():
		// For Kubernetes, JD is already running in the cluster, generate service URLs
		lggr.Info().Msg("Generating Kubernetes service URLs for Job Distributor (already running in cluster)")
		jdOutput, jdErr = infra.GenerateKubernetesJDOutput(&infraInput, lggr)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to generate Kubernetes JD output")
		}
	default:
		jdOutput, jdErr = jd.NewWithContext(ctx, jdConfig.InputRef())
		if jdErr != nil {
			jdErr = fmt.Errorf("failed to start JD container for image %s: %w", jdConfig.Image, jdErr)

			// useful end user messages
			if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
				jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
			}

			if err := framework.PrintFailedContainerLogs(30); err != nil {
				lggr.Error().Err(err).Msg("failed to print failed Docker container logs")
			}

			return nil, jdErr
		}
	}

	// Configure gRPC credentials for JD connection
	creds := getJDCredentials(lggr, infraInput, jdOutput)

	jdClientConfig := cldf_jd.JDConfig{
		GRPC:  jdOutput.ExternalGRPCUrl,
		WSRPC: jdOutput.ExternalWSRPCUrl,
		Creds: creds,
	}

	lggr.Info().Msgf("Connecting to JD GRPC at: %s", jdOutput.ExternalGRPCUrl)
	lggr.Info().
		Str("internalWSRPC", jdOutput.InternalWSRPCUrl).
		Str("externalWSRPC", jdOutput.ExternalWSRPCUrl).
		Msg("Resolved JD endpoints")

	jdClient, jdErr := cldf_jd.NewJDClient(jdClientConfig)
	if jdErr != nil {
		return nil, pkgerrors.Wrap(jdErr, "failed to create JD client")
	}

	lggr.Info().Msgf("Job Distributor started in %.2f seconds", time.Since(startTime).Seconds())

	return &StartedJD{
		JDOutput: jdOutput,
		Client:   jdClient,
	}, nil
}

func rewriteJDForDirectAccess(output *jd.Output, ec2HostIP string) error {
	if output == nil {
		return nil
	}
	if output.ExternalGRPCUrl != "" {
		rewritten, err := rewriteAddressHost(output.ExternalGRPCUrl, ec2HostIP)
		if err != nil {
			return err
		}
		output.ExternalGRPCUrl = rewritten
	}

	if output.ExternalWSRPCUrl != "" || output.InternalWSRPCUrl != "" {
		source := output.ExternalWSRPCUrl
		if source == "" {
			source = output.InternalWSRPCUrl
		}
		rewritten, err := rewriteAddressHost(source, ec2HostIP)
		if err != nil {
			return err
		}
		output.ExternalWSRPCUrl = rewritten
	}
	return nil
}
