package v2

import (
	"context"
	"encoding/base64"
	"errors"
	"math/big"
	"testing"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
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
	workflowKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowKey.PublicKey()))
	localNode, err := reg.LocalNode(t.Context())
	require.NoError(t, err)

	plainText1 := "encryptedShare1"
	k, n := 2, 3
	_, pk, privateShares, err := tdh2easy.GenerateKeys(k, n)
	require.NoError(t, err)

	cipher, err := tdh2easy.Encrypt(pk, []byte(plainText1))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)
	privateShare0Bytes, err := privateShares[0].Marshal()
	require.NoError(t, err)
	encryptedPrivateShare0, err := workflowKey.Encrypt(privateShare0Bytes)
	require.NoError(t, err)
	privateShare1Bytes, err := privateShares[1].Marshal()
	require.NoError(t, err)
	encryptedPrivateShare1, err := workflowKey.Encrypt(privateShare1Bytes)
	require.NoError(t, err)

	share0, err := tdh2easy.Decrypt(cipher, privateShares[0])
	require.NoError(t, err)
	share1, err := tdh2easy.Decrypt(cipher, privateShares[1])
	require.NoError(t, err)
	share2, err := tdh2easy.Decrypt(cipher, privateShares[2])
	require.NoError(t, err)
	shares := []*tdh2easy.DecryptionShare{share0, share1, share2}
	plaintext, err := tdh2easy.Aggregate(cipher, shares, n)
	require.NoError(t, err)
	assert.Equal(t, plainText1, string(plaintext))

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
								EncryptedValue: base64.StdEncoding.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											base64.StdEncoding.EncodeToString(encryptedPrivateShare0),
											base64.StdEncoding.EncodeToString(encryptedPrivateShare1),
										},
										EncryptionKey: base64.StdEncoding.EncodeToString(localNode.EncryptionPublicKey[:]),
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
								EncryptedValue: base64.StdEncoding.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											base64.StdEncoding.EncodeToString(encryptedPrivateShare0),
											base64.StdEncoding.EncodeToString(encryptedPrivateShare1),
										},
										EncryptionKey: base64.StdEncoding.EncodeToString(localNode.EncryptionPublicKey[:]),
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
	err = reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		workflowKey,
		pk,
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
	require.Nil(t, resp[0].GetError())
	r := resp[0].GetSecret()
	assert.Equal(t, keyFor("owner", "Bar", "Foo"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, plainText1, r.Value)

	assert.Nil(t, resp[1].GetError())
	r = resp[1].GetSecret()
	assert.Equal(t, keyFor("owner", "Bar", "Baz"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, plainText1, r.Value)
}

func TestSecretsFetcher_ReturnsErrorIfCapabilityNoFound(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	reg.SetLocalRegistry(CreateLocalRegistry(t, peer))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		vaultPublicKey,
	)

	_, err = sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
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

	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		vaultPublicKey,
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

	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		vaultPublicKey,
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

func TestSecretsFetcher_ReturnsErrorIfMissingEncryptionSharesForNode(t *testing.T) {
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
										Shares:        []string{"encryptedShare1"},
										EncryptionKey: base64.StdEncoding.EncodeToString([]byte{}),
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

	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		vaultPublicKey,
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
	assert.Contains(t, errVal.Error, "no shares found for this node's encryption key")
}

func TestSecretsFetcher_ReturnsErrorIfCantCombineShares(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	workflowKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowKey.PublicKey()))
	localNode, err := reg.LocalNode(t.Context())
	require.NoError(t, err)

	plainText1 := "encryptedShare1"
	k, n := 2, 3
	_, pk, privateShares, err := tdh2easy.GenerateKeys(k, n)
	require.NoError(t, err)

	cipher, err := tdh2easy.Encrypt(pk, []byte(plainText1))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)
	privateShare0Bytes, err := privateShares[0].Marshal()
	require.NoError(t, err)
	encryptedPrivateShare0, err := workflowKey.Encrypt(privateShare0Bytes)
	require.NoError(t, err)

	share0, err := tdh2easy.Decrypt(cipher, privateShares[0])
	require.NoError(t, err)
	share1, err := tdh2easy.Decrypt(cipher, privateShares[1])
	require.NoError(t, err)
	share2, err := tdh2easy.Decrypt(cipher, privateShares[2])
	require.NoError(t, err)
	shares := []*tdh2easy.DecryptionShare{share0, share1, share2}
	plaintext, err := tdh2easy.Aggregate(cipher, shares, n)
	require.NoError(t, err)
	assert.Equal(t, plainText1, string(plaintext))

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
								EncryptedValue: base64.StdEncoding.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											base64.StdEncoding.EncodeToString(encryptedPrivateShare0),
										},
										EncryptionKey: base64.StdEncoding.EncodeToString(localNode.EncryptionPublicKey[:]),
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
	err = reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		NewSemaphore[[]*sdkpb.SecretResponse](5),
		"owner",
		"workflowName",
		workflowKey,
		pk,
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

	require.Len(t, resp, 1)
	require.NotNil(t, resp[0].GetError())
	errVal := resp[0].GetError()
	assert.Contains(t, errVal.Error, "failed to aggregate decryption shares")
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
		map[string]registrysyncer.Capability{},
	)
	return &localRegistry
}

func CreateLocalRegistryWith1Node(t *testing.T, pid ragetypes.PeerID, encryptionPublicKey [32]byte) *registrysyncer.LocalRegistry {
	workflowDonNodes := []p2ptypes.PeerID{
		pid,
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
				EncryptionPublicKey: encryptionPublicKey,
			},
		},
		map[string]registrysyncer.Capability{},
	)
	return &localRegistry
}
