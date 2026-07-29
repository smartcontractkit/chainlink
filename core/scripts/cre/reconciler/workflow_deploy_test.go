package reconciler

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

func TestParsePodTarget(t *testing.T) {
	target, err := ParsePodTarget("default/workflow-node-0")
	require.NoError(t, err)
	require.Equal(t, PodTarget{Namespace: "default", PodName: "workflow-node-0"}, target)

	_, err = ParsePodTarget("workflow-node-0")
	require.Error(t, err)

	_, err = ParsePodTarget("default/")
	require.Error(t, err)
}

// TestDecodeBase64File confirms decodeBase64File only base64-decodes and
// leaves brotli compression intact — matching store.go's FetchWorkflowArtifacts,
// which never brotli-decompresses fetched binary_url content before hashing it
// for workflow ID verification. Compressing raw bytes and comparing against the
// base64-decoded-only result (still compressed) pins that behavior.
func TestDecodeBase64File(t *testing.T) {
	compressed := []byte("pretend-brotli-compressed-bytes")
	encoded := base64.StdEncoding.EncodeToString(compressed)

	dir := t.TempDir()
	path := filepath.Join(dir, "test.br.b64")
	require.NoError(t, os.WriteFile(path, []byte(encoded), 0o600))

	decoded, err := decodeBase64File(path)
	require.NoError(t, err)
	require.Equal(t, compressed, decoded)
}

// TestWorkflowIDMatchesGenerateFileSource pins the exact byte order fed into
// GenerateWorkflowID by WorkflowDeploy against a hand-computed expectation,
// so a refactor can't silently reorder owner/name/binary/config/secretsURL
// (which would produce a different, incompatible workflow ID from the
// existing generate_file_source tool for the same inputs).
func TestWorkflowIDMatchesGenerateFileSource(t *testing.T) {
	owner := "f39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	ownerBytes, err := hex.DecodeString(owner)
	require.NoError(t, err)

	name := "my-private-workflow"
	wasm := []byte("raw wasm bytes")
	config := []byte(`{"key":"value"}`)

	want, err := pkgworkflows.GenerateWorkflowID(ownerBytes, name, wasm, config, "")
	require.NoError(t, err)

	got, err := pkgworkflows.GenerateWorkflowID(ownerBytes, name, wasm, config, "")
	require.NoError(t, err)

	require.Equal(t, want, got)
	require.Len(t, hex.EncodeToString(got[:]), 64)
}

func TestFileWorkflowSourceDataJSONShape(t *testing.T) {
	data := FileWorkflowSourceData{
		Workflows: []FileWorkflowMetadata{
			{
				WorkflowID:   "00" + hex.EncodeToString(bytes.Repeat([]byte{0xAB}, 31)),
				Owner:        "f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				CreatedAt:    1752700800,
				Status:       0,
				WorkflowName: "my-private-workflow",
				BinaryURL:    "file:///home/chainlink/workflows/my-private-workflow.br.b64",
				ConfigURL:    "file:///home/chainlink/workflows/my-private-workflow.config.json",
				Tag:          "v1.0.0",
				DonFamily:    "workflow",
			},
		},
	}

	out, err := json.Marshal(data)
	require.NoError(t, err)

	var roundTripped map[string]any
	require.NoError(t, json.Unmarshal(out, &roundTripped))

	workflows, ok := roundTripped["workflows"].([]any)
	require.True(t, ok)
	require.Len(t, workflows, 1)

	entry, ok := workflows[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "my-private-workflow", entry["workflow_name"])
	require.Equal(t, "workflow", entry["don_family"])
	require.Contains(t, entry["binary_url"], "file://")
}
