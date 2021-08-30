package proof

// Contains logic/data for mandatorily mixing VRF seeds with the hash of the
// block in which a VRF request appeared

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ProofResponse is the data which is sent back to the VRFCoordinator, so that
// it can verify that the seed the oracle finally used is correct.
type ProofResponse struct {
	// Approximately the proof which will be checked on-chain. Note that this
	// contains the pre-seed in place of the final seed. That should be computed
	// as in FinalSeed.
	P        vrfkey.Proof
	PreSeed  Seed   // Seed received during VRF request
	BlockNum uint64 // Height of the block in which this request was made
}

type ProofResponseV2 struct {
	P        vrfkey.Proof
	PreSeed  Seed
	BlockNum uint64
	// V2 Only fields
	SubId            uint64         // Subscription ID to be charged for fulfillment
	CallbackGasLimit uint32         // Gas limit for consumer callback
	NumWords         uint32         // Number of random words to expand to
	Sender           common.Address // VRF consumer address
}

// OnChainResponseLength is the length of the MarshaledOnChainResponse. The
// extra 32 bytes are for blocknumber (as a uint256), which goes at the end. The
// seed is rewritten with the preSeed. (See MarshalForVRFCoordinator and
// ProofResponse#ActualProof.)
const OnChainResponseLength = ProofLength +
	32 // blocknum

const OnChainResponseLengthV2 = ProofLength +
	32 + // blocknum
	32 + // subID
	32 + // gaslimit
	32 + // numWords
	32 // sender

// MarshaledOnChainResponse is the flat bytes which are sent back to the
// VRFCoordinator.
type MarshaledOnChainResponse [OnChainResponseLength]byte
type MarshaledOnChainResponseV2 [OnChainResponseLengthV2]byte

// MarshalForVRFCoordinator constructs the flat bytes which are sent to the
// VRFCoordinator.
func (p *ProofResponse) MarshalForVRFCoordinator() (
	response MarshaledOnChainResponse, err error) {
	solidityProof, err := SolidityPrecalculations(&p.P)
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
	rl := copy(response[:], append(mProof[:], wireBlockNum...))
	if rl != OnChainResponseLength {
		return MarshaledOnChainResponse{}, errors.Errorf(
			"wrong length for response to VRFCoordinator")
	}
	return response, nil
}

func (p *ProofResponseV2) MarshalForVRFCoordinator() (
	response MarshaledOnChainResponseV2, err error) {
	solidityProof, err := SolidityPrecalculations(&p.P)
	if err != nil {
		return MarshaledOnChainResponseV2{}, errors.Wrap(err,
			"while marshaling proof for VRFCoordinatorV2")
	}
	solidityProof.P.Seed = common.BytesToHash(p.PreSeed[:]).Big()
	mProof := solidityProof.MarshalForSolidityVerifier()
	wireBlockNum := utils.EVMWordUint64(p.BlockNum)
	subId := utils.EVMWordUint64(p.SubId)
	callbackLimit := utils.EVMWordUint32(p.CallbackGasLimit)
	numWords := utils.EVMWordUint32(p.NumWords)
	sender := utils.EVMWordAddress(p.Sender)

	rl := copy(response[:], bytes.Join([][]byte{mProof[:], wireBlockNum, subId, callbackLimit, numWords, sender}, []byte{}))
	if rl != OnChainResponseLengthV2 {
		return MarshaledOnChainResponseV2{}, errors.Errorf(
			"wrong length for response to VRFCoordinatorV2")
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
func (p ProofResponse) CryptoProof(s PreSeedData) (vrfkey.Proof, error) {
	proof := p.P // Copy P, which has wrong seed value
	proof.Seed = FinalSeed(s)
	valid, err := proof.VerifyVRFProof()
	if err != nil {
		return vrfkey.Proof{}, errors.Wrap(err,
			"could not validate proof implied by on-chain response")
	}
	if !valid {
		return vrfkey.Proof{}, errors.Errorf(
			"proof implied by on-chain response is invalid")
	}
	return proof, nil
}

func GenerateProofResponseFromProof(proof vrfkey.Proof, s PreSeedData) (MarshaledOnChainResponse, error) {
	p := ProofResponse{P: proof, PreSeed: s.PreSeed, BlockNum: s.BlockNum}
	rv, err := p.MarshalForVRFCoordinator()
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return rv, nil
}

func GenerateProofResponseFromProofV2(proof vrfkey.Proof, s PreSeedDataV2) (MarshaledOnChainResponseV2, error) {
	p := ProofResponseV2{
		P:                proof,
		PreSeed:          s.PreSeed,
		BlockNum:         s.BlockNum,
		SubId:            s.SubId,            // Subscription ID to be charged for fulfillment
		CallbackGasLimit: s.CallbackGasLimit, // Gas limit for consumer callback
		NumWords:         s.NumWords,         // Number of random words to expand to
		Sender:           s.Sender,           // VRF consumer address
	}
	rv, err := p.MarshalForVRFCoordinator()
	if err != nil {
		return MarshaledOnChainResponseV2{}, err
	}
	return rv, nil
}

func GenerateProofResponse(keystore keystore.VRF, id string, s PreSeedData) (
	MarshaledOnChainResponse, error) {
	seed := FinalSeed(s)
	proof, err := keystore.GenerateProof(id, seed)
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return GenerateProofResponseFromProof(proof, s)
}

func generateProofResponseFromProofV2(proof vrfkey.Proof, s PreSeedDataV2) (MarshaledOnChainResponseV2, error) {
	p := ProofResponseV2{P: proof,
		PreSeed:          s.PreSeed,
		BlockNum:         s.BlockNum,
		SubId:            s.SubId,
		CallbackGasLimit: s.CallbackGasLimit,
		NumWords:         s.NumWords,
		Sender:           s.Sender,
	}
	rv, err := p.MarshalForVRFCoordinator()
	if err != nil {
		return MarshaledOnChainResponseV2{}, err
	}
	return rv, nil
}

func GenerateProofResponseV2(keystore keystore.VRF, id string, s PreSeedDataV2) (
	MarshaledOnChainResponseV2, error) {
	seedHashMsg := append(s.PreSeed[:], s.BlockHash.Bytes()...)
	seed := utils.MustHash(string(seedHashMsg)).Big()
	proof, err := keystore.GenerateProof(id, seed)
	if err != nil {
		return MarshaledOnChainResponseV2{}, err
	}
	return generateProofResponseFromProofV2(proof, s)
}
