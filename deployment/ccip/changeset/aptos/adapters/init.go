package adapters

import (
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
)

func init() {
	curseRegistry := fastcurse.GetCurseRegistry()
	curseRegistry.RegisterNewCurse(fastcurse.CurseRegistryInput{
		CursingFamily:       chainsel.FamilyAptos,
		CursingVersion:      semver.MustParse("1.6.0"),
		CurseAdapter:        NewCurseAdapter(),
		CurseSubjectAdapter: NewCurseAdapter(),
	})
}
