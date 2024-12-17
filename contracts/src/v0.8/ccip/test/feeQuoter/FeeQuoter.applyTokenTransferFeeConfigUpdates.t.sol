// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Pool} from "../../libraries/Pool.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_applyTokenTransferFeeConfigUpdates is FeeQuoterSetup {
  function testFuzz_ApplyTokenTransferFeeConfig_Success(
    FeeQuoter.TokenTransferFeeConfig[2] memory tokenTransferFeeConfigs
  ) public {
    // To prevent Invalid Fee Range error from the fuzzer, bound the results to a valid range that
    // where minFee < maxFee
    tokenTransferFeeConfigs[0].minFeeUSDCents =
      uint32(bound(tokenTransferFeeConfigs[0].minFeeUSDCents, 0, type(uint8).max));
    tokenTransferFeeConfigs[1].minFeeUSDCents =
      uint32(bound(tokenTransferFeeConfigs[1].minFeeUSDCents, 0, type(uint8).max));

    tokenTransferFeeConfigs[0].maxFeeUSDCents = uint32(
      bound(tokenTransferFeeConfigs[0].maxFeeUSDCents, tokenTransferFeeConfigs[0].minFeeUSDCents + 1, type(uint32).max)
    );
    tokenTransferFeeConfigs[1].maxFeeUSDCents = uint32(
      bound(tokenTransferFeeConfigs[1].maxFeeUSDCents, tokenTransferFeeConfigs[1].minFeeUSDCents + 1, type(uint32).max)
    );

    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(2, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    for (uint256 i = 0; i < tokenTransferFeeConfigArgs.length; ++i) {
      for (uint256 j = 0; j < tokenTransferFeeConfigs.length; ++j) {
        tokenTransferFeeConfigs[j].destBytesOverhead = uint32(
          bound(tokenTransferFeeConfigs[j].destBytesOverhead, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, type(uint32).max)
        );
        address feeToken = s_sourceTokens[j];
        tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs[j].token = feeToken;
        tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs[j].tokenTransferFeeConfig = tokenTransferFeeConfigs[j];

        vm.expectEmit();
        emit FeeQuoter.TokenTransferFeeConfigUpdated(
          tokenTransferFeeConfigArgs[i].destChainSelector, feeToken, tokenTransferFeeConfigs[j]
        );
      }
    }

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    for (uint256 i = 0; i < tokenTransferFeeConfigs.length; ++i) {
      _assertTokenTransferFeeConfigEqual(
        tokenTransferFeeConfigs[i],
        s_feeQuoter.getTokenTransferFeeConfig(
          tokenTransferFeeConfigArgs[0].destChainSelector,
          tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[i].token
        )
      );
    }
  }

  function test_ApplyTokenTransferFeeConfig() public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: 312,
      isEnabled: true
    });
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token = address(11);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 12,
      maxFeeUSDCents: 13,
      deciBps: 14,
      destGasOverhead: 15,
      destBytesOverhead: 394,
      isEnabled: true
    });

    vm.expectEmit();
    emit FeeQuoter.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig
    );
    vm.expectEmit();
    emit FeeQuoter.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig
    );

    FeeQuoter.TokenTransferFeeConfigRemoveArgs[] memory tokensToRemove =
      new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0);
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs, tokensToRemove);

    FeeQuoter.TokenTransferFeeConfig memory config0 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    FeeQuoter.TokenTransferFeeConfig memory config1 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig, config0
    );
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );

    // Remove only the first token and validate only the first token is removed
    tokensToRemove = new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](1);
    tokensToRemove[0] = FeeQuoter.TokenTransferFeeConfigRemoveArgs({
      destChainSelector: tokenTransferFeeConfigArgs[0].destChainSelector,
      token: tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    });

    vm.expectEmit();
    emit FeeQuoter.TokenTransferFeeConfigDeleted(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(new FeeQuoter.TokenTransferFeeConfigArgs[](0), tokensToRemove);

    config0 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    config1 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    FeeQuoter.TokenTransferFeeConfig memory emptyConfig;

    _assertTokenTransferFeeConfigEqual(emptyConfig, config0);
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );
  }

  function test_ApplyTokenTransferFeeZeroInput() public {
    vm.recordLogs();
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      new FeeQuoter.TokenTransferFeeConfigArgs[](0), new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_RevertWhen_OnlyCallableByOwnerOrAdmin() public {
    vm.startPrank(STRANGER);
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs;

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }

  function test_RevertWhen_InvalidDestBytesOverhead() public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES - 1),
      isEnabled: true
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        FeeQuoter.InvalidDestBytesOverhead.selector,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.destBytesOverhead
      )
    );

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }
}
