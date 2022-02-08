package monitoring

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/riferrei/srclient"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/avro"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
	"github.com/stretchr/testify/require"
)

const baseSchema = `
{"name": "person", "type": "record",  "fields": [
	{"name": "name", "type": "string"}
]}`

const extendedSchema = `
{"name": "person", "type": "record",  "fields": [
	{"name": "name", "type": "string"},
	{"name": "age", "default": null, "type": ["null","int"]}
]}`

func TestSchemaRegistry(t *testing.T) {
	//defer goleak.VerifyNone(t)

	t.Run("EnsureSchema with mock registry", func(t *testing.T) {
		client := srclient.CreateMockSchemaRegistryClient("http://127.0.0.1:6767")
		registry := &schemaRegistry{client, newNullLogger()}

		newSchema, err := registry.EnsureSchema("test_schema", baseSchema)
		require.NoError(t, err, "error when fetching a new schema")

		existingSchema, err := registry.EnsureSchema("test_schema", baseSchema)
		require.NoError(t, err, "error when fetching existing schema")
		require.Equal(t, newSchema.ID(), existingSchema.ID(), "should return the same schema ID")
		require.Equal(t, newSchema.Version(), existingSchema.Version(), "should return the same schema version")

		extendedSchema, err := registry.EnsureSchema("test_schema", extendedSchema)
		require.NoError(t, err, "error when extending existing schema")
		require.Equal(t, existingSchema.ID()+1, extendedSchema.ID(), "should bump the schema ID")
		require.Equal(t, existingSchema.Version()+1, extendedSchema.Version(), "should bump the version after a schema update")
	})
	t.Run("Encode/Decode", func(t *testing.T) {
		client := srclient.CreateMockSchemaRegistryClient("http://127.0.0.1:6767")
		registry := &schemaRegistry{client, newNullLogger()}
		_, err := client.CreateSchema("person", baseSchema, srclient.Avro)
		require.NoError(t, err)
		schema, err := registry.EnsureSchema("person", baseSchema)
		require.NoError(t, err)

		subject := map[string]interface{}{"name": "test"}
		expectedEncoded := []byte{
			0x0,                // "magic" byte
			0x0, 0x0, 0x0, 0x1, // 4 bytes for schema id
			0x8, 0x74, 0x65, 0x73, 0x74, // avro-encoded payload
		}
		encoded, err := schema.Encode(subject)
		require.NoError(t, err)
		require.Equal(t, expectedEncoded, encoded)

		decoded, err := schema.Decode(encoded)
		require.NoError(t, err)
		require.Equal(t, subject, decoded)
	})
	t.Run("live registry", func(t *testing.T) {
		if _, isPresent := os.LookupEnv("FEATURE_TEST_ONLY_ENV_RUNNING"); !isPresent {
			t.Skip()
		}
		srURL := os.Getenv("SCHEMA_REGISTRY_URL")
		srUsername := os.Getenv("SCHEMA_REGISTRY_USERNAME")
		srPassword := os.Getenv("SCHEMA_REGISTRY_PASSWORD")
		registry := NewSchemaRegistry(config.SchemaRegistry{srURL, srUsername, srPassword}, newNullLogger())

		t.Run("EnsureSchema", func(t *testing.T) {
			defer func() {
				backend := srclient.CreateSchemaRegistryClient(srURL)
				backend.SetCredentials(srUsername, srPassword)
				_ = backend.DeleteSubject("test_schema", true)
			}()

			newSchema, err := registry.EnsureSchema("test_schema", baseSchema)
			require.NoError(t, err, "error when fetching a new schema")

			existingSchema, err := registry.EnsureSchema("test_schema", baseSchema)
			require.NoError(t, err, "error when fetching existing schema")
			require.Equal(t, newSchema.ID(), existingSchema.ID(), "should return the same schema ID")
			require.Equal(t, newSchema.Version(), existingSchema.Version(), "should return the same schema version")

			extendedSchema, err := registry.EnsureSchema("test_schema", extendedSchema)
			require.NoError(t, err, "error when extending existing schema")
			require.Equal(t, existingSchema.ID(), extendedSchema.ID(), "same schema id")
			require.Equal(t, existingSchema.Version(), extendedSchema.Version(), "same version because schemas are full transitive")
		})
		t.Run("test schema compatibility", func(t *testing.T) {
			buf, err := json.Marshal(previousTransmissionAvroSchema)
			require.NoError(t, err, "failed to generate Avro schema for transmission")
			previousTransmissionCompiled := string(buf)

			buf, err = json.Marshal(previousConfigSetSimplifiedAvroSchema)
			require.NoError(t, err, "failed to generate Avro schema for configSimplified")
			previousConfigSetSimplifiedCompiled := string(buf)

			backend := srclient.CreateSchemaRegistryClient(srURL)
			backend.SetCredentials(srUsername, srPassword)

			var transmissionName = "transmission_test_only"
			var configSetSimplifiedName = "config_set_simplified_test_only"

			_ = backend.DeleteSubject(transmissionName, true)
			_ = backend.DeleteSubject(configSetSimplifiedName, true)
			defer func() {
				require.NoError(t, backend.DeleteSubject(transmissionName, true))
				require.NoError(t, backend.DeleteSubject(configSetSimplifiedName, true))
			}()

			// create schemas using the previous versions
			previousTransmissionSchema, err := registry.EnsureSchema(transmissionName, previousTransmissionCompiled)
			require.NoError(t, err)
			previousConfigSetSimplifiedSchema, err := registry.EnsureSchema(configSetSimplifiedName, previousConfigSetSimplifiedCompiled)
			require.NoError(t, err)

			// create schemas using the new production versions
			currentTransmissionSchema, err := registry.EnsureSchema(transmissionName, TransmissionAvroSchema)
			require.NoError(t, err)
			currentConfigSetSimplifiedSchema, err := registry.EnsureSchema(configSetSimplifiedName, ConfigSetSimplifiedAvroSchema)
			require.NoError(t, err)

			require.Equal(t, previousTransmissionSchema.ID(), currentTransmissionSchema.ID())
			require.Equal(t, previousTransmissionSchema.Version(), currentTransmissionSchema.Version())

			// For FULL_TRANSITIVE compatibility, versions don't change!
			require.Equal(t, previousConfigSetSimplifiedSchema.ID(), currentConfigSetSimplifiedSchema.ID())
			require.Equal(t, previousConfigSetSimplifiedSchema.Version(), currentConfigSetSimplifiedSchema.Version())
		})
	})
}

// This section contains previous versions of the schema in schemas.go
// Whenever schemas are updates, check for compatibility by pasting the previsous
// versions here running the test suite above against a running schema registry process.
// NOTE: you must set the FEATURE_TEST_ONLY_LIVE_SCHEMA_REGISTRY and SCHEMA_REGISTRY_URL env vars.

var previousTransmissionAvroSchema = avro.Record("transmission", avro.Opts{Namespace: "link.chain.ocr2"}, avro.Fields{
	avro.Field("block_number", avro.Opts{Doc: "uint64 big endian"}, avro.Bytes),
	avro.Field("answer", avro.Opts{}, avro.Record("answer", avro.Opts{}, avro.Fields{
		avro.Field("data", avro.Opts{Doc: "*big.Int"}, avro.Bytes),
		avro.Field("timestamp", avro.Opts{Doc: "uint32"}, avro.Long),
	})),
	avro.Field("solana_chain_config", avro.Opts{}, avro.Record("solana_chain_config", avro.Opts{}, avro.Fields{
		avro.Field("network_name", avro.Opts{}, avro.String),
		avro.Field("network_id", avro.Opts{}, avro.String),
		avro.Field("chain_id", avro.Opts{}, avro.String),
	})),
	avro.Field("feed_config", avro.Opts{}, avro.Record("feed_config", avro.Opts{}, avro.Fields{
		avro.Field("feed_name", avro.Opts{}, avro.String),
		avro.Field("feed_path", avro.Opts{}, avro.String),
		avro.Field("symbol", avro.Opts{}, avro.String),
		avro.Field("heartbeat_sec", avro.Opts{}, avro.Long),
		avro.Field("contract_type", avro.Opts{}, avro.String),
		avro.Field("contract_status", avro.Opts{}, avro.String),
		avro.Field("contract_address", avro.Opts{Doc: "[32]byte"}, avro.Bytes),
		avro.Field("transmissions_account", avro.Opts{Doc: "[32]byte"}, avro.Bytes),
		avro.Field("state_account", avro.Opts{Doc: "[32]byte"}, avro.Bytes),
	})),
})

var previousConfigSetSimplifiedAvroSchema = avro.Record("config_set_simplified", avro.Opts{Namespace: "link.chain.ocr2"}, avro.Fields{
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
	avro.Field("s", avro.Opts{Doc: "json encoded []int"}, avro.String),
	avro.Field("oracles", avro.Opts{Doc: "json encoded list of oracles' "}, avro.String),
	avro.Field("feed_state_account", avro.Opts{Doc: "[32]byte"}, avro.String),
})
