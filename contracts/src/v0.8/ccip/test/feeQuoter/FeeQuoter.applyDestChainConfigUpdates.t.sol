// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_applyDestChainConfigUpdates is FeeQuoterSetup {
  function testFuzz_applyDestChainConfigUpdates_Success(
    FeeQuoter.DestChainConfigArgs memory destChainConfigArgs
  ) public {
    vm.assume(destChainConfigArgs.destChainSelector != 0);
    vm.assume(destChainConfigArgs.destChainConfig.maxPerMsgGasLimit != 0);
    destChainConfigArgs.destChainConfig.defaultTxGasLimit = uint32(
      bound(
        destChainConfigArgs.destChainConfig.defaultTxGasLimit, 1, destChainConfigArgs.destChainConfig.maxPerMsgGasLimit
      )
    );
    destChainConfigArgs.destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;

    bool isNewChain = destChainConfigArgs.destChainSelector != DEST_CHAIN_SELECTOR;

    FeeQuoter.DestChainConfigArgs[] memory newDestChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    newDestChainConfigArgs[0] = destChainConfigArgs;

    if (isNewChain) {
      vm.expectEmit();
      emit FeeQuoter.DestChainAdded(destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig);
    } else {
      vm.expectEmit();
      emit FeeQuoter.DestChainConfigUpdated(destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig);
    }

    s_feeQuoter.applyDestChainConfigUpdates(newDestChainConfigArgs);

    _assertFeeQuoterDestChainConfigsEqual(
      destChainConfigArgs.destChainConfig, s_feeQuoter.getDestChainConfig(destChainConfigArgs.destChainSelector)
    );
  }

  function test_applyDestChainConfigUpdates_Success() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](2);
    destChainConfigArgs[0] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[0].destChainConfig.isEnabled = false;
    destChainConfigArgs[1] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    vm.expectEmit();
    emit FeeQuoter.DestChainConfigUpdated(DEST_CHAIN_SELECTOR, destChainConfigArgs[0].destChainConfig);
    vm.expectEmit();
    emit FeeQuoter.DestChainAdded(DEST_CHAIN_SELECTOR + 1, destChainConfigArgs[1].destChainConfig);

    vm.recordLogs();
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    FeeQuoter.DestChainConfig memory gotDestChainConfig0 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
    FeeQuoter.DestChainConfig memory gotDestChainConfig1 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);

    assertEq(vm.getRecordedLogs().length, 2);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[0].destChainConfig, gotDestChainConfig0);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[1].destChainConfig, gotDestChainConfig1);
  }

  function test_applyDestChainConfigUpdatesZeroInput_Success() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](0);

    vm.recordLogs();
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitEqZero_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.defaultTxGasLimit = 0;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitGtMaxPerMessageGasLimit_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    // Allow setting to the max value
    destChainConfigArg.destChainConfig.defaultTxGasLimit = destChainConfigArg.destChainConfig.maxPerMsgGasLimit;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    // Revert when exceeding max value
    destChainConfigArg.destChainConfig.defaultTxGasLimit = destChainConfigArg.destChainConfig.maxPerMsgGasLimit + 1;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestChainConfigDestChainSelectorEqZero_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainSelector = 0;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidChainFamilySelector_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.chainFamilySelector = bytes4(uint32(1));

    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }
}
