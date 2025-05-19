// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_applyDestChainConfigUpdates is FeeQuoterSetup {
  bytes4[] internal s_validChainFamilySelector =
    [Internal.CHAIN_FAMILY_SELECTOR_EVM, Internal.CHAIN_FAMILY_SELECTOR_SVM];

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

    for (uint256 i = 0; i < s_validChainFamilySelector.length; i++) {
      destChainConfigArgs.destChainConfig.chainFamilySelector = s_validChainFamilySelector[i];
      destChainConfigArgs.destChainSelector = uint64(
        uint256(keccak256(abi.encodePacked(destChainConfigArgs.destChainSelector, s_validChainFamilySelector[i])))
      );

      FeeQuoter.DestChainConfigArgs[] memory newDestChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
      newDestChainConfigArgs[0] = destChainConfigArgs;

      vm.expectEmit();
      emit FeeQuoter.DestChainAdded(destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig);

      s_feeQuoter.applyDestChainConfigUpdates(newDestChainConfigArgs);

      _assertFeeQuoterDestChainConfigsEqual(
        destChainConfigArgs.destChainConfig, s_feeQuoter.getDestChainConfig(destChainConfigArgs.destChainSelector)
      );
    }
  }

  function test_applyDestChainConfigUpdates() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](3);
    destChainConfigArgs[0] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[0].destChainConfig.isEnabled = false;

    destChainConfigArgs[1] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;
    destChainConfigArgs[1].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    destChainConfigArgs[2] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[2].destChainSelector = DEST_CHAIN_SELECTOR + 2;
    destChainConfigArgs[2].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_APTOS;

    vm.expectEmit();
    emit FeeQuoter.DestChainConfigUpdated(DEST_CHAIN_SELECTOR, destChainConfigArgs[0].destChainConfig);
    vm.expectEmit();
    emit FeeQuoter.DestChainAdded(DEST_CHAIN_SELECTOR + 1, destChainConfigArgs[1].destChainConfig);

    vm.recordLogs();
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    FeeQuoter.DestChainConfig memory gotDestChainConfig0 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
    FeeQuoter.DestChainConfig memory gotDestChainConfig1 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);
    FeeQuoter.DestChainConfig memory gotDestChainConfig2 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR + 2);

    assertEq(vm.getRecordedLogs().length, destChainConfigArgs.length);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[0].destChainConfig, gotDestChainConfig0);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[1].destChainConfig, gotDestChainConfig1);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[2].destChainConfig, gotDestChainConfig2);
  }

  function test_applyDestChainConfigUpdatesZeroInput() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](0);

    vm.recordLogs();
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_RevertWhen_applyDestChainConfigUpdatesDefaultTxGasLimitEqZero() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.defaultTxGasLimit = 0;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_RevertWhen_applyDestChainConfigUpdatesDefaultTxGasLimitGtMaxPerMessageGasLimit() public {
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

  function test_RevertWhen_InvalidDestChainConfigDestChainSelectorEqZero() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainSelector = 0;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_RevertWhen_InvalidChainFamilySelector() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.chainFamilySelector = bytes4(uint32(1));

    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }
}
