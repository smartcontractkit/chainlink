package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultMock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault/mock"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

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
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	workflowKeyBytes := workflowEncryptionKey.PublicKey()

	rawSecret := "Raw Secret Value"
	f, n := 2, 3
	_, vaultPublicKey, privateShares, err := tdh2easy.GenerateKeys(f, n)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))

	cipher, err := tdh2easy.Encrypt(vaultPublicKey, []byte(rawSecret))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)

	decryptionShare0, err := tdh2easy.Decrypt(cipher, privateShares[0])
	require.NoError(t, err)
	decryptionShare0Bytes, err := decryptionShare0.Marshal()
	require.NoError(t, err)
	decryptionShare1, err := tdh2easy.Decrypt(cipher, privateShares[1])
	require.NoError(t, err)
	decryptionShare1Bytes, err := decryptionShare1.Marshal()
	require.NoError(t, err)
	decryptionShare2, err := tdh2easy.Decrypt(cipher, privateShares[2])
	require.NoError(t, err)
	decryptionShare2Bytes, err := decryptionShare2.Marshal()
	require.NoError(t, err)

	// Sanity testing that we can decrypt the secret with just 2 shares
	twoDecryptionShares := []*tdh2easy.DecryptionShare{decryptionShare0, decryptionShare1}
	decryptedSecret, err := tdh2easy.Aggregate(cipher, twoDecryptionShares, n)
	require.NoError(t, err)
	assert.Equal(t, rawSecret, string(decryptedSecret))

	// Encrypt the decryption shares with the workflow key. This is the expected output from Vault capability.
	encryptedDecryptionShare0, err := workflowEncryptionKey.Encrypt(decryptionShare0Bytes)
	require.NoError(t, err)
	encryptedDecryptionShare1, err := workflowEncryptionKey.Encrypt(decryptionShare1Bytes)
	require.NoError(t, err)
	encryptedDecryptionShare2, err := workflowEncryptionKey.Encrypt(decryptionShare2Bytes)
	require.NoError(t, err)

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			resp := &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "R1",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											hex.EncodeToString(encryptedDecryptionShare0),
											hex.EncodeToString(encryptedDecryptionShare2),
											hex.EncodeToString(encryptedDecryptionShare1),
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
									},
								},
							},
						},
					},
					{
						Id: &vault.SecretIdentifier{
							Key:       "R2",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											hex.EncodeToString(encryptedDecryptionShare1),
											hex.EncodeToString(encryptedDecryptionShare0),
											hex.EncodeToString([]byte("junk value")),
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
									},
								},
							},
						},
					},
					{
						Id: &vault.SecretIdentifier{
							Key:       "R3",
							Namespace: "Bar",
							Owner:     "owner",
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											hex.EncodeToString(encryptedDecryptionShare0),
											// deliberately supplying less than threshold shares
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
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
		limits.WorkflowResourcePoolLimiter[int](5),
		"owner",
		"workflowName",
		workflowEncryptionKey,
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "R1",
				Namespace: "Bar",
			},
			{
				Id:        "R2",
				Namespace: "Bar",
			},
			{
				Id:        "R3",
				Namespace: "Bar",
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, resp, 3)
	require.Nil(t, resp[0].GetError())
	r := resp[0].GetSecret()
	assert.Equal(t, keyFor("owner", "Bar", "R1"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, rawSecret, r.Value)

	require.Nil(t, resp[1].GetError())
	r = resp[1].GetSecret()
	assert.Equal(t, keyFor("owner", "Bar", "R2"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, rawSecret, r.Value)

	assert.NotNil(t, resp[2].GetError())
	errVal := resp[2].GetError()
	assert.Contains(t, errVal.Error, "failed to aggregate decryption shares")
}

func TestSecretsFetcher_ReturnsErrorIfCapabilityNoFound(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
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

	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
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
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{},
			}, nil
		},
	}
	err := reg.Add(t.Context(), mc)
	require.NoError(t, err)

	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
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
	assert.Equal(t, "could not find response for the request: owner::Bar::Foo", errVal.Error)
}

func TestSecretsFetcher_ReturnsErrorIfMissingEncryptionSharesForNode(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
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
										EncryptionKey: hex.EncodeToString([]byte{}),
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

	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		"owner",
		"workflowName",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
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
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	workflowKeyBytes := workflowEncryptionKey.PublicKey()

	plainText1 := "encryptedShare1"
	f, n := 2, 3
	_, vaultPublicKey, privateShares, err := tdh2easy.GenerateKeys(f, n)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))

	cipher, err := tdh2easy.Encrypt(vaultPublicKey, []byte(plainText1))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)
	privateShare0Bytes, err := privateShares[0].Marshal()
	require.NoError(t, err)
	encryptedPrivateShare0, err := workflowEncryptionKey.Encrypt(privateShare0Bytes)
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
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											hex.EncodeToString(encryptedPrivateShare0),
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
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
		limits.WorkflowResourcePoolLimiter[int](5),
		"owner",
		"workflowName",
		workflowEncryptionKey,
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
		map[p2ptypes.PeerID]registrysyncer.NodeInfo{
			workflowDonNodes[0]: {
				NodeOperatorID:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pID:               workflowDonNodes[0],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
			workflowDonNodes[1]: {
				NodeOperatorID:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pID:               workflowDonNodes[1],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
			workflowDonNodes[2]: {
				NodeOperatorID:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pID:               workflowDonNodes[2],
				EncryptionPublicKey: coreCap.RandomUTF8BytesWord(),
			},
			workflowDonNodes[3]: {
				NodeOperatorID:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pID:               workflowDonNodes[3],
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

func CreateLocalRegistryWith1Node(t *testing.T, pid ragetypes.PeerID, workflowPublicKey [32]byte, vaultPublicKey []byte) *registrysyncer.LocalRegistry {
	workflowDonNodes := []p2ptypes.PeerID{
		pid,
	}

	valueMap, err := values.NewMap[string](map[string]string{
		"VaultPublicKey": hex.EncodeToString(vaultPublicKey),
	})
	require.NoError(t, err)
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(valueMap).GetMapValue(),
	}
	configb, err := proto.Marshal(config)

	require.NoError(t, err)
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
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					vault.CapabilityID: {
						Config: configb,
					},
				},
			},
		},
		map[p2ptypes.PeerID]registrysyncer.NodeInfo{
			workflowDonNodes[0]: {
				NodeOperatorID:      1,
				WorkflowDONId:       dID,
				Signer:              coreCap.RandomUTF8BytesWord(),
				P2pID:               workflowDonNodes[0],
				EncryptionPublicKey: workflowPublicKey,
			},
		},
		map[string]registrysyncer.Capability{
			vault.CapabilityID: {
				CapabilityType: capabilities.CapabilityTypeAction,
				ID:             vault.CapabilityID,
			},
		},
	)
	return &localRegistry
}
