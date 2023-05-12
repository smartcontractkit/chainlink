package plugins_test

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq" // postgres driver
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

var (
	db               *sql.DB
	dbName           = "loop_integration_new_test"
	pgUser           = "admin"
	pgPassword       = "sixteenCharacter"
	chainlinkBaseUrl string
	autoCleanupOpts  = func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	}
)

func TestMain(m *testing.M) {

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resourcePurger := newPurger(pool)
	// resources registered with the clean with be purged
	defer resourcePurger.cleanup()

	// parallelize node build and postgres startup
	// TODO: in ci can we use an existing chainlink image rather than building here? there ought to be one from the CI setup
	var (
		buildWg  sync.WaitGroup
		buildErr error
		rootDir  string // root directory of the repo to be used for docker run and build contexts
	)

	rootDir, err = filepath.Abs("..")
	if err != nil {
		log.Fatal("could not resolve root dir")
	}
	log.Println("root dir", rootDir)

	buildWg.Add(1)
	loopImageName := "loop-integration2-test"
	go func() {
		defer buildWg.Done()
		buildErr = buildChainlinkImage(rootDir, pool, loopImageName)
	}()

	// postgres has to be running before chainlink bc of DB dependence
	// pulls an image, creates a container based on it and runs it
	pgResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Name:       "pg-15-loop-test",
		Env: []string{
			"POSTGRES_PASSWORD=" + pgPassword,
			"POSTGRES_USER=" + pgUser,
			"POSTGRES_DB=" + dbName,
			"listen_addresses = '*'",
		},
	})
	if err != nil {
		log.Fatalf("Could not start pg resource: %s", err)
	}
	// register to purge postgres container
	resourcePurger.register(pgResource)

	hostAndPort := pgResource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUser, pgPassword, hostAndPort, dbName)
	log.Println("Connecting to database on url: ", databaseUrl)
	pgResource.Expire(300) // Tell docker to hard kill the container in 300 seconds. Acts as a hard cut off for the test suite, too

	pgMaxWait := 90 * time.Second
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = pgMaxWait
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to pg container: %s", err)
	}

	log.Println("connected to postgres at ", databaseUrl)
	log.Println("waiting for chainlink build...")
	buildWg.Wait()
	if buildErr != nil {
		log.Fatalf(buildErr.Error())
	}
	log.Println("starting for chainlink container...")

	chainlinkResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       loopImageName,
		Repository: loopImageName,
		Env: []string{
			// TODO this will need to change to work in CI, probably...
			// it is a hack to ensure the the two container can communicate over the physical host
			"CL_DATABASE_URL=" + strings.Replace(databaseUrl, "localhost", "host.docker.internal", 1),
			"CL_DEV=true",
			"CL_PASSWORD_KEYSTORE=ThisIsATestPassword123456"},
		// hackery to get the container to run the solana loop
		Entrypoint: []string{
			"chainlink",
			"-c", "/run/secrets/node/solana-config.toml",
			"-s", "/run/secrets/node/secure-secrets.toml",
			"node", "start",
			"-d",
			"-p", "/run/secrets/api/password.txt",
			"-a", "/run/secrets/api/apicredentials",
		},

		Mounts: []string{rootDir + "/tools/clroot:/run/secrets/api",
			rootDir + "/plugins/test_data:/run/secrets/node"},
	})

	if err != nil {
		log.Fatalf("failed to run chainlink image %s", err)
	}
	// comment out to keep container for debugging
	resourcePurger.register(chainlinkResource)

	port := chainlinkResource.GetPort("6688/tcp")
	if port == "" {
		log.Fatal("failed to resolve chainlink port 6688")
	}
	chainlinkBaseUrl = fmt.Sprintf("http://localhost:%s", port)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		url := chainlinkBaseUrl + "/health"
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to chainlink container: %s", err)
	}

	//Run tests
	code := m.Run()

	// defer'd call will not run with os.Exit, so if we are here, explicitly cleanup
	resourcePurger.cleanup()

	os.Exit(code)
}

func TestContainerEndpoints(t *testing.T) {

	resp, err := http.Get(chainlinkBaseUrl + "/health")
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = http.Get(chainlinkBaseUrl + "/metrics")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node metrics", string(b))

	resp, err = http.Get(chainlinkBaseUrl + "/discovery")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node discovery", string(b))
	// note that value `Solana` is created by the node (via the logger name today) and could be brittle
	require.Contains(t, string(b), "/plugins/Solana/metrics", "expected solana plugin metric endpoint in %s", b)

	resp, err = http.Get(chainlinkBaseUrl + "/plugins/Solana/metrics")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), "solana_txm_tx_pending", "expected solana specific metric in %s", b)

}

func buildChainlinkImage(ctxDir string, pool *dockertest.Pool, id string) error {
	err := pool.Client.BuildImage(docker.BuildImageOptions{
		Name:         id,
		Dockerfile:   "plugins/chainlink.Dockerfile",
		OutputStream: os.Stderr,
		ContextDir:   ctxDir,
	})

	if err != nil {
		return fmt.Errorf("failed to build chainlink image: %w", err)
	}
	return nil
}

type purger struct {
	pool      *dockertest.Pool
	mu        sync.Mutex
	resources []*dockertest.Resource
}

func newPurger(pool *dockertest.Pool) *purger {
	return &purger{
		pool:      pool,
		resources: make([]*dockertest.Resource, 0),
	}
}

// safe to call multiple times
// UNSAFE to address any resource afterward
func (p *purger) cleanup() error {
	var err error
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, r := range p.resources {
		if r != nil {
			if rerr := p.pool.Purge(r); rerr != nil {
				err = errors.Join(err, rerr)
				r = nil
			}
		}
	}
	return err
}

func (p *purger) register(r *dockertest.Resource) {
	p.mu.Lock()
	p.mu.Unlock()
	p.resources = append(p.resources, r)
}
