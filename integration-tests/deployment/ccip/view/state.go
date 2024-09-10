package view

type CCIPSnapShot struct {
	Chains map[string]Chain `json:"chains,omitempty"`
}

func NewCCIPSnapShot() CCIPSnapShot {
	return CCIPSnapShot{
		Chains: make(map[string]Chain),
	}
}
