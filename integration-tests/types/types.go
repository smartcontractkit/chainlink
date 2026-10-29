package types

import "github.com/ethereum/go-ethereum/common"

type MercuryServerType string

const (
	MSWSRPC MercuryServerType = "wsrpc"
	MSWS    MercuryServerType = "ws"
	MSREST  MercuryServerType = "rest"
	MSAll   MercuryServerType = "all"
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
		URL               string
		WriterInstanceURL string
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
	ID       string
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
