package vault

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCapability_CapabilityCall(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	requestAuthorizer := vaultcapmocks.NewRequestAuthorizer(t)
	capability := NewCapability(lggr, clock, expiry, handler, requestAuthorizer)
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
	requestAuthorizer := vaultcapmocks.NewRequestAuthorizer(t)
	capability := NewCapability(lggr, clock, expiry, handler, requestAuthorizer)
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

func TestCapability_CapabilityCall_ReturnsIncorrectType(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	requestAuthorizer := vaultcapmocks.NewRequestAuthorizer(t)
	capability := NewCapability(lggr, clock, expiry, handler, requestAuthorizer)
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
	requestAuthorizer := vaultcapmocks.NewRequestAuthorizer(t)
	capability := NewCapability(lggr, fakeClock, expiry, handler, requestAuthorizer)
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
	owner := "test-owner"
	requestID := owner + "::" + "test-request-id"
	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}

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
							EncryptedValue: "encrypted-value",
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
							EncryptedValue: "encrypted-value",
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
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
						},
						{
							Id:             sid,
							EncryptedValue: "encrypted-value",
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
							EncryptedValue: "encrypted-value",
						},
					},
				}
				return capability.UpdateSecrets(t.Context(), req)
			},
		},
		{
			name: "UpdateSecrets_InvalidSecretID",
			response: &vaulttypes.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "secret ID must have both key and owner set",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.UpdateSecretsRequest{
					RequestId: requestID,
					EncryptedSecrets: []*vault.EncryptedSecret{
						{
							Id: &vault.SecretIdentifier{
								Key:       "",
								Namespace: "Bar",
								Owner:     "",
							},
							EncryptedValue: "encrypted-value",
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
								Owner:     "Owner",
							},
							EncryptedValue: "encrypted-value",
						},
						{
							Id: &vault.SecretIdentifier{
								Key:       "Foo",
								Namespace: "Bar",
								Owner:     "Owner",
							},
							EncryptedValue: "encrypted-value",
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
			name:     "ListSecretIdentifiers_Invalid_OwnerMissing",
			response: nil,
			error:    "owner must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: requestID,
					Owner:     "",
				}
				return capability.ListSecretIdentifiers(t.Context(), req)
			},
		},
		{
			name:     "ListSecretIdentifiers_Invalid_RequestIDMissing",
			response: nil,
			error:    "request ID must not be empty",
			call: func(t *testing.T, capability *Capability) (*vaulttypes.Response, error) {
				req := &vault.ListSecretIdentifiersRequest{
					RequestId: "",
					Owner:     "owner",
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
			requestAuthorizer := vaultcapmocks.NewRequestAuthorizer(t)
			requestAuthorizer.On("AuthorizeRequest", t.Context(), mock.Anything).Return(true, owner, nil).Maybe()
			capability := NewCapability(lggr, clock, expiry, handler, requestAuthorizer)
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
