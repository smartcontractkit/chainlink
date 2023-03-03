package tcpreuse

import (
	"context"
	"fmt"
	"math/rand"
	"net"

	"github.com/libp2p/go-netroute"
)

type multiDialer struct {
	listeningAddresses []*net.TCPAddr
	loopback           []*net.TCPAddr
	unspecified        []*net.TCPAddr
	fallback           net.TCPAddr
}

func (d *multiDialer) Dial(network, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

func randAddr(addrs []*net.TCPAddr) *net.TCPAddr {
	if len(addrs) > 0 {
		return addrs[rand.Intn(len(addrs))]
	}
	return nil
}

// DialContext dials a target addr.
// Dialing preference is
// * If there is a listener on the local interface the OS expects to use to route towards addr, use that.
// * If there is a listener on a loopback address, addr is loopback, use that.
// * If there is a listener on an undefined address (0.0.0.0 or ::), use that.
// * Use the fallback IP specified during construction, with a port that's already being listened on, if one exists.
func (d *multiDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	ip := tcpAddr.IP
	if !ip.IsLoopback() && !ip.IsGlobalUnicast() {
		return nil, fmt.Errorf("undialable IP: %s", ip)
	}

	if router, err := netroute.New(); err == nil {
		if _, _, preferredSrc, err := router.Route(ip); err == nil {
			for _, optAddr := range d.listeningAddresses {
				if optAddr.IP.Equal(preferredSrc) {
					return reuseDial(ctx, optAddr, network, addr)
				}
			}
		}
	}

	if ip.IsLoopback() && len(d.loopback) > 0 {
		return reuseDial(ctx, randAddr(d.loopback), network, addr)
	}
	if len(d.unspecified) == 0 {
		return reuseDial(ctx, &d.fallback, network, addr)
	}

	return reuseDial(ctx, randAddr(d.unspecified), network, addr)
}

func newMultiDialer(unspec net.IP, listeners map[*listener]struct{}) (m dialer) {
	addrs := make([]*net.TCPAddr, 0)
	loopback := make([]*net.TCPAddr, 0)
	unspecified := make([]*net.TCPAddr, 0)
	existingPort := 0

	for l := range listeners {
		addr := l.Addr().(*net.TCPAddr)
		addrs = append(addrs, addr)
		if addr.IP.IsLoopback() {
			loopback = append(loopback, addr)
		} else if addr.IP.IsGlobalUnicast() && existingPort == 0 {
			existingPort = addr.Port
		} else if addr.IP.IsUnspecified() {
			unspecified = append(unspecified, addr)
		}
	}
	m = &multiDialer{
		listeningAddresses: addrs,
		loopback:           loopback,
		unspecified:        unspecified,
		fallback:           net.TCPAddr{IP: unspec, Port: existingPort},
	}
	return
}
