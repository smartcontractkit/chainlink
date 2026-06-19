package remote

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const registrationStatusUpdateCacheKey = "subscriber_registration"

type triggerRegistrationState struct {
	lggr              logger.Logger
	triggerResponseCh chan commoncap.TriggerResponse
	rawRequest        []byte

	statusUpdateCache *messagecache.MessageCache[string, p2ptypes.PeerID]

	statusMu                        sync.RWMutex
	registrationStatus              types.RegistrationStatus
	registrationError               caperrors.Error
	registrationFinalized           bool
	initialRegistrationResponseChan chan struct{}
}

func newTriggerRegistrationState(lggr logger.Logger, rawRequest []byte, workflowID, triggerID string) *triggerRegistrationState {
	return &triggerRegistrationState{
		lggr:                            logger.With(lggr, "workflowID", workflowID, "triggerID", triggerID),
		triggerResponseCh:               make(chan commoncap.TriggerResponse, sendChannelBufferSize),
		rawRequest:                      rawRequest,
		statusUpdateCache:               messagecache.NewMessageCache[string, p2ptypes.PeerID](),
		initialRegistrationResponseChan: make(chan struct{}),
	}
}

func (s *triggerRegistrationState) handleRegistrationStatusUpdate(
	sender p2ptypes.PeerID,
	msg *types.MessageBody,
	minResponsesToAggregate uint32,
	registrationExpiry time.Duration,
	publisherNodeCount int,
) {
	nowMs := time.Now().UnixMilli()
	s.statusUpdateCache.DeleteOlderThan(nowMs - registrationExpiry.Milliseconds())

	meta := msg.GetTriggerRegistrationMetadata()
	if meta == nil {
		s.lggr.Errorw("received trigger registration status with nil metadata", "sender", sender)
		return
	}

	if meta.Status == types.RegistrationStatus_REGISTRATION_ERROR {
		s.statusUpdateCache.Insert(registrationStatusUpdateCacheKey, sender, nowMs, []byte(meta.RegistrationError))
	} else {
		s.statusUpdateCache.Insert(registrationStatusUpdateCacheKey, sender, nowMs, nil)
	}

	ready, _, registrationResponses := s.statusUpdateCache.Ready(
		registrationStatusUpdateCacheKey,
		minResponsesToAggregate,
		nowMs-registrationExpiry.Milliseconds(),
		false,
	)
	if !ready {
		return
	}

	var successfulRegistrationCount uint32
	totalErrorCount := 0
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

	if successfulRegistrationCount >= minResponsesToAggregate {
		s.lggr.Infow("successful remote trigger registration", "sender", sender)
		s.setRegistrationStatus(types.RegistrationStatus_REGISTERED, nil)
		return
	}

	errStrs := make([]string, 0, len(errorToCount))
	for errStr := range errorToCount {
		errStrs = append(errStrs, errStr)
	}
	sort.Strings(errStrs)

	lastErr := ""
	for _, errStr := range errStrs {
		count := errorToCount[errStr]
		if count >= minResponsesToAggregate {
			capErr := caperrors.DeserializeErrorFromString(errStr)
			s.setRegistrationStatus(types.RegistrationStatus_REGISTRATION_ERROR, capErr)
			return
		}
		lastErr = errStr
	}

	if totalErrorCount >= publisherNodeCount-int(minResponsesToAggregate)+1 {
		s.setRegistrationStatus(
			types.RegistrationStatus_REGISTRATION_ERROR,
			caperrors.NewPublicSystemError(
				fmt.Errorf("received %d errors, last error: %s", totalErrorCount, SanitizeLogString(lastErr)),
				caperrors.ConsensusFailed,
			),
		)
		s.lggr.Warnw("failed to achieve consensus on trigger registration errors", "errors", errStrs)
	}
}

func (s *triggerRegistrationState) awaitRegistration(ctx context.Context) error {
	registrationStatus, registrationErr := s.getRegistrationStatus()
	switch registrationStatus {
	case types.RegistrationStatus_REGISTERED:
		return nil
	case types.RegistrationStatus_REGISTRATION_ERROR:
		return registrationErr
	default:
		select {
		case <-ctx.Done():
			return commoncap.ErrUnableToDetermineRegistrationStatus
		case <-s.initialRegistrationResponseChan:
			_, registrationErr = s.getRegistrationStatus()
			return registrationErr
		}
	}
}

func (s *triggerRegistrationState) setRegistrationStatus(status types.RegistrationStatus, registrationError caperrors.Error) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	if s.registrationStatus == types.RegistrationStatus_UNREGISTERED && status != types.RegistrationStatus_UNREGISTERED {
		close(s.initialRegistrationResponseChan)
	}
	s.registrationStatus = status
	s.registrationError = registrationError
	if status == types.RegistrationStatus_REGISTERED || status == types.RegistrationStatus_REGISTRATION_ERROR {
		s.registrationFinalized = true
	}
}

func (s *triggerRegistrationState) getRegistrationStatus() (types.RegistrationStatus, caperrors.Error) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.registrationStatus, s.registrationError
}

func (s *triggerRegistrationState) isRegistrationFinalized() bool {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.registrationFinalized
}

func (s *triggerRegistrationState) updateRegistrationRequest(rawRequest []byte) {
	s.rawRequest = rawRequest
}

func (s *triggerRegistrationState) getRawRequest() []byte {
	return s.rawRequest
}

func (s *triggerRegistrationState) getTriggerResponseChannel() chan commoncap.TriggerResponse {
	return s.triggerResponseCh
}

func (s *triggerRegistrationState) close() {
	close(s.triggerResponseCh)
}
