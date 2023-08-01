package envcommon

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"os"
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
	ReuseContainerName string
	Networks           []string
}

func NewEnvComponent(componentType string, opts EnvComponentOpts) EnvComponent {
	var name string
	if opts.ReuseContainerName != "" {
		name = opts.ReuseContainerName
	} else {
		name = fmt.Sprintf("%s-%s", componentType, uuid.NewString())
	}

	return EnvComponent{
		ContainerName: name,
		Networks:      opts.Networks,
	}
}

func ParseJSONFile(path string, v any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	b, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}
