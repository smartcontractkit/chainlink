package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"flag"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

func main() {
	lggr, _ := logger.NewLogger()

	n := flag.Int("n", 1, "how many key pairs to generate")
	flag.Parse()

	for i := 0; i < *n; i++ {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			lggr.Error("error generating key pair ", err)
			return
		}
		lggr.Info("key pair ", i, ":")
		lggr.Info("public key ", hex.EncodeToString(pubKey))
		lggr.Info("private key ", hex.EncodeToString(privKey))

		peerID, err := ragep2ptypes.PeerIDFromPrivateKey(privKey)
		if err != nil {
			lggr.Error("error generating peer ID ", err)
			return
		}
		lggr.Info("peer ID ", peerID.String())
	}
}
