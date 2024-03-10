package types

import "github.com/ethereum/go-ethereum/common"

type MercuryServerType string

const (
	MS_WSRPC MercuryServerType = "wsrpc"
	MS_WS    MercuryServerType = "ws"
	MS_REST  MercuryServerType = "rest"
	MS_ALL   MercuryServerType = "all"
)

type MercuryServerOpts struct {
	Server struct {
		DevMode             bool
		AutomaticMigrations bool
		Service             string
		Port                string
	}
	RPC struct {
		PrivateKey  string
		NodePubKeys []string
		Port        string
	}
	Database struct {
		Url               string
		WriterInstanceUrl string
		EncryptionKey     string
	}
	Bootstrap struct {
		Username string
		Password string
	}
	WSRPCUrlInternal string
	WSRPCUrlExternal string
}

type User struct {
	Id       string
	Username string
	Password string
}

type MercuryOCRConfig struct {
	Signers               []common.Address
	Transmitters          [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}
