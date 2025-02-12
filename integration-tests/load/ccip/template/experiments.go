package template

import (
	"context"
	"fmt"
	"time"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
)

func defaultListeners(l zerolog.Logger) []havoc.ChaosListener {
	return []havoc.ChaosListener{
		havoc.NewChaosLogger(l),
	}
}

type ChaosRunner struct {
	l zerolog.Logger
	c client.Client
}

func NewChaosRunner(l zerolog.Logger, c client.Client) *ChaosRunner {
	return &ChaosRunner{
		l: l,
		c: c,
	}
}

type PodPartitionCfg struct {
	Namespace             string
	Description           string
	LabelFromKey          string
	LabelFromValues       []string
	LabelToKey            string
	LabelToValues         []string
	InjectionDuration     time.Duration
	ExperimentCreateDelay time.Duration
}

func (cr *ChaosRunner) RunPodPartition(ctx context.Context, cfg PodPartitionCfg) (*havoc.Chaos, error) {
	experiment, err := havoc.NewChaos(havoc.ChaosOpts{
		Object: &v1alpha1.NetworkChaos{
			TypeMeta: metav1.TypeMeta{
				Kind:       string(v1alpha1.TypeNetworkChaos),
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("partition-%s", uuid.NewString()[0:5]),
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.NetworkChaosSpec{
				Action:   v1alpha1.PartitionAction,
				Duration: ptr.To[string]((cfg.InjectionDuration).String()),
				PodSelector: v1alpha1.PodSelector{
					Mode: v1alpha1.AllMode,
					Selector: v1alpha1.PodSelectorSpec{
						GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
							Namespaces: []string{cfg.Namespace},
							ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
								{
									Operator: "In",
									Key:      cfg.LabelFromKey,
									Values:   cfg.LabelFromValues,
								},
							},
						},
					},
				},
				Target: &v1alpha1.PodSelector{
					Mode: v1alpha1.AllMode,
					Selector: v1alpha1.PodSelectorSpec{
						GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
							Namespaces: []string{cfg.Namespace},
							ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
								{
									Operator: "In",
									Key:      cfg.LabelToKey,
									Values:   cfg.LabelToValues,
								},
							},
						},
					},
				},
			},
		},
		Listeners: defaultListeners(cr.l),
		Logger:    &cr.l,
		Client:    cr.c,
	})
	if err != nil {
		return nil, err
	}
	experiment.Create(ctx)
	return experiment, nil
}

type PodDelayCfg struct {
	Namespace             string
	Description           string
	Latency               time.Duration
	Jitter                time.Duration
	Correlation           string
	LabelKey              string
	LabelValues           []string
	InjectionDuration     time.Duration
	ExperimentCreateDelay time.Duration
}

func (cr *ChaosRunner) RunPodDelay(ctx context.Context, cfg PodDelayCfg) (*havoc.Chaos, error) {
	experiment, err := havoc.NewChaos(havoc.ChaosOpts{
		Object: &v1alpha1.NetworkChaos{
			TypeMeta: metav1.TypeMeta{
				Kind:       string(v1alpha1.TypeNetworkChaos),
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("delay-%s", uuid.NewString()[0:5]),
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.NetworkChaosSpec{
				Action:   v1alpha1.DelayAction,
				Duration: ptr.To[string]((cfg.InjectionDuration).String()),
				TcParameter: v1alpha1.TcParameter{
					Delay: &v1alpha1.DelaySpec{
						Latency:     cfg.Latency.String(),
						Correlation: cfg.Correlation,
						Jitter:      cfg.Jitter.String(),
					},
				},
				PodSelector: v1alpha1.PodSelector{
					Mode: v1alpha1.AllMode,
					Selector: v1alpha1.PodSelectorSpec{
						GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
							Namespaces: []string{cfg.Namespace},
							ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
								{
									Operator: "In",
									Key:      cfg.LabelKey,
									Values:   cfg.LabelValues,
								},
							},
						},
					},
				},
			},
		},
		Listeners: defaultListeners(cr.l),
		Logger:    &cr.l,
		Client:    cr.c,
	})
	if err != nil {
		return nil, err
	}
	experiment.Create(ctx)
	return experiment, nil
}

type PodFailCfg struct {
	Namespace             string
	Description           string
	LabelKey              string
	LabelValues           []string
	InjectionDuration     time.Duration
	ExperimentCreateDelay time.Duration
}

func (cr *ChaosRunner) RunPodFail(ctx context.Context, cfg PodFailCfg) (*havoc.Chaos, error) {
	experiment, err := havoc.NewChaos(havoc.ChaosOpts{
		Description: cfg.Description,
		DelayCreate: cfg.ExperimentCreateDelay,
		Object: &v1alpha1.PodChaos{
			TypeMeta: metav1.TypeMeta{
				Kind:       string(v1alpha1.TypePodChaos),
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("fail-%s", uuid.NewString()[0:5]),
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.PodChaosSpec{
				Action:   v1alpha1.PodFailureAction,
				Duration: ptr.To[string](cfg.InjectionDuration.String()),
				ContainerSelector: v1alpha1.ContainerSelector{
					PodSelector: v1alpha1.PodSelector{
						Mode: v1alpha1.AllMode,
						Selector: v1alpha1.PodSelectorSpec{
							GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
								Namespaces: []string{cfg.Namespace},
								ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
									{
										Operator: "In",
										Key:      cfg.LabelKey,
										Values:   cfg.LabelValues,
									},
								},
							},
						},
					},
				},
			},
		},
		Listeners: defaultListeners(cr.l),
		Logger:    &cr.l,
		Client:    cr.c,
	})
	if err != nil {
		return nil, err
	}
	experiment.Create(ctx)
	return experiment, nil
}

type NodeCPUStressConfig struct {
	Namespace               string
	Description             string
	Cores                   int
	CoreLoadPercentage      int // 0-100
	LabelKey                string
	LabelValues             []string
	InjectionDuration       time.Duration
	ExperimentTotalDuration time.Duration
	ExperimentCreateDelay   time.Duration
}

func (cr *ChaosRunner) RunPodStressCPU(ctx context.Context, cfg NodeCPUStressConfig) (*havoc.Schedule, error) {
	experiment, err := havoc.NewSchedule(havoc.ScheduleOpts{
		Description: cfg.Description,
		DelayCreate: cfg.ExperimentCreateDelay,
		Duration:    cfg.ExperimentTotalDuration,
		Object: &v1alpha1.Schedule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Schedule",
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "stress",
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.ScheduleSpec{
				Schedule:          "@every 1m",
				ConcurrencyPolicy: v1alpha1.ForbidConcurrent,
				Type:              v1alpha1.ScheduleTypeStressChaos,
				HistoryLimit:      2,
				ScheduleItem: v1alpha1.ScheduleItem{
					EmbedChaos: v1alpha1.EmbedChaos{
						StressChaos: &v1alpha1.StressChaosSpec{
							ContainerSelector: v1alpha1.ContainerSelector{
								PodSelector: v1alpha1.PodSelector{
									Mode: v1alpha1.AllMode,
									Selector: v1alpha1.PodSelectorSpec{
										GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
											Namespaces: []string{cfg.Namespace},
											ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
												{
													Operator: "In",
													Key:      cfg.LabelKey,
													Values:   cfg.LabelValues,
												},
											},
										},
									},
								},
							},
							Stressors: &v1alpha1.Stressors{
								CPUStressor: &v1alpha1.CPUStressor{
									Stressor: v1alpha1.Stressor{
										Workers: cfg.Cores,
									},
									Load: ptr.To[int](cfg.CoreLoadPercentage),
								},
							},
							Duration: ptr.To[string](cfg.InjectionDuration.String()),
						},
					},
				},
			},
		},
		Listeners: defaultListeners(cr.l),
		Logger:    &cr.l,
		Client:    cr.c,
	})
	if err != nil {
		return nil, err
	}
	experiment.Create(ctx)
	return experiment, nil
}
