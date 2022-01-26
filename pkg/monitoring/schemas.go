package monitoring

import (
	"encoding/json"
	"fmt"

	"github.com/linkedin/goavro"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/avro"
)

// See https://avro.apache.org/docs/current/spec.html#schemas

var transmissionAvroSchema = avro.Record("transmission", avro.Opts{Namespace: "link.chain.ocr2"}, avro.Fields{
	avro.Field("block_number", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("answer", avro.Opts{}, avro.Record("answer", avro.Opts{}, avro.Fields{
		avro.Field("data", avro.Opts{Doc: "*big.avro.Int"}, avro.Bytes),
		avro.Field("timestamp", avro.Opts{Doc: "uint32"}, avro.Long),
		// These fields are made "optional" for backwards compatibility, but they should be set in all cases.
		avro.Field("config_digest", avro.Opts{Doc: "[32]byte encoded as base64", Default: avro.Null}, avro.Union{avro.Null, avro.String}),
		avro.Field("epoch", avro.Opts{Doc: "uint32", Default: avro.Null}, avro.Union{avro.Null, avro.Long}),
		avro.Field("round", avro.Opts{Doc: "uint8", Default: avro.Null}, avro.Union{avro.Null, avro.Int}),
	})),
	avro.Field("chain_config", avro.Opts{Default: avro.Null}, avro.Union{
		avro.Null,
		avro.Record("chain_config", avro.Opts{}, avro.Fields{
			avro.Field("network_name", avro.Opts{}, avro.String),
			avro.Field("network_id", avro.Opts{}, avro.String),
			avro.Field("chain_id", avro.Opts{}, avro.String),
		}),
	}),
	avro.Field("feed_config", avro.Opts{}, avro.Record("feed_config", avro.Opts{}, avro.Fields{
		avro.Field("feed_name", avro.Opts{}, avro.String),
		avro.Field("feed_path", avro.Opts{}, avro.String),
		avro.Field("symbol", avro.Opts{}, avro.String),
		avro.Field("heartbeat_sec", avro.Opts{}, avro.Long),
		avro.Field("contract_type", avro.Opts{}, avro.String),
		avro.Field("contract_status", avro.Opts{}, avro.String),
		avro.Field("contract_address", avro.Opts{Doc: "[32]byte"}, avro.Bytes),
		avro.Field("transmissions_account", avro.Opts{Doc: "[32]byte", Default: avro.Null}, avro.Union{avro.Null, avro.Bytes}),
		avro.Field("state_account", avro.Opts{Doc: "[32]byte", Default: avro.Null}, avro.Union{avro.Null, avro.Bytes}),
	})),
	avro.Field("link_balance", avro.Opts{Default: avro.Null}, avro.Union{
		avro.Null,
		avro.Bytes,
	}),
})

var configSetSimplifiedAvroSchema = avro.Record("config_set_simplified", avro.Opts{Namespace: "link.chain.ocr2"}, avro.Fields{
	avro.Field("config_digest", avro.Opts{Doc: "[32]byte encoded as base64"}, avro.String),
	avro.Field("block_number", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("signers", avro.Opts{Doc: "json encoded array of base64-encoded signing keys"}, avro.String),
	avro.Field("transmitters", avro.Opts{Doc: "json encoded array of base64-encoded transmission keys"}, avro.String),
	avro.Field("f", avro.Opts{Doc: "uint8"}, avro.Int),
	avro.Field("delta_progress", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("delta_resend", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("delta_round", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("delta_grace", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("delta_stage", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("r_max", avro.Opts{Doc: "uint32"}, avro.Long),
	avro.Field("s", avro.Opts{Doc: "json encoded aray of ints"}, avro.String),
	avro.Field("oracles", avro.Opts{Doc: "json encoded list of oracles"}, avro.String),
	avro.Field("feed_state_account", avro.Opts{Doc: "[32]byte"}, avro.String),
})

var (
	// Avro schemas to sync with the registry
	TransmissionAvroSchema        string
	ConfigSetSimplifiedAvroSchema string

	// These codecs are used in tests
	transmissionCodec        *goavro.Codec
	configSetSimplifiedCodec *goavro.Codec
)

func init() {
	var err error
	var buf []byte

	buf, err = json.Marshal(transmissionAvroSchema)
	if err != nil {
		panic(fmt.Errorf("failed to generate Avro schema for transmission: %w", err))
	}
	TransmissionAvroSchema = string(buf)
	transmissionCodec, err = goavro.NewCodec(TransmissionAvroSchema)
	if err != nil {
		panic(fmt.Errorf("failed to parse Avro schema for the latest transmission: %w", err))
	}

	buf, err = json.Marshal(configSetSimplifiedAvroSchema)
	if err != nil {
		panic(fmt.Errorf("failed to generate Avro schema for configSimplified: %w", err))
	}
	ConfigSetSimplifiedAvroSchema = string(buf)
	configSetSimplifiedCodec, err = goavro.NewCodec(ConfigSetSimplifiedAvroSchema)
	if err != nil {
		panic(fmt.Errorf("failed to parse Avro schema for the latest configSetSimplified: %w", err))
	}

	// These codecs are used in tests but not in main, so the linter complains.
	_ = transmissionCodec
	_ = configSetSimplifiedCodec
}
