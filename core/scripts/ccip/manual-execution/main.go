package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/multierr"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"manual-execution/helpers"
)

const NumberOfBlocks = 20000

// Config represents configuration fields
type Config struct {
	SrcNodeURL       string `json:"src_rpc"`
	DestNodeURL      string `json:"dest_rpc"`
	DestOwner        string `json:"dest_owner_key"`
	CommitStore      string `json:"commit_store"`
	OffRamp          string `json:"off_ramp"`
	DestStartBlock   uint64 `json:"dest_start_block"`
	SourceChainTx    string `json:"source_chain_tx"`
	CCIPMsgID        string `json:"ccip_msg_id"`
	DestDeployedAt   uint64 `json:"dest_deployed_at"`
	GasLimitOverride uint64 `json:"gas_limit_override"`
}

type execArgs struct {
	cfg               Config
	seqNum            uint64
	msgID             [32]byte
	sourceChain       *ethclient.Client
	sourceChainId     *big.Int
	destChain         *ethclient.Client
	destUser          *bind.TransactOpts
	destChainId       *big.Int
	srcStartBlock     *big.Int
	destStartBlock    uint64
	destLatestBlock   uint64
	OnRamp            common.Address
	tokenGasOverrides []*big.Int
}

func main() {
	configPath := flag.String("configFile", "./config.json", "config for manually executing a failed ccip message "+
		"which has been successfully committed but failed to get executed")
	flag.Parse()

	if *configPath == "" {
		log.Println("config json is required")
		os.Exit(1)
	}
	cData, err := os.ReadFile(*configPath)
	if err != nil {
		log.Println("unable to read the json at ", *configPath, "error - ", err)
		os.Exit(1)
	}
	var cfg Config
	err = json.Unmarshal(cData, &cfg)
	if err != nil {
		log.Println("unable to marshal the json at ", *configPath, "error - ", err, `sample json
{
	"src_rpc": "",
	"dest_rpc": "",
	"dest_owner_key": "",
	"commit_store": "",
	"off_ramp": "",
	"dest_start_block": "",
	"ccip_send_tx": "",
	"source_start_block": "",
	"dest_deployed_at": 0,
	"gas_limit_override": 0,
}`)
		os.Exit(1)
	}
	// mandatory fields check
	err = cfg.verifyConfig()
	if err != nil {
		log.Println("config validation failed: \n", err)
		os.Exit(1)
	}
	args := &execArgs{cfg: cfg}
	err = args.populateValues()
	if err != nil {
		log.Println("error instantiating manual execution args ", err)
		os.Exit(1)
	}
	err = args.execute()
	if err != nil {
		log.Println("manual execution was not successful - ", err)
		os.Exit(1)
	}
}

func (cfg Config) verifyConfig() error {
	var allErr error
	if cfg.SrcNodeURL == "" {
		allErr = multierr.Append(allErr, fmt.Errorf("must set src_rpc - source chain rpc\n"))
	}
	if cfg.DestNodeURL == "" {
		allErr = multierr.Append(allErr, fmt.Errorf("must set dest_rpc - destination chain rpc\n"))
	}
	if cfg.DestOwner == "" {
		allErr = multierr.Append(allErr, fmt.Errorf("must set dest_owner_key - destination user private key\n"))
	}
	if cfg.SourceChainTx == "" {
		allErr = multierr.Append(allErr, fmt.Errorf("must set source_chain_tx - txHash of ccip-send request\n"))
	}

	if cfg.DestStartBlock == 0 && cfg.DestDeployedAt == 0 {
		allErr = multierr.Append(allErr, fmt.Errorf(`must set either of -
dest_deployed_at - the block number before destination contracts were deployed;
dest_start_block - the block number from which events will be filtered at destination chain.
`))
	}
	if cfg.GasLimitOverride == 0 {
		allErr = multierr.Append(allErr, fmt.Errorf("must set gas_limit_override - new value of gas limit for ccip-send request\n"))
	}
	err := helpers.VerifyAddress(cfg.CommitStore)
	if err != nil {
		allErr = multierr.Append(allErr, fmt.Errorf("check the commit_store address - %v\n", err))
	}
	err = helpers.VerifyAddress(cfg.OffRamp)
	if err != nil {
		allErr = multierr.Append(allErr, fmt.Errorf("check the off_ramp address - %v\n", err))
	}

	return allErr
}

func (args *execArgs) populateValues() error {
	var err error
	cfg := args.cfg
	args.sourceChain, err = ethclient.Dial(cfg.SrcNodeURL)
	if err != nil {
		return err
	}
	args.sourceChainId, err = args.sourceChain.ChainID(context.Background())
	if err != nil {
		return err
	}

	args.destChain, err = ethclient.Dial(cfg.DestNodeURL)
	if err != nil {
		return err
	}
	args.destChainId, err = args.destChain.ChainID(context.Background())
	if err != nil {
		return err
	}
	ownerKey, err := crypto.HexToECDSA(cfg.DestOwner)
	if err != nil {
		return err
	}

	args.destUser, err = bind.NewKeyedTransactorWithChainID(ownerKey, args.destChainId)
	if err != nil {
		return err
	}
	log.Println("--- Owner address---/n", args.destUser.From.Hex())

	var txReceipt *types.Receipt
	txReceipt, err = args.sourceChain.TransactionReceipt(context.Background(), common.HexToHash(cfg.SourceChainTx))
	if err != nil {
		return err
	}
	args.srcStartBlock = big.NewInt(0).Sub(txReceipt.BlockNumber, big.NewInt(NumberOfBlocks))
	args.destLatestBlock, err = args.destChain.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	err = args.seqNumFromCCIPSendRequested(txReceipt.Logs)
	if err != nil {
		return err
	}
	if args.cfg.DestStartBlock < 1 {
		err = args.approxDestStartBlock()
		if err != nil {
			return err
		}
	} else {
		args.destStartBlock = args.cfg.DestStartBlock
	}
	return nil
}

func (args *execArgs) execute() error {
	iterator, err := helpers.FilterReportAccepted(args.destChain, &bind.FilterOpts{Start: args.destStartBlock}, args.cfg.CommitStore)
	if err != nil {
		return err
	}

	var commitReport *helpers.ICommitStoreCommitReport
	for iterator.Next() {
		eventReport, err := iterator.CommitStoreReportAcceptedFromLog()
		if err != nil {
			return err
		}

		if eventReport.Report.Interval.Min <= args.seqNum && eventReport.Report.Interval.Max >= args.seqNum {
			commitReport = &eventReport.Report
			log.Println("Found root")
			break
		}
	}
	if commitReport == nil {
		return fmt.Errorf("unable to find seq num %d in commit report", args.seqNum)
	}
	log.Println("Executing request manually")
	seqNr := args.seqNum
	// Build a merkle tree for the report
	mctx := helpers.NewKeccakCtx()
	leafHasher := helpers.NewLeafHasher(
		GetCCIPChainSelector(args.sourceChainId.Uint64()),
		GetCCIPChainSelector(args.destChainId.Uint64()),
		args.OnRamp,
		mctx,
	)

	var leaves [][32]byte
	var curr, prove int
	var tokenData [][][]byte
	var msgs []helpers.InternalEVM2EVMMessage

	sendRequestedIterator, err := helpers.FilterCCIPSendRequested(args.sourceChain, &bind.FilterOpts{
		Start: args.srcStartBlock.Uint64(),
	}, args.OnRamp.Hex())
	if err != nil {
		return err
	}

	for sendRequestedIterator.Next() {
		event, err := sendRequestedIterator.SendRequestedEventFromLog()
		if err != nil {
			return err
		}
		if event.Message.SequenceNumber <= commitReport.Interval.Max &&
			event.Message.SequenceNumber >= commitReport.Interval.Min {
			log.Println("Found seq num in commit report", event.Message.SequenceNumber, commitReport.Interval)
			hash, err := leafHasher.HashLeaf(sendRequestedIterator.Raw)
			if err != nil {
				return err
			}
			leaves = append(leaves, hash)
			if event.Message.SequenceNumber == seqNr && event.Message.MessageID == args.msgID {
				log.Printf("Found proving %d %+v\n\n", curr, event.Message)
				msgs = append(msgs, event.Message)

				var msgTokenData [][]byte
				for range event.Message.TokenAmounts {
					msgTokenData = append(msgTokenData, []byte{})
				}

				tokenData = append(tokenData, msgTokenData)
				prove = curr
			}
			curr++
		}
	}

	sendRequestedIterator.Close()
	if len(msgs) == 0 {
		return fmt.Errorf("unable to find msg with seqNr %d", seqNr)
	}

	expectedNumberOfLeaves := int(commitReport.Interval.Max - commitReport.Interval.Min + 1)
	if len(leaves) != expectedNumberOfLeaves {
		return fmt.Errorf("not enough leaves gather to build a commit root - want %d got %d. Please set NumberOfBlocks const to a higher value", expectedNumberOfLeaves, len(leaves))
	}

	tree, err := helpers.NewTree(mctx, leaves)
	if err != nil {
		return err
	}
	if tree.Root() != commitReport.MerkleRoot {
		return errors.New("root doesn't match. cannot execute")
	}

	proof := tree.Prove([]int{prove})
	offRampProof := helpers.InternalExecutionReport{
		Messages:          msgs,
		Proofs:            proof.Hashes,
		OffchainTokenData: tokenData,
		ProofFlagBits:     helpers.ProofFlagsToBits(proof.SourceFlags),
	}

	gasLimitOverrides := make([]*helpers.EVM2EVMOffRampGasLimitOverride, len(offRampProof.Messages))

	for range offRampProof.Messages {
		evm2evmOffRampGasLimitOverride := &helpers.EVM2EVMOffRampGasLimitOverride{
			ReceiverExecutionGasLimit: big.NewInt(int64(args.cfg.GasLimitOverride)),
			TokenGasOverrides:         args.tokenGasOverrides,
		}
		gasLimitOverrides = append(gasLimitOverrides, evm2evmOffRampGasLimitOverride)
	}

	// GasLimit may need to be raised if the TX is reverting. Must be set to a value larger than the GasLimitOverride.
	// args.destUser.GasLimit = 5000000
	tx, err := helpers.ManuallyExecute(args.destChain, args.destUser, args.cfg.OffRamp, offRampProof, gasLimitOverrides)
	if err != nil {
		return err
	}
	// wait for tx confirmation
	err = helpers.WaitForSuccessfulTxReceipt(args.destChain, tx.Hash())
	if err != nil {
		log.Println("Failures may be due to insufficient gas, try increasing args.destUser.GasLimit.")
		return err
	}

	// check if the message got successfully delivered
	changed, err := helpers.FilterExecutionStateChanged(args.destChain, &bind.FilterOpts{
		Start: args.destStartBlock,
	}, args.cfg.OffRamp, []uint64{args.seqNum}, [][32]byte{args.msgID})
	if err != nil {
		return err
	}
	if changed != 2 {
		return errors.New("manual execution did not result in ExecutionStateChanged as success")
	}
	return nil
}

func (args *execArgs) seqNumFromCCIPSendRequested(logs []*types.Log) error {
	abi, err := abi.JSON(strings.NewReader(helpers.OnRampABI))
	if err != nil {
		return err
	}
	var topic0 common.Hash
	for name, abiEvent := range abi.Events {
		if name == "CCIPSendRequested" {
			topic0 = abiEvent.ID
			break
		}
	}
	if topic0 == (common.Hash{}) {
		return errors.New("no CCIPSendRequested event found in ABI")
	}
	var sendRequestedLogs []types.Log
	for _, sendReqLog := range logs {
		if sendReqLog.Topics[0] == topic0 && sendReqLog.TxHash == common.HexToHash(args.cfg.SourceChainTx) {
			args.OnRamp = sendReqLog.Address
			sendRequestedLogs = append(sendRequestedLogs, *sendReqLog)
		}
	}

	if len(sendRequestedLogs) == 0 {
		return fmt.Errorf("no CCIPSendRequested logs found for in txReceipt for txhash %s", args.cfg.SourceChainTx)
	}
	onRampContract := bind.NewBoundContract(args.OnRamp, abi, args.sourceChain, args.sourceChain, args.sourceChain)

	for _, sendReqLog := range sendRequestedLogs {
		var event helpers.SendRequestedEvent

		err = onRampContract.UnpackLog(&event, "CCIPSendRequested", sendReqLog)
		if err != nil {
			return err
		}

		if args.cfg.CCIPMsgID != "" &&
			"0x"+hex.EncodeToString(event.Message.MessageID[:]) != args.cfg.CCIPMsgID {
			continue
		}

		args.seqNum = event.Message.SequenceNumber
		args.msgID = event.Message.MessageID
		return nil
	}

	return errors.New("send request not found in logs")
}

func (args *execArgs) approxDestStartBlock() error {
	sourceBlockHdr, err := args.sourceChain.HeaderByNumber(context.Background(), args.srcStartBlock)
	if err != nil {
		return err
	}
	sendTxTime := sourceBlockHdr.Time
	maxBlockNum := args.destLatestBlock
	// setting this to an approx value of 1000 considering destination chain would have at least 1000 blocks before the transaction started
	minBlockNum := args.cfg.DestDeployedAt
	closestBlockNum := uint64(math.Floor((float64(maxBlockNum) + float64(minBlockNum)) / 2))
	var closestBlockHdr *types.Header
	closestBlockHdr, err = args.destChain.HeaderByNumber(context.Background(), big.NewInt(int64(closestBlockNum)))
	if err != nil {
		return err
	}
	// to reduce the number of RPC calls increase the value of blockOffset
	blockOffset := uint64(10)
	for {
		blockNum := closestBlockHdr.Number.Uint64()
		if minBlockNum > maxBlockNum {
			break
		}
		timeDiff := math.Abs(float64(closestBlockHdr.Time - sendTxTime))
		// break if the difference in timestamp is lesser than 1 minute
		if timeDiff < 60 {
			break
		}
		if closestBlockHdr.Time > sendTxTime {
			maxBlockNum = blockNum - 1
		} else {
			minBlockNum = blockNum + 1
		}
		closestBlockNum = uint64(math.Floor((float64(maxBlockNum) + float64(minBlockNum)) / 2))
		closestBlockHdr, err = args.destChain.HeaderByNumber(context.Background(), big.NewInt(int64(closestBlockNum)))
		if err != nil {
			return err
		}
	}

	for {
		if closestBlockHdr.Time <= sendTxTime {
			break
		}
		closestBlockNum -= blockOffset
		if closestBlockNum <= 0 {
			return errors.New("approx destination blocknumber not found")
		}
		closestBlockHdr, err = args.destChain.HeaderByNumber(context.Background(), big.NewInt(int64(closestBlockNum)))
		if err != nil {
			return err
		}
	}
	args.destStartBlock = closestBlockHdr.Number.Uint64()
	log.Printf("using approx destination start block number %d for filtering event", args.destStartBlock)
	return nil
}

func GetCCIPChainSelector(chainId uint64) uint64 {
	selector, err := chainselectors.SelectorFromChainId(chainId)
	if err != nil {
		panic(fmt.Sprintf("no chain selector for %d", chainId))
	}
	return selector
}
