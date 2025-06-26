//go:build wasip1

package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	httpaction "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http"
	evmcap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/chains/evm"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/cmd/cre/examples/v2/por/bindings"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
)

type EvmConfig struct {
	TokenAddress  string
	PorAddress    string
	ChainSelector uint32
	GasLimit      uint64
}

type Config struct {
	PublicKey string
	Schedule  string
	Url       string
	Evms      []EvmConfig
}

type ReserveInfo struct {
	LastUpdated  time.Time       `consensus_aggregation:"median" json:"lastUpdated"`
	TotalReserve decimal.Decimal `consensus_aggregation:"median" json:"totalReserve"`
}

type PorResponse struct {
	DataSignature string `json:"dataSignature"`
	Ripcord       bool   `json:"ripcord"`
	Data          string `json:"data"`
}

func RunSimpleCronWorkflow(runner sdk.Runner[[]byte]) {
	// cron := &croncap.Cron{}
	cfg := &croncap.Config{
		Schedule: "*/3 * * * * *", // every three seconds
	}

	// runner.Run(&sdk.[sdk.DonRuntime]{
	// 	Handlers: []sdk.Handler[sdk.DonRuntime]{
	// 		sdk.New(
	// 			cron.Trigger(cfg),
	// 			onTrigger,
	// 		),
	// 	},
	// })

	runner.Run(func(env *sdk.Environment[[]byte]) (sdk.Workflow[[]byte], error) {
		return sdk.Workflow[[]byte]{
			sdk.Handler(
				croncap.Trigger(cfg),
				onTrigger),
		}, nil
	})
}

func onTrigger(env *sdk.Environment[[]byte], runtime sdk.Runtime, outputs *croncap.Payload) (string, error) {
	// TODO: Fix
	return doPor(env, runtime, outputs.ScheduledExecutionTime.AsTime(), "http://localhost:3000", "publicKey", []EvmConfig{
		{
			TokenAddress:  "0x9b41EB05aC02c4fBD5eFFa657c627BfA1dC8f2e6",
			PorAddress:    "0x2f4f914826fEf265345A9752fa6B113594E4DD8b",
			ChainSelector: 1,
			GasLimit:      1000000,
		},
	})
}

func doPor(env *sdk.Environment[[]byte], runtime sdk.Runtime, runTime time.Time, url string, publicKey string, evms []EvmConfig) (string, error) {
	// Fetch Por
	env.Logger.Info("fetching por", "url", url, "publicKey", publicKey, "evms", evms)
	reserveInfo, err := sdk.RunInNodeMode(env, runtime, func(env *sdk.NodeEnvironment[[]byte], nodeRuntime sdk.NodeRuntime) (*ReserveInfo, error) {
		reserveInfo, err := fetchPor(url, publicKey, nodeRuntime)
		if err != nil {
			env.Logger.Error("error fetching por", "err", err)
			return nil, err
		}
		return reserveInfo, nil
	}, sdk.ConsensusAggregationFromTags[*ReserveInfo]()).Await()
	if err != nil {
		return "", err
	}

	env.Logger.Info("ReserveInfo", reserveInfo)

	if reserveInfo.LastUpdated.Before(runTime.Add(-time.Hour * 24)) {
		env.Logger.Warn("reserve time is too old", "time", reserveInfo.LastUpdated)
		// return "", errors.New("reserved time is too old")
	}

	// TODO: Make this work
	totalSupply, err := getTotalSupply(env, runtime, evms)
	if err != nil {
		return "", err
	}

	env.Logger.Info("TotalSupply", totalSupply)
	totalReserveScaled := reserveInfo.TotalReserve.Mul(decimal.NewFromUint64(1e18)).BigInt()
	env.Logger.Info("TotalReserveScaled", totalReserveScaled)

	if err = updateReserve(env, runtime, totalSupply, totalReserveScaled, evms); err != nil {
		return "", err
	}

	return reserveInfo.TotalReserve.String(), nil
}

func updateReserve(env *sdk.Environment[[]byte], runtime sdk.Runtime, totalSupply, totalReserveScaled *big.Int, evms []EvmConfig) error {
	reportWrites := make([]sdk.Promise[*evmcap.WriteReportReply], len(evms))
	for i, evmConfig := range evms {
		evmClient := &evmcap.Client{}

		// Address must be parsable or the workflow would fail to initialize the trigger.
		address, _ := hexToBytes(evmConfig.PorAddress)
		reserveManager := bindings.NewIReserveManager(bindings.ContractInputs{EVM: evmClient, Address: address, Options: &bindings.ContractOptions{
			GasConfig: &evm.GasConfig{
				GasLimit: evmConfig.GasLimit,
			},
		}})
		reportWrites[i] = reserveManager.WriteReportUpdateReserves(runtime, bindings.UpdateReservesStruct{
			TotalMinted:  totalSupply,
			TotalReserve: totalReserveScaled,
		}, nil)
	}

	var errs []error
	for i, promise := range reportWrites {
		_, err := promise.Await()
		if err == nil {
			// runtime.Logger().Debug("update reserve write report reply", "chain_selector", evms[i].ChainSelector, "tx hash", writeReportReply.TxHash)
		} else {
			selector := evms[i].ChainSelector
			env.Logger.Error("Could not write to contract", "contract_chain", selector, "err", err.Error())
			errs = append(errs, fmt.Errorf("failed to write report for chain %d: %w", selector, err))
			continue
		}
	}

	return errors.Join(errs...)
}

func getTotalSupply(env *sdk.Environment[[]byte], runtime sdk.Runtime, evms []EvmConfig) (*big.Int, error) {
	// Fetch supply from all EVMs in parallel
	supplyPromises := make([]sdk.Promise[*big.Int], len(evms))
	for i, evmConfig := range evms {
		evmClient := &evmcap.Client{}

		address, err := hexToBytes(evmConfig.TokenAddress)
		if err != nil {
			env.Logger.Error("failed to decode token address", "address", evmConfig.TokenAddress, "err", err)
			return nil, fmt.Errorf("failed to decode token address %s: %w", evmConfig.TokenAddress, err)
		}
		token := bindings.NewIERC20(bindings.ContractInputs{EVM: evmClient, Address: address})
		evmTotalSupplyPromise := token.Methods.TotalSupply.Call(runtime, nil)
		supplyPromises[i] = evmTotalSupplyPromise
	}

	// We can add sdk.AwaitAll that takes []sdk.Promise[T] and returns ([]T, error)
	totalSupply := big.NewInt(0)
	for _, promise := range supplyPromises {
		supply, err := promise.Await()
		if err != nil {
			// selector := evms[i].ChainSelector
			// runtime.Logger().Error("Could not read from contract", "contract_chain", selector, "err", err.Error())
			return nil, err
		}

		totalSupply = totalSupply.Add(totalSupply, supply)
	}

	return totalSupply, nil
}

func fetchPor(urlString string, publicKey string, runtime sdk.NodeRuntime) (*ReserveInfo, error) {
	httpAction := httpaction.Client{}
	httpActionOut, err := httpAction.SendRequest(runtime, &httpaction.Request{
		Method: "GET",
		Url:    urlString,
	}).Await()

	if err != nil {
		return nil, err
	}

	porResponse := &PorResponse{}
	if err = json.Unmarshal(httpActionOut.Body, porResponse); err != nil {
		return nil, err
	}

	// TODO: Make this work
	// err = verifySignature(porResponse, publicKey)
	// if err != nil {
	// 	return nil, err
	// }

	if porResponse.Ripcord {
		return nil, errors.New("ripcord is true")
	}

	reserveInfo := &ReserveInfo{}
	if err = json.Unmarshal([]byte(porResponse.Data), reserveInfo); err != nil {
		return nil, err
	}

	return reserveInfo, nil
}

func verifySignature(porResponse *PorResponse, publicKey string) error {
	// Decode the signature
	rawSig, err := base64.RawURLEncoding.DecodeString(porResponse.DataSignature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Parse the PEM public key
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("invalid PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	// Hash the payload
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(porResponse.Data))
	digest := hasher.Sum(nil)

	// Verify
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest, rawSig); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	return nil
}

func hexToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr[2:])
}

func main() {
	RunSimpleCronWorkflow(wasm.NewRunner(func(configBytes []byte) ([]byte, error) {
		return configBytes, nil
	}))
}
