package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client interface {
	bind.ContractBackend
	bind.DeployBackend
}

func MaybeAddStubCommitMethod(b Client) CommittingClient {
	client, ok := b.(CommittingClient)
	if !ok {
		client = noncommittalClient{b}
	}
	return client
}

type committer interface{ Commit() }

type CommittingClient interface {
	Client
	committer
}

type noncommittalClient struct{ Client }

var _ CommittingClient = noncommittalClient{}

func (n noncommittalClient) Commit() {}

func CheckStatus(
	ctx context.Context,
	tx *types.Transaction,
	client Client,
) (err error) {
	timeout := time.After(5 * time.Second)
	var buf bytes.Buffer
	var receipt *types.Receipt
	errMsg := "could not get byte encoding of tx with hash 0x%x while " +
		"reporting on %s while checking its status, due to error \"%w\""
	for receipt == nil {
		receipt, err = client.TransactionReceipt(ctx, tx.Hash())
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			berr := tx.EncodeRLP(&buf)
			if berr != nil {
				return WrapErrorf(berr, errMsg, "error", tx.Hash(), err)
			}
			b := buf.Bytes()
			return WrapErrorf(err, "could not retrieve receipt for tx 0x%x", b)
		}
		select {
		case <-time.After(10 * time.Millisecond):
			continue
		case <-timeout:
			berr := tx.EncodeRLP(&buf)
			if berr != nil {
				return WrapErrorf(berr, errMsg, "timeout", tx.Hash(), err)
			}
			return fmt.Errorf("timeout while checking status on tx 0x%x", buf.Bytes())
		case <-ctx.Done():
			return fmt.Errorf("context expired while checking tx status")
		}
	}
	if receipt.Status != 1 {
		berr := tx.EncodeRLP(&buf)
		if berr != nil {
			return WrapErrorf(berr, errMsg, "error", tx.Hash(), err)
		}
		h := tx.Hash()
		b := buf.Bytes()
		return fmt.Errorf("transaction failed: hash: 0x%x; contents: 0x%x", h, b)
	}
	return nil
}
