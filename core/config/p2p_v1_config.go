package config

import (
	"net"
	"time"
)

// P2PV1Networking is a subset of global config relevant to p2p v1 networking.
type P2PV1Networking interface {
	P2PAnnounceIP() net.IP
	P2PAnnouncePort() uint16
	P2PBootstrapPeers() ([]string, error)
	P2PDHTAnnouncementCounterUserPrefix() uint32
	P2PListenIP() net.IP
	P2PListenPort() uint16
	P2PListenPortRaw() string
	P2PNewStreamTimeout() time.Duration
	P2PBootstrapCheckInterval() time.Duration
	P2PDHTLookupInterval() int
	P2PPeerstoreWriteInterval() time.Duration
}
