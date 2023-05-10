package plugins_test

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
)

func TestInDocker(t *testing.T) {
	dbName := "loop_docker_test"
	var db *sql.DB
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Build and run the given Dockerfile
	require.NoError(t, os.Chdir("/Users/kreherma/git/cll/chainlink"))
	err = pool.Client.BuildImage(docker.BuildImageOptions{
		Name:         "loop-test1",
		Dockerfile:   "plugins/chainlink.Dockerfile",
		OutputStream: os.Stderr,
		ContextDir:   "/Users/kreherma/git/cll/chainlink",
	})

	require.NoError(t, err)

	//pool.BuildAndRun()
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "loop-test1",
		Repository: "loop-test1",
		Env:        []string{"CL_DATABASE_URL=postgres://postgres:admin@localhost:5432/chainlink_plugin_test?sslmode=disable", "CL_DEV=true", "CL_PASSWORD_KEYSTORE=ThisIsATestPassword123456"}})

	require.NoError(t, err)
	t.Cleanup(func() {
		// When you're done, kill and remove the container
		if resource != nil {
			require.NoError(t, pool.Purge(resource), "could not purge resource")
		}
	})

	port := resource.GetPort("6688/tcp")
	require.NotEmpty(t, port)
	base := fmt.Sprintf("http://localhost:%s", port)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	// the minio client does not do service discovery for you (i.e. it does not check if connection can be established), so we have to use the health check
	if err := pool.Retry(func() error {
		url := base + "/health"
		resp, err := http.Get(url)
		t.Logf("trying health check")
		if err != nil {

			t.Logf("retrying err %s %s", url, err)
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resp, err := http.Get(base + "/health")
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = http.Get(base + "/metrics")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node metrics", b)

	resp, err = http.Get("http://localhost:6688/discovery")
	require.NoError(t, err)
	require.NotNil(t, resp)
	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log("node discovery", b)

	/*
		if err = pool.Retry(func() error {
			var err error
			db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), dbName))
			if err != nil {
				return err
			}
			return db.Ping()
		}); err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}

		// When you're done, kill and remove the container
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	*/
}
