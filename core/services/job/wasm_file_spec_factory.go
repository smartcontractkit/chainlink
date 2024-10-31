package job

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/andybalholm/brotli"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type WasmFileSpecFactory struct{}

func (w WasmFileSpecFactory) Spec(ctx context.Context, workflow, configLocation string) (sdk.WorkflowSpec, []byte, string, error) {
	config, err := w.Config(ctx, configLocation)
	if err != nil {
		return sdk.WorkflowSpec{}, nil, "", err
	}

	compressedBinary, sha, err := w.rawSpecAndSha(workflow, config)
	if err != nil {
		return sdk.WorkflowSpec{}, nil, "", err
	}

	moduleConfig := &host.ModuleConfig{Logger: logger.NullLogger}
	spec, err := host.GetWorkflowSpec(ctx, moduleConfig, compressedBinary, config)
	if err != nil {
		return sdk.WorkflowSpec{}, nil, "", err
	} else if spec == nil {
		return sdk.WorkflowSpec{}, nil, "", errors.New("workflow spec not found when running wasm")
	}

	return *spec, compressedBinary, sha, nil
}

func (w WasmFileSpecFactory) RawSpec(ctx context.Context, workflow, configLocation string) ([]byte, error) {
	config, err := w.Config(ctx, configLocation)
	if err != nil {
		return nil, err
	}

	raw, _, err := w.rawSpecAndSha(workflow, config)
	return raw, err
}

func (w WasmFileSpecFactory) Config(_ context.Context, configLocation string) ([]byte, error) {
	config, err := os.ReadFile(configLocation)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// rawSpecAndSha returns the brotli compressed version of the raw wasm file, alongside the sha256 hash of the raw wasm file
func (w WasmFileSpecFactory) rawSpecAndSha(wf string, config []byte) ([]byte, string, error) {
	read, err := os.ReadFile(wf)
	if err != nil {
		return nil, "", err
	}

	extension := strings.ToLower(path.Ext(wf))
	switch extension {
	case ".wasm", "":
		return w.rawSpecAndShaFromWasm(read, config)
	case ".br":
		return w.rawSpecAndShaFromBrotli(read, config)
	default:
		return nil, "", fmt.Errorf("unsupported file type %s", extension)
	}
}

func (w WasmFileSpecFactory) rawSpecAndShaFromBrotli(wasm, config []byte) ([]byte, string, error) {
	brr := brotli.NewReader(bytes.NewReader(wasm))
	rawWasm, err := io.ReadAll(brr)
	if err != nil {
		return nil, "", err
	}

	return wasm, w.sha(rawWasm, config), nil
}

func (w WasmFileSpecFactory) rawSpecAndShaFromWasm(wasm, config []byte) ([]byte, string, error) {
	var b bytes.Buffer
	bwr := brotli.NewWriter(&b)
	if _, err := bwr.Write(wasm); err != nil {
		return nil, "", err
	}

	if err := bwr.Close(); err != nil {
		return nil, "", err
	}

	return b.Bytes(), w.sha(wasm, config), nil
}

func (w WasmFileSpecFactory) sha(wasm, config []byte) string {
	sum := sha256.New()
	sum.Write(wasm)
	sum.Write(config)
	return fmt.Sprintf("%x", sum.Sum(nil))
}

var _ WorkflowSpecFactory = (*WasmFileSpecFactory)(nil)
