package streams

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type StreamCacheORM interface {
	LoadStreams(ctx context.Context, lggr logger.Logger, runner Runner, m map[commontypes.StreamID]Stream) error
}

type StreamCache interface {
	Get(streamID commontypes.StreamID) (Stream, bool)
	Load(ctx context.Context, lggr logger.Logger, runner Runner) error
}

type streamCache struct {
	orm     StreamCacheORM
	streams map[commontypes.StreamID]Stream
}

func NewStreamCache(orm StreamCacheORM) StreamCache {
	return newStreamCache(orm)
}

func newStreamCache(orm StreamCacheORM) *streamCache {
	return &streamCache{
		orm,
		make(map[commontypes.StreamID]Stream),
	}
}

func (s *streamCache) Get(streamID commontypes.StreamID) (Stream, bool) {
	strm, exists := s.streams[streamID]
	return strm, exists
}

func (s *streamCache) Load(ctx context.Context, lggr logger.Logger, runner Runner) error {
	return s.orm.LoadStreams(ctx, lggr, runner, s.streams)
}
