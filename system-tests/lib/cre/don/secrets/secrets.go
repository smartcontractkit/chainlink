package secrets

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type nodeSecret struct {
	EthKeys         nodeEthKeyWrapper   `toml:"EVM"`
	SolKeys         nodeSolKeyWrapper   `toml:"Solana"`
	P2PKey          nodeP2PKey          `toml:"P2PKey"`
	DKGRecipientKey nodeDKGRecipientKey `toml:"DKGRecipientKey"`

	// Add more fields as needed to reflect 'Secrets' struct from /core/config/toml/types.go
	// We can't use the original struct, because it's using custom types that serialize secrets to 'xxxxx'
}

type nodeEthKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
	ChainID  uint64 `toml:"ID"`
}

type nodeSolKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
	ChainID  string `toml:"ID"`
}

type nodeP2PKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
}

type nodeDKGRecipientKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
}

type nodeEthKeyWrapper struct {
	EthKeys []nodeEthKey `toml:"Keys"`
}

type nodeSolKeyWrapper struct {
	SolKeys []nodeSolKey `toml:"Keys"`
}

type ChainFamily = string

type NodeKeys struct {
	CSAKey        *crypto.CSAKey
	EVM           map[uint64]*crypto.EVMKey
	Solana        map[string]*crypto.SolKey
	P2PKey        *crypto.P2PKey
	DKGKey        *crypto.DKGRecipientKey
	OCR2BundleIDs map[ChainFamily]string
}

func (n NodeKeys) PeerID() string {
	if n.P2PKey == nil {
		return ""
	}
	return n.P2PKey.PeerID.String()
}

func (n *NodeKeys) ToNodeSecretsTOML() (string, error) {
	ns := nodeSecret{}

	if n.P2PKey != nil {
		ns.P2PKey = nodeP2PKey{
			JSON:     string(n.P2PKey.EncryptedJSON),
			Password: n.P2PKey.Password,
		}
	}

	if n.DKGKey != nil {
		ns.DKGRecipientKey = nodeDKGRecipientKey{
			JSON:     string(n.DKGKey.EncryptedJSON),
			Password: n.DKGKey.Password,
		}
	}

	if n.EVM != nil {
		ns.EthKeys = nodeEthKeyWrapper{}
		for chainID, evmKeys := range n.EVM {
			ns.EthKeys.EthKeys = append(ns.EthKeys.EthKeys, nodeEthKey{
				JSON:     string(evmKeys.EncryptedJSON),
				Password: evmKeys.Password,
				ChainID:  chainID,
			})
		}
	}

	if n.Solana != nil {
		ns.SolKeys = nodeSolKeyWrapper{}
		for chainID, solKeys := range n.Solana {
			ns.SolKeys.SolKeys = append(ns.SolKeys.SolKeys, nodeSolKey{
				JSON:     string(solKeys.EncryptedJSON),
				Password: solKeys.Password,
				ChainID:  chainID,
			})
		}
	}

	nodeSecretString, err := toml.Marshal(ns)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal node secrets")
	}
	return string(nodeSecretString), nil
}

// secrets struct mirrors `Secrets` struct in "github.com/smartcontractkit/chainlink/v2/core/config/toml"
// we use a copy to avoid depending on the core config package, we consider it safe, because that struct changes very rarely
type secrets struct {
	EVM             ethKeys         `toml:",omitempty"` // choose EVM as the TOML field name to align with relayer config convention
	P2PKey          p2PKey          `toml:",omitempty"`
	Solana          solKeys         `toml:",omitempty"`
	DKGRecipientKey dkgRecipientKey `toml:",omitempty"`
}

type p2PKey struct {
	JSON     *string
	Password *string
}

type dkgRecipientKey struct {
	JSON     *string
	Password *string
}

type ethKeys struct {
	Keys []*ethKey
}

type ethKey struct {
	JSON     *string
	ID       *uint64
	Password *string
}

type solKeys struct {
	Keys []*solKey
}

type solKey struct {
	JSON     *string
	ID       *string
	Password *string
}

// struct required for reading "address" from this bit of encrypted JSON:
// JSON = '{"address":"e753ac0b6e175ce3a939c55433a0109c5a6f8777"}'
type evmJSON struct {
	Address string `json:"address"`
}

func publicEVMAddressFromEncryptedJSON(jsonString string) (string, error) {
	var eJSON evmJSON

	err := json.Unmarshal([]byte(jsonString), &eJSON)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal evm json")
	}

	return eJSON.Address, nil
}

// struct required for reading "address" from this bit of encrypted JSON:
// JSON = '{"publicKey":"22b4b2618de6dc8254d76276d51f6a9d53471d5b2465c8cae237f21425b10b7d"}'
type solJSON struct {
	PublicKey string `json:"publicKey"`
}

func publicSolKeyFromEncryptedJSON(jsonString string) (solana.PublicKey, error) {
	var eJSON solJSON
	err := json.Unmarshal([]byte(jsonString), &eJSON)
	if err != nil {
		return solana.PublicKey{}, errors.Wrap(err, "failed to unmarshal Solana json")
	}

	b, err := hex.DecodeString(eJSON.PublicKey)
	if err != nil {
		return solana.PublicKey{}, errors.Wrap(err, "invalid hex string for public key solana")
	}

	return solana.PublicKeyFromBytes(b), nil
}

// struct required for reading "peerID" from this bit of encrypted JSON:
// JSON = '{"keyType":"P2P","publicKey":"f3c458c9064bdde449a3904ba8d3f8f5ebf79623077430325252c3368f920199","peerID":"p2p_12D3KooWSDvtYVF3FoyGeMrmDxYeJZMzbEyMHRwmf5GUSqgJhST2"}'
type p2pJSON struct {
	PeerID string `json:"peerID"`
}

func publicP2PAddressFromEncryptedJSON(jsonString string) (string, error) {
	var pJSON p2pJSON

	err := json.Unmarshal([]byte(jsonString), &pJSON)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal p2p json")
	}

	return pJSON.PeerID, nil
}

type dkgRecipientJSON struct {
	PublicKey string `json:"publicKey"`
}

func publicDKGRecipientKeyFromEncryptedJSON(jsonString string) (dkgocrtypes.P256ParticipantPublicKey, error) {
	var dJSON dkgRecipientJSON

	err := json.Unmarshal([]byte(jsonString), &dJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dkg recipient key json")
	}

	return hex.DecodeString(dJSON.PublicKey)
}

func ImportNodeKeys(secretsToml string) (*NodeKeys, error) {
	keys := &NodeKeys{
		EVM:    make(map[uint64]*crypto.EVMKey),
		Solana: make(map[string]*crypto.SolKey),
	}

	var sSecrets secrets
	unmarshallErr := toml.Unmarshal([]byte(secretsToml), &sSecrets)
	if unmarshallErr != nil {
		return nil, errors.Wrap(unmarshallErr, "failed to unmarshal TOML secrets")
	}

	if sSecrets.P2PKey.JSON == nil || sSecrets.P2PKey.Password == nil {
		return nil, errors.New("P2P key or password is nil")
	}

	peerID, peerIDErr := publicP2PAddressFromEncryptedJSON(*sSecrets.P2PKey.JSON)
	if peerIDErr != nil {
		return nil, errors.Wrapf(peerIDErr, "failed to get public p2p address for node from encrypted JSON")
	}

	p := new(p2pkey.PeerID)
	if err := p.UnmarshalString(peerID); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal PeerID")
	}

	keys.P2PKey = &crypto.P2PKey{
		EncryptedJSON: []byte(*sSecrets.P2PKey.JSON),
		Password:      *sSecrets.P2PKey.Password,
		PeerID:        *p,
	}

	if sSecrets.DKGRecipientKey.JSON != nil {
		keys.DKGKey = &crypto.DKGRecipientKey{
			EncryptedJSON: []byte(*sSecrets.DKGRecipientKey.JSON),
			Password:      *sSecrets.DKGRecipientKey.Password,
		}
		dkgRecipientPubKey, err := publicDKGRecipientKeyFromEncryptedJSON(*sSecrets.DKGRecipientKey.JSON)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get public DKG recipient key from encrypted JSON")
		}
		keys.DKGKey.PubKey = dkgRecipientPubKey
	}

	for _, evmKey := range sSecrets.EVM.Keys {
		if evmKey.JSON == nil || evmKey.Password == nil || evmKey.ID == nil {
			return nil, errors.New("EVM key or password or ID is nil")
		}

		publicEVMAddress, publicEVMAddressErr := publicEVMAddressFromEncryptedJSON(*evmKey.JSON)
		if publicEVMAddressErr != nil {
			return nil, errors.Wrapf(publicEVMAddressErr, "failed to get public evm address from encrypted JSON")
		}

		if _, ok := keys.EVM[*evmKey.ID]; !ok {
			keys.EVM[*evmKey.ID] = &crypto.EVMKey{}
		}

		keys.EVM[*evmKey.ID] = &crypto.EVMKey{
			EncryptedJSON: []byte(*evmKey.JSON),
			PublicAddress: common.HexToAddress(publicEVMAddress),
			Password:      *evmKey.Password,
		}
	}

	for _, solKey := range sSecrets.Solana.Keys {
		if solKey.JSON == nil || solKey.Password == nil || solKey.ID == nil {
			return nil, errors.New("solana key or password or id is nil")
		}

		publicSolAddr, addrErr := publicSolKeyFromEncryptedJSON(*solKey.JSON)
		if addrErr != nil {
			return nil, errors.Wrapf(addrErr, "failed to get public Solana address from encrypted JSON")
		}

		keys.Solana[*solKey.ID] = &crypto.SolKey{
			EncryptedJSON: []byte(*solKey.JSON),
			PublicAddress: publicSolAddr,
			Password:      *solKey.Password,
		}
	}

	return keys, nil
}
