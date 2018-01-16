package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestWebURLUnmarshalJSONError(t *testing.T) {
	t.Parallel()
	j := []byte(`"NotAUrl"`)
	wurl := &models.WebURL{}
	err := json.Unmarshal(j, wurl)
	assert.NotNil(t, err)
}

func TestWebURLUnmarshalJSON(t *testing.T) {
	t.Parallel()
	j := []byte(`"http://www.duckduckgo.com"`)
	wurl := &models.WebURL{}
	err := json.Unmarshal(j, wurl)
	assert.Nil(t, err)
}

func TestTimeDurationFromNow(t *testing.T) {
	future := models.Time{time.Now().Add(time.Second)}
	duration := future.DurationFromNow()
	assert.True(t, 0 < duration)
}
