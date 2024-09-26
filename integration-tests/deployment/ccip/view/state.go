package view

type CCIPView struct {
	Chains map[string]ChainView `json:"chains,omitempty"`
}

func NewCCIPView() CCIPView {
	return CCIPView{
		Chains: make(map[string]ChainView),
	}
}
