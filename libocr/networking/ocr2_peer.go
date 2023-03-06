package networking

import (
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type ocr2BinaryNetworkEndpointFactory struct {
	*concretePeer
}

type ocr2BootstrapperFactory struct {
	*concretePeer
}

func (o *ocr2BinaryNetworkEndpointFactory) NewEndpoint(
	configDigest types.ConfigDigest,
	pids []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
	limits types.BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if !o.networkingStack.needsv2() {
		return nil, fmt.Errorf("OCR2 requires v2 networking, but current NetworkingStack is %v", o.networkingStack)
	}
	return o.newEndpoint(
		NetworkingStackV2,
		configDigest,
		pids,
		nil,
		v2bootstrappers,
		f,
		BinaryNetworkEndpointLimits(limits),
	)
}

func (o *ocr2BootstrapperFactory) NewBootstrapper(
	configDigest types.ConfigDigest,
	peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
) (commontypes.Bootstrapper, error) {
	if !o.networkingStack.needsv2() {
		return nil, fmt.Errorf("OCR2 requires v2 networking, but current NetworkingStack is %v", o.networkingStack)
	}
	return o.newBootstrapper(
		NetworkingStackV2,
		configDigest,
		peerIDs,
		nil,
		v2bootstrappers,
		f,
	)
}
