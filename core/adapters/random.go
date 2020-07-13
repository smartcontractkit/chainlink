package adapters

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Random adapter type implements VRF calculation in its Perform method.
//
// The VRFCoordinator.sol contract and its integration with the chainlink node
// will handle interaction with the Random adapter, but if you need to interact
// with it directly, its input to should be a JSON object with "seed" and
// "keyHash" fields containing the input seed as a hex-represented uint256, and
// the keccak256 hash of the UNCOMPRESSED REPRESENTATION(*) of the public key
// E.g., given the input
//
//   {
//     "seed":
//       "0x0000000000000000000000000000000000000000000000000000000000000001",
//     "keyHash":
//       "0xc0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf",
//   }
//
// the adapter will return a proof for the VRF output given seed 1, as long as
// the keccak256 hash of its public key matches the hash in the input.
// Otherwise, it will error.
//
// The adapter returns the hex representation of a solidity bytes array which
// can be verified on-chain by VRF.sol#randomValueFromVRFProof. (I.e., it is the
// proof expected by that method, prepended by its length as a uint256.)
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
	key, err := getKey(ra, input)
	if err != nil {
		return models.NewRunOutputError(errors.Wrapf(err, "bad key for vrf task"))
	}
	seed, err := getSeed(input)
	if err != nil {
		return models.NewRunOutputError(errors.Wrap(err, "bad seed for vrf task"))
	}
	solidityProof, err := store.VRFKeyStore.GenerateProof(key, seed)
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

// getSeed returns the numeric seed for the vrf task, or an error
func getSeed(input models.RunInput) (*big.Int, error) {
	rawSeed, err := extractHex(input, "seed")
	if err != nil {
		return nil, err
	}
	seed := big.NewInt(0).SetBytes(rawSeed)
	if err := utils.CheckUint256(seed); err != nil {
		return nil, err
	}
	return seed, nil
}

// getKey returns the public key for the VRF, or an error.
func getKey(ra *Random, input models.RunInput) (*vrfkey.PublicKey, error) {
	inputKeyHash, err := extractHex(input, "keyHash")
	if err != nil {
		return nil, err
	}
	key, err := vrfkey.NewPublicKeyFromHex(ra.PublicKey)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse %v as public key", ra.PublicKey)
	}
	keyHash, err := key.Hash()
	if err != nil {
		return nil, errors.Wrapf(err, "could not compute %v' hash", ra.PublicKey)
	}

	if keyHash != common.BytesToHash(inputKeyHash) {
		return nil, fmt.Errorf(
			"this task's keyHash %x does not match the input hash %x", keyHash, inputKeyHash)
	}
	return key, nil
}

// extractHex returns the bytes corresponding to the string input at the key
// field, or an error.
func extractHex(input models.RunInput, key string) ([]byte, error) {
	rawValue := input.Data().Get(key)
	if rawValue.Type != gjson.String {
		return nil, fmt.Errorf("%s %#+v is not a hex string", key, rawValue)
	}
	return hexutil.Decode(rawValue.String())
}
