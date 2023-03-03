package client

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
)

// Query constants
const (
	QueryMemoryUsage          = `100 * (1 - ((avg_over_time(node_memory_MemFree_bytes[%s]) + avg_over_time(node_memory_Cached_bytes[%s]) + avg_over_time(node_memory_Buffers_bytes[%s])) / avg_over_time(node_memory_MemTotal_bytes[%s])))`
	QueryAllCPUBusyPercentage = `100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[%s])) * 100)`
)

// ResourcesSummary Data for memory and cpu usage
type ResourcesSummary struct {
	MemoryUsage   float64
	CPUPercentage float64
}

// Prometheus ease of use for v1.API
type Prometheus struct {
	v1.API
}

// NewPrometheusClient returns a prometheus client
func NewPrometheusClient(url string) (*Prometheus, error) {
	client, err := api.NewClient(api.Config{
		Address: url,
	})
	if err != nil {
		return nil, err
	}
	return &Prometheus{
		API: v1.NewAPI(client),
	}, nil
}

func (p *Prometheus) printWarns(warns v1.Warnings) {
	if len(warns) > 0 {
		log.Info().Interface("Warnings", warns).Msg("Warnings found when performing prometheus query")
	}
}

func (p *Prometheus) validateNotEmptyVec(q string, val model.Value) bool {
	if len(val.(model.Vector)) == 0 {
		log.Warn().Str("query", q).Msg("empty response for prometheus query")
		return false
	}
	return true
}

// CPUBusyPercentage host CPU busy percentage
func (p *Prometheus) CPUBusyPercentage() (float64, error) {
	q := fmt.Sprintf(QueryAllCPUBusyPercentage, "2m")
	val, warns, err := p.API.Query(context.Background(), q, time.Now())
	if err != nil {
		return 0, err
	}
	p.printWarns(warns)
	if !p.validateNotEmptyVec(q, val) {
		return 0, nil
	}
	scalarVal := val.(model.Vector)[0].Value
	return float64(scalarVal), nil
}

// MemoryUsage total memory used by interval
func (p *Prometheus) MemoryUsage() (float64, error) {
	q := fmt.Sprintf(QueryMemoryUsage, "2m", "2m", "2m", "2m")
	val, warns, err := p.API.Query(context.Background(), q, time.Now())
	if err != nil {
		return 0, err
	}
	p.printWarns(warns)
	if !p.validateNotEmptyVec(q, val) {
		return 0, nil
	}
	scalarVal := val.(model.Vector)[0].Value
	return float64(scalarVal), nil
}

// ResourcesSummary returns a summary of memory and cpu used
func (p *Prometheus) ResourcesSummary() (float64, float64, error) {
	cpu, err := p.CPUBusyPercentage()
	if err != nil {
		return 0, 0, err
	}
	mem, err := p.MemoryUsage()
	if err != nil {
		return 0, 0, err
	}
	return cpu, mem, nil
}

// GetAlerts returns all firing alerts
func (p *Prometheus) GetAlerts() (v1.AlertsResult, error) {
	alerts, err := p.API.Alerts(context.Background())
	if err != nil {
		return v1.AlertsResult{}, err
	}
	return alerts, nil
}

// GetQuery returns the result of applying a PromQL query
func (p *Prometheus) GetQuery(query string) (model.Value, error) {
	value, _, err := p.API.Query(context.Background(), query, time.Now())
	if err != nil {
		return nil, err
	}
	return value, nil
}
