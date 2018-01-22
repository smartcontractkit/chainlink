package models_test

import (
	"encoding/json"
	"net/url"
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

func TestWebURLMarshalJSON(t *testing.T) {
	t.Parallel()

	str := "http://www.duckduckgo.com"
	parsed, err := url.ParseRequestURI(str)
	assert.Nil(t, err)
	wurl := &models.WebURL{parsed}
	b, err := json.Marshal(wurl)
	assert.Nil(t, err)
	assert.Equal(t, `"`+str+`"`, string(b))
}

func TestTimeDurationFromNow(t *testing.T) {
	future := models.Time{time.Now().Add(time.Second)}
	duration := future.DurationFromNow()
	assert.True(t, 0 < duration)
}
