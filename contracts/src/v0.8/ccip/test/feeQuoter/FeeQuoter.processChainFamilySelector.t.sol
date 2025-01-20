// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_processChainFamilySelector is FeeQuoterSetup {
  uint64 internal constant SVM_SELECTOR = SOURCE_CHAIN_SELECTOR;
  uint64 internal constant EVM_SELECTOR = DEST_CHAIN_SELECTOR;
  uint64 internal constant INVALID_SELECTOR = 99;

  function setUp() public virtual override {
    super.setUp();

    // 1. Configure an EVM chain
    FeeQuoter.DestChainConfig memory evmConfig;
    evmConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;
    evmConfig.defaultTxGasLimit = 500_000;
    evmConfig.maxPerMsgGasLimit = 1_000_000; // Example
    evmConfig.enforceOutOfOrder = false; // Example

    // 2. Configure an SVM chain
    FeeQuoter.DestChainConfig memory svmConfig;
    svmConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;
    svmConfig.defaultTxGasLimit = 2_000_000;
    svmConfig.maxPerMsgGasLimit = 3_000_000; // Example
    svmConfig.enforceOutOfOrder = true;

    // Apply both configs
    FeeQuoter.DestChainConfigArgs[] memory configs = new FeeQuoter.DestChainConfigArgs[](2);
    configs[0] = FeeQuoter.DestChainConfigArgs({destChainSelector: EVM_SELECTOR, destChainConfig: evmConfig});
    configs[1] = FeeQuoter.DestChainConfigArgs({destChainSelector: SVM_SELECTOR, destChainConfig: svmConfig});
    s_feeQuoter.applyDestChainConfigUpdates(configs);
  }

  // ----------------------------------------------------------------
  // TEST: EVM path
  // ----------------------------------------------------------------
  function test_processChainFamilySelectorEVM() public {
    Client.EVMExtraArgsV2 memory evmArgs = Client.EVMExtraArgsV2({gasLimit: 400_000, allowOutOfOrderExecution: true});
    bytes memory encodedEvmArgs = Client._argsToBytes(evmArgs);

    (bytes memory resultBytes, bool outOfOrder) = s_feeQuoter.processChainFamilySelector(
      EVM_SELECTOR,
      false, // isMessageWithTokenTransfer
      encodedEvmArgs
    );

    assertEq(resultBytes, encodedEvmArgs, "Should return the same EVM-encoded bytes");
    assertEq(outOfOrder, evmArgs.allowOutOfOrderExecution, "Out-of-order mismatch");
  }

  // ----------------------------------------------------------------
  // TEST: SVM path
  // ----------------------------------------------------------------
  function test_processChainFamilySelectorSVM_WithTokenTransfer() public {
    // Construct an SVMExtraArgsV1 with a non-zero tokenReceiver
    Client.SVMExtraArgsV1 memory svmArgs = Client.SVMExtraArgsV1({
      computeUnits: 1_500_000, // within the limit
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32("someReceiver"),
      allowOutOfOrderExecution: true,
      accounts: new bytes32[](0)
    });
    bytes memory encodedSvmArgs = Client._svmArgsToBytes(svmArgs);

    (bytes memory resultBytes, bool outOfOrder) = s_feeQuoter.processChainFamilySelector(
      SVM_SELECTOR,
      true, // isMessageWithTokenTransfer
      encodedSvmArgs
    );

    // The function should NOT revert since tokenReceiver != 0
    // Check that it returned the SVM-encoded bytes
    assertEq(resultBytes, encodedSvmArgs, "Should return the same SVM-encoded bytes");
    // The function always returns `true` for outOfOrder on SVM
    assertTrue(outOfOrder, "Out-of-order for SVM must be true");
  }

  function test_processChainFamilySelectorSVM_NoTokenTransfer() public {
    Client.SVMExtraArgsV1 memory svmArgs = Client.SVMExtraArgsV1({
      computeUnits: 2_000_000,
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32(0), // zero is fine if not transferring tokens
      allowOutOfOrderExecution: true,
      accounts: new bytes32[](0)
    });
    bytes memory encodedSvmArgs = Client._svmArgsToBytes(svmArgs);

    (bytes memory resultBytes, bool outOfOrder) = s_feeQuoter.processChainFamilySelector(
      SVM_SELECTOR,
      false, // no token transfer
      encodedSvmArgs
    );

    // Should succeed with outOfOrder = true
    assertEq(resultBytes, encodedSvmArgs, "Should return the SVM-encoded bytes");
    assertTrue(outOfOrder, "Out-of-order should be true for SVM");
  }

  // TEST: SVM path → reverts

  function test_processChainFamilySelector_RevertWhen_TokenTransferNoTokenReceiver() public {
    // Construct an SVMExtraArgsV1 with tokenReceiver == 0
    Client.SVMExtraArgsV1 memory svmArgs = Client.SVMExtraArgsV1({
      computeUnits: 1_500_000,
      accountIsWritableBitmap: 0,
      tokenReceiver: bytes32(0), // <-- zero
      allowOutOfOrderExecution: true,
      accounts: new bytes32[](0)
    });
    bytes memory encodedSvmArgs = Client._svmArgsToBytes(svmArgs);

    vm.expectRevert(FeeQuoter.InvalidTokenReceiver.selector);

    s_feeQuoter.processChainFamilySelector(
      SVM_SELECTOR,
      true, // token transfer
      encodedSvmArgs
    );
  }

  function test_processChainFamilySelector_RevertWhen_InvalidChainFamilySelector() public {
    // Provide random extraArgs
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.InvalidChainFamilySelector.selector, bytes4(0)));

    s_feeQuoter.processChainFamilySelector(INVALID_SELECTOR, false, "0x1234");
  }
}
