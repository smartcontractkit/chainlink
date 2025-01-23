package deployment

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// OCRSecrets are used to disseminate a shared secret to OCR nodes
// through the blockchain where OCR configuration is stored. Its a low value secret used
// to derive transmission order etc. They are extracted here such that they can common
// across signers when multiple signers are signing the same OCR config.
type OCRSecrets struct {
	SharedSecret [16]byte
	EphemeralSk  [32]byte
}

func (s OCRSecrets) IsEmpty() bool {
	return s.SharedSecret == [16]byte{} || s.EphemeralSk == [32]byte{}
}

func XXXGenerateTestOCRSecrets() OCRSecrets {
	var s OCRSecrets
	copy(s.SharedSecret[:], crypto.Keccak256([]byte("shared"))[:16])
	copy(s.EphemeralSk[:], crypto.Keccak256([]byte("ephemeral")))
	return s
}

// SimTransactOpts is useful to generate just the calldata for a given gethwrapper method.
func SimTransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		return transaction, nil
	}, From: common.HexToAddress("0x0"), NoSend: true, GasLimit: 1_000_000}
}

func GetErrorReasonFromTx(client bind.ContractBackend, from common.Address, tx *types.Transaction, receipt *types.Receipt) (string, error) {
	call := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Data:     tx.Data(),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	_, err := client.CallContract(context.Background(), call, receipt.BlockNumber)
	if err != nil {
		errorReason, err := parseError(err)
		if err == nil {
			return errorReason, nil
		}
	}
	return "", fmt.Errorf("tx %s reverted with no reason", tx.Hash().Hex())
}

func parseError(txError error) (string, error) {
	b, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	if json.Unmarshal(b, &callErr) != nil {
		return "", err
	}

	if callErr.Data == "" && strings.Contains(callErr.Message, "missing trie node") {
		return "", errors.Errorf("please use an archive node")
	}

	return callErr.Data, nil
}

func parseErrorFromABI(errorString string, contractABI string) (string, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return "", errors.Wrap(err, "error loading ABI")
	}
	errorString = strings.TrimPrefix(errorString, "Reverted ")
	errorString = strings.TrimPrefix(errorString, "0x")

	data, err := hex.DecodeString(errorString)
	if err != nil {
		return "", errors.Wrap(err, "error decoding error string")
	}
	for errorName, abiError := range parsedAbi.Errors {
		if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
			// Found a matching error
			v, err3 := abiError.Unpack(data)
			if err3 != nil {
				return "", errors.Wrap(err3, "error unpacking data")
			}
			return fmt.Sprintf("error -`%v` args %v", errorName, v), nil
		}
	}
	return "", errors.New("error not found in ABI")
}

// DecodeErr decodes an error from a contract call using the contract's ABI.
// If the error is not decodable, it returns the original error.
func DecodeErr(encodedABI string, err error) error {
	if err == nil {
		return nil
	}
	//revive:disable
	var d rpc.DataError
	ok := errors.As(err, &d)
	if ok {
		encErr, ok := d.ErrorData().(string)
		if !ok {
			return fmt.Errorf("error without error data: %s", d.Error())
		}
		errStr, parseErr := parseErrorFromABI(encErr, encodedABI)
		if parseErr != nil {
			return fmt.Errorf("failed to decode error '%s' with abi: %w", encErr, parseErr)
		}
		return fmt.Errorf("contract error: %s", errStr)
	}
	return fmt.Errorf("cannot decode error with abi: %w", err)
}

// ContractDeploy represents the result of an EVM contract deployment
// via an abigen Go binding. It contains all the return values
// as they are useful in different ways.
type ContractDeploy[C any] struct {
	Address  common.Address     // We leave this incase a Go binding doesn't have Address()
	Contract C                  // Expected to be a Go binding
	Tx       *types.Transaction // Incase the caller needs for example tx hash info for
	Tv       TypeAndVersion
	Err      error
}

// DeployContract deploys an EVM contract and
// records the address in the provided address book
// if the deployment was confirmed onchain.
// Deploying and saving the address is a very common pattern
// so this helps to reduce boilerplate.
// It returns an error if the deployment failed, the tx was not
// confirmed or the address could not be saved.
func DeployContract[C any](
	lggr logger.Logger,
	chain Chain,
	addressBook AddressBook,
	deploy func(chain Chain) ContractDeploy[C],
) (*ContractDeploy[C], error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "chain", chain.String(), "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	_, err := chain.Confirm(contractDeploy.Tx)
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "chain", chain.String(), "Contract", contractDeploy.Tv.String(), "err", err)
		return nil, err
	}
	lggr.Infow("Deployed contract", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.Address, "chain", chain.String())
	err = addressBook.Save(chain.Selector, contractDeploy.Address.String(), contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.Address, "chain", chain.String(), "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}

func IsValidChainSelector(cs uint64) error {
	if cs == 0 {
		return errors.New("chain selector must be set")
	}
	_, err := chain_selectors.GetSelectorFamily(cs)
	if err != nil {
		return err
	}
	return nil
}

func ChainInfo(cs uint64) (chain_selectors.ChainDetails, error) {
	id, err := chain_selectors.GetChainIDFromSelector(cs)
	if err != nil {
		return chain_selectors.ChainDetails{}, err
	}
	family, err := chain_selectors.GetSelectorFamily(cs)
	if err != nil {
		return chain_selectors.ChainDetails{}, err
	}
	info, err := chain_selectors.GetChainDetailsByChainIDAndFamily(id, family)
	if err != nil {
		return chain_selectors.ChainDetails{}, err
	}
	return info, nil
}

func ValidateSelectorsInEnvironment(e Environment, chains []uint64) error {
	for _, chain := range chains {
		_, evmOk := e.Chains[chain]
		_, solOk := e.SolChains[chain]
		if !evmOk && !solOk {
			return fmt.Errorf("chain %d not found in environment", chain)
		}
	}
	return nil
}
