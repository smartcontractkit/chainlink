package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-common/pkg/billing"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http/server"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

const (
	defaultMaxUncompressedBinarySize = 1000000000
	defaultRPS                       = 1000.0
	defaultBurst                     = 1000
	defaultWorkflowID                = "1111111111111111111111111111111111111111111111111111111111111111"
	defaultOwner                     = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	defaultName                      = "myworkflow"
)

var (
	defaultTimeout = 10 * time.Minute
)

func NewStandaloneEngine(
	ctx context.Context,
	lggr logger.Logger,
	registry *capabilities.Registry,
	binary, config, secrets []byte,
	billingClientAddr string,
	lifecycleHooks v2.LifecycleHooks,
	workflowName string,
) (services.Service, []*sdkpb.TriggerSubscription, error) {
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: defaultOwner, Workflow: defaultWorkflowID})
	labeler := custmsg.NewLabeler()
	moduleConfig := &host.ModuleConfig{
		Logger:                  lggr,
		Labeler:                 labeler,
		MaxCompressedBinarySize: defaultMaxUncompressedBinarySize,
		IsUncompressed:          true,
		Timeout:                 &defaultTimeout,
	}

	module, err := host.NewModule(moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create module from config: %w", err)
	}

	if workflowName == "" {
		workflowName = defaultName
	}

	name, err := types.NewWorkflowName(workflowName)
	if err != nil {
		return nil, nil, err
	}

	lf := limits.Factory{Logger: logger.Named(lggr, "Limits")}
	limiters, err := v2.NewLimiters(lf, nil)
	if err != nil {
		return nil, nil, err
	}
	rl, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      defaultRPS,
		GlobalBurst:    defaultBurst,
		PerSenderRPS:   defaultRPS,
		PerSenderBurst: defaultBurst,
	}, lf)
	if err != nil {
		return nil, nil, err
	}
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{
		Global:   1000000000,
		PerOwner: 1000000000,
	}, lf)
	if err != nil {
		return nil, nil, err
	}

	var billingClient billing.WorkflowClient
	if billingClientAddr != "" {
		clientOpts := []billing.WorkflowClientOpt{}
		if strings.HasPrefix(billingClientAddr, "https") {
			clientOpts = append(clientOpts, billing.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		billingClient, _ = billing.NewWorkflowClient(lggr, billingClientAddr, clientOpts...)
	}

	if module.IsLegacyDAG() {
		sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
		if err != nil {
			return nil, nil, err
		}

		cfg := workflows.Config{
			Lggr:                 lggr,
			Workflow:             *sdkSpec,
			WorkflowID:           defaultWorkflowID,
			WorkflowOwner:        defaultOwner,
			WorkflowName:         name,
			Registry:             registry,
			Store:                store.NewInMemoryStore(lggr, clockwork.NewRealClock()),
			Config:               config,
			Binary:               binary,
			SecretsFetcher:       SecretsFor,
			RateLimiter:          rl,
			WorkflowLimits:       workflowLimits,
			NewWorkerTimeout:     time.Minute,
			StepTimeout:          time.Minute,
			MaxExecutionDuration: time.Minute,
			BillingClient:        billingClient,
		}

		engine, err := workflows.NewEngine(ctx, cfg)
		if err != nil {
			return nil, nil, err
		}
		return engine, nil, nil
	}

	secretsFetcher, err := NewFileBasedSecrets(secrets)
	if err != nil {
		return nil, nil, err
	}

	cfg := &v2.EngineConfig{
		Lggr:            lggr,
		Module:          module,
		WorkflowConfig:  config,
		CapRegistry:     registry,
		DonTimeStore:    dontime.NewStore(dontime.DefaultRequestTimeout),
		ExecutionsStore: store.NewInMemoryStore(lggr, clockwork.NewRealClock()),

		WorkflowID:    defaultWorkflowID,
		WorkflowOwner: defaultOwner,
		WorkflowName:  name,
		WorkflowTag:   "workflowTag",

		LocalLimits:                       v2.EngineLimits{},
		LocalLimiters:                     limiters,
		GlobalExecutionConcurrencyLimiter: workflowLimits,
		GlobalExecutionRateLimiter:        rl,

		BeholderEmitter: custmsg.NewLabeler(),

		BillingClient: billingClient,
		Hooks:         lifecycleHooks,

		SecretsFetcher: secretsFetcher,
		DebugMode:      true,
	}

	engine, err := v2.NewEngine(cfg)
	if err != nil {
		return nil, nil, err
	}

	moduleExecuteMaxResponseSizeBytes, err := cfg.LocalLimiters.ExecutionResponse.Limit(ctx)
	if err != nil {
		return nil, nil, err
	}
	if moduleExecuteMaxResponseSizeBytes < 0 {
		return nil, nil, fmt.Errorf("invalid moduleExecuteMaxResponseSizeBytes; must not be negative: %d", moduleExecuteMaxResponseSizeBytes)
	}
	result, err := module.Execute(ctx, &sdkpb.ExecuteRequest{
		Request:         &sdkpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(moduleExecuteMaxResponseSizeBytes), //nolint:gosec // G115
		Config:          config,
	}, v2.NewDisallowedExecutionHelper(lggr, nil, &types.LocalTimeProvider{}, secretsFetcher))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute subscribe: %w", err)
	}
	if result.GetError() != "" {
		return nil, nil, fmt.Errorf("failed to execute subscribe: %s", result.GetError())
	}
	triggerSubscriptions := result.GetTriggerSubscriptions()

	return &serviceWithClosers{engine, []io.Closer{limiters, workflowLimits, rl}}, triggerSubscriptions.GetSubscriptions(), nil
}

type serviceWithClosers struct {
	services.Service
	closers []io.Closer
}

func (s *serviceWithClosers) Close() error {
	return errors.Join(s.Service.Close(), services.MultiCloser(s.closers).Close())
}

// yamlConfig represents the structure of your secrets.yaml file.
type yamlConfig struct {
	SecretsNames map[string][]string `yaml:"secretsNames"`
}

type fileBasedSecrets struct {
	secrets yamlConfig
}

func NewFileBasedSecrets(secrets []byte) (*fileBasedSecrets, error) {
	fbs := new(fileBasedSecrets)
	if err := yaml.Unmarshal(secrets, &fbs.secrets); err != nil {
		return nil, err
	}

	return fbs, nil
}

func (f *fileBasedSecrets) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	responses := make([]*sdkpb.SecretResponse, 0, len(request.Requests))
	for _, req := range request.Requests {
		values, ok := f.secrets.SecretsNames[req.Id]

		// Handle secret not found
		if !ok {
			responses = append(responses, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Error: "secret not found",
					},
				},
			})
			continue
		}

		// Handle secret found but no value associated
		if len(values) == 0 {
			responses = append(responses, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Error: "secret found but no value associated"},
				},
			})
			continue
		}

		// Secret found with value
		secret := &sdkpb.Secret{
			Id:        req.Id,
			Namespace: req.Namespace, // Use the namespace from the request
			Value:     values[0],     // Take the first value as the secret
		}
		responses = append(responses, &sdkpb.SecretResponse{
			Response: &sdkpb.SecretResponse_Secret{
				Secret: secret,
			},
		})
	}
	return responses, nil
}

// TODO support fetching secrets (from a local file)
func SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return map[string]string{}, nil
}

// NewCapabilities builds capabilities using latest standard capabilities where possible, otherwise filled in with faked capabilities.
// Capabilities are then registered with the capability registry.
func NewCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps, err := NewFakeCapabilities(ctx, lggr, registry)
	if err != nil {
		return nil, err
	}

	caps = append(caps, newStandardCapabilities(standardCapabilities, lggr, registry)...)

	return caps, nil
}

func NewFakeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0)

	streamsTrigger := fakes.NewFakeStreamsTrigger(lggr, 6)
	if err := registry.Add(ctx, streamsTrigger); err != nil {
		return nil, err
	}
	caps = append(caps, streamsTrigger)

	httpAction := fakes.NewDirectHTTPAction(lggr)
	if err := registry.Add(ctx, httpserver.NewClientServer(httpAction)); err != nil {
		return nil, err
	}
	caps = append(caps, httpAction)

	fakeConsensus, err := fakes.NewFakeConsensus(lggr, fakes.DefaultFakeConsensusConfig())
	if err != nil {
		return nil, err
	}
	if err := registry.Add(ctx, fakeConsensus); err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensus)

	// generate deterministic signers - need to be configured on the Forwarder contract
	nSigners := 4
	signers := []ocr2key.KeyBundle{}
	for range nSigners {
		signer := ocr2key.MustNewInsecure(fakes.SeedForKeys(), chaintype.EVM)
		lggr.Infow("Generated new consensus signer", "addrss", common.BytesToAddress(signer.PublicKey()))
		signers = append(signers, signer)
	}
	fakeConsensusNoDAG := fakes.NewFakeConsensusNoDAG(signers, lggr)
	if err := registry.Add(ctx, consensusserver.NewConsensusServer(fakeConsensusNoDAG)); err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensusNoDAG)

	writers := []string{"write_aptos-testnet@1.0.0"}
	for _, writer := range writers {
		writeCap := fakes.NewFakeWriteChain(lggr, writer)
		if err := registry.Add(ctx, writeCap); err != nil {
			return nil, err
		}
		caps = append(caps, writeCap)
	}

	return caps, nil
}
