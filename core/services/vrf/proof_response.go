package vrf

// Contains logic/data for mandatorily mixing VRF seeds with the hash of the
// block in which a VRF request appeared

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// ProofResponse is the data which is sent back to the VRFCoordinator, so that
// it can verify that the seed the oracle finally used is correct.
type ProofResponse struct {
	// Approximately the proof which will be checked on-chain. Note that this
	// contains the pre-seed in place of the final seed. That should be computed
	// as in FinalSeed.
	P        Proof
	PreSeed  Seed   // Seed received during VRF request
	BlockNum uint64 // Height of the block in which tihs request was made
	// V2 Only fields
	SubId            uint64
	CallbackGasLimit uint64
	NumWords         uint64
	Sender           common.Address
}

// OnChainResponseLength is the length of the MarshaledOnChainResponse. The
// extra 32 bytes are for blocknumber (as a uint256), which goes at the end. The
// seed is rewritten with the preSeed. (See MarshalForVRFCoordinator and
// ProofResponse#ActualProof.)
const OnChainResponseLength = ProofLength +
	32 + //blocknum
	32 + //subID
	32 + //gaslimit
	32 + //numWords
	32 //sender

// MarshaledOnChainResponse is the flat bytes which are sent back to the
// VRFCoordinator.
type MarshaledOnChainResponse [OnChainResponseLength]byte

// MarshalForVRFCoordinator constructs the flat bytes which are sent to the
// VRFCoordinator.
func (p *ProofResponse) MarshalForVRFCoordinator() (
	response MarshaledOnChainResponse, err error) {
	solidityProof, err := p.P.SolidityPrecalculations()
	if err != nil {
		return MarshaledOnChainResponse{}, errors.Wrap(err,
			"while marshaling proof for VRFCoordinator")
	}
	// Overwrite seed input to the VRF proof generator with the seed the
	// VRFCoordinator originally requested, so that it can identify the request
	// corresponding to this response, and compute the final seed itself using the
	// blockhash it infers from the block number.
	solidityProof.P.Seed = common.BytesToHash(p.PreSeed[:]).Big()
	mProof := solidityProof.MarshalForSolidityVerifier()
	wireBlockNum := utils.EVMWordUint64(p.BlockNum)
	subId := utils.EVMWordUint64(p.SubId)
	callbackLimit := utils.EVMWordUint64(p.CallbackGasLimit)
	numWords := utils.EVMWordUint64(p.NumWords)
	sender := utils.EVMWordAddress(p.Sender)
	fmt.Println("sender", p.Sender.String(), hex.EncodeToString(sender))
	rl := copy(response[:], append(append(append(append(append(mProof[:], wireBlockNum...), subId...), callbackLimit...), numWords...), sender...))
	if rl != OnChainResponseLength {
		return MarshaledOnChainResponse{}, errors.Errorf(
			"wrong length for response to VRFCoordinator")
	}
	return response, nil
}

// UnmarshalProofResponse returns the ProofResponse represented by the bytes in m
func UnmarshalProofResponse(m MarshaledOnChainResponse) (*ProofResponse, error) {
	blockNum := common.BytesToHash(m[ProofLength : ProofLength+32]).Big().Uint64()
	proof, err := UnmarshalSolidityProof(m[:ProofLength])
	if err != nil {
		return nil, errors.Wrap(err, "while parsing ProofResponse")
	}
	preSeed, err := BigToSeed(proof.Seed)
	if err != nil {
		return nil, errors.Wrap(err, "while converting seed to bytes representation")
	}
	return &ProofResponse{P: proof, PreSeed: preSeed, BlockNum: blockNum}, nil
}

// CryptoProof returns the proof implied by p, with the correct seed
func (p ProofResponse) CryptoProof(s PreSeedData) (Proof, error) {
	proof := p.P // Copy P, which has wrong seed value
	proof.Seed = FinalSeed(s)
	valid, err := proof.VerifyVRFProof()
	if err != nil {
		return Proof{}, errors.Wrap(err,
			"could not validate proof implied by on-chain response")
	}
	if !valid {
		return Proof{}, errors.Errorf(
			"proof implied by on-chain response is invalid")
	}
	return proof, nil
}

// GenerateProofResponse returns the marshaled proof of the VRF output given the
// secretKey and the seed computed from the s.PreSeed and the s.BlockHash
func GenerateProofResponse(secretKey common.Hash, s PreSeedData) (
	MarshaledOnChainResponse, error) {
	seed := FinalSeed(s)
	proof, err := GenerateProof(secretKey, common.BigToHash(seed))
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	p := ProofResponse{P: proof, PreSeed: s.PreSeed, BlockNum: s.BlockNum, SubId: s.SubId, CallbackGasLimit: s.CallbackGasLimit, NumWords: s.NumWords, Sender: s.Sender}
	rv, err := p.MarshalForVRFCoordinator()
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return rv, nil
}
