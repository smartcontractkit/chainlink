package networking

import (
	"fmt"
	"net"

	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

func makeAddrsFactory(ip net.IP, port uint16) (bhost.AddrsFactory, error) {
	if ip != nil && port != 0 {
		addr, err := makeMultiaddr(ip, port)
		if err != nil {
			return nil, errors.Wrap(err, "could not make multiaddr")
		}
		var addrs []ma.Multiaddr
		addrs = append(addrs, addr)
		return func([]ma.Multiaddr) []ma.Multiaddr {
			return addrs
		}, nil
	} else if ip != nil || port != 0 {
		return nil, errors.New("ip and port must both be specified, or both left unspecified")
	} else {
		return bhost.DefaultAddrsFactory, nil
	}
}

func makeMultiaddr(ip net.IP, port uint16) (ma.Multiaddr, error) {
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.Errorf("listen address must be a valid ipv4 address, got: %s", ip.String())
	}
	if port == 0 {
		return nil, errors.New("port must not be zero")
	}
	s := fmt.Sprintf("/ip4/%s/tcp/%d", ip4.String(), port)
	return ma.NewMultiaddr(s)
}
