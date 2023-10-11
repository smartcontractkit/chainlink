package soak

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

type SoakTest interface {
	Setup(logger zerolog.Logger) error
	Run() error
	Resume(logger zerolog.Logger) error

	NotifyStart() error
	NotifyResume() error
	NotifyEnd() error
	GatherReport() error

	Cleanup() error
	Monitor()

	OnInterrupt() error
	Interrupted() bool
}

func TestSoak(t *testing.T) {
	log := logging.GetTestLogger(t)
	testName := flag.String("test", "", "test name")
	flag.Parse()
	test, err := GetSoakTest(*testName)
	require.NoError(t, err, "Error finding soak test")

	t.Cleanup(func() {
		if err := test.Cleanup(); err != nil {
			log.Error().Err(err).Msg("Error tearing down soak test")
		}
	})

	err = test.Setup(log)
	require.NoError(t, err, "Error setting up soak test")

	// Monitor test state
	go test.Monitor()

	interruptSignal := make(chan os.Signal)
	signal.Notify(interruptSignal, os.Kill, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-interruptSignal
		log.Warn().Str("Signal", sig.String()).Msg("Interrupt signal received, calling OnInterrupt()")
		if err = test.OnInterrupt(); err != nil {
			log.Error().Err(err).Msg("Error while interrupting soak test")
		}

	}()

	if test.Interrupted() {
		if err = test.Resume(log); err != nil {
			log.Err(err).Msg("Error resuming soak test")
		}
	} else {
		if err = test.Run(); err != nil {
			log.Err(err).Msg("Error running soak test")
		}
	}

	if err = test.GatherReport(); err != nil {
		log.Err(err).Msg("Error gathering test report")
	}

}

var soakTests = map[string]SoakTest{}

func RegisterSoakTest(name string, test SoakTest) {
	soakTests[name] = test
}

// GetSoakTest returns the soak test registered with the given name
func GetSoakTest(name string) (SoakTest, error) {
	test, ok := soakTests[name]
	if !ok {
		return nil, fmt.Errorf(
			"no soak test registered with name '%s', did you forget to call RegisterSoakTest? valid names: %s",
			name,
			getSoakTestNames(),
		)
	}
	return test, nil
}

func getSoakTestNames() []string {
	names := make([]string, len(soakTests))
	i := 0
	for name := range soakTests {
		names[i] = name
		i++
	}
	return names
}
