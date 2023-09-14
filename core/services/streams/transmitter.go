package streams

// TODO: llo transmitter

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
)

type Transmitter interface {
	commontypes.StreamsTransmitter
	services.Service
}

type transmitter struct {
	services.StateMachine
	lggr        logger.Logger
	rpcClient   wsrpc.Client
	fromAccount string
}

func NewTransmitter(lggr logger.Logger, rpcClient wsrpc.Client, fromAccount ed25519.PublicKey) Transmitter {
	return &transmitter{
		services.StateMachine{},
		lggr,
		rpcClient,
		fmt.Sprintf("%x", fromAccount),
	}
}

func (t *transmitter) Start(ctx context.Context) error {
	// TODO
	return nil
}

func (t *transmitter) Close() error {
	// TODO
	return nil
}

func (t *transmitter) HealthReport() map[string]error {
	report := map[string]error{t.Name(): t.Healthy()}
	services.CopyHealth(report, t.rpcClient.HealthReport())
	// FIXME
	// services.CopyHealth(report, t.queue.HealthReport())
	return report
}

func (t *transmitter) Name() string { return t.lggr.Name() }

func (t *transmitter) Transmit(
	context.Context,
	types.ConfigDigest,
	uint64,
	ocr3types.ReportWithInfo[commontypes.StreamsReportInfo],
	[]types.AttributedOnchainSignature,
) error {
	panic("TODO")
}

// FromAccount returns the stringified (hex) CSA public key
func (t *transmitter) FromAccount() (ocrtypes.Account, error) {
	return ocrtypes.Account(t.fromAccount), nil
}
