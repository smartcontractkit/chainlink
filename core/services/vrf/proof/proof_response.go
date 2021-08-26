package proof

// Contains logic/data for mandatorily mixing VRF seeds with the hash of the
// block in which a VRF request appeared

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
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

const RequestCommitmentLength = 32 + // blocknum
	32 + // subID
	32 + // gaslimit
	32 + // numWords
	32 // sender

// MarshaledOnChainResponse is the flat bytes which are sent back to the
// VRFCoordinator.
type MarshaledOnChainResponse [OnChainResponseLength]byte

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
	[ProofLength]byte, [RequestCommitmentLength]byte, error) {
	var proof [ProofLength]byte
	var rc [RequestCommitmentLength]byte
	solidityProof, err := SolidityPrecalculations(&p.P)
	if err != nil {
		return proof, rc, errors.Wrap(err,
			"while marshaling proof for VRFCoordinatorV2")
	}
	solidityProof.P.Seed = common.BytesToHash(p.PreSeed[:]).Big()
	x, y := secp256k1.Coordinates(solidityProof.P.PublicKey)
	gx, gy := secp256k1.Coordinates(solidityProof.P.Gamma)
	cgx, cgy := secp256k1.Coordinates(solidityProof.CGammaWitness)
	shx, shy := secp256k1.Coordinates(solidityProof.SHashWitness)
	proofBytes, err := utils.GenericEncode([]string{
		"uint256[2]", // pk
		"uint256[2]", // gamma
		"uint256",    // c
		"uint256",    // s
		"uint256",    // seed
		"address",    // uWit
		"uint256[2]", // cGammaWitness
		"uint256[2]", // sHashWitness
		"uint256",    // zInv
	},
		[]interface{}{
			[2]*big.Int{x, y},
			[2]*big.Int{gx, gy},
			solidityProof.P.C,
			solidityProof.P.S,
			solidityProof.P.Seed,
			solidityProof.UWitness,
			[2]*big.Int{cgx, cgy},
			[2]*big.Int{shx, shy},
			solidityProof.ZInv,
		}...)
	if err != nil {
		return proof, rc, err
	}
	rcBytes, err := utils.GenericEncode([]string{
		"uint64",  // blocknum
		"uint64",  // subId
		"uint32",  // callbackLimit
		"uint32",  // numWords
		"address", // sender
	},
		[]interface{}{
			p.BlockNum,
			p.SubId,
			p.CallbackGasLimit,
			p.NumWords,
			p.Sender,
		}...)
	if err != nil {
		return proof, rc, err
	}
	if len(proofBytes) != ProofLength {
		return proof, rc, errors.Errorf(
			"wrong length of proof for VRFCoordinatorV2")
	}
	if len(rcBytes) != RequestCommitmentLength {
		return proof, rc, errors.Errorf(
			"wrong length of request commitment for VRFCoordinatorV2")
	}
	copy(rc[:], rcBytes[:])
	copy(proof[:], proofBytes[:])
	return proof, rc, nil
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

func GenerateProofResponseFromProofV2(proof vrfkey.Proof, s PreSeedDataV2) ([ProofLength]byte, [RequestCommitmentLength]byte, error) {
	p := ProofResponseV2{
		P:                proof,
		PreSeed:          s.PreSeed,
		BlockNum:         s.BlockNum,
		SubId:            s.SubId,            // Subscription ID to be charged for fulfillment
		CallbackGasLimit: s.CallbackGasLimit, // Gas limit for consumer callback
		NumWords:         s.NumWords,         // Number of random words to expand to
		Sender:           s.Sender,           // VRF consumer address
	}
	return p.MarshalForVRFCoordinator()
}

func GenerateProofResponse(keystore *keystore.VRF, key secp256k1.PublicKey, s PreSeedData) (
	MarshaledOnChainResponse, error) {
	seed := FinalSeed(s)
	proof, err := keystore.GenerateProof(key, seed)
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return GenerateProofResponseFromProof(proof, s)
}

func GenerateProofResponseV2(keystore *keystore.VRF, key secp256k1.PublicKey, s PreSeedDataV2) (
	[ProofLength]byte, [RequestCommitmentLength]byte, error) {
	seedHashMsg := append(s.PreSeed[:], s.BlockHash.Bytes()...)
	seed := utils.MustHash(string(seedHashMsg)).Big()
	proof, err := keystore.GenerateProof(key, seed)
	var p [ProofLength]byte
	var rc [RequestCommitmentLength]byte
	if err != nil {
		return p, rc, err
	}
	return GenerateProofResponseFromProofV2(proof, s)
}
