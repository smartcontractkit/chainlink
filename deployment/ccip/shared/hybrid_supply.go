package shared

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	aptoslib "github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"
)

// maxDecimalExponent is the largest power of ten that fits in a uint256, since
// 10^77 < 2^256-1 < 10^78.
//
// Scaling up multiplies by 10^exponent, so an exponent above this cannot produce an
// encodable amount. Scaling down only divides, where an exponent this large already
// yields a remainder and is rejected as truncation; the bound is kept there to avoid
// building a pointlessly large divisor.
const maxDecimalExponent = 77

// RemoteSupply is a token supply expressed in a remote chain's own decimal domain.
//
// A group update's RemoteChainSupply is denominated in the LOCAL token's base units,
// because updateGroups mints it directly through the external minter. Stating the remote
// supply separately lets the changeset derive the local amount and reject one that was
// passed through without conversion.
type RemoteSupply struct {
	// Decimals is the token's decimal precision on the remote chain.
	Decimals uint8 `json:"decimals"`

	// Amount is the supply in the remote chain's smallest indivisible unit, exactly as the
	// remote chain reports it, with no decimal point applied. For a 6-decimal token holding
	// 1,206,841,810.125844 tokens this is 1206841810125844.
	//
	// Read it from totalSupply() on EVM, getTokenSupply on SVM, or the fungible asset's
	// supply on Aptos.
	Amount *big.Int `json:"amount"`
}

// NormalizeRemoteSupply converts an amount from a remote chain's decimal domain into the
// local one, mirroring TokenPool._calculateLocalAmount so that migration backing and later
// redemptions agree.
//
// Scaling down must be exact: a remainder means the remote supply has no local
// representation, and truncating it would under-back the pool.
func NormalizeRemoteSupply(amount *big.Int, remoteDecimals, localDecimals uint8) (*big.Int, error) {
	if amount == nil {
		return nil, errors.New("remote supply amount must not be nil")
	}

	if amount.Sign() < 0 {
		return nil, fmt.Errorf("remote supply amount %s must not be negative", amount)
	}

	if remoteDecimals == localDecimals {
		return new(big.Int).Set(amount), nil
	}

	if remoteDecimals > localDecimals {
		return scaleDown(amount, remoteDecimals, localDecimals)
	}

	return scaleUp(amount, remoteDecimals, localDecimals)
}

func scaleDown(amount *big.Int, remoteDecimals, localDecimals uint8) (*big.Int, error) {
	exponent := remoteDecimals - localDecimals
	if exponent > maxDecimalExponent {
		return nil, fmt.Errorf("cannot convert from %d to %d decimals, the scale factor exceeds uint256", remoteDecimals, localDecimals)
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)

	quotient, remainder := new(big.Int).QuoRem(amount, divisor, new(big.Int))
	if remainder.Sign() != 0 {
		return nil, fmt.Errorf("remote supply %s at %d decimals has no exact representation at %d local decimals", amount, remoteDecimals, localDecimals)
	}

	return quotient, nil
}

func scaleUp(amount *big.Int, remoteDecimals, localDecimals uint8) (*big.Int, error) {
	exponent := localDecimals - remoteDecimals
	if exponent > maxDecimalExponent {
		return nil, fmt.Errorf("cannot convert from %d to %d decimals, the scale factor exceeds uint256", remoteDecimals, localDecimals)
	}

	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)

	localAmount := new(big.Int).Mul(amount, multiplier)
	if localAmount.BitLen() > 256 {
		return nil, fmt.Errorf("remote supply %s at %d decimals overflows uint256 at %d local decimals", amount, remoteDecimals, localDecimals)
	}

	return localAmount, nil
}

// ValidateRemoteChainSupply checks that suppliedLocalAmount is remoteSupply expressed in
// local base units.
//
// This rejects a remote supply passed straight into a group update: with a 6-decimal remote
// and an 18-decimal local token, the unconverted value under-backs the pool by 10^12.
func ValidateRemoteChainSupply(
	remoteSupply RemoteSupply,
	suppliedLocalAmount *big.Int,
	localDecimals uint8,
	remoteChainSelector uint64,
) error {
	if suppliedLocalAmount == nil {
		return fmt.Errorf("remote chain %d: RemoteChainSupply must not be nil", remoteChainSelector)
	}

	expected, err := NormalizeRemoteSupply(remoteSupply.Amount, remoteSupply.Decimals, localDecimals)
	if err != nil {
		return fmt.Errorf("remote chain %d: %w", remoteChainSelector, err)
	}

	if expected.Cmp(suppliedLocalAmount) != 0 {
		return fmt.Errorf(
			"remote chain %d: RemoteChainSupply must be in local base units (%d decimals), so remote supply %s at %d decimals requires %s, but the update supplies %s",
			remoteChainSelector, localDecimals,
			remoteSupply.Amount, remoteSupply.Decimals, expected, suppliedLocalAmount,
		)
	}

	return nil
}

// VerifyRemoteDecimals cross-checks a supplied remote decimal precision against the remote
// chain, and reports whether that check could actually be made.
//
// Decimals are not stored on the pool, so each family is read through its own interface.
// verified is false when the remote chain is absent from the environment or its family is
// unsupported, which leaves the supplied value as the only source of truth; the caller is
// told rather than given a silent pass.
func VerifyRemoteDecimals(
	ctx context.Context,
	env cldf.Environment,
	remoteChainSelector uint64,
	remoteTokenAddress []byte,
	declaredDecimals uint8,
) (verified bool, err error) {
	if len(remoteTokenAddress) == 0 {
		return false, nil
	}

	family, err := chain_selectors.GetSelectorFamily(remoteChainSelector)
	if err != nil {
		return false, fmt.Errorf("failed to resolve chain family for remote chain %d: %w", remoteChainSelector, err)
	}

	var onchainDecimals uint8

	switch family {
	case chain_selectors.FamilyEVM:
		onchainDecimals, verified, err = evmTokenDecimals(ctx, env, remoteChainSelector, remoteTokenAddress)
	case chain_selectors.FamilySolana:
		onchainDecimals, verified, err = solanaMintDecimals(ctx, env, remoteChainSelector, remoteTokenAddress)
	case chain_selectors.FamilyAptos:
		onchainDecimals, verified, err = aptosFungibleAssetDecimals(env, remoteChainSelector, remoteTokenAddress)
	default:
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if !verified {
		return false, nil
	}

	if onchainDecimals != declaredDecimals {
		return false, fmt.Errorf("remote chain %d: supplied %d decimals but the remote token reports %d", remoteChainSelector, declaredDecimals, onchainDecimals)
	}

	return true, nil
}

func evmTokenDecimals(
	ctx context.Context,
	env cldf.Environment,
	remoteChainSelector uint64,
	remoteTokenAddress []byte,
) (uint8, bool, error) {
	remoteChain, ok := env.BlockChains.EVMChains()[remoteChainSelector]
	if !ok {
		return 0, false, nil
	}

	// EVM token addresses are left-padded to 32 bytes on the pool.
	tokenAddress := common.BytesToAddress(remoteTokenAddress)
	if tokenAddress == (common.Address{}) {
		return 0, false, nil
	}

	token, err := erc20.NewERC20(tokenAddress, remoteChain.Client)
	if err != nil {
		return 0, false, fmt.Errorf("failed to bind erc20 at %s on %s: %w", tokenAddress, remoteChain.String(), err)
	}

	decimals, err := token.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, false, fmt.Errorf("failed to read decimals from token %s on %s: %w", tokenAddress, remoteChain.String(), err)
	}

	return decimals, true, nil
}

// solanaMintDecimals reads the decimals field of the SPL mint account.
func solanaMintDecimals(
	ctx context.Context,
	env cldf.Environment,
	remoteChainSelector uint64,
	remoteTokenAddress []byte,
) (uint8, bool, error) {
	remoteChain, ok := env.BlockChains.SolanaChains()[remoteChainSelector]
	if !ok {
		return 0, false, nil
	}

	if len(remoteTokenAddress) != solana.PublicKeyLength {
		return 0, false, fmt.Errorf("remote chain %d: expected a %d byte SVM mint address, got %d bytes", remoteChainSelector, solana.PublicKeyLength, len(remoteTokenAddress))
	}

	mintAddress := solana.PublicKeyFromBytes(remoteTokenAddress)

	var mint solToken.Mint
	if err := remoteChain.GetAccountDataBorshInto(ctx, mintAddress, &mint); err != nil {
		return 0, false, fmt.Errorf("failed to read mint %s on %s: %w", mintAddress, remoteChain.String(), err)
	}

	return mint.Decimals, true, nil
}

// aptosFungibleAssetDecimals calls 0x1::fungible_asset::decimals on the asset's metadata
// object, which is what the pool stores as the remote token address.
func aptosFungibleAssetDecimals(
	env cldf.Environment,
	remoteChainSelector uint64,
	remoteTokenAddress []byte,
) (uint8, bool, error) {
	remoteChain, ok := env.BlockChains.AptosChains()[remoteChainSelector]
	if !ok {
		return 0, false, nil
	}

	var metadataAddress aptoslib.AccountAddress
	if len(remoteTokenAddress) != len(metadataAddress) {
		return 0, false, fmt.Errorf("remote chain %d: expected a %d byte Aptos metadata address, got %d bytes", remoteChainSelector, len(metadataAddress), len(remoteTokenAddress))
	}

	copy(metadataAddress[:], remoteTokenAddress)

	values, err := remoteChain.Client.View(&aptoslib.ViewPayload{
		Module:   aptoslib.ModuleId{Address: aptoslib.AccountOne, Name: "fungible_asset"},
		Function: "decimals",
		ArgTypes: []aptoslib.TypeTag{{Value: &aptoslib.StructTag{
			Address: aptoslib.AccountOne,
			Module:  "fungible_asset",
			Name:    "Metadata",
		}}},
		Args: [][]byte{metadataAddress[:]},
	})
	if err != nil {
		return 0, false, fmt.Errorf("failed to read decimals for fungible asset %s on %s: %w", metadataAddress.String(), remoteChain.String(), err)
	}

	if len(values) == 0 {
		return 0, false, fmt.Errorf("remote chain %d: decimals view returned no value for fungible asset %s", remoteChainSelector, metadataAddress.String())
	}

	// Aptos view results arrive as JSON numbers.
	decimals, ok := values[0].(float64)
	if !ok {
		return 0, false, fmt.Errorf("remote chain %d: decimals view returned %T, want a number", remoteChainSelector, values[0])
	}

	if decimals < 0 || decimals > math.MaxUint8 {
		return 0, false, fmt.Errorf("remote chain %d: decimals view returned %v, which is out of range", remoteChainSelector, decimals)
	}

	return uint8(decimals), true, nil
}
