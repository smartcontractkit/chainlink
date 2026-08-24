// Package keyseed copies this node's keys into the keystore the CRE processes
// beside it read.
//
// The node keeps its keys in the legacy key ring; the p2p proxy and the
// capabilities it fronts read chainlink-common's keystore, which is a different
// table in the same database. They have to be the same keys: the proxy announces
// as this node's peer and signs rounds as this node's OCR identity, and the
// registry lists those public halves against this node.
//
// # This is a bootstrap, not a migration
//
// It converts by exporting each key the way the node already can - encrypted JSON,
// which is the supported way key material leaves the ring - decrypting it, and
// importing the raw key into the new store. That means it knows two things it
// would rather not: the password prefix each legacy export uses, and the byte
// layout inside it. Both are stable, and a change to either fails loudly here
// rather than quietly later, but neither is a contract.
//
// It exists because the two stores are not yet one. When the node's own keystore
// moves to chainlink-common's, this package goes: the keys will already be where
// the proxy looks, and nothing will need copying.
package keyseed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/keystore/ocr2offchain"
	"github.com/smartcontractkit/chainlink-common/keystore/pgstore"
	"github.com/smartcontractkit/chainlink-common/keystore/ragep2p"
	"github.com/smartcontractkit/chainlink-common/keystore/serialization"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// Names are what the copied keys are called in the new store, which addresses
// keys by name rather than by type.
//
// They must match what the reading side is configured with - crecore's
// --ocr.peer-key-name and --ocr.ocr-key-name - and these are the defaults on both
// sides.
type Names struct {
	Store string
	Peer  string
	OCR   string
}

// DefaultNames are the conventional names. They are duplicated rather than shared
// because the two sides are different repositories; the flags exist for a
// deployment that wants to say so once and differently.
var DefaultNames = Names{Store: "node", Peer: "p2p", OCR: "ocr2"}

// prefixOCR2Onchain and onchainSigning are where the onchain key sits. They are
// spelled here because the side that reads them is another repository (the
// capabilities repo's libs/standalone/nodekeys); the offchain and peer paths come
// from chainlink-common, whose helpers read those.
const (
	prefixOCR2Onchain = "ocr2_onchain"
	onchainSigning    = "ocr2_onchain_signing"
)

// Seed copies this node's P2P key, its EVM OCR2 bundle and its EVM keys into the
// keystore at names.Store, leaving anything already there alone.
//
// It is idempotent, and it is quiet about what it did not have to do: a node
// restarted after this ran once copies nothing.
func Seed(ctx context.Context, lggr logger.Logger, ks keystore.Master, ds sqlutil.DataSource, password string, names Names) error {
	if password == "" {
		return errors.New("the keystore password is required to copy this node's keys")
	}
	names = names.withDefaults()

	store, err := commonkeystore.LoadKeystore(ctx, emptyAsNew{pgstore.NewStorage(ds, names.Store)}, password)
	if err != nil {
		return fmt.Errorf("failed to open the keystore %q: %w", names.Store, err)
	}

	seeded := 0
	for _, seed := range []struct {
		what string
		fn   func() (int, error)
	}{
		{"the P2P key", func() (int, error) { return seedPeer(ctx, ks, store, password, names) }},
		{"the OCR2 keys", func() (int, error) { return seedOCR2(ctx, ks, store, password, names) }},
		{"the EVM keys", func() (int, error) { return seedEVM(ctx, ks, store, password) }},
	} {
		n, err := seed.fn()
		if err != nil {
			return fmt.Errorf("failed to copy %s into the keystore: %w", seed.what, err)
		}
		seeded += n
	}

	if seeded > 0 {
		lggr.Infow("Copied this node's keys into the CRE keystore", "keystore", names.Store, "keys", seeded)
	}
	return nil
}

func (n Names) withDefaults() Names {
	if n.Store == "" {
		n.Store = DefaultNames.Store
	}
	if n.Peer == "" {
		n.Peer = DefaultNames.Peer
	}
	if n.OCR == "" {
		n.OCR = DefaultNames.OCR
	}
	return n
}

// seedPeer copies the node's P2P key, which is the identity the proxy announces
// as. A node has one; if it somehow has several, the first is the one the proxy
// would have taken from the ring too.
func seedPeer(ctx context.Context, ks keystore.Master, store commonkeystore.Keystore, password string, names Names) (int, error) {
	keys, err := ks.P2P().GetAll()
	if err != nil {
		return 0, err
	}
	if len(keys) == 0 {
		return 0, nil
	}

	peer, err := peerKey(keys[0], password, names)
	if err != nil {
		return 0, err
	}
	return imported(ctx, store, password, peer)
}

// peerKey is the node's P2P key as the new store holds it.
func peerKey(legacy p2pkey.KeyV2, password string, names Names) (key, error) {
	raw, err := exported(func() ([]byte, error) {
		return legacy.ToEncryptedJSON(password, commonkeystore.DefaultScryptParams)
	}, func(b []byte) (gethkeystore.CryptoJSON, error) {
		var export p2pkey.EncryptedP2PKeyExport
		err := json.Unmarshal(b, &export)
		return export.Crypto, err
	}, "p2pkey"+password)
	if err != nil {
		return key{}, err
	}

	// The ring stores a libp2p-prefixed private key; what the new store wants is the
	// ed25519 private key itself, which is the last of those bytes.
	const ed25519PrivateKeySize = 64
	if len(raw) < ed25519PrivateKeySize {
		return key{}, fmt.Errorf("the P2P key is %d bytes, too short to hold an ed25519 private key", len(raw))
	}

	return key{
		name:       path(ragep2p.PrefixPeerKeyring, names.withDefaults().Peer),
		keyType:    commonkeystore.Ed25519,
		privateKey: raw[len(raw)-ed25519PrivateKeySize:],
	}, nil
}

// seedOCR2 copies the node's EVM OCR2 bundle, as the three keys the new store
// holds it as: the offchain signing and encryption keys under the name
// chainlink-common's helpers look for, and the onchain signing key beside them.
func seedOCR2(ctx context.Context, ks keystore.Master, store commonkeystore.Keystore, password string, names Names) (int, error) {
	bundles, err := ks.OCR2().GetAllOfType(corekeys.EVM)
	if err != nil {
		return 0, err
	}
	if len(bundles) == 0 {
		return 0, nil
	}

	keys, err := ocr2Keys(bundles[0], password, names)
	if err != nil {
		return 0, err
	}
	return imported(ctx, store, password, keys...)
}

// ocr2Keys is an OCR2 bundle as the three keys the new store holds it as: the
// offchain signing and encryption keys under the name chainlink-common's helpers
// look for, and the onchain signing key beside them.
func ocr2Keys(bundle ocr2key.KeyBundle, password string, names Names) ([]key, error) {
	names = names.withDefaults()

	raw, err := exported(func() ([]byte, error) {
		return ocr2key.ToEncryptedJSON(bundle, password, commonkeystore.DefaultScryptParams)
	}, func(b []byte) (gethkeystore.CryptoJSON, error) {
		var export ocr2key.EncryptedOCRKeyExport
		err := json.Unmarshal(b, &export)
		return export.Crypto, err
	}, "ocr2key"+password)
	if err != nil {
		return nil, err
	}

	// The bundle's own encoding: the two offchain keys as one blob, and the onchain
	// key beside it. Declared here rather than imported because the type that writes
	// it is not exported; the field names are what the encoding is.
	var encoded struct {
		OffchainKeyring []byte
		Keyring         []byte
	}
	if err := json.Unmarshal(raw, &encoded); err != nil {
		return nil, fmt.Errorf("failed to read the OCR2 bundle: %w", err)
	}

	// The offchain blob is an ed25519 private key followed by an x25519 scalar.
	const signingSize, encryptionSize = 64, 32
	if len(encoded.OffchainKeyring) != signingSize+encryptionSize {
		return nil, fmt.Errorf("the OCR2 offchain keys are %d bytes, want %d", len(encoded.OffchainKeyring), signingSize+encryptionSize)
	}
	if len(encoded.Keyring) == 0 {
		return nil, errors.New("the OCR2 bundle has no onchain key")
	}

	return []key{
		{
			name:       path(ocr2offchain.PrefixOCR2Offchain, names.OCR, ocr2offchain.OCR2OffchainSigning),
			keyType:    commonkeystore.Ed25519,
			privateKey: encoded.OffchainKeyring[:signingSize],
		},
		{
			name:       path(ocr2offchain.PrefixOCR2Offchain, names.OCR, ocr2offchain.OCR2OffchainEncryption),
			keyType:    commonkeystore.X25519,
			privateKey: encoded.OffchainKeyring[signingSize:],
		},
		{
			name:       path(prefixOCR2Onchain, names.OCR, onchainSigning),
			keyType:    commonkeystore.ECDSA_S256,
			privateKey: encoded.Keyring,
		},
	}, nil
}

// seedEVM copies the node's EVM keys, each named by its address.
//
// The name is the address because that is what asks for it: a chain capability
// transmits from an account, and the new store is addressed by name, so the two
// meet if the name is the account.
func seedEVM(ctx context.Context, ks keystore.Master, store commonkeystore.Keystore, password string) (int, error) {
	evmKeys, err := ks.Eth().GetAll(ctx)
	if err != nil {
		return 0, err
	}

	keys := make([]key, 0, len(evmKeys))
	for _, legacy := range evmKeys {
		k, err := evmKey(legacy, password)
		if err != nil {
			return 0, err
		}
		keys = append(keys, k)
	}
	return imported(ctx, store, password, keys...)
}

// evmKey is one EVM key as the new store holds it.
func evmKey(legacy ethkey.KeyV2, password string) (key, error) {
	privateKey, err := evmPrivateKey(legacy, password)
	if err != nil {
		return key{}, err
	}
	return key{
		name:       legacy.Address.Hex(),
		keyType:    commonkeystore.ECDSA_S256,
		privateKey: privateKey,
	}, nil
}

// evmPrivateKey is the 32 bytes an EVM key signs with. Its export is a standard
// geth keystore file, so this needs no knowledge of the ring's own encoding.
func evmPrivateKey(evmKey ethkey.KeyV2, password string) ([]byte, error) {
	exportJSON, err := evmKey.ToEncryptedJSON(password, commonkeystore.DefaultScryptParams)
	if err != nil {
		return nil, fmt.Errorf("failed to export the EVM key %s: %w", evmKey.Address, err)
	}
	decrypted, err := gethkeystore.DecryptKey(exportJSON, password)
	if err != nil {
		return nil, fmt.Errorf("failed to read the EVM key %s: %w", evmKey.Address, err)
	}
	return gethcrypto.FromECDSA(decrypted.PrivateKey), nil
}

// exported is the shape every legacy export shares: encrypted JSON whose envelope
// says what it is, and whose contents are the raw key under a password the export
// prefixed.
func exported(export func() ([]byte, error), crypto func([]byte) (gethkeystore.CryptoJSON, error), password string) ([]byte, error) {
	exportJSON, err := export()
	if err != nil {
		return nil, fmt.Errorf("failed to export the key: %w", err)
	}
	cryptoJSON, err := crypto(exportJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to read the exported key: %w", err)
	}
	raw, err := gethkeystore.DecryptDataV3(cryptoJSON, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt the exported key: %w", err)
	}
	return raw, nil
}

// key is one key on its way into the new store.
type key struct {
	name       string
	keyType    commonkeystore.KeyType
	privateKey []byte
}

// imported puts keys into the store, skipping any that are already there, and
// returns how many it added.
//
// Skipping rather than replacing is what makes this safe to run on every boot: the
// store is the reading side's, and a key it already has is one it may already have
// signed with.
func imported(ctx context.Context, store commonkeystore.Keystore, password string, keys ...key) (int, error) {
	requests := make([]commonkeystore.ImportKeyRequest, 0, len(keys))
	for _, k := range keys {
		present, err := has(ctx, store, k.name)
		if err != nil {
			return 0, err
		}
		if present {
			continue
		}

		encrypted, err := encryptedKey(k, password)
		if err != nil {
			return 0, err
		}
		requests = append(requests, commonkeystore.ImportKeyRequest{
			NewKeyName: k.name,
			Data:       encrypted,
			Password:   password,
		})
	}
	if len(requests) == 0 {
		return 0, nil
	}

	if _, err := store.ImportKeys(ctx, commonkeystore.ImportKeysRequest{Keys: requests}); err != nil {
		return 0, err
	}
	return len(requests), nil
}

// encryptedKey wraps a raw key the way ImportKeys reads it: the key as
// chainlink-common serialises it, encrypted as a geth keystore file.
func encryptedKey(k key, password string) ([]byte, error) {
	serialized, err := proto.Marshal(&serialization.Key{
		Name:       k.name,
		KeyType:    string(k.keyType),
		PrivateKey: k.privateKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to serialise the key %s: %w", k.name, err)
	}

	cryptoJSON, err := gethkeystore.EncryptDataV3(serialized, []byte(password),
		commonkeystore.DefaultScryptParams.N, commonkeystore.DefaultScryptParams.P)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt the key %s: %w", k.name, err)
	}
	return json.Marshal(cryptoJSON)
}

func has(ctx context.Context, store commonkeystore.Keystore, name string) (bool, error) {
	keys, err := store.GetKeys(ctx, commonkeystore.GetKeysRequest{KeyNames: []string{name}})
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to look for the key %s: %w", name, err)
	}
	return len(keys.Keys) > 0, nil
}

// isNotFound is "the store does not have it", which is not a failure here: it is
// the whole question being asked.
//
// The message is checked as well as the error, because a missing key does not
// always arrive as a wrapped ErrKeyNotFound - a request for several keys reports
// what it could not find by joining, and a joined error is not one this can unwrap.
func isNotFound(err error) bool {
	return errors.Is(err, commonkeystore.ErrKeyNotFound) ||
		strings.Contains(err.Error(), commonkeystore.ErrKeyNotFound.Error())
}

func path(segments ...string) string {
	return commonkeystore.NewKeyPath(segments...).String()
}

// emptyAsNew lets a keystore that is not there yet be opened as an empty one,
// which is what the first boot after this was added finds.
type emptyAsNew struct {
	commonkeystore.Storage
}

func (s emptyAsNew) GetEncryptedKeystore(ctx context.Context) ([]byte, error) {
	data, err := s.Storage.GetEncryptedKeystore(ctx)
	if err != nil {
		// The store distinguishes "no rows" from a real failure by returning an error
		// either way, so an absent keystore reads as empty here rather than as broken.
		return nil, nil //nolint:nilerr // an absent keystore is an empty one
	}
	return data, nil
}
