package stellar

import (
	"context"

	protocolrpc "github.com/stellar/go-stellar-sdk/protocols/rpc"
	"github.com/stellar/go-stellar-sdk/xdr"
)

type invocation struct {
	ContractID string
	Function   string
	Args       []xdr.ScVal
}

type fakeInvoker struct {
	calls []invocation
}

func (f *fakeInvoker) InvokeContract(_ context.Context, contractID, fn string, args []xdr.ScVal) (*xdr.ScVal, error) {
	f.calls = append(f.calls, invocation{contractID, fn, args})
	v := xdr.ScVal{Type: xdr.ScValTypeScvVoid}
	return &v, nil
}

func (f *fakeInvoker) SimulateContract(_ context.Context, contractID, fn string, args []xdr.ScVal) (*xdr.ScVal, error) {
	f.calls = append(f.calls, invocation{contractID, fn, args})
	v := xdr.ScVal{Type: xdr.ScValTypeScvVoid}
	return &v, nil
}

func (f *fakeInvoker) GetEvents(_ context.Context, _ string, _ uint32, _ []string) ([]protocolrpc.EventInfo, error) {
	return nil, nil
}

type deployCall struct {
	WasmPath string
	Salt     [32]byte
	Args     []xdr.ScVal
}

type fakeDeployer struct {
	deploys    []deployCall
	contractID string
	wasmHash   xdr.Hash
	uploads    []string // wasm paths passed to UploadContractWASM, in call order
}

func (f *fakeDeployer) DeployContract(ctx context.Context, wasmPath string, salt [32]byte) (string, error) {
	return f.DeployContractWithArgs(ctx, wasmPath, salt, nil)
}

func (f *fakeDeployer) DeployContractWithArgs(_ context.Context, wasmPath string, salt [32]byte, args []xdr.ScVal) (string, error) {
	f.deploys = append(f.deploys, deployCall{wasmPath, salt, args})
	return f.contractID, nil
}

func (f *fakeDeployer) UploadContractWASM(_ context.Context, wasmPath string) (xdr.Hash, error) {
	f.uploads = append(f.uploads, wasmPath)
	return f.wasmHash, nil
}
