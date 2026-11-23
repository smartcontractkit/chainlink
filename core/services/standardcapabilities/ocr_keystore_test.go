package standardcapabilities

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	capkeystore "github.com/smartcontractkit/capabilities/libs/capability/keystore"
	capocr "github.com/smartcontractkit/capabilities/libs/capability/keystore/ocr"
	ocrremote "github.com/smartcontractkit/capabilities/libs/capability/keystore/ocr/remote"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestOCRKeystore_ListAndGet(t *testing.T) {
	t.Parallel()
	evm, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	sol, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)

	ks := keystoremocks.NewOCR2(t)
	ks.On("GetAll").Return([]ocr2key.KeyBundle{evm, sol}, nil)
	s := ocrKeystore{ks: ks}

	all, err := s.List(t.Context())
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Less(t, all[0].ID(), all[1].ID(), "bundles are sorted by ID")

	solana, err := s.List(t.Context(), "solana")
	require.NoError(t, err)
	require.Len(t, solana, 1)
	assert.Equal(t, sol.ID(), solana[0].ID())
	assert.Equal(t, "solana", solana[0].ChainType())

	kb, err := s.Get(t.Context(), evm.ID())
	require.NoError(t, err)
	assert.Equal(t, "evm", kb.ChainType())
	assert.Equal(t, evm.PublicKey(), kb.PublicKey())

	_, err = s.Get(t.Context(), "missing")
	require.ErrorIs(t, err, capocr.ErrKeyBundleNotFound)
}

func TestOCRKeystore_Locked(t *testing.T) {
	t.Parallel()
	ks := keystoremocks.NewOCR2(t)
	ks.On("GetAll").Return([]ocr2key.KeyBundle{}, keystore.ErrLocked)
	s := ocrKeystore{ks: ks}

	_, err := s.List(t.Context())
	require.ErrorIs(t, err, keystore.ErrLocked)

	_, err = s.Get(t.Context(), "any")
	require.ErrorIs(t, err, keystore.ErrLocked)
	require.NotErrorIs(t, err, capocr.ErrKeyBundleNotFound)
}

func TestStandardCapabilitiesV2_ServesOCRKeystore(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	evm, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	ks := keystoremocks.NewOCR2(t)
	ks.On("GetAll").Return([]ocr2key.KeyBundle{evm}, nil)

	rc := ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{ConfigDigest: ocrtypes.ConfigDigest{1}, Epoch: 2, Round: 3}}
	// Larger than gRPC's 4 MiB default, as libocr allows.
	report := ocrtypes.Report(bytes.Repeat([]byte{7}, ocr3types.MaxMaxReportLength))

	var env map[string]string
	var kb capocr.KeyBundleClient
	var onchainSig, offchainSig []byte
	var secret [32]byte
	var callErr error
	registrar := &capturingRegistrar{onRegister: func(cfg plugins.CmdConfig) {
		env = envMap(cfg.Env)
		client := ocrremote.NewClient(lggr, dial(t, env["CL_KEYSTORE_URL"]))
		if kb, callErr = client.GetKeyBundle(t.Context(), evm.ID()); callErr != nil {
			return
		}
		if onchainSig, callErr = kb.Onchain.Sign(rc, report); callErr != nil {
			return
		}
		if offchainSig, callErr = kb.Offchain.OffchainSign([]byte("msg")); callErr != nil {
			return
		}
		secret, callErr = kb.Offchain.ConfigDiffieHellman([32]byte{4})
	}}

	s, err := NewV2(V2Params{
		Log:              lggr,
		Command:          "cronv2",
		App:              "cron",
		PluginRegistrar:  registrar,
		CapDonIDResolver: &staticCapDonIDResolver{},
		Keystores:        capkeystore.Keystores{OCR: ocrKeystore{ks: ks}},
	})
	require.NoError(t, err)
	require.NoError(t, s.Start(t.Context()))
	s.launchWG.Wait()
	require.NoError(t, s.Close())
	require.NoError(t, callErr)

	assert.Equal(t, "evm", kb.ChainType)
	assert.Equal(t, evm.PublicKey(), kb.Onchain.PublicKey())
	assert.Equal(t, evm.MaxSignatureLength(), kb.Onchain.MaxSignatureLength())
	assert.Equal(t, evm.OffchainPublicKey(), kb.Offchain.OffchainPublicKey())
	assert.Equal(t, evm.ConfigEncryptionPublicKey(), kb.Offchain.ConfigEncryptionPublicKey())

	assert.True(t, evm.Verify(evm.PublicKey(), rc, report, onchainSig))
	wantOffchainSig, err := evm.OffchainSign([]byte("msg"))
	require.NoError(t, err)
	assert.Equal(t, wantOffchainSig, offchainSig)
	wantSecret, err := evm.ConfigDiffieHellman([32]byte{4})
	require.NoError(t, err)
	assert.Equal(t, wantSecret, secret)
}
