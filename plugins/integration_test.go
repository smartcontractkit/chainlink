package plugins_test

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	cleaner := newPurger(pool)

	defer cleaner.cleanup()

	//TODO parallelize container builds and start
	var (
		buildWg sync.WaitGroup
		//errMu sync.Mutex
		buildErr error
	)

	buildWg.Add(1)
	loopImageName := "loop-integration2-test"
	go func() {
		defer buildWg.Done()
		buildErr = buildChainlinkImage(pool, loopImageName)
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
	defer pool.Purge(pgResource)

	cleaner.register(pgResource)
	hostAndPort := pgResource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUser, pgPassword, hostAndPort, dbName)
	log.Println("Connecting to database on url: ", databaseUrl)
	pgResource.Expire(120) // Tell docker to hard kill the container in 300 seconds

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
			"CL_DATABASE_URL=" + strings.Replace(databaseUrl, "localhost", "host.docker.internal", 1),
			"CL_DEV=true",
			"CL_PASSWORD_KEYSTORE=ThisIsATestPassword123456"},
		Entrypoint: []string{
			"chainlink", "-c", "/run/secrets/docker/config.toml", "-s", "/run/secrets/docker/dev-secrets.toml", "node",
			"start", "-d", "-p", "/run/secrets/clroot/password.txt", "-a", "/run/secrets/clroot/apicredentials",
		},
		Mounts: []string{"/Users/kreherma/git/cll/chainlink/tools/clroot:/run/secrets/clroot", "/Users/kreherma/git/cll/chainlink/tools/docker:/run/secrets/docker"},
	})

	if err != nil {
		log.Fatalf("failed to run chainlink image %s", err)
	}
	defer pool.Purge(chainlinkResource)
	// comment out to keep container for debugging
	//cleaner.register(chainlinkResource)

	port := chainlinkResource.GetPort("6688/tcp")
	if port == "" {
		log.Fatal("failed to resolve chainlink port 6688")
	}
	chainlinkBaseUrl = fmt.Sprintf("http://localhost:%s", port)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	// the minio client does not do service discovery for you (i.e. it does not check if connection can be established), so we have to use the health check
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
	cleaner.cleanup()

	os.Exit(code)
}

func buildChainlinkImage(pool *dockertest.Pool, id string) error {
	// Build and run the given Dockerfile
	err := os.Chdir("/Users/kreherma/git/cll/chainlink")
	if err != nil {
		return fmt.Errorf("failed to chdir for building image: %w", err)
	}

	err = pool.Client.BuildImage(docker.BuildImageOptions{
		Name:         id,
		Dockerfile:   "plugins/chainlink.Dockerfile",
		OutputStream: os.Stderr,
		ContextDir:   "/Users/kreherma/git/cll/chainlink",
	})

	if err != nil {
		return fmt.Errorf("failed to build chainlink image: %w", err)
	}
	return nil
}

func purge(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if resource != nil {
		log.Printf("purging resource %s, %s, %s", resource.Container.Name, resource.Container.Image, resource.Container.ID)
		return pool.Purge(resource)
	}
	return nil
}
func TestInDocker(t *testing.T) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resp, err := http.Get(chainlinkBaseUrl + "/health")
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = http.Get(chainlinkBaseUrl + "/metrics")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node metrics", b)

	resp, err = http.Get(chainlinkBaseUrl + "/discovery")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node discovery", b)

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

// safe to call multiple time
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

func TestRealbob(t *testing.T) {
	// all tests
}
