package test_env

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type PostgresDb struct {
	EnvComponent
	User     string
	Password string
	DbName   string
	Port     string
}

type PostgresDbOption = func(c *PostgresDb)

// Sets custom container name if name is not empty
func WithPostgresDbContainerName(name string) PostgresDbOption {
	return func(c *PostgresDb) {
		if name != "" {
			c.ContainerName = name
		}
	}
}

func WithPostgresDbPort(port string) PostgresDbOption {
	return func(c *PostgresDb) {
		c.Port = port
	}
}

func NewPostgresDb(networks []string, opts ...EnvComponentOption) *PostgresDb {
	pg := &PostgresDb{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "postgres-db", uuid.NewString()),
			Networks:      networks,
		},
		User:     "postgres",
		Password: "test",
		DbName:   "testdb",
		Port:     "5342",
	}
	for _, opt := range opts {
		opt(&pg.EnvComponent)
	}
	return pg
}

func (pg *PostgresDb) StartContainer(lw *logwatch.LogWatch) error {
	req := pg.getContainerRequest()
	c, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return err
	}
	pg.Container = c

	return nil
}

func (pg *PostgresDb) getContainerRequest() *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         pg.ContainerName,
		Image:        "postgres:15.3",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", pg.Port)},
		Env: map[string]string{
			"POSTGRES_USER":     pg.User,
			"POSTGRES_DB":       pg.DbName,
			"POSTGRES_PASSWORD": pg.Password,
		},
		Networks: pg.Networks,
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "localhost",
			"-U", pg.User, "-c", "select", "1", "-d", pg.DbName}).
			WithStartupTimeout(10 * time.Second),
	}
}
