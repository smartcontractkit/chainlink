package chainlink

import (
	"net"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type p2p struct {
	c toml.P2P
}

func (p *p2p) Enabled() bool {
	return p.V1().Enabled() || p.V2().Enabled()
}

func (p *p2p) NetworkStack() (n ocrnetworking.NetworkingStack) {
	return p.c.NetworkStack()
}

func (p *p2p) PeerID() p2pkey.PeerID {
	return *p.c.PeerID
}

func (p *p2p) TraceLogging() bool {
	return *p.c.TraceLogging
}

func (p *p2p) IncomingMessageBufferSize() int {
	return int(*p.c.IncomingMessageBufferSize)
}

func (p *p2p) OutgoingMessageBufferSize() int {
	return int(*p.c.OutgoingMessageBufferSize)
}

func (p *p2p) V2() config.V2 {
	return &p2pv2{p.c.V2}
}

func (p *p2p) V1() config.V1 {
	return &p2pv1{p.c.V1}
}

type p2pv1 struct {
	c toml.P2PV1
}

func (v *p2pv1) Enabled() bool {
	return *v.c.Enabled
}

func (v *p2pv1) AnnounceIP() net.IP {
	return *v.c.AnnounceIP
}

func (v *p2pv1) AnnouncePort() uint16 {
	return *v.c.AnnouncePort
}

func (v *p2pv1) DefaultBootstrapPeers() ([]string, error) {
	p := *v.c.DefaultBootstrapPeers
	if p == nil {
		p = []string{}
	}
	return p, nil
}

func (v *p2pv1) DHTAnnouncementCounterUserPrefix() uint32 {
	return *v.c.DHTAnnouncementCounterUserPrefix
}

func (v *p2pv1) ListenIP() net.IP {
	return *v.c.ListenIP
}

func (v *p2pv1) ListenPort() uint16 {
	p := *v.c.ListenPort
	return p
}

func (v *p2pv1) NewStreamTimeout() time.Duration {
	return v.c.NewStreamTimeout.Duration()
}

func (v *p2pv1) BootstrapCheckInterval() time.Duration {
	return v.c.BootstrapCheckInterval.Duration()
}

func (v *p2pv1) DHTLookupInterval() int {
	return int(*v.c.DHTLookupInterval)
}

func (v *p2pv1) PeerstoreWriteInterval() time.Duration {
	return v.c.PeerstoreWriteInterval.Duration()
}

type p2pv2 struct {
	c toml.P2PV2
}

func (v *p2pv2) Enabled() bool {
	return *v.c.Enabled
}

func (v *p2pv2) AnnounceAddresses() []string {
	if a := v.c.AnnounceAddresses; a != nil {
		return *a
	}
	return nil
}

func (v *p2pv2) DefaultBootstrappers() (locators []commontypes.BootstrapperLocator) {
	if d := v.c.DefaultBootstrappers; d != nil {
		return *d
	}
	return nil
}

func (v *p2pv2) DeltaDial() models.Duration {
	if d := v.c.DeltaDial; d != nil {
		return *d
	}
	return models.Duration{}
}

func (v *p2pv2) DeltaReconcile() models.Duration {
	if d := v.c.DeltaReconcile; d != nil {
		return *d

	}
	return models.Duration{}
}

func (v *p2pv2) ListenAddresses() []string {
	if l := v.c.ListenAddresses; l != nil {
		return *l
	}
	return nil
}
