package mailbox

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var mailboxLoad = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "mailbox_load_percent",
	Help: "Percent of mailbox capacity used",
},
	[]string{"appID", "name", "capacity"},
)

const mailboxPromInterval = 5 * time.Second

type Monitor struct {
	services.StateMachine
	appID string
	lggr  logger.Logger

	mailboxes sync.Map
	stopCh    services.StopChan
	done      chan struct{}
}

func NewMonitor(appID string, lggr logger.Logger) *Monitor {
	return &Monitor{appID: appID, lggr: logger.Named(lggr, "Monitor"), stopCh: make(services.StopChan), done: make(chan struct{})}
}

func (m *Monitor) Name() string { return m.lggr.Name() }

func (m *Monitor) Start(context.Context) error {
	return m.StartOnce("Monitor", func() error {
		go m.monitorLoop()
		return nil
	})
}

func (m *Monitor) Close() error {
	return m.StopOnce("Monitor", func() error {
		close(m.stopCh)
		<-m.done
		return nil
	})
}

func (m *Monitor) HealthReport() map[string]error {
	return map[string]error{m.Name(): m.Healthy()}
}

func (m *Monitor) monitorLoop() {
	defer close(m.done)
	t := services.NewTicker(mailboxPromInterval)
	defer t.Stop()
	for {
		select {
		case <-m.stopCh:
			return
		case <-t.C:
			m.mailboxes.Range(func(k, v any) bool {
				name, mb := k.(string), v.(mailbox)
				c, p := mb.load()
				capacity := strconv.FormatUint(c, 10)
				mailboxLoad.WithLabelValues(m.appID, name, capacity).Set(p)
				return true
			})
		}
	}
}

type mailbox interface {
	load() (capacity uint64, percent float64)
	onClose(func())
}

func (m *Monitor) Monitor(mb mailbox, name ...string) {
	n := strings.Join(name, ".")
	m.mailboxes.Store(n, mb)
	mb.onClose(func() { m.mailboxes.Delete(n) })
}
