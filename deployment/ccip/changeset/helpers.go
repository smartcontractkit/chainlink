package changeset

import (
	"fmt"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/rpc"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	solanaStateUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
)

func GetPingPongDemoContractAddress(env deployment.Environment, from uint64, to uint64) (string, error) {
	addrs, err := env.ExistingAddresses.AddressesForChain(from)
	if err != nil {
		return "", fmt.Errorf("failed to get addresses for chain %d: %w", from, err)
	}

	for addr, tv := range addrs {
		if tv.Type != PingPongDemo {
			continue
		}

		for label := range tv.Labels {
			if label == fmt.Sprintf("To - %d", to) {
				return addr, nil
			}
		}
	}

	return "", fmt.Errorf("no ping pong demo contract found for chain %d and destination %d", from, to)
}

func GetPaddedPingPongAddressBytes(env deployment.Environment, from uint64, to uint64, familyChainEncoding string) ([]byte, error) {
	counterpartAddressStr, err := GetPingPongDemoContractAddress(env, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get counterpart address: %w", err)
	}

	var addressBytes []byte

	switch familyChainEncoding {
	case chainsel.FamilyEVM:
		counterpartAddressBytes := make([]byte, 32)
		copy(counterpartAddressBytes[12:], common.HexToAddress(counterpartAddressStr).Bytes())

		addressBytes = counterpartAddressBytes
	case chainsel.FamilySolana:
		address, err := solana.PublicKeyFromBase58(counterpartAddressStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse address: %w", err)
		}

		addressBytes = address.Bytes()
	default:
		return nil, fmt.Errorf("unsupported chain family: %s", familyChainEncoding)
	}

	return addressBytes, nil
}

type pingPongPDAData struct {
	PPConfigPDA             solana.PublicKey
	FeeBillingSignerPDA     solana.PublicKey
	PPSendSignerPDA         solana.PublicKey
	PPFeeTokenAta           solana.PublicKey
	RouterFeeTokenReceiver  solana.PublicKey
	RouterDestChainStatePDA solana.PublicKey
	RouterNoncePDA          solana.PublicKey
	FqBillingTokenConfigPDA solana.PublicKey
	FqDestChainPDA          solana.PublicKey
	FqLinkTokenConfigPDA    solana.PublicKey
	RMNRemoteConfigPDA      solana.PublicKey
	RMNRemoteCursesPDA      solana.PublicKey
	NameVersionPDA          solana.PublicKey
	ProgramData             struct {
		DataType uint32
		Address  solana.PublicKey
	}
}

func LoadPingPongPDAData(
	env deployment.Environment,
	solChainSelector uint64,
	pingPongProgram solana.PublicKey,
	counterpartChainSelector uint64,
	feesTokenProgram solana.PublicKey,
	feesTokenMint solana.PublicKey,
	linkTokenMint solana.PublicKey,
) (pingPongPDAData, error) {
	s, err := LoadOnchainStateSolana(env)
	if err != nil {
		env.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return pingPongPDAData{}, err
	}

	chain := env.SolChains[solChainSelector]
	chainState := s.SolChains[solChainSelector]

	var data pingPongPDAData

	data.ProgramData, err = GetSolProgramData(env, chain, pingPongProgram)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to get program data: %w", err)
	}

	data.PPConfigPDA, _, err = solanaStateUtils.FindPingPongDemoConfigPDA(pingPongProgram)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find ping pong config PDA: %w", err)
	}

	data.FeeBillingSignerPDA, _, err = solanaStateUtils.FindFeeBillingSignerPDA(chainState.Router)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find fee billing signer PDA: %w", err)
	}

	data.PPSendSignerPDA, _, err = solanaStateUtils.FindPingPongCCIPSendSignerPDA(pingPongProgram)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find ping pong send signer PDA: %w", err)
	}

	data.PPFeeTokenAta, _, err = tokens.FindAssociatedTokenAddress(feesTokenProgram, feesTokenMint, data.PPSendSignerPDA)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find ping pong fee token ATA: %w", err)
	}

	data.RouterFeeTokenReceiver, _, err = tokens.FindAssociatedTokenAddress(feesTokenProgram, feesTokenMint, data.FeeBillingSignerPDA)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find router fee token receiver: %w", err)
	}

	data.RouterDestChainStatePDA, err = solanaStateUtils.FindDestChainStatePDA(counterpartChainSelector, chainState.Router)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find destination chain state PDA: %w", err)
	}

	data.RouterNoncePDA, err = solanaStateUtils.FindNoncePDA(counterpartChainSelector, data.PPSendSignerPDA, chainState.Router)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find router nonce PDA: %w", err)
	}

	data.FqBillingTokenConfigPDA, _, err = solanaStateUtils.FindFqBillingTokenConfigPDA(feesTokenMint, chainState.FeeQuoter)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find FQ billing token config PDA: %w", err)
	}

	data.FqDestChainPDA, _, err = solanaStateUtils.FindFqDestChainPDA(counterpartChainSelector, chainState.FeeQuoter)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find FQ destination chain PDA: %w", err)
	}

	data.FqLinkTokenConfigPDA, _, err = solanaStateUtils.FindFqBillingTokenConfigPDA(linkTokenMint, chainState.FeeQuoter)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find FQ link token config PDA: %w", err)
	}

	data.RMNRemoteConfigPDA, _, err = solanaStateUtils.FindRMNRemoteConfigPDA(chainState.RMNRemote)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find RMN remote config PDA: %w", err)
	}

	data.RMNRemoteCursesPDA, _, err = solanaStateUtils.FindRMNRemoteCursesPDA(chainState.RMNRemote)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find RMN remote curses PDA: %w", err)
	}

	data.NameVersionPDA, _, err = solanaStateUtils.FindNameAndVersionPDA(pingPongProgram)
	if err != nil {
		return pingPongPDAData{}, fmt.Errorf("failed to find name and version PDA: %w", err)
	}

	return data, nil
}

func GetSolProgramSize(e *deployment.Environment, chain deployment.SolChain, programID solana.PublicKey) (int, error) {
	accountInfo, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &rpc.GetAccountInfoOpts{
		Commitment: deployment.SolDefaultCommitment,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}
	if accountInfo == nil {
		return 0, fmt.Errorf("program account not found: %w", err)
	}
	programBytes := len(accountInfo.Value.Data.GetBinary())
	return programBytes, nil
}

func GetSolProgramData(e deployment.Environment, chain deployment.SolChain, programID solana.PublicKey) (struct {
	DataType uint32
	Address  solana.PublicKey
}, error) {
	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &solRpc.GetAccountInfoOpts{
		Commitment: solRpc.CommitmentConfirmed,
	})
	if err != nil {
		return programData, fmt.Errorf("failed to deploy program: %w", err)
	}

	err = solBinary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return programData, fmt.Errorf("failed to unmarshal program data: %w", err)
	}
	return programData, nil
}
