package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

type SpecResolver struct {
	j job.Job
}

func NewSpec(j job.Job) *SpecResolver {
	return &SpecResolver{j: j}
}

func (r *SpecResolver) ToDirectRequestSpec() (*DirectRequestSpecResolver, bool) {
	if r.j.Type != job.DirectRequest {
		return nil, false
	}

	return &DirectRequestSpecResolver{spec: *r.j.DirectRequestSpec}, true
}

func (r *SpecResolver) ToFluxMonitorSpec() (*FluxMonitorSpecResolver, bool) {
	if r.j.Type != job.FluxMonitor {
		return nil, false
	}

	return &FluxMonitorSpecResolver{spec: *r.j.FluxMonitorSpec}, true
}

type DirectRequestSpecResolver struct {
	spec job.DirectRequestSpec
}

// ContractAddress resolves the spec's contract address.
func (r *DirectRequestSpecResolver) ContractAddress() string {
	return r.spec.ContractAddress.String()
}

// CreatedAt resolves the spec's created at timestamp.
func (r *DirectRequestSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// EVMChainID resolves the spec's evm chain id.
func (r *DirectRequestSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// MinIncomingConfirmations resolves the spec's min incoming confirmations.
func (r *DirectRequestSpecResolver) MinIncomingConfirmations() *int32 {
	if r.spec.MinIncomingConfirmations.Valid {
		min := int32(r.spec.MinIncomingConfirmations.Uint32)
		return &min
	}

	return nil
}

// EVMChainID resolves the spec's evm chain id.
func (r *DirectRequestSpecResolver) MinIncomingConfirmationsEnv() bool {
	return r.spec.MinIncomingConfirmationsEnv
}

// MinContractPayment resolves the spec's evm chain id.
func (r *DirectRequestSpecResolver) MinContractPayment() string {
	return r.spec.MinContractPayment.String()
}

// Requesters resolves the spec's evm chain id.
func (r *DirectRequestSpecResolver) Requesters() []string {
	return r.spec.Requesters.ToStrings()
}

type FluxMonitorSpecResolver struct {
	spec job.FluxMonitorSpec
}

// AbsoluteThreshold resolves the spec's absolute deviation threshold.
func (r *FluxMonitorSpecResolver) AbsoluteThreshold() float64 {
	return float64(r.spec.AbsoluteThreshold)
}

// ContractAddress resolves the spec's contract address.
func (r *FluxMonitorSpecResolver) ContractAddress() string {
	return r.spec.ContractAddress.String()
}

// CreatedAt resolves the spec's created at timestamp.
func (r *FluxMonitorSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// AbsoluteThreshold resolves the spec's absolute threshold.
func (r *FluxMonitorSpecResolver) DrumbeatEnabled() bool {
	return r.spec.DrumbeatEnabled
}

// DrumbeatRandomDelay resolves the spec's drumbeat random delay.
func (r *FluxMonitorSpecResolver) DrumbeatRandomDelay() *string {
	var delay *string
	if r.spec.DrumbeatRandomDelay > 0 {
		drumbeatRandomDelay := r.spec.DrumbeatRandomDelay.String()
		delay = &drumbeatRandomDelay
	}

	return delay
}

// DrumbeatSchedule resolves the spec's drumbeat schedule.
func (r *FluxMonitorSpecResolver) DrumbeatSchedule() *string {
	if r.spec.DrumbeatEnabled {
		return &r.spec.DrumbeatSchedule
	}

	return nil
}

// EVMChainID resolves the spec's evm chain id.
func (r *FluxMonitorSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// IdleTimerDisabled resolves the spec's idle timer disabled flag.
func (r *FluxMonitorSpecResolver) IdleTimerDisabled() bool {
	return r.spec.IdleTimerDisabled
}

// IdleTimerPeriod resolves the spec's idle timer period.
func (r *FluxMonitorSpecResolver) IdleTimerPeriod() string {
	return r.spec.IdleTimerPeriod.String()
}

// MinPayment resolves the spec's min payment.
func (r *FluxMonitorSpecResolver) MinPayment() *string {
	if r.spec.MinPayment != nil {
		min := r.spec.MinPayment.String()

		return &min
	}
	return nil
}

// PollTimerDisabled resolves the spec's poll timer disabled flag.
func (r *FluxMonitorSpecResolver) PollTimerDisabled() bool {
	return r.spec.PollTimerDisabled
}

// PollTimerPeriod resolves the spec's poll timer period.
func (r *FluxMonitorSpecResolver) PollTimerPeriod() string {
	return r.spec.PollTimerPeriod.String()
}

// Threshold resolves the spec's deviation threshold.
func (r *FluxMonitorSpecResolver) Threshold() float64 {
	return float64(r.spec.Threshold)
}
