package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/extraargs"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
)

var vrfProofTemplate = `{
	pk: [
		%s,
		%s
	],
	gamma: [
		%s,
		%s
	],
	c: %s,
	s: %s,
	seed: %s,
	uWitness: %s,
	cGammaWitness: [
		%s,
		%s
	],
	sHashWitness: [
		%s,
		%s
	],
	zInv: %s
}
`

var rcTemplate = `{
	blockNum: %d,
	subId: %d,
	callbackGasLimit: %d,
	numWords: %d,
	sender: %s,
	extraArgs: %s
}
`

func generateProofForV2Plus(e helpers.Environment) {
	deployCmd := flag.NewFlagSet("generate-proof-v2-plus", flag.ExitOnError)

	keyHashString := deployCmd.String("key-hash", "", "key hash for VRF request")
	preSeedString := deployCmd.String("pre-seed", "", "pre-seed for VRF request")
	blockhashString := deployCmd.String("block-hash", "", "blockhash of VRF request")
	blockNum := deployCmd.Uint64("block-num", 0, "block number of VRF request")
	senderString := deployCmd.String("sender", "", "requestor of VRF request")
	secretKeyString := deployCmd.String("secret-key", "10", "secret key for VRF V2Key")
	subId := deployCmd.String("sub-id", "1", "subscription Id for VRF request")
	callbackGasLimit := deployCmd.Uint64("gas-limit", 1_000_000, "callback gas limit for VRF request")
	numWords := deployCmd.Uint64("num-words", 1, "number of words for VRF request")
	nativePayment := deployCmd.Bool("native-payment", false, "requestor of VRF request")

	helpers.ParseArgs(
		deployCmd, os.Args[2:], "key-hash", "pre-seed", "block-hash", "block-num", "sender",
	)

	// Generate V2Key from secret key.
	secretKey := decimal.RequireFromString(*secretKeyString).BigInt()
	key := vrfkey.MustNewV2XXXTestingOnly(secretKey)
	uncompressed, err := key.PublicKey.StringUncompressed()
	if err != nil {
		panic(err)
	}
	pk := key.PublicKey
	pkh := pk.MustHash()
	fmt.Println("Compressed: ", pk.String())
	fmt.Println("Uncompressed: ", uncompressed)
	fmt.Println("Hash: ", pkh.String())

	// Parse big ints and hexes.
	requestKeyHash := common.HexToHash(*keyHashString)
	requestPreSeed := decimal.RequireFromString(*preSeedString).BigInt()
	sender := common.HexToAddress(*senderString)
	blockHash := common.HexToHash(*blockhashString)

	// Ensure that the provided keyhash of the request matches the keyhash of the secret key.
	if !bytes.Equal(requestKeyHash[:], pkh[:]) {
		helpers.PanicErr(errors.New("invalid key hash"))
	}

	// Generate proof.
	preSeed, err := proof.BigToSeed(requestPreSeed)
	if err != nil {
		helpers.PanicErr(fmt.Errorf("unable to parse preseed: %w", err))
	}

	parsedSubId, ok := new(big.Int).SetString(*subId, 10)
	if !ok {
		helpers.PanicErr(fmt.Errorf("unable to parse subID: %s %w", *subId, err))
	}
	extraArgs, err := extraargs.EncodeV1(*nativePayment)
	helpers.PanicErr(err)
	preSeedData := proof.PreSeedDataV2Plus{
		PreSeed:          preSeed,
		BlockHash:        blockHash,
		BlockNum:         *blockNum,
		SubId:            parsedSubId,
		CallbackGasLimit: uint32(*callbackGasLimit),
		NumWords:         uint32(*numWords),
		Sender:           sender,
		ExtraArgs:        extraArgs,
	}
	finalSeed := proof.FinalSeedV2Plus(preSeedData)
	p, err := key.GenerateProof(finalSeed)
	if err != nil {
		helpers.PanicErr(fmt.Errorf("unable to generate proof: %w", err))
	}
	onChainProof, rc, err := proof.GenerateProofResponseFromProofV2Plus(p, preSeedData)
	if err != nil {
		helpers.PanicErr(fmt.Errorf("unable to generate proof response: %w", err))
	}

	// Print formatted VRF proof.
	fmt.Println("ON-CHAIN PROOF:")
	fmt.Printf(
		vrfProofTemplate,
		onChainProof.Pk[0],
		onChainProof.Pk[1],
		onChainProof.Gamma[0],
		onChainProof.Gamma[1],
		onChainProof.C,
		onChainProof.S,
		onChainProof.Seed,
		onChainProof.UWitness,
		onChainProof.CGammaWitness[0],
		onChainProof.CGammaWitness[1],
		onChainProof.SHashWitness[0],
		onChainProof.SHashWitness[1],
		onChainProof.ZInv,
	)

	// Print formatted request commitment.
	fmt.Println("\nREQUEST COMMITMENT:")
	fmt.Printf(
		rcTemplate,
		rc.BlockNum,
		rc.SubId,
		rc.CallbackGasLimit,
		rc.NumWords,
		rc.Sender,
		hexutil.Encode(rc.ExtraArgs),
	)
}
