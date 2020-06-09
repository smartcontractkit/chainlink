package vrfkey

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"testing"

	tvrf "github.com/smartcontractkit/chainlink/core/internal/cltest/vrf"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_verifier_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sk = 0xdeadbeefdeadbee
var k = mustNewPrivateKey(big.NewInt(int64(sk)))
var pkr = regexp.MustCompile(fmt.Sprintf(
	`PrivateKey\{k: <redacted>, PublicKey: 0x[[:xdigit:]]{%d}\}`,
	2*CompressedPublicKeyLength))

func TestPrintingDoesNotLeakKey(t *testing.T) {
	v := fmt.Sprintf("%v", k)
	assert.Equal(t, v+"\n", fmt.Sprintln(k))
	assert.Regexp(t, pkr, v)
	assert.NotContains(t, v, fmt.Sprintf("%x", sk))
	// Other verbs just give the corresponding encoding of .String()
	assert.Equal(t, fmt.Sprintf("%x", k), hex.EncodeToString([]byte(v)))
}

func TestMarshaledProof(t *testing.T) {
	blockHash := common.Hash{}
	blockNum := 0
	preSeed := big.NewInt(1)
	s := tvrf.SeedData(t, preSeed, blockHash, blockNum)
	proofResponse, err := k.MarshaledProof(s)
	require.NoError(t, err)
	goProof, err := vrf.UnmarshalProofResponse(proofResponse)
	require.NoError(t, err)
	actualProof, err := goProof.ActualProof(s)
	require.NoError(t, err)
	proof, err := actualProof.MarshalForSolidityVerifier()
	require.NoError(t, err)
	// NB: For changes to the VRF solidity code to be reflected here, "go generate"
	// must be run in core/services/vrf.
	ethereumKey, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(ethereumKey)
	genesisData := core.GenesisAlloc{auth.From: {Balance: big.NewInt(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	_, _, verifier, err := solidity_vrf_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
	backend.Commit()
	_, err = verifier.RandomValueFromVRFProof(nil, proof[:])
	require.NoError(t, err)
}

func mustNewPrivateKey(rawKey *big.Int) *PrivateKey {
	k, err := newPrivateKey(rawKey)
	if err != nil {
		panic(err)
	}
	return k
}
