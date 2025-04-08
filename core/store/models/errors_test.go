package models_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestNewJSONAPIErrors(t *testing.T) {
	t.Parallel()

	res := models.NewJSONAPIErrors()
	require.NotNil(t, res)
	require.NotNil(t, res.Errors)
	require.Empty(t, res.Errors)
}

func TestNewJSONAPIErrorsWith(t *testing.T) {
	t.Parallel()

	res := models.NewJSONAPIErrorsWith("foo")
	require.NotNil(t, res)
	require.NotNil(t, res.Errors)
	require.Len(t, res.Errors, 1)
	require.Equal(t, "foo", res.Errors[0].Detail)
}

func TestJSONAPIErrors_Error(t *testing.T) {
	t.Parallel()

	res := models.NewJSONAPIErrorsWith("foo")
	require.NotNil(t, res)
	require.Equal(t, "foo", res.Error())

	res.Add("bar")
	require.Equal(t, "foo,bar", res.Error())
}

func TestJSONAPIErrors_CoerceEmptyToNil(t *testing.T) {
	t.Parallel()

	res := models.NewJSONAPIErrors()
	require.NotNil(t, res)

	err := res.CoerceEmptyToNil()
	require.NoError(t, err)

	res = models.NewJSONAPIErrorsWith("foo")
	require.NotNil(t, res)

	err = res.CoerceEmptyToNil()
	require.Equal(t, res, err)
}

func TestJSONAPIErrors_Merge(t *testing.T) {
	t.Parallel()

	res1 := models.NewJSONAPIErrorsWith("foo")
	require.NotNil(t, res1)

	res2 := models.NewJSONAPIErrorsWith("bar")
	require.NotNil(t, res2)

	res1.Merge(res2)
	require.Equal(t, "foo,bar", res1.Error())

	res1.Merge(errors.New("zet"))
	require.Equal(t, "foo,bar,zet", res1.Error())
}
