package logger

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type clientMock struct {
	responseCode int
	req          string
}

func (c *clientMock) Do(req *http.Request) (*http.Response, error) {
	content, _ := ioutil.ReadAll(req.Body)
	c.req = string(content)
	resp := &http.Response{}
	resp.StatusCode = c.responseCode
	resp.Body = io.NopCloser(strings.NewReader("foo"))
	return resp, nil
}

func Test_splunkSink(t *testing.T) {
	client := &clientMock{
		responseCode: 200,
	}

	s := &splunkSink{
		HTTPClient:    client,
		URL:           "https://example.com/collector/event",
		Hostname:      "foo",
		Token:         "123-456",
		Source:        "mysource",
		SourceType:    "mysourcetype",
		Index:         "myindex",
		SkipTLSVerify: true,
		nowFn: func() int64 {
			return 0
		},
	}
	_, err := s.Write([]byte("{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"foo\"}"))
	assert.NoError(t, err)
	_, err = s.Write([]byte("{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"bar\"}"))
	assert.NoError(t, err)
	_, err = s.Write([]byte("{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"foobar\"}"))
	assert.NoError(t, err)
	err = s.Sync()
	assert.NoError(t, err)

	assert.Equal(t, "{\"time\":0,\"host\":\"foo\",\"source\":\"mysource\",\"sourcetype\":\"mysourcetype\",\"index\":\"myindex\",\"event\":{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"foo\"}}{\"time\":0,\"host\":\"foo\",\"source\":\"mysource\",\"sourcetype\":\"mysourcetype\",\"index\":\"myindex\",\"event\":{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"bar\"}}{\"time\":0,\"host\":\"foo\",\"source\":\"mysource\",\"sourcetype\":\"mysourcetype\",\"index\":\"myindex\",\"event\":{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"foobar\"}}", client.req)
	_, err = s.Write([]byte("{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"foo\"}"))
	assert.NoError(t, err)
	err = s.Sync()
	assert.NoError(t, err)
	assert.Equal(t, "{\"time\":0,\"host\":\"foo\",\"source\":\"mysource\",\"sourcetype\":\"mysourcetype\",\"index\":\"myindex\",\"event\":{\"level\":\"error\",\"ts\":\"2021-11-04T16:05:14.641-0700\",\"caller\":\"config/general_config.go:302\",\"msg\":\"foo\"}}", client.req)
}
