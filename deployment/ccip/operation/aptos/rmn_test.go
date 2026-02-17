package aptos

import (
	"context"
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	module_auth "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/auth"
	module_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	module_nonce_manager "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/nonce_manager"
	module_receiver_registry "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/receiver_registry"
	module_rmn_remote "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/rmn_remote"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/token_admin_registry"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
)

func TestCurseMultipleOp(t *testing.T) {
	restore := overrideSeams()
	t.Cleanup(restore)

	ccipAddr := aptos.AccountAddress{}
	subjects := [][]byte{{1, 2, 3}}
	deps := dependency.AptosDeps{AptosChain: cldf_aptos.Chain{Selector: 99}}

	rmn := newStubRMNRemote()
	stub := &stubCCIP{addr: ccipAddr, rmn: rmn}
	ccipBindFn = func(addr aptos.AccountAddress, _ aptos.AptosRpcClient) ccip.CCIP {
		require.Equal(t, ccipAddr, addr)
		return stub
	}

	wantTx := mcmstypes.Transaction{}
	generateMCMSTxFn = func(to aptos.AccountAddress, mi bind.ModuleInformation, fn string, args [][]byte) (mcmstypes.Transaction, error) {
		require.Equal(t, ccipAddr, to)
		require.Equal(t, "rmn_remote", mi.ModuleName)
		require.Equal(t, "curse_multiple", fn)
		require.Equal(t, subjects, args)
		return wantTx, nil
	}

	b := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())

	report, err := cldf_ops.ExecuteOperation(b, CurseMultipleOp, deps, CurseMultipleInput{
		CCIPAddress: ccipAddr,
		Subjects:    subjects,
	})
	require.NoError(t, err)
	require.Equal(t, wantTx, report.Output)
	require.Equal(t, subjects, rmn.enc.lastSubjects)
}

func TestUncurseMultipleOp(t *testing.T) {
	restore := overrideSeams()
	t.Cleanup(restore)

	ccipAddr := aptos.AccountAddress{}
	subjects := [][]byte{{9, 8}}
	deps := dependency.AptosDeps{AptosChain: cldf_aptos.Chain{Selector: 42}}

	rmn := newStubRMNRemote()
	stub := &stubCCIP{addr: ccipAddr, rmn: rmn}
	ccipBindFn = func(addr aptos.AccountAddress, _ aptos.AptosRpcClient) ccip.CCIP {
		require.Equal(t, ccipAddr, addr)
		return stub
	}

	wantTx := mcmstypes.Transaction{}
	generateMCMSTxFn = func(to aptos.AccountAddress, mi bind.ModuleInformation, fn string, args [][]byte) (mcmstypes.Transaction, error) {
		require.Equal(t, ccipAddr, to)
		require.Equal(t, "rmn_remote", mi.ModuleName)
		require.Equal(t, "uncurse_multiple", fn)
		require.Equal(t, subjects, args)
		return wantTx, nil
	}

	b := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())

	report, err := cldf_ops.ExecuteOperation(b, UncurseMultipleOp, deps, UncurseMultipleInput{
		CCIPAddress: ccipAddr,
		Subjects:    subjects,
	})
	require.NoError(t, err)
	require.Equal(t, wantTx, report.Output)
	require.Equal(t, subjects, rmn.enc.lastSubjects)
}

func TestIsSubjectCursed(t *testing.T) {
	restore := overrideSeams()
	t.Cleanup(restore)

	ccipAddr := aptos.AccountAddress{}
	deps := dependency.AptosDeps{AptosChain: cldf_aptos.Chain{Selector: 7}}
	wantSubject := []byte{5, 5, 5}

	rmn := newStubRMNRemote()
	rmn.cursed = true
	stub := &stubCCIP{addr: ccipAddr, rmn: rmn}
	ccipBindFn = func(addr aptos.AccountAddress, _ aptos.AptosRpcClient) ccip.CCIP {
		require.Equal(t, ccipAddr, addr)
		return stub
	}

	isCursed, err := IsSubjectCursed(deps, ccipAddr, wantSubject)
	require.NoError(t, err)
	require.True(t, isCursed)
	require.Equal(t, wantSubject, rmn.lastSubject)
}

// --- test seams and stubs ---

func overrideSeams() func() {
	origBind := ccipBindFn
	origGen := generateMCMSTxFn
	return func() {
		ccipBindFn = origBind
		generateMCMSTxFn = origGen
	}
}

type stubCCIP struct {
	addr aptos.AccountAddress
	rmn  *stubRMNRemote
}

func (s *stubCCIP) Address() aptos.AccountAddress   { return s.addr }
func (s *stubCCIP) Auth() module_auth.AuthInterface { return nil }
func (s *stubCCIP) FeeQuoter() module_fee_quoter.FeeQuoterInterface {
	return nil
}
func (s *stubCCIP) NonceManager() module_nonce_manager.NonceManagerInterface { return nil }
func (s *stubCCIP) ReceiverRegistry() module_receiver_registry.ReceiverRegistryInterface {
	return nil
}
func (s *stubCCIP) TokenAdminRegistry() module_token_admin_registry.TokenAdminRegistryInterface {
	return nil
}
func (s *stubCCIP) RMNRemote() module_rmn_remote.RMNRemoteInterface { return s.rmn }

type stubRMNRemote struct {
	enc         *stubEncoder
	cursed      bool
	lastSubject []byte
}

func newStubRMNRemote() *stubRMNRemote {
	return &stubRMNRemote{enc: &stubEncoder{}}
}

func (s *stubRMNRemote) Encoder() module_rmn_remote.RMNRemoteEncoder { return s.enc }
func (s *stubRMNRemote) IsCursed(_ *bind.CallOpts, subject []byte) (bool, error) {
	s.lastSubject = subject
	return s.cursed, nil
}

// Unused interface methods – keep minimal stubs to satisfy compiler.
func (s *stubRMNRemote) TypeAndVersion(*bind.CallOpts) (string, error) { return "", nil }
func (s *stubRMNRemote) Verify(*bind.CallOpts, aptos.AccountAddress, []uint64, [][]byte, []uint64, []uint64, [][]byte, [][]byte) (bool, error) {
	return false, nil
}
func (s *stubRMNRemote) GetArm(*bind.CallOpts) (aptos.AccountAddress, error) {
	return aptos.AccountAddress{}, nil
}
func (s *stubRMNRemote) GetVersionedConfig(*bind.CallOpts) (uint32, module_rmn_remote.Config, error) {
	return 0, module_rmn_remote.Config{}, nil
}
func (s *stubRMNRemote) GetLocalChainSelector(*bind.CallOpts) (uint64, error) { return 0, nil }
func (s *stubRMNRemote) GetReportDigestHeader(*bind.CallOpts) ([]byte, error) { return nil, nil }
func (s *stubRMNRemote) GetCursedSubjects(*bind.CallOpts) ([][]byte, error)   { return nil, nil }
func (s *stubRMNRemote) IsCursedGlobal(*bind.CallOpts) (bool, error)          { return false, nil }
func (s *stubRMNRemote) IsCursedU128(*bind.CallOpts, *big.Int) (bool, error)  { return false, nil }
func (s *stubRMNRemote) Initialize(*bind.TransactOpts, uint64) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) SetConfig(*bind.TransactOpts, []byte, [][]byte, []uint64, uint64) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) Curse(*bind.TransactOpts, []byte) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) CurseMultiple(*bind.TransactOpts, [][]byte) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) Uncurse(*bind.TransactOpts, []byte) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) UncurseMultiple(*bind.TransactOpts, [][]byte) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) MCMSEntrypoint(aptos.AccountAddress) (*api.PendingTransaction, error) {
	return nil, nil
}
func (s *stubRMNRemote) RegisterMCMSEntrypoint() (*api.PendingTransaction, error) {
	return nil, nil
}

// encoder methods
type stubEncoder struct {
	lastSubjects [][]byte
}

func (s *stubEncoder) CurseMultiple(subjects [][]byte) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	s.lastSubjects = subjects
	return bind.ModuleInformation{ModuleName: "rmn_remote"}, "curse_multiple", nil, subjects, nil
}

func (s *stubEncoder) UncurseMultiple(subjects [][]byte) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	s.lastSubjects = subjects
	return bind.ModuleInformation{ModuleName: "rmn_remote"}, "uncurse_multiple", nil, subjects, nil
}
func (s *stubEncoder) TypeAndVersion() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) Verify(aptos.AccountAddress, []uint64, [][]byte, []uint64, []uint64, [][]byte, [][]byte) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) GetArm() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) GetVersionedConfig() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) GetLocalChainSelector() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) GetReportDigestHeader() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) GetCursedSubjects() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) IsCursedGlobal() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) IsCursed([]byte) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) IsCursedU128(*big.Int) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) Initialize(uint64) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) SetConfig([]byte, [][]byte, []uint64, uint64) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) Curse([]byte) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) Uncurse([]byte) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) MCMSEntrypoint(aptos.AccountAddress) (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
func (s *stubEncoder) RegisterMCMSEntrypoint() (bind.ModuleInformation, string, []aptos.TypeTag, [][]byte, error) {
	return bind.ModuleInformation{}, "", nil, nil, nil
}
