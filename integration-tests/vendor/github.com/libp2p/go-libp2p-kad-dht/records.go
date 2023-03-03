package dht

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	ci "github.com/libp2p/go-libp2p-core/crypto"
)

type pubkrs struct {
	pubk ci.PubKey
	err  error
}

// GetPublicKey gets the public key when given a Peer ID. It will extract from
// the Peer ID if inlined or ask the node it belongs to or ask the DHT.
func (dht *IpfsDHT) GetPublicKey(ctx context.Context, p peer.ID) (ci.PubKey, error) {
	logger.Debugf("getPublicKey for: %s", p)

	// Check locally. Will also try to extract the public key from the peer
	// ID itself if possible (if inlined).
	pk := dht.peerstore.PubKey(p)
	if pk != nil {
		return pk, nil
	}

	// Try getting the public key both directly from the node it identifies
	// and from the DHT, in parallel
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	resp := make(chan pubkrs, 2)
	go func() {
		pubk, err := dht.getPublicKeyFromNode(ctx, p)
		resp <- pubkrs{pubk, err}
	}()

	// Note that the number of open connections is capped by the dial
	// limiter, so there is a chance that getPublicKeyFromDHT(), which
	// potentially opens a lot of connections, will block
	// getPublicKeyFromNode() from getting a connection.
	// Currently this doesn't seem to cause an issue so leaving as is
	// for now.
	go func() {
		pubk, err := dht.getPublicKeyFromDHT(ctx, p)
		resp <- pubkrs{pubk, err}
	}()

	// Wait for one of the two go routines to return
	// a public key (or for both to error out)
	var err error
	for i := 0; i < 2; i++ {
		r := <-resp
		if r.err == nil {
			// Found the public key
			err := dht.peerstore.AddPubKey(p, r.pubk)
			if err != nil {
				logger.Errorw("failed to add public key to peerstore", "peer", p)
			}
			return r.pubk, nil
		}
		err = r.err
	}

	// Both go routines failed to find a public key
	return nil, err
}

func (dht *IpfsDHT) getPublicKeyFromDHT(ctx context.Context, p peer.ID) (ci.PubKey, error) {
	// Only retrieve one value, because the public key is immutable
	// so there's no need to retrieve multiple versions
	pkkey := routing.KeyForPublicKey(p)
	val, err := dht.GetValue(ctx, pkkey, Quorum(1))
	if err != nil {
		return nil, err
	}

	pubk, err := ci.UnmarshalPublicKey(val)
	if err != nil {
		logger.Errorf("Could not unmarshal public key retrieved from DHT for %v", p)
		return nil, err
	}

	// Note: No need to check that public key hash matches peer ID
	// because this is done by GetValues()
	logger.Debugf("Got public key for %s from DHT", p)
	return pubk, nil
}

func (dht *IpfsDHT) getPublicKeyFromNode(ctx context.Context, p peer.ID) (ci.PubKey, error) {
	// check locally, just in case...
	pk := dht.peerstore.PubKey(p)
	if pk != nil {
		return pk, nil
	}

	// Get the key from the node itself
	pkkey := routing.KeyForPublicKey(p)
	pmes, err := dht.getValueSingle(ctx, p, pkkey)
	if err != nil {
		return nil, err
	}

	// node doesn't have key :(
	record := pmes.GetRecord()
	if record == nil {
		return nil, fmt.Errorf("node %v not responding with its public key", p)
	}

	pubk, err := ci.UnmarshalPublicKey(record.GetValue())
	if err != nil {
		logger.Errorf("Could not unmarshal public key for %v", p)
		return nil, err
	}

	// Make sure the public key matches the peer ID
	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		logger.Errorf("Could not extract peer id from public key for %v", p)
		return nil, err
	}
	if id != p {
		return nil, fmt.Errorf("public key %v does not match peer %v", id, p)
	}

	logger.Debugf("Got public key from node %v itself", p)
	return pubk, nil
}
