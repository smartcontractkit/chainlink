package vault

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	vault2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
)

func TestCapability_CapabilityCall(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vault2.Request]()
	handler := requests.NewHandler[*vault2.Request, *vault2.Response](lggr, store, clock, expiry)
	capability := NewCapability(lggr, clock, expiry, handler)
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
					req.SendResponse(t.Context(), &vault2.Response{
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
	store := requests.NewStore[*vault2.Request]()
	handler := requests.NewHandler[*vault2.Request, *vault2.Response](lggr, store, clock, expiry)
	capability := NewCapability(lggr, clock, expiry, handler)
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
					req.SendResponse(t.Context(), &vault2.Response{
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
	store := requests.NewStore[*vault2.Request]()
	handler := requests.NewHandler[*vault2.Request, *vault2.Response](lggr, store, clock, expiry)
	capability := NewCapability(lggr, clock, expiry, handler)
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
					req.SendResponse(t.Context(), &vault2.Response{
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
	store := requests.NewStore[*vault2.Request]()
	handler := requests.NewHandler[*vault2.Request, *vault2.Response](lggr, store, fakeClock, expiry)
	capability := NewCapability(lggr, fakeClock, expiry, handler)
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
	requestID := "test-request-id"
	owner := "test-owner"
	sid := &vault.SecretIdentifier{
		Key:       "Foo",
		Namespace: "Bar",
		Owner:     owner,
	}

	testCases := []struct {
		name     string
		error    string
		response *vault2.Response
		call     func(t *testing.T, capability *Capability) (*vault2.Response, error)
	}{
		{
			name: "CreateSecrets",
			response: &vault2.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			call: func(t *testing.T, capability *Capability) (*vault2.Response, error) {
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
			response: &vault2.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			call: func(t *testing.T, capability *Capability) (*vault2.Response, error) {
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
			response: &vault2.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "request batch size exceeds maximum of 10",
			call: func(t *testing.T, capability *Capability) (*vault2.Response, error) {
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
			response: &vault2.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "request ID must not be empty",
			call: func(t *testing.T, capability *Capability) (*vault2.Response, error) {
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
			response: &vault2.Response{
				ID:      "response-id",
				Payload: []byte("hello world"),
				Format:  "protobuf",
			},
			error: "secret ID must have both key and owner set",
			call: func(t *testing.T, capability *Capability) (*vault2.Response, error) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			clock := clockwork.NewFakeClock()
			expiry := 10 * time.Second
			store := requests.NewStore[*vault2.Request]()
			handler := requests.NewHandler[*vault2.Request, *vault2.Response](lggr, store, clock, expiry)
			capability := NewCapability(lggr, clock, expiry, handler)
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
			wait()

			if tc.error != "" {
				assert.ErrorContains(t, err, tc.error)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.response, resp)
			}
		})
	}
}
