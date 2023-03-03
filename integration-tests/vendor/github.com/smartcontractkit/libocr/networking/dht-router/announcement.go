package dhtrouter

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/smartcontractkit/libocr/networking/dht-router/serialization"
	"google.golang.org/protobuf/proto"
)

// monotonic counter
type announcementCounter struct {
	userPrefix uint32 // should be 0 most of the time, but when needed user can bump the counter manually. See announcementCounter.Gt().
	value      uint64 // counter value persisted in a db
}

func (n announcementCounter) Gt(other announcementCounter) bool {
	if n.userPrefix > other.userPrefix {
		return true
	}
	return n.userPrefix == other.userPrefix && n.value > other.value
}

type announcement struct {
	Addrs   []ma.Multiaddr      // addresses of a peer
	Counter announcementCounter // counter
}

type signedAnnouncement struct {
	announcement
	PublicKey p2pcrypto.PubKey // PublicKey used to verify Sig
	Sig       []byte           // sig over announcement
}

const (
	// The maximum number of addr an announcement may broadcast
	maxAddrInAnnouncements = 10
	// domain separator for signatures
	announcementDomainSeparator = "announcement OCR v1.0.0"
)

func serdeError(field string) error {
	return fmt.Errorf("invalid pm: %s", field)
}

// Validate and serialize an announcement. Return error on invalid announcements.
func (ann signedAnnouncement) serialize() ([]byte, error) {
	// Require all fields to be non-nil and addrs shorter than maxAddrInAnnouncements
	if ann.Addrs == nil || ann.PublicKey == nil || ann.Sig == nil || len(ann.Addrs) > maxAddrInAnnouncements {
		return nil, errors.New("invalid announcement")
	}

	// verify the signature
	err := ann.verify()
	if err != nil {
		return nil, err
	}

	// addr
	var addrs [][]byte
	for _, a := range ann.Addrs {
		addrBytes, err := a.MarshalBinary()
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addrBytes)
	}

	pkBytes, err := p2pcrypto.MarshalPublicKey(ann.PublicKey)
	if err != nil {
		return nil, err
	}

	pm := serialization.SignedAnnouncement{
		Addrs:      addrs,
		UserPrefix: ann.Counter.userPrefix,
		Counter:    ann.Counter.value,
		PublicKey:  pkBytes,
		Sig:        ann.Sig,
	}

	return proto.Marshal(&pm)
}

func deserializeSignedAnnouncement(binary []byte) (signedAnnouncement, error) {
	pm := serialization.SignedAnnouncement{}
	err := proto.Unmarshal(binary, &pm)
	if err != nil {
		return signedAnnouncement{}, err
	}

	// addr
	if len(pm.Addrs) == 0 {
		return signedAnnouncement{}, serdeError("addrs is empty array")
	}

	if len(pm.Addrs) > maxAddrInAnnouncements {
		return signedAnnouncement{}, serdeError("more addr than maxAddrInAnnouncements")
	}

	var addrs []ma.Multiaddr
	for _, addr := range pm.Addrs {
		mAddr, err := ma.NewMultiaddrBytes(addr)
		if err != nil {
			return signedAnnouncement{}, err
		}
		addrs = append(addrs, mAddr)
	}

	publicKey, err := p2pcrypto.UnmarshalPublicKey(pm.PublicKey)
	if err != nil {
		return signedAnnouncement{}, err
	}

	return signedAnnouncement{
		announcement{
			addrs,
			announcementCounter{
				pm.UserPrefix,
				pm.Counter,
			},
		},
		publicKey,
		pm.Sig,
	}, nil
}

func (ann signedAnnouncement) String() string {
	pkStr := "can't stringify PK"
	if b, err := ann.PublicKey.Bytes(); err == nil {
		pkStr = base64.StdEncoding.EncodeToString(b)
	}
	return fmt.Sprintf("addrs=%s, pk=%s, sig=%s",
		ann.Addrs,
		pkStr,
		base64.StdEncoding.EncodeToString(ann.Sig))
}

// digest returns a deterministic digest used for signing
func (ann announcement) digest() ([]byte, error) {
	// serialize only addrs and the counter
	if ann.Addrs == nil || len(ann.Addrs) > maxAddrInAnnouncements {
		return nil, errors.New("invalid announcement")
	}

	hasher := sha256.New()
	hasher.Write([]byte(announcementDomainSeparator))

	// encode addr length
	err := binary.Write(hasher, binary.LittleEndian, uint32(len(ann.Addrs)))
	if err != nil {
		return nil, err
	}
	// encode addr
	for _, a := range ann.Addrs {
		addr, err := a.MarshalBinary()
		if err != nil {
			return nil, err
		}
		err = binary.Write(hasher, binary.LittleEndian, uint32(len(addr)))
		if err != nil {
			return nil, err
		}
		hasher.Write(addr)
	}

	// counter
	err = binary.Write(hasher, binary.LittleEndian, ann.Counter.userPrefix)
	if err != nil {
		return nil, err
	}
	err = binary.Write(hasher, binary.LittleEndian, ann.Counter.value)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func (ann *announcement) sign(sk p2pcrypto.PrivKey) (signedAnnouncement, error) {
	digest, err := ann.digest()
	if err != nil {
		return signedAnnouncement{}, err
	}

	sig, err := sk.Sign(digest)
	if err != nil {
		return signedAnnouncement{}, err
	}

	return signedAnnouncement{
		*ann,
		sk.GetPublic(),
		sig,
	}, nil
}

func (ann signedAnnouncement) verify() error {
	if ann.Sig == nil {
		return errors.New("nil sig")
	}

	msg, err := ann.digest()
	if err != nil {
		return err
	}

	verified, err := ann.PublicKey.Verify(msg, ann.Sig)
	if err != nil {
		return err
	}

	if !verified {
		return errors.New("invalid signature")
	}

	return nil
}
