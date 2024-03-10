package ocrcommon

import ocrnetworking "github.com/smartcontractkit/libocr/networking"

func (p *SingletonPeerWrapper) PeerConfig() (ocrnetworking.PeerConfig, error) {
	return p.peerConfig()
}
