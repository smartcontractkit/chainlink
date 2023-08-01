package envcommon

import (
	"fmt"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type PgOpts struct {
	User     string
	Password string
	DbName   string
	Networks []string
	Port     string
}

func NewDefaultPgOpts(dbName string, networks []string) PgOpts {
	return PgOpts{
		Port:     "5432",
		User:     "postgres",
		Password: "test",
		DbName:   dbName,
		Networks: networks,
	}
}

func GetPgContainerRequest(name string, opts PgOpts) *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         name,
		Image:        "postgres:15.3",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", opts.Port)},
		Env: map[string]string{
			"POSTGRES_USER":     opts.User,
			"POSTGRES_DB":       opts.DbName,
			"POSTGRES_PASSWORD": opts.Password,
		},
		Networks: opts.Networks,
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "localhost",
			"-U", opts.User, "-c", "select", "1", "-d", opts.DbName}).
			WithStartupTimeout(10 * time.Second),
	}
}

type EnvComponent struct {
	ContainerName string
	Container     tc.Container
	Networks      []string
}

type EnvComponentOpts struct {
	Reuse        bool
	Name         string
	ReplicaIndex int
	ID           string
	Networks     []string
}

func NewReusableName(name string, idx int, id string, reuse bool) string {
	if reuse {
		return fmt.Sprintf("%s-%d", name, idx)
	} else {
		return fmt.Sprintf("%s-%d-%s", name, idx, id)
	}
}

func NewEnvComponent(opts EnvComponentOpts) EnvComponent {
	return EnvComponent{
		ContainerName: NewReusableName(opts.Name, opts.ReplicaIndex, opts.ID, opts.Reuse),
		Networks:      opts.Networks,
	}
}
