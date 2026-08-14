package report

import (
	"github.com/stellar/go-stellar-sdk/xdr"
)

func EncodeEntries(dataID [32]byte, answer int64, timestamp uint64) ([]byte, error) {
	sign := answer >> 63
	answerVal, err := xdr.NewScVal(xdr.ScValTypeScvI256, xdr.Int256Parts{
		HiHi: xdr.Int64(sign),
		HiLo: xdr.Uint64(sign),
		LoHi: xdr.Uint64(sign),
		LoLo: xdr.Uint64(answer),
	})
	if err != nil {
		return nil, err
	}
	dataIDVal, err := xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(dataID[:]))
	if err != nil {
		return nil, err
	}
	tsVal, err := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(timestamp))
	if err != nil {
		return nil, err
	}

	entries := xdr.ScMap{}
	for _, kv := range []struct {
		k string
		v xdr.ScVal
	}{{"answer", answerVal}, {"data_id", dataIDVal}, {"timestamp", tsVal}} {
		key, kerr := xdr.NewScVal(xdr.ScValTypeScvSymbol, xdr.ScSymbol(kv.k))
		if kerr != nil {
			return nil, kerr
		}
		entries = append(entries, xdr.ScMapEntry{Key: key, Val: kv.v})
	}
	entryVal, err := xdr.NewScVal(xdr.ScValTypeScvMap, &entries)
	if err != nil {
		return nil, err
	}
	vec := xdr.ScVec{entryVal}
	v, err := xdr.NewScVal(xdr.ScValTypeScvVec, &vec)
	if err != nil {
		return nil, err
	}
	return v.MarshalBinary()
}
