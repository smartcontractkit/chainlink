package chainlink

import (
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/libocr/commontypes"
)

type p2p struct {
	c toml.P2P
}

func (p *p2p) Enabled() bool {
	return p.V2().Enabled()
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

func (v *p2pv2) DeltaDial() commonconfig.Duration {
	if d := v.c.DeltaDial; d != nil {
		return *d
	}
	return commonconfig.Duration{}
}

func (v *p2pv2) DeltaReconcile() commonconfig.Duration {
	if d := v.c.DeltaReconcile; d != nil {
		return *d
	}
	return commonconfig.Duration{}
}

func (v *p2pv2) ListenAddresses() []string {
	if l := v.c.ListenAddresses; l != nil {
		return *l
	}
	return nil
}
