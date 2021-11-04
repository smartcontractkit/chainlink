package logger

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type event struct {
	Time       int64           `json:"time"`                 // epoch time in nano-seconds
	Host       string          `json:"host"`                 // hostname
	Source     string          `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string          `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string          `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      json.RawMessage `json:"event"`                // zap json logger
}

type splunkSink struct {
	HTTPClient    httpClient
	nowFn         func() int64
	URL           string // URL to Splunk HTTP Event Collector
	Hostname      string // Hostname is the human-readable machine name sending logs
	Token         string // Splunk HTTP Event Collector token
	Source        string // Splunk source field value, description of the source of the event
	SourceType    string // Splunk source type, optional name of a sourcetype field value
	Index         string // Splunk Index, optional name of the Splunk index to store the event in
	SkipTLSVerify bool   // Skip verifying the certificate of the HTTP Event Collector
	events        [][]byte
	sync.Mutex
}

func envVarOrDefault(key string, defaultValue string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}

func initialize(s *splunkSink) *splunkSink {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: s.SkipTLSVerify}}
	s.HTTPClient = &http.Client{Timeout: time.Second * 20, Transport: tr}

	if s.Hostname == "" {
		s.Hostname, _ = os.Hostname()
	}
	if s.nowFn == nil {
		s.nowFn = func() int64 {
			return time.Now().UnixNano()
		}
	}

	return s
}

func newSplunkSink() (*splunkSink, error) {
	s := &splunkSink{}
	s.URL = os.Getenv("SPLUNK_URL")
	s.Source = envVarOrDefault("SPLUNK_SOURCE", "chainlink")
	s.SourceType = envVarOrDefault("SPLUNK_SOURCETYPE", "chainlink")
	s.Index = os.Getenv("SPLUNK_INDEX")
	s.Hostname = os.Getenv("SPLUNK_HOSTNAME")
	s.Token = os.Getenv("SPLUNK_TOKEN")
	s.SkipTLSVerify = "true" == os.Getenv("SPLUNK_SKIP_TLS_VERIFY")

	return initialize(s), nil
}

func (p *splunkSink) Close() error {
	return nil
}

// Write implement zap.Sink func Write
func (p *splunkSink) Write(b []byte) (n int, err error) {
	e := &event{
		Time:       p.nowFn(),
		Host:       p.Hostname,
		Source:     p.Source,
		SourceType: p.SourceType,
		Index:      p.Index,
		Event:      json.RawMessage(b),
	}
	encoded, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}
	p.Lock()
	defer p.Unlock()
	p.events = append(p.events, encoded)
	return len(b), nil
}

// Sync implement zap.Sink func Sync
func (p *splunkSink) Sync() error {
	eventsLen := len(p.events)
	if eventsLen > 0 {
		p.Lock()
		eventsToLog := p.events[0:eventsLen]
		p.events = p.events[eventsLen:]
		p.Unlock()
		return p.logEvents(eventsToLog)

	}
	return nil
}

func (p *splunkSink) logEvents(events [][]byte) error {
	buf := new(bytes.Buffer)
	for _, e := range events {
		buf.Write(e)
	}

	return p.doRequest(buf)
}

func (p *splunkSink) doRequest(b *bytes.Buffer) error {
	url := p.URL
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Splunk "+p.Token)

	res, err := p.HTTPClient.Do(req)
	if err != nil && err != io.EOF {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		responseBody := buf.String()
		err = errors.New(responseBody)
	} else {
		_, _ = io.Copy(ioutil.Discard, res.Body)
	}
	return err
}
