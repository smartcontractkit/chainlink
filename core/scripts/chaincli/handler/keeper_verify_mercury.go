package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/feed_lookup_compatible_interface"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	evm "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"
)

const (
	BlockNumber        = "blockNumber" // valid for v0.2
	FeedID             = "feedID"      // valid for v0.3
	FeedIDHex          = "feedIDHex"   // valid for v0.2
	MercuryHostV2      = "https://mercury-arbitrum-testnet.chain.link"
	MercuryHostV3      = ""
	MercuryPathV2      = "/client?"
	MercuryPathV3      = "/v1/reports?"
	MercuryBatchPathV3 = "/v1/reports/bulk?"
	Retry              = "retry"
	RetryDelay         = 600 * time.Millisecond
	Timestamp          = "timestamp" // valid for v0.3
	TotalAttempt       = 3
	UserId             = "userId"
)

type FeedLookup struct {
	feedParamKey string
	feeds        []string
	timeParamKey string
	time         *big.Int
	extraData    []byte
}

type MercuryResponse struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type MercuryBytes struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     []byte
}

func (k *Keeper) VerifyFeedLookup(ctx context.Context) {

	log.Println("======================== Mercury Request ========================")
	blockNumber, err := k.client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get block number: %v", err)
	}
	hc := http.DefaultClient
	feeds := []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"}

	log.Println("======================== Allow List ========================")
	registry21, err := iregistry21.NewIKeeperRegistryMaster(common.HexToAddress(k.cfg.RegistryAddress), k.client)
	if err != nil {
		log.Fatalf("cannot create registry 2.1: %v", err)
	}
	v, err := registry21.TypeAndVersion(nil)
	if err != nil {
		log.Fatalf("failed to fetch type and version from registry 2.1: %v", err)
	}
	log.Printf("Version is %s", v)

	upkeepIds, err := registry21.GetActiveUpkeepIDs(nil, big.NewInt(0), big.NewInt(5))
	if err != nil {
		log.Fatalf("failed to fetch active upkeep ids from registry 2.1: %v", err)
	}
	log.Printf("active upkeep ids: %v", upkeepIds)

	for _, id := range upkeepIds {
		cfg, err := registry21.GetUpkeepAdminOffchainConfig(nil, id)
		if err != nil {
			log.Fatalf("failed to get upkeep admin offchain config for upkeep ID %s: %v", id, err)
		}

		var a evm.AdminOffchainConfig
		err = json.Unmarshal(cfg, &a)
		if err != nil {
			log.Fatalf("failed to unmarshal admin offchain config for upkeep ID %s: %v", id, err)
		}

		log.Printf("upkeep ID %s is mercury enabled: %v", id, a.MercuryEnabled)
	}

	log.Println("======================== check Callback ========================")
	keeperRegistryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		log.Fatalf("failed to create ABI: %v", err)
	}

	value1 := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	value2 := []byte{10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15}
	values := [][]byte{value1, value2}
	ed := []byte{120, 111, 101, 122, 90, 54, 44, 1}

	upkeepId := big.NewInt(0)
	upkeepId.SetString(k.cfg.UpkeepID, 10)

	payload, err := keeperRegistryABI.Pack("checkCallback", upkeepId, values, ed)
	if err != nil {
		log.Fatalf("failed to pack: %v", err)
	}

	var b hexutil.Bytes
	args := map[string]interface{}{
		"to":   k.cfg.RegistryAddress,
		"data": hexutil.Bytes(payload),
	}

	log.Printf("======================== for block %d ========================\n", blockNumber)
	err = k.client.Client().CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(blockNumber))
	if err != nil {
		log.Fatalf("eth call failed: %v", err)
	}

	log.Printf("checkCallback input: %s\n", hexutil.Encode(b))
	resp, err := hexutil.Decode(hexutil.Encode(b))
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	log.Printf("checkCallback input: %v\n", resp)

	out, err := keeperRegistryABI.Methods["checkCallback"].Outputs.UnpackValues(b)
	if err != nil {
		log.Fatalf("%v: unpack checkUpkeep return: %s", err, hexutil.Encode(b))
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason := *abi.ConvertType(out[2], new(uint8)).(*uint8)
	gasUsed := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	log.Printf("upkeepNeeded: %v\n", upkeepNeeded)
	log.Printf("rawPerformData: %v\n", rawPerformData)
	log.Printf("failureReason: %d\n", failureReason)
	log.Printf("gasUsed: %d\n", gasUsed)

	log.Printf("======================== for block %d ========================\n", blockNumber+1)
	err = k.client.Client().CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(blockNumber+1))
	if err != nil {
		log.Fatalf("eth call failed: %v", err)
	}

	log.Printf("checkCallback input: %s\n", hexutil.Encode(b))
	resp, err = hexutil.Decode(hexutil.Encode(b))
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	log.Printf("checkCallback input: %v\n", resp)

	out, err = keeperRegistryABI.Methods["checkCallback"].Outputs.UnpackValues(b)
	if err != nil {
		log.Fatalf("%v: unpack checkUpkeep return: %s", err, hexutil.Encode(b))
	}

	upkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	gasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	log.Printf("upkeepNeeded: %v\n", upkeepNeeded)
	log.Printf("rawPerformData: %v\n", rawPerformData)
	log.Printf("failureReason: %d\n", failureReason)
	log.Printf("gasUsed: %d\n", gasUsed)

	var (
		checkReqs    = make([]rpc.BatchElem, 1)
		checkResults = make([]*string, 1)
	)

	block := blockNumber

	opts, err := buildCallOpts(ctx, big.NewInt(int64(block)))
	if err != nil {
		log.Fatalf("failed to build call opts: %v", err)
	}

	payload, err = keeperRegistryABI.Pack("checkUpkeep", upkeepId, []byte{})
	if err != nil {
		log.Fatalf("failed to check upkeep: %v", err)
	}

	var result string
	checkReqs[0] = rpc.BatchElem{
		Method: "eth_call",
		Args: []interface{}{
			map[string]interface{}{
				"to":   k.cfg.RegistryAddress,
				"data": hexutil.Bytes(payload),
			},
			hexutil.EncodeBig(opts.BlockNumber),
		},
		Result: &result,
	}

	checkResults[0] = &result

	if err := k.client.Client().BatchCallContext(ctx, checkReqs); err != nil {
		log.Fatalf("failed to batch call: %v", err)
	}

	raw := *checkResults[0]
	b, err = hexutil.Decode(raw)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}

	out, err = keeperRegistryABI.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		log.Fatalf("failed to unpack: %v", err)
	}

	result21 := evm.EVMAutomationUpkeepResult21{
		Block:            uint32(blockNumber),
		ID:               upkeepId,
		Eligible:         true,
		CheckBlockNumber: uint32(blockNumber),
		CheckBlockHash:   [32]byte{},
	}

	upkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	result21.FailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	result21.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	result21.FastGasWei = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	result21.LinkNative = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	if !upkeepNeeded {
		result21.Eligible = false
	}
	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if result21.FailureReason == evm.UPKEEP_FAILURE_REASON_NONE || (result21.FailureReason == evm.UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED && len(rawPerformData) > 0) {
		result21.PerformData = rawPerformData
	}

	// This is a default placeholder which is used since we do not get the execute gas
	// from checkUpkeep result. This field is overwritten later from the execute gas
	// we have for an upkeep in memory. TODO (AUTO-1482): Refactor this
	result21.ExecuteGas = 5_000_000

	log.Printf("result21 FailureReason: %v", result21.FailureReason)
	log.Printf("result21 PerformData: %v", result21.PerformData)
	log.Printf("result21 Eligible: %v", result21.Eligible)
	log.Printf("result21 ID: %v", result21.ID)

	feedLookupCompatibleABI, err := abi.JSON(strings.NewReader(feed_lookup_compatible_interface.FeedLookupCompatibleInterfaceABI))
	if err != nil {
		log.Fatalf("failed to get ABI: %v", err)
	}

	e := feedLookupCompatibleABI.Errors["FeedLookup"]
	unpack, err := e.Unpack(result21.PerformData)
	if err != nil {
		log.Fatalf("failed to unpack: %v", err)
	}
	errorParameters := unpack.([]interface{})

	fl := FeedLookup{
		feedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		timeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}
	log.Printf("feedParamKey: %s", fl.feedParamKey)
	log.Printf("feeds: %v", fl.feeds)
	log.Printf("timeParamKey: %s", fl.timeParamKey)
	log.Printf("time: %s", fl.time)
	log.Printf("extraData: %v", fl.extraData)

	resultLen := len(feeds)
	ch := make(chan MercuryBytes, resultLen)
	if fl.feedParamKey == FeedIDHex && fl.timeParamKey == BlockNumber {
		// only mercury v0.2
		for i := range feeds {
			go k.singleFeedRequest(ctx, hc, ch, upkeepId, i, fl, job.MercuryV02)
		}
	}

	var reqErr error
	results := make([][]byte, len(fl.feeds))
	retryable := true
	allSuccess := true
	for i := 0; i < len(results); i++ {
		m := <-ch
		if m.Error != nil {
			reqErr = errors.Join(reqErr, m.Error)
			retryable = retryable && m.Retryable
			allSuccess = false
		}
		results[m.Index] = m.Bytes
	}
	log.Printf("FeedLookup upkeep %s retryable %v reqErr %v", upkeepId.String(), retryable && !allSuccess, reqErr)
	// only retry when not all successful AND none are not retryable
	log.Printf("results[0]: %v", results[0])
	log.Printf("results[1]: %v", results[1])
	log.Printf("retryable: %v", retryable)
	log.Printf("allSuccess: %v", allSuccess)
	log.Printf("reqErr: %v", reqErr)
}

func buildCallOpts(ctx context.Context, block *big.Int) (*bind.CallOpts, error) {
	opts := bind.CallOpts{
		Context:     ctx,
		BlockNumber: block,
	}

	return &opts, nil
}

func generateHMAC(method string, path string, body []byte, clientId string, secret string, ts int64) string {
	bodyHash := sha256.New()
	bodyHash.Write(body)
	hashString := fmt.Sprintf("%s %s %s %s %d",
		method,
		path,
		hex.EncodeToString(bodyHash.Sum(nil)),
		clientId,
		ts)
	signedMessage := hmac.New(sha256.New, []byte(secret))
	signedMessage.Write([]byte(hashString))
	userHmac := hex.EncodeToString(signedMessage.Sum(nil))
	return userHmac
}

func (k *Keeper) singleFeedRequest(ctx context.Context, hc *http.Client, ch chan<- MercuryBytes, upkeepId *big.Int, index int, ml FeedLookup, mv job.MercuryVersion) {
	q := url.Values{
		ml.feedParamKey: {ml.feeds[index]},
		ml.timeParamKey: {ml.time.String()},
		UserId:          {upkeepId.String()},
	}
	mercuryURL := MercuryHostV2
	path := MercuryPathV2
	if mv == job.MercuryV03 {
		mercuryURL = MercuryHostV3
		path = MercuryPathV3
	}
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, path, q.Encode())
	log.Printf("FeedLookup request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryBytes{Index: index, Error: err}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := generateHMAC(http.MethodGet, path+q.Encode(), []byte{}, k.cfg.MercuryID, k.cfg.MercuryKey, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", k.cfg.MercuryID)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	retryable := false
	retryErr := retry.Do(
		func() error {
			resp, err1 := hc.Do(req)
			if err1 != nil {
				log.Printf("FeedLookup upkeep %s block %s GET request fails for feed %s: %v", upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				log.Printf("FeedLookup upkeep %s block %s fails to read response body for feed %s: %v", upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				log.Printf("FeedLookup upkeep %s block %s received status code %d for feed %s", upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
				retryable = true
				return errors.New(Retry)
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for feed %s", upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				log.Printf("FeedLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for feed %s: %v", upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				return err1
			}
			log.Printf("ChainlinkBlob %d: %s", index, m.ChainlinkBlob)
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				log.Printf("FeedLookup upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v", upkeepId.String(), ml.time.String(), m.ChainlinkBlob, ml.feeds[index], err1)
				return err1
			}
			ch <- MercuryBytes{Index: index, Bytes: blobBytes}
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == Retry
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	// if all retries fail, return the error and ask the caller to handle cool down and heavyweight retry
	if retryErr != nil {
		mb := MercuryBytes{
			Index:     index,
			Retryable: retryable,
			Error:     retryErr,
		}
		ch <- mb
	}
}
