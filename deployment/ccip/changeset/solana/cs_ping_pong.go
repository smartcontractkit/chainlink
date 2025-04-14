package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	solPingPong "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ping_pong_demo"
	solanaStateUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[DeployPingPongContractsConfig] = DeployPingPongContractChangeset
var _ deployment.ChangeSet[StartPingPongConfig] = StartPingPongContractChangeset
var _ deployment.ChangeSet[SetPausePingPongConfig] = SetPausePingPongChangeset

type DeployPingPongContractsConfig struct {
	FromChainSelector  uint64
	ToChainSelector    uint64
	CounterpartAddress []byte
	IsPaused           bool
	ExtraArgs          []byte
	DeployerKey        solana.PublicKey

	FeesTokenProgram solana.PublicKey
	FeesTokenMint    solana.PublicKey
}

func (c DeployPingPongContractsConfig) Validate(e deployment.Environment) error {
	return nil
}
func DeployPingPongContractChangeset(
	e deployment.Environment,
	c DeployPingPongContractsConfig,
) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployPingPongContractsConfig: %w", err)
	}

	s, err := ccipChangeset.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	chainSel := c.FromChainSelector

	chain := e.SolChains[chainSel]
	chainState := s.SolChains[chainSel]

	newAddresses := deployment.NewMemoryAddressBook()

	programName := getTypeToProgramDeployName()[changeset.PingPongDemo]

	// programID := "BjFYvj71HHzrAVjzvLSWuxPXybqAVQXLtnCvq3D9wyY3"
	programID, err := chain.DeployProgram(e.Logger, programName, false)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
	}

	programAddress := solana.MustPublicKeyFromBase58(programID)

	solPingPong.SetProgramID(programAddress)

	e.Logger.Infow("Deployed program", "Program", changeset.PingPongDemo, "addr", programID, "chain", chain.String())

	tv := deployment.NewTypeAndVersion(changeset.PingPongDemo, deployment.Version1_0_0)
	tv.Labels.Add(fmt.Sprintf("From - %d", c.FromChainSelector))
	tv.Labels.Add(fmt.Sprintf("To - %d", c.ToChainSelector))

	err = newAddresses.Save(chain.Selector, programID, tv)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address: %w", err)
	}

	programData, err := solProgramData(e, chain, programAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get solana ping pong program data: %w", err)
	}

	ppConfigPDA, _, err := solanaStateUtils.FindPingPongDemoConfigPDA(programAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	initTx, err := solPingPong.NewInitializeConfigInstruction(
		chainState.Router,
		c.ToChainSelector,
		c.CounterpartAddress,
		c.IsPaused,
		c.ExtraArgs,
		ppConfigPDA,
		c.FeesTokenMint,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		programAddress,
		programData.Address,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initTx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm initialize config instruction: %w", err)
	}

	feeBillingSignerPDA, _, err := solanaStateUtils.FindFeeBillingSignerPDA(chainState.Router)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	ppSendSignerPDA, _, err := solanaStateUtils.FindPingPongCCIPSendSignerPDA(programAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	ppFeeTokenAta, _, err := tokens.FindAssociatedTokenAddress(c.FeesTokenProgram, c.FeesTokenMint, ppSendSignerPDA)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	ppNameVersion, _, err := solanaStateUtils.FindNameAndVersionPDA(programAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong name version pda: %w", err)
	}

	initialize, err := solPingPong.NewInitializeInstruction(
		ppConfigPDA,
		ppNameVersion,
		feeBillingSignerPDA,
		c.FeesTokenProgram,
		c.FeesTokenMint,
		ppFeeTokenAta,
		ppSendSignerPDA,
		c.DeployerKey,
		solana.SPLAssociatedTokenAccountProgramID,
		solana.SystemProgramID,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initialize}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm initialize instruction: %w", err)
	}

	e.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

type StartPingPongConfig struct {
	DeployerKey     solana.PublicKey
	ChainSelector   uint64
	PingPongProgram solana.PublicKey

	FeesTokenProgram solana.PublicKey
	FeesTokenMint    solana.PublicKey

	LinkTokenMint solana.PublicKey
}

func (c StartPingPongConfig) Validate(e deployment.Environment) error {
	return nil
}
func StartPingPongContractChangeset(
	e deployment.Environment,
	c StartPingPongConfig,
) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid StartPingPongConfig: %w", err)
	}

	s, err := ccipChangeset.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[c.ChainSelector]
	chainState := s.SolChains[c.ChainSelector]

	solPingPong.SetProgramID(c.PingPongProgram)

	ppConfigPDA, _, err := solanaStateUtils.FindPingPongDemoConfigPDA(c.PingPongProgram)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	ppSendSignerPDA, _, err := solanaStateUtils.FindPingPongCCIPSendSignerPDA(c.PingPongProgram)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	destChainStatePDA, err := solanaStateUtils.FindDestChainStatePDA(c.ChainSelector, chainState.Router)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	routerNoncePDA, err := solanaStateUtils.FindNoncePDA(c.ChainSelector, ppSendSignerPDA, chainState.Router)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	feeBillingSignerPDA, _, err := solanaStateUtils.FindFeeBillingSignerPDA(chainState.Router)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	fqBillingTokenConfigPDA, _, err := solanaStateUtils.FindFqBillingTokenConfigPDA(c.FeesTokenMint, chainState.FeeQuoter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	rmnRemoteConfigPDA, _, err := solanaStateUtils.FindRMNRemoteConfigPDA(chainState.RMNRemote)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	rmnRemoteCursesPDA, _, err := solanaStateUtils.FindRMNRemoteCursesPDA(chainState.RMNRemote)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	ppFeeTokenAta, _, err := tokens.FindAssociatedTokenAddress(c.FeesTokenProgram, c.FeesTokenMint, ppSendSignerPDA)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	routerFeeTokenReceiver, _, err := tokens.FindAssociatedTokenAddress(c.FeesTokenProgram, c.FeesTokenMint, feeBillingSignerPDA)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find router fee token receiver: %w", err)
	}

	fqDestChainPDA, _, err := solanaStateUtils.FindFqDestChainPDA(c.ChainSelector, chainState.FeeQuoter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find fq dest chain pda: %w", err)
	}

	fqLinkTokenConfigPDA, _, err := solanaStateUtils.FindFqBillingTokenConfigPDA(c.LinkTokenMint, chainState.FeeQuoter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find link config pda: %w", err)
	}

	initTx, err := solPingPong.NewStartPingPongInstruction(
		ppConfigPDA,
		c.DeployerKey,
		ppSendSignerPDA,
		c.FeesTokenProgram,
		c.FeesTokenMint,
		ppFeeTokenAta,
		chainState.Router,
		chainState.RouterConfigPDA,
		destChainStatePDA,
		routerNoncePDA,
		routerFeeTokenReceiver,
		feeBillingSignerPDA,
		chainState.FeeQuoter,
		chainState.FeeQuoterConfigPDA,
		fqDestChainPDA,
		fqBillingTokenConfigPDA,
		fqLinkTokenConfigPDA,
		chainState.RMNRemote,
		rmnRemoteCursesPDA,
		rmnRemoteConfigPDA,
		solana.SystemProgramID,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initTx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm start instruction: %w", err)
	}

	e.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{}, nil
}

type SetPausePingPongConfig struct {
	ChainSelector   uint64
	PingPongProgram solana.PublicKey
	IsPaused        bool
}

func (c SetPausePingPongConfig) Validate(e deployment.Environment) error {
	return nil
}
func SetPausePingPongChangeset(
	e deployment.Environment,
	c SetPausePingPongConfig,
) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid SetPausePingPongConfig: %w", err)
	}

	chain := e.SolChains[c.ChainSelector]

	ppConfigPDA, _, err := solanaStateUtils.FindPingPongDemoConfigPDA(c.PingPongProgram)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	initTx, err := solPingPong.NewSetPausedInstruction(
		c.IsPaused,
		ppConfigPDA,
		chain.DeployerKey.PublicKey(),
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initTx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instruction: %w", err)
	}

	e.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{}, nil
}
