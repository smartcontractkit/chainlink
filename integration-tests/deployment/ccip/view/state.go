package view

import "encoding/json"

type CCIPView struct {
	Chains map[string]ChainView `json:"chains,omitempty"`
	Nops   map[string]NopView   `json:"nops,omitempty"`
}

func (v CCIPView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}
