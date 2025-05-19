package ocrcommon

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type ocrLoggerService struct {
	stopCh services.StopChan
	ocrtypes.Logger
}

func NewOCRWrapper(l logger.Logger, trace bool, saveError func(context.Context, string)) *ocrLoggerService {
	stopCh := make(services.StopChan)
	return &ocrLoggerService{
		stopCh: stopCh,
		Logger: logger.NewOCRWrapper(l, trace, func(s string) {
			ctx, cancel := stopCh.NewCtx()
			defer cancel()
			saveError(ctx, s)
		}),
	}
}

func (*ocrLoggerService) Start(context.Context) error { return nil }
func (s *ocrLoggerService) Close() error              { close(s.stopCh); return nil }
