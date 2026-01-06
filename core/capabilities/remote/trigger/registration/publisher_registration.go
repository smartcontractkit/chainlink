package registration

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const messageCacheKey = "publisher_registration"

type PublisherRegistration struct {
	lggr logger.Logger

	wg     sync.WaitGroup
	stopCh services.StopChan

	sendToDispatcher func(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error

	triggerID       string
	workflowID      string
	capabilityID    string
	capabilityDonID uint32

	underlyingTrigger commoncap.TriggerCapability

	eventForwardingChannel               chan commoncap.TriggerResponse
	underlyingTriggerRegistrationRequest commoncap.TriggerRegistrationRequest

	peerRequestsCache *messagecache.MessageCache[string, p2ptypes.PeerID]

	registrationRequests map[p2ptypes.PeerID]uint32

	registrationStatus       types.RegistrationStatus
	registrationErrorMessage string

	createdAt time.Time
}

func NewPublisherRegistration(lggr logger.Logger,
	underlying commoncap.TriggerCapability, // This is passed in and not read from dynamic config, reading from dynamic config introduces the possibility that
// the registration would register on a different underlying trigger than it unregisters from, this is a bug in the legacy code.
// If the underlying trigger changes, then this would take effect on a new registration only, that is to say the existing registration
// would need to be allowed to expire, this is the same behaviour as the legacy code.
	triggerID string,
	workflowID string,
	capabilityID string,
	sendToDispatcher func(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error) *PublisherRegistration {

	return &PublisherRegistration{
		lggr:                   lggr,
		stopCh:                 make(services.StopChan),
		underlyingTrigger:      underlying,
		registrationRequests:   make(map[p2ptypes.PeerID]uint32),
		triggerID:              triggerID,
		workflowID:             workflowID,
		capabilityID:           capabilityID,
		sendToDispatcher:       sendToDispatcher,
		peerRequestsCache:      messagecache.NewMessageCache[string, p2ptypes.PeerID](),
		eventForwardingChannel: make(chan commoncap.TriggerResponse),
		createdAt:              time.Now(),
		registrationStatus:     types.RegistrationStatus_UNREGISTERED,
	}
}

func (pr *PublisherRegistration) EventChannel() <-chan commoncap.TriggerResponse {
	return pr.eventForwardingChannel
}

// IsLive determines whether the publisher registration is considered live based on whether sufficient
// registration requests have been received recently from peers in the caller DON.
func (pr *PublisherRegistration) IsLive(registrationExpiry time.Duration, callerDon commoncap.DON) bool {

	// If it's still unregistered after registrationExpiry from creation time, it is considered to be dead
	if pr.registrationStatus == types.RegistrationStatus_UNREGISTERED {
		if time.Since(pr.createdAt) > registrationExpiry {
			return false
		}
	}

	now := time.Now().UnixMilli()
	ready, _, _ := pr.peerRequestsCache.Ready(messageCacheKey, uint32(2*callerDon.F+1), now-registrationExpiry.Milliseconds(), false)
	return ready
}

func (pr *PublisherRegistration) AddRegistrationRequest(ctx context.Context, sender p2ptypes.PeerID, rawRequest []byte,
	callerDon commoncap.DON, registrationExpiry time.Duration) {
	nowMs := time.Now().UnixMilli()
	// First remote any expired messages from the cache, the cache is not expected to be large so this is acceptable
	// and uses no more resource that the ready call in this case so need to clean-up on a background goroutine.
	pr.peerRequestsCache.DeleteOlderThan(nowMs - registrationExpiry.Milliseconds())

	// Record the request to ensure liveness of the registration

	pr.peerRequestsCache.Insert(messageCacheKey, sender, nowMs, rawRequest)

	// If registration is already successful, send a status update immediately
	if pr.registrationStatus == types.RegistrationStatus_REGISTERED {
		pr.lggr.Debugw("already registered on trigger", "workflowId", pr.workflowID)
		pr.sendTriggerRegistrationStatus(sender, callerDon.ID)
		return
	}

	// Retry registration if not yet successful - in the case where registration fails it means each request from each calling
	// capability peer will result in a re-registration attempt (assuming sufficient requests have been received to aggregate)
	// This is the same behaviour as the legacy code.  TODO Possible we may want limit this.

	// NOTE: require 2F+1 by default, introduce different strategies later (KS-76)
	minRequired := uint32(2*callerDon.F + 1)
	ready, requestingPeers, payloads := pr.peerRequestsCache.Ready(messageCacheKey, minRequired, nowMs-registrationExpiry.Milliseconds(), false)
	if !ready {
		pr.lggr.Debugw("not ready to aggregate yet", "workflowId", pr.workflowID, "minRequired", minRequired)
		return
	}
	aggregated, err := aggregation.AggregateModeRaw(payloads, uint32(callerDon.F+1))
	if err != nil {
		pr.lggr.Errorw("failed to aggregate trigger registrations", "workflowId", pr.workflowID, "err", err)
		return
	}
	unmarshalled, err := pb.UnmarshalTriggerRegistrationRequest(aggregated)
	if err != nil {
		pr.lggr.Errorw("failed to unmarshal request", "err", err)
		return
	}

	callbackCh, err := pr.underlyingTrigger.RegisterTrigger(ctx, unmarshalled)
	registrationStatus := types.RegistrationStatus_REGISTERED
	var errMsg string
	if err != nil {
		registrationStatus = types.RegistrationStatus_REGISTRATION_ERROR
		var capErr caperrors.Error
		if errors.As(err, &capErr) {
			errMsg = capErr.SerializeToRemoteString()
		} else {
			errMsg = caperrors.NewPublicSystemError(err, caperrors.Unknown).SerializeToRemoteString()
		}
		pr.lggr.Warnw("failed to register trigger", "workflowId", pr.workflowID, "err", err)
	} else {
		pr.forwardEvents(callbackCh)
		pr.lggr.Debugw("updated trigger registration", "workflowId", pr.workflowID)
	}

	// If transitioning from unregistered send the status update all peers who requested registration thus far to ensure
	// a timely response on initial registration.
	sendToAllPeers := pr.registrationStatus == types.RegistrationStatus_UNREGISTERED

	pr.registrationStatus = registrationStatus
	pr.registrationErrorMessage = errMsg

	if sendToAllPeers {
		for _, peerID := range requestingPeers {
			pr.sendTriggerRegistrationStatus(peerID, callerDon.ID)
		}
	} else {
		// Send registration status only to the current requester
		// Other peers will get an updated registration status when they re-request registration.
		// The alternative would be to send the status to all requesting peers every time,
		// but that could lead to a lot of duplicate messages being sent if registration keeps failing.
		pr.sendTriggerRegistrationStatus(sender, callerDon.ID)
	}
}

func (pr *PublisherRegistration) forwardEvents(eventsCh <-chan commoncap.TriggerResponse) {
	pr.wg.Add(1)
	go func() {
		defer pr.wg.Done()
		for {
			select {
			case <-pr.stopCh:
				return
			case event, ok := <-eventsCh:
				if !ok {
					close(pr.eventForwardingChannel)
					return
				}
				select {
				case pr.eventForwardingChannel <- event:
				case <-pr.stopCh:
					return
				}
			}
		}
	}()
}

func (pr *PublisherRegistration) Close(ctx context.Context) error {
	close(pr.stopCh)
	pr.wg.Wait()

	if pr.registrationStatus == types.RegistrationStatus_REGISTERED {
		unregisterErr := pr.underlyingTrigger.UnregisterTrigger(ctx, pr.underlyingTriggerRegistrationRequest)
		if unregisterErr != nil {
			return fmt.Errorf("failed to unregister from underlying trigger, workflow ID %s, trigger ID %s: %w",
				pr.workflowID, pr.triggerID, unregisterErr)
		}
	}

	return nil
}

// sendTriggerRegistrationStatus sends a trigger registration response back to the caller DON with an optional error message
func (pr *PublisherRegistration) sendTriggerRegistrationStatus(peerID p2ptypes.PeerID, callerDonID uint32) {
	registrationResponseMessage := &types.MessageBody{
		CapabilityId:    pr.capabilityID,
		CapabilityDonId: pr.capabilityDonID,
		CallerDonId:     callerDonID,
		Method:          types.TriggerRegistrationStatus,
		Metadata: &types.MessageBody_TriggerRegistrationMetadata{
			TriggerRegistrationMetadata: &types.TriggerRegistrationMetadata{
				TriggerId:         pr.triggerID,
				WorkflowId:        pr.workflowID,
				RegistrationError: pr.registrationErrorMessage,
				Status:            pr.registrationStatus,
			},
		},
	}
	err := pr.sendToDispatcher(peerID, registrationResponseMessage)
	if err != nil {
		pr.lggr.Errorw("failed to send trigger registration response", "peerID", peerID, "err", err)
	}
}
