package dashboard_lib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

type Dashboard struct {
	Name       string
	DeployOpts EnvConfig
	/* SDK panels that are missing in Grabana */
	SDKPanels []map[string]interface{}
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
	var client *grabana.Client
	if m.DeployOpts.GrafanaBasicAuthUser != "" && m.DeployOpts.GrafanaBasicAuthPassword != "" {
		L.Info().Msg("Authorizing using BasicAuth")
		client = grabana.NewClient(
			&http.Client{},
			m.DeployOpts.GrafanaURL,
			grabana.WithBasicAuth(m.DeployOpts.GrafanaBasicAuthUser, m.DeployOpts.GrafanaBasicAuthPassword),
		)
	} else {
		L.Info().Msg("Authorizing using Bearer token")
		client = grabana.NewClient(
			&http.Client{},
			m.DeployOpts.GrafanaURL,
			grabana.WithAPIToken(m.DeployOpts.GrafanaToken),
		)
	}
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

func (m *Dashboard) AddSDKPanel(panel map[string]interface{}) {
	m.SDKPanels = append(m.SDKPanels, panel)
}

func (m *Dashboard) Delete() error {
	ctx := context.Background()
	var client *grabana.Client
	if m.DeployOpts.GrafanaBasicAuthUser != "" && m.DeployOpts.GrafanaBasicAuthPassword != "" {
		L.Info().Msg("Authorizing using BasicAuth")
		client = grabana.NewClient(
			&http.Client{},
			m.DeployOpts.GrafanaURL,
			grabana.WithBasicAuth(m.DeployOpts.GrafanaBasicAuthUser, m.DeployOpts.GrafanaBasicAuthPassword),
		)
	} else {
		L.Info().Msg("Authorizing using Bearer token")
		client = grabana.NewClient(
			&http.Client{},
			m.DeployOpts.GrafanaURL,
			grabana.WithAPIToken(m.DeployOpts.GrafanaToken),
		)
	}
	db, err := client.GetDashboardByTitle(ctx, m.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get the dashboard")
	}
	fmt.Println(db.UID)
	errDelete := client.DeleteDashboard(ctx, db.UID)
	if errDelete != nil {
		return errors.Wrap(errDelete, "failed to delete the dashboard")
	}
	return nil
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

// TODO: re-write after forking Grabana, inject foundation SDK components from official schema
func (m *Dashboard) injectSDKPanels(b dashboard.Builder) (dashboard.Builder, error) {
	data, err := b.MarshalIndentJSON()
	if err != nil {
		return dashboard.Builder{}, err
	}
	var asMap map[string]interface{}
	if err := json.Unmarshal(data, &asMap); err != nil {
		return dashboard.Builder{}, err
	}
	asMap["rows"].([]interface{})[0].(map[string]interface{})["panels"] = append(asMap["rows"].([]interface{})[0].(map[string]interface{})["panels"].([]interface{}), m.SDKPanels[0])
	d, err := json.Marshal(asMap)
	if err != nil {
		return dashboard.Builder{}, err
	}
	if err := os.WriteFile("generated_ccip_dashboard.json", d, os.ModePerm); err != nil {
		return dashboard.Builder{}, err
	}
	return b, nil
}
