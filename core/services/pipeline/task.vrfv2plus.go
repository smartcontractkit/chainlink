package pipeline

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
)

var (
	vrfCoordinatorV2PlusABI = evmtypes.MustGetABI(vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalABI)
)

// VRFTaskV2Plus is identical to VRFTaskV2 except that it uses the V2Plus VRF
// request commitment, which includes a boolean indicating whether native or
// link payment was used.
type VRFTaskV2Plus struct {
	BaseTask           `mapstructure:",squash"`
	PublicKey          string `json:"publicKey"`
	RequestBlockHash   string `json:"requestBlockHash"`
	RequestBlockNumber string `json:"requestBlockNumber"`
	Topics             string `json:"topics"`

	keyStore VRFKeyStore
}

var _ Task = (*VRFTaskV2Plus)(nil)

func (t *VRFTaskV2Plus) Type() TaskType {
	return TaskTypeVRFV2Plus
}

func (t *VRFTaskV2Plus) Run(_ context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	if len(inputs) != 1 {
		return Result{Error: ErrWrongInputCardinality}, runInfo
	}
	if inputs[0].Error != nil {
		return Result{Error: ErrInputTaskErrored}, runInfo
	}
	logValues, ok := inputs[0].Value.(map[string]interface{})
	if !ok {
		return Result{Error: errors.Wrap(ErrBadInput, "expected map input")}, runInfo
	}
	var (
		pubKey             BytesParam
		requestBlockHash   BytesParam
		requestBlockNumber Uint64Param
		topics             HashSliceParam
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&pubKey, From(VarExpr(t.PublicKey, vars))), "publicKey"),
		errors.Wrap(ResolveParam(&requestBlockHash, From(VarExpr(t.RequestBlockHash, vars))), "requestBlockHash"),
		errors.Wrap(ResolveParam(&requestBlockNumber, From(VarExpr(t.RequestBlockNumber, vars))), "requestBlockNumber"),
		errors.Wrap(ResolveParam(&topics, From(VarExpr(t.Topics, vars))), "topics"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	requestKeyHash, ok := logValues["keyHash"].([32]byte)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid keyHash")}, runInfo
	}
	requestPreSeed, ok := logValues["preSeed"].(*big.Int)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid preSeed")}, runInfo
	}
	requestId, ok := logValues["requestId"].(*big.Int)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid requestId")}, runInfo
	}
	subID, ok := logValues["subId"].(*big.Int)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid subId")}, runInfo
	}
	callbackGasLimit, ok := logValues["callbackGasLimit"].(uint32)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid callbackGasLimit")}, runInfo
	}
	numWords, ok := logValues["numWords"].(uint32)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid numWords")}, runInfo
	}
	sender, ok := logValues["sender"].(common.Address)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid sender")}, runInfo
	}
	extraArgs, ok := logValues["extraArgs"].([]byte)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid extraArgs")}, runInfo
	}
	pk, err := secp256k1.NewPublicKeyFromBytes(pubKey)
	if err != nil {
		return Result{Error: fmt.Errorf("failed to create PublicKey from bytes %v", err)}, runInfo
	}
	pkh := pk.MustHash()
	// Validate the key against the spec
	if !bytes.Equal(requestKeyHash[:], pkh[:]) {
		return Result{Error: fmt.Errorf("invalid key hash %v expected %v", hex.EncodeToString(requestKeyHash[:]), hex.EncodeToString(pkh[:]))}, runInfo
	}
	preSeed, err := proof.BigToSeed(requestPreSeed)
	if err != nil {
		return Result{Error: fmt.Errorf("unable to parse preseed %v", preSeed)}, runInfo
	}
	if len(requestBlockHash) != common.HashLength {
		return Result{Error: fmt.Errorf("invalid BlockHash length %d expected %d", len(requestBlockHash), common.HashLength)}, runInfo
	}
	preSeedData := proof.PreSeedDataV2Plus{
		PreSeed:          preSeed,
		BlockHash:        common.BytesToHash(requestBlockHash),
		BlockNum:         uint64(requestBlockNumber),
		SubId:            subID,
		CallbackGasLimit: callbackGasLimit,
		NumWords:         numWords,
		Sender:           sender,
		ExtraArgs:        extraArgs,
	}
	finalSeed := proof.FinalSeedV2Plus(preSeedData)
	id := hexutil.Encode(pk[:])
	p, err := t.keyStore.GenerateProof(id, finalSeed)
	if err != nil {
		return Result{Error: err}, retryableRunInfo()
	}
	onChainProof, rc, err := proof.GenerateProofResponseFromProofV2Plus(p, preSeedData)
	if err != nil {
		return Result{Error: err}, retryableRunInfo()
	}
	// onlyPremium is false because this task assumes that chainlink node fulfills the VRF request
	// gas cost should be billed to the requesting subscription
	b, err := vrfCoordinatorV2PlusABI.Pack("fulfillRandomWords", onChainProof, rc, false /* onlyPremium */)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	results := make(map[string]interface{})
	output := hexutil.Encode(b)
	results["output"] = output
	// RequestID needs to be a [32]byte for EvmTxMeta.
	results["requestID"] = hexutil.Encode(requestId.Bytes())

	// store vrf proof and request commitment separately so they can be used in a batch fashion
	results["proof"] = onChainProof
	results["requestCommitment"] = rc

	lggr.Debugw("Completed VRF V2 task run", "reqID", requestId.String(), "output", output)

	return Result{Value: results}, runInfo
}
