package avro

import (
	"encoding/json"
	"fmt"

	"github.com/linkedin/goavro"
)

func ParseSchema(schema Schema) (jsonEncodedSchema string, codec *goavro.Codec, err error) {
	buf, err := json.Marshal(schema)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encode Avro schema to JSON: %w", err)
	}
	jsonEncodedSchema = string(buf)
	codec, err = goavro.NewCodec(jsonEncodedSchema)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse JSON-encoded Avro schema into a codec: %w", err)
	}
	return jsonEncodedSchema, codec, nil
}
