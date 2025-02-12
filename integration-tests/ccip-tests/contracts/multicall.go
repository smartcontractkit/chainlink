package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
)

const (
	MultiCallABI = "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowFailure\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Call3[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate3\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowFailure\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Call3Value[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate3Value\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"blockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBasefee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"basefee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainid\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockCoinbase\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"coinbase\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"difficulty\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockGasLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"gaslimit\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryAggregate\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryBlockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"struct Multicall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]"
	MultiCallBIN = "0x608060405234801561001057600080fd5b50610ee0806100206000396000f3fe6080604052600436106100f35760003560e01c80634d2301cc1161008a578063a8b0574e11610059578063a8b0574e1461025a578063bce38bd714610275578063c3077fa914610288578063ee82ac5e1461029b57600080fd5b80634d2301cc146101ec57806372425d9d1461022157806382ad56cb1461023457806386d516e81461024757600080fd5b80633408e470116100c65780633408e47014610191578063399542e9146101a45780633e64a696146101c657806342cbb15c146101d957600080fd5b80630f28c97d146100f8578063174dea711461011a578063252dba421461013a57806327e86d6e1461015b575b600080fd5b34801561010457600080fd5b50425b6040519081526020015b60405180910390f35b61012d610128366004610a85565b6102ba565b6040516101119190610bbe565b61014d610148366004610a85565b6104ef565b604051610111929190610bd8565b34801561016757600080fd5b50437fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0140610107565b34801561019d57600080fd5b5046610107565b6101b76101b2366004610c60565b610690565b60405161011193929190610cba565b3480156101d257600080fd5b5048610107565b3480156101e557600080fd5b5043610107565b3480156101f857600080fd5b50610107610207366004610ce2565b73ffffffffffffffffffffffffffffffffffffffff163190565b34801561022d57600080fd5b5044610107565b61012d610242366004610a85565b6106ab565b34801561025357600080fd5b5045610107565b34801561026657600080fd5b50604051418152602001610111565b61012d610283366004610c60565b61085a565b6101b7610296366004610a85565b610a1a565b3480156102a757600080fd5b506101076102b6366004610d18565b4090565b60606000828067ffffffffffffffff8111156102d8576102d8610d31565b60405190808252806020026020018201604052801561031e57816020015b6040805180820190915260008152606060208201528152602001906001900390816102f65790505b5092503660005b8281101561047757600085828151811061034157610341610d60565b6020026020010151905087878381811061035d5761035d610d60565b905060200281019061036f9190610d8f565b6040810135958601959093506103886020850185610ce2565b73ffffffffffffffffffffffffffffffffffffffff16816103ac6060870187610dcd565b6040516103ba929190610e32565b60006040518083038185875af1925050503d80600081146103f7576040519150601f19603f3d011682016040523d82523d6000602084013e6103fc565b606091505b50602080850191909152901515808452908501351761046d577f08c379a000000000000000000000000000000000000000000000000000000000600052602060045260176024527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060445260846000fd5b5050600101610325565b508234146104e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4d756c746963616c6c333a2076616c7565206d69736d6174636800000000000060448201526064015b60405180910390fd5b50505092915050565b436060828067ffffffffffffffff81111561050c5761050c610d31565b60405190808252806020026020018201604052801561053f57816020015b606081526020019060019003908161052a5790505b5091503660005b8281101561068657600087878381811061056257610562610d60565b90506020028101906105749190610e42565b92506105836020840184610ce2565b73ffffffffffffffffffffffffffffffffffffffff166105a66020850185610dcd565b6040516105b4929190610e32565b6000604051808303816000865af19150503d80600081146105f1576040519150601f19603f3d011682016040523d82523d6000602084013e6105f6565b606091505b5086848151811061060957610609610d60565b602090810291909101015290508061067d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060448201526064016104dd565b50600101610546565b5050509250929050565b43804060606106a086868661085a565b905093509350939050565b6060818067ffffffffffffffff8111156106c7576106c7610d31565b60405190808252806020026020018201604052801561070d57816020015b6040805180820190915260008152606060208201528152602001906001900390816106e55790505b5091503660005b828110156104e657600084828151811061073057610730610d60565b6020026020010151905086868381811061074c5761074c610d60565b905060200281019061075e9190610e76565b925061076d6020840184610ce2565b73ffffffffffffffffffffffffffffffffffffffff166107906040850185610dcd565b60405161079e929190610e32565b6000604051808303816000865af19150503d80600081146107db576040519150601f19603f3d011682016040523d82523d6000602084013e6107e0565b606091505b506020808401919091529015158083529084013517610851577f08c379a000000000000000000000000000000000000000000000000000000000600052602060045260176024527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060445260646000fd5b50600101610714565b6060818067ffffffffffffffff81111561087657610876610d31565b6040519080825280602002602001820160405280156108bc57816020015b6040805180820190915260008152606060208201528152602001906001900390816108945790505b5091503660005b82811015610a105760008482815181106108df576108df610d60565b602002602001015190508686838181106108fb576108fb610d60565b905060200281019061090d9190610e42565b925061091c6020840184610ce2565b73ffffffffffffffffffffffffffffffffffffffff1661093f6020850185610dcd565b60405161094d929190610e32565b6000604051808303816000865af19150503d806000811461098a576040519150601f19603f3d011682016040523d82523d6000602084013e61098f565b606091505b506020830152151581528715610a07578051610a07576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060448201526064016104dd565b506001016108c3565b5050509392505050565b6000806060610a2b60018686610690565b919790965090945092505050565b60008083601f840112610a4b57600080fd5b50813567ffffffffffffffff811115610a6357600080fd5b6020830191508360208260051b8501011115610a7e57600080fd5b9250929050565b60008060208385031215610a9857600080fd5b823567ffffffffffffffff811115610aaf57600080fd5b610abb85828601610a39565b90969095509350505050565b6000815180845260005b81811015610aed57602081850181015186830182015201610ad1565b81811115610aff576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600082825180855260208086019550808260051b84010181860160005b84811015610bb1578583037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001895281518051151584528401516040858501819052610b9d81860183610ac7565b9a86019a9450505090830190600101610b4f565b5090979650505050505050565b602081526000610bd16020830184610b32565b9392505050565b600060408201848352602060408185015281855180845260608601915060608160051b870101935082870160005b82811015610c52577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018452610c40868351610ac7565b95509284019290840190600101610c06565b509398975050505050505050565b600080600060408486031215610c7557600080fd5b83358015158114610c8557600080fd5b9250602084013567ffffffffffffffff811115610ca157600080fd5b610cad86828701610a39565b9497909650939450505050565b838152826020820152606060408201526000610cd96060830184610b32565b95945050505050565b600060208284031215610cf457600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610bd157600080fd5b600060208284031215610d2a57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112610dc357600080fd5b9190910192915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610e0257600080fd5b83018035915067ffffffffffffffff821115610e1d57600080fd5b602001915036819003821315610a7e57600080fd5b8183823760009101908152919050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1833603018112610dc357600080fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa1833603018112610dc357600080fdfea2646970667358221220bb2b5c71a328032f97c676ae39a1ec2148d3e5d6f73d95e9b17910152d61f16264736f6c634300080c0033"
)

type CallWithValue struct {
	Target       common.Address
	AllowFailure bool
	Value        *big.Int
	CallData     []byte
}

type Call struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

type Result struct {
	Success    bool
	ReturnData []byte
}
type CCIPMsgData struct {
	RouterAddr    common.Address
	ChainSelector uint64
	Msg           router.ClientEVM2AnyMessage
	Fee           *big.Int
}

func TransferTokenCallData(to common.Address, amount *big.Int) ([]byte, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(erc20.ERC20ABI))
	if err != nil {
		return nil, err
	}
	transferToken := erc20ABI.Methods["transfer"]
	inputs, err := transferToken.Inputs.Pack(to, amount)
	if err != nil {
		return nil, err
	}
	inputs = append(transferToken.ID[:], inputs...)
	return inputs, nil
}

// ApproveTokenCallData returns the call data for approving a token with approve function of erc20 contract
func ApproveTokenCallData(to common.Address, amount *big.Int) ([]byte, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(erc20.ERC20ABI))
	if err != nil {
		return nil, err
	}
	approveToken := erc20ABI.Methods["approve"]
	inputs, err := approveToken.Inputs.Pack(to, amount)
	if err != nil {
		return nil, err
	}
	inputs = append(approveToken.ID[:], inputs...)
	return inputs, nil
}

// CCIPSendCallData returns the call data for sending a CCIP message with ccipSend function of router contract
func CCIPSendCallData(msg CCIPMsgData) ([]byte, error) {
	routerABI, err := abi.JSON(strings.NewReader(router.RouterABI))
	if err != nil {
		return nil, err
	}
	ccipSend := routerABI.Methods["ccipSend"]
	sendID := ccipSend.ID
	inputs, err := ccipSend.Inputs.Pack(
		msg.ChainSelector,
		msg.Msg,
	)
	if err != nil {
		return nil, err
	}
	inputs = append(sendID[:], inputs...)
	return inputs, nil
}

func WaitForSuccessfulTxMined(evmClient blockchain.EVMClient, tx *types.Transaction) error {
	log.Info().Str("tx", tx.Hash().Hex()).Msg("waiting for tx to be mined")
	receipt, err := bind.WaitMined(context.Background(), evmClient.DeployBackend(), tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		// TODO: Add error reason from receipt/tx
		return fmt.Errorf("tx failed %s", tx.Hash().Hex())
	}
	log.Info().Str("tx", tx.Hash().Hex()).Str("Network", evmClient.GetNetworkName()).Msg("tx mined successfully")
	return nil
}

// MultiCallCCIP sends multiple CCIP messages in a single transaction
// if native is true, it will send msg with native as fee. In this case the msg should be sent with a
// msg.value equivalent to the total fee with the help of aggregate3Value
//
// if native is false, it will send msg with fee in specific feetoken. In this case the msg should be sent without value with the help of aggregate3.
// In both cases, if there are any bridge tokens included in ccip transfer, the amount for corresponding token should be approved to the router contract as spender.
// The approval should be done by calling approval function as part of the call data of aggregate3 or aggregate3Value
// If feetoken is used as fee, the amount for feetoken should be approved to the router contract as spender and should be done as part of the call data of aggregate3
// In case of native as fee, there is no need for fee amount approval
func MultiCallCCIP(
	evmClient blockchain.EVMClient,
	address string,
	msgData []CCIPMsgData,
	native bool,
) (*types.Transaction, error) {
	contractAddress := common.HexToAddress(address)
	multiCallABI, err := abi.JSON(strings.NewReader(MultiCallABI))
	if err != nil {
		return nil, err
	}
	boundContract := bind.NewBoundContract(contractAddress, multiCallABI, evmClient.Backend(), evmClient.Backend(), evmClient.Backend())

	// if native, use aggregate3Value to send msg with value
	if native {
		var callData []CallWithValue
		allValue := big.NewInt(0)
		// create call data for each msg
		for _, msg := range msgData {
			if msg.Msg.FeeToken != (common.Address{}) {
				return nil, fmt.Errorf("fee token should be %s for native as fee", common.HexToAddress("0x0").Hex())
			}
			// approve bridge token
			for _, tokenAndAmount := range msg.Msg.TokenAmounts {
				inputs, err := ApproveTokenCallData(msg.RouterAddr, tokenAndAmount.Amount)
				if err != nil {
					return nil, err
				}
				data := CallWithValue{Target: tokenAndAmount.Token, AllowFailure: false, Value: big.NewInt(0), CallData: inputs}
				callData = append(callData, data)
			}
			inputs, err := CCIPSendCallData(msg)
			if err != nil {
				return nil, err
			}
			data := CallWithValue{Target: msg.RouterAddr, AllowFailure: false, Value: msg.Fee, CallData: inputs}
			callData = append(callData, data)
			allValue.Add(allValue, msg.Fee)
		}

		opts, err := evmClient.TransactionOpts(evmClient.GetDefaultWallet())
		if err != nil {
			return nil, err
		}
		// the value of transactionOpts is the sum of the value of all msg, which is the total fee of all ccip-sends
		opts.Value = allValue

		// call aggregate3Value to group all msg call data and send them in a single transaction
		tx, err := boundContract.Transact(opts, "aggregate3Value", callData)
		if err != nil {
			return nil, err
		}
		err = evmClient.MarkTxAsSentOnL2(tx)
		if err != nil {
			return nil, err
		}
		err = WaitForSuccessfulTxMined(evmClient, tx)
		if err != nil {
			return nil, errors.Wrapf(err, "multicall failed for ccip-send; multicall %s", contractAddress.Hex())
		}
		return tx, nil
	}
	// if with feetoken, use aggregate3 to send msg without value
	var callData []Call
	// create call data for each msg
	for _, msg := range msgData {
		isFeeTokenAndBridgeTokenSame := false
		// approve bridge token
		for _, tokenAndAmount := range msg.Msg.TokenAmounts {
			var inputs []byte
			// if feetoken is same as bridge token, approve total amount including transfer amount + fee amount
			if tokenAndAmount.Token == msg.Msg.FeeToken {
				isFeeTokenAndBridgeTokenSame = true
				inputs, err = ApproveTokenCallData(msg.RouterAddr, new(big.Int).Add(msg.Fee, tokenAndAmount.Amount))
				if err != nil {
					return nil, err
				}
			} else {
				inputs, err = ApproveTokenCallData(msg.RouterAddr, tokenAndAmount.Amount)
				if err != nil {
					return nil, err
				}
			}

			callData = append(callData, Call{Target: tokenAndAmount.Token, AllowFailure: false, CallData: inputs})
		}
		// approve fee token if not already approved
		if msg.Fee != nil && msg.Fee.Cmp(big.NewInt(0)) > 0 && !isFeeTokenAndBridgeTokenSame {
			inputs, err := ApproveTokenCallData(msg.RouterAddr, msg.Fee)
			if err != nil {
				return nil, err
			}
			callData = append(callData, Call{Target: msg.Msg.FeeToken, AllowFailure: false, CallData: inputs})
		}

		inputs, err := CCIPSendCallData(msg)
		if err != nil {
			return nil, err
		}
		callData = append(callData, Call{Target: msg.RouterAddr, AllowFailure: false, CallData: inputs})
	}
	opts, err := evmClient.TransactionOpts(evmClient.GetDefaultWallet())
	if err != nil {
		return nil, err
	}

	// call aggregate3 to group all msg call data and send them in a single transaction
	tx, err := boundContract.Transact(opts, "aggregate3", callData)
	if err != nil {
		return nil, err
	}
	err = WaitForSuccessfulTxMined(evmClient, tx)
	if err != nil {
		return tx, errors.Wrapf(err, "multicall failed for ccip-send; router %s", contractAddress.Hex())
	}
	return tx, nil
}

func TransferTokens(
	evmClient blockchain.EVMClient,
	contractAddress common.Address,
	tokens []*ERC20Token,
) error {
	multiCallABI, err := abi.JSON(strings.NewReader(MultiCallABI))
	if err != nil {
		return err
	}
	var callData []Call
	boundContract := bind.NewBoundContract(contractAddress, multiCallABI, evmClient.Backend(), evmClient.Backend(), evmClient.Backend())
	for _, token := range tokens {
		var inputs []byte
		balance, err := token.BalanceOf(context.Background(), contractAddress.Hex())
		if err != nil {
			return err
		}
		inputs, err = TransferTokenCallData(common.HexToAddress(evmClient.GetDefaultWallet().Address()), balance)
		if err != nil {
			return err
		}
		data := Call{Target: token.ContractAddress, AllowFailure: false, CallData: inputs}
		callData = append(callData, data)
	}

	opts, err := evmClient.TransactionOpts(evmClient.GetDefaultWallet())
	if err != nil {
		return err
	}

	// call aggregate3 to group all msg call data and send them in a single transaction
	tx, err := boundContract.Transact(opts, "aggregate3", callData)
	if err != nil {
		return err
	}
	err = WaitForSuccessfulTxMined(evmClient, tx)
	if err != nil {
		return errors.Wrapf(err, "token transfer failed for token; router %s", contractAddress.Hex())
	}
	return nil
}
