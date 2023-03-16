package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
)

const bufferCapacity = 2048
const webRequestTimeout = 10

type Data = map[string]any

type AuditLogger interface {
	services.ServiceCtx

	Audit(eventID EventID, data Data)
}

type AuditLoggerConfig struct {
	Enabled        *bool
	ForwardToUrl   *models.URL
	JsonWrapperKey *string
	Headers        *[]ServiceHeader
}

func (p *AuditLoggerConfig) SetFrom(f *AuditLoggerConfig) {
	if v := f.Enabled; v != nil {
		p.Enabled = v
	}
	if v := f.ForwardToUrl; v != nil {
		p.ForwardToUrl = v
	}
	if v := f.JsonWrapperKey; v != nil {
		p.JsonWrapperKey = v
	}
	if v := f.Headers; v != nil {
		p.Headers = v
	}

}

// ServiceHeader is an HTTP header to include in POST to log service.
type ServiceHeader struct {
	Header string
	Value  string
}

func (h *ServiceHeader) UnmarshalText(input []byte) error {
	parts := strings.SplitN(string(input), ":", 2)
	h.Header = parts[0]
	if len(parts) > 1 {
		h.Value = strings.TrimSpace(parts[1])
	}
	return h.validate()
}

func (h *ServiceHeader) MarshalText() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s: %s", h.Header, h.Value)
	return b.Bytes(), nil
}

// We act slightly more strictly than the HTTP specifications
// technically allow instead following the guidelines of
// cloudflare transforms.
// https://developers.cloudflare.com/rules/transform/request-header-modification/reference/header-format
var (
	headerNameRegex  = regexp.MustCompile(`^[A-Za-z\-]+$`)
	headerValueRegex = regexp.MustCompile("^[A-Za-z_ :;.,\\/\"'?!(){}[\\]@<>=\\-+*#$&`|~^%]+$")
)

func (h ServiceHeader) validate() (err error) {
	if !headerNameRegex.MatchString(h.Header) {
		err = multierr.Append(err, errors.Errorf("invalid header name: %s", h.Header))
	}

	if !headerValueRegex.MatchString(h.Value) {
		err = multierr.Append(err, errors.Errorf("invalid header value: %s", h.Value))
	}
	return
}

type ServiceHeaders []ServiceHeader

func (sh *ServiceHeaders) UnmarshalText(input []byte) error {
	if sh == nil {
		return errors.New("Cannot unmarshal to a nil receiver")
	}

	headers := string(input)

	var parsedHeaders []ServiceHeader
	if headers != "" {
		headerLines := strings.Split(headers, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return errors.Errorf("invalid headers provided for the audit logger. Value, single pair split on || required, got: %s", keyValue)
			}
			h := ServiceHeader{
				Header: keyValue[0],
				Value:  keyValue[1],
			}

			if err := h.validate(); err != nil {
				return err
			}
			parsedHeaders = append(parsedHeaders, h)
		}
	}

	*sh = parsedHeaders
	return nil
}

func (sh *ServiceHeaders) MarshalText() ([]byte, error) {
	if sh == nil {
		return nil, errors.New("Cannot marshal to a nil receiver")
	}

	sb := strings.Builder{}
	for _, header := range *sh {
		sb.WriteString(header.Header)
		sb.WriteString("||")
		sb.WriteString(header.Value)
		sb.WriteString("\\")
	}

	serialized := sb.String()

	if len(serialized) > 0 {
		serialized = serialized[:len(serialized)-1]
	}

	return []byte(serialized), nil
}

var AuditLoggerHeaders = envvar.New("AuditLoggerHeaders", func(s string) (ServiceHeaders, error) {
	sh := make(ServiceHeaders, 0)
	err := sh.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return sh, nil
})

type Config interface {
	AuditLoggerEnabled() bool
	AuditLoggerForwardToUrl() (models.URL, error)
	AuditLoggerEnvironment() string
	AuditLoggerJsonWrapperKey() string
	AuditLoggerHeaders() (ServiceHeaders, error)
}

type HTTPAuditLoggerInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuditLoggerService struct {
	logger          logger.Logger            // The standard logger configured in the node
	enabled         bool                     // Whether the audit logger is enabled or not
	forwardToUrl    models.URL               // Location we are going to send logs to
	headers         []ServiceHeader          // Headers to be sent along with logs for identification/authentication
	jsonWrapperKey  string                   // Wrap audit data as a map under this key if present
	environmentName string                   // Decorate the environment this is coming from
	hostname        string                   // The self-reported hostname of the machine
	localIP         string                   // A non-loopback IP address as reported by the machine
	loggingClient   HTTPAuditLoggerInterface // Abstract type for sending logs onward

	loggingChannel chan wrappedAuditLog
	chStop         chan struct{}
	chDone         chan struct{}
}

type wrappedAuditLog struct {
	eventID EventID
	data    Data
}

var NoopLogger AuditLogger = &AuditLoggerService{}

// NewAuditLogger returns a buffer push system that ingests audit log events and
// asynchronously pushes them up to an HTTP log service.
// Parses and validates the AUDIT_LOGS_* environment values and returns an enabled
// AuditLogger instance. If the environment variables are not set, the logger
// is disabled and short circuits execution via enabled flag.
func NewAuditLogger(logger logger.Logger, config Config) (AuditLogger, error) {
	// If the unverified config is nil, then we assume this came from the
	// configuration system and return a nil logger.
	if config == nil || !config.AuditLoggerEnabled() {
		return &AuditLoggerService{}, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Errorf("initialization error - unable to get hostname: %s", err)
	}

	forwardToUrl, err := config.AuditLoggerForwardToUrl()
	if err != nil {
		return &AuditLoggerService{}, nil
	}

	headers, err := config.AuditLoggerHeaders()
	if err != nil {
		return &AuditLoggerService{}, nil
	}

	loggingChannel := make(chan wrappedAuditLog, bufferCapacity)

	// Create new AuditLoggerService
	auditLogger := AuditLoggerService{
		logger:          logger.Helper(1),
		enabled:         true,
		forwardToUrl:    forwardToUrl,
		headers:         headers,
		jsonWrapperKey:  config.AuditLoggerJsonWrapperKey(),
		environmentName: config.AuditLoggerEnvironment(),
		hostname:        hostname,
		localIP:         getLocalIP(),
		loggingClient:   &http.Client{Timeout: time.Second * webRequestTimeout},

		loggingChannel: loggingChannel,
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
	}

	return &auditLogger, nil
}

func (l *AuditLoggerService) SetLoggingClient(newClient HTTPAuditLoggerInterface) {
	l.loggingClient = newClient
}

// Entrypoint for new audit logs. This buffers all logs that come in they will
// sent out by the goroutine that was started when the AuditLoggerService was
// created. If this service was not enabled, this immeidately returns.
//
// This function never blocks.
func (l *AuditLoggerService) Audit(eventID EventID, data Data) {
	if !l.enabled {
		return
	}

	wrappedLog := wrappedAuditLog{
		eventID: eventID,
		data:    data,
	}

	select {
	case l.loggingChannel <- wrappedLog:
	default:
		l.logger.Errorf("buffer is full. Dropping log with eventID: %s", eventID)
	}
}

// Start the audit logger and begin processing logs on the channel
func (l *AuditLoggerService) Start(context.Context) error {
	if !l.enabled {
		return errors.New("The audit logger is not enabled")
	}

	go l.runLoop()
	return nil
}

// Stops the logger and will close the channel.
func (l *AuditLoggerService) Close() error {
	if !l.enabled {
		return errors.New("The audit logger is not enabled")
	}

	l.logger.Warnf("Disabled the audit logger service")
	close(l.chStop)
	<-l.chDone

	return nil
}

func (l *AuditLoggerService) Healthy() error {
	if !l.enabled {
		return errors.New("the audit logger is not enabled")
	}

	if len(l.loggingChannel) == bufferCapacity {
		return errors.New("buffer is full")
	}

	return nil
}

func (l *AuditLoggerService) Name() string {
	return l.logger.Name()
}

func (l *AuditLoggerService) HealthReport() map[string]error {
	return map[string]error{l.logger.Name(): l.Healthy()}
}

func (l *AuditLoggerService) Ready() error {
	if !l.enabled {
		return errors.New("the audit logger is not enabled")
	}

	return nil
}

// Entrypoint for our log handling goroutine. This waits on the channel and sends out
// logs as they come in.
//
// This function calls postLogToLogService which blocks.
func (l *AuditLoggerService) runLoop() {
	defer close(l.chDone)

	for {
		select {
		case <-l.chStop:
			l.logger.Warn("The audit logger is shutting down")
			return
		case event := <-l.loggingChannel:
			l.postLogToLogService(event.eventID, event.data)
		}
	}
}

// Takes an EventID and associated data and sends it to the configured logging
// endpoint. This function blocks on the send by timesout after a period of
// several seconds. This helps us prevent getting stuck on a single log
// due to transient network errors.
//
// This function blocks when called.
func (l *AuditLoggerService) postLogToLogService(eventID EventID, data Data) {
	// Audit log JSON data
	logItem := map[string]interface{}{
		"eventID":  eventID,
		"hostname": l.hostname,
		"localIP":  l.localIP,
		"env":      l.environmentName,
		"data":     data,
	}

	// Optionally wrap audit log data into JSON object to help dynamically structure for an HTTP log service call
	if l.jsonWrapperKey != "" {
		logItem = map[string]interface{}{l.jsonWrapperKey: logItem}
	}

	serializedLog, err := json.Marshal(logItem)
	if err != nil {
		l.logger.Errorw("unable to serialize wrapped audit log item to JSON", "err", err, "logItem", logItem)
		return
	}
	ctx, cancel := utils.ContextFromChan(l.chStop)
	defer cancel()

	// Send to remote service
	req, err := http.NewRequestWithContext(ctx, "POST", (*url.URL)(&l.forwardToUrl).String(), bytes.NewReader(serializedLog))
	if err != nil {
		l.logger.Error("failed to create request to remote logging service!")
	}
	for _, header := range l.headers {
		req.Header.Add(header.Header, header.Value)
	}
	resp, err := l.loggingClient.Do(req)
	if err != nil {
		l.logger.Errorw("failed to send audit log to HTTP log service", "err", err, "logItem", logItem)
		return
	}
	if resp.StatusCode != 200 {
		if resp.Body == nil {
			l.logger.Errorw("no body to read. Possibly an error occurred sending", "logItem", logItem)
			return
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			l.logger.Errorw("error reading errored HTTP log service webhook response body", "err", err, "logItem", logItem)
			return
		}
		l.logger.Errorw("error sending log to HTTP log service", "statusCode", resp.StatusCode, "bodyString", string(bodyBytes))
		return

	}
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
