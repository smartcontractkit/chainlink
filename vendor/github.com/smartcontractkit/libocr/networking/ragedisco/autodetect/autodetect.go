package autodetect

import (
	"fmt"
	"net"
	"net/netip"
	"sort"
)

type addrSet struct {
	public   []netip.Addr
	private  []netip.Addr
	loopback []netip.Addr
}

func (as *addrSet) add(ip netip.Addr) {
	switch {
	case ip.IsLoopback():
		as.loopback = append(as.loopback, ip)
	case ip.IsPrivate() || ip.IsLinkLocalUnicast(): // we need the second check for IPv6
		as.private = append(as.private, ip)
	default:
		as.public = append(as.public, ip)
	}
}

func sortIPs(ips []netip.Addr) {
	sort.Slice(ips, func(i, j int) bool {
		return ips[i].String() < ips[j].String()
	})
}

func (as *addrSet) all() []netip.Addr {
	// We sort each subset of IPs because the order in which they are returned from net.InterfaceAddrs can be
	// unpredictable.
	sortIPs(as.public)
	sortIPs(as.private)
	sortIPs(as.loopback)
	var ips []netip.Addr
	// We return IPs in order of decreasing reachability from the outside.
	ips = append(ips, as.public...)
	ips = append(ips, as.private...)
	ips = append(ips, as.loopback...)
	return ips
}

func parseIPFromNetAddr(netAddr string) (netip.Addr, error) {
	ipPrefix, err := netip.ParsePrefix(netAddr)
	if err != nil {
		// NOTE: This might be unnecessary but it seems like the type returned by net.InterfaceAddrs is
		// OS-dependent so we don't take any chances.
		return netip.ParseAddr(netAddr)
	} else {
		return ipPrefix.Addr(), nil
	}
}

func prioritizeNetAddrs(netAddrs []string) ([]netip.Addr, []netip.Addr, error) {
	var v4, v6 addrSet

	for _, netAddr := range netAddrs {
		ip, err := parseIPFromNetAddr(netAddr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse IP (%q): %w", netAddr, err)
		}
		if ip.Is4() {
			v4.add(ip)
		} else if ip.Is6() {
			v6.add(ip)
		} else {
			return nil, nil, fmt.Errorf("invalid ip: %s", ip)
		}
	}

	return v4.all(), v6.all(), nil
}

func AutodetectIPs() ([]netip.Addr, []netip.Addr, error) {
	netAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, nil, err
	}
	var addrs []string
	for _, netAddr := range netAddrs {
		addrs = append(addrs, netAddr.String())
	}
	return prioritizeNetAddrs(addrs)
}
