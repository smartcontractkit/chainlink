package pipeline

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf/proof"
	"go.uber.org/multierr"
)

type VRFTask struct {
	BaseTask           `mapstructure:",squash"`
	PublicKey          string `json:"publicKey"`
	RequestBlockHash   string `json:"requestBlockHash"`
	RequestBlockNumber string `json:"requestBlockNumber"`
	Topics             string `json:"topics"`
	ProofGenerator     string `json:"proofGenerator"`
}

var _ Task = (*VRFTask)(nil)

func (t *VRFTask) Type() TaskType {
	return TaskTypeVRF
}

func (t *VRFTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	// TODO: clean up error handling, fix all these hard asserts
	if len(inputs) != 1 {
		return Result{Error: errors.New("invalid inputs")}
	}
	if inputs[0].Error != nil {
		return Result{Error: errors.New("input errored")}
	}
	logValues, ok := inputs[0].Value.(map[string]interface{})
	if !ok {
		return Result{Error: errors.New("expected map input")}
	}
	var (
		pubKey             BytesParam
		requestBlockHash   BytesParam
		requestBlockNumber Uint64Param
		topics             HashSliceParam
		proofGenerator     FunctionParam
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&pubKey, From(VarExpr(t.PublicKey, vars))), "publicKey"),
		errors.Wrap(ResolveParam(&requestBlockHash, From(VarExpr(t.RequestBlockHash, vars))), "requestBlockHash"),
		errors.Wrap(ResolveParam(&requestBlockNumber, From(VarExpr(t.RequestBlockNumber, vars))), "requestBlockNumber"),
		errors.Wrap(ResolveParam(&topics, From(VarExpr(t.Topics, vars))), "topics"),
		errors.Wrap(ResolveParam(&proofGenerator, From(VarExpr(t.ProofGenerator, vars))), "proofGenerator"),
	)
	if err != nil {
		return Result{Error: err}
	}

	requestKeyHash := logValues["keyHash"].([32]byte)
	requestPreSeed := logValues["seed"].(*big.Int)
	requestJobID := logValues["jobID"].([32]byte)
	var pk secp256k1.PublicKey
	copy(pk[:], pubKey[:])
	pkh := pk.MustHash()
	// Validate the key against the spec
	if !bytes.Equal(requestKeyHash[:], pkh[:]) {
		return Result{Error: fmt.Errorf("invalid key hash %v expected %v", hex.EncodeToString(requestKeyHash[:]), hex.EncodeToString(pkh[:]))}
	}
	preSeed, err := proof.BigToSeed(requestPreSeed)
	if err != nil {
		return Result{Error: errors.New("unable to parse preseed")}
	}
	if !bytes.Equal(topics[0][:], requestJobID[:]) && !bytes.Equal(topics[1][:], requestJobID[:]) {
		return Result{Error: fmt.Errorf("request jobID %v doesn't match expected %v or %v", requestJobID[:], topics[0][:], topics[1][:])}
	}
	seed := proof.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: common.BytesToHash(requestBlockHash),
		BlockNum:  uint64(requestBlockNumber),
	}
	solidityProofIntf, err := proofGenerator([]interface{}{seed})
	if err != nil {
		return Result{Error: err}
	}
	resp := solidityProofIntf.(proof.MarshaledOnChainResponse)
	return Result{Value: resp[:]}
}
