package streams

import (
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type StreamID = string

type Registry interface {
	Get(streamID StreamID) (strm Stream, exists bool)
	Register(streamID StreamID, spec pipeline.Spec, rrs ResultRunSaver) error
	Unregister(streamID StreamID)
}

type streamRegistry struct {
	sync.RWMutex
	lggr    logger.Logger
	runner  Runner
	streams map[StreamID]Stream
}

func NewRegistry(lggr logger.Logger, runner Runner) Registry {
	return newRegistry(lggr, runner)
}

func newRegistry(lggr logger.Logger, runner Runner) *streamRegistry {
	return &streamRegistry{
		sync.RWMutex{},
		lggr.Named("Registry"),
		runner,
		make(map[StreamID]Stream),
	}
}

func (s *streamRegistry) Get(streamID StreamID) (strm Stream, exists bool) {
	s.RLock()
	defer s.RUnlock()
	strm, exists = s.streams[streamID]
	return
}

func (s *streamRegistry) Register(streamID StreamID, spec pipeline.Spec, rrs ResultRunSaver) error {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.streams[streamID]; exists {
		return fmt.Errorf("stream already registered for id: %q", streamID)
	}
	s.streams[streamID] = NewStream(s.lggr, streamID, spec, s.runner, rrs)
	return nil
}

func (s *streamRegistry) Unregister(streamID StreamID) {
	s.Lock()
	defer s.Unlock()
	delete(s.streams, streamID)
}
