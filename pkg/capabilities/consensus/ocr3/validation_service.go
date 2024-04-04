package ocr3

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ commontypes.ValidationService = (*validationService)(nil)

type validationService struct {
	lggr logger.Logger
	services.StateMachine
}

func (v *validationService) ValidateConfig(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (v *validationService) Start(ctx context.Context) error {
	return v.StartOnce("OCR3ReportingPluginValidation", func() error {
		return nil
	})
}

func (v *validationService) Close() error {
	return v.StopOnce("OCR3ReportingPluginValidation", func() error {
		return nil
	})
}

func (v *validationService) Name() string { return v.lggr.Name() }

func (v *validationService) HealthReport() map[string]error {
	return map[string]error{v.Name(): v.Healthy()}
}
