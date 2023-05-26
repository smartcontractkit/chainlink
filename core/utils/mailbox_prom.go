package utils

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var mailboxLoad = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "mailbox_load_percent",
	Help: "Percent of mailbox capacity used",
},
	[]string{"appID", "name", "capacity"},
)

const mailboxPromInterval = 5 * time.Second

type MailboxMonitor struct {
	StartStopOnce
	appID string

	mailboxes sync.Map
	stop      func()
}

func NewMailboxMonitor(appID string) *MailboxMonitor {
	return &MailboxMonitor{appID: appID}
}

func (m *MailboxMonitor) Name() string { return "MailboxMonitor" }

func (m *MailboxMonitor) Start(context.Context) error {
	return m.StartOnce("MailboxMonitor", func() error {
		t := time.NewTicker(WithJitter(mailboxPromInterval))
		m.stop = t.Stop
		go m.monitorLoop(t.C)
		return nil
	})
}

func (m *MailboxMonitor) Close() error {
	return m.StopOnce("MailboxMonitor", func() error {
		m.stop()
		return nil
	})
}

func (m *MailboxMonitor) HealthReport() map[string]error {
	return map[string]error{m.Name(): m.StartStopOnce.Healthy()}
}

func (m *MailboxMonitor) monitorLoop(c <-chan time.Time) {
	for range c {
		m.mailboxes.Range(func(k, v any) bool {
			name, mb := k.(string), v.(mailbox)
			c, p := mb.load()
			capacity := strconv.FormatUint(c, 10)
			mailboxLoad.WithLabelValues(m.appID, name, capacity).Set(p)
			return true
		})
	}
}

type mailbox interface {
	load() (capacity uint64, percent float64)
	onClose(func())
}

func (m *MailboxMonitor) Monitor(mb mailbox, name ...string) {
	n := strings.Join(name, ".")
	m.mailboxes.Store(n, mb)
	mb.onClose(func() { m.mailboxes.Delete(n) })
}
