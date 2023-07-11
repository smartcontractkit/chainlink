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
}

type LokiReporter struct {
	host    string
	auth    string
	command string
	now     func() time.Time
}

func (l *LokiReporter) createRequest(flakeyTests map[string][]string) (pushRequest, error) {
	vs := [][]string{}
	now := l.now()
	for pkg, tests := range flakeyTests {
		for _, t := range tests {
			d, err := json.Marshal(flakeyTest{
				Package:    pkg,
				TestName:   t,
				FQTestName: fmt.Sprintf("%s:%s", pkg, t),
			})
			if err != nil {
				return pushRequest{}, err
			}
			vs = append(vs, []string{fmt.Sprintf("%d", now.UnixNano()), string(d)})
		}
	}

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

func (l *LokiReporter) Report(flakeyTests map[string][]string) error {
	if len(flakeyTests) == 0 {
		return nil
	}

	pushReq, err := l.createRequest(flakeyTests)
	if err != nil {
		return err
	}

	return l.makeRequest(pushReq)
}

func NewLokiReporter(host, auth, command string) *LokiReporter {
	return &LokiReporter{host: host, auth: auth, command: command, now: time.Now}
}
