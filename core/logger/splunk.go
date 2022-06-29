package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger/audit"

	"go.uber.org/zap/zapcore"
)

type splunkLogger struct {
	logger            Logger
	logBuffer         chan splunkLogItem
	serviceDoneSignal chan struct{}
	splunkToken       string
	splunkURL         string
	environmentName   string
	hostname          string
	localIP           string
}

type splunkLogItem struct {
	eventID audit.EventID
	data    map[string]interface{}
}

func newSplunkLogger(logger Logger, splunkToken string, splunkURL string, hostname string, environment string) Logger {
	// This logger implements a single goroutine buffer/queue system to enable fire and forget
	// dispatch of audit log events within web controllers. The async http post to collector has
	// a timeout of 30 seconds, so internal API responses won't hang indefinitely.
	logBuff := make(chan splunkLogItem, 10)
	doneSignal := make(chan struct{})

	sLogger := splunkLogger{
		logger:            logger.Helper(1),
		logBuffer:         logBuff,
		serviceDoneSignal: doneSignal,
		splunkToken:       splunkToken,
		splunkURL:         splunkURL,
		environmentName:   environment,
		hostname:          hostname,
		localIP:           getLocalIP(),
	}

	// Start single async forwarder thread
	go sLogger.StartLogForwarder()

	// Initialize and return Splunk logger struct with required state for HEC calls
	return &sLogger
}

// StartLogForwarder is a goroutine with buffer limit to process and forward log events async
func (l *splunkLogger) StartLogForwarder() {
	for {
		select {
		case event := <-l.logBuffer:
			l.postLogToSplunk(event.eventID, event.data)
		case <-l.serviceDoneSignal:
			return
		}
	}
}

func (l *splunkLogger) Audit(eventID audit.EventID, data map[string]interface{}) {
	// Queue event data in logBuffer for forwarder to pick up and send
	event := splunkLogItem{
		eventID: eventID,
		data:    data,
	}
	l.logBuffer <- event
	l.logger.Audit(eventID, data)
}

// getLocalIP returns the first non-loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// filter and return address types for first non loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (l *splunkLogger) postLogToSplunk(eventID audit.EventID, data map[string]interface{}) {
	// Splunk JSON data
	splunkLog := map[string]interface{}{
		"eventID":  eventID,
		"hostname": l.hostname,
		"localIP":  l.localIP,
		"env":      l.environmentName,
	}
	if len(data) != 0 {
		splunkLog["data"] = data
	}

	// Wrap serialized audit log map into JSON object `event` for API call
	serializedArgs, _ := json.Marshal(splunkLog)
	splunkLog = map[string]interface{}{"event": string(serializedArgs)}
	serializedSplunkLog, _ := json.Marshal(splunkLog)

	// Send up to HEC log collector
	httpClient := &http.Client{Timeout: time.Second * 30}
	req, _ := http.NewRequest("POST", l.splunkURL, bytes.NewReader(serializedSplunkLog))
	req.Header.Add("Authorization", "Splunk "+l.splunkToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		l.logger.Errorw("Failed to send audit log to Splunk", "err", err, "splunkLog", splunkLog)
	}
	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.logger.Errorw("Error reading errored Splunk webhook response body", "err", err, "splunkLog", splunkLog)
		}
		l.logger.Errorw("Error sending log to Splunk", "statusCode", resp.StatusCode, "bodyString", string(bodyBytes))
	}
}

func (l *splunkLogger) With(args ...interface{}) Logger {
	return &splunkLogger{
		logger:          l.logger.With(args...),
		splunkToken:     l.splunkToken,
		splunkURL:       l.splunkURL,
		environmentName: l.environmentName,
		hostname:        l.hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) Named(name string) Logger {
	return &splunkLogger{
		logger:          l.logger.Named(name),
		splunkToken:     l.splunkToken,
		splunkURL:       l.splunkURL,
		environmentName: l.environmentName,
		hostname:        l.hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) SetLogLevel(level zapcore.Level) {
	l.logger.SetLogLevel(level)
}

func (l *splunkLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *splunkLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *splunkLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *splunkLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *splunkLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *splunkLogger) Critical(args ...interface{}) {
	l.logger.Critical(args...)
}

func (l *splunkLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *splunkLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *splunkLogger) Tracef(format string, values ...interface{}) {
	l.logger.Tracef(format, values...)
}

func (l *splunkLogger) Debugf(format string, values ...interface{}) {
	l.logger.Debugf(format, values...)
}

func (l *splunkLogger) Infof(format string, values ...interface{}) {
	l.logger.Infof(format, values...)
}

func (l *splunkLogger) Warnf(format string, values ...interface{}) {
	l.logger.Warnf(format, values...)
}

func (l *splunkLogger) Errorf(format string, values ...interface{}) {
	l.logger.Errorf(format, values...)
}

func (l *splunkLogger) Criticalf(format string, values ...interface{}) {
	l.logger.Criticalf(format, values...)
}

func (l *splunkLogger) Panicf(format string, values ...interface{}) {
	l.logger.Panicf(format, values...)
}

func (l *splunkLogger) Fatalf(format string, values ...interface{}) {
	l.logger.Fatalf(format, values...)
}

func (l *splunkLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.logger.Tracew(msg, keysAndValues...)
}

func (l *splunkLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *splunkLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *splunkLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *splunkLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *splunkLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	l.logger.Criticalw(msg, keysAndValues...)
}

func (l *splunkLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *splunkLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *splunkLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.logger.Errorw(msg, "err", err)
	}
}

func (l *splunkLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		l.logger.Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (l *splunkLogger) Sync() error {
	return l.logger.Sync()
}

func (l *splunkLogger) Helper(add int) Logger {
	return &splunkLogger{
		logger:          l.logger.Helper(add),
		splunkToken:     l.splunkToken,
		splunkURL:       l.splunkURL,
		environmentName: l.environmentName,
		hostname:        l.hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) Recover(panicErr interface{}) {
	l.logger.Recover(panicErr)
}
