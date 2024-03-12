package flakeytests

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	messageType_flakeyTest   = "flakey_test"
	messageType_runReport    = "run_report"
	messageType_packagePanic = "package_panic"
)

type pushRequest struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type BaseMessage struct {
	MessageType string `json:"message_type"`
	Context
}

type flakeyTest struct {
	BaseMessage
	Package    string `json:"package"`
	TestName   string `json:"test_name"`
	FQTestName string `json:"fq_test_name"`
}

type packagePanic struct {
	BaseMessage
	Package string `json:"package"`
}

type runReport struct {
	BaseMessage
	NumPackagePanics int `json:"num_package_panics"`
	NumFlakes        int `json:"num_flakes"`
	NumCombined      int `json:"num_combined"`
}

type Context struct {
	CommitSHA      string `json:"commit_sha"`
	PullRequestURL string `json:"pull_request_url,omitempty"`
	Repository     string `json:"repository"`
	Type           string `json:"event_type"`
	RunURL         string `json:"run_url,omitempty"`
}

type LokiReporter struct {
	host    string
	auth    string
	orgId   string
	command string
	now     func() time.Time
	ctx     Context
}

func (l *LokiReporter) createRequest(report *Report) (pushRequest, error) {
	vs := [][]string{}
	now := l.now()
	nows := fmt.Sprintf("%d", now.UnixNano())

	for pkg, tests := range report.tests {
		for t := range tests {
			d, err := json.Marshal(flakeyTest{
				BaseMessage: BaseMessage{
					MessageType: messageType_flakeyTest,
					Context:     l.ctx,
				},
				Package:    pkg,
				TestName:   t,
				FQTestName: fmt.Sprintf("%s:%s", pkg, t),
			})
			if err != nil {
				return pushRequest{}, err
			}
			vs = append(vs, []string{nows, string(d)})
		}
	}

	// Flakes are stored in a map[string][]string, so to count them, we can't just do len(flakeyTests),
	// as that will get us the number of flakey packages, not the number of flakes tests.
	// However, we do emit one log line per flakey test above, so use that to count our flakes.
	numFlakes := len(vs)

	for pkg := range report.packagePanics {
		d, err := json.Marshal(packagePanic{
			BaseMessage: BaseMessage{
				MessageType: messageType_packagePanic,
				Context:     l.ctx,
			},
			Package: pkg,
		})
		if err != nil {
			return pushRequest{}, err
		}

		vs = append(vs, []string{nows, string(d)})
	}

	f, err := json.Marshal(runReport{
		BaseMessage: BaseMessage{
			MessageType: messageType_runReport,
			Context:     l.ctx,
		},
		NumFlakes:        numFlakes,
		NumPackagePanics: len(report.packagePanics),
		NumCombined:      numFlakes + len(report.packagePanics),
	})
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

func (l *LokiReporter) makeRequest(ctx context.Context, pushReq pushRequest) error {
	body, err := json.Marshal(pushReq)
	if err != nil {
		return err
	}

	u := url.URL{Scheme: "https", Host: l.host, Path: "loki/api/v1/push"}
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(l.auth))),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Scope-OrgID", l.orgId)
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

func (l *LokiReporter) Report(ctx context.Context, report *Report) error {
	pushReq, err := l.createRequest(report)
	if err != nil {
		return err
	}

	return l.makeRequest(ctx, pushReq)
}

func NewLokiReporter(host, auth, orgId, command string, ctx Context) *LokiReporter {
	return &LokiReporter{host: host, auth: auth, orgId: orgId, command: command, now: time.Now, ctx: ctx}
}
