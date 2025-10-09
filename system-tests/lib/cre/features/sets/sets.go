package sets

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	consensus_v1_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus/v1"
	don_time_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/don_time"
	evm_v1_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm/v1"
	http_actions_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/http_action"
	http_trigger_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/http_trigger"
	vault_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/vault"
	web_api_target_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/web_api_target"
	web_api_trigger_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/web_api_trigger"
)

func New() cre.Features {
	return cre.NewFeatures(
		&consensus_v1_feature.Consensus{},
		&don_time_feature.DONTime{},
		&evm_v1_feature.EVM{},
		&http_actions_feature.HTTPAction{},
		&http_trigger_feature.HTTPTrigger{},
		&web_api_target_feature.WebAPITarget{},
		&web_api_trigger_feature.WebAPITrigger{},
		&vault_feature.Vault{},
	)
}
