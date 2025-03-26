// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_processChainFamilySelector is FeeQuoterSetup {
  uint64 internal constant SVM_SELECTOR = SOURCE_CHAIN_SELECTOR;
  uint64 internal constant EVM_SELECTOR = DEST_CHAIN_SELECTOR;
  uint64 internal constant APTOS_SELECTOR = DEST_CHAIN_SELECTOR + 1;
  uint64 internal constant INVALID_SELECTOR = 99;

  function setUp() public virtual override {
    super.setUp();

    // 1. Configure an EVM chain
    FeeQuoter.DestChainConfig memory evmConfig;
    evmConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;
    evmConfig.defaultTxGasLimit = 500_000;
    evmConfig.maxPerMsgGasLimit = 1_000_000;
    evmConfig.enforceOutOfOrder = false;

    // 2. Configure an SVM chain
    FeeQuoter.DestChainConfig memory svmConfig;
    svmConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;
    svmConfig.defaultTxGasLimit = 2_000_000;
    svmConfig.maxPerMsgGasLimit = 3_000_000;
    svmConfig.enforceOutOfOrder = true;

    // 2. Configure an SVM chain
    FeeQuoter.DestChainConfig memory aptosConfig;
    aptosConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_APTOS;
    aptosConfig.defaultTxGasLimit = 2_000_000;
    aptosConfig.maxPerMsgGasLimit = 3_000_000;
    aptosConfig.enforceOutOfOrder = true;

    // Apply both configs
    FeeQuoter.DestChainConfigArgs[] memory configs = new FeeQuoter.DestChainConfigArgs[](3);
    configs[0] = FeeQuoter.DestChainConfigArgs({destChainSelector: EVM_SELECTOR, destChainConfig: evmConfig});
    configs[1] = FeeQuoter.DestChainConfigArgs({destChainSelector: SVM_SELECTOR, destChainConfig: svmConfig});
    configs[2] = FeeQuoter.DestChainConfigArgs({destChainSelector: APTOS_SELECTOR, destChainConfig: aptosConfig});
    s_feeQuoter.applyDestChainConfigUpdates(configs);
  }

  function test_processChainFamilySelector_EVM() public view {
    Client.GenericExtraArgsV2 memory evmArgs =
      Client.GenericExtraArgsV2({gasLimit: 400_000, allowOutOfOrderExecution: true});
    bytes memory encodedArgs = Client._argsToBytes(evmArgs);

    (bytes memory resultBytes, bool outOfOrder, bytes memory tokenReceiver) =
      s_feeQuoter.processChainFamilySelector(EVM_SELECTOR, MESSAGE_RECEIVER, encodedArgs);

    assertEq(resultBytes, encodedArgs, "Should return the same EVM-encoded bytes");
    assertEq(outOfOrder, evmArgs.allowOutOfOrderExecution, "Out-of-order mismatch");
    assertEq(tokenReceiver, MESSAGE_RECEIVER, "Token receiver mismatch");
  }

  function test_processChainFamilySelector_Aptos() public view {
    Client.GenericExtraArgsV2 memory genericExtraArgs =
      Client.GenericExtraArgsV2({gasLimit: 400_000, allowOutOfOrderExecution: true});
    bytes memory encodedArgs = Client._argsToBytes(genericExtraArgs);

    (bytes memory resultBytes, bool outOfOrder, bytes memory tokenReceiver) =
      s_feeQuoter.processChainFamilySelector(APTOS_SELECTOR, MESSAGE_RECEIVER, encodedArgs);

    assertEq(resultBytes, encodedArgs, "Should return the same EVM-encoded bytes");
    assertEq(outOfOrder, genericExtraArgs.allowOutOfOrderExecution, "Out-of-order mismatch");
    assertEq(tokenReceiver, MESSAGE_RECEIVER, "Token receiver mismatch");
  }

  function test_processChainFamilySelector_SVM_WithTokenTransfer() public view {
    // Construct an SVMExtraArgsV1 with a non-zero tokenReceiver
    Client.SVMExtraArgsV1 memory svmArgs = Client.SVMExtraArgsV1({
      computeUnits: 1_500_000, // within the limit
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32("someReceiver"),
      allowOutOfOrderExecution: true,
      accounts: new bytes32[](0)
    });
    bytes memory encodedSvmArgs = Client._svmArgsToBytes(svmArgs);

    (bytes memory resultBytes, bool outOfOrder, bytes memory tokenReceiver) =
      s_feeQuoter.processChainFamilySelector(SVM_SELECTOR, MESSAGE_RECEIVER, encodedSvmArgs);

    // The function should NOT revert since tokenReceiver != 0
    // Check that it returned the SVM-encoded bytes
    assertEq(resultBytes, encodedSvmArgs, "Should return the same SVM-encoded bytes");
    // The function always returns `true` for outOfOrder on SVM
    assertTrue(outOfOrder, "Out-of-order for SVM must be true");
    assertEq(tokenReceiver, abi.encode(bytes32("someReceiver")));
  }

  function test_processChainFamilySelector_SVM_NoTokenTransfer() public view {
    Client.SVMExtraArgsV1 memory svmArgs = Client.SVMExtraArgsV1({
      computeUnits: 2_000_000,
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32(0), // zero is fine if not transferring tokens
      allowOutOfOrderExecution: true,
      accounts: new bytes32[](0)
    });
    bytes memory encodedSvmArgs = Client._svmArgsToBytes(svmArgs);

    (bytes memory resultBytes, bool outOfOrder,) =
      s_feeQuoter.processChainFamilySelector(SVM_SELECTOR, MESSAGE_RECEIVER, encodedSvmArgs);

    // Should succeed with outOfOrder = true
    assertEq(resultBytes, encodedSvmArgs, "Should return the SVM-encoded bytes");
    assertTrue(outOfOrder, "Out-of-order should be true for SVM");
  }

  function test_processChainFamilySelector_RevertWhen_InvalidChainFamilySelector() public {
    // Provide random extraArgs
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.InvalidChainFamilySelector.selector, bytes4(0)));

    s_feeQuoter.processChainFamilySelector(INVALID_SELECTOR, MESSAGE_RECEIVER, "0x1234");
  }
}
