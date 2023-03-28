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

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
)

type VRFTask struct {
	BaseTask           `mapstructure:",squash"`
	PublicKey          string `json:"publicKey"`
	RequestBlockHash   string `json:"requestBlockHash"`
	RequestBlockNumber string `json:"requestBlockNumber"`
	Topics             string `json:"topics"`

	keyStore VRFKeyStore
}

type VRFKeyStore interface {
	GenerateProof(id string, seed *big.Int) (vrfkey.Proof, error)
}

var _ Task = (*VRFTask)(nil)

func (t *VRFTask) Type() TaskType {
	return TaskTypeVRF
}

func (t *VRFTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
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
	requestPreSeed, ok := logValues["seed"].(*big.Int)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid preSeed")}, runInfo
	}
	requestJobID, ok := logValues["jobID"].([32]byte)
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "invalid requestJobID")}, runInfo
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
	if !bytes.Equal(topics[0][:], requestJobID[:]) && !bytes.Equal(topics[1][:], requestJobID[:]) {
		return Result{Error: fmt.Errorf("request jobID %v doesn't match expected %v or %v", requestJobID[:], topics[0][:], topics[1][:])}, runInfo
	}
	if len(requestBlockHash) != common.HashLength {
		return Result{Error: fmt.Errorf("invalid BlockHash length %d expected %d", len(requestBlockHash), common.HashLength)}, runInfo
	}
	preSeedData := proof.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: common.BytesToHash(requestBlockHash),
		BlockNum:  uint64(requestBlockNumber),
	}
	finalSeed := proof.FinalSeed(preSeedData)
	p, err := t.keyStore.GenerateProof(pk.String(), finalSeed)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	onChainProof, err := proof.GenerateProofResponseFromProof(p, preSeedData)
	if err != nil {
		return Result{Error: err}, retryableRunInfo()
	}
	var results = make(map[string]interface{})
	results["onChainProof"] = hexutil.Encode(onChainProof[:])

	return Result{Value: hexutil.Encode(onChainProof[:])}, runInfo
}
