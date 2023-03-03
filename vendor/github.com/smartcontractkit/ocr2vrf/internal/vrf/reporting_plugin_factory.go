package vrf

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	dkg_contract "github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	vrf_types "github.com/smartcontractkit/ocr2vrf/types"
)

type vrfReportingPluginFactory struct {
	l *localArgs
}

var _ types.ReportingPluginFactory = (*vrfReportingPluginFactory)(nil)

type localArgs struct {
	keyID              dkg_contract.KeyID
	coordinator        vrf_types.CoordinatorInterface
	keyProvider        KeyProvider
	serializer         vrf_types.ReportSerializer
	juelsPerFeeCoin    vrf_types.JuelsPerFeeCoin
	reasonableGasPrice vrf_types.ReasonableGasPrice
	period             uint16

	logger     commontypes.Logger
	randomness io.Reader
}

func (v *vrfReportingPluginFactory) NewReportingPlugin(
	c types.ReportingPluginConfig,
) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	if c.N > int(player_idx.MaxPlayer) {
		return nil, types.ReportingPluginInfo{},
			errors.Errorf("too many players: %d > %d", c.N, player_idx.MaxPlayer)
	}
	players, err := player_idx.PlayerIdxs(player_idx.Int(c.N))
	if err != nil {
		return nil, types.ReportingPluginInfo{},
			errors.Wrap(err, "could not determine local player DKG index")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	confDelays, err := v.l.coordinator.ConfirmationDelays(ctx)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, errors.Wrap(err, "could not get confirmation delays")
	}
	confDelaysSet := make(map[uint32]struct{})
	for _, d := range confDelays {
		confDelaysSet[d] = struct{}{}
	}

	err = v.l.coordinator.SetOffChainConfig(c.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, errors.Wrap(err, "could not set offchain config")
	}

	tbls, err := newSigRequest(
		v.l.keyID,
		v.l.keyProvider,
		player_idx.Int(c.N),
		player_idx.Int(c.F),
		common.Hash(c.ConfigDigest),
		*players[c.OracleID],
		&altbn_128.PairingSuite{},
		v.l.serializer,
		time.Hour,
		v.l.logger,
		v.l.juelsPerFeeCoin,
		v.l.reasonableGasPrice,
		v.l.coordinator,
		confDelaysSet,
		v.l.period,
	)
	if err != nil {
		return nil, types.ReportingPluginInfo{},
			errors.Wrap(err, "could not create new VRF Beacon reporting plugin")
	}
	return tbls, types.ReportingPluginInfo{
		Name: fmt.Sprintf("vrf instance %v", tbls.i),
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       200000,
			MaxObservationLength: 200000,
			MaxReportLength:      200000,
		},
	}, nil
}
