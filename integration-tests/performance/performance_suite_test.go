package performance_test

import (
	"os"
	"testing"

	"github.com/smartcontractkit/integrations-framework/config"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestPerformance(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	RegisterFailHandler(Fail)
	conf, err := config.NewConfig("./")
	if err != nil {
		Fail("failed to load config")
	}
	log.Logger = log.Logger.Level(zerolog.Level(conf.Logging.Level))
	junitReporter := reporters.NewJUnitReporter("../logs/tests-performance.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Performance suite", []Reporter{junitReporter})
}
