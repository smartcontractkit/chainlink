package streams

import (
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// alias for easier refactoring
type StreamID = llo.StreamID

type Registry interface {
	Getter
	Register(jb job.Job, rrs ResultRunSaver) error
	Unregister(jobID int32)
}

type Getter interface {
	Get(streamID StreamID) (p Pipeline, exists bool)
}

type streamRegistry struct {
	sync.RWMutex
	lggr   logger.Logger
	runner Runner
	// keyed by stream ID
	pipelines map[StreamID]Pipeline
	// keyed by job ID
	pipelinesByJobID map[int32]Pipeline
}

func NewRegistry(lggr logger.Logger, runner Runner) Registry {
	return newRegistry(lggr, runner)
}

func newRegistry(lggr logger.Logger, runner Runner) *streamRegistry {
	return &streamRegistry{
		sync.RWMutex{},
		lggr.Named("Registry"),
		runner,
		make(map[StreamID]Pipeline),
		make(map[int32]Pipeline),
	}
}

func (s *streamRegistry) Get(streamID StreamID) (p Pipeline, exists bool) {
	s.RLock()
	defer s.RUnlock()
	p, exists = s.pipelines[streamID]
	return
}

func (s *streamRegistry) Register(jb job.Job, rrs ResultRunSaver) error {
	if jb.Type != job.Stream {
		return fmt.Errorf("cannot register job type %s; only Stream jobs are supported", jb.Type)
	}
	p, err := NewMultiStreamPipeline(s.lggr, jb, s.runner, rrs)
	if err != nil {
		return fmt.Errorf("cannot register job with ID: %d; %w", jb.ID, err)
	}
	s.Lock()
	defer s.Unlock()
	if _, exists := s.pipelinesByJobID[jb.ID]; exists {
		return fmt.Errorf("cannot register job with ID: %d; it is already registered", jb.ID)
	}
	for _, strmID := range p.StreamIDs() {
		if _, exists := s.pipelines[strmID]; exists {
			return fmt.Errorf("cannot register job with ID: %d; stream id %d is already registered", jb.ID, strmID)
		}
	}
	s.pipelinesByJobID[jb.ID] = p
	streamIDs := p.StreamIDs()
	for _, strmID := range streamIDs {
		s.pipelines[strmID] = p
	}
	return nil
}

func (s *streamRegistry) Unregister(jobID int32) {
	s.Lock()
	defer s.Unlock()
	p, exists := s.pipelinesByJobID[jobID]
	if !exists {
		return
	}
	streamIDs := p.StreamIDs()
	for _, id := range streamIDs {
		delete(s.pipelines, id)
	}
}
