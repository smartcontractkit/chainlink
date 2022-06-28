package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidatedKeeperSpec(t *testing.T) {
	t.Parallel()

	type args struct {
		tomlString string
	}

	type want struct {
		id           int32
		contractAddr string
		fromAddr     string
		createdAt    time.Time
		updatedAt    time.Time
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "valid job spec",
			args: args{
				tomlString: `
                                                type                        = "keeper"
                                                name                        = "example keeper spec"
                                                contractAddress             = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
                                                fromAddress                 = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
                                                externalJobID               =  "123e4567-e89b-12d3-a456-426655440002"
					    `,
			},
			want: want{
				id:           0,
				contractAddr: "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba",
				fromAddr:     "0xa8037A20989AFcBC51798de9762b351D63ff462e",
				createdAt:    time.Time{},
				updatedAt:    time.Time{},
			},
			wantErr: false,
		},

		{
			name: "valid job spec with reordered fields",
			args: args{
				tomlString: `
						    type                        = "keeper"
						    name                        = "example keeper spec"
						    contractAddress             = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
						    fromAddress                 = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
						    evmChainID                  = 4
						    externalJobID               =  "123e4567-e89b-12d3-a456-426655440002"
					    `,
			},
			want: want{
				id:           0,
				contractAddr: "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba",
				fromAddr:     "0xa8037A20989AFcBC51798de9762b351D63ff462e",
				createdAt:    time.Time{},
				updatedAt:    time.Time{},
			},
			wantErr: false,
		},

		{
			name: "invalid job spec because of type",
			args: args{
				tomlString: `
						type            = "vrf"
						name            = "invalid keeper spec example 1"
						contractAddress = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
						fromAddress     = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
						evmChainID      = 4
						externalJobID   = "123e4567-e89b-12d3-a456-426655440002"
					`,
			},
			want:    want{},
			wantErr: true,
		},

		{
			name: "invalid job spec because observation source is passed as parameter (lowercase)",
			args: args{
				tomlString: `
                                                type              = "keeper"
                                                name              = "invalid keeper spec example 2"
                                                contractAddress   = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
                                                fromAddress       = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
                                                evmChainID        = 4
                                                externalJobID     = "123e4567-e89b-12d3-a456-426655440002"
                                                observationSource = "
                                                                    encode_check_upkeep_tx   [type=ethabiencode abi="checkUpkeep(uint256 id, address from)"
                                                                                              data="{\"id\":$(jobSpec.upkeepID),\"from\":$(jobSpec.fromAddress)}"]
                                                                    check_upkeep_tx          [type=ethcall
                                                                                              failEarly=false
                                                                                              extractRevertReason=false
                                                                                              evmChainID="$(jobSpec.evmChainID)"
                                                                                              contract="$(jobSpec.contractAddress)"
                                                                                              gas="$(jobSpec.checkUpkeepGasLimit)"
                                                                                              gasPrice="$(jobSpec.gasPrice)"
                                                                                              gasTipCap="$(jobSpec.gasTipCap)"
                                                                                              gasFeeCap="$(jobSpec.gasFeeCap)"
                                                                                              data="$(encode_check_upkeep_tx)"]
                                                                    decode_check_upkeep_tx   [type=ethabidecode
                                                                                              abi="bytes memory performData, uint256 maxLinkPayment,
                                                                                              uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
                                                                    encode_perform_upkeep_tx [type=ethabiencode
                                                                                              abi="performUpkeep(uint256 id, bytes calldata performData)"
                                                                                              data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
                                                                    perform_upkeep_tx        [type=ethtx
                                                                                              minConfirmations=2
                                                                                              to="$(jobSpec.contractAddress)"
                                                                                              from="[$(jobSpec.fromAddress)]"
                                                                                              evmChainID="$(jobSpec.evmChainID)"
                                                                                              data="$(encode_perform_upkeep_tx)"
                                                                                              gasLimit="$(jobSpec.performUpkeepGasLimit)"
                                                                                              txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
                                                                    encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx"
					      `,
			},
			want:    want{},
			wantErr: true,
		},

		{
			name: "invalid job spec because observation source is passed as parameter (uppercase)",
			args: args{
				tomlString: `
                                                type              = "keeper"
                                                name              = "invalid keeper spec example 2"
                                                contractAddress   = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
                                                fromAddress       = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
                                                evmChainID        = 4
                                                externalJobID     = "123e4567-e89b-12d3-a456-426655440002"
                                                ObservationSource = "
                                                                    encode_check_upkeep_tx   [type=ethabiencode abi="checkUpkeep(uint256 id, address from)"
                                                                                              data="{\"id\":$(jobSpec.upkeepID),\"from\":$(jobSpec.fromAddress)}"]
                                                                    check_upkeep_tx          [type=ethcall
                                                                                              failEarly=false
                                                                                              extractRevertReason=false
                                                                                              evmChainID="$(jobSpec.evmChainID)"
                                                                                              contract="$(jobSpec.contractAddress)"
                                                                                              gas="$(jobSpec.checkUpkeepGasLimit)"
                                                                                              gasPrice="$(jobSpec.gasPrice)"
                                                                                              gasTipCap="$(jobSpec.gasTipCap)"
                                                                                              gasFeeCap="$(jobSpec.gasFeeCap)"
                                                                                              data="$(encode_check_upkeep_tx)"]
                                                                    decode_check_upkeep_tx   [type=ethabidecode
                                                                                              abi="bytes memory performData, uint256 maxLinkPayment,
                                                                                              uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
                                                                    encode_perform_upkeep_tx [type=ethabiencode
                                                                                              abi="performUpkeep(uint256 id, bytes calldata performData)"
                                                                                              data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
                                                                    perform_upkeep_tx        [type=ethtx
                                                                                              minConfirmations=2
                                                                                              to="$(jobSpec.contractAddress)"
                                                                                              from="[$(jobSpec.fromAddress)]"
                                                                                              evmChainID="$(jobSpec.evmChainID)"
                                                                                              data="$(encode_perform_upkeep_tx)"
                                                                                              gasLimit="$(jobSpec.performUpkeepGasLimit)"
                                                                                              txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
                                                                    encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx"
					      `,
			},
			want:    want{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatedKeeperSpec(tt.args.tomlString)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tt.want.id, got.ID)
			require.Equal(t, tt.want.contractAddr, got.KeeperSpec.ContractAddress.Hex())
			require.Equal(t, tt.want.fromAddr, got.KeeperSpec.FromAddress.Hex())
			require.Equal(t, tt.want.createdAt, got.KeeperSpec.CreatedAt)
			require.Equal(t, tt.want.updatedAt, got.KeeperSpec.UpdatedAt)
		})
	}

}
