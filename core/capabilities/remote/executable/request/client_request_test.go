package request_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	pbvalues "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	stepRef1             = "stepRef1"

	testDispatcherChanCap = 100
)

func Test_ClientRequest_MessageValidation(t *testing.T) {
	numWorkflowPeers := 2
	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
	}

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)
	require.NoError(t, err)

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "1000ms",
	})
	require.NoError(t, err)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
			ReferenceID:         stepRef1,
		},
		Inputs: executeInputs,
		Config: transmissionSchedule,
	}

	m, err := values.NewMap(map[string]any{"response": "response1"})
	require.NoError(t, err)
	capabilityResponse := commoncap.CapabilityResponse{
		Value: m,
	}

	rawResponse, err := pb.MarshalCapabilityResponse(capabilityResponse)
	require.NoError(t, err)

	t.Run("Send second message with different response", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 2, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		defer req.Cancel(errors.New("test end"))

		require.NoError(t, err)

		nm, err := values.NewMap(map[string]any{"response": "response2"})
		require.NoError(t, err)
		capabilityResponse2 := commoncap.CapabilityResponse{
			Value: nm,
		}

		rawResponse2, err := pb.MarshalCapabilityResponse(capabilityResponse2)
		require.NoError(t, err)
		msg2 := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse2,
			MessageId:       []byte("messageID"),
		}

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg2.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msg2)
		require.NoError(t, err)

		select {
		case resp := <-req.ResponseChan():
			require.Error(t, resp.Err)
			var capErr caperrors.Error
			require.ErrorAs(t, resp.Err, &capErr)
			require.Equal(t, caperrors.OriginSystem, capErr.Origin())
			require.Equal(t, caperrors.ConsensusFailed, capErr.Code())
			assert.Contains(t, capErr.Error(),
				"[100]ConsensusFailed: response quorum unreachable: not enough matching capability responses: received 2/2 peer responses with 2 unique payloads; best match count 1, need 2 (0 responses pending)")

		default:
			t.Fatal("expected early quorum unreachable response")
		}
	})

	t.Run("Send second message from non calling Don peer", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 2, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		nonDonPeer := NewP2PPeerID(t)
		msg.Sender = nonDonPeer[:]
		err = req.OnMessage(ctx, msg)
		require.Error(t, err)

		select {
		case <-req.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message from same peer as first message", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 2, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)
		err = req.OnMessage(ctx, msg)
		require.Error(t, err)

		select {
		case <-req.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message with same error as first", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		drainInitialPeerSends(t, dispatcher, len(capabilityPeers))

		msgWithError := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        assert.AnError.Error(),
		}

		msgWithError.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		msgWithError.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		response := <-req.ResponseChan()

		assert.Equal(t, fmt.Sprintf("%s : %s", types.Error_INTERNAL_ERROR, assert.AnError.Error()), response.Err.Error())

		var capErr caperrors.Error
		require.ErrorAs(t, response.Err, &capErr)
		assert.Equal(t, caperrors.OriginSystem, capErr.Origin(), "non-serialized ErrorMsg falls back to private system capability error")
		assert.Equal(t, caperrors.VisibilityPrivate, capErr.Visibility())
		assert.Equal(t, caperrors.Unknown, capErr.Code())
	})

	t.Run("Error response with serialized caperrors unwraps correctly as usererror", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		drainInitialPeerSends(t, dispatcher, len(capabilityPeers))

		serialized := caperrors.NewPublicUserError(errors.New("rpc error: EVM error invalid argument"), caperrors.FailedPrecondition).SerializeToRemoteString()
		msgWithError := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        serialized,
		}

		msgWithError.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		msgWithError.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		response := <-req.ResponseChan()

		wantDisplay := fmt.Sprintf("%s : %s", types.Error_INTERNAL_ERROR, serialized)
		assert.Equal(t, wantDisplay, response.Err.Error(), "It should be equal to 'Public:User:FailedPrecondition:rpc error: EVM error invalid argument'")

		var capErr caperrors.Error
		require.ErrorAs(t, response.Err, &capErr)
		assert.Equal(t, caperrors.OriginUser, capErr.Origin())
		assert.Equal(t, caperrors.VisibilityPublic, capErr.Visibility())
		assert.Equal(t, caperrors.FailedPrecondition, capErr.Code())
	})

	t.Run("Send three messages with different errors", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		drainInitialPeerSends(t, dispatcher, len(capabilityPeers))

		msgWithError := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error",
			Sender:          capabilityPeers[0][:],
		}

		msgWithError2 := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error2",
			Sender:          capabilityPeers[1][:],
		}

		msgWithError3 := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error3",
			Sender:          capabilityPeers[2][:],
		}

		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)
		err = req.OnMessage(ctx, msgWithError2)
		require.NoError(t, err)
		err = req.OnMessage(ctx, msgWithError3)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		assert.Equal(t, "received 3 errors, last error INTERNAL_ERROR : an error3", response.Err.Error())
	})

	t.Run("Execute Request", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		drainInitialPeerSends(t, dispatcher, len(capabilityPeers))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))
	})
	t.Run("Execute Request With Valid Attestation", func(t *testing.T) {
		const F = 1
		const N = 3*F + 1
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, N, F)

		configDigest := ocrtypes.ConfigDigest{1, 2, 3, 4, 5}
		kb1, err := ocr2key.New(corekeys.EVM)
		require.NoError(t, err)
		kb2, err := ocr2key.New(corekeys.EVM)
		require.NoError(t, err)

		seqNr := uint64(100)

		payload, err := values.NewMap(map[string]int{
			"number": 42,
		})
		require.NoError(t, err)
		payloadAsAny, err := anypb.New(values.Proto(payload))
		require.NoError(t, err)

		spendUnit, spendValue := "testunit", "42"
		capResponse := commoncap.CapabilityResponse{
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{
					{SpendUnit: spendUnit, SpendValue: spendValue},
				},
			},
			Payload: payloadAsAny,
		}

		reportData, err := commoncap.ResponseToReportData(capabilityRequest.Metadata.WorkflowExecutionID, capabilityRequest.Metadata.ReferenceID, payloadAsAny.Value, capResponse.Metadata)
		require.NoError(t, err)

		sig1, err := kb1.Sign3(configDigest, seqNr, reportData[:])
		require.NoError(t, err)
		sig2, err := kb2.Sign3(configDigest, seqNr, reportData[:])
		require.NoError(t, err)

		capResponse.OCRAttestation = &commoncap.OCRAttestation{
			ConfigDigest:   configDigest,
			SequenceNumber: seqNr,
			Sigs: []commoncap.AttributedSignature{
				{Signer: 0, Signature: sig1},
				{Signer: 1, Signature: sig2},
			},
		}

		rawResponseWithAttestation, err := pb.MarshalCapabilityResponse(capResponse)
		require.NoError(t, err)

		ocrSigners := [][]byte{kb1.PublicKey(), kb2.PublicKey()}

		assertValidResponse := func(t *testing.T, result []byte) {
			capResponse, err := pb.UnmarshalCapabilityResponse(result)
			require.NoError(t, err)

			var pbValue pbvalues.Value
			require.NoError(t, capResponse.Payload.UnmarshalTo(&pbValue))
			receivedValue, err := values.FromProto(&pbValue)
			require.NoError(t, err)

			var receivedMap map[string]int
			require.NoError(t, receivedValue.UnwrapTo(&receivedMap))

			assert.Equal(t, 42, receivedMap["number"])
			require.GreaterOrEqual(t, len(capResponse.Metadata.Metering), 1)
			require.Equal(t, spendUnit, capResponse.Metadata.Metering[0].SpendUnit)
			require.Equal(t, spendValue, capResponse.Metadata.Metering[0].SpendValue)
		}

		t.Run("succeeds on first peer with valid attestation", func(t *testing.T) {
			ctx := t.Context()

			dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
			req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
				workflowDonInfo, dispatcher, 10*time.Minute, nil, "", ocrSigners)
			require.NoError(t, err)
			defer req.Cancel(errors.New("test end"))

			for range N {
				<-dispatcher.msgs
			}

			assert.Empty(t, dispatcher.msgs)

			msg := &types.MessageBody{
				CapabilityId:    capInfo.ID,
				CapabilityDonId: capDonInfo.ID,
				CallerDonId:     workflowDonInfo.ID,
				Method:          types.MethodExecute,
				Payload:         rawResponseWithAttestation,
				MessageId:       []byte("messageID"),
			}
			msg.Sender = capabilityPeers[0][:]
			err = req.OnMessage(ctx, msg)
			require.NoError(t, err)

			response := <-req.ResponseChan()
			assertValidResponse(t, response.Result)
		})
		t.Run("attestation is not valid, but we fallback to identical responses", func(t *testing.T) {
			ctx := t.Context()

			dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
			req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
				workflowDonInfo, dispatcher, 10*time.Minute, nil, "", ocrSigners)
			require.NoError(t, err)
			defer req.Cancel(errors.New("test end"))

			for range N {
				<-dispatcher.msgs
			}

			assert.Empty(t, dispatcher.msgs)

			for i := range F + 1 {
				respInvalidAtt := commoncap.CapabilityResponse{
					Metadata: commoncap.ResponseMetadata{
						Metering: []commoncap.MeteringNodeDetail{
							{SpendUnit: spendUnit, SpendValue: spendValue},
						},
					},
					OCRAttestation: &commoncap.OCRAttestation{
						ConfigDigest: configDigest,
						// make the sequence number invalid
						SequenceNumber: seqNr + uint64(i) + 1, // #nosec G115 -- i is non-negative and within uint64 range
						Sigs: []commoncap.AttributedSignature{
							{Signer: 0, Signature: sig1},
							{Signer: 1, Signature: sig2},
						},
					},
					Payload: payloadAsAny,
				}

				rawRespInvalidAtt, err := pb.MarshalCapabilityResponse(respInvalidAtt)
				require.NoError(t, err)

				msg := &types.MessageBody{
					CapabilityId:    capInfo.ID,
					CapabilityDonId: capDonInfo.ID,
					CallerDonId:     workflowDonInfo.ID,
					Method:          types.MethodExecute,
					Payload:         rawRespInvalidAtt,
					MessageId:       []byte("messageID"),
				}
				msg.Sender = capabilityPeers[i][:]
				err = req.OnMessage(ctx, msg)
				require.NoError(t, err)
			}

			response := <-req.ResponseChan()
			assertValidResponse(t, response.Result)
		})

		t.Run("2F peers return ErrResponsePayloadNotAvailable then success", func(t *testing.T) {
			ctx := t.Context()
			dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
			req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
				workflowDonInfo, dispatcher, 10*time.Minute, nil, "", ocrSigners)
			require.NoError(t, err)
			defer req.Cancel(errors.New("test end"))

			for range N {
				<-dispatcher.msgs
			}

			assert.Empty(t, dispatcher.msgs)

			for i := range 2 * F {
				msgNA := &types.MessageBody{
					CapabilityId:    capInfo.ID,
					CapabilityDonId: capDonInfo.ID,
					CallerDonId:     workflowDonInfo.ID,
					Method:          types.MethodExecute,
					MessageId:       []byte("messageID"),
					Error:           types.Error_INTERNAL_ERROR,
					ErrorMsg:        commoncap.ErrResponsePayloadNotAvailable.Error(),
				}
				msgNA.Sender = capabilityPeers[i][:]
				require.NoError(t, req.OnMessage(ctx, msgNA))
			}

			msgOK := &types.MessageBody{
				CapabilityId:    capInfo.ID,
				CapabilityDonId: capDonInfo.ID,
				CallerDonId:     workflowDonInfo.ID,
				Method:          types.MethodExecute,
				Payload:         rawResponseWithAttestation,
				MessageId:       []byte("messageID"),
			}
			msgOK.Sender = capabilityPeers[2*F][:]
			require.NoError(t, req.OnMessage(ctx, msgOK))

			response := <-req.ResponseChan()
			assertValidResponse(t, response.Result)
		})

		t.Run("2F+1 peers return ErrResponsePayloadNotAvailable", func(t *testing.T) {
			ctx := t.Context()
			dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
			req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
				workflowDonInfo, dispatcher, 10*time.Minute, nil, "", ocrSigners)
			require.NoError(t, err)
			defer req.Cancel(errors.New("test end"))

			for range N {
				<-dispatcher.msgs
			}

			assert.Empty(t, dispatcher.msgs)

			noPayloadMsg := types.MessageBody{
				CapabilityId:    capInfo.ID,
				CapabilityDonId: capDonInfo.ID,
				CallerDonId:     workflowDonInfo.ID,
				Method:          types.MethodExecute,
				MessageId:       []byte("messageID"),
				Error:           types.Error_INTERNAL_ERROR,
				ErrorMsg:        commoncap.ErrResponsePayloadNotAvailable.Error(),
			}

			for i := range 2 * F {
				noPayloadMsg.Sender = capabilityPeers[i][:]
				require.NoError(t, req.OnMessage(ctx, &noPayloadMsg))
			}

			noPayloadMsg.Sender = capabilityPeers[2*F][:]
			require.Error(t, req.OnMessage(ctx, &noPayloadMsg))
		})
	})

	t.Run("Executes full schedule", func(t *testing.T) {
		beholderTester := beholdertest.NewObserver(t)
		lggr, obs := logger.TestObserved(t, zapcore.DebugLevel)

		capPeers, capDonInfo, capInfo := capabilityDon(t, 3, 1)

		ctx := t.Context()
		ctxWithCancel, cancelFn := context.WithCancel(t.Context())

		// cancel the context immediately so we can verify
		// that the schedule is still executed entirely.
		cancelFn()

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(
			ctxWithCancel,
			lggr,
			capabilityRequest,
			capInfo,
			workflowDonInfo,
			dispatcher,
			10*time.Minute,
			nil,
			"",
			nil,
		)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		// Despite the context being cancelled,
		// we still send the full schedule.
		drainInitialPeerSends(t, dispatcher, len(capPeers))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capPeers[1][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))

		logs := obs.FilterMessage("sending request to peers").All()
		assert.Len(t, logs, 1)

		log := logs[0]
		for _, k := range log.Context {
			if k.Key == "originalTimeout" {
				assert.Equal(t, int64(0), k.Integer)
			}

			if k.Key == "effectiveTimeout" {
				assert.Equal(t, k.Integer, int64(10*time.Second))
			}
		}

		// Verify the TransmissionsScheduledEvent data
		assert.Equal(t, 1, beholderTester.Len(t, "beholder_entity", fmt.Sprintf("%v.%v", request.TransmissionEventProtoPkg, request.TransmissionEventEntity)))

		// Get the messages for the transmission event
		messages := beholderTester.Messages(t, "beholder_entity", fmt.Sprintf("%v.%v", request.TransmissionEventProtoPkg, request.TransmissionEventEntity))
		assert.Len(t, messages, 1)

		// Unmarshal the message to verify its contents
		var event events.TransmissionsScheduledEvent
		err = proto.Unmarshal(messages[0].Body, &event)
		require.NoError(t, err)

		// Verify the event fields
		assert.Equal(t, transmission.Schedule_AllAtOnce, event.ScheduleType)
		assert.Equal(t, workflowExecutionID1, event.WorkflowExecutionID)
		assert.Equal(t, "cap_id@1.0.0", event.CapabilityID)
		assert.Equal(t, stepRef1, event.StepRef)
		assert.Equal(t, fmt.Sprintf("Execute:%v:%v", workflowExecutionID1, stepRef1), event.TransmissionID)
		assert.NotEmpty(t, event.Timestamp)

		// Verify the peer delays
		assert.Len(t, event.PeerTransmissionDelays, 3)

		// Convert map to slice of delays and sort them
		var delays []int64
		for _, delay := range event.PeerTransmissionDelays {
			delays = append(delays, delay)
		}
		slices.Sort(delays)

		// Verify delays are sorted and increment by 1000ms
		for i := 1; i < len(delays); i++ {
			assert.Equal(t, delays[i-1], delays[i], "delays should be the same")
		}

		// Verify each peer ID exists in capability peers
		for peerID := range event.PeerTransmissionDelays {
			found := false
			for _, peer := range capPeers {
				if peer.String() == peerID {
					found = true
					break
				}
			}
			assert.True(t, found, "peer ID %s not found in capability peers", peerID)
		}
	})

	t.Run("Uses passed in time out if larger than schedule", func(t *testing.T) {
		lggr, obs := logger.TestObserved(t, zapcore.DebugLevel)

		capPeers, capDonInfo, capInfo := capabilityDon(t, 3, 1)

		ctx := t.Context()
		ctx, cancelFn := context.WithTimeout(ctx, 15*time.Second)
		defer cancelFn()

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(
			ctx,
			lggr,
			capabilityRequest,
			capInfo,
			workflowDonInfo,
			dispatcher,
			10*time.Minute,
			nil,
			"",
			nil,
		)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		// Despite the context being cancelled,
		// we still send the full schedule.
		drainInitialPeerSends(t, dispatcher, len(capPeers))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capPeers[1][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))

		logs := obs.FilterMessage("sending request to peers").All()
		assert.Len(t, logs, 1)

		log := logs[0]
		for _, k := range log.Context {
			if k.Key == "effectiveTimeout" {
				// Greater than what it would otherwise be
				// i.e. 2 *deltaStage + margin = 12s
				assert.Greater(t, k.Integer, int64(12*time.Second))
			}
		}
	})

	// tests that (once added) metering data in the capability responses
	// will not cause the identical response calculation to break;
	// also locks in no validation of SpendUnit/SpendValue at that layer.
	t.Run("with metering metadata", func(t *testing.T) {
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		capabilityResponseWithMetering1 := commoncap.CapabilityResponse{
			Value: m,
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{
					{SpendUnit: "testunit_a", SpendValue: "15"},
				},
			},
		}

		capabilityResponseWithMetering2 := commoncap.CapabilityResponse{
			Value: m,
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{
					{SpendUnit: "testunit_b", SpendValue: "17"},
				},
			},
		}

		payload1, err2 := pb.MarshalCapabilityResponse(capabilityResponseWithMetering1)
		require.NoError(t, err2)

		payload2, err2 := pb.MarshalCapabilityResponse(capabilityResponseWithMetering2)
		require.NoError(t, err2)

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Payload = payload1

		ctx := t.Context()

		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(ctx, logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capabilityPeers[1][:]
		msg.Payload = payload2
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]
		assert.Equal(t, resp, values.NewString("response1"))

		assert.Len(t, capResponse.Metadata.Metering, 2)
		spendUnit := capResponse.Metadata.Metering[0].SpendUnit
		spendValue := capResponse.Metadata.Metering[0].SpendValue
		p2pID := capResponse.Metadata.Metering[0].Peer2PeerID

		assert.Equal(t, "testunit_a", spendUnit)
		assert.Equal(t, "15", spendValue)
		assert.Equal(t, capabilityPeers[0].String(), p2pID)

		spendUnit = capResponse.Metadata.Metering[1].SpendUnit
		spendValue = capResponse.Metadata.Metering[1].SpendValue
		p2pID = capResponse.Metadata.Metering[1].Peer2PeerID

		assert.Equal(t, "testunit_b", spendUnit)
		assert.Equal(t, "17", spendValue)
		assert.Equal(t, capabilityPeers[1].String(), p2pID)
	})

	capabilityRequestV2 := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
			ReferenceID:         stepRef1,
		},
		// No Inputs or Config, including transmission schedule
	}

	t.Run("Executes full schedule for a V2 request", func(t *testing.T) {
		beholderTester := beholdertest.NewObserver(t)
		lggr, obs := logger.TestObserved(t, zapcore.DebugLevel)
		capPeers, capDonInfo, capInfo := capabilityDon(t, 3, 1)
		dispatcher := newClientRequestTestDispatcher()
		req, err := request.NewClientExecuteRequest(
			t.Context(),
			lggr,
			capabilityRequestV2,
			capInfo,
			workflowDonInfo,
			dispatcher,
			10*time.Minute,
			&transmission.TransmissionConfig{
				Schedule:   transmission.Schedule_OneAtATime,
				DeltaStage: 1000 * time.Millisecond,
			},
			"",
			nil,
		)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		drainInitialPeerSends(t, dispatcher, len(capPeers))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capPeers[0][:]
		require.NoError(t, req.OnMessage(t.Context(), msg))
		msg.Sender = capPeers[1][:]
		require.NoError(t, req.OnMessage(t.Context(), msg))

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]
		assert.Equal(t, resp, values.NewString("response1"))
		assert.Len(t, obs.FilterMessage("sending request to peers").All(), 1)

		// Verify the TransmissionsScheduledEvent data
		assert.Equal(t, 1, beholderTester.Len(t, "beholder_entity", fmt.Sprintf("%v.%v", request.TransmissionEventProtoPkg, request.TransmissionEventEntity)))

		// Get the messages for the transmission event
		messages := beholderTester.Messages(t, "beholder_entity", fmt.Sprintf("%v.%v", request.TransmissionEventProtoPkg, request.TransmissionEventEntity))
		assert.Len(t, messages, 1)

		// Unmarshal the message to verify its contents
		var event events.TransmissionsScheduledEvent
		err = proto.Unmarshal(messages[0].Body, &event)
		require.NoError(t, err)

		// Verify the event fields
		assert.Equal(t, transmission.Schedule_AllAtOnce, event.ScheduleType)
		assert.Equal(t, workflowExecutionID1, event.WorkflowExecutionID)
		assert.Equal(t, "cap_id@1.0.0", event.CapabilityID)
		assert.Equal(t, stepRef1, event.StepRef)
		assert.Equal(t, fmt.Sprintf("Execute:%v:%v", workflowExecutionID1, stepRef1), event.TransmissionID)
		assert.NotEmpty(t, event.Timestamp)

		// Verify the peer delays
		assert.Len(t, event.PeerTransmissionDelays, 3)

		// Convert map to slice of delays and sort them
		var delays []int64
		for _, delay := range event.PeerTransmissionDelays {
			delays = append(delays, delay)
		}
		slices.Sort(delays)

		// Verify delays are sorted and increment by 1000ms
		for i := 1; i < len(delays); i++ {
			assert.Equal(t, delays[i-1], delays[i], "v2 capabilities should be all at once")
		}
	})
}

func newClientRequestTestDispatcher() *clientRequestTestDispatcher {
	return &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, testDispatcherChanCap)}
}

func drainInitialPeerSends(t *testing.T, d *clientRequestTestDispatcher, numCapabilityPeers int) {
	t.Helper()
	require.Eventually(t, func() bool {
		return len(d.msgs) == numCapabilityPeers
	}, 2*time.Second, time.Millisecond, "timed out waiting for %d buffered outbound messages", numCapabilityPeers)
	require.Len(t, d.msgs, numCapabilityPeers, "dispatcher outbound buffer before draining initial peer sends")
	for range numCapabilityPeers {
		<-d.msgs
	}
	require.Empty(t, d.msgs)
}

func capabilityDon(t *testing.T, numCapabilityPeers int, f uint8) ([]p2ptypes.PeerID, commoncap.DON, commoncap.CapabilityInfo) {
	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       f,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}
	return capabilityPeers, capDonInfo, capInfo
}

type clientRequestTestDispatcher struct {
	msgs chan *types.MessageBody
}

func (t *clientRequestTestDispatcher) Name() string {
	return "clientRequestTestDispatcher"
}

func (t *clientRequestTestDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *clientRequestTestDispatcher) Close() error {
	return nil
}

func (t *clientRequestTestDispatcher) Ready() error {
	return nil
}

func (t *clientRequestTestDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *clientRequestTestDispatcher) SetReceiver(capabilityID string, donID uint32, receiver types.Receiver) error {
	return nil
}

func (t *clientRequestTestDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (t *clientRequestTestDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, methodName string, receiver types.Receiver) error {
	return nil
}

func (t *clientRequestTestDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, methodName string) {
}

func (t *clientRequestTestDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	t.msgs <- msgBody
	return nil
}
