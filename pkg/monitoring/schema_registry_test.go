package monitoring

import (
	"testing"

	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/require"
)

const baseSchema = `
{"name": "person", "type": "record",  "fields": [
	{"name": "name", "type": "string"}
]}`

const extendedSchema = `
{"name": "person", "type": "record",  "fields": [
	{"name": "name", "type": "string"},
	{"name": "age", "type": "int"}
]}`

func TestSchemaRegistry(t *testing.T) {
	t.Run("EnsureSchema", func(t *testing.T) {
		client := srclient.CreateMockSchemaRegistryClient("http://127.0.0.1:6767")
		registry := &schemaRegistry{client, newNullLogger()}

		// Note: because the mock schema registry panics(!) upon calling GetLatestSchema()
		// when querying for an inexistent subject, we can't test the case where the
		// schema does not exist and is first created by EnsureSchema!
		newSchema, err := client.CreateSchema("config_set", baseSchema, srclient.Avro)
		require.NoError(t, err, "no error when creating the schema")

		existingSchema, err := registry.EnsureSchema("config_set", baseSchema)
		require.NoError(t, err, "no error when fetching existing schema")
		require.Equal(t, newSchema.ID(), existingSchema.ID(), "should return the same schema ID")
		require.Equal(t, newSchema.Version(), existingSchema.Version(), "should return the same schema version")

		extendedSchema, err := registry.EnsureSchema("config_set", extendedSchema)
		require.NoError(t, err, "no error when extending existing schema")
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
}
