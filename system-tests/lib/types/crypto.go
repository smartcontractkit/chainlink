package types

import "github.com/ethereum/go-ethereum/common"

type EVMKeys struct {
	EncryptedJSONs  [][]byte
	PublicAddresses []common.Address
	Password        string
	ChainID         int
}

type P2PKeys struct {
	EncryptedJSONs [][]byte
	PeerIDs        []string
	Password       string
}
