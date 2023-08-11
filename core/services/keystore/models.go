package keystore

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type encryptedKeyRing struct {
	UpdatedAt     time.Time
	EncryptedKeys []byte
}

func (ekr encryptedKeyRing) Decrypt(password string) (*keyRing, error) {
	if len(ekr.EncryptedKeys) == 0 {
		return newKeyRing(), nil
	}
	var cryptoJSON gethkeystore.CryptoJSON
	err := json.Unmarshal(ekr.EncryptedKeys, &cryptoJSON)
	if err != nil {
		return nil, err
	}
	marshalledRawKeyRingJson, err := gethkeystore.DecryptDataV3(cryptoJSON, adulteratedPassword(password))
	if err != nil {
		return nil, err
	}
	var rawKeys rawKeyRing
	err = json.Unmarshal(marshalledRawKeyRingJson, &rawKeys)
	if err != nil {
		return nil, err
	}
	ring, err := rawKeys.keys()
	if err != nil {
		return nil, err
	}

	err = rawKeys.LegacyKeys.StoreUnsupported(marshalledRawKeyRingJson, ring)
	if err != nil {
		return nil, err
	}
	ring.LegacyKeys = rawKeys.LegacyKeys

	return ring, nil
}

type keyStates struct {
	// Key ID => chain ID => state
	KeyIDChainID map[string]map[string]*ethkey.State
	// Chain ID => Key ID => state
	ChainIDKeyID map[string]map[string]*ethkey.State
	All          []*ethkey.State
}

func newKeyStates() *keyStates {
	return &keyStates{
		KeyIDChainID: make(map[string]map[string]*ethkey.State),
		ChainIDKeyID: make(map[string]map[string]*ethkey.State),
	}
}

// warning: not thread-safe! caller must sync
// adds or replaces a state
func (ks *keyStates) add(state *ethkey.State) {
	cid := state.EVMChainID.String()
	kid := state.KeyID()

	keyStates, exists := ks.KeyIDChainID[kid]
	if !exists {
		keyStates = make(map[string]*ethkey.State)
		ks.KeyIDChainID[kid] = keyStates
	}
	keyStates[cid] = state

	chainStates, exists := ks.ChainIDKeyID[cid]
	if !exists {
		chainStates = make(map[string]*ethkey.State)
		ks.ChainIDKeyID[cid] = chainStates
	}
	chainStates[kid] = state

	exists = false
	for i, existingState := range ks.All {
		if existingState.ID == state.ID {
			ks.All[i] = state
			exists = true
			break
		}
	}
	if !exists {
		ks.All = append(ks.All, state)
	}
}

// warning: not thread-safe! caller must sync
func (ks *keyStates) get(addr common.Address, chainID *big.Int) *ethkey.State {
	chainStates, exists := ks.KeyIDChainID[addr.Hex()]
	if !exists {
		return nil
	}
	return chainStates[chainID.String()]
}

// warning: not thread-safe! caller must sync
func (ks *keyStates) disable(addr common.Address, chainID *big.Int, updatedAt time.Time) {
	state := ks.get(addr, chainID)
	state.Disabled = true
	state.UpdatedAt = updatedAt
}

// warning: not thread-safe! caller must sync
func (ks *keyStates) enable(addr common.Address, chainID *big.Int, updatedAt time.Time) {
	state := ks.get(addr, chainID)
	state.Disabled = false
	state.UpdatedAt = updatedAt
}

// warning: not thread-safe! caller must sync
func (ks *keyStates) delete(addr common.Address) {
	var chainIDs []*big.Int
	for i := len(ks.All) - 1; i >= 0; i-- {
		if ks.All[i].Address.Address() == addr {
			chainIDs = append(chainIDs, ks.All[i].EVMChainID.ToInt())
			ks.All = append(ks.All[:i], ks.All[i+1:]...)
		}
	}
	for _, cid := range chainIDs {
		delete(ks.KeyIDChainID[addr.Hex()], cid.String())
		delete(ks.ChainIDKeyID[cid.String()], addr.Hex())
	}
}

type keyRing struct {
	CSA        map[string]csakey.KeyV2
	Eth        map[string]ethkey.KeyV2
	OCR        map[string]ocrkey.KeyV2
	OCR2       map[string]ocr2key.KeyBundle
	P2P        map[string]p2pkey.KeyV2
	Cosmos     map[string]cosmoskey.Key
	Solana     map[string]solkey.Key
	StarkNet   map[string]starkkey.Key
	VRF        map[string]vrfkey.KeyV2
	DKGSign    map[string]dkgsignkey.Key
	DKGEncrypt map[string]dkgencryptkey.Key
	LegacyKeys LegacyKeyStorage
}

func newKeyRing() *keyRing {
	return &keyRing{
		CSA:        make(map[string]csakey.KeyV2),
		Eth:        make(map[string]ethkey.KeyV2),
		OCR:        make(map[string]ocrkey.KeyV2),
		OCR2:       make(map[string]ocr2key.KeyBundle),
		P2P:        make(map[string]p2pkey.KeyV2),
		Cosmos:     make(map[string]cosmoskey.Key),
		Solana:     make(map[string]solkey.Key),
		StarkNet:   make(map[string]starkkey.Key),
		VRF:        make(map[string]vrfkey.KeyV2),
		DKGSign:    make(map[string]dkgsignkey.Key),
		DKGEncrypt: make(map[string]dkgencryptkey.Key),
	}
}

func (kr *keyRing) Encrypt(password string, scryptParams utils.ScryptParams) (ekr encryptedKeyRing, err error) {
	marshalledRawKeyRingJson, err := json.Marshal(kr.raw())
	if err != nil {
		return ekr, err
	}

	marshalledRawKeyRingJson, err = kr.LegacyKeys.UnloadUnsupported(marshalledRawKeyRingJson)
	if err != nil {
		return encryptedKeyRing{}, err
	}

	cryptoJSON, err := gethkeystore.EncryptDataV3(
		marshalledRawKeyRingJson,
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return ekr, errors.Wrapf(err, "could not encrypt key ring")
	}
	encryptedKeys, err := json.Marshal(&cryptoJSON)
	if err != nil {
		return ekr, errors.Wrapf(err, "could not encode cryptoJSON")
	}
	return encryptedKeyRing{
		EncryptedKeys: encryptedKeys,
	}, nil
}

func (kr *keyRing) raw() (rawKeys rawKeyRing) {
	for _, csaKey := range kr.CSA {
		rawKeys.CSA = append(rawKeys.CSA, csaKey.Raw())
	}
	for _, ethKey := range kr.Eth {
		rawKeys.Eth = append(rawKeys.Eth, ethKey.Raw())
	}
	for _, ocrKey := range kr.OCR {
		rawKeys.OCR = append(rawKeys.OCR, ocrKey.Raw())
	}
	for _, ocr2key := range kr.OCR2 {
		rawKeys.OCR2 = append(rawKeys.OCR2, ocr2key.Raw())
	}
	for _, p2pKey := range kr.P2P {
		rawKeys.P2P = append(rawKeys.P2P, p2pKey.Raw())
	}
	for _, cosmoskey := range kr.Cosmos {
		rawKeys.Cosmos = append(rawKeys.Cosmos, cosmoskey.Raw())
	}
	for _, solkey := range kr.Solana {
		rawKeys.Solana = append(rawKeys.Solana, solkey.Raw())
	}
	for _, starkkey := range kr.StarkNet {
		rawKeys.StarkNet = append(rawKeys.StarkNet, starkkey.Raw())
	}
	for _, vrfKey := range kr.VRF {
		rawKeys.VRF = append(rawKeys.VRF, vrfKey.Raw())
	}
	for _, dkgSignKey := range kr.DKGSign {
		rawKeys.DKGSign = append(rawKeys.DKGSign, dkgSignKey.Raw())
	}
	for _, dkgEncryptKey := range kr.DKGEncrypt {
		rawKeys.DKGEncrypt = append(rawKeys.DKGEncrypt, dkgEncryptKey.Raw())
	}
	return rawKeys
}

func (kr *keyRing) logPubKeys(lggr logger.Logger) {
	lggr = lggr.Named("KeyRing")
	var csaIDs []string
	for _, CSAKey := range kr.CSA {
		csaIDs = append(csaIDs, CSAKey.ID())
	}
	var ethIDs []string
	for _, ETHKey := range kr.Eth {
		ethIDs = append(ethIDs, ETHKey.ID())
	}
	var ocrIDs []string
	for _, OCRKey := range kr.OCR {
		ocrIDs = append(ocrIDs, OCRKey.ID())
	}
	var ocr2IDs []string
	for _, OCR2Key := range kr.OCR2 {
		ocr2IDs = append(ocr2IDs, OCR2Key.ID())
	}
	var p2pIDs []string
	for _, P2PKey := range kr.P2P {
		p2pIDs = append(p2pIDs, P2PKey.ID())
	}
	var cosmosIDs []string
	for _, cosmosKey := range kr.Cosmos {
		cosmosIDs = append(cosmosIDs, cosmosKey.ID())
	}
	var solanaIDs []string
	for _, solanaKey := range kr.Solana {
		solanaIDs = append(solanaIDs, solanaKey.ID())
	}
	var starknetIDs []string
	for _, starkkey := range kr.StarkNet {
		starknetIDs = append(starknetIDs, starkkey.ID())
	}
	var vrfIDs []string
	for _, VRFKey := range kr.VRF {
		vrfIDs = append(vrfIDs, VRFKey.ID())
	}
	var dkgSignIDs []string
	for _, dkgSignKey := range kr.DKGSign {
		dkgSignIDs = append(dkgSignIDs, dkgSignKey.ID())
	}
	var dkgEncryptIDs []string
	for _, dkgEncryptKey := range kr.DKGEncrypt {
		dkgEncryptIDs = append(dkgEncryptIDs, dkgEncryptKey.ID())
	}
	if len(csaIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d CSA keys", len(csaIDs)), "keys", csaIDs)
	}
	if len(ethIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d ETH keys", len(ethIDs)), "keys", ethIDs)
	}
	if len(ocrIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d OCR keys", len(ocrIDs)), "keys", ocrIDs)
	}
	if len(ocr2IDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d OCR2 keys", len(ocr2IDs)), "keys", ocr2IDs)
	}
	if len(p2pIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d P2P keys", len(p2pIDs)), "keys", p2pIDs)
	}
	if len(cosmosIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d Cosmos keys", len(cosmosIDs)), "keys", cosmosIDs)
	}
	if len(solanaIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d Solana keys", len(solanaIDs)), "keys", solanaIDs)
	}
	if len(starknetIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d StarkNet keys", len(starknetIDs)), "keys", starknetIDs)
	}
	if len(vrfIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d VRF keys", len(vrfIDs)), "keys", vrfIDs)
	}
	if len(dkgSignIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d DKGSign keys", len(dkgSignIDs)), "keys", dkgSignIDs)
	}
	if len(dkgEncryptIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d DKGEncrypt keys", len(dkgEncryptIDs)), "keys", dkgEncryptIDs)
	}
	if len(kr.LegacyKeys.legacyRawKeys) > 0 {
		lggr.Infow(fmt.Sprintf("%d keys stored in legacy system", kr.LegacyKeys.legacyRawKeys.len()))
	}
}

// rawKeyRing is an intermediate struct for encrypting / decrypting keyRing
// it holds only the essential key information to avoid adding unnecessary data
// (like public keys) to the database
type rawKeyRing struct {
	Eth        []ethkey.Raw
	CSA        []csakey.Raw
	OCR        []ocrkey.Raw
	OCR2       []ocr2key.Raw
	P2P        []p2pkey.Raw
	Cosmos     []cosmoskey.Raw
	Solana     []solkey.Raw
	StarkNet   []starkkey.Raw
	VRF        []vrfkey.Raw
	DKGSign    []dkgsignkey.Raw
	DKGEncrypt []dkgencryptkey.Raw
	LegacyKeys LegacyKeyStorage `json:"-"`
}

func (rawKeys rawKeyRing) keys() (*keyRing, error) {
	keyRing := newKeyRing()
	for _, rawCSAKey := range rawKeys.CSA {
		csaKey := rawCSAKey.Key()
		keyRing.CSA[csaKey.ID()] = csaKey
	}
	for _, rawETHKey := range rawKeys.Eth {
		ethKey := rawETHKey.Key()
		keyRing.Eth[ethKey.ID()] = ethKey
	}
	for _, rawOCRKey := range rawKeys.OCR {
		ocrKey := rawOCRKey.Key()
		keyRing.OCR[ocrKey.ID()] = ocrKey
	}
	for _, rawOCR2Key := range rawKeys.OCR2 {
		if ocr2Key := rawOCR2Key.Key(); ocr2Key != nil {
			keyRing.OCR2[ocr2Key.ID()] = ocr2Key
		}
	}
	for _, rawP2PKey := range rawKeys.P2P {
		p2pKey := rawP2PKey.Key()
		keyRing.P2P[p2pKey.ID()] = p2pKey
	}
	for _, rawCosmosKey := range rawKeys.Cosmos {
		cosmosKey := rawCosmosKey.Key()
		keyRing.Cosmos[cosmosKey.ID()] = cosmosKey
	}
	for _, rawSolKey := range rawKeys.Solana {
		solKey := rawSolKey.Key()
		keyRing.Solana[solKey.ID()] = solKey
	}
	for _, rawStarkNetKey := range rawKeys.StarkNet {
		starkKey := rawStarkNetKey.Key()
		keyRing.StarkNet[starkKey.ID()] = starkKey
	}
	for _, rawVRFKey := range rawKeys.VRF {
		vrfKey := rawVRFKey.Key()
		keyRing.VRF[vrfKey.ID()] = vrfKey
	}
	for _, rawDKGSignKey := range rawKeys.DKGSign {
		dkgSignKey := rawDKGSignKey.Key()
		keyRing.DKGSign[dkgSignKey.ID()] = dkgSignKey
	}
	for _, rawDKGEncryptKey := range rawKeys.DKGEncrypt {
		dkgEncryptKey := rawDKGEncryptKey.Key()
		keyRing.DKGEncrypt[dkgEncryptKey.ID()] = dkgEncryptKey
	}

	keyRing.LegacyKeys = rawKeys.LegacyKeys
	return keyRing, nil
}

// adulteration prevents the password from getting used in the wrong place
func adulteratedPassword(password string) string {
	return "master-password-" + password
}
