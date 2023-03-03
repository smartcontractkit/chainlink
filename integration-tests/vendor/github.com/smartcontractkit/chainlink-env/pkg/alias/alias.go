package alias

import (
	"fmt"
	"strings"
	"time"

	jsii "github.com/aws/jsii-runtime-go"
	"github.com/smartcontractkit/chainlink-env/imports/k8s"
)

func Str(value string) *string {
	return jsii.String(value)
}

func Num(value float64) *float64 {
	return jsii.Number(value)
}

// ShortDur is a helper method for kube-janitor duration format
func ShortDur(d time.Duration) *string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return Str(s)
}

func ConvertLabels(labels []string) (*map[string]*string, error) {
	cdk8sLabels := make(map[string]*string)
	for _, s := range labels {
		a := strings.Split(s, "=")
		if len(a) != 2 {
			return nil, fmt.Errorf("invalid label '%s' provided, please provide labels in format key=value", a)
		}
		cdk8sLabels[a[0]] = Str(a[1])
	}
	return &cdk8sLabels, nil
}

// EnvVarStr quick shortcut for string/string key/value var
func EnvVarStr(k, v string) *k8s.EnvVar {
	return &k8s.EnvVar{
		Name:  Str(k),
		Value: Str(v),
	}
}

// ContainerResources container resource requirements
func ContainerResources(reqCPU, reqMEM, limCPU, limMEM string) *k8s.ResourceRequirements {
	return &k8s.ResourceRequirements{
		Requests: &map[string]k8s.Quantity{
			"cpu":    k8s.Quantity_FromString(Str(reqCPU)),
			"memory": k8s.Quantity_FromString(Str(reqMEM)),
		},
		Limits: &map[string]k8s.Quantity{
			"cpu":    k8s.Quantity_FromString(Str(limCPU)),
			"memory": k8s.Quantity_FromString(Str(limMEM)),
		},
	}
}
