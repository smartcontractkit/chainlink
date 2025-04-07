package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/olekukonko/tablewriter"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepers20 "github.com/smartcontractkit/chainlink-automation/pkg/v2"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"
	evm "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v20"
)

type OCR2ReportDataElem struct {
	Err                string
	From               string
	To                 string
	ChainID            string
	BlockNumber        string
	PerformKeys        string
	PerformBlockChecks string
}

// JsonError is a rpc.jsonError interface
type JsonError interface {
	Error() string
	// ErrorCode() int
	ErrorData() interface{}
}

func OCR2AutomationReports(hdlr *baseHandler, txs []string) error {
	latestBlock, err := hdlr.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	fmt.Println("")
	fmt.Printf("latest block: %s\n", latestBlock.Number())
	fmt.Println("")

	txRes, txErr, err := getTransactionDetailForHashes(hdlr, txs)
	if err != nil {
		return fmt.Errorf("batch call error: %w", err)
	}

	ocr2Txs := make([]*OCR2TransmitTx, len(txRes))
	elements := make([]OCR2ReportDataElem, len(txRes))
	simBatch := make([]rpc.BatchElem, len(txRes))
	for i := range txRes {
		if txErr[i] != nil {
			elements[i].Err = txErr[i].Error()
			continue
		}

		if txRes[i] == nil {
			elements[i].Err = "nil response"
			continue
		}

		ocr2Txs[i], err = NewOCR2TransmitTx(*txRes[i])
		if err != nil {
			elements[i].Err = fmt.Sprintf("failed to create ocr2 transaction: %s", err)
			continue
		}

		ocr2Txs[i].SetStaticValues(&elements[i])
		simBatch[i], err = ocr2Txs[i].BatchElem()
		if err != nil {
			return err
		}
	}

	txRes, txErr, err = getSimulationsForTxs(hdlr, simBatch)
	if err != nil {
		return err
	}
	for i := range txRes {
		if txErr[i] == nil {
			continue
		}

		err2, ok := txErr[i].(JsonError) //nolint:errorlint
		if ok {
			decoded, err := hexutil.Decode(err2.ErrorData().(string))
			if err != nil {
				elements[i].Err = err.Error()
				continue
			}

			elements[i].Err = ocr2Txs[i].DecodeError(decoded)
		} else if err2 != nil {
			elements[i].Err = err2.Error()
		}
	}

	data := make([][]string, len(elements))
	for i, elem := range elements {
		data[i] = []string{
			txs[i],
			elem.ChainID,
			elem.BlockNumber,
			elem.Err,
			elem.From,
			elem.To,
			elem.PerformKeys,
			elem.PerformBlockChecks,
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][2] > data[j][2]
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Hash", "ChainID", "Block", "Error", "From", "To", "Keys", "CheckBlocks"})
	// table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return nil
}

func getTransactionDetailForHashes(hdlr *baseHandler, txs []string) ([]*map[string]interface{}, []error, error) {
	var (
		txReqs = make([]rpc.BatchElem, len(txs))
		txRes  = make([]*map[string]interface{}, len(txs))
		txErr  = make([]error, len(txs))
	)

	for i, txHash := range txs {
		b, err := common.ParseHexOrString(txHash)
		if err != nil {
			return txRes, txErr, fmt.Errorf("failed to parse transaction hash: %s", txHash)
		}

		var result map[string]interface{}
		txReqs[i] = rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args: []interface{}{
				common.BytesToHash(b),
			},
			Result: &result,
		}

		txRes[i] = &result
	}

	err := hdlr.rpcClient.BatchCallContext(context.Background(), txReqs)

	for i := range txReqs {
		txErr[i] = txReqs[i].Error
	}

	return txRes, txErr, err
}

func getSimulationsForTxs(hdlr *baseHandler, txReqs []rpc.BatchElem) ([]*map[string]interface{}, []error, error) {
	var (
		txRes = make([]*map[string]interface{}, len(txReqs))
		txErr = make([]error, len(txReqs))
	)

	for i := range txReqs {
		var result map[string]interface{}
		txReqs[i].Result = &result
		txRes[i] = &result
	}

	err := hdlr.rpcClient.BatchCallContext(context.Background(), txReqs)

	for i := range txReqs {
		txErr[i] = txReqs[i].Error
	}

	return txRes, txErr, err
}

func NewOCR2Transaction(raw map[string]interface{}) (*OCR2Transaction, error) {
	contract, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	txBytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var tx types.Transaction
	if err := json.Unmarshal(txBytes, &tx); err != nil {
		return nil, err
	}

	return &OCR2Transaction{
		encoder: evm.EVMAutomationEncoder20{},
		abi:     contract,
		raw:     raw,
		tx:      &tx,
	}, nil
}

type OCR2Transaction struct {
	encoder evm.EVMAutomationEncoder20
	abi     abi.ABI
	raw     map[string]interface{}
	tx      *types.Transaction
}

func (t *OCR2Transaction) TransactionHash() common.Hash {
	return t.tx.Hash()
}

func (t *OCR2Transaction) ChainId() *big.Int {
	return t.tx.ChainId()
}

func (t *OCR2Transaction) BlockNumber() (uint64, error) {
	if bl, ok := t.raw["blockNumber"]; ok {
		var blStr string
		blStr, ok = bl.(string)
		if ok {
			block, err := hexutil.DecodeUint64(blStr)
			if err != nil {
				return 0, fmt.Errorf("failed to parse block number: %w", err)
			}
			return block, nil
		}
		return 0, errors.New("not a string")
	}
	return 0, errors.New("not found")
}

func (t *OCR2Transaction) To() *common.Address {
	return t.tx.To()
}

func (t *OCR2Transaction) From() (common.Address, error) {
	switch t.tx.Type() {
	case 2:
		from, err := types.Sender(types.NewLondonSigner(t.tx.ChainId()), t.tx)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to get from addr: %w", err)
		} else {
			return from, nil
		}
	}

	return common.Address{}, errors.New("from address not found")
}

func (t *OCR2Transaction) Method() (*abi.Method, error) {
	return t.abi.MethodById(t.tx.Data()[0:4])
}

func (t *OCR2Transaction) DecodeError(b []byte) string {
	j := common.Bytes2Hex(b)

	for _, e := range t.abi.Errors {
		if bytes.Equal(e.ID[:4], b[:4]) {
			return e.Name
		}
	}

	return j
}

func NewOCR2TransmitTx(raw map[string]interface{}) (*OCR2TransmitTx, error) {
	tx, err := NewOCR2Transaction(raw)
	if err != nil {
		return nil, err
	}

	return &OCR2TransmitTx{
		OCR2Transaction: *tx,
	}, nil
}

type OCR2TransmitTx struct {
	OCR2Transaction
}

func (t *OCR2TransmitTx) UpkeepsInTransmit() ([]ocr2keepers20.UpkeepResult, error) {
	txData := t.tx.Data()

	// recover Method from signature and ABI
	method, err := t.abi.MethodById(txData[0:4])
	if err != nil {
		return nil, fmt.Errorf("failed to get method from sig: %w", err)
	}

	vals := make(map[string]interface{})
	if err := t.abi.Methods[method.Name].Inputs.UnpackIntoMap(vals, txData[4:]); err != nil {
		return nil, fmt.Errorf("unpacking error: %w", err)
	}

	reportData, ok := vals["rawReport"]
	if !ok {
		return nil, errors.New("raw report data missing from input")
	}

	reportBytes, ok := reportData.([]byte)
	if !ok {
		return nil, fmt.Errorf("report data not bytes: %T", reportData)
	}

	return t.encoder.DecodeReport(reportBytes)
}

func (t *OCR2TransmitTx) SetStaticValues(elem *OCR2ReportDataElem) {
	if t.To() != nil {
		elem.To = t.To().String()
	}

	elem.ChainID = t.ChainId().String()

	from, err := t.From()
	if err != nil {
		elem.Err = err.Error()
		return
	}
	elem.From = from.String()

	block, err := t.BlockNumber()
	if err != nil {
		elem.Err = err.Error()
		return
	}
	elem.BlockNumber = strconv.FormatUint(block, 10)

	upkeeps, err := t.UpkeepsInTransmit()
	if err != nil {
		elem.Err = err.Error()
	}

	keys := []string{}
	chkBlocks := []string{}

	for _, u := range upkeeps {
		val, ok := u.(evm.EVMAutomationUpkeepResult20)
		if !ok {
			panic("unrecognized upkeep result type")
		}

		keys = append(keys, val.ID.String())
		chkBlocks = append(chkBlocks, strconv.FormatUint(uint64(val.CheckBlockNumber), 10))
	}

	elem.PerformKeys = strings.Join(keys, "\n")
	elem.PerformBlockChecks = strings.Join(chkBlocks, "\n")
}

func (t *OCR2TransmitTx) BatchElem() (rpc.BatchElem, error) {
	bn, err := t.BlockNumber()
	if err != nil {
		return rpc.BatchElem{}, err
	}

	from, err := t.From()
	if err != nil {
		return rpc.BatchElem{}, err
	}

	return rpc.BatchElem{
		Method: "eth_call",
		Args: []interface{}{
			map[string]interface{}{
				"from": from.Hex(),
				"to":   t.To().Hex(),
				"data": hexutil.Bytes(t.tx.Data()),
			},
			hexutil.EncodeBig(big.NewInt(int64(bn) - 1)),
		},
	}, nil
}

func NewBaseOCR2Tx(tx *types.Transaction) (*BaseOCR2Tx, error) {
	contract, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &BaseOCR2Tx{
		abi:         contract,
		Transaction: *tx,
	}, nil
}

type BaseOCR2Tx struct {
	abi abi.ABI
	types.Transaction
}

func (tx *BaseOCR2Tx) Method() (*abi.Method, error) {
	return tx.abi.MethodById(tx.Data()[0:4])
}

func (tx *BaseOCR2Tx) DataMap() (map[string]interface{}, error) {
	txData := tx.Data()

	// recover Method from signature and ABI
	method, err := tx.abi.MethodById(txData[0:4])
	if err != nil {
		return nil, fmt.Errorf("failed to get method from sig: %w", err)
	}

	vals := make(map[string]interface{})
	if err := tx.abi.Methods[method.Name].Inputs.UnpackIntoMap(vals, txData[4:]); err != nil {
		return nil, fmt.Errorf("unpacking error: %w", err)
	}

	return vals, nil
}

func NewOCR2SetConfigTx(tx *types.Transaction) (*OCR2SetConfigTx, error) {
	base, err := NewBaseOCR2Tx(tx)
	if err != nil {
		return nil, err
	}

	return &OCR2SetConfigTx{
		BaseOCR2Tx: *base,
	}, nil
}

type OCR2SetConfigTx struct {
	BaseOCR2Tx
}

func (tx *OCR2SetConfigTx) Config() (ocrtypes.ContractConfig, error) {
	conf := ocrtypes.ContractConfig{}

	vals, err := tx.DataMap()
	if err != nil {
		return conf, err
	}

	if fVal, ok := vals["f"]; ok {
		conf.F = fVal.(uint8)
	}

	if onVal, ok := vals["onchainConfig"]; ok {
		conf.OnchainConfig = onVal.([]byte)
	}

	if vVal, ok := vals["offchainConfigVersion"]; ok {
		conf.OffchainConfigVersion = vVal.(uint64)
	}

	if onVal, ok := vals["offchainConfig"]; ok {
		conf.OffchainConfig = onVal.([]byte)
	}

	if sVal, ok := vals["signers"]; ok {
		for _, s := range sVal.([]common.Address) {
			conf.Signers = append(conf.Signers, s.Bytes())
		}
	}

	if tVal, ok := vals["transmitters"]; ok {
		for _, t := range tVal.([]common.Address) {
			conf.Transmitters = append(conf.Transmitters, ocrtypes.Account(t.Hex()))
		}
	}

	return conf, nil
}
