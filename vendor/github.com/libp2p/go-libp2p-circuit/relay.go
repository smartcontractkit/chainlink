package relay

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/libp2p/go-libp2p-circuit/pb"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"

	pool "github.com/libp2p/go-buffer-pool"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"

	logging "github.com/ipfs/go-log"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("relay")

const ProtoID = "/libp2p/circuit/relay/0.1.0"

const maxMessageSize = 4096

var (
	RelayAcceptTimeout   = 10 * time.Second
	HopConnectTimeout    = 30 * time.Second
	StopHandshakeTimeout = 1 * time.Minute

	HopStreamBufferSize = 4096
	HopStreamLimit      = 1 << 19 // 512K hops for 1M goroutines
)

// Relay is the relay transport and service.
type Relay struct {
	host     host.Host
	upgrader *tptu.Upgrader
	ctx      context.Context
	self     peer.ID

	active bool
	hop    bool

	incoming chan *Conn

	// atomic counters
	streamCount  int32
	liveHopCount int32

	// per peer hop counters
	mx       sync.Mutex
	hopCount map[peer.ID]int
}

// RelayOpts are options for configuring the relay transport.
type RelayOpt int

var (
	// OptActive configures the relay transport to actively establish
	// outbound connections on behalf of clients. You probably don't want to
	// enable this unless you know what you're doing.
	OptActive = RelayOpt(0)
	// OptHop configures the relay transport to accept requests to relay
	// traffic on behalf of third-parties. Unless OptActive is specified,
	// this will only relay traffic between peers already connected to this
	// node.
	OptHop = RelayOpt(1)
	// OptDiscovery is a no-op. It was introduced as a way to probe new
	// peers to see if they were willing to act as a relays. However, in
	// practice, it's useless. While it does test to see if these peers are
	// relays, it doesn't (and can't), check to see if these peers are
	// _active_ relays (i.e., will actively dial the target peer).
	//
	// This option may be re-enabled in the future but for now you shouldn't
	// use it.
	OptDiscovery = RelayOpt(2)
)

type RelayError struct {
	Code pb.CircuitRelay_Status
}

func (e RelayError) Error() string {
	return fmt.Sprintf("error opening relay circuit: %s (%d)", pb.CircuitRelay_Status_name[int32(e.Code)], e.Code)
}

// NewRelay constructs a new relay.
func NewRelay(ctx context.Context, h host.Host, upgrader *tptu.Upgrader, opts ...RelayOpt) (*Relay, error) {
	r := &Relay{
		upgrader: upgrader,
		host:     h,
		ctx:      ctx,
		self:     h.ID(),
		incoming: make(chan *Conn),
		hopCount: make(map[peer.ID]int),
	}

	for _, opt := range opts {
		switch opt {
		case OptActive:
			r.active = true
		case OptHop:
			r.hop = true
		case OptDiscovery:
			log.Errorf(
				"circuit.OptDiscovery is now a no-op: %s",
				"dialing peers with a random relay is no longer supported",
			)
		default:
			return nil, fmt.Errorf("unrecognized option: %d", opt)
		}
	}

	h.SetStreamHandler(ProtoID, r.handleNewStream)

	return r, nil
}

// Increment the live hop count and increment the connection manager tags by 1 for the two
// sides of the hop stream. This ensures that connections with many hop streams will be protected
// from pruning, thus minimizing disruption from connection trimming in a relay node.
func (r *Relay) addLiveHop(from, to peer.ID) {
	atomic.AddInt32(&r.liveHopCount, 1)
	r.host.ConnManager().UpsertTag(from, "relay-hop-stream", incrementTag)
	r.host.ConnManager().UpsertTag(to, "relay-hop-stream", incrementTag)
}

// Decrement the live hpo count and decrement the connection manager tags for the two sides
// of the hop stream.
func (r *Relay) rmLiveHop(from, to peer.ID) {
	atomic.AddInt32(&r.liveHopCount, -1)
	r.host.ConnManager().UpsertTag(from, "relay-hop-stream", decrementTag)
	r.host.ConnManager().UpsertTag(to, "relay-hop-stream", decrementTag)

}

func (r *Relay) GetActiveHops() int32 {
	return atomic.LoadInt32(&r.liveHopCount)
}

func (r *Relay) DialPeer(ctx context.Context, relay peer.AddrInfo, dest peer.AddrInfo) (*Conn, error) {

	log.Debugf("dialing peer %s through relay %s", dest.ID, relay.ID)

	if len(relay.Addrs) > 0 {
		r.host.Peerstore().AddAddrs(relay.ID, relay.Addrs, peerstore.TempAddrTTL)
	}

	s, err := r.host.NewStream(ctx, relay.ID, ProtoID)
	if err != nil {
		return nil, err
	}

	rd := newDelimitedReader(s, maxMessageSize)
	wr := newDelimitedWriter(s)
	defer rd.Close()

	var msg pb.CircuitRelay

	msg.Type = pb.CircuitRelay_HOP.Enum()
	msg.SrcPeer = peerInfoToPeer(r.host.Peerstore().PeerInfo(r.self))
	msg.DstPeer = peerInfoToPeer(dest)

	err = wr.WriteMsg(&msg)
	if err != nil {
		s.Reset()
		return nil, err
	}

	msg.Reset()

	err = rd.ReadMsg(&msg)
	if err != nil {
		s.Reset()
		return nil, err
	}

	if msg.GetType() != pb.CircuitRelay_STATUS {
		s.Reset()
		return nil, fmt.Errorf("unexpected relay response; not a status message (%d)", msg.GetType())
	}

	if msg.GetCode() != pb.CircuitRelay_SUCCESS {
		s.Reset()
		return nil, RelayError{msg.GetCode()}
	}

	return &Conn{stream: s, remote: dest, host: r.host, relay: r}, nil
}

func (r *Relay) Matches(addr ma.Multiaddr) bool {
	// TODO: Look at the prefix transport as well.
	_, err := addr.ValueForProtocol(P_CIRCUIT)
	return err == nil
}

// Queries a peer for support of hop relay
func CanHop(ctx context.Context, host host.Host, id peer.ID) (bool, error) {
	s, err := host.NewStream(ctx, id, ProtoID)
	if err != nil {
		return false, err
	}
	defer s.Close()

	rd := newDelimitedReader(s, maxMessageSize)
	wr := newDelimitedWriter(s)
	defer rd.Close()

	var msg pb.CircuitRelay

	msg.Type = pb.CircuitRelay_CAN_HOP.Enum()

	if err := wr.WriteMsg(&msg); err != nil {
		s.Reset()
		return false, err
	}

	msg.Reset()

	if err := rd.ReadMsg(&msg); err != nil {
		s.Reset()
		return false, err
	}

	if msg.GetType() != pb.CircuitRelay_STATUS {
		return false, fmt.Errorf("unexpected relay response; not a status message (%d)", msg.GetType())
	}

	return msg.GetCode() == pb.CircuitRelay_SUCCESS, nil
}

func (r *Relay) CanHop(ctx context.Context, id peer.ID) (bool, error) {
	return CanHop(ctx, r.host, id)
}

func (r *Relay) handleNewStream(s network.Stream) {
	log.Infof("new relay stream from: %s", s.Conn().RemotePeer())

	rd := newDelimitedReader(s, maxMessageSize)
	defer rd.Close()

	var msg pb.CircuitRelay

	err := rd.ReadMsg(&msg)
	if err != nil {
		r.handleError(s, pb.CircuitRelay_MALFORMED_MESSAGE)
		return
	}

	switch msg.GetType() {
	case pb.CircuitRelay_HOP:
		r.handleHopStream(s, &msg)
	case pb.CircuitRelay_STOP:
		r.handleStopStream(s, &msg)
	case pb.CircuitRelay_CAN_HOP:
		r.handleCanHop(s, &msg)
	default:
		log.Warnf("unexpected relay handshake: %d", msg.GetType())
		r.handleError(s, pb.CircuitRelay_MALFORMED_MESSAGE)
	}
}

func (r *Relay) handleHopStream(s network.Stream, msg *pb.CircuitRelay) {
	if !r.hop {
		r.handleError(s, pb.CircuitRelay_HOP_CANT_SPEAK_RELAY)
		return
	}

	streamCount := atomic.AddInt32(&r.streamCount, 1)
	liveHopCount := atomic.LoadInt32(&r.liveHopCount)
	defer atomic.AddInt32(&r.streamCount, -1)

	if (streamCount + liveHopCount) > int32(HopStreamLimit) {
		log.Warn("hop stream limit exceeded; resetting stream")
		s.Reset()
		return
	}

	src, err := peerToPeerInfo(msg.GetSrcPeer())
	if err != nil {
		r.handleError(s, pb.CircuitRelay_HOP_SRC_MULTIADDR_INVALID)
		return
	}

	if src.ID != s.Conn().RemotePeer() {
		r.handleError(s, pb.CircuitRelay_HOP_SRC_MULTIADDR_INVALID)
		return
	}

	dst, err := peerToPeerInfo(msg.GetDstPeer())
	if err != nil {
		r.handleError(s, pb.CircuitRelay_HOP_DST_MULTIADDR_INVALID)
		return
	}

	if dst.ID == r.self {
		r.handleError(s, pb.CircuitRelay_HOP_CANT_RELAY_TO_SELF)
		return
	}

	// open stream
	ctx, cancel := context.WithTimeout(r.ctx, HopConnectTimeout)
	defer cancel()

	if !r.active {
		ctx = network.WithNoDial(ctx, "relay hop")
	} else if len(dst.Addrs) > 0 {
		r.host.Peerstore().AddAddrs(dst.ID, dst.Addrs, peerstore.TempAddrTTL)
	}

	bs, err := r.host.NewStream(ctx, dst.ID, ProtoID)
	if err != nil {
		log.Debugf("error opening relay stream to %s: %s", dst.ID.Pretty(), err.Error())
		if err == network.ErrNoConn {
			r.handleError(s, pb.CircuitRelay_HOP_NO_CONN_TO_DST)
		} else {
			r.handleError(s, pb.CircuitRelay_HOP_CANT_DIAL_DST)
		}
		return
	}

	// stop handshake
	rd := newDelimitedReader(bs, maxMessageSize)
	wr := newDelimitedWriter(bs)
	defer rd.Close()

	// set handshake deadline
	bs.SetDeadline(time.Now().Add(StopHandshakeTimeout))

	msg.Type = pb.CircuitRelay_STOP.Enum()

	err = wr.WriteMsg(msg)
	if err != nil {
		log.Debugf("error writing stop handshake: %s", err.Error())
		bs.Reset()
		r.handleError(s, pb.CircuitRelay_HOP_CANT_OPEN_DST_STREAM)
		return
	}

	msg.Reset()

	err = rd.ReadMsg(msg)
	if err != nil {
		log.Debugf("error reading stop response: %s", err.Error())
		bs.Reset()
		r.handleError(s, pb.CircuitRelay_HOP_CANT_OPEN_DST_STREAM)
		return
	}

	if msg.GetType() != pb.CircuitRelay_STATUS {
		log.Debugf("unexpected relay stop response: not a status message (%d)", msg.GetType())
		bs.Reset()
		r.handleError(s, pb.CircuitRelay_HOP_CANT_OPEN_DST_STREAM)
		return
	}

	if msg.GetCode() != pb.CircuitRelay_SUCCESS {
		log.Debugf("relay stop failure: %d", msg.GetCode())
		bs.Reset()
		r.handleError(s, msg.GetCode())
		return
	}

	err = r.writeResponse(s, pb.CircuitRelay_SUCCESS)
	if err != nil {
		log.Debugf("error writing relay response: %s", err.Error())
		bs.Reset()
		s.Reset()
		return
	}

	// relay connection
	log.Infof("relaying connection between %s and %s", src.ID.Pretty(), dst.ID.Pretty())

	// reset deadline
	bs.SetDeadline(time.Time{})

	r.addLiveHop(src.ID, dst.ID)

	goroutines := new(int32)
	*goroutines = 2
	done := func() {
		if atomic.AddInt32(goroutines, -1) == 0 {
			s.Close()
			bs.Close()
			r.rmLiveHop(src.ID, dst.ID)
		}
	}

	// Don't reset streams after finishing or the other side will get an
	// error, not an EOF.
	go func() {
		defer done()

		buf := pool.Get(HopStreamBufferSize)
		defer pool.Put(buf)

		count, err := io.CopyBuffer(s, bs, buf)
		if err != nil {
			log.Debugf("relay copy error: %s", err)
			// Reset both.
			s.Reset()
			bs.Reset()
		} else {
			// propagate the close
			s.CloseWrite()
		}
		log.Debugf("relayed %d bytes from %s to %s", count, dst.ID.Pretty(), src.ID.Pretty())
	}()

	go func() {
		defer done()

		buf := pool.Get(HopStreamBufferSize)
		defer pool.Put(buf)

		count, err := io.CopyBuffer(bs, s, buf)
		if err != nil {
			log.Debugf("relay copy error: %s", err)
			// Reset both.
			bs.Reset()
			s.Reset()
		} else {
			// propagate the close
			bs.CloseWrite()
		}
		log.Debugf("relayed %d bytes from %s to %s", count, src.ID.Pretty(), dst.ID.Pretty())
	}()
}

func (r *Relay) handleStopStream(s network.Stream, msg *pb.CircuitRelay) {
	src, err := peerToPeerInfo(msg.GetSrcPeer())
	if err != nil {
		r.handleError(s, pb.CircuitRelay_STOP_SRC_MULTIADDR_INVALID)
		return
	}

	dst, err := peerToPeerInfo(msg.GetDstPeer())
	if err != nil || dst.ID != r.self {
		r.handleError(s, pb.CircuitRelay_STOP_DST_MULTIADDR_INVALID)
		return
	}

	log.Infof("relay connection from: %s", src.ID)

	if len(src.Addrs) > 0 {
		r.host.Peerstore().AddAddrs(src.ID, src.Addrs, peerstore.TempAddrTTL)
	}

	select {
	case r.incoming <- &Conn{stream: s, remote: src, host: r.host, relay: r}:
	case <-time.After(RelayAcceptTimeout):
		r.handleError(s, pb.CircuitRelay_STOP_RELAY_REFUSED)
	}
}

func (r *Relay) handleCanHop(s network.Stream, msg *pb.CircuitRelay) {
	var err error

	if r.hop {
		err = r.writeResponse(s, pb.CircuitRelay_SUCCESS)
	} else {
		err = r.writeResponse(s, pb.CircuitRelay_HOP_CANT_SPEAK_RELAY)
	}

	if err != nil {
		s.Reset()
		log.Debugf("error writing relay response: %s", err.Error())
	} else {
		s.Close()
	}
}

func (r *Relay) handleError(s network.Stream, code pb.CircuitRelay_Status) {
	log.Warnf("relay error: %s (%d)", pb.CircuitRelay_Status_name[int32(code)], code)
	err := r.writeResponse(s, code)
	if err != nil {
		s.Reset()
		log.Debugf("error writing relay response: %s", err.Error())
	} else {
		s.Close()
	}
}

func (r *Relay) writeResponse(s network.Stream, code pb.CircuitRelay_Status) error {
	wr := newDelimitedWriter(s)

	var msg pb.CircuitRelay
	msg.Type = pb.CircuitRelay_STATUS.Enum()
	msg.Code = code.Enum()

	return wr.WriteMsg(&msg)
}
