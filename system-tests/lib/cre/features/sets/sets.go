package sets

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	consensus_v1_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus/v1"
	consensus_v2_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus/v2"
	cron_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/cron"
	custom_compute_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/custom_compute"
	don_time_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/don_time"
	evm_v1_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm/v1"
	evm_v2_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm/v2"
	http_actions_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/http_action"
	http_trigger_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/http_trigger"
	log_event_trigger_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/log_event_trigger"
	mock_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/mock"
	read_contract_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/read_contract"
	solana_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/solana"
	vault_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/vault"
	web_api_target_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/web_api_target"
	web_api_trigger_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/web_api_trigger"
)

func New() cre.Features {
	return cre.NewFeatures(
		&consensus_v1_feature.Consensus{},
		&consensus_v2_feature.Consensus{},
		&cron_feature.Cron{},
		&custom_compute_feature.CustomCompute{},
		&don_time_feature.DONTime{},
		&evm_v1_feature.EVM{},
		&evm_v2_feature.EVM{},
		&http_actions_feature.HTTPAction{},
		&http_trigger_feature.HTTPTrigger{},
		&log_event_trigger_feature.LogEventTrigger{},
		&mock_feature.Mock{},
		&read_contract_feature.ReadContract{},
		&web_api_target_feature.WebAPITarget{},
		&web_api_trigger_feature.WebAPITrigger{},
		&vault_feature.Vault{},
		&solana_feature.Solana{},
	)
}
