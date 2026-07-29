package reconciler

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

// FileWorkflowMetadata mirrors the JSON shape read by the node's v2
// file-based workflow source (core/services/workflows/syncer/v2/file_workflow_source.go).
type FileWorkflowMetadata struct {
	WorkflowID   string `json:"workflow_id"`
	Owner        string `json:"owner"`
	CreatedAt    uint64 `json:"created_at"`
	Status       uint8  `json:"status"`
	WorkflowName string `json:"workflow_name"`
	BinaryURL    string `json:"binary_url"`
	ConfigURL    string `json:"config_url"`
	Tag          string `json:"tag"`
	DonFamily    string `json:"don_family"`
}

// FileWorkflowSourceData is the root structure of the registry JSON file.
type FileWorkflowSourceData struct {
	Workflows []FileWorkflowMetadata `json:"workflows"`
}

// PodTarget identifies a single pod to push workflow artifacts to.
type PodTarget struct {
	Namespace string
	PodName   string
}

// WorkflowDeployInputs are the parameters for WorkflowDeploy.
type WorkflowDeployInputs struct {
	WorkflowFilePath string // path to the workflow source file (.go or .ts)
	WorkflowName     string
	Owner            string // hex-encoded owner address, with or without 0x
	DonFamily        string
	ConfigPath       string // optional path to a workflow config file
	Tag              string // registry entry tag, default "v1.0.0" if empty
	RemoteDir        string // directory on the pod both files are copied into
	Container        string // container name to exec into
	Pods             []PodTarget
}

// WorkflowDeployResult is returned by WorkflowDeploy.
type WorkflowDeployResult struct {
	WorkflowID string
	Failures   map[string]error // pod identity ("namespace/pod") -> error, only for failed copies
}

// WorkflowDeploy compiles a workflow to WASM, brotli-compresses it, computes
// its workflow ID, builds a single-entry private file-registry JSON, and
// copies the registry file + compiled binary (+ optional config) into
// RemoteDir on every target pod. No pod restart is triggered — the node's
// v2 file workflow source re-reads the registry file and re-fetches the
// binary/config on its next poll tick.
func WorkflowDeploy(ctx context.Context, k8s *infra.K8sClient, in WorkflowDeployInputs, log zerolog.Logger) (*WorkflowDeployResult, error) {
	if len(in.Pods) == 0 {
		return nil, errors.New("no target pods provided")
	}
	if in.RemoteDir == "" {
		return nil, errors.New("remote dir is required")
	}

	tag := in.Tag
	if tag == "" {
		tag = "v1.0.0"
	}

	buildDir, err := os.MkdirTemp("", "griddle-workflow-deploy-*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create build dir")
	}
	defer func() { _ = os.RemoveAll(buildDir) }()

	// (a) compile -> wasm -> brotli+base64
	compressedPath, err := creworkflow.CompileWorkflowToDir(ctx, in.WorkflowFilePath, in.WorkflowName, buildDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile workflow")
	}
	log.Info().Str("path", compressedPath).Msg("Compiled workflow to brotli-compressed WASM")

	// (b) recover the binary bytes used for workflow ID computation, and read config.
	//
	// This must match core/services/workflows/artifacts/v2/store.go's
	// FetchWorkflowArtifacts, which only base64-decodes the fetched binary_url
	// content before hashing it for verification in handler.go's tryEngineCreate —
	// it never brotli-decompresses. So the ID must be computed over the
	// brotli-compressed-but-base64-decoded bytes, not the fully decompressed WASM,
	// or the node's own recomputed hash will never match ours.
	compressedWasm, err := decodeBase64File(compressedPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode compiled workflow binary")
	}

	var config []byte
	if in.ConfigPath != "" {
		config, err = os.ReadFile(in.ConfigPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read config file %s", in.ConfigPath)
		}
	}

	ownerHex := strings.TrimPrefix(in.Owner, "0x")
	ownerBytes, err := hex.DecodeString(ownerHex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode owner hex")
	}

	workflowID, err := pkgworkflows.GenerateWorkflowID(ownerBytes, in.WorkflowName, compressedWasm, config, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate workflow ID")
	}
	workflowIDHex := hex.EncodeToString(workflowID[:])
	log.Info().Str("workflowID", workflowIDHex).Msg("Computed workflow ID")

	// (c) build the single-workflow registry file
	binaryFilename := in.WorkflowName + ".br.b64"
	binaryDestPath := filepath.Join(buildDir, binaryFilename)
	if err := copyFile(compressedPath, binaryDestPath); err != nil {
		return nil, errors.Wrap(err, "failed to stage binary for registry")
	}

	localPaths := []string{binaryDestPath}
	remoteDir := strings.TrimSuffix(in.RemoteDir, "/")

	entry := FileWorkflowMetadata{
		WorkflowID:   workflowIDHex,
		Owner:        ownerHex,
		CreatedAt:    uint64(nowUnix()),
		Status:       0,
		WorkflowName: in.WorkflowName,
		BinaryURL:    "file://" + remoteDir + "/" + binaryFilename,
		Tag:          tag,
		DonFamily:    in.DonFamily,
	}

	if in.ConfigPath != "" {
		configFilename := in.WorkflowName + ".config" + filepath.Ext(in.ConfigPath)
		configDestPath := filepath.Join(buildDir, configFilename)
		if err := copyFile(in.ConfigPath, configDestPath); err != nil {
			return nil, errors.Wrap(err, "failed to stage config for registry")
		}
		localPaths = append(localPaths, configDestPath)
		entry.ConfigURL = "file://" + remoteDir + "/" + configFilename
	}

	registryData := FileWorkflowSourceData{Workflows: []FileWorkflowMetadata{entry}}
	registryJSON, err := json.MarshalIndent(registryData, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal registry JSON")
	}

	registryPath := filepath.Join(buildDir, "registry.json")
	if err := os.WriteFile(registryPath, registryJSON, 0o600); err != nil {
		return nil, errors.Wrap(err, "failed to write registry file")
	}
	localPaths = append(localPaths, registryPath)

	// (d) copy registry + binary (+ config) to every pod
	result := &WorkflowDeployResult{WorkflowID: workflowIDHex, Failures: map[string]error{}}
	for _, pod := range in.Pods {
		podID := pod.Namespace + "/" + pod.PodName
		log.Info().Str("pod", podID).Msg("Copying workflow artifacts to pod")
		if err := k8s.CopyFilesToPod(ctx, pod.Namespace, pod.PodName, in.Container, in.RemoteDir, localPaths); err != nil {
			log.Warn().Err(err).Str("pod", podID).Msg("Failed to copy workflow artifacts")
			result.Failures[podID] = err
			continue
		}
		log.Info().Str("pod", podID).Msg("Copied workflow artifacts")
	}

	if len(result.Failures) > 0 {
		return result, fmt.Errorf("failed to copy workflow artifacts to %d/%d pods", len(result.Failures), len(in.Pods))
	}

	return result, nil
}

// decodeBase64File base64-decodes a .br.b64 file's contents, yielding the
// brotli-compressed (but not decompressed) binary bytes — the same
// representation core/services/workflows/artifacts/v2/store.go's
// FetchWorkflowArtifacts produces from binary_url, which is what the node
// actually hashes for workflow ID verification.
func decodeBase64File(path string) ([]byte, error) {
	encoded, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", path)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, errors.Wrap(err, "failed to base64-decode binary")
	}
	return decoded, nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", src)
	}
	return os.WriteFile(dst, data, 0o644) //nolint:gosec // artifact readable by node process
}

// nowUnix is a thin indirection so tests can stub the clock if needed.
var nowUnix = func() int64 { return time.Now().Unix() }

// ParsePodTarget parses a "namespace/podName" string into a PodTarget.
func ParsePodTarget(s string) (PodTarget, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return PodTarget{}, fmt.Errorf("invalid pod target %q: expected format namespace/podName", s)
	}
	return PodTarget{Namespace: parts[0], PodName: parts[1]}, nil
}
