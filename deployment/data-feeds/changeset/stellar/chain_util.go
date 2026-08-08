package stellar

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/stellar/go-stellar-sdk/xdr"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-stellar/bindings"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	proxy "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_proxy"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"
)

type StellarDeps struct {
	Deploy  SorobanContractDeployer
	Invoker bindings.Invoker
}

type SorobanContractDeployer interface {
	DeployContractWithArgs(ctx context.Context, wasmPath string, salt [32]byte, ctorArgs []xdr.ScVal) (string, error)
	UploadContractWASM(ctx context.Context, wasmPath string) (xdr.Hash, error)
}

type void struct{}

var opVersion = semver.MustParse("1.0.0")

// generated cache and proxy clients.
type contractAdmin interface {
	TransferOwnership(ctx context.Context, newOwner string, liveUntilLedger uint32) error
	AcceptOwnership(ctx context.Context) error
	Upgrade(ctx context.Context, newWasmHash [32]byte) error
	RecoverTokens(ctx context.Context, token, to string, amount int64) error
}

func adminClient(d StellarDeps, contractID string, isProxy bool) contractAdmin {
	if isProxy {
		return proxy.NewDataFeedsProxyClient(d.Invoker, contractID)
	}
	return cache.NewDataFeedsCacheClient(d.Invoker, contractID)
}

var newStellarDeps = func(ch cldfstellar.Chain) (StellarDeps, error) {
	d, err := stellardeploy.NewDeployerFromChain(ch)
	if err != nil {
		return StellarDeps{}, err
	}
	return StellarDeps{Deploy: d, Invoker: d}, nil
}

func chainDeps(env cldf.Environment, chainSel uint64) (cldfstellar.Chain, StellarDeps, error) {
	ch, ok := env.BlockChains.StellarChains()[chainSel]
	if !ok {
		return cldfstellar.Chain{}, StellarDeps{}, fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
	}
	deps, err := newStellarDeps(ch)
	return ch, deps, err
}

func ownerOrSigner(ch cldfstellar.Chain, owner string) (string, error) {
	if owner != "" {
		return owner, nil
	}
	if ch.Signer == nil {
		return "", errors.New("owner not set and chain has no signer")
	}
	return ch.Signer.Address(), nil
}

type stellarApplyDeps struct {
	deps       StellarDeps
	contractID string
}

func verifyContractRef(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) error {
	if _, ok := env.BlockChains.StellarChains()[chainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
	}
	_, err := getAddressRef(env, chainSel, contractType, qualifier, version)
	return err
}

func getAddressRef(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (datastore.AddressRef, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("invalid version %q: %w", version, err)
	}
	ref, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(chainSel, contractType, v, qualifier),
	)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("%s address ref not found for qualifier %q: %w", contractType, qualifier, err)
	}
	return ref, nil
}

func resolveContractDeps(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (stellarApplyDeps, datastore.AddressRef, error) {
	_, deps, err := chainDeps(env, chainSel)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, err
	}
	ref, err := getAddressRef(env, chainSel, contractType, qualifier, version)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, err
	}
	return stellarApplyDeps{deps: deps, contractID: ref.Address}, ref, nil
}

// dataIDsToBytes converts hex feed IDs to [16]byte. Data IDs are canonically
// left-aligned, so short values are left-justified with trailing zero padding.
func dataIDsToBytes(ids []string) ([][16]byte, error) {
	out := make([][16]byte, 0, len(ids))
	for _, id := range ids {
		v, ok := new(big.Int).SetString(id, 0)
		if !ok {
			return nil, fmt.Errorf("invalid data_id: %q", id)
		}
		if v.BitLen() > 128 {
			return nil, fmt.Errorf("data_id too long: %q (%d bits)", id, v.BitLen())
		}
		var b [16]byte
		copy(b[:], v.Bytes())
		out = append(out, b)
	}
	return out, nil
}

// workflowNameToBytes right-pads an ASCII workflow name into [10]byte.
func workflowNameToBytes(s string) ([10]byte, error) {
	var out [10]byte
	if len(s) > len(out) {
		return out, fmt.Errorf("workflow name %q exceeds %d bytes", s, len(out))
	}
	copy(out[:], s)
	return out, nil
}

// workflowOwnerToBytes decodes a 20-byte hex workflow owner.
func workflowOwnerToBytes(hexStr string) ([20]byte, error) {
	var out [20]byte
	b, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	if err != nil {
		return out, fmt.Errorf("invalid workflow owner %q: %w", hexStr, err)
	}
	if len(b) != len(out) {
		return out, fmt.Errorf("workflow owner must be %d bytes, got %d", len(out), len(b))
	}
	copy(out[:], b)
	return out, nil
}
