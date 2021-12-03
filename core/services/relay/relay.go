package relay

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
)

type Network string

var (
	Ethereum          Network = "ethereum"
	Solana            Network = "solana"
	SupportedRelayers         = map[Network]struct{}{
		Ethereum: {},
		Solana:   {},
	}
)

type Relayers map[Network]Relayer

type Relayer interface {
	service.Service
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2Provider, error)
}

type OCR2Provider interface {
	service.Service
	OffchainKeyring() types.OffchainKeyring
	OnchainKeyring() types.OnchainKeyring
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}

type Config struct {
	DB       *sqlx.DB
	Keystore keystore.Master
	ChainSet evm.ChainSet
	Lggr     logger.Logger
}
