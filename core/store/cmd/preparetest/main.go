package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store"
)

var (
	urlStr   = os.Getenv("CL_DATABASE_URL")
	force    = flag.Bool("force", false, "set to true to force the reset by dropping any existing connections to the database")
	userOnly = flag.Bool("user-only", false, "only include test user fixture")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	lggr, err := logger.NewWith(func(z *zap.Config) {
		z.OutputPaths = []string{"stdout"}
		z.Encoding = "console"
		z.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		z.EncoderConfig.TimeKey = ""
		z.EncoderConfig.CallerKey = ""
	})
	if err != nil {
		log.Fatalln(err)
	}

	if urlStr == "" {
		lggr.Fatal("CL_DATABASE_URL is required")
	}
	dbURL, err := url.Parse(urlStr)
	if err != nil {
		lggr.Fatalf("Unable to parse URL %q: %v", urlStr, err)
	}

	if dbname := dbURL.Path[1:]; !strings.HasSuffix(dbname, "_test") {
		lggr.Fatal("Cannot reset database that does not end in _test:", dbURL)
	}

	cfg := config{u: *dbURL}

	if err = store.ResetDatabase(ctx, lggr, cfg, *force); err != nil {
		lggr.Fatal("Failed to reset database:", err)
	}

	if err = store.PrepareTestDB(lggr, *dbURL, *userOnly); err != nil {
		lggr.Fatal("Failed to prepare test database:", err)
	}
}

var _ store.Config = config{}

type config struct{ u url.URL }

func (c config) DefaultIdleInTxSessionTimeout() time.Duration { return time.Hour }

func (c config) DefaultLockTimeout() time.Duration { return 15 * time.Second }

func (c config) MaxOpenConns() int { return 100 }

func (c config) MaxIdleConns() int { return 10 }

func (c config) URL() url.URL { return c.u }

func (c config) DriverName() string { return pg.DriverPostgres }
