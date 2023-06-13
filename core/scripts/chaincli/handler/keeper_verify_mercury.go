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

	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
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
	mb := MercuryBytes{}
	hc := http.DefaultClient
	feed := "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"
	feeds := []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"}
	q := url.Values{
		FeedIDHex:   {feed},
		BlockNumber: {strconv.FormatUint(blockNumber, 10)},
		UserId:      {k.cfg.UpkeepID},
		FeedID:      {strings.Join(feeds, ",")},
	}
	mercuryURL := MercuryHostV2
	path := MercuryPathV2
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, path, q.Encode())
	fmt.Printf("FeedLookup request URL: %s\n", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		fmt.Printf("cannot create request: %v\n", err)
		return
	}

	username := k.cfg.MercuryID
	password := k.cfg.MercuryKey
	if username == "" || password == "" {
		fmt.Print("username and password are empty\n")
		return
	}
	fmt.Println(username)

	ts := time.Now().UTC().UnixMilli()
	signature := generateHMAC(http.MethodGet, path+q.Encode(), []byte{}, username, password, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", username)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	retryable := false
	retryErr := retry.Do(
		func() error {
			resp, err1 := hc.Do(req)
			if err1 != nil {
				fmt.Printf("FeedLookup upkeep %s block %s GET request fails for feed %s: %v\n", k.cfg.UpkeepID, blockNumber, feed, err1)
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				fmt.Printf("FeedLookup upkeep %s block %s fails to read response body for feed %s: %v\n", k.cfg.UpkeepID, blockNumber, feed, err1)
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				fmt.Printf("FeedLookup upkeep %s block %s received status code %d for feed %s\n", k.cfg.UpkeepID, blockNumber, resp.StatusCode, feed)
				retryable = true
				return errors.New(Retry)
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for feed %s\n", k.cfg.UpkeepID, blockNumber, resp.StatusCode, feed)
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				fmt.Printf("FeedLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for feed %s: %v\n", k.cfg.UpkeepID, blockNumber, feed, err1)
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				fmt.Printf("FeedLookup upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v\n", k.cfg.UpkeepID, blockNumber, m.ChainlinkBlob, feed, err1)
				return err1
			}
			mb.Bytes = blobBytes
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == Retry
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	fmt.Printf("retryable: %v\n", retryable)
	fmt.Printf("retryErr: %v\n", retryErr)
	fmt.Printf("blob: %v\n", mb.Bytes)
	fmt.Printf("hex: %s\n", hexutil.Encode(mb.Bytes))

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

	payload, err = keeperRegistryABI.Pack("checkUpkeep", upkeepId)
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

	//var (
	//	multiErr error
	//	results  = make([]EVMAutomationUpkeepResult21, len(1))
	//)
	//
	//for i, req := range checkReqs {
	//	if req.Error != nil {
	//		r.lggr.Debugf("error encountered for key %s with message '%s' in check", keys[i], req.Error)
	//		multierr.AppendInto(&multiErr, req.Error)
	//	} else {
	//		var err error
	//		r.lggr.Debugf("UnpackCheckResult key %s checkResult: %s", string(keys[i]), *checkResults[i])
	//		results[i], err = r.packer.UnpackCheckResult(keys[i], *checkResults[i])
	//		if err != nil {
	//			return nil, errors.Wrap(err, "failed to unpack check result")
	//		}
	//	}
	//}
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
