package main

import (
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func main() {
	switch os.Args[1] {
	case "gen-vrf-key":
		cmd := flag.NewFlagSet("gen-vrf-key", flag.ExitOnError)
		password := cmd.String("pw", "", "password to encrypt key with")
		outfile := cmd.String("o", "key.json", "path to output file")
		helpers.ParseArgs(cmd, os.Args[2:], "pw")
		key, err := vrfkey.NewV2()
		helpers.PanicErr(err)
		exportJSON, err := key.ToEncryptedJSON(*password, utils.DefaultScryptParams)
		helpers.PanicErr(err)
		err = os.WriteFile(*outfile, exportJSON, 0600)
		helpers.PanicErr(err)
		fmt.Println("generated vrf key", key.PublicKey.String(), "and saved encrypted in", *outfile)
	case "gen-vrf-numbers":
		cmd := flag.NewFlagSet("gen-vrf-numbers", flag.ExitOnError)
		keyPath := cmd.String("key", "key.json", "path to encrypted key contents")
		numCount := cmd.Int("n", 100, "how many numbers to generate")
		password := cmd.String("pw", "", "password to decrypt the key with")
		outfile := cmd.String("o", "randomnumbers.csv", "path to output file")
		// preseed information
		senderAddr := cmd.String("sender", "", "sender of the requestRandomWords tx")
		subID := cmd.Uint64("subid", 1, "sub id")
		// seed information - can be fetched from a real chain's explorer
		blockhashStr := cmd.String("blockhash", "", "blockhash the request is in")
		blockNum := cmd.Uint64("blocknum", 10, "block number the request is in")
		cbGasLimit := cmd.Uint("cb-gas-limit", 100_000, "callback gas limit")
		numWords := cmd.Uint("num-words", 1, "num words")
		numWorkers := cmd.Uint64("num-workers", uint64(runtime.NumCPU()), "num workers")

		helpers.ParseArgs(cmd, os.Args[2:], "pw", "sender", "blockhash")

		fileBytes, err := os.ReadFile(*keyPath)
		helpers.PanicErr(err)
		key, err := vrfkey.FromEncryptedJSON(fileBytes, *password)
		helpers.PanicErr(err)

		keyHash := key.PublicKey.MustHash()
		sender := common.HexToAddress(*senderAddr)
		blockhash := common.HexToHash(*blockhashStr)

		// columns:
		// (keyHashHex, senderAddrHex, subID, nonce) preseed info
		// (preSeed, blockhash, blocknum, subID, cbGasLimit, numWords, senderAddrHex)
		// pubKeyHex, keyHashHex, senderAddrHex, subID, nonce, preSeed, blockhash, blocknum, cbGasLimit, numWords, finalSeed, proof..., randomNumber
		header := []string{
			"keyHashHex", "senderAddrHex", "subID", "nonce", "preSeed", "blockhash",
			"blocknum", "cbGasLimit", "numWords", "finalSeed",
			"proofPubKey", "proofGamma", "proofC", "proofS", "proofSeed",
			"randomNumber",
		}

		genProofs := func(
			nonceRange []uint64,
			outChan chan []string) {
			numIters := 0
			for nonce := nonceRange[0]; nonce <= nonceRange[1]; nonce++ {
				var record []string

				// construct preseed using typical preseed data
				preSeed := preseed(keyHash, sender, *subID, nonce)
				record = append(record,
					keyHash.String(), sender.String(), // keyHash, sender addr
					fmt.Sprintf("%d", *subID), fmt.Sprintf("%d", nonce), hexutil.Encode(preSeed[:]), // subId, nonce, preseed
					*blockhashStr, fmt.Sprintf("%d", *blockNum), // blockhash, blocknum
					fmt.Sprintf("%d", *cbGasLimit), fmt.Sprintf("%d", *numWords)) // cb gas limit, num words

				preseedData := proof.PreSeedDataV2{
					PreSeed:          preSeed,
					BlockHash:        blockhash,
					BlockNum:         *blockNum,
					SubId:            *subID,
					CallbackGasLimit: uint32(*cbGasLimit),
					NumWords:         uint32(*numWords),
					Sender:           sender,
				}
				finalSeed := proof.FinalSeedV2(preseedData)

				record = append(record, finalSeed.String())

				// generate proof
				pf, err2 := key.GenerateProof(finalSeed)
				helpers.PanicErr(err2)

				record = append(record,
					hex.EncodeToString(secp256k1.LongMarshal(pf.PublicKey)), // pub key
					hex.EncodeToString(secp256k1.LongMarshal(pf.Gamma)),     // gamma
					pf.C.String(), pf.S.String(), // c, s
					pf.Seed.String(), pf.Output.String()) // seed, output

				if len(record) != len(header) {
					panic("record length doesn't match header length - update one of them?")
				}
				outChan <- record
				numIters++
			}
			fmt.Println("genProofs worker wrote", numIters, "records to channel")
		}

		outFile, err := os.Create(*outfile)
		wc := utils.NewDeferableWriteCloser(outFile)
		defer wc.Close()
		helpers.PanicErr(err)

		csvWriter := csv.NewWriter(outFile)
		helpers.PanicErr(csvWriter.Write(header))
		gather := func(outChan chan []string) {
			for {
				select {
				case row := <-outChan:
					helpers.PanicErr(csvWriter.Write(row))
				case <-time.After(500 * time.Millisecond):
					// if no work is produced in this much time, we're probably done
					return
				}
			}
		}

		ranges := nonceRanges(1, uint64(*numCount), *numWorkers)

		fmt.Println("nonce ranges:", ranges, "generating proofs...")

		outC := make(chan []string)

		for _, nonceRange := range ranges {
			go genProofs(
				nonceRange,
				outC)
		}

		gather(outC)
		csvWriter.Flush()
		if csvWriter.Error() != nil {
			helpers.PanicErr(err)
		}
		if err := wc.Close(); err != nil {
			helpers.PanicErr(err)
		}
	case "verify-vrf-proofs":
		cmd := flag.NewFlagSet("verify-vrf-proofs", flag.ExitOnError)
		csvPath := cmd.String("p", "randomnumbers.csv", "path to csv file generated by gen-vrf-numbers")
		numWorkers := cmd.Int("num-workers", runtime.NumCPU()-1, "number of workers to run verification")

		helpers.ParseArgs(cmd, os.Args[2:])

		numValid := &atomic.Int64{}
		proofsChan := make(chan *vrfkey.Proof)

		verify := func(inC chan *vrfkey.Proof) {
			for {
				select {
				case pf := <-inC:
					valid, err := pf.VerifyVRFProof()
					helpers.PanicErr(err)
					if !valid {
						fmt.Println("proof", pf.String(), "is not valid", "total valid proofs:", numValid)
						panic("found invalid proof")
					}
					numValid.Add(1)
				case <-time.After(250 * time.Millisecond):
					return
				}
			}
		}

		for i := 0; i < *numWorkers; i++ {
			go verify(proofsChan)
		}

		f, err := os.Open(*csvPath)
		helpers.PanicErr(err)
		defer f.Close()
		reader := csv.NewReader(f)
		_, err = reader.Read() // read the column titles
		helpers.PanicErr(err)
		for {
			rec, err := reader.Read()
			if err != nil {
				break
			}
			proofPubKey, err := secp256k1.LongUnmarshal(hexutil.MustDecode("0x" + rec[len(rec)-6]))
			helpers.PanicErr(err)
			proofGamma, err := secp256k1.LongUnmarshal(hexutil.MustDecode("0x" + rec[len(rec)-5]))
			helpers.PanicErr(err)
			proofC := decimal.RequireFromString(rec[len(rec)-4]).BigInt()
			proofS := decimal.RequireFromString(rec[len(rec)-3]).BigInt()
			proofSeed := decimal.RequireFromString(rec[len(rec)-2]).BigInt()
			proofOutput := decimal.RequireFromString(rec[len(rec)-1]).BigInt()

			pf := &vrfkey.Proof{
				PublicKey: proofPubKey,
				Gamma:     proofGamma,
				C:         proofC,
				S:         proofS,
				Seed:      proofSeed,
				Output:    proofOutput,
			}

			proofsChan <- pf
		}
		fmt.Println("all proofs valid! num proofs:", numValid)
	}
}

func preseed(keyHash common.Hash, sender common.Address, subID, nonce uint64) [32]byte {
	encoded, err := utils.ABIEncode(
		`[{"type":"bytes32"}, {"type":"address"}, {"type":"uint64"}, {"type", "uint64"}]`,
		keyHash,
		sender,
		subID,
		nonce)
	helpers.PanicErr(err)
	preSeed := crypto.Keccak256(encoded)
	var preSeedSized [32]byte
	copy(preSeedSized[:], preSeed)
	return preSeedSized
}

func nonceRanges(start, end, numWorkers uint64) (ranges [][]uint64) {
	rangeSize := (end - start) / numWorkers
	for i := start; i <= end; i += rangeSize + 1 {
		j := i + rangeSize
		if j > end {
			j = end
		}

		ranges = append(ranges, []uint64{i, j})
	}
	return
}
