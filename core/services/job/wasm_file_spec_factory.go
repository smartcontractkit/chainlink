package job

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type WasmFileSpecFactory struct{}

func (w WasmFileSpecFactory) Spec(ctx context.Context, lggr logger.Logger, rawSpec, config []byte) (sdk.WorkflowSpec, error) {
	moduleConfig := &host.ModuleConfig{Logger: lggr}
	spec, err := host.GetWorkflowSpec(moduleConfig, rawSpec, config)
	if err != nil {
		return sdk.WorkflowSpec{}, err
	} else if spec == nil {
		return sdk.WorkflowSpec{}, errors.New("workflow spec not found when running wasm")
	}

	return *spec, nil
}

func (w WasmFileSpecFactory) RawSpec(_ context.Context, wf string) ([]byte, error) {
	isRawWasm := strings.HasSuffix(wf, ".wasm")
	read, err := os.ReadFile(wf)
	if err != nil || !isRawWasm {
		return read, err
	}

	var b bytes.Buffer
	bwr := brotli.NewWriter(&b)
	if _, err = bwr.Write(read); err != nil {
		return nil, err
	}

	if err = bwr.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

var _ SDKWorkflowSpecFactory = (*WasmFileSpecFactory)(nil)
