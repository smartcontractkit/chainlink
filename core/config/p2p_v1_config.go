package config

import (
	"net"
	"time"
)

type V1 interface {
	Enabled() bool
	AnnounceIP() net.IP
	AnnouncePort() uint16
	DefaultBootstrapPeers() ([]string, error)
	DHTAnnouncementCounterUserPrefix() uint32
	ListenIP() net.IP
	ListenPort() uint16
	NewStreamTimeout() time.Duration
	BootstrapCheckInterval() time.Duration
	DHTLookupInterval() int
	PeerstoreWriteInterval() time.Duration
}
