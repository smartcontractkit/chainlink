package ocr2keepers

import (
	"bytes"
	"encoding/json"
)

// encode is a convenience method that uses json encoding to
// encode any value to an array of bytes
// gob encoding was compared to json encoding where gob encoding
// was shown to be 8x slower
func encode[T any](value T) ([]byte, error) {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(value); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// decode is a convenience method that uses json encoding to
// decode any value from an array of bytes
func decode[T any](b []byte, value *T) error {
	bts := bytes.NewReader(b)
	dec := json.NewDecoder(bts)
	return dec.Decode(value)
}

func limitedLengthEncode(obs Observation, limit int) ([]byte, error) {
	if len(obs.UpkeepIdentifiers) == 0 {
		return encode(obs)
	}

	var res []byte
	for i := range obs.UpkeepIdentifiers {
		b, err := encode(Observation{
			BlockKey:          obs.BlockKey,
			UpkeepIdentifiers: obs.UpkeepIdentifiers[:i+1],
		})
		if err != nil {
			return nil, err
		}
		if len(b) > limit {
			break
		}
		res = b
	}

	return res, nil
}
