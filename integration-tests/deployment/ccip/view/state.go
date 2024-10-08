package view

type CCIPView struct {
	Chains        map[string]ChainView `json:"chains,omitempty"`
	NodeOperators NopsView             `json:"nodeOperators,omitempty"`
}

func NewCCIPView() CCIPView {
	return CCIPView{
		Chains: make(map[string]ChainView),
	}
}
