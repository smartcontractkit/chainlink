package view

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/os"
)

type CCIPView struct {
	Chains map[string]ChainView `json:"chains,omitempty"`
}

func NewCCIPView() CCIPView {
	return CCIPView{
		Chains: make(map[string]ChainView),
	}
}

func SaveView(view CCIPView) error {
	b, err := json.MarshalIndent(view, "", "	")
	if err != nil {
		return err
	}
	filepath := fmt.Sprintf("ccip_view_%s.json", time.Now().Format("2006-01-02T15:04:05"))
	return os.WriteFile(filepath, b, 0644)
}
