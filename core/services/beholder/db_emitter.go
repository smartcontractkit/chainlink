package beholder

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type DbEmitter struct {
	ds sqlutil.DataSource
}

var _ beholder.Emitter = (*DbEmitter)(nil)

func NewDbEmitter(ds sqlutil.DataSource) *DbEmitter {
	return &DbEmitter{
		ds: ds,
	}
}

func (emitter *DbEmitter) Emit(ctx context.Context, body []byte, attrKVs ...any) error {
	_, err := emitter.BatchEmit(ctx, []beholder.Message{
		beholder.NewMessage(body, attrKVs...),
	})
	return err
}

func (emitter *DbEmitter) BatchEmit(ctx context.Context, messages []beholder.Message, options ...beholder.BatchEmitOption) ([]*chipingress.PublishResult, error) {
	return nil, sqlutil.TransactDataSource(ctx, emitter.ds, nil, func(ds sqlutil.DataSource) error {
		for _, msg := range messages {
			attrs, err := json.Marshal(msg.Attrs)
			if err != nil {
				return fmt.Errorf("marshaling attributes: %w", err)
			}
			_, err = ds.ExecContext(ctx, "INSERT INTO queue (payload, attributes) VALUES ($1, $2)", msg.Body, attrs)
			if err != nil {
				return fmt.Errorf("inserting message: %w", err)
			}
		}
		return nil
	})
}

func (emitter *DbEmitter) Close() error {
	return nil
}
