package validationtest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var ValidationService = staticValidationService{}

var GoodPluginConfig = map[string]interface{}{
	"isGoodConfig":  true,
	"someFieldName": "someFieldValue",
}

var _ core.ValidationService = (*staticValidationService)(nil)

type staticValidationService struct {
	services.Service
}

func (t staticValidationService) Evaluate(ctx context.Context, other core.ValidationService) error {
	return other.ValidateConfig(ctx, GoodPluginConfig)
}

func (t staticValidationService) ValidateConfig(ctx context.Context, config map[string]interface{}) error {
	if !reflect.DeepEqual(GoodPluginConfig, config) {
		return fmt.Errorf("expected %+v but got %+v", GoodPluginConfig, config)
	}
	return nil
}
