package ragep2p

import (
	"math"
	"sync"

	"github.com/smartcontractkit/libocr/ragep2p/internal/msgbuf"
	"github.com/smartcontractkit/libocr/ragep2p/internal/ratelimit"
)

type shouldPushResult int

const (
	_ shouldPushResult = iota
	shouldPushResultYes
	shouldPushResultMessageTooBig
	shouldPushResultMessagesLimitExceeded
	shouldPushResultBytesLimitExceeded
	shouldPushResultUnknownStream
)

type pushResult int

const (
	_ pushResult = iota
	pushResultSuccess
	pushResultDropped
	pushResultUnknownStream
)

type demuxerStream struct {
	buffer          *msgbuf.MessageBuffer
	chSignal        chan struct{}
	maxMessageSize  int
	messagesLimiter ratelimit.TokenBucket
	bytesLimiter    ratelimit.TokenBucket
}

type demuxer struct {
	mutex   sync.Mutex
	streams map[streamID]*demuxerStream
}

func newDemuxer() *demuxer {
	return &demuxer{
		sync.Mutex{},
		map[streamID]*demuxerStream{},
	}
}

func makeRateLimiter(params TokenBucketParams) ratelimit.TokenBucket {
	tb := ratelimit.TokenBucket{}
	tb.SetRate(ratelimit.MillitokensPerSecond(math.Ceil(params.Rate * 1000)))
	tb.SetCapacity(params.Capacity)
	tb.AddTokens(params.Capacity)
	return tb
}

func (d *demuxer) AddStream(
	sid streamID,
	incomingBufferSize int,
	maxMessageSize int,
	messagesLimit TokenBucketParams,
	bytesLimit TokenBucketParams,
) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.streams[sid]; ok {
		return false
	}

	d.streams[sid] = &demuxerStream{
		msgbuf.NewMessageBuffer(incomingBufferSize),
		make(chan struct{}, 1),
		maxMessageSize,
		makeRateLimiter(messagesLimit),
		makeRateLimiter(bytesLimit),
	}
	return true
}

func (d *demuxer) RemoveStream(sid streamID) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.streams, sid)
}

func (d *demuxer) ShouldPush(sid streamID, size int) shouldPushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return shouldPushResultUnknownStream
	}

	if size > s.maxMessageSize {
		return shouldPushResultMessageTooBig
	}

	messagesLimiterAllow := s.messagesLimiter.RemoveTokens(1)
	bytesLimiterAllow := s.bytesLimiter.RemoveTokens(uint32(size))

	if !messagesLimiterAllow {
		return shouldPushResultMessagesLimitExceeded
	}

	if !bytesLimiterAllow {
		return shouldPushResultBytesLimitExceeded
	}

	return shouldPushResultYes
}

func (d *demuxer) PushMessage(sid streamID, msg []byte) pushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return pushResultUnknownStream
	}

	var result pushResult
	if s.buffer.Push(msg) == nil {
		result = pushResultSuccess
	} else {
		result = pushResultDropped
	}

	select {
	case s.chSignal <- struct{}{}:
	default:
	}

	return result
}

func (d *demuxer) PopMessage(sid streamID) []byte {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return nil
	}

	result := s.buffer.Pop()

	if s.buffer.Peek() != nil {
		select {
		case s.chSignal <- struct{}{}:
		default:
		}
	}

	return result
}

func (d *demuxer) SignalPending(sid streamID) <-chan struct{} {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return nil
	}

	return s.chSignal
}
