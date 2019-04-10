package forms_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/forms"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormsNewUpdateBridgeType(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()

	_, bt := cltest.NewBridgeType("bridgea")
	require.NoError(t, s.CreateBridgeType(bt))

	_, err := forms.NewUpdateBridgeType(s, "idontexist")
	assert.Equal(t, err, orm.ErrorNotFound)

	_, err = forms.NewUpdateBridgeType(s, "bridgea")
	assert.NoError(t, err)
}

func TestFormsUpdateBridgeType_Save(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()

	_, bt := cltest.NewBridgeType("bridgea", "http://bridge")
	require.NoError(t, s.CreateBridgeType(bt))

	form, err := forms.NewUpdateBridgeType(s, "bridgea")
	assert.NoError(t, err)
	assert.NoError(t, form.Save())

	ubt, err := s.FindBridge("bridgea")
	assert.Equal(t, cltest.WebURL("http://bridge"), ubt.URL)
	assert.Equal(t, uint64(0), ubt.Confirmations)
	assert.Nil(t, ubt.MinimumContractPayment)

	form.URL = cltest.WebURL("http://updatedbridge")
	form.Confirmations = uint64(10)
	form.MinimumContractPayment = assets.NewLink(100)
	assert.NoError(t, form.Save())

	ubt, err = s.FindBridge("bridgea")
	assert.NoError(t, err)
	assert.Equal(t, cltest.WebURL("http://updatedbridge"), ubt.URL)
	assert.Equal(t, uint64(10), ubt.Confirmations)
	assert.Equal(t, assets.NewLink(100), ubt.MinimumContractPayment)
}
