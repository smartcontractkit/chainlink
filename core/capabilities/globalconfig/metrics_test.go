package globalconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig/globalconfigtest"
)

func TestMetrics_DesignNamesAndLabels(t *testing.T) {
	t.Parallel()

	m, reader := globalconfigtest.NewMetrics(t)
	m.RecordAppliedVersion(t.Context(), "cre", "prod", 201)
	m.RecordValidationError(t.Context(), "cre", "prod")
	m.RecordValidationError(t.Context(), "cre", "prod")
	m.RecordApplyError(t.Context(), "cre", "staging")

	got := globalconfigtest.Collect(t, reader)
	assert.Equal(t, int64(201), got["platform_cap_config_applied_version"][globalconfigtest.Series("cre", "prod")])
	assert.Equal(t, int64(2), got["platform_cap_config_validation_errors_total"][globalconfigtest.Series("cre", "prod")])
	assert.Equal(t, int64(1), got["platform_cap_config_apply_errors_total"][globalconfigtest.Series("cre", "staging")])

	var nilMetrics *globalconfig.Metrics
	nilMetrics.RecordApplyError(t.Context(), "", "") // no-op, must not panic
}

func TestPayloadLabels(t *testing.T) {
	t.Parallel()
	d, e := globalconfig.PayloadLabels(`{"domain":"cre","env":"prod","version":"nope","bogus":1}`)
	assert.Equal(t, "cre", d)
	assert.Equal(t, "prod", e)
	d, e = globalconfig.PayloadLabels(`not json`)
	assert.Empty(t, d)
	assert.Empty(t, e)
}
