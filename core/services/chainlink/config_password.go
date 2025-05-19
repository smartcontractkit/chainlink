package chainlink

type passwordConfig struct {
	keystore func() string
	vrf      func() string
}

func (p *passwordConfig) Keystore() string { return p.keystore() }

func (p *passwordConfig) VRF() string { return p.vrf() }
