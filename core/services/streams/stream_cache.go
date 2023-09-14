package streams

import (
	"context"

	"github.com/smartcontractkit/chainlink-data-streams/streams"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type StreamCacheORM interface {
	LoadStreams(ctx context.Context, lggr logger.Logger, runner Runner, m map[streams.StreamID]Stream) error
}

type StreamCache interface {
	Get(streamID streams.StreamID) (Stream, bool)
	Load(ctx context.Context, lggr logger.Logger, runner Runner) error
}

type streamCache struct {
	orm     StreamCacheORM
	streams map[streams.StreamID]Stream
}

func NewStreamCache(orm StreamCacheORM) StreamCache {
	return newStreamCache(orm)
}

func newStreamCache(orm StreamCacheORM) *streamCache {
	return &streamCache{
		orm,
		make(map[streams.StreamID]Stream),
	}
}

func (s *streamCache) Get(streamID streams.StreamID) (Stream, bool) {
	strm, exists := s.streams[streamID]
	return strm, exists
}

func (s *streamCache) Load(ctx context.Context, lggr logger.Logger, runner Runner) error {
	return s.orm.LoadStreams(ctx, lggr, runner, s.streams)
}
