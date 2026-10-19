package solana

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	df_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_store"
	df1_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/store"
)

// DF1 feed (Transmissions) account layout: discriminator, 192 byte header, then the live ring of
// 48 byte transmissions { slot: u64, timestamp: u32, _padding: u32, answer: i128, _padding: [u64; 2] }.
const (
	df1HeaderSize       = 8 + 192
	df1TransmissionSize = 48
)

type CutoverDF1FeedsRequest struct {
	// The data feeds store whose authority becomes the writer. MCMS, when set, means the DF1 feeds'
	// owner is the timelock, so set_writer is proposed instead of sent.
	StoreTarget
	DF1StoreProgramID solana.PublicKey
	Feeds             []DF1Feed

	// The store's latest report must be observed within MaxStaleness, and its answer, rescaled to
	// the DF1 feed's decimals, within MaxDeviationBps of the DF1 feed's latest round.
	MaxStaleness    time.Duration
	MaxDeviationBps uint64
}

var _ cldf.ChangeSetV2[*CutoverDF1FeedsRequest] = CutoverDF1Feeds{}

// CutoverDF1Feeds makes the data feeds store authority the writer of DF1 feeds, so the store's
// report callback writes them from then on. Each feed must already be written through
// (ConfigureDF1ReportCallbacks) and the store's value must be fresh and match the DF1 feed's.
// Feeds already cut over are skipped. Stop the feeds' OCR2 jobs once it has landed.
//
// Under MCMS the preconditions hold when the proposal is built, not when it executes.
type CutoverDF1Feeds struct{}

func (cs CutoverDF1Feeds) VerifyPreconditions(env cldf.Environment, req *CutoverDF1FeedsRequest) error {
	if req.DF1StoreProgramID.IsZero() {
		return errors.New("DF1StoreProgramID cannot be zero")
	}
	if req.MaxStaleness <= 0 {
		return errors.New("MaxStaleness must be positive")
	}
	if err := verifyDF1Feeds(req.Feeds); err != nil {
		return err
	}
	if err := req.verify(env); err != nil {
		return err
	}
	_, err := cs.plan(env, req)
	return err
}

func (cs CutoverDF1Feeds) Apply(env cldf.Environment, req *CutoverDF1FeedsRequest) (cldf.ChangesetOutput, error) {
	ixs, err := cs.plan(env, req)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.send(env, "cut over DF1 feeds to the data feeds store", ixs)
}

func (cs CutoverDF1Feeds) plan(env cldf.Environment, req *CutoverDF1FeedsRequest) ([]solana.Instruction, error) {
	ctx := env.GetContext()
	client := env.BlockChains.SolanaChains()[req.ChainSel].Client
	store, signer, err := req.load(env)
	if err != nil {
		return nil, err
	}
	storeAuthority, err := store.authority()
	if err != nil {
		return nil, err
	}

	var ixs []solana.Instruction
	for _, f := range req.Feeds {
		dataID, err := ParseDataID(f.DataID)
		if err != nil {
			return nil, err
		}
		df1, err := loadDF1Feed(ctx, client, req.DF1StoreProgramID, f.Feed)
		if err != nil {
			return nil, err
		}
		if df1.Header.Writer == storeAuthority {
			env.Logger.Infow("DF1 feed already cut over, skipping", "feed", f.Feed, "dataID", f.DataID)
			continue
		}
		if err := df1.verifyOwner(signer); err != nil {
			return nil, err
		}
		if err := verifyDF1Callback(ctx, client, store, dataID, req.DF1StoreProgramID, f.Feed); err != nil {
			return nil, fmt.Errorf("%s: %w", f.DataID, err)
		}

		feedAddress, err := store.feed(dataID)
		if err != nil {
			return nil, err
		}
		feedInfo, err := getAccount(ctx, client, feedAddress)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to read the store feed: %w", f.DataID, err)
		}
		feed, err := df_store.ParseAccount_FeedData(feedInfo.Data.GetBinary())
		if err != nil {
			return nil, fmt.Errorf("%s: failed to parse the store feed: %w", f.DataID, err)
		}
		if feed.RoundId == 0 {
			return nil, fmt.Errorf("%s: the store has no report yet", f.DataID)
		}
		latest, err := df1.latest()
		if err != nil {
			return nil, err
		}
		if err := checkCutoverValue(feed.LatestAnswer.BigInt(), feed.LatestObservationsTimestamp, df1.Header.Decimals,
			latest.Answer, time.Now(), req.MaxStaleness, req.MaxDeviationBps); err != nil {
			return nil, fmt.Errorf("%s: DF1 feed %s: %w", f.DataID, f.Feed, err)
		}

		ix, err := df1.setWriter(req.DF1StoreProgramID, storeAuthority, signer)
		if err != nil {
			return nil, err
		}
		ixs = append(ixs, ix)
	}
	return ixs, nil
}

// checkCutoverValue checks the store's latest answer, at 18 decimals, is observed within
// maxStaleness of now and, rescaled to the DF1 feed's decimals, is within maxDeviationBps of
// the DF1 feed's latest answer.
func checkCutoverValue(storeAnswer *big.Int, storeObservedAt uint32, df1Decimals uint8, df1Answer *big.Int,
	now time.Time, maxStaleness time.Duration, maxDeviationBps uint64,
) error {
	age := now.Sub(time.Unix(int64(storeObservedAt), 0))
	if age > maxStaleness {
		return fmt.Errorf("the store's latest report is %s old, above %s", age.Truncate(time.Second), maxStaleness)
	}

	rescaled := rescaleAnswer(storeAnswer, df1Decimals)
	// |rescaled - df1| * 10_000 <= maxDeviationBps * |df1|
	diff := new(big.Int).Abs(new(big.Int).Sub(rescaled, df1Answer))
	lhs := diff.Mul(diff, big.NewInt(10_000))
	rhs := new(big.Int).Mul(new(big.Int).SetUint64(maxDeviationBps), new(big.Int).Abs(df1Answer))
	if lhs.Cmp(rhs) > 0 {
		return fmt.Errorf("the store's answer %s (rescaled to %d decimals) deviates from the DF1 answer %s by more than %d bps",
			rescaled, df1Decimals, df1Answer, maxDeviationBps)
	}
	return nil
}

// rescaleAnswer rescales an 18 decimal answer to the given decimals, truncating toward zero like
// the DF1 store's report callback.
func rescaleAnswer(answer *big.Int, decimals uint8) *big.Int {
	if decimals >= storedDecimals {
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals-storedDecimals)), nil)
		return new(big.Int).Mul(answer, scale)
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(storedDecimals-decimals)), nil)
	return new(big.Int).Quo(answer, scale)
}

func verifyDF1Callback(ctx context.Context, client *rpc.Client, store storeAddresses, dataID [16]byte, df1Program, df1Feed solana.PublicKey) error {
	feedAddress, err := store.feed(dataID)
	if err != nil {
		return err
	}
	info, err := getAccount(ctx, client, feedAddress)
	if err != nil {
		return fmt.Errorf("failed to read the store feed: %w", err)
	}
	feed, err := df_store.ParseAccount_FeedData(info.Data.GetBinary())
	if err != nil {
		return fmt.Errorf("failed to parse the store feed: %w", err)
	}
	config := feed.Settings
	if config.CallbackProgram != df1Program {
		return fmt.Errorf("report callback program is %s, expected the DF1 store %s", config.CallbackProgram, df1Program)
	}
	if config.CallbackAccountsLen < 1 || config.CallbackAccounts[0] != df1Feed || config.CallbackWritableMask&1 == 0 {
		return fmt.Errorf("report callback does not write the DF1 feed %s", df1Feed)
	}
	return nil
}

type RollbackDF1Feed struct {
	Feed solana.PublicKey
	// OCR2 state of the feed's aggregator, its store authority PDA ["store", state] becomes the writer
	OCR2State solana.PublicKey
}

type RollbackDF1FeedsRequest struct {
	ChainSel uint64
	// if set, the DF1 feeds' owner is the timelock and set_writer is proposed instead of sent
	MCMS              *cldfproposalutils.TimelockConfig
	DF1StoreProgramID solana.PublicKey
	OCR2ProgramID     solana.PublicKey
	Feeds             []RollbackDF1Feed
}

var _ cldf.ChangeSetV2[*RollbackDF1FeedsRequest] = RollbackDF1Feeds{}

// RollbackDF1Feeds makes OCR2 the writer of DF1 feeds again, undoing CutoverDF1Feeds. Restart the
// feeds' OCR2 jobs once it has landed. Feeds already written by OCR2 are skipped.
type RollbackDF1Feeds struct{}

func (cs RollbackDF1Feeds) VerifyPreconditions(env cldf.Environment, req *RollbackDF1FeedsRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	if req.DF1StoreProgramID.IsZero() || req.OCR2ProgramID.IsZero() {
		return errors.New("DF1StoreProgramID and OCR2ProgramID are required")
	}
	if len(req.Feeds) == 0 {
		return errors.New("feeds cannot be empty")
	}
	seen := map[solana.PublicKey]bool{}
	for _, f := range req.Feeds {
		if f.Feed.IsZero() || f.OCR2State.IsZero() {
			return errors.New("feed and OCR2 state are required")
		}
		if seen[f.Feed] {
			return fmt.Errorf("DF1 feed %s listed twice", f.Feed)
		}
		seen[f.Feed] = true
	}
	if err := verifyMCMS(env, req.ChainSel, req.MCMS); err != nil {
		return err
	}
	_, err := cs.plan(env, req)
	return err
}

func (cs RollbackDF1Feeds) Apply(env cldf.Environment, req *RollbackDF1FeedsRequest) (cldf.ChangesetOutput, error) {
	ixs, err := cs.plan(env, req)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return sendInstructions(env, req.ChainSel, req.MCMS, "roll DF1 feeds back to OCR2", ixs)
}

func (cs RollbackDF1Feeds) plan(env cldf.Environment, req *RollbackDF1FeedsRequest) ([]solana.Instruction, error) {
	ctx := env.GetContext()
	client := env.BlockChains.SolanaChains()[req.ChainSel].Client
	signer, err := authority(env, req.ChainSel, req.MCMS)
	if err != nil {
		return nil, err
	}

	var ixs []solana.Instruction
	for _, f := range req.Feeds {
		ocr2State, err := getAccount(ctx, client, f.OCR2State)
		if err != nil {
			return nil, fmt.Errorf("failed to read OCR2 state %s: %w", f.OCR2State, err)
		}
		if ocr2State.Owner != req.OCR2ProgramID {
			return nil, fmt.Errorf("OCR2 state %s is owned by %s, not the OCR2 program %s", f.OCR2State, ocr2State.Owner, req.OCR2ProgramID)
		}
		ocr2Authority, err := OCR2StoreAuthority(req.OCR2ProgramID, f.OCR2State)
		if err != nil {
			return nil, err
		}

		df1, err := loadDF1Feed(ctx, client, req.DF1StoreProgramID, f.Feed)
		if err != nil {
			return nil, err
		}
		if df1.Header.Writer == ocr2Authority {
			env.Logger.Infow("DF1 feed already written by OCR2, skipping", "feed", f.Feed)
			continue
		}
		if err := df1.verifyOwner(signer); err != nil {
			return nil, err
		}
		ix, err := df1.setWriter(req.DF1StoreProgramID, ocr2Authority, signer)
		if err != nil {
			return nil, err
		}
		ixs = append(ixs, ix)
	}
	return ixs, nil
}

// OCR2StoreAuthority is the OCR2 PDA ["store", state] that writes an aggregator's DF1 feed.
func OCR2StoreAuthority(ocr2Program, ocr2State solana.PublicKey) (solana.PublicKey, error) {
	pda, _, err := solana.FindProgramAddress([][]byte{[]byte("store"), ocr2State.Bytes()}, ocr2Program)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to derive OCR2 store authority: %w", err)
	}
	return pda, nil
}

type df1Feed struct {
	Address solana.PublicKey
	Header  *df1_store.Transmissions
	// Owner whose signature set_writer needs: the feed owner, or the owner of the store account
	// that owns the feed
	EffectiveOwner solana.PublicKey
	data           []byte
}

type df1Round struct {
	Timestamp uint32
	Answer    *big.Int
}

func loadDF1Feed(ctx context.Context, client *rpc.Client, df1Program, address solana.PublicKey) (*df1Feed, error) {
	info, err := getAccount(ctx, client, address)
	if err != nil {
		return nil, fmt.Errorf("failed to read DF1 feed %s: %w", address, err)
	}
	if info.Owner != df1Program {
		return nil, fmt.Errorf("DF1 feed %s is owned by %s, not the DF1 store %s", address, info.Owner, df1Program)
	}
	data := info.Data.GetBinary()
	header, err := df1_store.ParseAccount_Transmissions(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DF1 feed %s: %w", address, err)
	}

	f := &df1Feed{Address: address, Header: header, EffectiveOwner: header.Owner, data: data}
	// A feed owned by a store account is administered by the store's owner
	if owner, err := client.GetAccountInfoWithOpts(ctx, header.Owner, &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}); err == nil &&
		owner.Value != nil && owner.Value.Owner == df1Program {
		if store, err := df1_store.ParseAccount_Store(owner.Value.Data.GetBinary()); err == nil {
			f.EffectiveOwner = store.Owner
		}
	}
	return f, nil
}

func (f *df1Feed) verifyOwner(signer solana.PublicKey) error {
	if f.EffectiveOwner != signer {
		return fmt.Errorf("DF1 feed %s is administered by %s, not by the signer %s", f.Address, f.EffectiveOwner, signer)
	}
	return nil
}

func (f *df1Feed) latest() (df1Round, error) {
	h := f.Header
	if h.LatestRoundId == 0 || h.LiveLength == 0 {
		return df1Round{}, fmt.Errorf("DF1 feed %s has no round", f.Address)
	}
	i := (h.LiveCursor + h.LiveLength - 1) % h.LiveLength
	offset := df1HeaderSize + int(i)*df1TransmissionSize
	if len(f.data) < offset+df1TransmissionSize {
		return df1Round{}, fmt.Errorf("DF1 feed %s is too short for its live ring", f.Address)
	}
	t := f.data[offset : offset+df1TransmissionSize]
	return df1Round{
		Timestamp: binary.LittleEndian.Uint32(t[8:12]),
		Answer:    fromLittleEndianInt128(t[16:32]),
	}, nil
}

func (f *df1Feed) setWriter(df1Program, writer, signer solana.PublicKey) (solana.Instruction, error) {
	return storeInstruction(df1Program, nil)(df1_store.NewSetWriterInstruction(writer, f.Address, f.Header.Owner, signer))
}

func fromLittleEndianInt128(b []byte) *big.Int {
	be := make([]byte, 16)
	for i := range be {
		be[i] = b[15-i]
	}
	v := new(big.Int).SetBytes(be)
	if be[0]&0x80 != 0 {
		v.Sub(v, new(big.Int).Lsh(big.NewInt(1), 128))
	}
	return v
}

func getAccount(ctx context.Context, client *rpc.Client, address solana.PublicKey) (*rpc.Account, error) {
	res, err := client.GetAccountInfoWithOpts(ctx, address, &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Value == nil {
		return nil, fmt.Errorf("account %s not found", address)
	}
	return res.Value, nil
}
