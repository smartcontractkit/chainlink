// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";

import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_parseSVMExtraArgsFromBytes is FeeQuoterSetup {
  FeeQuoter.DestChainConfig private s_destChainConfig;

  /// @dev a Valid pubkey is one that is 32 bytes long, and that's it since no other validation can be performed
  /// within the constraints of the EVM.
  bytes32 internal constant VALID_SOL_PUBKEY = keccak256("SOL_PUBKEY");

  function setUp() public virtual override {
    super.setUp();
    s_destChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
    s_destChainConfig.enforceOutOfOrder = true; // Enforcing out of order execution for messages to SVM
    s_destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: s_destChainConfig});
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigs);
  }

  function test_SVMExtraArgsV1TagSelector() public pure {
    assertEq(Client.SVM_EXTRA_ARGS_V1_TAG, bytes4(keccak256("CCIP SVMExtraArgsV1")));
  }

  function test_SVMExtraArgsV1() public view {
    bytes32[] memory solAccounts = new bytes32[](1);
    solAccounts[0] = VALID_SOL_PUBKEY;

    Client.SVMExtraArgsV1 memory inputArgs = Client.SVMExtraArgsV1({
      computeUnits: GAS_LIMIT,
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32(0),
      allowOutOfOrderExecution: true,
      accounts: solAccounts
    });

    bytes memory inputExtraArgs = Client._svmArgsToBytes(inputArgs);

    Client.SVMExtraArgsV1 memory expectedOutputArgs = Client.SVMExtraArgsV1({
      computeUnits: GAS_LIMIT,
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32(0),
      allowOutOfOrderExecution: true,
      accounts: solAccounts
    });

    vm.assertEq(
      abi.encode(s_feeQuoter.parseSVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig)),
      abi.encode(expectedOutputArgs)
    );
  }

  // Reverts
  function test_RevertWhen_ExtraArgsAreEmpty() public {
    bytes memory inputExtraArgs = new bytes(0);
    vm.expectRevert(FeeQuoter.InvalidExtraArgsData.selector);
    s_feeQuoter.parseSVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_RevertWhen_InvalidExtraArgsTag() public {
    bytes memory inputExtraArgs = abi.encodeWithSelector(bytes4(0));

    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);
    s_feeQuoter.parseSVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_RevertWhen_SVMMessageGasLimitTooHigh() public {
    Client.SVMExtraArgsV1 memory inputArgs = Client.SVMExtraArgsV1({
      computeUnits: s_destChainConfig.maxPerMsgGasLimit + 1,
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32(0),
      allowOutOfOrderExecution: true,
      accounts: new bytes32[](0)
    });

    bytes memory inputExtraArgs = Client._svmArgsToBytes(inputArgs);

    vm.expectRevert(FeeQuoter.MessageComputeUnitLimitTooHigh.selector);
    s_feeQuoter.parseSVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_RevertWhen_ExtraArgOutOfOrderExecutionIsFalse() public {
    bytes memory inputExtraArgs = abi.encodeWithSelector(
      Client.SVM_EXTRA_ARGS_V1_TAG,
      Client.SVMExtraArgsV1({
        computeUnits: 1_000_000,
        accountIsWritableBitmap: 0,
        tokenReceiver: bytes32(0),
        allowOutOfOrderExecution: false, // mismatch with enforceOutOfOrder = true
        accounts: new bytes32[](0)
      })
    );

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_feeQuoter.parseSVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }
}
