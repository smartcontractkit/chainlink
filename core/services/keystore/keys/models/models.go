package models

import (
	"encoding/json"
	"math/big"
	"time"

	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
)

type EncryptedKeyRing struct {
	UpdatedAt     time.Time
	EncryptedKeys []byte
}

func (ekr EncryptedKeyRing) Decrypt(password string) (*KeyRing, error) {
	if len(ekr.EncryptedKeys) == 0 {
		return NewKeyRing(), nil
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

type KeyStates struct {
	// Key ID => chain ID => state
	KeyIDChainID map[string]map[string]*ethkey.State
	// Chain ID => Key ID => state
	ChainIDKeyID map[string]map[string]*ethkey.State
	All          []*ethkey.State
}

func NewKeyStates() *KeyStates {
	return &KeyStates{
		KeyIDChainID: make(map[string]map[string]*ethkey.State),
		ChainIDKeyID: make(map[string]map[string]*ethkey.State),
	}
}

// warning: not thread-safe! caller must sync
// adds or replaces a state
func (ks *KeyStates) Add(state *ethkey.State) {
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
func (ks *KeyStates) Get(addr common.Address, chainID *big.Int) *ethkey.State {
	chainStates, exists := ks.KeyIDChainID[addr.Hex()]
	if !exists {
		return nil
	}
	return chainStates[chainID.String()]
}

// warning: not thread-safe! caller must sync
func (ks *KeyStates) Disable(addr common.Address, chainID *big.Int, updatedAt time.Time) {
	state := ks.Get(addr, chainID)
	state.Disabled = true
	state.UpdatedAt = updatedAt
}

// warning: not thread-safe! caller must sync
func (ks *KeyStates) Enable(addr common.Address, chainID *big.Int, updatedAt time.Time) {
	state := ks.Get(addr, chainID)
	state.Disabled = false
	state.UpdatedAt = updatedAt
}

// warning: not thread-safe! caller must sync
func (ks *KeyStates) Delete(addr common.Address) {
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

type KeyRing struct {
	CSA          map[string]csakey.KeyV2
	Eth          map[string]ethkey.KeyV2
	OCR          map[string]ocrkey.KeyV2
	OCR2         map[string]ocr2key.KeyBundle
	P2P          map[string]p2pkey.KeyV2
	Cosmos       map[string]cosmoskey.Key
	Solana       map[string]solkey.Key
	StarkNet     map[string]starkkey.Key
	Sui          map[string]suikey.Key
	Aptos        map[string]aptoskey.Key
	Tron         map[string]tronkey.Key
	TON          map[string]tonkey.Key
	VRF          map[string]vrfkey.KeyV2
	Workflow     map[string]workflowkey.Key
	DKGRecipient map[string]dkgrecipientkey.Key
	LegacyKeys   LegacyKeyStorage
}

func NewKeyRing() *KeyRing {
	return &KeyRing{
		CSA:          make(map[string]csakey.KeyV2),
		Eth:          make(map[string]ethkey.KeyV2),
		OCR:          make(map[string]ocrkey.KeyV2),
		OCR2:         make(map[string]ocr2key.KeyBundle),
		P2P:          make(map[string]p2pkey.KeyV2),
		Cosmos:       make(map[string]cosmoskey.Key),
		Solana:       make(map[string]solkey.Key),
		StarkNet:     make(map[string]starkkey.Key),
		Sui:          make(map[string]suikey.Key),
		Aptos:        make(map[string]aptoskey.Key),
		Tron:         make(map[string]tronkey.Key),
		TON:          make(map[string]tonkey.Key),
		VRF:          make(map[string]vrfkey.KeyV2),
		Workflow:     make(map[string]workflowkey.Key),
		DKGRecipient: make(map[string]dkgrecipientkey.Key),
	}
}

func (kr *KeyRing) Encrypt(password string, scryptParams keystore.ScryptParams) (ekr EncryptedKeyRing, err error) {
	marshalledRawKeyRingJson, err := json.Marshal(kr.raw())
	if err != nil {
		return ekr, err
	}

	marshalledRawKeyRingJson, err = kr.LegacyKeys.UnloadUnsupported(marshalledRawKeyRingJson)
	if err != nil {
		return EncryptedKeyRing{}, err
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
	return EncryptedKeyRing{
		EncryptedKeys: encryptedKeys,
	}, nil
}

func (kr *KeyRing) raw() (rawKeys rawKeyRing) {
	for _, csaKey := range kr.CSA {
		rawKeys.CSA = append(rawKeys.CSA, internal.RawBytes(csaKey))
	}
	for _, ethKey := range kr.Eth {
		rawKeys.Eth = append(rawKeys.Eth, internal.RawBytes(ethKey))
	}
	for _, ocrKey := range kr.OCR {
		rawKeys.OCR = append(rawKeys.OCR, internal.RawBytes(ocrKey))
	}
	for _, ocr2key := range kr.OCR2 {
		rawKeys.OCR2 = append(rawKeys.OCR2, internal.RawBytes(ocr2key))
	}
	for _, p2pKey := range kr.P2P {
		rawKeys.P2P = append(rawKeys.P2P, internal.RawBytes(p2pKey))
	}
	for _, cosmoskey := range kr.Cosmos {
		rawKeys.Cosmos = append(rawKeys.Cosmos, internal.RawBytes(cosmoskey))
	}
	for _, solkey := range kr.Solana {
		rawKeys.Solana = append(rawKeys.Solana, internal.RawBytes(solkey))
	}
	for _, starkkey := range kr.StarkNet {
		rawKeys.StarkNet = append(rawKeys.StarkNet, internal.RawBytes(starkkey))
	}
	for _, aptoskey := range kr.Aptos {
		rawKeys.Aptos = append(rawKeys.Aptos, internal.RawBytes(aptoskey))
	}
	for _, tronkey := range kr.Tron {
		rawKeys.Tron = append(rawKeys.Tron, internal.RawBytes(tronkey))
	}
	for _, tonkey := range kr.TON {
		rawKeys.TON = append(rawKeys.TON, internal.RawBytes(tonkey))
	}
	for _, suikey := range kr.Sui {
		rawKeys.Sui = append(rawKeys.Sui, internal.RawBytes(suikey))
	}
	for _, vrfKey := range kr.VRF {
		rawKeys.VRF = append(rawKeys.VRF, internal.RawBytes(vrfKey))
	}
	for _, workflowKey := range kr.Workflow {
		rawKeys.Workflow = append(rawKeys.Workflow, internal.RawBytes(workflowKey))
	}
	for _, dkgRecipientKey := range kr.DKGRecipient {
		rawKeys.DKGRecipient = append(rawKeys.DKGRecipient, internal.RawBytes(dkgRecipientKey))
	}
	return rawKeys
}

// rawKeyRing is an intermediate struct for encrypting / decrypting KeyRing
// it holds only the essential key information to avoid adding unnecessary data
// (like public keys) to the database
type rawKeyRing struct {
	Eth          [][]byte
	CSA          [][]byte
	OCR          [][]byte
	OCR2         [][]byte
	P2P          [][]byte
	Cosmos       [][]byte
	Solana       [][]byte
	StarkNet     [][]byte
	Sui          [][]byte
	Aptos        [][]byte
	Tron         [][]byte
	TON          [][]byte
	VRF          [][]byte
	Workflow     [][]byte
	DKGRecipient [][]byte
	LegacyKeys   LegacyKeyStorage `json:"-"`
}

func (rawKeys rawKeyRing) keys() (*KeyRing, error) {
	keyRing := NewKeyRing()
	for _, rawCSAKey := range rawKeys.CSA {
		csaKey := csakey.KeyFor(internal.NewRaw(rawCSAKey))
		keyRing.CSA[csaKey.ID()] = csaKey
	}
	for _, rawETHKey := range rawKeys.Eth {
		ethKey := ethkey.KeyFor(internal.NewRaw(rawETHKey))
		keyRing.Eth[ethKey.ID()] = ethKey
	}
	for _, rawOCRKey := range rawKeys.OCR {
		ocrKey := ocrkey.KeyFor(internal.NewRaw(rawOCRKey))
		keyRing.OCR[ocrKey.ID()] = ocrKey
	}
	for _, rawOCR2Key := range rawKeys.OCR2 {
		if ocr2Key := ocr2key.KeyFor(internal.NewRaw(rawOCR2Key)); ocr2Key != nil {
			keyRing.OCR2[ocr2Key.ID()] = ocr2Key
		}
	}
	for _, rawP2PKey := range rawKeys.P2P {
		p2pKey := p2pkey.KeyFor(internal.NewRaw(rawP2PKey))
		keyRing.P2P[p2pKey.ID()] = p2pKey
	}
	for _, rawCosmosKey := range rawKeys.Cosmos {
		cosmosKey := cosmoskey.KeyFor(internal.NewRaw(rawCosmosKey))
		keyRing.Cosmos[cosmosKey.ID()] = cosmosKey
	}
	for _, rawSolKey := range rawKeys.Solana {
		solKey := solkey.KeyFor(internal.NewRaw(rawSolKey))
		keyRing.Solana[solKey.ID()] = solKey
	}
	for _, rawStarkNetKey := range rawKeys.StarkNet {
		starkKey := starkkey.KeyFor(internal.NewRaw(rawStarkNetKey))
		keyRing.StarkNet[starkKey.ID()] = starkKey
	}
	for _, rawAptosKey := range rawKeys.Aptos {
		aptosKey := aptoskey.KeyFor(internal.NewRaw(rawAptosKey))
		keyRing.Aptos[aptosKey.ID()] = aptosKey
	}
	for _, rawTronKey := range rawKeys.Tron {
		tronKey := tronkey.KeyFor(internal.NewRaw(rawTronKey))
		keyRing.Tron[tronKey.ID()] = tronKey
	}
	for _, rawTONKey := range rawKeys.TON {
		tonKey := tonkey.KeyFor(internal.NewRaw(rawTONKey))
		keyRing.TON[tonKey.ID()] = tonKey
	}
	for _, rawSuiKey := range rawKeys.Sui {
		suiKey := suikey.KeyFor(internal.NewRaw(rawSuiKey))
		keyRing.Sui[suiKey.ID()] = suiKey
	}
	for _, rawVRFKey := range rawKeys.VRF {
		vrfKey := vrfkey.KeyFor(internal.NewRaw(rawVRFKey))
		keyRing.VRF[vrfKey.ID()] = vrfKey
	}
	for _, rawWorkflowKey := range rawKeys.Workflow {
		workflowKey := workflowkey.KeyFor(internal.NewRaw(rawWorkflowKey))
		keyRing.Workflow[workflowKey.ID()] = workflowKey
	}
	for _, rawDKGRecipientKey := range rawKeys.DKGRecipient {
		dkgRecipientKey := dkgrecipientkey.KeyFor(internal.NewRaw(rawDKGRecipientKey))
		keyRing.DKGRecipient[dkgRecipientKey.ID()] = dkgRecipientKey
	}

	keyRing.LegacyKeys = rawKeys.LegacyKeys
	return keyRing, nil
}

// adulteration prevents the password from getting used in the wrong place
func adulteratedPassword(password string) string {
	return "master-password-" + password
}
