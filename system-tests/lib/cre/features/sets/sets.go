package sets

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	aptos_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/aptos"
	consensus_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus/v2"
	cron_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/cron"
	don_time_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/don_time"
	evm_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm/v2"
	http_actions_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/http_action"
	http_trigger_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/http_trigger"
	solana_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/solana/v2"
	vault_feature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/vault"
)

func New() cre.Features {
	return cre.NewFeatures(
		&consensus_feature.Consensus{},
		&cron_feature.Cron{},
		&don_time_feature.DONTime{},
		&evm_feature.EVM{},
		&http_actions_feature.HTTPAction{},
		&http_trigger_feature.HTTPTrigger{},
		&aptos_feature.Aptos{},
		&solana_feature.Solana{},
		&vault_feature.Vault{},
	)
}
