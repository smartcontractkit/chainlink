package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

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

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
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

type metadataCapturingVault struct {
	metadata capabilities.RequestMetadata
	response *vault.GetSecretsResponse
}

func (m *metadataCapturingVault) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.CapabilityInfo{
		ID:             vault.CapabilityID,
		CapabilityType: capabilities.CapabilityTypeAction,
		IsLocal:        true,
	}, nil
}

func (m *metadataCapturingVault) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	vr := &vault.GetSecretsRequest{}
	if err := req.Payload.UnmarshalTo(vr); err != nil {
		return capabilities.CapabilityResponse{}, errors.New("received unexpected payload: want *vault.GetSecretsRequest")
	}
	if req.Method != vault.MethodGetSecrets {
		return capabilities.CapabilityResponse{}, errors.New("received unexpected method: want vault.MethodGetSecrets")
	}
	m.metadata = req.Metadata

	anyvresp, err := anypb.New(m.response)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}
	return capabilities.CapabilityResponse{Payload: anyvresp}, nil
}

func (m *metadataCapturingVault) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return errors.New("not used")
}

func (m *metadataCapturingVault) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return errors.New("not used")
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

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			resp := &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "R1",
							Namespace: "Bar",
							Owner:     normalizedOwner,
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
							Owner:     normalizedOwner,
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
							Owner:     normalizedOwner,
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
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowEncryptionKey,
		nil,
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
	assert.Equal(t, keyFor(normalizedOwner, "Bar", "R1"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, rawSecret, r.Value)

	require.Nil(t, resp[1].GetError())
	r = resp[1].GetSecret()
	assert.Equal(t, keyFor(normalizedOwner, "Bar", "R2"), keyFor(r.Owner, r.Namespace, r.Id))
	assert.Equal(t, rawSecret, r.Value)

	assert.NotNil(t, resp[2].GetError())
	errVal := resp[2].GetError()
	assert.Contains(t, errVal.Error, "failed to aggregate decryption shares")
}

func TestSecretsFetcher_DecryptsBinaryShares(t *testing.T) {
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

	encryptedDecryptionShare0, err := workflowEncryptionKey.Encrypt(decryptionShare0Bytes)
	require.NoError(t, err)
	encryptedDecryptionShare1, err := workflowEncryptionKey.Encrypt(decryptionShare1Bytes)
	require.NoError(t, err)

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "R1",
							Namespace: "Bar",
							Owner:     normalizedOwner,
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										BinaryShares: [][]byte{
											encryptedDecryptionShare0,
											encryptedDecryptionShare1,
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}
	require.NoError(t, reg.Add(t.Context(), mc))

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowEncryptionKey,
		nil,
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{{Id: "R1", Namespace: "Bar"}},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Nil(t, resp[0].GetError())
	assert.Equal(t, rawSecret, resp[0].GetSecret().Value)
}

func TestEncryptedDecryptionShareBytes(t *testing.T) {
	t.Run("prefers binary shares", func(t *testing.T) {
		binary := [][]byte{{1, 2}}
		got, err := encryptedDecryptionShareBytes(binary, []string{"deadbeef"})
		require.NoError(t, err)
		require.Equal(t, binary, got)
	})

	t.Run("falls back to hex shares", func(t *testing.T) {
		got, err := encryptedDecryptionShareBytes(nil, []string{"0102", "0304"})
		require.NoError(t, err)
		require.Equal(t, [][]byte{{1, 2}, {3, 4}}, got)
	})

	t.Run("invalid hex share", func(t *testing.T) {
		_, err := encryptedDecryptionShareBytes(nil, []string{"not-hex"})
		require.Error(t, err)
	})
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
	owner := "1234567890abcdef1234567890abcdef12345678"

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		nil,
	)

	_, err = sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	assert.ErrorContains(t, err, "no compatible capability found")
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

	owner := "1234567890abcdef1234567890abcdef12345678"
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		nil,
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

func TestSecretsFetcher_VaultCapabilityRequestOmitsWorkflowIDMetadata(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()

	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	capture := &metadataCapturingVault{
		response: &vault.GetSecretsResponse{
			Responses: []*vault.SecretResponse{
				{
					Id: &vault.SecretIdentifier{
						Key:       "Foo",
						Namespace: "Bar",
						Owner:     normalizedOwner,
					},
					Result: &vault.SecretResponse_Error{Error: "not found"},
				},
			},
		},
	}
	err = reg.Add(t.Context(), capture)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"org-123",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowEncryptionKey,
		nil,
	)

	_, err = sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{
				Id:        "Foo",
				Namespace: "Bar",
			},
		},
	})
	require.NoError(t, err)
	assert.Empty(t, capture.metadata.WorkflowID, "vault Execute must omit WorkflowID from RequestMetadata for mixed-version DON consensus")
}

func TestSecretsFetcher_VaultBatchLeavesOrgAndWorkflowIdentityUnset(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()

	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	var capturedReq *vault.GetSecretsRequest
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			capturedReq = proto.Clone(req).(*vault.GetSecretsRequest)
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     normalizedOwner,
						},
						Result: &vault.SecretResponse_Error{Error: "not found"},
					},
				},
			}, nil
		},
	}
	err = reg.Add(t.Context(), mc)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"org-123",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowEncryptionKey,
		nil,
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
	require.NotNil(t, capturedReq)
	assert.Equal(t, normalizedOwner, capturedReq.Requests[0].Id.Owner)
	require.Len(t, resp, 1)
	require.NotNil(t, resp[0].GetError())
	assert.Contains(t, resp[0].GetError().Error, "not found")
	assert.NotContains(t, resp[0].GetError().Error, "could not find response")
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

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		nil,
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
	assert.Equal(t, fmt.Sprintf("could not find response for the request: %s::Bar::Foo", normalizedOwner), errVal.Error)
}

func TestSecretsFetcher_ReturnsErrorIfMissingEncryptionSharesForNode(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     normalizedOwner,
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
	err = reg.Add(t.Context(), mc)
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
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		nil,
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

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			resp := &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     normalizedOwner,
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
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		workflowEncryptionKey,
		nil,
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
	assert.Contains(t, errVal.Error, "not enough decryption shares to decrypt the secret")
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

	valueMap, err := values.Wrap(VaultCapabilityRegistryConfig{
		VaultPublicKey: hex.EncodeToString(vaultPublicKey),
		Threshold:      1,
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

func TestSecretsFetcher_EnforcesSecretsCallsLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))

	f, n := 2, 3
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(f, n)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))

	semaphore := limits.WorkflowResourcePoolLimiter[int](5)
	// bound limiter of 1 call allowed
	secretsCallsLimit := limits.NewUpperBoundLimiter[int](1)

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		semaphore,
		secretsCallsLimit,
		nil,
		"",
		"0x1111111111111111111111111111111111111111",
		"wf",
		"wfID",
		"phaseID",
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		nil,
	)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	req := &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{
			{Id: "R1"},
		},
	}

	// 1st call to occupy the only available slot in the limiter
	_, _ = sf.GetSecrets(ctx, req)

	// second call should fail due to exceeding the bound limiter (limit == 1)
	_, err = sf.GetSecrets(ctx, req)
	require.ErrorContains(t, err, "limited: cannot use 2, limit is 1")
}

func TestSecretsFetcher_VaultFirstThenLocalOverridesForVaultFailures(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	workflowKeyBytes := workflowEncryptionKey.PublicKey()

	rawSecret := "vault-secret"
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

	encryptedDecryptionShare0, err := workflowEncryptionKey.Encrypt(decryptionShare0Bytes)
	require.NoError(t, err)
	encryptedDecryptionShare1, err := workflowEncryptionKey.Encrypt(decryptionShare1Bytes)
	require.NoError(t, err)

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	var gotVaultReq *vault.GetSecretsRequest
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			gotVaultReq = proto.Clone(req).(*vault.GetSecretsRequest)
			return &vault.GetSecretsResponse{
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{
							Key:       "vault-only",
							Namespace: vaulttypes.DefaultNamespace,
							Owner:     normalizedOwner,
						},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										Shares: []string{
											hex.EncodeToString(encryptedDecryptionShare0),
											hex.EncodeToString(encryptedDecryptionShare1),
										},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}
	err = reg.Add(t.Context(), mc)
	require.NoError(t, err)

	local := NewLocalSecretsFetcher(owner, map[string]string{"local-only": "local-value"})
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"phaseID",
		workflowEncryptionKey,
		local,
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		CallbackId: 42,
		Requests: []*sdkpb.SecretRequest{
			{Id: "local-only"},
			{Id: "vault-only"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 2)

	require.NotNil(t, gotVaultReq)
	require.Len(t, gotVaultReq.Requests, 2)
	require.Equal(t, "local-only", gotVaultReq.Requests[0].GetId().GetKey())
	require.Equal(t, vaulttypes.DefaultNamespace, gotVaultReq.Requests[0].GetId().GetNamespace())
	require.Equal(t, normalizedOwner, gotVaultReq.Requests[0].GetId().GetOwner())
	require.Equal(t, "vault-only", gotVaultReq.Requests[1].GetId().GetKey())
	require.Equal(t, vaulttypes.DefaultNamespace, gotVaultReq.Requests[1].GetId().GetNamespace())
	require.Equal(t, normalizedOwner, gotVaultReq.Requests[1].GetId().GetOwner())

	localSecret := resp[0].GetSecret()
	require.NotNil(t, localSecret)
	require.Equal(t, "local-only", localSecret.GetId())
	require.Equal(t, normalizedOwner, localSecret.GetOwner())
	require.Equal(t, "local-value", localSecret.GetValue())

	vaultSecret := resp[1].GetSecret()
	require.NotNil(t, vaultSecret)
	require.Equal(t, "vault-only", vaultSecret.GetId())
	require.Equal(t, vaulttypes.DefaultNamespace, vaultSecret.GetNamespace())
	require.Equal(t, normalizedOwner, vaultSecret.GetOwner())
	require.Equal(t, rawSecret, vaultSecret.GetValue())
}

func TestSecretsFetcher_LocalOverridesWhenVaultExecuteFails(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := coreCap.NewRegistry(lggr)
	peer := coreCap.RandomUTF8BytesWord()
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))

	f, n := 2, 3
	_, vaultPublicKey, _, err := tdh2easy.GenerateKeys(f, n)
	require.NoError(t, err)
	vaultPublicKeyBytes, err := vaultPublicKey.Marshal()
	require.NoError(t, err)
	reg.SetLocalRegistry(CreateLocalRegistryWith1Node(t, peer, workflowEncryptionKey.PublicKey(), vaultPublicKeyBytes))

	var vaultCalls int
	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			vaultCalls++
			require.Len(t, req.Requests, 2)
			require.Equal(t, "secret-a", req.Requests[0].GetId().GetKey())
			require.Equal(t, "secret-b", req.Requests[1].GetId().GetKey())
			return nil, errors.New("vault capability execute failed")
		},
	}
	err = reg.Add(t.Context(), mc)
	require.NoError(t, err)

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedLocalOwner, err := normalizeOwner(owner)
	require.NoError(t, err)
	local := NewLocalSecretsFetcher(owner, map[string]string{
		"secret-a": "value-a",
		"secret-b": "value-b",
	})
	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"phaseID",
		workflowEncryptionKey,
		local,
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		CallbackId: 7,
		Requests: []*sdkpb.SecretRequest{
			{Id: "secret-a"},
			{Id: "secret-b"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1, vaultCalls)
	require.Len(t, resp, 2)

	s0 := resp[0].GetSecret()
	require.NotNil(t, s0)
	require.Equal(t, "secret-a", s0.GetId())
	require.Equal(t, normalizedLocalOwner, s0.GetOwner())
	require.Equal(t, "value-a", s0.GetValue())

	s1 := resp[1].GetSecret()
	require.NotNil(t, s1)
	require.Equal(t, "secret-b", s1.GetId())
	require.Equal(t, normalizedLocalOwner, s1.GetOwner())
	require.Equal(t, "value-b", s1.GetValue())
}
