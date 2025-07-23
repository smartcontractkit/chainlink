package v2

import (
	"context"
	"errors"
	"testing"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultMock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault/mock"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func MetricsLabelerTest(t *testing.T) *monitoring.WorkflowsMetricLabeler {
	m, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	l := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), m)
	return l
}

func TestSecretsFetcher_BulkFetchesSecretsFromCapability(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			resp := &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{"encryptedShare1"},
									},
								},
							},
						},
					},
					{
						Id: &vault.SecretIdentifier{
							Key:       "Baz",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{"encryptedShare2"},
									},
								},
							},
						},
					},
				},
			}
			return resp, nil
		},
	}
	err := reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		func(shares []string) (string, error) {
			if shares[0] == "encryptedShare1" {
				return "revealedShare1", nil
			}
			if shares[0] == "encryptedShare2" {
				return "revealedShare2", nil
			}

			return "", errors.New("unexpected shares")
		},
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
			{
				Id:        "Baz",
				Namespace: "Bar",
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, resp, 2)
	assert.Nil(t, resp[0].GetError())
	r := resp[0].GetSecret()
	assert.Equal(t, keyFor("owner", "Bar", "Foo"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, "revealedShare1", r.Value)

	assert.Nil(t, resp[1].GetError())
	r = resp[1].GetSecret()
	assert.Equal(t, keyFor("owner", "Bar", "Baz"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, "revealedShare2", r.Value)
}

func TestSecretsFetcher_ReturnsErrorIfCapabilityNoFound(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		func(shares []string) (string, error) { return "", nil },
	)

	_, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	assert.ErrorContains(t, err, "capability not found")
}

func TestSecretsFetcher_ReturnsErrorIfCapabilityErrors(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return nil, errors.New("could not authorize the request")
		},
	}
	err := reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		func(shares []string) (string, error) {
			return "", nil
		},
	)

	_, err = sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	require.ErrorContains(t, err, "could not authorize the request")
}

func TestSecretsFetcher_ReturnsErrorIfNoResponseForRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{},
			}, nil
		},
	}
	err := reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		func(shares []string) (string, error) {
			return "", nil
		},
	)
	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, resp, 1)
	assert.NotNil(t, resp[0].GetError())
	errVal := resp[0].GetError()
	assert.Equal(t, "could not find secret for owner::Bar::Foo", errVal.Error)
}

func TestSecretsFetcher_ReturnsErrorIfTooManyDecryptionShares(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{"encryptedShare1"},
									},
									{
										Shares: []string{"encryptedShare2"},
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}
	err := reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		func(shares []string) (string, error) {
			return "", nil
		},
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, resp, 1)
	assert.NotNil(t, resp[0].GetError())
	errVal := resp[0].GetError()
	assert.Equal(t, "unexpected error when getting secret for owner::Bar::Foo", errVal.Error)
}

func TestSecretsFetcher_ReturnsErrorIfCantCombineShares(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{"encryptedShare1"},
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}
	err := reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		func(shares []string) (string, error) {
			return "", errors.New("could not combine shares")
		},
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, resp, 1)
	assert.NotNil(t, resp[0].GetError())
	errVal := resp[0].GetError()
	assert.Equal(t, "unexpected error when getting secret for owner::Bar::Foo", errVal.Error)
}

func CreateLocalRegistry(t *testing.T, pid ragetypes.PeerID) *registrysyncer.LocalRegistry {
	workflowDonNodes := []p2ptypes.PeerID{
		pid,
		coreCap.RandomUTF8BytesWord(),
		coreCap.RandomUTF8BytesWord(),
		coreCap.RandomUTF8BytesWord(),
	}

	dID := uint32(1)
	localRegistry := registrysyncer.NewLocalRegistry(
		logger.TestLogger(t),
		func() (p2ptypes.PeerID, error) { return pid, nil },
		map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(2),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
		},
		map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			workflowDonNodes[0]: {
				NodeOperatorId:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pId:               workflowDonNodes[0],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
			workflowDonNodes[1]: {
				NodeOperatorId:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pId:               workflowDonNodes[1],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
			workflowDonNodes[2]: {
				NodeOperatorId:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pId:               workflowDonNodes[2],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
			workflowDonNodes[3]: {
				NodeOperatorId:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pId:               workflowDonNodes[3],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
		},
		map[string]registrysyncer.Capability{
			"test-target@1.0.0": {
				CapabilityType: capabilities.CapabilityTypeTarget,
				ID:             "write-chain@1.0.1",
			},
		},
	)
	return &localRegistry
}
