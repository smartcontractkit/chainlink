package servicetest

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

func TestRunHealthy(t *testing.T) {
	for _, test := range []struct {
		name    string
		s       services.Service
		expFail bool
	}{
		{name: "healty", s: &fakeService{}},
		{name: "start", s: &fakeService{start: errors.New("test")}, expFail: true},
		{name: "close", s: &fakeService{close: errors.New("test")}, expFail: true},
		{name: "unready", s: &fakeService{ready: errors.New("test")}, expFail: true},
		{name: "unhealthy", s: &fakeService{healthReport: map[string]error{
			"foo.bar": errors.New("baz"),
		}}, expFail: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, failed := runFake(func(t TestingT) {
				RunHealthy(t, test.s)
			})
			assert.Equal(t, test.expFail, failed)
		})
	}

}

func runFake(fn func(t TestingT)) ([]string, bool) {
	var t fakeTest
	func() {
		defer func() {
			for i := len(t.cleanup) - 1; i >= 0; i-- {
				t.cleanup[i]()
			}
			if r := recover(); r != nil {
				if _, ok := r.(failNow); ok {
					return
				}
				panic(r)
			}
		}()
		fn(&t)
	}()
	return t.errors, t.failed
}

type failNow struct{}

type fakeTest struct {
	cleanup []func()
	errors  []string
	failed  bool
}

func (f *fakeTest) Cleanup(fn func()) {
	f.cleanup = append(f.cleanup, fn)
}

func (f *fakeTest) Errorf(format string, args ...interface{}) {
	f.errors = append(f.errors, fmt.Sprintf(format, args...))
	f.failed = true
}

func (f *fakeTest) FailNow() {
	if f.failed == true {
		return // only panic the first time
	}
	f.failed = true
	panic(failNow{})
}

func (f *fakeTest) Helper() {}

type fakeService struct {
	start        error
	close        error
	ready        error
	healthReport map[string]error
}

func (h *fakeService) Name() string { return "fakeService" }

func (h *fakeService) Start(ctx context.Context) error { return h.start }

func (h *fakeService) Close() error { return h.close }

func (h *fakeService) Ready() error { return h.ready }

func (h *fakeService) HealthReport() map[string]error { return h.healthReport }
