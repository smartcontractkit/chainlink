package environment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	cldf_jd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
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
	tunnelManager tunnel.Manager,
	rewriteInternalForLocalNodes bool,
) (*StartedJD, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")
	if jdConfig == nil {
		return nil, errors.New("jd configuration is nil")
	}

	var jdOutput *jd.Output
	var jdErr error

	if jdConfig.Target == config.TargetRemote {
		startClient, err := newStartComponentClient(lggr, tunnelManager)
		if err != nil {
			return nil, err
		}
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeJD,
			JD:            jdConfig.InputRef(),
			ReusePolicy:   string(jdConfig.RemoteStartPolicy),
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to encode jd payload")
		}
		response, err := startClient.StartComponent(ctx, agent.StartComponentEnvelope{
			SchemaVersion: agent.SchemaVersionV1,
			Operation:     agent.OperationStartComponent,
			Payload:       payloadBytes,
		})
		if err != nil {
			return nil, err
		}
		if response.ComponentType != componentTypeJD {
			return nil, fmt.Errorf("unexpected component type in start response: %s", response.ComponentType)
		}
		for _, logLine := range response.AgentLogs {
			pretty := prettifyAgentLogLine(logLine)
			if pretty == "" {
				continue
			}
			lggr.Info().Msgf("[agent] %s", pretty)
		}
		jdOutput, err = agent.DecodeFromTransport[jd.Output](response.Output)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to decode jd transport payload")
		}
		if err := rewriteRemoteJDOutputForLocalAccess(ctx, lggr, tunnelManager, jdOutput, rewriteInternalForLocalNodes); err != nil {
			return nil, err
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
		jdOutput, jdErr = jd.NewWithContext(ctx, jdConfig.InputRef())
		if jdErr != nil {
			jdErr = fmt.Errorf("failed to start JD container for image %s: %w", jdConfig.Image, jdErr)

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

	jdClientConfig := cldf_jd.JDConfig{
		GRPC:  jdOutput.ExternalGRPCUrl,
		WSRPC: jdOutput.InternalWSRPCUrl,
		Creds: creds,
	}

	lggr.Info().Msgf("Connecting to JD GRPC at: %s", jdOutput.ExternalGRPCUrl)

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

func rewriteRemoteJDOutputForLocalAccess(
	ctx context.Context,
	lggr zerolog.Logger,
	tunnelManager tunnel.Manager,
	output *jd.Output,
	rewriteInternalForLocalNodes bool,
) error {
	if output == nil {
		return nil
	}
	if tunnelManager == nil {
		return errors.New("tunnel manager is required for remote jd target")
	}

	refs, err := describeJDEndpoints(output)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to describe jd tunnel endpoints")
	}
	bindings, err := tunnelManager.Start(ctx, refs)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to start tunnels for jd output")
	}
	for _, binding := range bindings {
		lggr.Info().
			Str("componentID", binding.ComponentID).
			Str("endpointName", binding.EndpointName).
			Str("originalURL", binding.OriginalURL).
			Str("localURL", binding.LocalURL).
			Msg("Established endpoint tunnel")
	}
	return rewriteJDWithBindings(output, bindings, rewriteInternalForLocalNodes)
}

func describeJDEndpoints(output *jd.Output) ([]tunnel.EndpointRef, error) {
	refs := make([]tunnel.EndpointRef, 0, 2)
	componentID := tunnel.CanonicalComponentID(tunnel.KindJD, 0, "job-distributor")

	grpcRef, err := jdEndpointFromAddress(componentID, "grpc", output.ExternalGRPCUrl)
	if err != nil {
		return nil, err
	}
	if grpcRef != nil {
		refs = append(refs, *grpcRef)
	}

	wsrpcRef, err := jdEndpointFromAddress(componentID, "wsrpc", output.ExternalWSRPCUrl)
	if err != nil {
		return nil, err
	}
	if wsrpcRef != nil {
		refs = append(refs, *wsrpcRef)
	}

	return refs, nil
}

func rewriteJDWithBindings(output *jd.Output, bindings []tunnel.TunnelBinding, rewriteInternalForLocalNodes bool) error {
	byName := make(map[string]tunnel.TunnelBinding, len(bindings))
	for _, binding := range bindings {
		byName[binding.EndpointName] = binding
	}

	if output.ExternalGRPCUrl != "" {
		binding, ok := byName["grpc"]
		if !ok {
			return fmt.Errorf("missing tunnel binding for jd grpc endpoint")
		}
		output.ExternalGRPCUrl = net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", binding.LocalPort))
		if rewriteInternalForLocalNodes {
			dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
			output.InternalGRPCUrl = net.JoinHostPort(dockerHost, fmt.Sprintf("%d", binding.LocalPort))
		}
	}

	if output.ExternalWSRPCUrl != "" || output.InternalWSRPCUrl != "" {
		binding, ok := byName["wsrpc"]
		if !ok {
			return fmt.Errorf("missing tunnel binding for jd wsrpc endpoint")
		}
		output.ExternalWSRPCUrl = net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", binding.LocalPort))
		if rewriteInternalForLocalNodes {
			dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
			output.InternalWSRPCUrl = net.JoinHostPort(dockerHost, fmt.Sprintf("%d", binding.LocalPort))
		}
	}

	return nil
}

func jdEndpointFromAddress(componentID, endpointName, rawAddress string) (*tunnel.EndpointRef, error) {
	trimmed := strings.TrimSpace(rawAddress)
	if trimmed == "" {
		return nil, nil
	}

	host := ""
	port := ""

	if strings.Contains(trimmed, "://") {
		parsedURL, err := url.Parse(trimmed)
		if err != nil {
			return nil, fmt.Errorf("failed to parse jd endpoint %q: %w", rawAddress, err)
		}
		host = parsedURL.Hostname()
		port = parsedURL.Port()
	} else {
		parsedHost, parsedPort, err := net.SplitHostPort(trimmed)
		if err != nil {
			return nil, fmt.Errorf("failed to parse jd host:port endpoint %q: %w", rawAddress, err)
		}
		host = parsedHost
		port = parsedPort
	}

	if host == "" || port == "" {
		return nil, fmt.Errorf("jd endpoint %q must contain host and port", rawAddress)
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber <= 0 || portNumber > 65535 {
		return nil, fmt.Errorf("jd endpoint %q has invalid port %q", rawAddress, port)
	}

	return &tunnel.EndpointRef{
		ComponentID:  componentID,
		EndpointName: endpointName,
		Scheme:       "tcp",
		Host:         host,
		Port:         portNumber,
		OriginalURL:  trimmed,
	}, nil
}
