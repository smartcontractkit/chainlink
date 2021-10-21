package presenters

// FeatureResource represents a Feature JSONAPI resource.
type FeatureResource struct {
	JAID
	Enabled bool `json:"enabled"`
}

// GetName implements the api2go EntityNamer interface
func (r FeatureResource) GetName() string {
	return "features"
}

// NewFeedsManagerResource constructs a new FeedsManagerResource.
func NewFeatureResource(name string, enabled bool) *FeatureResource {
	return &FeatureResource{
		JAID:    NewJAID(name),
		Enabled: enabled,
	}
}
