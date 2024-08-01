package cosmwasm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	cosmosSDK "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
)

type OCR2Reader struct {
	address     cosmosSDK.AccAddress
	chainReader client.Reader
	lggr        logger.Logger
}

func NewOCR2Reader(addess cosmosSDK.AccAddress, chainReader client.Reader, lggr logger.Logger) *OCR2Reader {
	return &OCR2Reader{
		address:     addess,
		chainReader: chainReader,
		lggr:        lggr,
	}
}

func (r *OCR2Reader) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	resp, err := r.chainReader.ContractState(
		r.address,
		[]byte(`{"latest_config_details":{}}`),
	)
	if err != nil {
		return
	}
	var config ConfigDetails
	if err = json.Unmarshal(resp, &config); err != nil {
		return
	}
	changedInBlock = config.BlockNumber
	configDigest = config.ConfigDigest
	return
}

func (r *OCR2Reader) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	// previously we queried with constraint "wasm-set_config._contract_address='address'" directly, but that does not
	// work with wasmd 0.41.0, which is at cosmos-sdk v0.47.4, which contains the following regex for each event query string:
	// https://github.com/cosmos/cosmos-sdk/blob/3b509c187e1643757f5ef8a0b5ae3decca0c7719/x/auth/tx/service.go#L49
	query := []string{fmt.Sprintf("tx.height=%d", changedInBlock), fmt.Sprintf("wasm._contract_address='%s'", r.address)}
	res, err := r.chainReader.TxsEvents(query, nil)
	if err != nil {
		return types.ContractConfig{}, err
	}
	if len(res.TxResponses) == 0 {
		return types.ContractConfig{}, fmt.Errorf("No transactions found for block %d, query %v", changedInBlock, query)
	}

	// Use the first matching tx we find, since results are in descending order.
	for _, txResponse := range res.TxResponses {
		if len(txResponse.Logs) == 0 {
			continue
		}
		for _, event := range txResponse.Logs[0].Events {
			if event.Type == "wasm-set_config" {
				cc, unknown, err := parseAttributes(event.Attributes)
				if len(unknown) > 0 {
					r.lggr.Warnf("wasm-set_config event contained unrecognized attributes: %v", unknown)
				}
				return cc, err
			}
		}
	}
	return types.ContractConfig{}, fmt.Errorf("No set_config event found in block %d", changedInBlock)
}

// parseAttributes returns a ContractConfig parsed from attrs.
// An error will be returned if any of the 8 required attributes are not present, or if any duplicates are found for
// unique attributes.
// unknownKeys contains counts of any unrecognized keys, which are otherwise ignored.
func parseAttributes(attrs []cosmosSDK.Attribute) (output types.ContractConfig, unknownKeys map[string]int, err error) {
	const uniqueKeys = 8
	known := make(map[string]struct{}, uniqueKeys)
	first := func(key string) bool {
		_, ok := known[key]
		if ok {
			return false
		}
		known[key] = struct{}{}
		return true
	}
	for _, attr := range attrs {
		key, value := attr.Key, attr.Value
		switch key {
		case "latest_config_digest":
			if !first(key) {
				err = ErrAttrDupe(key)
				return
			}
			// parse byte array encoded as hex string
			if err = HexToConfigDigest(value, &output.ConfigDigest); err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
		case "config_count":
			if !first(key) {
				err = ErrAttrDupe(key)
				return
			}
			var i int64
			i, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
			output.ConfigCount = uint64(i)
		case "signers":
			known[key] = struct{}{}
			// this assumes the value will be a hex encoded string which each signer 32 bytes and each signer will be a separate parameter
			var v []byte
			if err = HexToByteArray(value, &v); err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
			if len(v) != 32 {
				err = fmt.Errorf("failed to parse attribute %q: length '%d' != 32", key, len(v))
				return
			}
			output.Signers = append(output.Signers, v)
		case "transmitters":
			known[key] = struct{}{}
			// this assumes the return value be a string for each transmitter and each transmitter will be separate
			output.Transmitters = append(output.Transmitters, types.Account(attr.Value))
		case "f":
			if !first(key) {
				err = ErrAttrDupe(key)
				return
			}
			var i int64
			i, err = strconv.ParseInt(value, 10, 8)
			if err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
			output.F = uint8(i)
		case "onchain_config":
			if !first(key) {
				err = ErrAttrDupe(key)
				return
			}
			var config []byte
			// parse byte array encoded as base64
			config, err = base64.StdEncoding.DecodeString(value)
			if err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
			output.OnchainConfig = config
		case "offchain_config_version":
			if !first(key) {
				err = ErrAttrDupe(key)
				return
			}
			var i int64
			i, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
			output.OffchainConfigVersion = uint64(i)
		case "offchain_config":
			if !first(key) {
				err = ErrAttrDupe(key)
				return
			}
			var bytes []byte
			// parse byte array encoded as base64
			bytes, err = base64.StdEncoding.DecodeString(value)
			if err != nil {
				err = &ErrAttrInvalid{Err: err, Key: key}
				return
			}
			output.OffchainConfig = bytes
		default:
			if unknownKeys == nil {
				unknownKeys = make(map[string]int)
			}
			unknownKeys[key]++
		}
	}
	if len(known) != uniqueKeys {
		err = fmt.Errorf("expected %d types of known keys, but found %d: %v", uniqueKeys, len(known), known)
	}
	return
}

// ErrAttrInvalid is returned when parsing fails.
type ErrAttrInvalid struct {
	Key string
	Err error
}

func (e *ErrAttrInvalid) Error() string {
	return fmt.Sprintf("failed to parse attribute %q: %s", e.Key, e.Err.Error())
}

func (e *ErrAttrInvalid) Unwrap() error { return e.Err }

// ErrAttrDupe is returned when a duplicate attribute is found for a unique key.
type ErrAttrDupe string

func (e ErrAttrDupe) Error() string {
	return fmt.Sprintf("duplicate attributes for %q", string(e))
}

// LatestTransmissionDetails fetches the latest transmission details from address state
func (r *OCR2Reader) LatestTransmissionDetails(ctx context.Context) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	resp, err := r.chainReader.ContractState(r.address, []byte(`{"latest_transmission_details":{}}`))
	if err != nil {
		// Handle the 500 error that occurs when there has not been a submission
		// "rpc error: code = Unknown desc = ocr2::state::Transmission not found: contract query failed: unknown request"
		// which is thrown if this map lookup fails https://github.com/smartcontractkit/chainlink-cosmos/blob/main/contracts/ocr2/src/contract.rs#L759
		if strings.Contains(fmt.Sprint(err), "ocr2::state::Transmission not found") {
			r.lggr.Infof("No transmissions found when fetching `latest_transmission_details` attempting with `latest_config_digest_and_epoch`")
			digest, epoch, err2 := r.LatestConfigDigestAndEpoch(ctx)

			// In the case that there have been no transmissions, we expect the epoch to be zero.
			// We return just the contract digest here and set the rest of the
			// transmission details to their zero value.
			if err2 == nil {
				if epoch != 0 {
					r.lggr.Errorf("unexpected non-zero epoch %v and no transmissions found contract %v", epoch, r.address)
				}
				return digest, epoch, 0, big.NewInt(0), time.Unix(0, 0), nil
			}
			r.lggr.Errorf("error reading latest config digest and epoch err %v contract %v", err2, r.address)
		}

		// default response if there actually is an error
		return types.ConfigDigest{}, 0, 0, big.NewInt(0), time.Now(), err
	}

	// unmarshal
	var details LatestTransmissionDetails
	if err := json.Unmarshal(resp, &details); err != nil {
		return types.ConfigDigest{}, 0, 0, big.NewInt(0), time.Now(), err
	}

	// set answer big int
	ans := new(big.Int)
	if _, success := ans.SetString(details.LatestAnswer, 10); !success {
		return types.ConfigDigest{}, 0, 0, big.NewInt(0), time.Now(), fmt.Errorf("Could not create *big.Int from %s", details.LatestAnswer)
	}

	return details.LatestConfigDigest, details.Epoch, details.Round, ans, time.Unix(details.LatestTimestamp, 0), nil
}

// LatestRoundRequested fetches the latest round requested by filtering event logs
//func (cc *OCR2Reader) LatestRoundRequested(ctx context.Context, lookback time.Duration) (
//	configDigest types.ConfigDigest,
//	epoch uint32,
//	round uint8,
//	err error,
//) {
//	// calculate start block
//	latestBlock, blkErr := cc.chainReader.LatestBlock()
//	if blkErr != nil {
//		err = blkErr
//		return
//	}
//	blockNum := uint64(latestBlock.Block.Header.Height) - uint64(lookback/cc.cfg.BlockRate())
//	res, err := cc.chainReader.TxsEvents([]string{fmt.Sprintf("tx.height>=%d", blockNum+1), fmt.Sprintf("wasm-new_round.contract_address='%s'", cc.address.String())}, nil)
//	if err != nil {
//		return
//	}
//	if len(res.TxResponses) == 0 {
//		return
//	}
//	if len(res.TxResponses[0].Logs) == 0 {
//		err = fmt.Errorf("No logs found for tx %s", res.TxResponses[0].TxHash)
//		return
//	}
//	// First tx is the latest.
//	if len(res.TxResponses[0].Logs[0].Events) == 0 {
//		err = fmt.Errorf("No events found for tx %s", res.TxResponses[0].TxHash)
//		return
//	}
//
//	for _, event := range res.TxResponses[0].Logs[0].Events {
//		if event.Type == "wasm-new_round" {
//			// TODO: confirm event parameters
//			// https://github.com/smartcontractkit/chainlink-cosmos/issues/22
//			for _, attr := range event.Attributes {
//				key, value := string(attr.Key), string(attr.Value)
//				switch key {
//				case "latest_config_digest":
//					// parse byte array encoded as hex string
//					if err := HexToConfigDigest(value, &configDigest); err != nil {
//						return configDigest, epoch, round, err
//					}
//				case "epoch":
//					epochU64, err := strconv.ParseUint(value, 10, 32)
//					if err != nil {
//						return configDigest, epoch, round, err
//					}
//					epoch = uint32(epochU64)
//				case "round":
//					roundU64, err := strconv.ParseUint(value, 10, 8)
//					if err != nil {
//						return configDigest, epoch, round, err
//					}
//					round = uint8(roundU64)
//				}
//			}
//			return // exit once all parameters are processed
//		}
//	}
//	return
//}

// LatestConfigDigestAndEpoch fetches the latest details from address state
func (r *OCR2Reader) LatestConfigDigestAndEpoch(ctx context.Context) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	resp, err := r.chainReader.ContractState(
		r.address, []byte(`{"latest_config_digest_and_epoch":{}}`),
	)
	if err != nil {
		return types.ConfigDigest{}, 0, err
	}

	var digest LatestConfigDigestAndEpoch
	if err := json.Unmarshal(resp, &digest); err != nil {
		return types.ConfigDigest{}, 0, err
	}

	return digest.ConfigDigest, digest.Epoch, nil
}
