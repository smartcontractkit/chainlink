package dashboardlib

import (
	"context"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/pkg/errors"
	"net/http"
)

type Dashboard struct {
	Name       string
	DeployOpts EnvConfig
	/* generated dashboard opts and builder */
	builder dashboard.Builder
	Opts    []dashboard.Option
}

func NewDashboard(
	name string,
	deployOpts EnvConfig,
	opts []dashboard.Option,
) *Dashboard {
	return &Dashboard{
		Name:       name,
		DeployOpts: deployOpts,
		Opts:       opts,
	}
}

func (m *Dashboard) Deploy() error {
	ctx := context.Background()
	b, err := m.build()
	if err != nil {
		return err
	}
	client := grabana.NewClient(&http.Client{}, m.DeployOpts.GrafanaURL, grabana.WithAPIToken(m.DeployOpts.GrafanaToken))
	fo, folderErr := client.FindOrCreateFolder(ctx, m.DeployOpts.GrafanaFolder)
	if folderErr != nil {
		return errors.Wrap(err, "could not find or create Grafana folder")
	}
	if _, err := client.UpsertDashboard(ctx, fo, b); err != nil {
		return errors.Wrap(err, "failed to upsert the dashboard")
	}
	return nil
}

func (m *Dashboard) Add(opts []dashboard.Option) {
	m.Opts = append(m.Opts, opts...)
}

func (m *Dashboard) build() (dashboard.Builder, error) {
	b, err := dashboard.New(
		m.Name,
		m.Opts...,
	)
	if err != nil {
		return dashboard.Builder{}, errors.Wrap(err, "failed to build the dashboard")
	}
	return b, nil
}
