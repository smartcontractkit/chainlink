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
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq" // postgres driver
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	db                      *sql.DB
	dbName                  = "loop_integration_new_test"
	pgUser                  = "admin"
	pgPassword              = "sixteenCharacter"
	pgDockerHostName        = "pg_host"
	pgDockerPort            = 5432
	chainlinkDockerHostName = "chainlink_host"
	chainlinkBaseUrl        string
	autoCleanupOpts         = func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	}
)

func dburl(hostAndPort string) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUser, pgPassword, hostAndPort, dbName)
}
func TestMain(m *testing.M) {
	env := os.Environ()
	sort.Strings(env)
	log.Printf("env %v", strings.Join(env, "\n"))
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
	defer func() {
		if cerr := resourcePurger.cleanup(); cerr != nil {
			log.Fatalf("failed to cleanup resources: %s", cerr)
		}
	}()

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

	var (
		network *docker.Network
	)
	network, err = pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: "pg_chainlink_network"})
	if err != nil {
		log.Fatalf("could not create a network to pg and chainlink %s", err)
	}
	n2, err := pool.CreateNetwork("dt-net")
	if err != nil {
		log.Fatalf("error creating test network: %s", err)
	}
	log.Printf("setup docker network %+v %s %s", *network, network.ID, network.Name)
	defer func() {
		pruned, err := pool.Client.PruneImages(docker.PruneImagesOptions{})
		if err != nil {
			log.Printf("failed to prune images after build: %s", err)
		}
		log.Printf("pruned %+v", *pruned)
		pool.Client.RemoveNetwork(network.ID)
	}()

	pgResource, err := runPostgresContainer(pool, n2.Network)
	if err != nil {
		log.Fatal(err)
	}
	resourcePurger.register(pgResource)

	log.Println("waiting for chainlink build...")
	buildWg.Wait()
	if buildErr != nil {
		log.Fatalf(buildErr.Error())
	}

	dbUrlForNode := dburl(fmt.Sprintf("%s:%d", pgDockerHostName, pgDockerPort))

	log.Printf("starting for chainlink container with db %s ...", dbUrlForNode)

	nodeContainerStdout := newStreamHack("node.container.stdout")
	defer nodeContainerStdout.Close()
	nodeContainerStderr := newStreamHack("node.container.stderr")
	defer nodeContainerStderr.Close()
	lerr := pool.Client.Logs(docker.LogsOptions{
		Container:         loopImageName,
		OutputStream:      nodeContainerStdout,
		ErrorStream:       nodeContainerStderr,
		InactivityTimeout: 2 * time.Second,
		Stdout:            true,
		Stderr:            true,
		Timestamps:        true,
		RawTerminal:       false,
	})
	if lerr != nil {
		log.Printf("failed to get chainlink container logs %s", lerr)
	}
	chainlinkResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       loopImageName,
		Repository: loopImageName,
		Env: []string{
			// TODO this will need to change to work in CI, probably...
			// it is a hack to ensure the the two container can communicate over the physical host
			//"CL_DATABASE_URL=" + strings.Replace(databaseUrl, "localhost", "host.docker.internal", 1),
			"CL_DATABASE_URL=" + dbUrlForNode,
		},
		NetworkID: n2.Network.ID,
		Hostname:  "chainlink_host",

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
	//resourcePurger.register(chainlinkResource)
	err = chainlinkResource.Expire(600) // Tell docker to hard kill the container in 300 seconds. Acts as a hard cut off for the test suite, too
	if err != nil {
		log.Fatalf("failed to set pg expiry: %s", err)
	}

	port := chainlinkResource.GetPort("6688/tcp")
	if port == "" {
		log.Fatal("failed to resolve chainlink port 6688")
	}
	chainlinkBaseUrl = fmt.Sprintf("http://localhost:%s", port)

	log.Println("polling for start...")
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		log.Println("...attempting to connect...")
		url := chainlinkBaseUrl + "/health"
		resp, gerr := http.Get(url) // nolint
		if gerr != nil {
			return gerr
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		// dump the container logs
		nodeContainerStdout := newStreamHack("node.container.stdout")
		defer nodeContainerStdout.Close()
		nodeContainerStderr := newStreamHack("node.container.stderr")
		defer nodeContainerStderr.Close()
		lerr := pool.Client.Logs(docker.LogsOptions{
			Container:         loopImageName,
			OutputStream:      nodeContainerStdout,
			ErrorStream:       nodeContainerStderr,
			InactivityTimeout: 2 * time.Second,
			Stdout:            true,
			Stderr:            true,
			Timestamps:        true,
			RawTerminal:       false,
		})
		if lerr != nil {
			log.Printf("failed to get chainlink container logs %s", lerr)
		}
		log.Fatalf("Could not connect to chainlink container: %s", err)
	}

	//Run tests
	code := m.Run()

	// defer'd call will not run with os.Exit, so if we are here, explicitly cleanup
	if err := resourcePurger.cleanup(); err != nil {
		log.Fatalf("failed to cleanup resources after tests finished with code %d: %s", code, err)
	}

	os.Exit(code)
}

type streamHack struct {
	l       logger.Logger
	closeFn func() error
}

func (s *streamHack) Write(b []byte) (int, error) {
	s.l.Info(string(b))
	return len(b), nil
}

func newStreamHack(name string) *streamHack {
	l, closeFn := logger.NewLogger()

	return &streamHack{
		l:       l.Named(name),
		closeFn: closeFn,
	}
}

func (s *streamHack) Close() error {
	return s.closeFn()
}

func TestContainerEndpoints(t *testing.T) {

	resp, err := http.Get(chainlinkBaseUrl + "/health") //nolint
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = http.Get(chainlinkBaseUrl + "/metrics") //nolint
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node metrics", string(b))

	resp, err = http.Get(chainlinkBaseUrl + "/discovery") //nolint
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node discovery", string(b))
	// note that value `Solana` is created by the node (via the logger name today) and could be brittle
	require.Contains(t, string(b), "/plugins/Solana/metrics", "expected solana plugin metric endpoint in %s", b)

	resp, err = http.Get(chainlinkBaseUrl + "/plugins/Solana/metrics") //nolint
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), "solana_txm_tx_pending", "expected solana specific metric in %s", b)

}

func runPostgresContainer(pool *dockertest.Pool, network *docker.Network) (*dockertest.Resource, error) {
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
		NetworkID:    network.ID,
		Hostname:     "pg_host",
		ExposedPorts: []string{"5432/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5423/tcp": {{HostIP: "localhost", HostPort: "5432/tcp"}},
		},
	}, autoCleanupOpts)
	if err != nil {
		return nil, fmt.Errorf("Could not start pg resource: %w", err)
	}
	// register to purge postgres container
	//resourcePurger.register(pgResource)

	hostAndPort := pgResource.GetHostPort("5432/tcp")
	databaseUrl := dburl(hostAndPort)
	log.Println("Connecting to database on url: ", databaseUrl)
	err = pgResource.Expire(300) // Tell docker to hard kill the container in 300 seconds. Acts as a hard cut off for the test suite, too
	if err != nil {
		return pgResource, fmt.Errorf("failed to set pg expiry: %w", err)
	}

	pgMaxWait := 120 * time.Second
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = pgMaxWait
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return pgResource, fmt.Errorf("Could not connect to pg container: %w", err)
	}

	log.Println("connected to postgres container at ", databaseUrl)
	return pgResource, nil
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
	resources map[string]*wrappedResource
}

type wrappedResource struct {
	purged   bool
	resource *dockertest.Resource
}

func newPurger(pool *dockertest.Pool) *purger {
	return &purger{
		pool:      pool,
		resources: make(map[string]*wrappedResource),
	}
}

// safe to call multiple times
func (p *purger) cleanup() error {
	var err error
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, r := range p.resources {
		if !r.purged {
			if rerr := p.pool.Purge(r.resource); rerr != nil {
				err = errors.Join(err, rerr)
				r.purged = true
			}
		}
	}
	return err
}

func (p *purger) register(r *dockertest.Resource) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.resources[r.Container.Name] = &wrappedResource{purged: false, resource: r}
}
