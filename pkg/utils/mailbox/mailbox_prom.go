package mailbox

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
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

	mailboxes sync.Map
	stop      func()
	done      chan struct{}
}

func NewMonitor(appID string) *Monitor {
	return &Monitor{appID: appID}
}

func (m *Monitor) Name() string { return "Monitor" }

func (m *Monitor) Start(context.Context) error {
	return m.StartOnce("Monitor", func() error {
		t := time.NewTicker(utils.WithJitter(mailboxPromInterval))
		ctx, cancel := context.WithCancel(context.Background())
		m.stop = func() {
			t.Stop()
			cancel()
		}
		m.done = make(chan struct{})
		go m.monitorLoop(ctx, t.C)
		return nil
	})
}

func (m *Monitor) Close() error {
	return m.StopOnce("Monitor", func() error {
		m.stop()
		<-m.done
		return nil
	})
}

func (m *Monitor) HealthReport() map[string]error {
	return map[string]error{m.Name(): m.Healthy()}
}

func (m *Monitor) monitorLoop(ctx context.Context, c <-chan time.Time) {
	defer close(m.done)
	for {
		select {
		case <-ctx.Done():
			return
		case <-c:
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
