package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

type PeerID peer.ID

func (p PeerID) String() string {
	return peer.ID(p).String()
}

func (p *PeerID) UnmarshalText(bs []byte) error {
	peerID, err := peer.Decode(string(bs))
	if err != nil {
		return errors.Wrapf(err, `PeerID#UnmarshalText("%v")`, string(bs))
	}
	*p = PeerID(peerID)
	return nil
}

func (p *PeerID) Scan(value interface{}) error {
	s, is := value.(string)
	if !is {
		return errors.Errorf("PeerID#Scan got %T, expected string", value)
	}
	*p = PeerID("")
	return p.UnmarshalText([]byte(s))
}

func (p PeerID) Value() (driver.Value, error) {
	return peer.Encode(peer.ID(p)), nil
}

func (p PeerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(peer.Encode(peer.ID(p)))
}

func (p *PeerID) UnmarshalJSON(input []byte) error {
	var result string
	if err := json.Unmarshal(input, &result); err != nil {
		return err
	}

	peerId, err := peer.Decode(result)
	if err != nil {
		return err
	}

	*p = PeerID(peerId)
	return nil
}
