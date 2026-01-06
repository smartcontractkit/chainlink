package registration

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/log"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

const (
	// This default is taken from the legacy code, it is the buffer size of the channel
	sendChannelBufferSize = 1000
)

const registrationStatusUpdateCacheKey = "subscriber_registration"

type SubscriberRegistration struct {
	lggr              logger.Logger
	triggerResponseCh chan commoncap.TriggerResponse
	rawRequest        []byte

	registrationStatusUpdateCache *messagecache.MessageCache[string, p2ptypes.PeerID]

	statusMutex        sync.RWMutex
	registrationStatus types.RegistrationStatus
	registrationError  caperrors.Error

	// Used to signal that this subscription is no longer waiting for the initial registration response
	initialRegistrationResponseChan chan struct{}
}

func NewSubscriberRegistration(lggr logger.Logger, rawRequest []byte,
	workflowID string, triggerID string) *SubscriberRegistration {
	return &SubscriberRegistration{
		lggr:                            logger.With(lggr, "WorkflowID", workflowID, "TriggerID", triggerID),
		triggerResponseCh:               make(chan commoncap.TriggerResponse, sendChannelBufferSize),
		rawRequest:                      rawRequest,
		registrationStatusUpdateCache:   messagecache.NewMessageCache[string, p2ptypes.PeerID](),
		initialRegistrationResponseChan: make(chan struct{}),
	}
}

func (sr *SubscriberRegistration) HandleTriggerRegistrationStatusUpdate(sender p2ptypes.PeerID, msg *types.MessageBody, minResponseToAggregate uint32,
	registrationExpiry time.Duration, publisherNodeCount int) {
	nowMs := time.Now().UnixMilli()
	// First remove any expired messages from the cache, the cache is not expected to be large so this is acceptable
	// and uses no more resource that the ready call in this case so need to clean-up on a background goroutine.
	sr.registrationStatusUpdateCache.DeleteOlderThan(nowMs - registrationExpiry.Milliseconds())

	meta := msg.GetTriggerRegistrationMetadata()

	if meta.Status == types.RegistrationStatus_REGISTRATION_ERROR {
		sr.registrationStatusUpdateCache.Insert(registrationStatusUpdateCacheKey, sender, nowMs, []byte(meta.RegistrationError))
	} else {
		sr.registrationStatusUpdateCache.Insert(registrationStatusUpdateCacheKey, sender, nowMs, nil)
	}

	ready, _, registrationResponses := sr.registrationStatusUpdateCache.Ready(registrationStatusUpdateCacheKey, minResponseToAggregate, nowMs-registrationExpiry.Milliseconds(), false)

	if ready {
		var successfulRegistrationCount uint32
		var totalErrorCount int

		// aggregate errors by message
		errorToCount := map[string]uint32{}
		for _, responseError := range registrationResponses {
			if len(responseError) > 0 {
				errorStr := string(responseError)
				errorToCount[errorStr]++
				totalErrorCount++
			} else {
				successfulRegistrationCount++
			}
		}

		if successfulRegistrationCount >= minResponseToAggregate {
			// Successful registration
			sr.lggr.Infow("successful trigger registration", "sender", sender)
			sr.setRegistrationStatus(types.RegistrationStatus_REGISTERED, nil)
		} else {
			errStrs := make([]string, 0, len(errorToCount))
			for errStr := range errorToCount {
				errStrs = append(errStrs, errStr)
			}
			sort.Strings(errStrs)

			// Registration failed - set error response
			// Is there a consensus error?  if so set that
			lastErr := ""

			for _, errStr := range errStrs {
				count := errorToCount[errStr]
				if count >= minResponseToAggregate {
					caperr := caperrors.DeserializeErrorFromString(errStr)
					sr.setRegistrationStatus(types.RegistrationStatus_REGISTRATION_ERROR, caperr)
					return
				}

				lastErr = errStr
			}

			// The check on when to give up waiting for minResponseToAggregate identical errors assumes that all received
			// errors so far, and any future that will be received are and will be distinct.  It's the same logic
			// used to aggregate errors for remote executable capabilities.  For the purposes of error handling it is
			// sufficient.
			if totalErrorCount >= publisherNodeCount-int(minResponseToAggregate)+1 {
				// There is no consensus error, return a generic error message
				sr.setRegistrationStatus(types.RegistrationStatus_REGISTRATION_ERROR,
					caperrors.NewPublicSystemError(fmt.Errorf("received %d errors, last error %s : %s", totalErrorCount, msg.Error, log.SanitizeLogString(lastErr)), caperrors.ConsensusFailed))

				// Log all the errors received to help diagnose why consensus failed.
				// First sanitize the error messages and then log them in a single log line to ensure they are not interleaved with other log messages.
				var sanitizedErrStrs []string
				for _, errStr := range errStrs {
					sanitizedErrStrs = append(sanitizedErrStrs, log.SanitizeLogString(errStr))
				}
				sr.lggr.Warnw("failed to achieve consensus on trigger registration errors", "errors", sanitizedErrStrs)
			}
		}
	}
}

func (sr *SubscriberRegistration) GetTriggerResponseChannel() chan commoncap.TriggerResponse {
	return sr.triggerResponseCh
}

// AwaitRegistration waits for the registration until the context is cancelled.  Returns nil if registration is
// successful, returns a capability error if registration failed.
// If the registration status is unable to be determined before the context is cancelled ErrUnableToDetermineRegistrationStatus is returned.
func (sr *SubscriberRegistration) AwaitRegistration(ctx context.Context) error {
	registrationStatus, registrationErr := sr.getRegistrationStatus()
	switch registrationStatus {
	case types.RegistrationStatus_REGISTERED:
		return nil
	case types.RegistrationStatus_REGISTRATION_ERROR:
		return registrationErr
	default:
		// Wait for registration response or context cancellation
		select {
		case <-ctx.Done():
			return commoncap.ErrUnableToDetermineRegistrationStatus
		case <-sr.initialRegistrationResponseChan:
			_, registrationErr = sr.getRegistrationStatus()
			return registrationErr
		}
	}
}

func (sr *SubscriberRegistration) setRegistrationStatus(status types.RegistrationStatus, registrationError caperrors.Error) {
	sr.statusMutex.Lock()
	defer sr.statusMutex.Unlock()
	if sr.registrationStatus == types.RegistrationStatus_UNREGISTERED && status != types.RegistrationStatus_UNREGISTERED {
		// Signal that the initial registration response is ready
		close(sr.initialRegistrationResponseChan)
	}
	sr.registrationStatus = status
	sr.registrationError = registrationError
}

func (sr *SubscriberRegistration) getRegistrationStatus() (types.RegistrationStatus, caperrors.Error) {
	sr.statusMutex.RLock()
	defer sr.statusMutex.RUnlock()
	return sr.registrationStatus, sr.registrationError
}

func (sr *SubscriberRegistration) UpdateRegistrationRequest(rawRequest []byte) {
	sr.rawRequest = rawRequest
}

func (sr *SubscriberRegistration) GetRawRequest() []byte {
	return sr.rawRequest
}

func (sr *SubscriberRegistration) SendAggregatedEvent(aggregatedEvent commoncap.TriggerResponse) {
	sr.triggerResponseCh <- aggregatedEvent
}

func (sr *SubscriberRegistration) Close() {
	close(sr.triggerResponseCh)
}
