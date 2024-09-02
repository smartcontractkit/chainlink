package solclient

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog/log"
	ocr_2 "github.com/smartcontractkit/chainlink-solana/contracts/generated/ocr_2"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type OCRv2 struct {
	Client                   *Client
	ContractDeployer         *ContractDeployer
	State                    *solana.Wallet
	Authorities              map[string]*Authority
	Payees                   []*solana.Wallet
	Owner                    *solana.Wallet
	Proposal                 *solana.Wallet
	Mint                     *solana.Wallet
	OCRVaultAssociatedPubKey solana.PublicKey
	ProgramWallet            *solana.Wallet
}

func (m *OCRv2) ProgramAddress() string {
	return m.ProgramWallet.PublicKey().String()
}

func (m *OCRv2) writeOffChainConfig(ocConfigBytes []byte) error {
	payer := m.Client.DefaultWallet
	return m.Client.TXSync(
		"Write OffChain config chunk",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			ocr_2.NewWriteOffchainConfigInstruction(
				ocConfigBytes,
				m.Proposal.PublicKey(),
				m.Owner.PublicKey(),
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			if key.Equals(m.Proposal.PublicKey()) {
				return &m.Proposal.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
}

func (m *OCRv2) acceptProposal(digest []byte) error {
	payer := m.Client.DefaultWallet
	va := m.ContractDeployer.Accounts.Authorities["vault"]
	return m.Client.TXSync(
		"Accept OffChain config proposal",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			ocr_2.NewAcceptProposalInstruction(
				digest,
				m.State.PublicKey(),
				m.Proposal.PublicKey(),
				m.Owner.PublicKey(),
				m.OCRVaultAssociatedPubKey,
				m.Owner.PublicKey(),
				m.OCRVaultAssociatedPubKey,
				va.PublicKey,
				solana.TokenProgramID,
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			if key.Equals(m.Proposal.PublicKey()) {
				return &m.Proposal.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
}

// SetBilling sets default billing to oracles
func (m *OCRv2) SetBilling(observationPayment uint32, transmissionPayment uint32, controllerAddr string) error {
	payer := m.Client.DefaultWallet
	billingACPubKey, err := solana.PublicKeyFromBase58(controllerAddr)
	if err != nil {
		return nil
	}
	va := m.ContractDeployer.Accounts.Authorities["vault"]
	err = m.Client.TXSync(
		"Set billing",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			ocr_2.NewSetBillingInstruction(
				observationPayment,
				transmissionPayment,
				m.State.PublicKey(),
				m.Owner.PublicKey(),
				m.Owner.PublicKey(),
				billingACPubKey,
				m.OCRVaultAssociatedPubKey, // token vault
				va.PublicKey,               // vault authority
				solana.TokenProgramID,      // token program
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *OCRv2) finalizeOffChainConfig() error {
	payer := m.Client.DefaultWallet
	return m.Client.TXSync(
		"Finalize OffChain config",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			ocr_2.NewFinalizeProposalInstruction(
				m.Proposal.PublicKey(),
				m.Owner.PublicKey(),
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			if key.Equals(m.Proposal.PublicKey()) {
				return &m.Proposal.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
}

func (m *OCRv2) makeDigest() ([]byte, error) {
	proposal, err := m.fetchProposalAccount()
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	hasher.Write(append([]byte{}, uint8(proposal.Oracles.Len)))
	for _, oracle := range proposal.Oracles.Xs[:proposal.Oracles.Len] {
		hasher.Write(oracle.Signer.Key[:])
		hasher.Write(oracle.Transmitter.Bytes())
		hasher.Write(oracle.Payee.Bytes())
	}

	hasher.Write(append([]byte{}, proposal.F))
	hasher.Write(proposal.TokenMint.Bytes())
	header := make([]byte, 8+4)
	binary.BigEndian.PutUint64(header, proposal.OffchainConfig.Version)
	binary.BigEndian.PutUint32(header[8:], uint32(proposal.OffchainConfig.Len))
	hasher.Write(header)
	hasher.Write(proposal.OffchainConfig.Xs[:proposal.OffchainConfig.Len])
	return hasher.Sum(nil), nil
}

func (m *OCRv2) fetchProposalAccount() (*ocr_2.Proposal, error) {
	var proposal ocr_2.Proposal
	resp, err := m.Client.RPC.GetAccountInfoWithOpts(
		context.Background(),
		m.Proposal.PublicKey(),
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		return nil, err
	}
	err = bin.NewBinDecoder(resp.Value.Data.GetBinary()).Decode(&proposal)
	if err != nil {
		return nil, err
	}
	log.Debug().Interface("Proposal", proposal).Msg("OCR2 Proposal dump")
	return &proposal, nil
}

func (m *OCRv2) createProposal(version uint64) error {
	payer := m.Client.DefaultWallet
	programWallet := m.Client.ProgramWallets["ocr2-keypair.json"]
	proposalAccInstruction, err := m.Client.CreateAccInstr(m.Proposal.PublicKey(), OCRProposalAccountSize, programWallet.PublicKey())
	if err != nil {
		return err
	}
	return m.Client.TXSync(
		"Create proposal",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			proposalAccInstruction,
			ocr_2.NewCreateProposalInstruction(
				version,
				m.Proposal.PublicKey(),
				m.Owner.PublicKey(),
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(m.Proposal.PublicKey()) {
				return &m.Proposal.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
}

//// Configure sets offchain config in multiple transactions
//func (m *OCRv2) Configure(cfg contracts.OffChainAggregatorV2Config) error {
//	_, _, _, _, version, cfgBytes, err := confighelper.ContractSetConfigArgsForTests(
//		cfg.DeltaProgress,
//		cfg.DeltaResend,
//		cfg.DeltaRound,
//		cfg.DeltaGrace,
//		cfg.DeltaStage,
//		cfg.RMax,
//		cfg.S,
//		cfg.Oracles,
//		cfg.ReportingPluginConfig,
//		cfg.MaxDurationQuery,
//		cfg.MaxDurationObservation,
//		cfg.MaxDurationReport,
//		cfg.MaxDurationShouldAcceptFinalizedReport,
//		cfg.MaxDurationShouldTransmitAcceptedReport,
//		cfg.F,
//		cfg.OnchainConfig,
//	)
//	if err != nil {
//		return fmt.Errorf("config args: %w", err)
//	}
//	chunks := utils.ChunkSlice(cfgBytes, 1000)
//	if err = m.createProposal(version); err != nil {
//		return fmt.Errorf("createProposal: %w", err)
//	}
//	if err = m.proposeConfig(cfg); err != nil {
//		return fmt.Errorf("proposeConfig: %w", err)
//	}
//	for i, cfgChunk := range chunks {
//		if err = m.writeOffChainConfig(cfgChunk); err != nil {
//			return fmt.Errorf("writeOffchainConfig: (chunk %d) %w", i, err)
//		}
//	}
//	if err = m.finalizeOffChainConfig(); err != nil {
//		return fmt.Errorf("finalizeOffchainConfig: %w", err)
//	}
//	digest, err := m.makeDigest()
//	if err != nil {
//		return fmt.Errorf("makeDigest: %w", err)
//	}
//
//	if err = m.acceptProposal(digest); err != nil {
//		return fmt.Errorf("acceptProposal: %w", err)
//	}
//	return nil
//}

// DumpState dumps all OCR accounts state
func (m *OCRv2) DumpState() error {
	var stateDump ocr_2.State
	err := m.Client.RPC.GetAccountDataInto(
		context.Background(),
		m.State.PublicKey(),
		&stateDump,
	)
	if err != nil {
		return err
	}
	log.Debug().Interface("State", stateDump).Msg("OCR2 State dump")
	return nil
}

func (m *OCRv2) GetContractData(ctx context.Context) (*contracts.OffchainAggregatorData, error) {
	panic("implement me")
}

// ProposeConfig sets oracles with payee addresses
func (m *OCRv2) proposeConfig(ocConfig contracts.OffChainAggregatorV2Config) error {
	log.Info().Str("Program Address", m.ProgramWallet.PublicKey().String()).Msg("Proposing new config")
	payer := m.Client.DefaultWallet
	oracles := make([]ocr_2.NewOracle, 0)
	for _, oc := range ocConfig.Oracles {
		oracle := oc.OracleIdentity
		var keyArr [20]byte
		copy(keyArr[:], oracle.OnchainPublicKey)
		transmitter, err := solana.PublicKeyFromBase58(string(oracle.TransmitAccount))
		if err != nil {
			return err
		}
		oracles = append(oracles, ocr_2.NewOracle{
			Signer:      keyArr,
			Transmitter: transmitter,
		})
	}
	err := m.Client.TXSync(
		"Propose new config",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			ocr_2.NewProposeConfigInstruction(
				oracles,
				uint8(ocConfig.F),
				m.Proposal.PublicKey(),
				m.Owner.PublicKey(),
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(m.Proposal.PublicKey()) {
				return &m.Proposal.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return err
	}
	// set one payee for all
	instr := make([]solana.Instruction, 0)
	// TODO: get associated addr
	payee := solana.NewWallet()
	if err := m.ContractDeployer.AddNewAssociatedAccInstr(payee.PublicKey(), m.Owner.PublicKey(), payee.PublicKey(), &instr); err != nil {
		return err
	}
	payees := make([]solana.PublicKey, 0)
	for i := 0; i < len(oracles); i++ {
		payees = append(payees, payee.PublicKey())
	}
	proposeInstr := ocr_2.NewProposePayeesInstruction(
		m.Mint.PublicKey(),
		m.Proposal.PublicKey(),
		m.Owner.PublicKey())
	// Add payees as remaining accounts
	for i := 0; i < len(payees); i++ {
		proposeInstr.Append(solana.Meta(payees[i]))
	}
	instr = append(instr, proposeInstr.Build())
	return m.Client.TXSync(
		"Set payees",
		rpc.CommitmentConfirmed,
		instr,
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(payee.PublicKey()) {
				return &payee.PrivateKey
			}
			if key.Equals(m.Owner.PublicKey()) {
				return &m.Owner.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
}

func (m *OCRv2) RequestNewRound() error {
	panic("implement me")
}

func (m *OCRv2) Address() string {
	return m.State.PublicKey().String()
}

func (m *OCRv2) TransferOwnership(to string) error {
	panic("implement me")
}

func (m *OCRv2) GetLatestConfigDetails() (map[string]interface{}, error) {
	panic("implement me")
}

func (m *OCRv2) GetOwedPayment(transmitterAddr string) (map[string]interface{}, error) {
	panic("implement me")
}
