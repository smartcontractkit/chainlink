package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input Duration
		want  string
	}{
		{"zero", *MustNewDuration(0), `"0s"`},
		{"one second", *MustNewDuration(time.Second), `"1s"`},
		{"one minute", *MustNewDuration(time.Minute), `"1m0s"`},
		{"one hour", *MustNewDuration(time.Hour), `"1h0m0s"`},
		{"one hour thirty minutes", *MustNewDuration(time.Hour + 30*time.Minute), `"1h30m0s"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(b))
		})
	}
}

func TestDuration_Scan_Value(t *testing.T) {
	t.Parallel()

	d := MustNewDuration(100)
	require.NotNil(t, d)

	val, err := d.Value()
	require.NoError(t, err)

	dNew := MustNewDuration(0)
	err = dNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, d, dNew)
}

func TestDuration_MarshalJSON_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	d := MustNewDuration(100)
	require.NotNil(t, d)

	json, err := d.MarshalJSON()
	require.NoError(t, err)

	dNew := MustNewDuration(0)
	err = dNew.UnmarshalJSON(json)
	require.NoError(t, err)

	require.Equal(t, d, dNew)
}

func TestDuration_MakeDurationFromString(t *testing.T) {
	t.Parallel()

	d, err := ParseDuration("1s")
	require.NoError(t, err)
	require.Equal(t, 1*time.Second, d.Duration())

	_, err = ParseDuration("xyz")
	require.Error(t, err)
}
