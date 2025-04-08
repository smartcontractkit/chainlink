// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_resolveGasLimitForDestination is FeeQuoterSetup {
  FeeQuoter.DestChainConfig private s_destChainConfig;

  function setUp() public virtual override {
    super.setUp();
    s_destChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
  }

  function test_EVMExtraArgsV1TagSelector() public pure {
    assertEq(Client.EVM_EXTRA_ARGS_V1_TAG, bytes4(keccak256("CCIP EVMExtraArgsV1")));
  }

  function test_EVMExtraArgsV2TagSelector() public pure {
    assertEq(Client.GENERIC_EXTRA_ARGS_V2_TAG, bytes4(keccak256("CCIP EVMExtraArgsV2")));
  }

  function test_EVMExtraArgsV1() public view {
    Client.EVMExtraArgsV1 memory inputArgs = Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    Client.GenericExtraArgsV2 memory expectedOutputArgs =
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});

    vm.assertEq(
      abi.encode(s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, DEST_CHAIN_SELECTOR)),
      abi.encode(expectedOutputArgs)
    );
  }

  function test_EVMExtraArgsV2() public view {
    Client.GenericExtraArgsV2 memory inputArgs =
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

    vm.assertEq(
      abi.encode(s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, DEST_CHAIN_SELECTOR)), abi.encode(inputArgs)
    );
  }

  function test_EVMExtraArgsDefault() public view {
    Client.GenericExtraArgsV2 memory expectedOutputArgs =
      Client.GenericExtraArgsV2({gasLimit: s_destChainConfig.defaultTxGasLimit, allowOutOfOrderExecution: false});

    vm.assertEq(
      abi.encode(s_feeQuoter.parseEVMExtraArgsFromBytes("", DEST_CHAIN_SELECTOR)), abi.encode(expectedOutputArgs)
    );
  }

  // Reverts

  function test_RevertWhen_EVMExtraArgsInvalidExtraArgsTag() public {
    Client.GenericExtraArgsV2 memory inputArgs =
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    // Invalidate selector
    inputExtraArgs[0] = bytes1(uint8(0));

    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);
    s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, DEST_CHAIN_SELECTOR);
  }

  function test_RevertWhen_EVMExtraArgsEnforceOutOfOrder() public {
    Client.GenericExtraArgsV2 memory inputArgs =
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    s_destChainConfig.enforceOutOfOrder = true;

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, DEST_CHAIN_SELECTOR, true);
  }

  function test_RevertWhen_EVMExtraArgsGasLimitTooHigh() public {
    Client.GenericExtraArgsV2 memory inputArgs =
      Client.GenericExtraArgsV2({gasLimit: s_destChainConfig.maxPerMsgGasLimit + 1, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

    vm.expectRevert(FeeQuoter.MessageGasLimitTooHigh.selector);
    s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, DEST_CHAIN_SELECTOR);
  }
}
