package p2pkey

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

const peerIDPrefix = "p2p_"

type PeerID peer.ID

func MakePeerID(s string) (PeerID, error) {
	var peerID PeerID
	return peerID, peerID.UnmarshalString(s)
}

func (p PeerID) String() string {
	// Handle a zero peerID more gracefully, i.e. print it as empty string rather
	// than `p2p_`
	raw := p.Raw()
	if raw == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", peerIDPrefix, raw)
}

func (p PeerID) Raw() string {
	return peer.ID(p).String()
}

func (p *PeerID) UnmarshalString(s string) error {
	return p.UnmarshalText([]byte(s))
}

func (p *PeerID) MarshalText() ([]byte, error) {
	return []byte(p.Raw()), nil
}

func (p *PeerID) UnmarshalText(bs []byte) error {
	input := string(bs)
	if strings.HasPrefix(input, peerIDPrefix) {
		input = string(bs[len(peerIDPrefix):])
	}

	if input == "" {
		return nil
	}

	peerID, err := peer.Decode(input)
	if err != nil {
		return errors.Wrapf(err, `PeerID#UnmarshalText("%v")`, input)
	}
	*p = PeerID(peerID)
	return nil
}

func (p *PeerID) Scan(value interface{}) error {
	*p = PeerID("")
	switch s := value.(type) {
	case string:
		if s != "" {
			return p.UnmarshalText([]byte(s))
		}
	case nil:
	default:
		return errors.Errorf("PeerID#Scan got %T, expected string", value)
	}
	return nil
}

func (p PeerID) Value() (driver.Value, error) {
	return peer.Encode(peer.ID(p)), nil
}

func (p PeerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *PeerID) UnmarshalJSON(input []byte) error {
	var result string
	if err := json.Unmarshal(input, &result); err != nil {
		return err
	}

	return p.UnmarshalText([]byte(result))
}
