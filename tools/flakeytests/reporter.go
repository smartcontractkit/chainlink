package flakeytests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type pushRequest struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type flakeyTest struct {
	Package    string `json:"package"`
	TestName   string `json:"test_name"`
	FQTestName string `json:"fq_test_name"`
	Context
}

type numFlakes struct {
	NumFlakes int `json:"num_flakes"`
	Context
}

type Context struct {
	CommitSHA      string `json:"commit_sha,omitempty"`
	PullRequestURL string `json:"pull_request_url,omitempty"`
}

type LokiReporter struct {
	host    string
	auth    string
	command string
	now     func() time.Time
	ctx     Context
}

func (l *LokiReporter) createRequest(flakeyTests map[string]map[string]struct{}) (pushRequest, error) {
	vs := [][]string{}
	now := l.now()
	nows := fmt.Sprintf("%d", now.UnixNano())
	for pkg, tests := range flakeyTests {
		for t := range tests {
			d, err := json.Marshal(flakeyTest{
				Package:    pkg,
				TestName:   t,
				FQTestName: fmt.Sprintf("%s:%s", pkg, t),
				Context:    l.ctx,
			})
			if err != nil {
				return pushRequest{}, err
			}
			vs = append(vs, []string{nows, string(d)})
		}
	}

	// Flakes are store in a map[string][]string, so to count them, we can't just do len(flakeyTests),
	// as that will get us the number of flakey packages, not the number of flakes tests.
	// However, we do emit one log line per flakey test above, so use that to count our flakes.
	f, err := json.Marshal(numFlakes{NumFlakes: len(vs), Context: l.ctx})
	if err != nil {
		return pushRequest{}, nil
	}

	vs = append(vs, []string{nows, string(f)})

	pr := pushRequest{
		Streams: []stream{
			{
				Stream: map[string]string{
					"app":     "flakey-test-reporter",
					"command": l.command,
				},
				Values: vs,
			},
		},
	}
	return pr, nil
}

func (l *LokiReporter) makeRequest(pushReq pushRequest) error {
	body, err := json.Marshal(pushReq)
	if err != nil {
		return err
	}

	u := url.URL{Scheme: "https", Host: l.host, Path: "loki/api/v1/push"}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(l.auth))),
	)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusNoContent {
		b, berr := io.ReadAll(resp.Body)
		if berr != nil {
			return fmt.Errorf("error decoding body for failed push request: %w", berr)
		}
		return fmt.Errorf("push request failed: status=%d, body=%s", resp.StatusCode, b)
	}
	return err
}

func (l *LokiReporter) Report(flakeyTests map[string]map[string]struct{}) error {
	pushReq, err := l.createRequest(flakeyTests)
	if err != nil {
		return err
	}

	return l.makeRequest(pushReq)
}

func NewLokiReporter(host, auth, command string, ctx Context) *LokiReporter {
	return &LokiReporter{host: host, auth: auth, command: command, now: time.Now, ctx: ctx}
}
