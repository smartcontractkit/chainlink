package resolver

type ConfigV2PayloadResolver struct {
	user      string
	effective string
}

func NewConfigV2Payload(user, effective string) *ConfigV2PayloadResolver {
	return &ConfigV2PayloadResolver{user, effective}
}

func (r *ConfigV2PayloadResolver) User() string {
	return r.user
}

func (r *ConfigV2PayloadResolver) Effective() string {
	return r.effective
}
