package forms_test

import (
	"testing"

	"github.com/asdine/storm"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/forms"
	"github.com/stretchr/testify/assert"
)

func TestForms_NewUpdateBridgeType(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()

	bt := cltest.NewBridgeType("bridgea")
	assert.Nil(t, s.Save(&bt))

	_, err := forms.NewUpdateBridgeType(s, "idontexist")
	assert.Equal(t, err, storm.ErrNotFound)

	_, err = forms.NewUpdateBridgeType(s, "bridgea")
	assert.NoError(t, err)
}

func TestForms_UpdateBridgeType_Save(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()

	bt := cltest.NewBridgeType("bridgea", "http://bridge")
	assert.Nil(t, s.Save(&bt))

	form, err := forms.NewUpdateBridgeType(s, "bridgea")
	assert.NoError(t, err)
	assert.NoError(t, form.Save())

	ubt, err := s.FindBridge("bridgea")
	assert.Equal(t, cltest.WebURL("http://bridge"), ubt.URL)
	assert.Equal(t, uint64(0), ubt.Confirmations)
	assert.Equal(t, *assets.NewLink(0), ubt.MinimumContractPayment)

	form.URL = cltest.WebURL("http://updatedbridge")
	form.Confirmations = uint64(10)
	form.MinimumContractPayment = *assets.NewLink(100)
	assert.NoError(t, form.Save())

	ubt, err = s.FindBridge("bridgea")
	assert.NoError(t, err)
	assert.Equal(t, cltest.WebURL("http://updatedbridge"), ubt.URL)
	assert.Equal(t, uint64(10), ubt.Confirmations)
	assert.Equal(t, *assets.NewLink(100), ubt.MinimumContractPayment)
}
