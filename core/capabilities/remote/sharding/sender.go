package sharding

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

type ExecutionStatusUpdateSender struct {
	services.StateMachine
	stopCh     services.StopChan
	dispatcher remotetypes.Dispatcher
	secondary  commoncap.DON
	primaryID  uint32
	lggr       logger.SugaredLogger
}

func NewExecutionStatusUpdateSender(dispatcher remotetypes.Dispatcher, primaryID uint32, secondary commoncap.DON, lggr logger.Logger) *ExecutionStatusUpdateSender {
	return &ExecutionStatusUpdateSender{
		stopCh:     make(services.StopChan),
		dispatcher: dispatcher,
		secondary:  secondary,
		primaryID:  primaryID,
		lggr:       logger.Sugared(logger.With(lggr, "component", "ExecutionStatusUpdateSender")),
	}
}

func (s *ExecutionStatusUpdateSender) Start(ctx context.Context) error {
	return s.StartOnce(s.Name(), func() error { return nil })
}

func (s *ExecutionStatusUpdateSender) Close() error {
	return s.StopOnce(s.Name(), func() error {
		close(s.stopCh)
		return nil
	})
}

func (s *ExecutionStatusUpdateSender) Send(ctx context.Context, msg *ringpb.ExecutionStatusUpdate) {
	payload, err := proto.Marshal(msg)
	if err != nil {
		s.lggr.Errorw("failed to marshal ExecutionStatusUpdate", "err", err)
		return
	}

	messageID := fmt.Sprintf("%s:%s:%d", msg.WorkflowId, msg.TriggerEventId, msg.TriggerIndex)
	for _, peerID := range s.secondary.Members {
		body := &remotetypes.MessageBody{
			Method:          remotetypes.MethodExecutionStatusUpdate,
			Payload:         payload,
			CallerDonId:     s.primaryID,
			CapabilityDonId: s.secondary.ID,
			MessageId:       []byte(messageID),
		}
		if err := s.dispatcher.Send(peerID, body); err != nil {
			s.lggr.Errorw("failed to send ExecutionStatusUpdate", "peerID", peerID, "err", err)
		}
	}
}

func (s *ExecutionStatusUpdateSender) Name() string {
	return fmt.Sprintf("ExecutionStatusUpdateSender-primary-%d", s.primaryID)
}

func (s *ExecutionStatusUpdateSender) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

type ShardHeartbeatSender struct {
	services.StateMachine
	stopCh     services.StopChan
	dispatcher remotetypes.Dispatcher
	secondary  commoncap.DON
	primaryID  uint32
	interval   time.Duration
	lggr       logger.SugaredLogger
	wg         sync.WaitGroup
}

func NewShardHeartbeatSender(dispatcher remotetypes.Dispatcher, primaryID uint32, secondary commoncap.DON, interval time.Duration, lggr logger.Logger) *ShardHeartbeatSender {
	if interval == 0 {
		interval = 30 * time.Second
	}
	return &ShardHeartbeatSender{
		stopCh:     make(services.StopChan),
		dispatcher: dispatcher,
		secondary:  secondary,
		primaryID:  primaryID,
		interval:   interval,
		lggr:       logger.Sugared(logger.With(lggr, "component", "ShardHeartbeatSender")),
	}
}

func (s *ShardHeartbeatSender) Start(ctx context.Context) error {
	return s.StartOnce(s.Name(), func() error {
		s.wg.Add(1)
		go s.heartbeatLoop()
		return nil
	})
}

func (s *ShardHeartbeatSender) Close() error {
	return s.StopOnce(s.Name(), func() error {
		close(s.stopCh)
		s.wg.Wait()
		return nil
	})
}

func (s *ShardHeartbeatSender) heartbeatLoop() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sendHeartbeat(ctx)
		}
	}
}

func (s *ShardHeartbeatSender) sendHeartbeat(ctx context.Context) {
	hb := &ringpb.ShardHeartbeat{
		PrimaryDonId: s.primaryID,
		Timestamp:      time.Now().Unix(),
	}
	payload, err := proto.Marshal(hb)
	if err != nil {
		s.lggr.Errorw("failed to marshal ShardHeartbeat", "err", err)
		return
	}
	messageID := fmt.Sprintf("heartbeat:%d", s.primaryID)
	for _, peerID := range s.secondary.Members {
		body := &remotetypes.MessageBody{
			Method:          remotetypes.MethodShardHeartbeat,
			Payload:         payload,
			CallerDonId:     s.primaryID,
			CapabilityDonId: s.secondary.ID,
			MessageId:       []byte(messageID),
		}
		if err := s.dispatcher.Send(peerID, body); err != nil {
			s.lggr.Errorw("failed to send ShardHeartbeat", "peerID", peerID, "err", err)
		}
	}
}

func (s *ShardHeartbeatSender) Name() string {
	return fmt.Sprintf("ShardHeartbeatSender-primary-%d", s.primaryID)
}

func (s *ShardHeartbeatSender) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}
