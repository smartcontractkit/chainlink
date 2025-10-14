package sui

import (
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/blake2b"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

var _ cldf.ChangeSetV2[v1_6.SetOCR3OffRampConfig] = SetOCR3Offramp{}

type SetOCR3Offramp struct{}

// Ed25519Scheme Ed25519 signature scheme flag
// https://docs.sui.io/concepts/cryptography/transaction-auth/keys-addresses#address-format
const Ed25519Scheme byte = 0x00

// Apply implements deployment.ChangeSetV2.
func (s SetOCR3Offramp) Apply(e cldf.Environment, config v1_6.SetOCR3OffRampConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	suiState, err := sui_deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()

	for _, remoteSelector := range config.RemoteChainSels {
		suiChains := e.BlockChains.SuiChains()
		suiChain := suiChains[remoteSelector]
		suiSigner := suiChain.Signer

		deps := Deps{
			AB: ab,
			SuiChain: sui_ops.OpTxDeps{
				Client: suiChain.Client,
				Signer: suiSigner,
				GetCallOpts: func() *bind.CallOpts {
					b := uint64(400_000_000)
					return &bind.CallOpts{
						WaitForExecution: true,
						GasBudget:        &b,
					}
				},
			},
			CCIPOnChainState: state,
		}

		// DonIds for the chain
		donID, err := internal.DonIDForChain(deps.CCIPOnChainState.Chains[config.HomeChainSel].CapabilityRegistry,
			deps.CCIPOnChainState.Chains[config.HomeChainSel].CCIPHome,
			remoteSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		ocr3Args, err := internal.BuildSetOCR3ConfigArgsSui(
			donID,
			deps.CCIPOnChainState.Chains[config.HomeChainSel].CCIPHome,
			remoteSelector,
			globals.ConfigTypeActive,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		var commitArgs *internal.MultiOCR3BaseOCRConfigArgsSui
		var execArgs *internal.MultiOCR3BaseOCRConfigArgsSui
		for _, ocr3Arg := range ocr3Args {
			switch ocr3Arg.OcrPluginType {
			case uint8(types.PluginTypeCCIPCommit):
				commitArgs = &ocr3Arg
			case uint8(types.PluginTypeCCIPExec):
				execArgs = &ocr3Arg
			default:
				return cldf.ChangesetOutput{}, err
			}
		}

		// convert transmitters to account address
		var commitTransmitters []string

		for _, transmitter := range commitArgs.Transmitters {
			// 1) Strip any “0x” prefix
			clean := strings.TrimPrefix(transmitter, "0x")

			// 2) Decode the clean hex into bytes
			pubKeyBytes, err := hex.DecodeString(clean)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to decode transmitter %q: %w", transmitter, err)
			}
			flagged := append([]byte{Ed25519Scheme}, pubKeyBytes...)

			hash := blake2b.Sum256(flagged)
			addr := "0x" + hex.EncodeToString(hash[:])
			commitTransmitters = append(commitTransmitters, addr)
		}

		setOCR3ConfigCommitInput := offrampops.SetOCR3ConfigInput{
			OffRampPackageId: suiState[remoteSelector].OffRampAddress,
			OffRampStateId:   suiState[remoteSelector].OffRampStateObjectId,
			OwnerCapObjectId: suiState[remoteSelector].OffRampOwnerCapId,
			CCIPObjectRefId:  suiState[remoteSelector].CCIPObjectRef,
			// commit plugin config
			ConfigDigest:                   commitArgs.ConfigDigest[:],
			OCRPluginType:                  commitArgs.OcrPluginType,
			BigF:                           commitArgs.F,
			IsSignatureVerificationEnabled: commitArgs.IsSignatureVerificationEnabled,
			Signers:                        commitArgs.Signers,
			Transmitters:                   commitTransmitters,
		}

		// convert exec transmitters to account address
		var execTransmitters []string

		for _, transmitter := range execArgs.Transmitters {
			// 1) Strip any “0x” prefix
			clean := strings.TrimPrefix(transmitter, "0x")

			// 2) Decode the clean hex into bytes
			pubKeyBytes, err := hex.DecodeString(clean)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to decode transmitter %q: %w", transmitter, err)
			}

			flagged := append([]byte{Ed25519Scheme}, pubKeyBytes...)

			hash := blake2b.Sum256(flagged)
			addr := "0x" + hex.EncodeToString(hash[:])
			execTransmitters = append(execTransmitters, addr)
		}
		_, err = operations.ExecuteOperation(e.OperationsBundle, offrampops.SetOCR3ConfigOp, deps.SuiChain, setOCR3ConfigCommitInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		setOCR3ConfigExecInput := offrampops.SetOCR3ConfigInput{
			OffRampPackageId: suiState[remoteSelector].OffRampAddress,
			OffRampStateId:   suiState[remoteSelector].OffRampStateObjectId,
			OwnerCapObjectId: suiState[remoteSelector].OffRampOwnerCapId,
			CCIPObjectRefId:  suiState[remoteSelector].CCIPObjectRef,
			// exec plugin config
			ConfigDigest:                   execArgs.ConfigDigest[:],
			OCRPluginType:                  execArgs.OcrPluginType,
			BigF:                           execArgs.F,
			IsSignatureVerificationEnabled: execArgs.IsSignatureVerificationEnabled,
			Signers:                        execArgs.Signers,
			Transmitters:                   execTransmitters,
		}

		_, err = operations.ExecuteOperation(e.OperationsBundle, offrampops.SetOCR3ConfigOp, deps.SuiChain, setOCR3ConfigExecInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (s SetOCR3Offramp) VerifyPreconditions(e cldf.Environment, config v1_6.SetOCR3OffRampConfig) error {
	return nil
}
