package adapters

import (
	"fmt"
	"math"

	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Random adapter type implements VRF calculation in its Perform method.
//
// The VRFCoordinator.sol contract and its integration with the chainlink node
// will handle interaction with the Random adapter, but if you need to interact
// with it directly, its input to should be a JSON object with "preSeed",
// "blockHash", "blockNum", and "keyHash" fields containing, respectively,
//
// - The input seed as a hex-represented uint256 (this is the preSeed generated
//   by VRFCoordinator#requestRandomness)
// - The hex-represented hash of the block in which request appeared
// - The number of the block in which the request appeared, as a JSON number
// - The keccak256 hash of the UNCOMPRESSED REPRESENTATION(*) of the public key
//
// E.g., given the input
//
//   {
//     "preSeed":
//       "0x0000000000000000000000000000000000000000000000000000000000000001",
//     "blockHash":
//       "0x31dcb7c2e3f80ce552bf730d5c1a7ed7f9b42c17aff254729b5be081394617e6",
//     "blockNum": 10000000,
//     "keyHash":
//       "0xc0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf",
//   }
//
// The adapter will return a proof for the VRF output given these values, as
// long as the keccak256 hash of its public key matches the hash in the input.
// Otherwise, it will error.
//
// The seed which is actually passed to the VRF cryptographic module is
// controlled by vrf.FinalSeed, and is computed from the above inputs.
//
// The adapter returns the hex representation of a solidity bytes array suitable
// for passing to VRFCoordinator#fulfillRandomnessRequest, a
// vrf.MarshaledOnChainResponse.
//
// (*) I.e., the 64-byte concatenation of the point's x- and y- ordinates as
// uint256's
type Random struct {
	// Compressed hex representation public key used in Random's VRF proofs
	//
	// This is just a hex string because Random is instantiated by json.Unmarshal.
	// (See adapters.For function.)
	PublicKey string `json:"publicKey"`
}

// TaskType returns the type of Adapter.
func (ra *Random) TaskType() models.TaskType {
	return TaskTypeRandom
}

// Perform returns the the proof for the VRF output given seed, or an error.
func (ra *Random) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	key, i, err := getInputs(ra, input, store)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	solidityProof, err := store.VRFKeyStore.GenerateProof(key, i)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	vrfCoordinatorArgs, err := models.VRFFulfillMethod().Inputs.PackValues(
		[]interface{}{
			solidityProof[:], // geth expects slice, even if arg is constant-length
		})
	if err != nil {
		return models.NewRunOutputError(errors.Wrapf(err,
			"while packing VRF proof %s as argument to "+
				"VRFCoordinator.fulfillRandomnessRequest", solidityProof))
	}
	return models.NewRunOutputCompleteWithResult(fmt.Sprintf("0x%x",
		vrfCoordinatorArgs))
}

// getInputs parses the JSON input for the values needed by the random adapter,
// or returns an error.
func getInputs(ra *Random, input models.RunInput, store *store.Store) (
	vrfkey.PublicKey, vrf.PreSeedData, error) {
	key, err := getKey(ra, input)
	if err != nil {
		return vrfkey.PublicKey{}, vrf.PreSeedData{}, errors.Wrapf(err,
			"bad key for vrf task")
	}
	preSeed, err := getPreSeed(input)
	if err != nil {
		return vrfkey.PublicKey{}, vrf.PreSeedData{}, errors.Wrap(err,
			"bad seed for vrf task")
	}
	block, err := getBlockData(input)
	if err != nil {
		return vrfkey.PublicKey{}, vrf.PreSeedData{}, err
	}
	s := vrf.PreSeedData{PreSeed: preSeed, BlockHash: block.hash, BlockNum: block.num}
	return key, s, nil
}

// block contains information about the block containing the VRF request, which
// is to be mixed with the request's seed
type block struct {
	hash common.Hash // Block hash
	num  uint64      // Cardinal number of block
}

// getBlockData parses the block-related data from the JSON input passed to the
// random adapter
func getBlockData(input models.RunInput) (block, error) {
	hashBytes, err := extractHex(input, "blockHash")
	if err != nil {
		return block{}, errors.Wrap(err, "bad blockHash for vrf task")
	}
	bHash := common.BytesToHash(hashBytes)

	rawBlockNum := input.Data().Get("blockNum")
	if rawBlockNum.Type != gjson.Number {
		return block{}, errors.Errorf("blockNum field has no number: %+v",
			rawBlockNum)
	}
	if rawBlockNum.Float() >= math.MaxUint64 {
		return block{}, errors.Errorf("blockNum %f too big for Int64",
			rawBlockNum.Float())
	}
	directBlockNum := uint64(rawBlockNum.Float())
	if float64(directBlockNum) != rawBlockNum.Float() {
		return block{}, errors.Errorf("blockNum %f is not a natural number",
			rawBlockNum.Float())
	}
	return block{bHash, directBlockNum}, nil
}

// getPreSeed returns the numeric seed for the vrf task, or an error
func getPreSeed(input models.RunInput) (vrf.Seed, error) {
	rawSeed, err := extractHex(input, "seed")
	if err != nil {
		return vrf.Seed{}, err
	}
	rv, err := vrf.BytesToSeed(rawSeed)
	if err != nil {
		return vrf.Seed{}, err
	}
	if rv == nil {
		return vrf.Seed{}, errors.Errorf("nil pre-seed from %+v", rawSeed)
	}
	return *rv, nil
}

func checkKeyHash(key vrfkey.PublicKey, inputKeyHash []byte) error {
	keyHash, err := key.Hash()
	if err != nil {
		return errors.Wrapf(err, "could not compute %v' hash", key)
	}

	if keyHash != common.BytesToHash(inputKeyHash) {
		return fmt.Errorf("this task's keyHash %x does not match the input hash %x",
			keyHash, inputKeyHash)
	}
	return nil
}

var failedKey = vrfkey.PublicKey{}

// getKey returns the public key for the VRF, or an error.
func getKey(ra *Random, input models.RunInput) (vrfkey.PublicKey, error) {
	key, err := vrfkey.NewPublicKeyFromHex(ra.PublicKey)
	if err != nil {
		return failedKey, errors.Wrapf(err, "could not parse %v as public key",
			ra.PublicKey)
	}
	if key.IsZero() {
		return failedKey, errors.Wrapf(err, "zero public key!")
	}
	inputKeyHash, err := extractHex(input, "keyHash")
	if err != nil {
		return failedKey, err
	}
	if err = checkKeyHash(key, inputKeyHash); err != nil {
		return failedKey, err
	}
	return key, nil
}

// Max length of a solidity bytes32 as 0x-hex string
const bytes32HexRepresentationLength = /* 0x */ 2 +
	/* num bytes */ 32*2 // two nybbles per byte

// extractHex returns the bytes corresponding to the string input at the key
// field, or an error if the string cannot be interpreted as a hex string or
// represents more than 32 bytes
func extractHex(input models.RunInput, key string) ([]byte, error) {
	rawValue := input.Data().Get(key)
	if rawValue.Type != gjson.String {
		return nil, fmt.Errorf("%s %#+v is not a hex string", key, rawValue)
	}
	if len(rawValue.String()) > bytes32HexRepresentationLength {
		return nil, fmt.Errorf("%s should be a hex string representing at most "+
			"32 bytes", rawValue.String())
	}
	return hexutil.Decode(rawValue.String())
}
