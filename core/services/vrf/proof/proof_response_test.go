package proof_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	proof2 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
)

func TestMarshaledProof(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	key := cltest.DefaultVRFKey
	require.NoError(t, keyStore.VRF().Add(ctx, key))
	blockHash := common.Hash{}
	blockNum := 0
	preSeed := big.NewInt(1)
	s := proof2.TestXXXSeedData(t, preSeed, blockHash, blockNum)
	proofResponse, err := proof2.GenerateProofResponse(keyStore.VRF(), key.ID(), s)
	require.NoError(t, err)
	goProof, err := proof2.UnmarshalProofResponse(proofResponse)
	require.NoError(t, err)
	actualProof, err := goProof.CryptoProof(s)
	require.NoError(t, err)
	_, err = proof2.MarshalForSolidityVerifier(&actualProof)
	require.NoError(t, err)
}
