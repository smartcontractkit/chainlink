// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

var StakingEventsMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"AlertRaised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"feedOperators\",\"type\":\"address[]\"}],\"name\":\"FeedOperatorsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"MaxCommunityStakeAmountIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"MaxOperatorStakeAmountIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newMerkleRoot\",\"type\":\"bytes32\"}],\"name\":\"MerkleRootChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"delegationReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Migrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"MigrationTargetAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"MigrationTargetProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"PoolConcluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"PoolOpened\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxPoolSize\",\"type\":\"uint256\"}],\"name\":\"PoolSizeIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountAdded\",\"type\":\"uint256\"}],\"name\":\"RewardAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"name\":\"RewardInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"RewardRateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"operator\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"slashedBaseRewards\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"slashedDelegatedRewards\",\"type\":\"uint256[]\"}],\"name\":\"RewardSlashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newStake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"delegationReward\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"emitAlertRaised\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"feedOperators\",\"type\":\"address[]\"}],\"name\":\"emitFeedOperatorsSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"emitMaxCommunityStakeAmountIncreased\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"emitMaxOperatorStakeAmountIncreased\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newMerkleRoot\",\"type\":\"bytes32\"}],\"name\":\"emitMerkleRootChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseReward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationReward\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"emitMigrated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"emitMigrationTargetAccepted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"emitMigrationTargetProposed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"emitOperatorAdded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitOperatorRemoved\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitPoolConcluded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitPoolOpened\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxPoolSize\",\"type\":\"uint256\"}],\"name\":\"emitPoolSizeIncreased\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountAdded\",\"type\":\"uint256\"}],\"name\":\"emitRewardAdded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"name\":\"emitRewardInitialized\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"emitRewardRateChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operator\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"slashedBaseRewards\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"slashedDelegatedRewards\",\"type\":\"uint256[]\"}],\"name\":\"emitRewardSlashed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitRewardWithdrawn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"}],\"name\":\"emitStaked\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitUnpaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseReward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationReward\",\"type\":\"uint256\"}],\"name\":\"emitUnstaked\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610f7b806100206000396000f3fe608060405234801561001057600080fd5b50600436106101985760003560e01c80639ec3ce4b116100e3578063b019b4e81161008c578063f7420bc211610066578063f7420bc214610318578063fa2d7fd01461032b578063fe60093f1461033e57600080fd5b8063b019b4e8146102df578063e0275fae146102f2578063eb6289de1461030557600080fd5b8063ab8711bc116100bd578063ab8711bc146102b1578063add49a96146102c4578063aeac9600146102cc57600080fd5b80639ec3ce4b146102785780639f676e9f1461028b578063a5f3d6831461029e57600080fd5b80637be5c756116101455780639a8d2df71161011f5780639a8d2df71461023f5780639b022d73146102525780639e81ace31461026557600080fd5b80637be5c756146102065780637e31a64b146102195780639652ab7b1461022c57600080fd5b80632752e639116101765780632752e639146101d85780633e8f1e05146101eb5780635d21e09a146101f357600080fd5b8063086c1c4a1461019d5780631351da48146101b257806313c664e8146101c5575b600080fd5b6101b06101ab366004610b25565b610351565b005b6101b06101c0366004610a73565b6103b6565b6101b06101d3366004610d04565b610403565b6101b06101e6366004610b5e565b610433565b6101b0610479565b6101b0610201366004610d04565b6104a4565b6101b0610214366004610a73565b6104d4565b6101b0610227366004610d04565b61051a565b6101b061023a366004610a73565b61054a565b6101b061024d366004610a73565b610590565b6101b0610260366004610d04565b6105d6565b6101b0610273366004610c7c565b610606565b6101b0610286366004610a73565b610645565b6101b0610299366004610d1d565b61068b565b6101b06102ac366004610c3f565b6106d0565b6101b06102bf366004610af2565b6106ff565b6101b0610753565b6101b06102da366004610d04565b61077e565b6101b06102ed366004610a95565b6107ae565b6101b0610300366004610af2565b61080c565b6101b0610313366004610ac8565b610860565b6101b0610326366004610a95565b6108b3565b6101b0610339366004610d04565b610911565b6101b061034c366004610d04565b610941565b6040805173ffffffffffffffffffffffffffffffffffffffff8616815260208101859052908101839052606081018290527f204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00906080015b60405180910390a150505050565b60405173ffffffffffffffffffffffffffffffffffffffff821681527fac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d906020015b60405180910390a150565b6040518181527f816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52906020016103f8565b7f667838b33bdc898470de09e0e746990f2adc11b965b7fe6828e502ebc39e0434858585858560405161046a959493929190610dd0565b60405180910390a15050505050565b6040517fded6ebf04e261e1eb2f3e3b268a2e6aee5b478c15b341eba5cf18b9bc80c2e6390600090a1565b6040518181527fde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d906020016103f8565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020016103f8565b6040518181527fb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df7906020016103f8565b60405173ffffffffffffffffffffffffffffffffffffffff821681527ffa33c052bbee754f3c0482a89962daffe749191fa33c696a61e947fbfd68bd84906020016103f8565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f5c74c441be501340b2713817a6c6975e6f3d4a4ae39fa1ac0bf75d3c54a0cad3906020016103f8565b6040518181527f7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c0906020016103f8565b7e635ea9da6e262e92bb713d71840af7c567807ff35bf73e927490c61283248083838360405161063893929190610e89565b60405180910390a1505050565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020016103f8565b6040805185815260208101859052908101839052606081018290527f125fc8494f786b470e3c39d0932a62e9e09e291ebd81ea19c57604f6d2b1d167906080016103a8565b7f40aed8e423b39a56b445ae160f4c071fc2cfb48ee0b6dcd5ffeb6bc5b18d10d0816040516103f89190610e76565b6040805173ffffffffffffffffffffffffffffffffffffffff85168152602081018490529081018290527f1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee9090606001610638565b6040517ff7d0e0f15586495da8c687328ead30fb829d9da55538cb0ef73dd229e517cdb890600090a1565b6040518181527f1b930366dfeaa7eb3b325021e4ae81e36527063452ee55b86c95f85b36f4c31c906020016103f8565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6040805173ffffffffffffffffffffffffffffffffffffffff85168152602081018490529081018290527fd2720e8f454493f612cc97499fe8cbce7fa4d4c18d346fe7104e9042df1c1edd90606001610638565b6040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f2360404a74478febece1a14f11275f22ada88d19ef96f7d785913010bfff4479910160405180910390a15050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b6040518181527f150a6ec0e6f4e9ddcaaaa1674f157d91165a42d60653016f87a9fc870a39f050906020016103f8565b6040518181527f1e3be2efa25bca5bff2215c7b30b31086e703d6aa7d9b9a1f8ba62c5291219ad906020016103f8565b803573ffffffffffffffffffffffffffffffffffffffff8116811461099557600080fd5b919050565b600082601f8301126109ab57600080fd5b813560206109c06109bb83610f1b565b610ecc565b80838252828201915082860187848660051b89010111156109e057600080fd5b60005b85811015610a06576109f482610971565b845292840192908401906001016109e3565b5090979650505050505050565b600082601f830112610a2457600080fd5b81356020610a346109bb83610f1b565b80838252828201915082860187848660051b8901011115610a5457600080fd5b60005b85811015610a0657813584529284019290840190600101610a57565b600060208284031215610a8557600080fd5b610a8e82610971565b9392505050565b60008060408385031215610aa857600080fd5b610ab183610971565b9150610abf60208401610971565b90509250929050565b60008060408385031215610adb57600080fd5b610ae483610971565b946020939093013593505050565b600080600060608486031215610b0757600080fd5b610b1084610971565b95602085013595506040909401359392505050565b60008060008060808587031215610b3b57600080fd5b610b4485610971565b966020860135965060408601359560600135945092505050565b600080600080600060a08688031215610b7657600080fd5b610b7f86610971565b945060208087013594506040870135935060608701359250608087013567ffffffffffffffff80821115610bb257600080fd5b818901915089601f830112610bc657600080fd5b813581811115610bd857610bd8610f3f565b610c08847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610ecc565b91508082528a84828501011115610c1e57600080fd5b80848401858401376000848284010152508093505050509295509295909350565b600060208284031215610c5157600080fd5b813567ffffffffffffffff811115610c6857600080fd5b610c748482850161099a565b949350505050565b600080600060608486031215610c9157600080fd5b833567ffffffffffffffff80821115610ca957600080fd5b610cb58783880161099a565b94506020860135915080821115610ccb57600080fd5b610cd787838801610a13565b93506040860135915080821115610ced57600080fd5b50610cfa86828701610a13565b9150509250925092565b600060208284031215610d1657600080fd5b5035919050565b60008060008060808587031215610d3357600080fd5b5050823594602084013594506040840135936060013592509050565b600081518084526020808501945080840160005b83811015610d9557815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610d63565b509495945050505050565b600081518084526020808501945080840160005b83811015610d9557815187529582019590820190600101610db4565b73ffffffffffffffffffffffffffffffffffffffff8616815260006020868184015285604084015284606084015260a0608084015283518060a085015260005b81811015610e2c5785810183015185820160c001528201610e10565b81811115610e3e57600060c083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160c001979650505050505050565b602081526000610a8e6020830184610d4f565b606081526000610e9c6060830186610d4f565b8281036020840152610eae8186610da0565b90508281036040840152610ec28185610da0565b9695505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610f1357610f13610f3f565b604052919050565b600067ffffffffffffffff821115610f3557610f35610f3f565b5060051b60200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var StakingEventsMockABI = StakingEventsMockMetaData.ABI

var StakingEventsMockBin = StakingEventsMockMetaData.Bin

func DeployStakingEventsMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingEventsMock, error) {
	parsed, err := StakingEventsMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingEventsMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingEventsMock{StakingEventsMockCaller: StakingEventsMockCaller{contract: contract}, StakingEventsMockTransactor: StakingEventsMockTransactor{contract: contract}, StakingEventsMockFilterer: StakingEventsMockFilterer{contract: contract}}, nil
}

type StakingEventsMock struct {
	address common.Address
	abi     abi.ABI
	StakingEventsMockCaller
	StakingEventsMockTransactor
	StakingEventsMockFilterer
}

type StakingEventsMockCaller struct {
	contract *bind.BoundContract
}

type StakingEventsMockTransactor struct {
	contract *bind.BoundContract
}

type StakingEventsMockFilterer struct {
	contract *bind.BoundContract
}

type StakingEventsMockSession struct {
	Contract     *StakingEventsMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type StakingEventsMockCallerSession struct {
	Contract *StakingEventsMockCaller
	CallOpts bind.CallOpts
}

type StakingEventsMockTransactorSession struct {
	Contract     *StakingEventsMockTransactor
	TransactOpts bind.TransactOpts
}

type StakingEventsMockRaw struct {
	Contract *StakingEventsMock
}

type StakingEventsMockCallerRaw struct {
	Contract *StakingEventsMockCaller
}

type StakingEventsMockTransactorRaw struct {
	Contract *StakingEventsMockTransactor
}

func NewStakingEventsMock(address common.Address, backend bind.ContractBackend) (*StakingEventsMock, error) {
	abi, err := abi.JSON(strings.NewReader(StakingEventsMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindStakingEventsMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingEventsMock{address: address, abi: abi, StakingEventsMockCaller: StakingEventsMockCaller{contract: contract}, StakingEventsMockTransactor: StakingEventsMockTransactor{contract: contract}, StakingEventsMockFilterer: StakingEventsMockFilterer{contract: contract}}, nil
}

func NewStakingEventsMockCaller(address common.Address, caller bind.ContractCaller) (*StakingEventsMockCaller, error) {
	contract, err := bindStakingEventsMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockCaller{contract: contract}, nil
}

func NewStakingEventsMockTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingEventsMockTransactor, error) {
	contract, err := bindStakingEventsMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockTransactor{contract: contract}, nil
}

func NewStakingEventsMockFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingEventsMockFilterer, error) {
	contract, err := bindStakingEventsMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockFilterer{contract: contract}, nil
}

func bindStakingEventsMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingEventsMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_StakingEventsMock *StakingEventsMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingEventsMock.Contract.StakingEventsMockCaller.contract.Call(opts, result, method, params...)
}

func (_StakingEventsMock *StakingEventsMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.StakingEventsMockTransactor.contract.Transfer(opts)
}

func (_StakingEventsMock *StakingEventsMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.StakingEventsMockTransactor.contract.Transact(opts, method, params...)
}

func (_StakingEventsMock *StakingEventsMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingEventsMock.Contract.contract.Call(opts, result, method, params...)
}

func (_StakingEventsMock *StakingEventsMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.contract.Transfer(opts)
}

func (_StakingEventsMock *StakingEventsMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.contract.Transact(opts, method, params...)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitAlertRaised(opts *bind.TransactOpts, alerter common.Address, roundId *big.Int, rewardAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitAlertRaised", alerter, roundId, rewardAmount)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitAlertRaised(alerter common.Address, roundId *big.Int, rewardAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitAlertRaised(&_StakingEventsMock.TransactOpts, alerter, roundId, rewardAmount)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitAlertRaised(alerter common.Address, roundId *big.Int, rewardAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitAlertRaised(&_StakingEventsMock.TransactOpts, alerter, roundId, rewardAmount)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitFeedOperatorsSet(opts *bind.TransactOpts, feedOperators []common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitFeedOperatorsSet", feedOperators)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitFeedOperatorsSet(feedOperators []common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitFeedOperatorsSet(&_StakingEventsMock.TransactOpts, feedOperators)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitFeedOperatorsSet(feedOperators []common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitFeedOperatorsSet(&_StakingEventsMock.TransactOpts, feedOperators)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitMaxCommunityStakeAmountIncreased(opts *bind.TransactOpts, maxStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitMaxCommunityStakeAmountIncreased", maxStakeAmount)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitMaxCommunityStakeAmountIncreased(maxStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMaxCommunityStakeAmountIncreased(&_StakingEventsMock.TransactOpts, maxStakeAmount)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitMaxCommunityStakeAmountIncreased(maxStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMaxCommunityStakeAmountIncreased(&_StakingEventsMock.TransactOpts, maxStakeAmount)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitMaxOperatorStakeAmountIncreased(opts *bind.TransactOpts, maxStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitMaxOperatorStakeAmountIncreased", maxStakeAmount)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitMaxOperatorStakeAmountIncreased(maxStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMaxOperatorStakeAmountIncreased(&_StakingEventsMock.TransactOpts, maxStakeAmount)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitMaxOperatorStakeAmountIncreased(maxStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMaxOperatorStakeAmountIncreased(&_StakingEventsMock.TransactOpts, maxStakeAmount)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitMerkleRootChanged(opts *bind.TransactOpts, newMerkleRoot [32]byte) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitMerkleRootChanged", newMerkleRoot)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitMerkleRootChanged(newMerkleRoot [32]byte) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMerkleRootChanged(&_StakingEventsMock.TransactOpts, newMerkleRoot)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitMerkleRootChanged(newMerkleRoot [32]byte) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMerkleRootChanged(&_StakingEventsMock.TransactOpts, newMerkleRoot)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitMigrated(opts *bind.TransactOpts, staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int, data []byte) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitMigrated", staker, principal, baseReward, delegationReward, data)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitMigrated(staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int, data []byte) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMigrated(&_StakingEventsMock.TransactOpts, staker, principal, baseReward, delegationReward, data)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitMigrated(staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int, data []byte) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMigrated(&_StakingEventsMock.TransactOpts, staker, principal, baseReward, delegationReward, data)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitMigrationTargetAccepted(opts *bind.TransactOpts, migrationTarget common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitMigrationTargetAccepted", migrationTarget)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitMigrationTargetAccepted(migrationTarget common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMigrationTargetAccepted(&_StakingEventsMock.TransactOpts, migrationTarget)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitMigrationTargetAccepted(migrationTarget common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMigrationTargetAccepted(&_StakingEventsMock.TransactOpts, migrationTarget)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitMigrationTargetProposed(opts *bind.TransactOpts, migrationTarget common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitMigrationTargetProposed", migrationTarget)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitMigrationTargetProposed(migrationTarget common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMigrationTargetProposed(&_StakingEventsMock.TransactOpts, migrationTarget)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitMigrationTargetProposed(migrationTarget common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitMigrationTargetProposed(&_StakingEventsMock.TransactOpts, migrationTarget)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitOperatorAdded(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitOperatorAdded", operator)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitOperatorAdded(operator common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOperatorAdded(&_StakingEventsMock.TransactOpts, operator)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitOperatorAdded(operator common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOperatorAdded(&_StakingEventsMock.TransactOpts, operator)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitOperatorRemoved(opts *bind.TransactOpts, operator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitOperatorRemoved", operator, amount)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitOperatorRemoved(operator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOperatorRemoved(&_StakingEventsMock.TransactOpts, operator, amount)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitOperatorRemoved(operator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOperatorRemoved(&_StakingEventsMock.TransactOpts, operator, amount)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOwnershipTransferRequested(&_StakingEventsMock.TransactOpts, from, to)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOwnershipTransferRequested(&_StakingEventsMock.TransactOpts, from, to)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOwnershipTransferred(&_StakingEventsMock.TransactOpts, from, to)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitOwnershipTransferred(&_StakingEventsMock.TransactOpts, from, to)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitPaused", account)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPaused(&_StakingEventsMock.TransactOpts, account)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPaused(&_StakingEventsMock.TransactOpts, account)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitPoolConcluded(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitPoolConcluded")
}

func (_StakingEventsMock *StakingEventsMockSession) EmitPoolConcluded() (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPoolConcluded(&_StakingEventsMock.TransactOpts)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitPoolConcluded() (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPoolConcluded(&_StakingEventsMock.TransactOpts)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitPoolOpened(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitPoolOpened")
}

func (_StakingEventsMock *StakingEventsMockSession) EmitPoolOpened() (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPoolOpened(&_StakingEventsMock.TransactOpts)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitPoolOpened() (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPoolOpened(&_StakingEventsMock.TransactOpts)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitPoolSizeIncreased(opts *bind.TransactOpts, maxPoolSize *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitPoolSizeIncreased", maxPoolSize)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitPoolSizeIncreased(maxPoolSize *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPoolSizeIncreased(&_StakingEventsMock.TransactOpts, maxPoolSize)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitPoolSizeIncreased(maxPoolSize *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitPoolSizeIncreased(&_StakingEventsMock.TransactOpts, maxPoolSize)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitRewardAdded(opts *bind.TransactOpts, amountAdded *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitRewardAdded", amountAdded)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitRewardAdded(amountAdded *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardAdded(&_StakingEventsMock.TransactOpts, amountAdded)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitRewardAdded(amountAdded *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardAdded(&_StakingEventsMock.TransactOpts, amountAdded)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitRewardInitialized(opts *bind.TransactOpts, rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitRewardInitialized", rate, available, startTimestamp, endTimestamp)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitRewardInitialized(rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardInitialized(&_StakingEventsMock.TransactOpts, rate, available, startTimestamp, endTimestamp)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitRewardInitialized(rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardInitialized(&_StakingEventsMock.TransactOpts, rate, available, startTimestamp, endTimestamp)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitRewardRateChanged(opts *bind.TransactOpts, rate *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitRewardRateChanged", rate)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitRewardRateChanged(rate *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardRateChanged(&_StakingEventsMock.TransactOpts, rate)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitRewardRateChanged(rate *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardRateChanged(&_StakingEventsMock.TransactOpts, rate)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitRewardSlashed(opts *bind.TransactOpts, operator []common.Address, slashedBaseRewards []*big.Int, slashedDelegatedRewards []*big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitRewardSlashed", operator, slashedBaseRewards, slashedDelegatedRewards)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitRewardSlashed(operator []common.Address, slashedBaseRewards []*big.Int, slashedDelegatedRewards []*big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardSlashed(&_StakingEventsMock.TransactOpts, operator, slashedBaseRewards, slashedDelegatedRewards)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitRewardSlashed(operator []common.Address, slashedBaseRewards []*big.Int, slashedDelegatedRewards []*big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardSlashed(&_StakingEventsMock.TransactOpts, operator, slashedBaseRewards, slashedDelegatedRewards)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitRewardWithdrawn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitRewardWithdrawn", amount)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitRewardWithdrawn(amount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardWithdrawn(&_StakingEventsMock.TransactOpts, amount)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitRewardWithdrawn(amount *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitRewardWithdrawn(&_StakingEventsMock.TransactOpts, amount)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitStaked(opts *bind.TransactOpts, staker common.Address, newStake *big.Int, totalStake *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitStaked", staker, newStake, totalStake)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitStaked(staker common.Address, newStake *big.Int, totalStake *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitStaked(&_StakingEventsMock.TransactOpts, staker, newStake, totalStake)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitStaked(staker common.Address, newStake *big.Int, totalStake *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitStaked(&_StakingEventsMock.TransactOpts, staker, newStake, totalStake)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitUnpaused", account)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitUnpaused(&_StakingEventsMock.TransactOpts, account)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitUnpaused(&_StakingEventsMock.TransactOpts, account)
}

func (_StakingEventsMock *StakingEventsMockTransactor) EmitUnstaked(opts *bind.TransactOpts, staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.contract.Transact(opts, "emitUnstaked", staker, principal, baseReward, delegationReward)
}

func (_StakingEventsMock *StakingEventsMockSession) EmitUnstaked(staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitUnstaked(&_StakingEventsMock.TransactOpts, staker, principal, baseReward, delegationReward)
}

func (_StakingEventsMock *StakingEventsMockTransactorSession) EmitUnstaked(staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int) (*types.Transaction, error) {
	return _StakingEventsMock.Contract.EmitUnstaked(&_StakingEventsMock.TransactOpts, staker, principal, baseReward, delegationReward)
}

type StakingEventsMockAlertRaisedIterator struct {
	Event *StakingEventsMockAlertRaised

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockAlertRaisedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockAlertRaised)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockAlertRaised)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockAlertRaisedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockAlertRaisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockAlertRaised struct {
	Alerter      common.Address
	RoundId      *big.Int
	RewardAmount *big.Int
	Raw          types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterAlertRaised(opts *bind.FilterOpts) (*StakingEventsMockAlertRaisedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "AlertRaised")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockAlertRaisedIterator{contract: _StakingEventsMock.contract, event: "AlertRaised", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchAlertRaised(opts *bind.WatchOpts, sink chan<- *StakingEventsMockAlertRaised) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "AlertRaised")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockAlertRaised)
				if err := _StakingEventsMock.contract.UnpackLog(event, "AlertRaised", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseAlertRaised(log types.Log) (*StakingEventsMockAlertRaised, error) {
	event := new(StakingEventsMockAlertRaised)
	if err := _StakingEventsMock.contract.UnpackLog(event, "AlertRaised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockFeedOperatorsSetIterator struct {
	Event *StakingEventsMockFeedOperatorsSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockFeedOperatorsSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockFeedOperatorsSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockFeedOperatorsSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockFeedOperatorsSetIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockFeedOperatorsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockFeedOperatorsSet struct {
	FeedOperators []common.Address
	Raw           types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterFeedOperatorsSet(opts *bind.FilterOpts) (*StakingEventsMockFeedOperatorsSetIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "FeedOperatorsSet")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockFeedOperatorsSetIterator{contract: _StakingEventsMock.contract, event: "FeedOperatorsSet", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchFeedOperatorsSet(opts *bind.WatchOpts, sink chan<- *StakingEventsMockFeedOperatorsSet) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "FeedOperatorsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockFeedOperatorsSet)
				if err := _StakingEventsMock.contract.UnpackLog(event, "FeedOperatorsSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseFeedOperatorsSet(log types.Log) (*StakingEventsMockFeedOperatorsSet, error) {
	event := new(StakingEventsMockFeedOperatorsSet)
	if err := _StakingEventsMock.contract.UnpackLog(event, "FeedOperatorsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockMaxCommunityStakeAmountIncreasedIterator struct {
	Event *StakingEventsMockMaxCommunityStakeAmountIncreased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockMaxCommunityStakeAmountIncreasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockMaxCommunityStakeAmountIncreased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockMaxCommunityStakeAmountIncreased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockMaxCommunityStakeAmountIncreasedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockMaxCommunityStakeAmountIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockMaxCommunityStakeAmountIncreased struct {
	MaxStakeAmount *big.Int
	Raw            types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterMaxCommunityStakeAmountIncreased(opts *bind.FilterOpts) (*StakingEventsMockMaxCommunityStakeAmountIncreasedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "MaxCommunityStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockMaxCommunityStakeAmountIncreasedIterator{contract: _StakingEventsMock.contract, event: "MaxCommunityStakeAmountIncreased", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchMaxCommunityStakeAmountIncreased(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMaxCommunityStakeAmountIncreased) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "MaxCommunityStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockMaxCommunityStakeAmountIncreased)
				if err := _StakingEventsMock.contract.UnpackLog(event, "MaxCommunityStakeAmountIncreased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseMaxCommunityStakeAmountIncreased(log types.Log) (*StakingEventsMockMaxCommunityStakeAmountIncreased, error) {
	event := new(StakingEventsMockMaxCommunityStakeAmountIncreased)
	if err := _StakingEventsMock.contract.UnpackLog(event, "MaxCommunityStakeAmountIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockMaxOperatorStakeAmountIncreasedIterator struct {
	Event *StakingEventsMockMaxOperatorStakeAmountIncreased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockMaxOperatorStakeAmountIncreasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockMaxOperatorStakeAmountIncreased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockMaxOperatorStakeAmountIncreased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockMaxOperatorStakeAmountIncreasedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockMaxOperatorStakeAmountIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockMaxOperatorStakeAmountIncreased struct {
	MaxStakeAmount *big.Int
	Raw            types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterMaxOperatorStakeAmountIncreased(opts *bind.FilterOpts) (*StakingEventsMockMaxOperatorStakeAmountIncreasedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "MaxOperatorStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockMaxOperatorStakeAmountIncreasedIterator{contract: _StakingEventsMock.contract, event: "MaxOperatorStakeAmountIncreased", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchMaxOperatorStakeAmountIncreased(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMaxOperatorStakeAmountIncreased) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "MaxOperatorStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockMaxOperatorStakeAmountIncreased)
				if err := _StakingEventsMock.contract.UnpackLog(event, "MaxOperatorStakeAmountIncreased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseMaxOperatorStakeAmountIncreased(log types.Log) (*StakingEventsMockMaxOperatorStakeAmountIncreased, error) {
	event := new(StakingEventsMockMaxOperatorStakeAmountIncreased)
	if err := _StakingEventsMock.contract.UnpackLog(event, "MaxOperatorStakeAmountIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockMerkleRootChangedIterator struct {
	Event *StakingEventsMockMerkleRootChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockMerkleRootChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockMerkleRootChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockMerkleRootChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockMerkleRootChangedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockMerkleRootChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockMerkleRootChanged struct {
	NewMerkleRoot [32]byte
	Raw           types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterMerkleRootChanged(opts *bind.FilterOpts) (*StakingEventsMockMerkleRootChangedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "MerkleRootChanged")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockMerkleRootChangedIterator{contract: _StakingEventsMock.contract, event: "MerkleRootChanged", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchMerkleRootChanged(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMerkleRootChanged) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "MerkleRootChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockMerkleRootChanged)
				if err := _StakingEventsMock.contract.UnpackLog(event, "MerkleRootChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseMerkleRootChanged(log types.Log) (*StakingEventsMockMerkleRootChanged, error) {
	event := new(StakingEventsMockMerkleRootChanged)
	if err := _StakingEventsMock.contract.UnpackLog(event, "MerkleRootChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockMigratedIterator struct {
	Event *StakingEventsMockMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockMigratedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockMigrated struct {
	Staker           common.Address
	Principal        *big.Int
	BaseReward       *big.Int
	DelegationReward *big.Int
	Data             []byte
	Raw              types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterMigrated(opts *bind.FilterOpts) (*StakingEventsMockMigratedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "Migrated")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockMigratedIterator{contract: _StakingEventsMock.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMigrated) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "Migrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockMigrated)
				if err := _StakingEventsMock.contract.UnpackLog(event, "Migrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseMigrated(log types.Log) (*StakingEventsMockMigrated, error) {
	event := new(StakingEventsMockMigrated)
	if err := _StakingEventsMock.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockMigrationTargetAcceptedIterator struct {
	Event *StakingEventsMockMigrationTargetAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockMigrationTargetAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockMigrationTargetAccepted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockMigrationTargetAccepted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockMigrationTargetAcceptedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockMigrationTargetAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockMigrationTargetAccepted struct {
	MigrationTarget common.Address
	Raw             types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterMigrationTargetAccepted(opts *bind.FilterOpts) (*StakingEventsMockMigrationTargetAcceptedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "MigrationTargetAccepted")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockMigrationTargetAcceptedIterator{contract: _StakingEventsMock.contract, event: "MigrationTargetAccepted", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchMigrationTargetAccepted(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMigrationTargetAccepted) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "MigrationTargetAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockMigrationTargetAccepted)
				if err := _StakingEventsMock.contract.UnpackLog(event, "MigrationTargetAccepted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseMigrationTargetAccepted(log types.Log) (*StakingEventsMockMigrationTargetAccepted, error) {
	event := new(StakingEventsMockMigrationTargetAccepted)
	if err := _StakingEventsMock.contract.UnpackLog(event, "MigrationTargetAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockMigrationTargetProposedIterator struct {
	Event *StakingEventsMockMigrationTargetProposed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockMigrationTargetProposedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockMigrationTargetProposed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockMigrationTargetProposed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockMigrationTargetProposedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockMigrationTargetProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockMigrationTargetProposed struct {
	MigrationTarget common.Address
	Raw             types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterMigrationTargetProposed(opts *bind.FilterOpts) (*StakingEventsMockMigrationTargetProposedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "MigrationTargetProposed")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockMigrationTargetProposedIterator{contract: _StakingEventsMock.contract, event: "MigrationTargetProposed", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchMigrationTargetProposed(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMigrationTargetProposed) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "MigrationTargetProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockMigrationTargetProposed)
				if err := _StakingEventsMock.contract.UnpackLog(event, "MigrationTargetProposed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseMigrationTargetProposed(log types.Log) (*StakingEventsMockMigrationTargetProposed, error) {
	event := new(StakingEventsMockMigrationTargetProposed)
	if err := _StakingEventsMock.contract.UnpackLog(event, "MigrationTargetProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockOperatorAddedIterator struct {
	Event *StakingEventsMockOperatorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockOperatorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockOperatorAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockOperatorAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockOperatorAddedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockOperatorAdded struct {
	Operator common.Address
	Raw      types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*StakingEventsMockOperatorAddedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockOperatorAddedIterator{contract: _StakingEventsMock.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockOperatorAdded)
				if err := _StakingEventsMock.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseOperatorAdded(log types.Log) (*StakingEventsMockOperatorAdded, error) {
	event := new(StakingEventsMockOperatorAdded)
	if err := _StakingEventsMock.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockOperatorRemovedIterator struct {
	Event *StakingEventsMockOperatorRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockOperatorRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockOperatorRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockOperatorRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockOperatorRemovedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockOperatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockOperatorRemoved struct {
	Operator common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterOperatorRemoved(opts *bind.FilterOpts) (*StakingEventsMockOperatorRemovedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "OperatorRemoved")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockOperatorRemovedIterator{contract: _StakingEventsMock.contract, event: "OperatorRemoved", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchOperatorRemoved(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOperatorRemoved) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "OperatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockOperatorRemoved)
				if err := _StakingEventsMock.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseOperatorRemoved(log types.Log) (*StakingEventsMockOperatorRemoved, error) {
	event := new(StakingEventsMockOperatorRemoved)
	if err := _StakingEventsMock.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockOwnershipTransferRequestedIterator struct {
	Event *StakingEventsMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingEventsMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockOwnershipTransferRequestedIterator{contract: _StakingEventsMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockOwnershipTransferRequested)
				if err := _StakingEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*StakingEventsMockOwnershipTransferRequested, error) {
	event := new(StakingEventsMockOwnershipTransferRequested)
	if err := _StakingEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockOwnershipTransferredIterator struct {
	Event *StakingEventsMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingEventsMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockOwnershipTransferredIterator{contract: _StakingEventsMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockOwnershipTransferred)
				if err := _StakingEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseOwnershipTransferred(log types.Log) (*StakingEventsMockOwnershipTransferred, error) {
	event := new(StakingEventsMockOwnershipTransferred)
	if err := _StakingEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockPausedIterator struct {
	Event *StakingEventsMockPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockPausedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterPaused(opts *bind.FilterOpts) (*StakingEventsMockPausedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockPausedIterator{contract: _StakingEventsMock.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPaused) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockPaused)
				if err := _StakingEventsMock.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParsePaused(log types.Log) (*StakingEventsMockPaused, error) {
	event := new(StakingEventsMockPaused)
	if err := _StakingEventsMock.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockPoolConcludedIterator struct {
	Event *StakingEventsMockPoolConcluded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockPoolConcludedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockPoolConcluded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockPoolConcluded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockPoolConcludedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockPoolConcludedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockPoolConcluded struct {
	Raw types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterPoolConcluded(opts *bind.FilterOpts) (*StakingEventsMockPoolConcludedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "PoolConcluded")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockPoolConcludedIterator{contract: _StakingEventsMock.contract, event: "PoolConcluded", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchPoolConcluded(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPoolConcluded) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "PoolConcluded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockPoolConcluded)
				if err := _StakingEventsMock.contract.UnpackLog(event, "PoolConcluded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParsePoolConcluded(log types.Log) (*StakingEventsMockPoolConcluded, error) {
	event := new(StakingEventsMockPoolConcluded)
	if err := _StakingEventsMock.contract.UnpackLog(event, "PoolConcluded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockPoolOpenedIterator struct {
	Event *StakingEventsMockPoolOpened

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockPoolOpenedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockPoolOpened)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockPoolOpened)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockPoolOpenedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockPoolOpenedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockPoolOpened struct {
	Raw types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterPoolOpened(opts *bind.FilterOpts) (*StakingEventsMockPoolOpenedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "PoolOpened")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockPoolOpenedIterator{contract: _StakingEventsMock.contract, event: "PoolOpened", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchPoolOpened(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPoolOpened) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "PoolOpened")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockPoolOpened)
				if err := _StakingEventsMock.contract.UnpackLog(event, "PoolOpened", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParsePoolOpened(log types.Log) (*StakingEventsMockPoolOpened, error) {
	event := new(StakingEventsMockPoolOpened)
	if err := _StakingEventsMock.contract.UnpackLog(event, "PoolOpened", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockPoolSizeIncreasedIterator struct {
	Event *StakingEventsMockPoolSizeIncreased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockPoolSizeIncreasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockPoolSizeIncreased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockPoolSizeIncreased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockPoolSizeIncreasedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockPoolSizeIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockPoolSizeIncreased struct {
	MaxPoolSize *big.Int
	Raw         types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterPoolSizeIncreased(opts *bind.FilterOpts) (*StakingEventsMockPoolSizeIncreasedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "PoolSizeIncreased")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockPoolSizeIncreasedIterator{contract: _StakingEventsMock.contract, event: "PoolSizeIncreased", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchPoolSizeIncreased(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPoolSizeIncreased) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "PoolSizeIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockPoolSizeIncreased)
				if err := _StakingEventsMock.contract.UnpackLog(event, "PoolSizeIncreased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParsePoolSizeIncreased(log types.Log) (*StakingEventsMockPoolSizeIncreased, error) {
	event := new(StakingEventsMockPoolSizeIncreased)
	if err := _StakingEventsMock.contract.UnpackLog(event, "PoolSizeIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockRewardAddedIterator struct {
	Event *StakingEventsMockRewardAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockRewardAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockRewardAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockRewardAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockRewardAddedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockRewardAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockRewardAdded struct {
	AmountAdded *big.Int
	Raw         types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterRewardAdded(opts *bind.FilterOpts) (*StakingEventsMockRewardAddedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "RewardAdded")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockRewardAddedIterator{contract: _StakingEventsMock.contract, event: "RewardAdded", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchRewardAdded(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardAdded) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "RewardAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockRewardAdded)
				if err := _StakingEventsMock.contract.UnpackLog(event, "RewardAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseRewardAdded(log types.Log) (*StakingEventsMockRewardAdded, error) {
	event := new(StakingEventsMockRewardAdded)
	if err := _StakingEventsMock.contract.UnpackLog(event, "RewardAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockRewardInitializedIterator struct {
	Event *StakingEventsMockRewardInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockRewardInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockRewardInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockRewardInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockRewardInitializedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockRewardInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockRewardInitialized struct {
	Rate           *big.Int
	Available      *big.Int
	StartTimestamp *big.Int
	EndTimestamp   *big.Int
	Raw            types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterRewardInitialized(opts *bind.FilterOpts) (*StakingEventsMockRewardInitializedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "RewardInitialized")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockRewardInitializedIterator{contract: _StakingEventsMock.contract, event: "RewardInitialized", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchRewardInitialized(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardInitialized) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "RewardInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockRewardInitialized)
				if err := _StakingEventsMock.contract.UnpackLog(event, "RewardInitialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseRewardInitialized(log types.Log) (*StakingEventsMockRewardInitialized, error) {
	event := new(StakingEventsMockRewardInitialized)
	if err := _StakingEventsMock.contract.UnpackLog(event, "RewardInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockRewardRateChangedIterator struct {
	Event *StakingEventsMockRewardRateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockRewardRateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockRewardRateChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockRewardRateChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockRewardRateChangedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockRewardRateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockRewardRateChanged struct {
	Rate *big.Int
	Raw  types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterRewardRateChanged(opts *bind.FilterOpts) (*StakingEventsMockRewardRateChangedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "RewardRateChanged")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockRewardRateChangedIterator{contract: _StakingEventsMock.contract, event: "RewardRateChanged", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchRewardRateChanged(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardRateChanged) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "RewardRateChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockRewardRateChanged)
				if err := _StakingEventsMock.contract.UnpackLog(event, "RewardRateChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseRewardRateChanged(log types.Log) (*StakingEventsMockRewardRateChanged, error) {
	event := new(StakingEventsMockRewardRateChanged)
	if err := _StakingEventsMock.contract.UnpackLog(event, "RewardRateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockRewardSlashedIterator struct {
	Event *StakingEventsMockRewardSlashed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockRewardSlashedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockRewardSlashed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockRewardSlashed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockRewardSlashedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockRewardSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockRewardSlashed struct {
	Operator                []common.Address
	SlashedBaseRewards      []*big.Int
	SlashedDelegatedRewards []*big.Int
	Raw                     types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterRewardSlashed(opts *bind.FilterOpts) (*StakingEventsMockRewardSlashedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "RewardSlashed")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockRewardSlashedIterator{contract: _StakingEventsMock.contract, event: "RewardSlashed", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchRewardSlashed(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardSlashed) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "RewardSlashed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockRewardSlashed)
				if err := _StakingEventsMock.contract.UnpackLog(event, "RewardSlashed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseRewardSlashed(log types.Log) (*StakingEventsMockRewardSlashed, error) {
	event := new(StakingEventsMockRewardSlashed)
	if err := _StakingEventsMock.contract.UnpackLog(event, "RewardSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockRewardWithdrawnIterator struct {
	Event *StakingEventsMockRewardWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockRewardWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockRewardWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockRewardWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockRewardWithdrawnIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockRewardWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockRewardWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterRewardWithdrawn(opts *bind.FilterOpts) (*StakingEventsMockRewardWithdrawnIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "RewardWithdrawn")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockRewardWithdrawnIterator{contract: _StakingEventsMock.contract, event: "RewardWithdrawn", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchRewardWithdrawn(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardWithdrawn) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "RewardWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockRewardWithdrawn)
				if err := _StakingEventsMock.contract.UnpackLog(event, "RewardWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseRewardWithdrawn(log types.Log) (*StakingEventsMockRewardWithdrawn, error) {
	event := new(StakingEventsMockRewardWithdrawn)
	if err := _StakingEventsMock.contract.UnpackLog(event, "RewardWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockStakedIterator struct {
	Event *StakingEventsMockStaked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockStakedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockStaked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockStaked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockStakedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockStaked struct {
	Staker     common.Address
	NewStake   *big.Int
	TotalStake *big.Int
	Raw        types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterStaked(opts *bind.FilterOpts) (*StakingEventsMockStakedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "Staked")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockStakedIterator{contract: _StakingEventsMock.contract, event: "Staked", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakingEventsMockStaked) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "Staked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockStaked)
				if err := _StakingEventsMock.contract.UnpackLog(event, "Staked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseStaked(log types.Log) (*StakingEventsMockStaked, error) {
	event := new(StakingEventsMockStaked)
	if err := _StakingEventsMock.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockUnpausedIterator struct {
	Event *StakingEventsMockUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockUnpausedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterUnpaused(opts *bind.FilterOpts) (*StakingEventsMockUnpausedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockUnpausedIterator{contract: _StakingEventsMock.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StakingEventsMockUnpaused) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockUnpaused)
				if err := _StakingEventsMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseUnpaused(log types.Log) (*StakingEventsMockUnpaused, error) {
	event := new(StakingEventsMockUnpaused)
	if err := _StakingEventsMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StakingEventsMockUnstakedIterator struct {
	Event *StakingEventsMockUnstaked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StakingEventsMockUnstakedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingEventsMockUnstaked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(StakingEventsMockUnstaked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *StakingEventsMockUnstakedIterator) Error() error {
	return it.fail
}

func (it *StakingEventsMockUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StakingEventsMockUnstaked struct {
	Staker           common.Address
	Principal        *big.Int
	BaseReward       *big.Int
	DelegationReward *big.Int
	Raw              types.Log
}

func (_StakingEventsMock *StakingEventsMockFilterer) FilterUnstaked(opts *bind.FilterOpts) (*StakingEventsMockUnstakedIterator, error) {

	logs, sub, err := _StakingEventsMock.contract.FilterLogs(opts, "Unstaked")
	if err != nil {
		return nil, err
	}
	return &StakingEventsMockUnstakedIterator{contract: _StakingEventsMock.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakingEventsMockUnstaked) (event.Subscription, error) {

	logs, sub, err := _StakingEventsMock.contract.WatchLogs(opts, "Unstaked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StakingEventsMockUnstaked)
				if err := _StakingEventsMock.contract.UnpackLog(event, "Unstaked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_StakingEventsMock *StakingEventsMockFilterer) ParseUnstaked(log types.Log) (*StakingEventsMockUnstaked, error) {
	event := new(StakingEventsMockUnstaked)
	if err := _StakingEventsMock.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_StakingEventsMock *StakingEventsMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _StakingEventsMock.abi.Events["AlertRaised"].ID:
		return _StakingEventsMock.ParseAlertRaised(log)
	case _StakingEventsMock.abi.Events["FeedOperatorsSet"].ID:
		return _StakingEventsMock.ParseFeedOperatorsSet(log)
	case _StakingEventsMock.abi.Events["MaxCommunityStakeAmountIncreased"].ID:
		return _StakingEventsMock.ParseMaxCommunityStakeAmountIncreased(log)
	case _StakingEventsMock.abi.Events["MaxOperatorStakeAmountIncreased"].ID:
		return _StakingEventsMock.ParseMaxOperatorStakeAmountIncreased(log)
	case _StakingEventsMock.abi.Events["MerkleRootChanged"].ID:
		return _StakingEventsMock.ParseMerkleRootChanged(log)
	case _StakingEventsMock.abi.Events["Migrated"].ID:
		return _StakingEventsMock.ParseMigrated(log)
	case _StakingEventsMock.abi.Events["MigrationTargetAccepted"].ID:
		return _StakingEventsMock.ParseMigrationTargetAccepted(log)
	case _StakingEventsMock.abi.Events["MigrationTargetProposed"].ID:
		return _StakingEventsMock.ParseMigrationTargetProposed(log)
	case _StakingEventsMock.abi.Events["OperatorAdded"].ID:
		return _StakingEventsMock.ParseOperatorAdded(log)
	case _StakingEventsMock.abi.Events["OperatorRemoved"].ID:
		return _StakingEventsMock.ParseOperatorRemoved(log)
	case _StakingEventsMock.abi.Events["OwnershipTransferRequested"].ID:
		return _StakingEventsMock.ParseOwnershipTransferRequested(log)
	case _StakingEventsMock.abi.Events["OwnershipTransferred"].ID:
		return _StakingEventsMock.ParseOwnershipTransferred(log)
	case _StakingEventsMock.abi.Events["Paused"].ID:
		return _StakingEventsMock.ParsePaused(log)
	case _StakingEventsMock.abi.Events["PoolConcluded"].ID:
		return _StakingEventsMock.ParsePoolConcluded(log)
	case _StakingEventsMock.abi.Events["PoolOpened"].ID:
		return _StakingEventsMock.ParsePoolOpened(log)
	case _StakingEventsMock.abi.Events["PoolSizeIncreased"].ID:
		return _StakingEventsMock.ParsePoolSizeIncreased(log)
	case _StakingEventsMock.abi.Events["RewardAdded"].ID:
		return _StakingEventsMock.ParseRewardAdded(log)
	case _StakingEventsMock.abi.Events["RewardInitialized"].ID:
		return _StakingEventsMock.ParseRewardInitialized(log)
	case _StakingEventsMock.abi.Events["RewardRateChanged"].ID:
		return _StakingEventsMock.ParseRewardRateChanged(log)
	case _StakingEventsMock.abi.Events["RewardSlashed"].ID:
		return _StakingEventsMock.ParseRewardSlashed(log)
	case _StakingEventsMock.abi.Events["RewardWithdrawn"].ID:
		return _StakingEventsMock.ParseRewardWithdrawn(log)
	case _StakingEventsMock.abi.Events["Staked"].ID:
		return _StakingEventsMock.ParseStaked(log)
	case _StakingEventsMock.abi.Events["Unpaused"].ID:
		return _StakingEventsMock.ParseUnpaused(log)
	case _StakingEventsMock.abi.Events["Unstaked"].ID:
		return _StakingEventsMock.ParseUnstaked(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (StakingEventsMockAlertRaised) Topic() common.Hash {
	return common.HexToHash("0xd2720e8f454493f612cc97499fe8cbce7fa4d4c18d346fe7104e9042df1c1edd")
}

func (StakingEventsMockFeedOperatorsSet) Topic() common.Hash {
	return common.HexToHash("0x40aed8e423b39a56b445ae160f4c071fc2cfb48ee0b6dcd5ffeb6bc5b18d10d0")
}

func (StakingEventsMockMaxCommunityStakeAmountIncreased) Topic() common.Hash {
	return common.HexToHash("0xb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df7")
}

func (StakingEventsMockMaxOperatorStakeAmountIncreased) Topic() common.Hash {
	return common.HexToHash("0x816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52")
}

func (StakingEventsMockMerkleRootChanged) Topic() common.Hash {
	return common.HexToHash("0x1b930366dfeaa7eb3b325021e4ae81e36527063452ee55b86c95f85b36f4c31c")
}

func (StakingEventsMockMigrated) Topic() common.Hash {
	return common.HexToHash("0x667838b33bdc898470de09e0e746990f2adc11b965b7fe6828e502ebc39e0434")
}

func (StakingEventsMockMigrationTargetAccepted) Topic() common.Hash {
	return common.HexToHash("0xfa33c052bbee754f3c0482a89962daffe749191fa33c696a61e947fbfd68bd84")
}

func (StakingEventsMockMigrationTargetProposed) Topic() common.Hash {
	return common.HexToHash("0x5c74c441be501340b2713817a6c6975e6f3d4a4ae39fa1ac0bf75d3c54a0cad3")
}

func (StakingEventsMockOperatorAdded) Topic() common.Hash {
	return common.HexToHash("0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d")
}

func (StakingEventsMockOperatorRemoved) Topic() common.Hash {
	return common.HexToHash("0x2360404a74478febece1a14f11275f22ada88d19ef96f7d785913010bfff4479")
}

func (StakingEventsMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (StakingEventsMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (StakingEventsMockPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (StakingEventsMockPoolConcluded) Topic() common.Hash {
	return common.HexToHash("0xf7d0e0f15586495da8c687328ead30fb829d9da55538cb0ef73dd229e517cdb8")
}

func (StakingEventsMockPoolOpened) Topic() common.Hash {
	return common.HexToHash("0xded6ebf04e261e1eb2f3e3b268a2e6aee5b478c15b341eba5cf18b9bc80c2e63")
}

func (StakingEventsMockPoolSizeIncreased) Topic() common.Hash {
	return common.HexToHash("0x7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c0")
}

func (StakingEventsMockRewardAdded) Topic() common.Hash {
	return common.HexToHash("0xde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d")
}

func (StakingEventsMockRewardInitialized) Topic() common.Hash {
	return common.HexToHash("0x125fc8494f786b470e3c39d0932a62e9e09e291ebd81ea19c57604f6d2b1d167")
}

func (StakingEventsMockRewardRateChanged) Topic() common.Hash {
	return common.HexToHash("0x1e3be2efa25bca5bff2215c7b30b31086e703d6aa7d9b9a1f8ba62c5291219ad")
}

func (StakingEventsMockRewardSlashed) Topic() common.Hash {
	return common.HexToHash("0x00635ea9da6e262e92bb713d71840af7c567807ff35bf73e927490c612832480")
}

func (StakingEventsMockRewardWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x150a6ec0e6f4e9ddcaaaa1674f157d91165a42d60653016f87a9fc870a39f050")
}

func (StakingEventsMockStaked) Topic() common.Hash {
	return common.HexToHash("0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90")
}

func (StakingEventsMockUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (StakingEventsMockUnstaked) Topic() common.Hash {
	return common.HexToHash("0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00")
}

func (_StakingEventsMock *StakingEventsMock) Address() common.Address {
	return _StakingEventsMock.address
}

type StakingEventsMockInterface interface {
	EmitAlertRaised(opts *bind.TransactOpts, alerter common.Address, roundId *big.Int, rewardAmount *big.Int) (*types.Transaction, error)

	EmitFeedOperatorsSet(opts *bind.TransactOpts, feedOperators []common.Address) (*types.Transaction, error)

	EmitMaxCommunityStakeAmountIncreased(opts *bind.TransactOpts, maxStakeAmount *big.Int) (*types.Transaction, error)

	EmitMaxOperatorStakeAmountIncreased(opts *bind.TransactOpts, maxStakeAmount *big.Int) (*types.Transaction, error)

	EmitMerkleRootChanged(opts *bind.TransactOpts, newMerkleRoot [32]byte) (*types.Transaction, error)

	EmitMigrated(opts *bind.TransactOpts, staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int, data []byte) (*types.Transaction, error)

	EmitMigrationTargetAccepted(opts *bind.TransactOpts, migrationTarget common.Address) (*types.Transaction, error)

	EmitMigrationTargetProposed(opts *bind.TransactOpts, migrationTarget common.Address) (*types.Transaction, error)

	EmitOperatorAdded(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error)

	EmitOperatorRemoved(opts *bind.TransactOpts, operator common.Address, amount *big.Int) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitPoolConcluded(opts *bind.TransactOpts) (*types.Transaction, error)

	EmitPoolOpened(opts *bind.TransactOpts) (*types.Transaction, error)

	EmitPoolSizeIncreased(opts *bind.TransactOpts, maxPoolSize *big.Int) (*types.Transaction, error)

	EmitRewardAdded(opts *bind.TransactOpts, amountAdded *big.Int) (*types.Transaction, error)

	EmitRewardInitialized(opts *bind.TransactOpts, rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) (*types.Transaction, error)

	EmitRewardRateChanged(opts *bind.TransactOpts, rate *big.Int) (*types.Transaction, error)

	EmitRewardSlashed(opts *bind.TransactOpts, operator []common.Address, slashedBaseRewards []*big.Int, slashedDelegatedRewards []*big.Int) (*types.Transaction, error)

	EmitRewardWithdrawn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	EmitStaked(opts *bind.TransactOpts, staker common.Address, newStake *big.Int, totalStake *big.Int) (*types.Transaction, error)

	EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitUnstaked(opts *bind.TransactOpts, staker common.Address, principal *big.Int, baseReward *big.Int, delegationReward *big.Int) (*types.Transaction, error)

	FilterAlertRaised(opts *bind.FilterOpts) (*StakingEventsMockAlertRaisedIterator, error)

	WatchAlertRaised(opts *bind.WatchOpts, sink chan<- *StakingEventsMockAlertRaised) (event.Subscription, error)

	ParseAlertRaised(log types.Log) (*StakingEventsMockAlertRaised, error)

	FilterFeedOperatorsSet(opts *bind.FilterOpts) (*StakingEventsMockFeedOperatorsSetIterator, error)

	WatchFeedOperatorsSet(opts *bind.WatchOpts, sink chan<- *StakingEventsMockFeedOperatorsSet) (event.Subscription, error)

	ParseFeedOperatorsSet(log types.Log) (*StakingEventsMockFeedOperatorsSet, error)

	FilterMaxCommunityStakeAmountIncreased(opts *bind.FilterOpts) (*StakingEventsMockMaxCommunityStakeAmountIncreasedIterator, error)

	WatchMaxCommunityStakeAmountIncreased(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMaxCommunityStakeAmountIncreased) (event.Subscription, error)

	ParseMaxCommunityStakeAmountIncreased(log types.Log) (*StakingEventsMockMaxCommunityStakeAmountIncreased, error)

	FilterMaxOperatorStakeAmountIncreased(opts *bind.FilterOpts) (*StakingEventsMockMaxOperatorStakeAmountIncreasedIterator, error)

	WatchMaxOperatorStakeAmountIncreased(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMaxOperatorStakeAmountIncreased) (event.Subscription, error)

	ParseMaxOperatorStakeAmountIncreased(log types.Log) (*StakingEventsMockMaxOperatorStakeAmountIncreased, error)

	FilterMerkleRootChanged(opts *bind.FilterOpts) (*StakingEventsMockMerkleRootChangedIterator, error)

	WatchMerkleRootChanged(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMerkleRootChanged) (event.Subscription, error)

	ParseMerkleRootChanged(log types.Log) (*StakingEventsMockMerkleRootChanged, error)

	FilterMigrated(opts *bind.FilterOpts) (*StakingEventsMockMigratedIterator, error)

	WatchMigrated(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMigrated) (event.Subscription, error)

	ParseMigrated(log types.Log) (*StakingEventsMockMigrated, error)

	FilterMigrationTargetAccepted(opts *bind.FilterOpts) (*StakingEventsMockMigrationTargetAcceptedIterator, error)

	WatchMigrationTargetAccepted(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMigrationTargetAccepted) (event.Subscription, error)

	ParseMigrationTargetAccepted(log types.Log) (*StakingEventsMockMigrationTargetAccepted, error)

	FilterMigrationTargetProposed(opts *bind.FilterOpts) (*StakingEventsMockMigrationTargetProposedIterator, error)

	WatchMigrationTargetProposed(opts *bind.WatchOpts, sink chan<- *StakingEventsMockMigrationTargetProposed) (event.Subscription, error)

	ParseMigrationTargetProposed(log types.Log) (*StakingEventsMockMigrationTargetProposed, error)

	FilterOperatorAdded(opts *bind.FilterOpts) (*StakingEventsMockOperatorAddedIterator, error)

	WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOperatorAdded) (event.Subscription, error)

	ParseOperatorAdded(log types.Log) (*StakingEventsMockOperatorAdded, error)

	FilterOperatorRemoved(opts *bind.FilterOpts) (*StakingEventsMockOperatorRemovedIterator, error)

	WatchOperatorRemoved(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOperatorRemoved) (event.Subscription, error)

	ParseOperatorRemoved(log types.Log) (*StakingEventsMockOperatorRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingEventsMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*StakingEventsMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingEventsMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*StakingEventsMockOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*StakingEventsMockPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*StakingEventsMockPaused, error)

	FilterPoolConcluded(opts *bind.FilterOpts) (*StakingEventsMockPoolConcludedIterator, error)

	WatchPoolConcluded(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPoolConcluded) (event.Subscription, error)

	ParsePoolConcluded(log types.Log) (*StakingEventsMockPoolConcluded, error)

	FilterPoolOpened(opts *bind.FilterOpts) (*StakingEventsMockPoolOpenedIterator, error)

	WatchPoolOpened(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPoolOpened) (event.Subscription, error)

	ParsePoolOpened(log types.Log) (*StakingEventsMockPoolOpened, error)

	FilterPoolSizeIncreased(opts *bind.FilterOpts) (*StakingEventsMockPoolSizeIncreasedIterator, error)

	WatchPoolSizeIncreased(opts *bind.WatchOpts, sink chan<- *StakingEventsMockPoolSizeIncreased) (event.Subscription, error)

	ParsePoolSizeIncreased(log types.Log) (*StakingEventsMockPoolSizeIncreased, error)

	FilterRewardAdded(opts *bind.FilterOpts) (*StakingEventsMockRewardAddedIterator, error)

	WatchRewardAdded(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardAdded) (event.Subscription, error)

	ParseRewardAdded(log types.Log) (*StakingEventsMockRewardAdded, error)

	FilterRewardInitialized(opts *bind.FilterOpts) (*StakingEventsMockRewardInitializedIterator, error)

	WatchRewardInitialized(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardInitialized) (event.Subscription, error)

	ParseRewardInitialized(log types.Log) (*StakingEventsMockRewardInitialized, error)

	FilterRewardRateChanged(opts *bind.FilterOpts) (*StakingEventsMockRewardRateChangedIterator, error)

	WatchRewardRateChanged(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardRateChanged) (event.Subscription, error)

	ParseRewardRateChanged(log types.Log) (*StakingEventsMockRewardRateChanged, error)

	FilterRewardSlashed(opts *bind.FilterOpts) (*StakingEventsMockRewardSlashedIterator, error)

	WatchRewardSlashed(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardSlashed) (event.Subscription, error)

	ParseRewardSlashed(log types.Log) (*StakingEventsMockRewardSlashed, error)

	FilterRewardWithdrawn(opts *bind.FilterOpts) (*StakingEventsMockRewardWithdrawnIterator, error)

	WatchRewardWithdrawn(opts *bind.WatchOpts, sink chan<- *StakingEventsMockRewardWithdrawn) (event.Subscription, error)

	ParseRewardWithdrawn(log types.Log) (*StakingEventsMockRewardWithdrawn, error)

	FilterStaked(opts *bind.FilterOpts) (*StakingEventsMockStakedIterator, error)

	WatchStaked(opts *bind.WatchOpts, sink chan<- *StakingEventsMockStaked) (event.Subscription, error)

	ParseStaked(log types.Log) (*StakingEventsMockStaked, error)

	FilterUnpaused(opts *bind.FilterOpts) (*StakingEventsMockUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StakingEventsMockUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*StakingEventsMockUnpaused, error)

	FilterUnstaked(opts *bind.FilterOpts) (*StakingEventsMockUnstakedIterator, error)

	WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakingEventsMockUnstaked) (event.Subscription, error)

	ParseUnstaked(log types.Log) (*StakingEventsMockUnstaked, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
