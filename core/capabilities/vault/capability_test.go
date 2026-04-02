package vault

import (
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCapability_CapabilityCall(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := limits.Factory{Settings: cresettings.DefaultGetter}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, nil, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	owner := "test-owner"
	workflowID := "test-workflow-id"
	workflowExecutionID := "test-workflow-execution-id"
	referenceID := "test-reference-id"

	requestID := fmt.Sprintf("%s::%s::%s", workflowID, workflowExecutionID, referenceID)

	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}

	gsr := &vault.GetSecretsRequest{
		WorkflowOwner: owner,
		Requests: []*vault.SecretRequest{
			{
				Id:             sid,
				EncryptionKeys: []string{"key"},
			},
		},
	}

	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	expectedResponse := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: sid,
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{Shares: []string{"share1", "share2"}},
							{Shares: []string{"share3", "share4"}},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(expectedResponse)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) == 1 {
					req := reqs[0]
					req.SendResponse(t.Context(), &vaulttypes.Response{
						ID:      requestID,
						Payload: data,
					})
					return
				}
			}
		}
	}()

	resp, err := capability.Execute(t.Context(), capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          workflowID,
			WorkflowExecutionID: workflowExecutionID,
			ReferenceID:         referenceID,
		},
	})
	wg.Wait()

	require.NoError(t, err)
	typedResponse := &vault.GetSecretsResponse{}
	err = resp.Payload.UnmarshalTo(typedResponse)
	require.NoError(t, err)
	assert.True(t, proto.Equal(expectedResponse, typedResponse))
}

func TestCapability_CapabilityCall_DuringSubscriptionPhase(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := limits.Factory{Settings: cresettings.DefaultGetter}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, nil, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	owner := "test-owner"
	workflowID := "test-workflow-id"
	referenceID := "0"

	requestID := fmt.Sprintf("%s::%s::%s", workflowID, "subscription", referenceID)

	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}

	gsr := &vault.GetSecretsRequest{
		WorkflowOwner: owner,
		Requests: []*vault.SecretRequest{
			{
				Id:             sid,
				EncryptionKeys: []string{"key"},
			},
		},
	}

	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	expectedResponse := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: sid,
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{Shares: []string{"share1", "share2"}},
							{Shares: []string{"share3", "share4"}},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(expectedResponse)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) == 1 {
					req := reqs[0]
					req.SendResponse(t.Context(), &vaulttypes.Response{
						ID:      requestID,
						Payload: data,
					})
					return
				}
			}
		}
	}()

	resp, err := capability.Execute(t.Context(), capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          workflowID,
			WorkflowExecutionID: "", // Empty execution ID indicates subscription phase
			ReferenceID:         referenceID,
		},
	})
	wg.Wait()

	require.NoError(t, err)
	typedResponse := &vault.GetSecretsResponse{}
	err = resp.Payload.UnmarshalTo(typedResponse)
	require.NoError(t, err)
	assert.True(t, proto.Equal(expectedResponse, typedResponse))
}

func TestCapability_CapabilityCall_SecretIdentifierOwnerMismatch(t *testing.T) {
	testCases := []struct {
		name          string
		workflowOwner string
		secretOwners  []string
		shouldReject  bool
	}{
		{
			name:          "mismatched owner",
			workflowOwner: "0xABCDef1234567890abcdef1234567890abcdef12",
			secretOwners:  []string{"0x1111111111111111111111111111111111111111"},
			shouldReject:  true,
		},
		{
			name:          "second secret has mismatched owner",
			workflowOwner: "0xABCDef1234567890abcdef1234567890abcdef12",
			secretOwners: []string{
				"0xABCDef1234567890abcdef1234567890abcdef12",
				"0x1111111111111111111111111111111111111111",
			},
			shouldReject: true,
		},
		{
			name:          "matching with different casing",
			workflowOwner: "0xABCDEF1234567890ABCDEF1234567890ABCDEF12",
			secretOwners:  []string{"0xabcdef1234567890abcdef1234567890abcdef12"},
			shouldReject:  false,
		},
		{
			name:          "matching with 0x prefix vs without",
			workflowOwner: "0xabcdef1234567890abcdef1234567890abcdef12",
			secretOwners:  []string{"abcdef1234567890abcdef1234567890abcdef12"},
			shouldReject:  false,
		},
		{
			name:          "matching without 0x prefix vs with",
			workflowOwner: "abcdef1234567890abcdef1234567890abcdef12",
			secretOwners:  []string{"0xabcdef1234567890abcdef1234567890abcdef12"},
			shouldReject:  false,
		},
		{
			name:          "matching with mixed casing and prefix difference",
			workflowOwner: "0xAbCdEf1234567890AbCdEf1234567890AbCdEf12",
			secretOwners:  []string{"abcdef1234567890abcdef1234567890abcdef12"},
			shouldReject:  false,
		},
		{
			name:          "both without prefix and same case",
			workflowOwner: "abcdef1234567890abcdef1234567890abcdef12",
			secretOwners:  []string{"abcdef1234567890abcdef1234567890abcdef12"},
			shouldReject:  false,
		},
	}

	expectedResponse := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: &vault.SecretIdentifier{
					Key:       "Foo",
					Namespace: "Bar",
					Owner:     "placeholder",
				},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
					},
				},
			},
		},
	}
	responseData, err := proto.Marshal(expectedResponse)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			clock := clockwork.NewFakeClock()
			expiry := 10 * time.Second
			store := requests.NewStore[*vaulttypes.Request]()
			handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
			reg := coreCapabilities.NewRegistry(lggr)
			lf := limits.Factory{Settings: cresettings.DefaultGetter}
			capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, nil, lf)
			require.NoError(t, err)
			servicetest.Run(t, capability)

			requestID := fmt.Sprintf("%s::%s::%s", "wf-id", "exec-id", "ref-id")

			reqs := []*vault.SecretRequest{}
			for _, s := range tc.secretOwners {
				reqs = append(reqs, &vault.SecretRequest{
					Id: &vault.SecretIdentifier{
						Key:       "Foo",
						Namespace: "Bar",
						Owner:     s,
					},
					EncryptionKeys: []string{"key"},
				})
			}

			gsr := &vault.GetSecretsRequest{
				WorkflowOwner: tc.workflowOwner,
				Requests: reqs,
			}
			anyproto, err := anypb.New(gsr)
			require.NoError(t, err)

			if !tc.shouldReject {
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						select {
						case <-t.Context().Done():
							return
						default:
							reqs := store.GetByIDs([]string{requestID})
							if len(reqs) == 1 {
								reqs[0].SendResponse(t.Context(), &vaulttypes.Response{
									ID:      requestID,
									Payload: responseData,
								})
								return
							}
						}
					}
				}()
				defer wg.Wait()
			}

			_, err = capability.Execute(t.Context(), capabilities.CapabilityRequest{
				Payload: anyproto,
				Method:  vault.MethodGetSecrets,
				Metadata: capabilities.RequestMetadata{
					WorkflowOwner:       tc.workflowOwner,
					WorkflowID:          "wf-id",
					WorkflowExecutionID: "exec-id",
					ReferenceID:         "ref-id",
				},
			})

			if tc.shouldReject {
				require.ErrorContains(t, err, "secret identifier owner")
				require.ErrorContains(t, err, "does not match workflow owner")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCapability_CapabilityCall_UsesMetadataWorkflowOwner(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := newVaultOrgIDAsSecretOwnerLimitsFactory(t, true)
	resolver := &testOrgResolver{orgID: "org-123"}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	workflowOwner := "0xABCDef1234567890abcdef1234567890abcdef12"
	requestID := "wf-id::exec-id::ref-id"
	gsr := &vault.GetSecretsRequest{
		OrgId:         "",
		WorkflowOwner: workflowOwner,
		Requests: []*vault.SecretRequest{
			{
				Id: &vault.SecretIdentifier{
					Key:       "Foo",
					Namespace: "Bar",
					Owner:     workflowOwner,
				},
				EncryptionKeys: []string{"key"},
			},
		},
	}

	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	expectedResponse := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: &vault.SecretIdentifier{Key: "Foo", Namespace: "Bar", Owner: workflowOwner},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{EncryptedValue: "encrypted-value"},
				},
			},
		},
	}
	data, err := proto.Marshal(expectedResponse)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) == 1 {
					reqs[0].SendResponse(t.Context(), &vaulttypes.Response{ID: requestID, Payload: data})
					return
				}
			}
		}
	}()

	_, err = capability.Execute(t.Context(), capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       workflowOwner,
			WorkflowID:          "wf-id",
			WorkflowExecutionID: "exec-id",
			ReferenceID:         "ref-id",
		},
	})
	require.NoError(t, err)
	wg.Wait()
	assert.Empty(t, resolver.calledWith)
}

func TestCapability_CapabilityCall_ForwardsRequestGetSecretsIdentity(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := newVaultOrgIDAsSecretOwnerLimitsFactory(t, true)
	resolver := &testOrgResolver{orgID: "org-123"}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	requestID := "wf-id::exec-id::ref-id"
	gsr := &vault.GetSecretsRequest{
		OrgId:         "org-123",
		WorkflowOwner: "0xABCDef1234567890abcdef1234567890abcdef12",
		Requests: []*vault.SecretRequest{
			{
				Id: &vault.SecretIdentifier{
					Key:       "Foo",
					Namespace: "Bar",
					Owner:     "0xABCDef1234567890abcdef1234567890abcdef12",
				},
				EncryptionKeys: []string{"key"},
			},
		},
	}

	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)
	responsePayload, err := proto.Marshal(&vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: &vault.SecretIdentifier{
					Key:       "Foo",
					Namespace: "Bar",
					Owner:     "owner",
				},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{EncryptedValue: "encrypted-value"},
				},
			},
		},
	})
	require.NoError(t, err)

	var (
		wg          sync.WaitGroup
		forward     *vault.GetSecretsRequest
		forwardedOK bool
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) != 1 {
					continue
				}
				payload, ok := reqs[0].Payload.(*vault.GetSecretsRequest)
				if !ok {
					return
				}
				cloned, ok := proto.Clone(payload).(*vault.GetSecretsRequest)
				if !ok {
					return
				}
				forward = cloned
				forwardedOK = true
				reqs[0].SendResponse(t.Context(), &vaulttypes.Response{ID: requestID, Payload: responsePayload})
				return
			}
		}
	}()

	_, err = capability.Execute(t.Context(), capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       "0xABCDef1234567890abcdef1234567890abcdef12",
			WorkflowID:          "wf-id",
			WorkflowExecutionID: "exec-id",
			ReferenceID:         "ref-id",
		},
	})
	require.NoError(t, err)
	wg.Wait()
	require.True(t, forwardedOK)
	require.NotNil(t, forward)
	assert.Equal(t, "org-123", forward.OrgId)
	assert.Equal(t, "0xABCDef1234567890abcdef1234567890abcdef12", forward.WorkflowOwner)
	assert.Empty(t, resolver.calledWith)
}

func TestCapability_CapabilityCall_ReturnsIncorrectType(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := limits.Factory{Settings: cresettings.DefaultGetter}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, nil, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	owner := "test-owner"
	workflowID := "test-workflow-id"
	workflowExecutionID := "test-workflow-execution-id"
	referenceID := "test-reference-id"

	requestID := fmt.Sprintf("%s::%s::%s", workflowID, workflowExecutionID, referenceID)

	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}

	gsr := &vault.GetSecretsRequest{
		WorkflowOwner: owner,
		Requests: []*vault.SecretRequest{
			{
				Id:             sid,
				EncryptionKeys: []string{"key"},
			},
		},
	}

	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) == 1 {
					req := reqs[0]
					req.SendResponse(t.Context(), &vaulttypes.Response{
						ID:      requestID,
						Payload: []byte("invalid data"),
					})
					return
				}
			}
		}
	}()

	_, err = capability.Execute(t.Context(), capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          workflowID,
			WorkflowExecutionID: workflowExecutionID,
			ReferenceID:         referenceID,
		},
	})

	wg.Wait()
	assert.ErrorContains(t, err, "cannot parse invalid wire-format data")
}

func TestCapability_CapabilityCall_TimeOut(t *testing.T) {
	lggr := logger.TestLogger(t)
	fakeClock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, fakeClock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := limits.Factory{Settings: cresettings.DefaultGetter}
	capability, err := NewCapability(lggr, fakeClock, expiry, handler, reg, nil, nil, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	owner := "test-owner"
	workflowID := "test-workflow-id"
	workflowExecutionID := "test-workflow-execution-id"
	referenceID := "test-reference-id"

	requestID := fmt.Sprintf("%s::%s::%s", workflowID, workflowExecutionID, referenceID)

	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}

	gsr := &vault.GetSecretsRequest{
		WorkflowOwner: owner,
		Requests: []*vault.SecretRequest{
			{
				Id:             sid,
				EncryptionKeys: []string{"key"},
			},
		},
	}

	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) == 1 {
					fakeClock.Advance(1 * time.Hour)
					return
				}
			}
		}
	}()

	_, err = capability.Execute(t.Context(), capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          workflowID,
			WorkflowExecutionID: workflowExecutionID,
			ReferenceID:         referenceID,
		},
	})

	wg.Wait()
	assert.ErrorContains(t, err, "timeout exceeded")
}

func TestCapability_CRUD(t *testing.T) {
	owner := "0x0001020304050607080900010203040506070809"
	requestID := owner + "::" + "test-request-id"
	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}
	lpk := NewLazyPublicKey()
	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	lpk.Set(pk)
	rawSecret := "raw secret string"
	ownerAddr := common.HexToAddress(owner) // canonical 20-byte address
	var label [32]byte
	copy(label[12:], ownerAddr.Bytes()) // left-pad with 12 zero bytes
	cipher, err := tdh2easy.EncryptWithLabel(pk, []byte(rawSecret), label)
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)
	encryptedSecret := hex.EncodeToString(cipherBytes)

	testCases := []struct {
		name     string
		error    string
		response *vaulttypes.Response
		call     func(t *testing.T, capability *Capability) (*vaulttypes.Response, error)
	}{
		{
			name: "CreateSecrets",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.CreateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.CreateSecrets(t.Context(), req)
			},
		},
		{
			name:     "CreateSecrets_Missing_Key",
			response: nil,
			error:    "secret ID must have key, namespace and owner set",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.CreateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "",
								Namespace: "Bar",
								Owner:     owner,
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.CreateSecrets(t.Context(), req)
			},
		},
		{
			name:     "CreateSecrets_Missing_Namespace",
			response: nil,
			error:    "secret ID must have key, namespace and owner set",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.CreateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "a",
								Namespace: "",
								Owner:     owner,
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.CreateSecrets(t.Context(), req)
			},
		},
		{
			name:     "CreateSecrets_Missing_Owner",
			response: nil,
			error:    "secret ID must have key, namespace and owner set",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.CreateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "a",
								Namespace: "Bar",
								Owner:     "",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.CreateSecrets(t.Context(), req)
			},
		},
		{
			name:     "CreateSecrets_Invalid_Owner",
			response: nil,
			error:    "Encrypted Secret at index [0] doesn't have owner as the label.",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.CreateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "a",
								Namespace: "Bar",
								Owner:     "a",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.CreateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_BatchTooBig",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "request batch size exceeds maximum of 10",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_EmptyRequestID",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "request ID must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: "",
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id:             sid,
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_Missing_Key",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "secret ID must have key, namespace and owner set at index",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "",
								Namespace: "Bar",
								Owner:     "a",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_Missing_Namespace",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "secret ID must have key, namespace and owner set at index",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "w",
								Namespace: "",
								Owner:     "a",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_Missing_Owner",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "secret ID must have key, namespace and owner set at index",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "w",
								Namespace: "na",
								Owner:     "",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_Invalid_Owner",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "Encrypted Secret at index [0] doesn't have owner as the label.",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "w",
								Namespace: "na",
								Owner:     "random",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_InvalidEncryptedSecret",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "failed to verify encrypted value",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id:             sid,
							EncryptedValue: "abcd1234",
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_InvalidRequests_DuplicateIDs",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "duplicate secret ID found",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "Foo",
								Namespace: "Bar",
								Owner:     owner,
							},
							EncryptedValue: encryptedSecret,
						},
						{
							Id: &vault.SecretIdentifier{
								Key:       "Foo",
								Namespace: "Bar",
								Owner:     owner,
							},
							EncryptedValue: encryptedSecret,
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name:     "DeleteSecrets_Invalid_BatchTooBig",
			response: nil,
			error:    "request batch size exceeds maximum of 10",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: requestID,
					Ids: []*vault.SecretIdentifier{
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
					},
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name:     "DeleteSecrets_Invalid_RequestIDMissing",
			response: nil,
			error:    "request ID must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: "",
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name: "DeleteSecrets",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: requestID,
					Ids: []*vault.SecretIdentifier{
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
					},
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name:     "DeleteSecrets_Missing_Owner",
			response: nil,
			error:    "secret ID must have key, namespace and owner set at index 0",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: requestID,
					Ids: []*vault.SecretIdentifier{
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     "",
						},
					},
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name:     "DeleteSecrets_Missing_Namespace",
			response: nil,
			error:    "secret ID must have key, namespace and owner set at index 0",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: requestID,
					Ids: []*vault.SecretIdentifier{
						{
							Key:       "Foo",
							Namespace: "",
							Owner:     "random",
						},
					},
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name:     "DeleteSecrets_Missing_Key",
			response: nil,
			error:    "secret ID must have key, namespace and owner set at index 0",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: requestID,
					Ids: []*vault.SecretIdentifier{
						{
							Key:       "",
							Namespace: "namespace",
							Owner:     "random",
						},
					},
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name:  "DeleteSecrets_Invalid_Duplicates",
			error: "duplicate secret ID found",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.DeleteSecretsRequest{
					RequestId: requestID,
					Ids: []*vault.SecretIdentifier{
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
						{
							Key:       "Foo",
							Namespace: "Bar",
							Owner:     owner,
						},
					},
				}
				return capability.DeleteSecrets(t.Context(), req)
			},
		},
		{
			name:     "ListSecretIdentifiers_OwnerMissing",
			response: nil,
			error:    "requestID, owner or namespace must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: requestID,
					Owner:     "",
				}
				return capability.ListSecretIdentifiers(t.Context(), req)
			},
		},
		{
			name:     "ListSecretIdentifiers_RequestID_Missing",
			response: nil,
			error:    "requestID, owner or namespace must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: "",
					Owner:     "owner",
					Namespace: "namespace",
				}
				return capability.ListSecretIdentifiers(t.Context(), req)
			},
		},
		{
			name:     "ListSecretIdentifiers_Owner_Missing",
			response: nil,
			error:    "requestID, owner or namespace must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: "kk",
					Owner:     "",
					Namespace: "namespace",
				}
				return capability.ListSecretIdentifiers(t.Context(), req)
			},
		},
		{
			name:     "ListSecretIdentifiers_Namespace_Missing",
			response: nil,
			error:    "requestID, owner or namespace must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: "kk",
					Owner:     "owner",
					Namespace: "",
				}
				return capability.ListSecretIdentifiers(t.Context(), req)
			},
		},
		{
			name: "ListSecretIdentifiers",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: requestID,
					Owner:     owner,
					Namespace: "namespace",
				}
				return capability.ListSecretIdentifiers(t.Context(), req)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			clock := clockwork.NewFakeClock()
			expiry := 10 * time.Second
			store := requests.NewStore[*vaulttypes.Request]()
			handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
			reg := coreCapabilities.NewRegistry(lggr)
			lf := limits.Factory{Settings: cresettings.DefaultGetter}
			capability, err := NewCapability(lggr, clock, expiry, handler, reg, lpk, nil, lf)
			require.NoError(t, err)
			servicetest.Run(t, capability)

			wait := func() {}
			if tc.error == "" {
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						select {
						case <-t.Context().Done():
							return
						default:
							reqs := store.GetByIDs([]string{requestID})
							if len(reqs) == 1 {
								req := reqs[0]
								req.SendResponse(t.Context(), tc.response)
								return
							}
						}
					}
				}()
				wait = wg.Wait
			}

			resp, err := tc.call(t, capability)

			if tc.error != "" {
				assert.ErrorContains(t, err, tc.error)
			} else {
				require.NoError(t, err)
				wait()
				assert.Equal(t, tc.response, resp)
			}
		})
	}
}

func TestCapability_Lifecycle(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lf := limits.Factory{Settings: cresettings.DefaultGetter}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, nil, lf)
	require.NoError(t, err)

	_, err = reg.GetExecutable(t.Context(), vault.CapabilityID)
	require.ErrorContains(t, err, "no compatible capability found for id vault@1.0.0")

	require.NoError(t, capability.Start(t.Context()))

	_, err = reg.GetExecutable(t.Context(), vault.CapabilityID)
	require.NoError(t, err)

	require.NoError(t, capability.Close())

	got, err := reg.GetExecutable(t.Context(), vault.CapabilityID)
	require.NoError(t, err)
	loader, ok := got.(interface {
		Load() *capabilities.ExecutableCapability
	})
	require.True(t, ok)
	require.Nil(t, loader.Load())
}

func TestCapability_PublicKeyGet(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lpk := NewLazyPublicKey()
	lf := limits.Factory{Settings: cresettings.DefaultGetter}
	capability, err := NewCapability(lggr, clock, expiry, handler, reg, lpk, nil, lf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	_, err = capability.GetPublicKey(t.Context(), nil)
	require.ErrorContains(t, err, "could not get public key: is the plugin initialized?")

	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	lpk.Set(pk)

	pkb, err := pk.Marshal()
	require.NoError(t, err)

	hpkb := hex.EncodeToString(pkb)

	resp, err := capability.GetPublicKey(t.Context(), nil)
	require.NoError(t, err)

	assert.Equal(t, hpkb, resp.PublicKey)
}
