package solana

import (
	"context"
	"errors"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

type TransmissionsCache struct {
	*client.Cache[Answer]
}

func NewTransmissionsCache(transmissionsID solana.PublicKey, chainID string, cfg config.Config, reader client.Reader, lggr logger.Logger) *TransmissionsCache {
	name := "ocr2_median_transmissions"
	getter := func(ctx context.Context) (Answer, uint64, error) {
		return GetLatestTransmission(ctx, reader, transmissionsID, cfg.Commitment())
	}
	return &TransmissionsCache{client.NewCache(name, transmissionsID, chainID, cfg, getter, logger.With(lggr, "cache", name))}
}

func GetLatestTransmission(ctx context.Context, reader client.AccountReader, account solana.PublicKey, commitment rpc.CommitmentType) (Answer, uint64, error) {
	// query for transmission header
	headerStart := AccountDiscriminatorLen // skip account discriminator
	headerLen := TransmissionsHeaderLen
	res, err := reader.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: commitment,
		DataSlice: &rpc.DataSlice{
			Offset: &headerStart,
			Length: &headerLen,
		},
	})
	if err != nil {
		return Answer{}, 0, fmt.Errorf("error on rpc.GetAccountInfo [cursor]: %w", err)
	}

	// check for nil pointers
	if res == nil || res.Value == nil || res.Value.Data == nil {
		return Answer{}, 0, errors.New("nil pointer returned in GetLatestTransmission.GetAccountInfoWithOpts.Header")
	}

	// parse header
	var header TransmissionsHeader
	if err = bin.NewBinDecoder(res.Value.Data.GetBinary()).Decode(&header); err != nil {
		return Answer{}, 0, fmt.Errorf("failed to decode transmission account header: %w", err)
	}

	if header.Version != 2 {
		return Answer{}, 0, fmt.Errorf("can't parse feed version %v: %w", header.Version, err)
	}

	cursor := header.LiveCursor
	liveLength := header.LiveLength

	if cursor == 0 { // handle array wrap
		cursor = liveLength
	}
	cursor-- // cursor indicates index for new answer, latest answer is in previous index

	// setup transmissionLen
	transmissionLen := TransmissionLen

	transmissionOffset := AccountDiscriminatorLen + TransmissionsHeaderMaxSize + (uint64(cursor) * transmissionLen)

	res, err = reader.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: commitment,
		DataSlice: &rpc.DataSlice{
			Offset: &transmissionOffset,
			Length: &transmissionLen,
		},
	})
	if err != nil {
		return Answer{}, 0, fmt.Errorf("error on rpc.GetAccountInfo [transmission]: %w", err)
	}
	// check for nil pointers
	if res == nil || res.Value == nil || res.Value.Data == nil {
		return Answer{}, 0, errors.New("nil pointer returned in GetLatestTransmission.GetAccountInfoWithOpts.Transmission")
	}

	// parse tranmission
	var t Transmission
	if err := bin.NewBinDecoder(res.Value.Data.GetBinary()).Decode(&t); err != nil {
		return Answer{}, 0, fmt.Errorf("failed to decode transmission: %w", err)
	}

	return Answer{
		Data:      t.Answer.BigInt(),
		Timestamp: t.Timestamp,
	}, res.RPCContext.Context.Slot, nil
}
