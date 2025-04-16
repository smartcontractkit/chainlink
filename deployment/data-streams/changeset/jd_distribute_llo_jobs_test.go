package changeset

import (
	"context"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

func TestDistributeLLOJobSpecs(t *testing.T) {
	t.Parallel()

	const donID = 1
	const donName = "don"
	const envName = "env"

	e := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      false,
		ShouldDeployLinkToken: false,
		NumNodes:              1,
		NumBootstrapNodes:     1,
		NodeLabels: []*ptypes.Label{
			{
				Key:   "product",
				Value: pointer.To(jd.ProductLabel),
			},
			{
				Key:   "environment",
				Value: pointer.To(envName),
			},
			{
				Key: utils.DonIdentifier(donID, donName),
			},
		},
	}).Environment

	// Partially mock the JD client, so it's ProposeJob just return.
	e.Offchain = NewAdHocOffchainMock(e.Offchain)

	// pick the first EVM chain selector
	chainSelector := e.AllChainSelectors()[0]

	// insert a Configurator address for the given DON
	configuratorAddr := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err := e.ExistingAddresses.Save(chainSelector, configuratorAddr,
		deployment.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  deployment.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	oracleSpec := `name = 'don | 1'
type = 'offchainreporting2'
schemaVersion = 1
contractID = '0x4170ed0880ac9a755fd29b2688956bd959f923f4'
ocrKeyBundleID = 'cee9d802bf0e28bc74c78d7512e44b25ce6580bf5c45ed15186ae871a3437eb1'
maxTaskDuration = '1s'
contractConfigTrackerPollInterval = '1s'
relay = 'evm'
pluginType = 'llo'

[relayConfig]
chainID = '90000001'
lloConfigMode = 'bluegreen'
lloDonID = 1

[pluginConfig]
channelDefinitionsContractAddress = '0x000000000000000000000000000000000000dEaD'
channelDefinitionsContractFromBlock = 0
donID = 1
servers = {'mercury-pipeline-testnet-producer.TEST.cldev.cloud:1340' = '0000005187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c'}
`

	bootstrapSpec := `name = 'bootstrap'
type = 'bootstrap'
schemaVersion = 1
contractID = '0x4170ed0880ac9a755fd29b2688956bd959f923f4'
donID = 1
relay = 'evm'

[relayConfig]
chainID = '90000001'
`

	config := CsDistributeLLOJobSpecsConfig{
		ChainSelectorEVM: chainSelector,
		Filter: &jd.ListFilter{
			DONID:    donID,
			DONName:  donName,
			EnvLabel: envName,
			Size:     1,
		},
		FromBlock:                   0,
		ConfigMode:                  "bluegreen",
		ChannelConfigStoreAddr:      common.HexToAddress("DEAD"),
		ChannelConfigStoreFromBlock: 0,
		ConfiguratorAddress:         configuratorAddr,
		Servers: map[string]string{
			"mercury-pipeline-testnet-producer.TEST.cldev.cloud:1340": "0000005187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c",
		},
	}

	tests := []struct {
		name              string
		env               deployment.Environment
		config            CsDistributeLLOJobSpecsConfig
		prepConfFn        func(CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig
		wantErr           *string
		wantOracleSpec    string
		wantBootstrapSpec string
	}{
		{
			name:              "success",
			env:               e,
			config:            config,
			wantOracleSpec:    oracleSpec,
			wantBootstrapSpec: bootstrapSpec,
		},
		{
			name:   "missing channel config store",
			env:    e,
			config: config,
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.ChannelConfigStoreAddr = common.Address{}
				return c
			},
			wantErr: pointer.To("channel config store address is required"),
		},
		{
			name:   "missing servers",
			env:    e,
			config: config,
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.Servers = nil
				return c
			},
			wantErr: pointer.To("servers map is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := tt.config
			if tt.prepConfFn != nil {
				conf = tt.prepConfFn(tt.config)
			}
			_, out, err := changeset.ApplyChangesetsV2(t,
				tt.env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(CsDistributeLLOJobSpecs{}, conf),
				},
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), *tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, 2)

			t.Log(out[0].Jobs[0].Spec)
			t.Log(out[0].Jobs[1].Spec)
			// These are lines with dynamic values which we cannot compare.
			linesToStrip := []string{"externalJobID", "transmitterID", "p2pv2Bootstrappers", "ocrKeyBundleID"}
			spec1 := testutil.StripLineContaining(out[0].Jobs[0].Spec, linesToStrip)
			spec2 := testutil.StripLineContaining(out[0].Jobs[1].Spec, linesToStrip)
			wantBootstrapSpec := testutil.StripLineContaining(tt.wantBootstrapSpec, linesToStrip)
			wantOracleSpec := testutil.StripLineContaining(tt.wantOracleSpec, linesToStrip)

			if strings.Contains(spec1, "bootstrap") {
				require.Equal(t, wantBootstrapSpec, spec1)
				require.Equal(t, wantOracleSpec, spec2)
			} else {
				require.Equal(t, wantOracleSpec, spec1)
				require.Equal(t, wantBootstrapSpec, spec2)
			}
		})
	}
}

type adHocOffchainMock struct {
	RealOffchain deployment.OffchainClient
}

func NewAdHocOffchainMock(real deployment.OffchainClient) *adHocOffchainMock {
	return &adHocOffchainMock{
		RealOffchain: real,
	}
}

// JobServiceClient interface
func (a *adHocOffchainMock) GetJob(ctx context.Context, in *job.GetJobRequest, opts ...grpc.CallOption) (*job.GetJobResponse, error) {
	return a.RealOffchain.GetJob(ctx, in, opts...)
}
func (a *adHocOffchainMock) GetProposal(ctx context.Context, in *job.GetProposalRequest, opts ...grpc.CallOption) (*job.GetProposalResponse, error) {
	return a.RealOffchain.GetProposal(ctx, in, opts...)
}
func (a *adHocOffchainMock) ListJobs(ctx context.Context, in *job.ListJobsRequest, opts ...grpc.CallOption) (*job.ListJobsResponse, error) {
	return a.RealOffchain.ListJobs(ctx, in, opts...)
}
func (a *adHocOffchainMock) ListProposals(ctx context.Context, in *job.ListProposalsRequest, opts ...grpc.CallOption) (*job.ListProposalsResponse, error) {
	return a.RealOffchain.ListProposals(ctx, in, opts...)
}
func (a *adHocOffchainMock) ProposeJob(ctx context.Context, in *job.ProposeJobRequest, opts ...grpc.CallOption) (*job.ProposeJobResponse, error) {
	return &job.ProposeJobResponse{
		Proposal: &job.Proposal{
			JobId: uuid.New().String(), // Maybe replace this with a fixed value for testing
			Spec:  in.Spec,
		},
	}, nil
}
func (a *adHocOffchainMock) BatchProposeJob(ctx context.Context, in *job.BatchProposeJobRequest, opts ...grpc.CallOption) (*job.BatchProposeJobResponse, error) {
	return a.RealOffchain.BatchProposeJob(ctx, in, opts...)
}
func (a *adHocOffchainMock) RevokeJob(ctx context.Context, in *job.RevokeJobRequest, opts ...grpc.CallOption) (*job.RevokeJobResponse, error) {
	return a.RealOffchain.RevokeJob(ctx, in, opts...)
}
func (a *adHocOffchainMock) DeleteJob(ctx context.Context, in *job.DeleteJobRequest, opts ...grpc.CallOption) (*job.DeleteJobResponse, error) {
	return a.RealOffchain.DeleteJob(ctx, in, opts...)
}
func (a *adHocOffchainMock) UpdateJob(ctx context.Context, in *job.UpdateJobRequest, opts ...grpc.CallOption) (*job.UpdateJobResponse, error) {
	return a.RealOffchain.UpdateJob(ctx, in, opts...)
}

// NodeServiceClient interface
func (a *adHocOffchainMock) DisableNode(ctx context.Context, in *node.DisableNodeRequest, opts ...grpc.CallOption) (*node.DisableNodeResponse, error) {
	return a.RealOffchain.DisableNode(ctx, in, opts...)
}
func (a *adHocOffchainMock) EnableNode(ctx context.Context, in *node.EnableNodeRequest, opts ...grpc.CallOption) (*node.EnableNodeResponse, error) {
	return a.RealOffchain.EnableNode(ctx, in, opts...)
}
func (a *adHocOffchainMock) GetNode(ctx context.Context, in *node.GetNodeRequest, opts ...grpc.CallOption) (*node.GetNodeResponse, error) {
	return a.RealOffchain.GetNode(ctx, in, opts...)
}
func (a *adHocOffchainMock) ListNodes(ctx context.Context, in *node.ListNodesRequest, opts ...grpc.CallOption) (*node.ListNodesResponse, error) {
	return a.RealOffchain.ListNodes(ctx, in, opts...)
}
func (a *adHocOffchainMock) ListNodeChainConfigs(ctx context.Context, in *node.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*node.ListNodeChainConfigsResponse, error) {
	return a.RealOffchain.ListNodeChainConfigs(ctx, in, opts...)
}
func (a *adHocOffchainMock) RegisterNode(ctx context.Context, in *node.RegisterNodeRequest, opts ...grpc.CallOption) (*node.RegisterNodeResponse, error) {
	return a.RealOffchain.RegisterNode(ctx, in, opts...)
}
func (a *adHocOffchainMock) UpdateNode(ctx context.Context, in *node.UpdateNodeRequest, opts ...grpc.CallOption) (*node.UpdateNodeResponse, error) {
	return a.RealOffchain.UpdateNode(ctx, in, opts...)
}

// CSAServiceClient interface
func (a *adHocOffchainMock) GetKeypair(ctx context.Context, in *csa.GetKeypairRequest, opts ...grpc.CallOption) (*csa.GetKeypairResponse, error) {
	return a.RealOffchain.GetKeypair(ctx, in, opts...)
}
func (a *adHocOffchainMock) ListKeypairs(ctx context.Context, in *csa.ListKeypairsRequest, opts ...grpc.CallOption) (*csa.ListKeypairsResponse, error) {
	return a.RealOffchain.ListKeypairs(ctx, in, opts...)
}
