package test_env

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func NewPostgresDb(networks []string, opts ...PostgresDbOption) *PostgresDb {
	pg := &PostgresDb{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "postgres-db", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		User:     "postgres",
		Password: "mysecretpassword",
		DbName:   "testdb",
		Port:     "5432",
	}
	for _, opt := range opts {
		opt(pg)
	}
	return pg
}

func (pg *PostgresDb) StartContainer() error {
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

	log.Info().Str("containerName", pg.ContainerName).
		Msg("Started Postgres DB container")

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
