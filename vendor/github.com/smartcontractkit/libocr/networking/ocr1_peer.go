package networking

import (
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/configdigesthelper"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type ocr1BinaryNetworkEndpointFactory struct {
	*concretePeer
}

var _ ocr1types.BinaryNetworkEndpointFactory = (*ocr1BinaryNetworkEndpointFactory)(nil)

const (
	// MaxOCRMsgLength is the maximum allowed length for a data payload in bytes
	// This is exported as serialization tests depend on it.
	// NOTE: This is slightly larger than 2x of the largest message we can
	// possibly send, assuming N=31.
	MaxOCRMsgLength = 10000
)

func (o *ocr1BinaryNetworkEndpointFactory) NewEndpoint(
	configDigest ocr1types.ConfigDigest,
	pids []string,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
	messagesRatePerOracle float64,
	messagesCapacityPerOracle int,
) (commontypes.BinaryNetworkEndpoint, error) {
	return o.newEndpoint(
		o.networkingStack,
		configdigesthelper.OCR1ToOCR2(configDigest),
		pids,
		v1bootstrappers,
		v2bootstrappers,
		f,
		BinaryNetworkEndpointLimits{
			MaxOCRMsgLength,
			messagesRatePerOracle,
			messagesCapacityPerOracle,
			messagesRatePerOracle * MaxOCRMsgLength,
			messagesCapacityPerOracle * MaxOCRMsgLength,
		},
	)
}

type ocr1BootstrapperFactory struct {
	*concretePeer
}

func (o *ocr1BootstrapperFactory) NewBootstrapper(
	configDigest ocr1types.ConfigDigest,
	peerIDs []string,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
) (commontypes.Bootstrapper, error) {
	return o.newBootstrapper(
		o.networkingStack,
		configdigesthelper.OCR1ToOCR2(configDigest),
		peerIDs,
		v1bootstrappers,
		v2bootstrappers,
		f,
	)
}
