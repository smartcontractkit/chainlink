package metadata

import "fmt"

type OffchainConfig struct {
	DeltaGrace   string `json:"deltaGrace"`
	DeltaInitial string `json:"deltaInitial"`
}
type DonMetadata struct {
	ID                  string `json:"id"`
	ConfiguratorAddress string `json:"configuratorAddress"`
	OffchainConfig      OffchainConfig
	Streams             []int `json:"streams"`
}

// DataStreamsMetadata is a struct that can be used as a default metadata type.
// This is just a placeholder for the actual metadata structure to show how it can be used.
type DataStreamsMetadata struct {
	DONs []DonMetadata
}

// DefaultMetadata implements the Cloneable interface
func (d DataStreamsMetadata) Clone() DataStreamsMetadata { return d }

func (d DataStreamsMetadata) GetDonByID(id string) (DonMetadata, error) {
	for _, don := range d.DONs {
		if don.ID == id {
			return don, nil
		}
	}
	return DonMetadata{}, fmt.Errorf("don with id %s not found", id)
}
