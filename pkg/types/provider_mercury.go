package types

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

// MercuryProvider provides components needed for a mercury OCR2 plugin.
// Mercury requires config tracking but does not transmit on-chain.
type MercuryProvider interface {
	PluginProvider

	ReportCodecV1() v1.ReportCodec
	ReportCodecV2() v2.ReportCodec
	ReportCodecV3() v3.ReportCodec
	OnchainConfigCodec() mercury.OnchainConfigCodec
	MercuryServerFetcher() mercury.ServerFetcher
	MercuryChainReader() mercury.ChainReader
}

type PluginMercury interface {
	NewMercuryV1Factory(ctx context.Context, provider MercuryProvider, dataSource v1.DataSource) (MercuryPluginFactory, error)
	NewMercuryV2Factory(ctx context.Context, provider MercuryProvider, dataSource v2.DataSource) (MercuryPluginFactory, error)
	NewMercuryV3Factory(ctx context.Context, provider MercuryProvider, dataSource v3.DataSource) (MercuryPluginFactory, error)
}

type MercuryPluginFactory interface {
	Service
	ocr3types.MercuryPluginFactory
}
