package account

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

type Keystore interface {
	Sign(ctx context.Context, id string, msgHash *big.Int) (x *big.Int, y *big.Int, err error)
}

// MemKeystore implements the Keystore interface and is intended for example and test code.
type MemKeystore struct {
	mu   sync.Mutex
	keys map[string]*big.Int
}

// NewMemKeystore initializes and returns a new instance of MemKeystore.
//
// Parameters:
//  none
// Returns:
// - *MemKeystore: a pointer to MemKeystore.
func NewMemKeystore() *MemKeystore {
	return &MemKeystore{
		keys: make(map[string]*big.Int),
	}
}

// SetNewMemKeystore returns a new instance of MemKeystore and sets the given public key and private key in it.
//
// Parameters:
// - pub: a string representing the public key
// - priv: a pointer to a big.Int representing the private key
// Returns:
// - *MemKeystore: a pointer to the newly created MemKeystore instance
func SetNewMemKeystore(pub string, priv *big.Int) *MemKeystore {
	ks := NewMemKeystore()
	ks.Put(pub, priv)
	return ks
}

// Put stores the given key in the keystore for the specified sender address.
//
// Parameters:
// - senderAddress: the address of the sender
// - k: the key to be stored
func (ks *MemKeystore) Put(senderAddress string, k *big.Int) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys[senderAddress] = k
}

var ErrSenderNoExist = errors.New("sender does not exist")

// Get retrieves the value associated with the senderAddress from the MemKeystore.
//
// Parameter:
// - senderAddress: The address of the sender
// Returns:
// - *big.Int: The value associated with the senderAddress
// - error: An error if the senderAddress does not exist in the keystore
func (ks *MemKeystore) Get(senderAddress string) (*big.Int, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	k, exists := ks.keys[senderAddress]
	if !exists {
		return nil, fmt.Errorf("error getting key for sender %s: %w", senderAddress, ErrSenderNoExist)
	}
	return k, nil
}

// Sign signs a message hash using the given key in the MemKeystore.
//
// Parameters:
// - ctx: the context of the operation.
// - id: is the identifier of the key.
// - msgHash: is the message hash to be signed.
// Returns:
// - *big.Int: the R component of the signature as *big.Int
// - *big.Int: the S component of the signature as *big.Int
// - error: an error if any
func (ks *MemKeystore) Sign(ctx context.Context, id string, msgHash *big.Int) (*big.Int, *big.Int, error) {

	k, err := ks.Get(id)
	if err != nil {
		return nil, nil, err
	}

	return sign(ctx, msgHash, k)
}

// sign signs the given message hash with the provided key using the Curve.
// illustrates one way to handle context cancellation
//
// Parameters:
// - ctx: the context.Context object for cancellation and timeouts
// - msgHash: the message hash to be signed as a *big.Int
// - key: the private key as a *big.Int
// Returns:
// - x: the X coordinate of the signature point as a *big.Int
// - y: the Y coordinate of the signature point as a *big.Int
// - err: an error object if any error occurred during the signing process
func sign(ctx context.Context, msgHash *big.Int, key *big.Int) (x *big.Int, y *big.Int, err error) {

	select {
	case <-ctx.Done():
		x = nil
		y = nil
		err = ctx.Err()

	default:
		x, y, err = curve.Curve.Sign(msgHash, key)
	}
	return x, y, err
}

// GetRandomKeys gets a random set of pub-priv keys.
// Note: This should be used for testing purposes only, do NOT send real funds to these addresses.
// Parameters:
//  none
// Returns:
// - *MemKeystore: a pointer to a MemKeystore instance
// - *felt.Felt: a pointer to a public key as a felt.Felt
// - *felt.Felt: a pointer to a private key as a felt.Felt
func GetRandomKeys() (*MemKeystore, *felt.Felt, *felt.Felt) {
	// Get random keys
	privateKey, err := curve.Curve.GetRandomPrivateKey()
	if err != nil {
		fmt.Println("can't get random private key:", err)
		os.Exit(1)
	}
	pubX, _, err := curve.Curve.PrivateToPoint(privateKey)
	if err != nil {
		fmt.Println("can't generate public key:", err)
		os.Exit(1)
	}
	privFelt := utils.BigIntToFelt(privateKey)
	pubFelt := utils.BigIntToFelt(pubX)

	// set up keystore
	ks := NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privFelt.String(), 0)
	if !ok {
		panic("Error setting up account key store")
	}
	ks.Put(pubFelt.String(), fakePrivKeyBI)

	return ks, pubFelt, privFelt
}
