package cmd

// JAID represents a JSON API ID.
//
// It implements the api2go MarshalIdentifier and UnmarshalIdentitier interface.
//
// When you embed a JSONAPI resource into a presenter, it will not render the
// ID into the JSON object when you perform a json.Marshal. Instead we use this
// to override the ID field of the resource with a JSON tag that will render.
//
// Embed this into a Presenter to render the ID. For example
//
//	type JobPresenter struct {
//		JAID
//		presenters.JobResource
//	}
type JAID struct {
	ID string `json:"id"`
}

func NewJAID(id string) JAID {
	return JAID{ID: id}
}

// GetID implements the api2go MarshalIdentifier interface.
func (jaid JAID) GetID() string {
	return jaid.ID
}

// SetID implements the api2go UnmarshalIdentitier interface.
func (jaid *JAID) SetID(value string) error {
	jaid.ID = value

	return nil
}
