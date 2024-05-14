package workflows

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sigs.k8s.io/yaml"
)

var fixtureDir = "./testdata/fixtures/workflows/"

// yamlFixtureReaderObj reads a yaml fixture file and returns the parsed object
func yamlFixtureReaderObj(t *testing.T, testCase string) func(name string) any {
	testFixtureReader := yamlFixtureReaderBytes(t, testCase)

	return func(name string) any {
		testFileBytes := testFixtureReader(name)

		var testFileYaml any
		err := yaml.Unmarshal(testFileBytes, &testFileYaml)
		require.NoError(t, err)

		return testFileYaml
	}
}

// yamlFixtureReaderBytes reads a yaml fixture file and returns the bytes
func yamlFixtureReaderBytes(t *testing.T, testCase string) func(name string) []byte {
	return func(name string) []byte {
		testFileBytes, err := os.ReadFile(fmt.Sprintf(fixtureDir+"%s/%s.yaml", testCase, name))
		require.NoError(t, err)

		return testFileBytes
	}
}

var transformJSON = cmp.FilterValues(func(x, y []byte) bool {
	return json.Valid(x) && json.Valid(y)
}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
	if err := json.Unmarshal(in, &out); err != nil {
		panic(err) // should never occur given previous filter to ensure valid JSON
	}
	return out
}))

func TestWorkflowSpecMarshalling(t *testing.T) {
	t.Parallel()
	fixtureReader := yamlFixtureReaderBytes(t, "marshalling")

	t.Run("Type coercion", func(t *testing.T) {
		workflowBytes := fixtureReader("workflow_1")

		spec := workflowSpecYaml{}
		err := yaml.Unmarshal(workflowBytes, &spec)
		require.NoError(t, err)

		// Test that our workflowSpec still keeps all of the original data
		var rawSpec interface{}
		err = yaml.Unmarshal(workflowBytes, &rawSpec)
		require.NoError(t, err)

		workflowspecJson, err := json.MarshalIndent(spec, "", "  ")
		require.NoError(t, err)
		rawWorkflowSpecJson, err := json.MarshalIndent(rawSpec, "", "  ")
		require.NoError(t, err)

		if diff := cmp.Diff(rawWorkflowSpecJson, workflowspecJson, transformJSON); diff != "" {
			t.Errorf("ParseWorkflowWorkflowSpecFromString() mismatch (-want +got):\n%s", diff)
			t.FailNow()
		}

		// Spot check some fields
		consensusConfig := spec.Consensus[0].Config
		v, ok := consensusConfig["aggregation_config"]
		require.True(t, ok, "expected aggregation_config to be present in consensus config")

		// the type of the keys present in v should be string rather than a number
		// this is because JSON keys are always strings
		_, ok = v.(map[string]any)
		require.True(t, ok, "expected map[string]interface{} but got %T", v)

		// Make sure we dont have any weird type coercion with possible boolean values
		booleanCoercions, ok := spec.Triggers[0].Config["boolean_coercion"].(map[string]any)
		require.True(t, ok, "expected boolean_coercion to be present in triggers config")

		// check bools
		bools, ok := booleanCoercions["bools"]
		require.True(t, ok, "expected bools to be present in boolean_coercions")
		for _, v := range bools.([]interface{}) {
			_, ok = v.(bool)
			require.True(t, ok, "expected bool but got %T", v)
		}

		// check strings
		strings, ok := booleanCoercions["strings"]
		require.True(t, ok, "expected strings to be present in boolean_coercions")
		for _, v := range strings.([]interface{}) {
			_, ok = v.(string)
			require.True(t, ok, "expected string but got %T", v)
		}

		// check numbers
		numbers, ok := booleanCoercions["numbers"]
		require.True(t, ok, "expected numbers to be present in boolean_coercions")
		for _, v := range numbers.([]interface{}) {
			_, ok = v.(int64)
			require.True(t, ok, "expected int64 but got %T", v)
		}
	})

	t.Run("Table and string capability id", func(t *testing.T) {
		workflowBytes := fixtureReader("workflow_2")

		spec := workflowSpecYaml{}
		err := yaml.Unmarshal(workflowBytes, &spec)
		require.NoError(t, err)

		// Test that our workflowSpec still keeps all of the original data
		var rawSpec interface{}
		err = yaml.Unmarshal(workflowBytes, &rawSpec)
		require.NoError(t, err)

		workflowspecJson, err := json.MarshalIndent(spec, "", "  ")
		require.NoError(t, err)
		rawWorkflowSpecJson, err := json.MarshalIndent(rawSpec, "", "  ")
		require.NoError(t, err)

		if diff := cmp.Diff(rawWorkflowSpecJson, workflowspecJson, transformJSON); diff != "" {
			t.Errorf("ParseWorkflowWorkflowSpecFromString() mismatch (-want +got):\n%s", diff)
			t.FailNow()
		}
	})

	t.Run("Yaml spec to spec", func(t *testing.T) {
		expectedSpecPath := fixtureDir + "marshalling/" + "workflow_2_spec.json"
		workflowBytes := fixtureReader("workflow_2")

		workflowYaml := &workflowSpecYaml{}
		err := yaml.Unmarshal(workflowBytes, workflowYaml)
		require.NoError(t, err)

		workflowSpec := workflowYaml.toWorkflowSpec()
		workflowSpecBytes, err := json.MarshalIndent(workflowSpec, "", "  ")
		require.NoError(t, err)

		// change this to update golden file
		shouldUpdateWorkflowSpec := false
		if shouldUpdateWorkflowSpec {
			err = os.WriteFile(expectedSpecPath, workflowSpecBytes, 0600)
			require.NoError(t, err)
		}

		expectedSpecBytes, err := os.ReadFile(expectedSpecPath)
		require.NoError(t, err)
		diff := cmp.Diff(expectedSpecBytes, workflowSpecBytes, transformJSON)
		if diff != "" {
			t.Errorf("WorkflowYamlSpecToWorkflowSpec() mismatch (-want +got):\n%s", diff)
			t.FailNow()
		}
	})
}

func TestJsonSchema(t *testing.T) {
	t.Parallel()
	t.Run("GenerateJsonSchema", func(t *testing.T) {
		expectedSchemaPath := fixtureDir + "workflow_schema.json"
		generatedSchema, err := GenerateJsonSchema()
		require.NoError(t, err)

		// change this to update golden file
		shouldUpdateSchema := false
		if shouldUpdateSchema {
			err = os.WriteFile(expectedSchemaPath, generatedSchema, 0600)
			require.NoError(t, err)
		}

		expectedSchema, err := os.ReadFile(expectedSchemaPath)
		require.NoError(t, err)
		diff := cmp.Diff(expectedSchema, generatedSchema, transformJSON)
		if diff != "" {
			t.Errorf("GenerateJsonSchema() mismatch (-want +got):\n%s", diff)
			t.FailNow()
		}
	})

	t.Run("ValidateJsonSchema", func(t *testing.T) {
		generatedSchema, err := GenerateJsonSchema()
		require.NoError(t, err)

		// test version regex
		// for keystone, we should support major versions only along with prereleases and build metadata
		t.Run("version", func(t *testing.T) {
			readVersionFixture := yamlFixtureReaderObj(t, "versioning")
			failingFixture1 := readVersionFixture("failing_1")
			failingFixture2 := readVersionFixture("failing_2")
			passingFixture1 := readVersionFixture("passing_1")
			jsonSchema, err := jsonschema.CompileString("github.com/smartcontractkit/chainlink", string(generatedSchema))
			require.NoError(t, err)

			err = jsonSchema.Validate(failingFixture1)
			require.Error(t, err)

			err = jsonSchema.Validate(failingFixture2)
			require.Error(t, err)

			err = jsonSchema.Validate(passingFixture1)
			require.NoError(t, err)
		})

		// test ref regex
		t.Run("ref", func(t *testing.T) {
			readRefFixture := yamlFixtureReaderObj(t, "references")
			failingFixture1 := readRefFixture("failing_1")
			passingFixture1 := readRefFixture("passing_1")
			jsonSchema, err := jsonschema.CompileString("github.com/smartcontractkit/chainlink", string(generatedSchema))
			require.NoError(t, err)

			err = jsonSchema.Validate(failingFixture1)
			require.Error(t, err)

			err = jsonSchema.Validate(passingFixture1)
			require.NoError(t, err)
		})
	})
}

func TestParsesIntsCorrectly(t *testing.T) {
	wf, err := Parse(hardcodedWorkflow)
	require.NoError(t, err)

	n, err := wf.Vertex("evm_median")
	require.NoError(t, err)

	assert.Equal(t, int64(3600), n.Config["aggregation_config"].(map[string]any)["0x1111111111111111111100000000000000000000000000000000000000000000"].(map[string]any)["heartbeat"])
}

func TestMappingCustomType(t *testing.T) {
	m := mapping(map[string]any{})
	data := `
{
	"foo": 100,
	"bar": 100.00,
	"baz": { "gnat": 11.10 }
}`

	err := m.UnmarshalJSON([]byte(data))
	require.NoError(t, err)
	assert.Equal(t, int64(100), m["foo"], m)
	assert.Equal(t, decimal.NewFromFloat(100.00), m["bar"], m)
	assert.Equal(t, decimal.NewFromFloat(11.10), m["baz"].(map[string]any)["gnat"], m)
}
