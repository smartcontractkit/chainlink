package vrf

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	dkg_contract "github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	vrf_types "github.com/smartcontractkit/ocr2vrf/types"
)

var _ types.ReportingPlugin = (*sigRequest)(nil)

type sigRequest struct {
	keyID       dkg_contract.KeyID
	keyProvider KeyProvider
	n           player_idx.Int

	t            player_idx.Int
	configDigest common.Hash
	i            player_idx.PlayerIdx
	pairing      pairing.Suite
	serializer   vrf_types.ReportSerializer
	blockProofs  map[vrf_types.Block]kyber.Point
	proofLock    sync.RWMutex

	logger commontypes.Logger

	retransmissionDelay time.Duration
	juelsPerFeeCoin     vrf_types.JuelsPerFeeCoin
	reasonableGasPrice  vrf_types.ReasonableGasPrice
	confirmationDelays  map[uint32]struct{}

	period      uint16
	coordinator vrf_types.CoordinatorInterface
	reports     map[types.ReportTimestamp]report
	reportsLock sync.RWMutex
}

func newSigRequest(
	keyID dkg_contract.KeyID,
	keyProvider KeyProvider,
	n player_idx.Int,
	t player_idx.Int,
	configDigest common.Hash,
	i player_idx.PlayerIdx,
	pairing pairing.Suite,
	serializer vrf_types.ReportSerializer,
	retransmissionDelay time.Duration,
	logger commontypes.Logger,
	juelsPerFeeCoin vrf_types.JuelsPerFeeCoin,
	reasonableGasPrice vrf_types.ReasonableGasPrice,
	coordinator vrf_types.CoordinatorInterface,
	confirmationDelays map[uint32]struct{},
	period uint16,
) (*sigRequest, error) {
	if n <= t {
		return nil, errors.Errorf(
			"committee size must be larger than the fault-tolerance threshold",
		)
	}
	return &sigRequest{
		keyID,
		keyProvider,
		n,
		t,
		configDigest,
		i,
		pairing,
		serializer,
		map[vrf_types.Block]kyber.Point{},
		sync.RWMutex{},
		logger,
		retransmissionDelay,
		juelsPerFeeCoin,
		reasonableGasPrice,
		confirmationDelays,
		period,
		coordinator,
		make(map[types.ReportTimestamp]report),
		sync.RWMutex{},
	}, nil
}

type report struct {
	r vrf_types.AbstractReport
	s []byte
}
