package ragedisco

import (
	"fmt"
	"net/netip"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking/ragedisco/autodetect"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type autodetector func() ([]netip.Addr, []netip.Addr, error)

// combinedAnnounceAddrs takes in the user-provided string addresses and converts them to proper ragep2p addresses.
// Unspecified addresses such as 0.0.0.0 are replaced with specified auto-detected addresses. This function never
// returns duplicate addresses. If autodetection fails and there is no specified (non-unspecified) address available,
// it returns returns ok=false.
func combinedAnnounceAddrs(logger commontypes.Logger, addrStrs []string, autodetectFunc autodetector) ([]ragetypes.Address, bool) {
	var addrs []ragetypes.Address
	ifaceV4, ifaceV6, autodetectErr := autodetectFunc()
	for _, addrStr := range addrStrs {
		addrPort, err := parseAddrPortForAnnouncement(addrStr)
		if err != nil {
			logger.Critical("Invalid announce address provided", commontypes.LogFields{"address": addrStr, "error": err})
			return nil, false
		}

		ip, port := addrPort.Addr(), addrPort.Port()
		if !ip.IsUnspecified() {
			addrs = append(addrs, ragetypes.Address(addrStr))
		} else {
			if ip.Is4() {
				if autodetectErr == nil {
					addrs = append(addrs, joinIPsPort(ifaceV4, port)...)
				}
			} else if ip.Is6() {
				if autodetectErr == nil {
					addrs = append(addrs, joinIPsPort(ifaceV6, port)...)
				}
			} else {
				logger.Critical("We ended up with an announce IP that is neither IPv4 nor IPv6. This should never happen!", commontypes.LogFields{"ip": ip})
				return nil, false
			}
		}
	}

	addrs = dedup(addrs)
	if autodetectErr != nil {
		if len(addrs) > 0 {
			logger.Critical("Could not autodetect announce addresses, using only specified addresses", commontypes.LogFields{
				"announceAddresses": addrs,
				"error":             autodetectErr,
			})
		} else {
			logger.Critical("No specified announce addresses were supplied and failed to autodetect interface IPs", commontypes.LogFields{
				"error": autodetectErr,
			})
			return nil, false
		}
	}
	if len(addrs) > maxAddrsInAnnouncement {
		logger.Critical("Announce addresses length is more than the allowed max, trimming", commontypes.LogFields{
			"length":           len(addrs),
			"maxAllowedLength": maxAddrsInAnnouncement,
		})
		addrs = addrs[:maxAddrsInAnnouncement]
	}

	// Sanity check, better to fail here than to produce announcements that
	// would be not accepted by other peers due to invalid addresses.
	for _, addr := range addrs {
		if !isValidForAnnouncement(addr) {
			logger.Critical("Produced announce addresses contain an invalid address. This should never happen!", commontypes.LogFields{
				"invalidAddress":          addr,
				"announceAddresses":       addrs,
				"configAnnounceAddresses": addrStrs,
			})
			return nil, false
		}
	}
	return addrs, true
}

func combinedAnnounceAddrsForDiscoverer(logger commontypes.Logger, addrStrs []string) ([]ragetypes.Address, bool) {
	return combinedAnnounceAddrs(logger, addrStrs, autodetect.AutodetectIPs)
}

func dedup(addrs []ragetypes.Address) []ragetypes.Address {
	m := make(map[ragetypes.Address]struct{})
	var ret []ragetypes.Address
	for _, addr := range addrs {
		if _, exists := m[addr]; exists {
			continue
		}
		ret = append(ret, addr)
		m[addr] = struct{}{}
	}
	return ret
}

const maxAddrPortValidForAnnouncementSize = len("[0000:0000:0000:0000:0000:ffff:255.255.255.255]:65535")

// Decoupled to aid in fuzzing, to ensure that really
// maxAddrPortValidForAnnouncementSize is the correct size limit.
func parseAddrPortForAnnouncementNoSizeLimit(s string) (netip.AddrPort, error) {
	addrPort, err := netip.ParseAddrPort(s)
	if err != nil {
		return netip.AddrPort{}, err
	}
	addr := addrPort.Addr()
	if addr.Zone() != "" {
		return netip.AddrPort{}, fmt.Errorf("address %q contains IPv6 zone", s)
	}
	if !(addr.Is4() || addr.Is6()) {
		return netip.AddrPort{}, fmt.Errorf("address %q should be either IPv4 or IPv6", s)
	}
	return addrPort, err
}

func parseAddrPortForAnnouncement(s string) (netip.AddrPort, error) {
	if len(s) > maxAddrPortValidForAnnouncementSize {
		return netip.AddrPort{}, fmt.Errorf("address %q larger than %d bytes", s, maxAddrPortValidForAnnouncementSize)
	}
	return parseAddrPortForAnnouncementNoSizeLimit(s)
}

// isValidForAnnouncement checks that the provided address is in the form ip:port.
// Hostnames or domain names are not allowed.
func isValidForAnnouncement(a ragetypes.Address) bool {
	_, err := parseAddrPortForAnnouncement(string(a))
	return err == nil
}

func joinIPPort(ip netip.Addr, port uint16) ragetypes.Address {
	return ragetypes.Address(netip.AddrPortFrom(ip, port).String())
}

func joinIPsPort(ips []netip.Addr, port uint16) []ragetypes.Address {
	var addrs []ragetypes.Address
	for _, ip := range ips {
		addrs = append(addrs, joinIPPort(ip, port))
	}
	return addrs
}
