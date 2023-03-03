package blockscout

import (
	"fmt"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/imports/k8s"
	"github.com/smartcontractkit/chainlink-env/pkg"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
)

const (
	URLsKey = "blockscout"
)

type Chart struct {
	Props *Props
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	if m.Props == nil || m.Props.Name == "" {
		return "blockscout"
	}
	return m.Props.Name
}

func (m Chart) GetProps() interface{} {
	return m.Props
}

func (m Chart) GetPath() string {
	return ""
}

func (m Chart) GetVersion() string {
	return ""
}

func (m Chart) GetValues() *map[string]interface{} {
	return nil
}

func (m Chart) ExportData(e *environment.Environment) error {
	bsURL, err := e.Fwd.FindPort(
		fmt.Sprintf("%s:0", m.GetName()),
		fmt.Sprintf("%s-node", m.GetName()), "explorer").
		As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	log.Info().Str("URL", bsURL).Msg("Blockscout explorer")
	e.URLs[URLsKey] = []string{bsURL}
	return nil
}

func New(props *Props) func(root cdk8s.Chart) environment.ConnectedChart {
	return func(root cdk8s.Chart) environment.ConnectedChart {
		dp := defaultProps()
		config.MustMerge(dp, props)
		c := &Chart{
			Props: dp,
		}
		vars := vars{
			Labels: &map[string]*string{
				"app": a.Str(c.GetName()),
			},
			ConfigMapName: fmt.Sprintf("%s-cm", c.GetName()),
			BaseName:      c.GetName(),
			Port:          4000,
			Props:         dp,
		}
		service(root, vars)
		deployment(root, vars)
		return c
	}
}

type Props struct {
	Name    string
	HttpURL string `envconfig:"http_url"`
	WsURL   string `envconfig:"ws_url"`
}

func defaultProps() *Props {
	return &Props{
		HttpURL: "http://geth:8544",
		WsURL:   "ws://geth:8546",
	}
}

// vars some shared labels/selectors and names that must match in resources
type vars struct {
	Labels        *map[string]*string
	BaseName      string
	ConfigMapName string
	Port          float64
	Props         *Props
}

func service(chart cdk8s.Chart, vars vars) {
	k8s.NewKubeService(chart, a.Str(fmt.Sprintf("%s-service", vars.BaseName)), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name: a.Str(vars.BaseName),
		},
		Spec: &k8s.ServiceSpec{
			Ports: &[]*k8s.ServicePort{
				{
					Name:       a.Str("explorer"),
					Port:       a.Num(vars.Port),
					TargetPort: k8s.IntOrString_FromNumber(a.Num(4000)),
				},
			},
			Selector: vars.Labels,
		},
	})
}

func postgresContainer(p vars) *k8s.Container {
	return &k8s.Container{
		Name:  a.Str(fmt.Sprintf("%s-db", p.BaseName)),
		Image: a.Str("postgres:13.6"),
		Ports: &[]*k8s.ContainerPort{
			{
				Name:          a.Str("postgres"),
				ContainerPort: a.Num(5432),
			},
		},
		Env: &[]*k8s.EnvVar{
			a.EnvVarStr("POSTGRES_PASSWORD", "postgres"),
			a.EnvVarStr("POSTGRES_DB", "blockscout"),
		},
		LivenessProbe: &k8s.Probe{
			Exec: &k8s.ExecAction{
				Command: pkg.PGIsReadyCheck()},
			InitialDelaySeconds: a.Num(60),
			PeriodSeconds:       a.Num(60),
		},
		ReadinessProbe: &k8s.Probe{
			Exec: &k8s.ExecAction{
				Command: pkg.PGIsReadyCheck()},
			InitialDelaySeconds: a.Num(2),
			PeriodSeconds:       a.Num(2),
		},
		Resources: a.ContainerResources("1000m", "2048Mi", "1000m", "2048Mi"),
	}
}

func deployment(chart cdk8s.Chart, vars vars) {
	k8s.NewKubeDeployment(
		chart,
		a.Str(fmt.Sprintf("%s-deployment", vars.BaseName)),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name: a.Str(vars.BaseName),
			},
			Spec: &k8s.DeploymentSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: vars.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: vars.Labels,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: a.Str("default"),
						Containers: &[]*k8s.Container{
							container(vars),
							postgresContainer(vars),
						},
					},
				},
			},
		})
}

func container(vars vars) *k8s.Container {
	return &k8s.Container{
		Name:            a.Str(fmt.Sprintf("%s-node", vars.BaseName)),
		Image:           a.Str("f4hrenh9it/blockscout:v1"),
		ImagePullPolicy: a.Str("Always"),
		Command:         &[]*string{a.Str(`/bin/bash`)},
		Args: &[]*string{
			a.Str("-c"),
			a.Str("mix ecto.create && mix ecto.migrate && mix phx.server"),
		},
		Ports: &[]*k8s.ContainerPort{
			{
				Name:          a.Str("explorer"),
				ContainerPort: a.Num(vars.Port),
			},
		},
		ReadinessProbe: &k8s.Probe{
			HttpGet: &k8s.HttpGetAction{
				Port: k8s.IntOrString_FromNumber(a.Num(vars.Port)),
				Path: a.Str("/"),
			},
			InitialDelaySeconds: a.Num(20),
			PeriodSeconds:       a.Num(5),
		},
		Env: &[]*k8s.EnvVar{
			a.EnvVarStr("MIX_ENV", "prod"),
			a.EnvVarStr("ECTO_USE_SSL", "'false'"),
			a.EnvVarStr("COIN", "DAI"),
			a.EnvVarStr("ETHEREUM_JSONRPC_VARIANT", "geth"),
			a.EnvVarStr("ETHEREUM_JSONRPC_HTTP_URL", vars.Props.HttpURL),
			a.EnvVarStr("ETHEREUM_JSONRPC_WS_URL", vars.Props.WsURL),
			a.EnvVarStr("DATABASE_URL", "postgresql://postgres:@localhost:5432/blockscout?ssl=false"),
		},
		Resources: a.ContainerResources("2000m", "2048Mi", "2000m", "2048Mi"),
	}
}
