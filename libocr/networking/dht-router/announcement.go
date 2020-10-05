package dhtrouter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	ma "github.com/multiformats/go-multiaddr"
)

type Announcement struct {
	Addrs     []ma.Multiaddr 	timestamp int64          	Pk        ic.PubKey      	Sig       []byte         }

func (ann Announcement) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{})
	var addrs []string
	for _, a := range ann.Addrs {
		addrs = append(addrs, a.String())
	}
	out["addrs"] = addrs
	out["timestamp"] = ann.timestamp

	pkBytes, err := ic.MarshalPublicKey(ann.Pk)
	if err != nil {
		return nil, err
	}

	out["pk"] = pkBytes
	out["sig"] = ann.Sig
	return json.Marshal(out)
}

func (ann *Announcement) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}

	d := json.NewDecoder(bytes.NewBuffer(b))
	d.UseNumber()
	if err := d.Decode(&data); err != nil {
		panic(err)
	}

	addrs := data["addrs"].([]interface{})
	for _, a := range addrs {
		ann.Addrs = append(ann.Addrs, ma.StringCast(a.(string)))
	}

	timestamp := data["timestamp"].(json.Number)
	i64, _ := strconv.ParseInt(string(timestamp), 10, 64)
	ann.timestamp = i64

	pkBytes, err := base64.StdEncoding.DecodeString(data["pk"].(string))
	if err != nil {
	}

	pk, err := ic.UnmarshalPublicKey(pkBytes)
	if err != nil {
		return nil
	}
	ann.Pk = pk

	sigBytes, err := base64.StdEncoding.DecodeString(data["sig"].(string))
	if err != nil {
	}
	ann.Sig = sigBytes

	return nil
}

func (ann Announcement) serializeForSign() ([]byte, error) {
	b1, err := json.Marshal(ann.Addrs)
	if err != nil {
		return nil, err
	}

	b2, err := json.Marshal(ann.timestamp)
	if err != nil {
		return nil, err
	}

		return append(b1, b2...), nil
}

func (ann *Announcement) SelfSign(sk ic.PrivKey) error {
	b, err := ann.serializeForSign()
	if err != nil {
		return err
	}

	sig, err := sk.Sign(b)
	if err != nil {
		return err
	}

	ann.Sig = append([]byte{}, sig...)
	return nil
}

func (ann Announcement) SelfVerify() (verified bool, err error) {
	verified = false

	b, err := ann.serializeForSign()
	if err != nil {
		return verified, err
	}

	verified, err = ann.Pk.Verify(b, ann.Sig)
	if err != nil {
		return verified, err
	}

	return verified, nil
}

func (ann Announcement) String() string {
	b, e := ann.Pk.Bytes()
	if e != nil {
		panic(e)
	}
	return fmt.Sprintf("addrs=%s, pk=%s, sig=%s",
		ann.Addrs,
		base64.StdEncoding.EncodeToString(b),
		base64.StdEncoding.EncodeToString(ann.Sig))
}
