package relay

import (
	"fmt"
	"net"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

// HopTagWeight is the connection manager weight for connections carrying relay hop streams
var HopTagWeight = 5

type Conn struct {
	stream network.Stream
	remote peer.AddrInfo
	host   host.Host
	relay  *Relay
}

type NetAddr struct {
	Relay  string
	Remote string
}

func (n *NetAddr) Network() string {
	return "libp2p-circuit-relay"
}

func (n *NetAddr) String() string {
	return fmt.Sprintf("relay[%s-%s]", n.Remote, n.Relay)
}

func (c *Conn) Close() error {
	c.untagHop()
	return c.stream.Reset()
}

func (c *Conn) Read(buf []byte) (int, error) {
	return c.stream.Read(buf)
}

func (c *Conn) Write(buf []byte) (int, error) {
	return c.stream.Write(buf)
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}

func (c *Conn) RemoteAddr() net.Addr {
	return &NetAddr{
		Relay:  c.stream.Conn().RemotePeer().Pretty(),
		Remote: c.remote.ID.Pretty(),
	}
}

// Increment the underlying relay connection tag by 1, thus increasing its protection from
// connection pruning. This ensures that connections to relays are not accidentally closed,
// by the connection manager, taking with them all the relayed connections (that may themselves
// be protected).
func (c *Conn) tagHop() {
	c.relay.mx.Lock()
	defer c.relay.mx.Unlock()

	p := c.stream.Conn().RemotePeer()
	c.relay.hopCount[p]++
	if c.relay.hopCount[p] == 1 {
		c.host.ConnManager().TagPeer(p, "relay-hop-stream", HopTagWeight)
	}
}

// Decrement the underlying relay connection tag by 1; this is performed when we close the
// relayed connection.
func (c *Conn) untagHop() {
	c.relay.mx.Lock()
	defer c.relay.mx.Unlock()

	p := c.stream.Conn().RemotePeer()
	c.relay.hopCount[p]--
	if c.relay.hopCount[p] == 0 {
		c.host.ConnManager().UntagPeer(p, "relay-hop-stream")
		delete(c.relay.hopCount, p)
	}
}

// TODO: is it okay to cast c.Conn().RemotePeer() into a multiaddr? might be "user input"
func (c *Conn) RemoteMultiaddr() ma.Multiaddr {
	// TODO: We should be able to do this directly without converting to/from a string.
	relayAddr, err := ma.NewComponent(
		ma.ProtocolWithCode(ma.P_P2P).Name,
		c.stream.Conn().RemotePeer().Pretty(),
	)
	if err != nil {
		panic(err)
	}
	return ma.Join(c.stream.Conn().RemoteMultiaddr(), relayAddr, circuitAddr)
}

func (c *Conn) LocalMultiaddr() ma.Multiaddr {
	return c.stream.Conn().LocalMultiaddr()
}

func (c *Conn) LocalAddr() net.Addr {
	na, err := manet.ToNetAddr(c.stream.Conn().LocalMultiaddr())
	if err != nil {
		log.Error("failed to convert local multiaddr to net addr:", err)
		return nil
	}
	return na
}
