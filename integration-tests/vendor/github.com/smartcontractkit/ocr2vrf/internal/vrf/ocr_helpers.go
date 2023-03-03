package vrf

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg"
	"github.com/smartcontractkit/ocr2vrf/internal/util"
	"github.com/smartcontractkit/ocr2vrf/internal/vrf/protobuf"
	vrf_types "github.com/smartcontractkit/ocr2vrf/types"

	kshare "go.dedis.ch/kyber/v3/share"

	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/libocr/commontypes"
)

func (s *sigRequest) ocrsSynced(ctx context.Context) error {
	deployedDKG, deployedVRF, err := s.coordinator.DKGVRFCommittees(ctx)
	if err != nil {
		return errors.Wrap(err, failedRetrieveOCRCommitteesMsg)
	}
	if len(deployedDKG.Signers) != len(deployedVRF.Signers) ||
		len(deployedDKG.Transmitters) != len(deployedVRF.Transmitters) {
		return errors.Errorf(
			committeesWithDifferentSizesMsg+" %s != %s", deployedDKG, deployedVRF,
		)
	}
	for i, s := range deployedDKG.Signers {
		if s != deployedVRF.Signers[i] {
			return errors.Errorf(
				signersMismatchMsg+" %s != %s", s, deployedVRF.Signers[i],
			)
		}
	}
	for i, s := range deployedDKG.Transmitters {
		if s != deployedVRF.Transmitters[i] {
			return errors.Errorf(
				transmittersMismatchMsg+" %s != %s", s, deployedVRF.Transmitters[i],
			)
		}
	}
	keyData := s.keyProvider.KeyLookup(s.keyID)
	if !keyData.Present {
		return errors.Errorf(noDistributedKeyMsg)
	}
	keyBytes, err := keyData.PublicKey.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, failedSerializeLocalKey)
	}
	onchainKeyHash, err := s.coordinator.ProvingKeyHash(ctx)
	if err != nil {
		return errors.Wrap(err, failedRetrieveOnchainKeyMsg)
	}
	localKeyHash := common.BytesToHash(crypto.Keccak256(keyBytes))
	if localKeyHash != onchainKeyHash {
		return errors.Errorf(incorrectPublicKeyMsg+" : 0x%x != 0x%x ", localKeyHash, onchainKeyHash)
	}
	if keyData.SecretShare == nil {
		return errors.Errorf(noLocalShareMsg)
	}
	return nil
}

func callbackHash(c *protobuf.CostedCallback) (common.Hash, error) {
	s, err := proto.Marshal(c)
	if err != nil {
		return common.Hash{}, errors.Wrap(
			err, "could not serialize callback for indexing",
		)
	}
	return crypto.Keccak256Hash(s), nil
}

func validateAndAddCallback(
	callbacks map[common.Hash]vrf_types.AbstractCostedCallbackRequest,
	c *protobuf.CostedCallback,
	oracle commontypes.OracleID, confDelays map[uint32]struct{}, beaconPeriod uint16,
	l commontypes.Logger,
) (common.Hash, error) {
	if err := sanityCheckCallback(c, l, oracle, confDelays, beaconPeriod); err != nil {
		return common.Hash{}, err
	}
	h, err := callbackHash(c)
	if err != nil {

		l.Error(
			"could not add callback",
			commontypes.LogFields{"err": err, "source": oracle, "callback": c},
		)
		return common.Hash{}, errors.Wrap(err, "could not add callback")
	}
	callbacks[h] = vrf_types.AbstractCostedCallbackRequest{
		BeaconHeight:      c.Callback.Height,
		ConfirmationDelay: c.Callback.ConfDelay,
		SubscriptionID:    big.NewInt(0).SetBytes(c.Callback.SubscriptionID),
		Price:             big.NewInt(0).SetBytes(c.Price),
		RequestID:         c.Callback.RequestId,
		NumWords:          uint16(c.Callback.NumWords),
		Requester:         common.BytesToAddress(c.Callback.Requester),
		Arguments:         c.Callback.Arguments,
		GasAllowance:      big.NewInt(0).SetBytes(c.GasAllowance),
		GasPrice:          big.NewInt(0).SetBytes(c.GasPrice),
		WeiPerUnitLink:    big.NewInt(0).SetBytes(c.WeiPerUnitLink),
	}
	return h, nil
}

func (s *sigRequest) storeCallbacksByBlocks(
	costedCallbacks []*protobuf.CostedCallback,
	callbacksByBlock map[heightDelay]map[common.Hash]struct{},
	callbackCounts map[common.Hash]uint64,
	callbacks map[common.Hash]vrf_types.AbstractCostedCallbackRequest,
	observer commontypes.OracleID,
) {

	seenCallbackHashes := make(
		map[common.Hash]struct{}, len(costedCallbacks),
	)
	for _, cb := range costedCallbacks {

		h, err := validateAndAddCallback(
			callbacks, cb, observer, s.confirmationDelays, s.period, s.logger,
		)
		if err != nil {

			continue
		}
		if _, present := seenCallbackHashes[h]; present {
			s.logger.Warn("duplicate callback received", commontypes.LogFields{
				"oracleID": observer, "duplicate callback": cb,
			})
			continue
		}
		seenCallbackHashes[h] = struct{}{}
		callbackCounts[h]++
		cbBlock := heightDelay{cb.Callback.Height, cb.Callback.ConfDelay}
		if _, present := callbacksByBlock[cbBlock]; !present {
			callbacksByBlock[cbBlock] = make(map[common.Hash]struct{})
		}
		callbacksByBlock[cbBlock][h] = struct{}{}
	}
}

func (s *sigRequest) parseAndStoreVRFProofs(
	proofs []*protobuf.VRFResponse,
	vrfContributions map[vrf_types.Block]map[commontypes.OracleID]kshare.PubShare,
	observer commontypes.OracleID,
	player *player_idx.PlayerIdx,
	kd dkg.KeyData,
) {

	pubShare := player.Index(kd.Shares).(kshare.PubShare)

	seenBlocks := make(map[heightDelay]struct{}, len(proofs))
	for _, output := range proofs {
		if _, present := s.confirmationDelays[output.Delay]; !present {
			s.logger.Warn(
				unknownConfirmationDelayInBlockMsg,
				commontypes.LogFields{
					"oracleID": observer, "delay": output.Delay,
					"known delays": s.confirmationDelays,
				})
			continue
		}
		if output.Height%uint64(s.period) != 0 {
			s.logger.Warn(
				nonBeaconHeightInBlockMsg,
				commontypes.LogFields{
					"oracleID": observer, "height": output.Height,
					"period": s.period,
				})
			continue
		}
		blockhash := common.BytesToHash(output.Blockhash)
		b := vrf_types.Block{output.Height, output.Delay, blockhash}
		hd := heightDelay{b.Height, b.ConfirmationDelay}
		if _, p := seenBlocks[hd]; p {
			s.logger.Warn(
				"multiple outputs requested for same block/delay pair",
				commontypes.LogFields{"oracleID": observer, "block": b})
			continue
		}
		seenBlocks[hd] = struct{}{}
		contribution := s.pairing.G1().Point()
		if err := contribution.UnmarshalBinary(output.Sig.Sig); err != nil {
			s.logger.Warn(failedReadContributionMsg, commontypes.LogFields{
				"oracleID": observer, "error": err,
				"contribution": fmt.Sprintf("0x%x", output.Sig.Sig),
			})
			continue
		}

		hashPoint := blsSeed(s.configDigest, b, kd.PublicKey)
		if _, present := vrfContributions[b]; !present {
			vrfContributions[b] = make(map[commontypes.OracleID]kshare.PubShare)
		}

		if !validateSignature(s.pairing, hashPoint, pubShare.V, contribution) {
			s.logger.Warn(wrongShare, commontypes.LogFields{
				"oracleID": observer, "sigShare": contribution,
				"keyShare": pubShare.V, "hashPoint": hashPoint,
				"pubKey": kd.PublicKey, "configDigest": s.configDigest, "block": b,
			})
			continue
		}
		vrfContributions[b][observer] = player.PubShare(contribution)
	}
}

func (s *sigRequest) aggregateOutputs(
	blocks vrf_types.Blocks,
	vrfContributions map[vrf_types.Block]map[commontypes.OracleID]kshare.PubShare,
	callbacksByBlock map[heightDelay]map[common.Hash]struct{},
	callbackCounts map[common.Hash]uint64,
	callbacks map[common.Hash]vrf_types.AbstractCostedCallbackRequest,
) (outputs []vrf_types.AbstractVRFOutput, err error) {
	outputs = make([]vrf_types.AbstractVRFOutput, 0, len(vrfContributions))
	for _, b := range blocks {
		hd := heightDelay{b.Height, b.ConfirmationDelay}

		if len(vrfContributions[b]) <= int(s.t) {
			s.logger.Debug(
				notEnoughContributions,
				commontypes.LogFields{
					"block": b, "num contributions": len(vrfContributions[b]),
				})

			continue
		}
		shares := make([]*kshare.PubShare, 0, len(vrfContributions[b]))
		player_indices, err := player_idx.PlayerIdxs(s.n)
		if err != nil {

			return nil, util.WrapError(
				err,
				"could not construct player indices for share reconstruction",
			)
		}
		for i, c := range vrfContributions[b] {
			pubShare := player_indices[i].PubShare(c.V)
			shares = append(shares, &pubShare)
		}
		kd := s.keyProvider.KeyLookup(s.keyID)

		output, err := kshare.RecoverCommit(
			s.pairing.G1(), shares, int(s.t)+1, len(shares),
		)
		if err != nil {

			s.logger.Error(
				"failed to recover distributed VRF output",
				commontypes.LogFields{"error": err, "shares": shares, "t": s.t},
			)

			continue
		}

		hpoint := blsSeed(s.configDigest, b, kd.PublicKey)

		if !validateSignature(s.pairing, hpoint, kd.PublicKey, output) {
			s.logger.Error(
				failedVerifyVRFOutput,
				commontypes.LogFields{"distributed signature": output},
			)

			continue
		}
		proof, err := output.MarshalBinary()
		if err != nil {

			s.logger.Error(
				"could not serialize VRF output for onchain transmission",
				commontypes.LogFields{"error": err, "output": output},
			)

		}
		ccallbacks := make(
			[]vrf_types.AbstractCostedCallbackRequest, 0, len(callbacksByBlock[hd]),
		)

		chashes := make([]string, 0, len(callbacksByBlock[hd]))
		for ch := range callbacksByBlock[hd] {
			chashes = append(chashes, ch.Hex())
		}
		sort.Strings(chashes)
		for _, chs := range chashes {
			ch := common.HexToHash(chs)

			if callbackCounts[ch] > uint64(s.t) {
				ccallbacks = append(ccallbacks, callbacks[ch])
			} else {
				s.logger.Error(
					notEnoughAppearancesCallback,
					commontypes.LogFields{"callback hash": ch, "t": s.t, "count": callbackCounts[ch]},
				)
			}
		}
		outputs = append(outputs, vrf_types.AbstractVRFOutput{
			b.Height,
			b.ConfirmationDelay,
			common.BytesToHash(proof),
			ccallbacks,
		})
	}
	return
}

func sanityCheckCallback(
	c *protobuf.CostedCallback, l commontypes.Logger, oracle commontypes.OracleID,
	confDelays map[uint32]struct{}, beaconPeriod uint16,
) error {
	if rem := c.Callback.Height % uint64(beaconPeriod); rem != 0 {
		l.Warn(nonBeaconHeightInCallbackMsg, commontypes.LogFields{
			"height": c.Callback.Height, "period": beaconPeriod, "remainder": rem})
		return errors.Errorf(
			nonBeaconHeightInCallbackMsg+" : %d ∤ %d", beaconPeriod, c.Callback.Height,
		)
	}
	if _, present := confDelays[c.Callback.ConfDelay]; !present {
		l.Warn(
			unknownConfirmationDelayMsg,
			commontypes.LogFields{
				"delay": c.Callback.ConfDelay, "good delays": confDelays,
				"source": oracle, "callback": c,
			},
		)
		return errors.Errorf(unknownConfirmationDelayMsg)
	}
	price := big.NewInt(0).SetBytes(c.Price)
	if price.Cmp(MaxPrice) > 0 {
		l.Warn(priceTooLargeMsg, commontypes.LogFields{
			"price": price, "max": MaxPrice, "callback": c, "source": oracle,
		})
		return errors.Errorf(priceTooLargeMsg)
	}
	if c.Callback.RequestId > MaxRequestID.Uint64() {
		l.Warn(requestIdTooLargeMsg, commontypes.LogFields{
			"requestID": c.Callback.RequestId, "max": maxUint48, "callback": c,
			"source": oracle,
		})
		return errors.Errorf(requestIdTooLargeMsg)
	}
	if uint64(c.Callback.NumWords) > MaxNumWords.Uint64() {

		l.Warn("numWords too large", commontypes.LogFields{
			"numWords": c.Callback.NumWords, "max": MaxNumWords, "callback": c,
			"source": oracle,
		})
		return errors.Errorf("numWords too large")
	}
	if len(c.Callback.Requester) > 20 {
		l.Warn("requester bytes too long to be address", commontypes.LogFields{
			"requester": c.Callback.Requester, "maxlen": 20, "source": oracle,
		})
		return errors.Errorf("requester bytes too long")
	}
	allowance := big.NewInt(0).SetBytes(c.GasAllowance)
	if allowance.Cmp(MaxGasAllowance) > 0 {
		l.Warn(
			excessGasAllowanceMsg,
			commontypes.LogFields{
				"allowance":     allowance,
				"max allowance": MaxGasAllowance,
				"source":        oracle,
			})
		return errors.Errorf(excessGasAllowanceMsg)
	}

	if len(c.Callback.Arguments) > MaxArgumentsLen {
		l.Warn(
			tooLongArgumentsMsg,
			commontypes.LogFields{
				"arguments length":     len(c.Callback.Arguments),
				"max arguments length": MaxArgumentsLen,
				"source":               oracle,
			})
		return errors.Errorf(tooLongArgumentsMsg)
	}
	return nil
}

func sortBigInt(l []*big.Int) []*big.Int {
	sort.Sort(byValue(l))
	return l
}

type byValue []*big.Int

func (a byValue) Len() int           { return len(a) }
func (a byValue) Less(i, j int) bool { return a[i].Cmp(a[j]) < 0 }
func (a byValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func medianBigInt(l []*big.Int) *big.Int {
	sortBigInt(l)
	midPoint := len(l) / 2
	if len(l)%2 == 1 {
		return l[midPoint]
	}
	if len(l) == 0 {

		panic("list must be populated")
	}

	midPointTotal := big.NewInt(0).Add(l[midPoint-1], l[midPoint])
	return midPointTotal.Div(midPointTotal, big.NewInt(2))
}

func callbacksEqual(c1, c2 vrf_types.AbstractCostedCallbackRequest) bool {
	return c1.BeaconHeight == c2.BeaconHeight &&
		c1.ConfirmationDelay == c2.ConfirmationDelay &&
		c1.SubscriptionID.Cmp(c2.SubscriptionID) == 0 &&
		c1.Price.Cmp(c2.Price) == 0 &&
		c1.RequestID == c2.RequestID &&
		c1.NumWords == c2.NumWords &&
		c1.Requester == c2.Requester &&
		bytes.Equal(c1.Arguments, c2.Arguments) &&
		c1.GasAllowance.Cmp(c2.GasAllowance) == 0 &&
		c1.GasPrice.Cmp(c2.GasPrice) == 0 &&
		c1.WeiPerUnitLink.Cmp(c2.WeiPerUnitLink) == 0
}

func getAbstractCallbackFromCallback(
	c *protobuf.CostedCallback,
) vrf_types.AbstractCostedCallbackRequest {
	return vrf_types.AbstractCostedCallbackRequest{
		BeaconHeight:      c.Callback.Height,
		ConfirmationDelay: c.Callback.ConfDelay,
		SubscriptionID:    big.NewInt(0).SetBytes(c.Callback.SubscriptionID),
		Price:             big.NewInt(0).SetBytes(c.Price),
		RequestID:         c.Callback.RequestId,
		NumWords:          uint16(c.Callback.NumWords),
		Requester:         common.BytesToAddress(c.Callback.Requester),
		Arguments:         c.Callback.Arguments,
		GasAllowance:      big.NewInt(0).SetBytes(c.GasAllowance),
		GasPrice:          big.NewInt(0).SetBytes(c.GasPrice),
		WeiPerUnitLink:    big.NewInt(0).SetBytes(c.WeiPerUnitLink),
	}
}

var (
	maxUint16 = big.NewInt(0).SetUint64(math.MaxUint16)
	maxUint24 = big.NewInt(0).SetBytes(bytes.Repeat([]byte{0xff}, 3))
	maxUint48 = big.NewInt(0).SetBytes(bytes.Repeat([]byte{0xff}, 6))
	maxUint32 = big.NewInt(0).SetUint64(math.MaxUint32)
	maxUint64 = big.NewInt(0).SetUint64(math.MaxUint64)
	maxUint96 = big.NewInt(0).SetBytes(bytes.Repeat([]byte{0xff}, 12))
)

var (
	MaxNumWords               = maxUint16
	MaxConfirmationDelay      = maxUint24
	MaxRequestID              = maxUint48
	MaxPrice                  = maxUint96
	MaxGasAllowance           = maxUint96
	MaxSubscriptionID         = maxUint64
	MaxArgumentsLen           = 62_500
	MaxBlocksInObservation    = 100
	MaxCallbacksInObservation = 100
)

func init() {

	if MaxNumWords.Cmp(maxUint32) > 0 {
		panic("MaxNumWords needs new backing type")
	}
	if MaxConfirmationDelay.Cmp(maxUint32) > 0 {
		panic("MaxConfirmationDelay needs new backing type")
	}
	if MaxRequestID.Cmp(maxUint64) > 0 {
		panic("MaxRequestID needs new backing type")
	}
	if MaxSubscriptionID.Cmp(maxUint64) > 0 {
		panic("MaxSubcriptionID needs new backing type")
	}
}

const (
	excessGasAllowanceMsg              = "gas allowance too large"
	unknownConfirmationDelayMsg        = "unknown confirmation delay"
	nonBeaconHeightInCallbackMsg       = "callback with non-beacon height"
	priceTooLargeMsg                   = "price too large"
	requestIdTooLargeMsg               = "requestID too large"
	noLocalShareMsg                    = "No local secret keyshare available"
	incorrectPublicKeyMsg              = "keyHash mismatch"
	noDistributedKeyMsg                = "no distributed key available"
	failedSerializeLocalKey            = "could not serialize local view of key"
	failedRetrieveOCRCommitteesMsg     = "failed to retrieve OCR committees"
	committeesWithDifferentSizesMsg    = "committee sizes differ"
	signersMismatchMsg                 = "committee signers differ"
	transmittersMismatchMsg            = "committee transmitters differ"
	failedRetrieveOnchainKeyMsg        = "could not retrieve onchain view of key hash"
	failedReadContributionMsg          = "could not read VRF contribution"
	nonBeaconHeightInBlockMsg          = "block output provided for non-beacon height"
	unknownConfirmationDelayInBlockMsg = "block output provided for unknown confirmation delay"
	tooLongArgumentsMsg                = "arguments field in callback is too long"
)
