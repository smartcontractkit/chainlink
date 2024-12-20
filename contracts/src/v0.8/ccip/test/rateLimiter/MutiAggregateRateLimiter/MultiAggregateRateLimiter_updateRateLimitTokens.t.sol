// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";

import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract MultiAggregateRateLimiter_updateRateLimitTokens is MultiAggregateRateLimiterSetup {
  function setUp() public virtual override {
    super.setUp();

    // Clear rate limit tokens state
    MultiAggregateRateLimiter.LocalRateLimitToken[] memory removes =
      new MultiAggregateRateLimiter.LocalRateLimitToken[](s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      removes[i] = MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[i]
      });
    }
    s_rateLimiter.updateRateLimitTokens(removes, new MultiAggregateRateLimiter.RateLimitTokenArgs[](0));
  }

  function test_UpdateRateLimitTokensSingleChain() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](2);
    adds[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[0]
      }),
      remoteToken: abi.encode(s_sourceTokens[0])
    });
    adds[1] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[1]
      }),
      remoteToken: abi.encode(s_sourceTokens[1])
    });

    for (uint256 i = 0; i < adds.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitAdded(
        CHAIN_SELECTOR_1, adds[i].remoteToken, adds[i].localTokenArgs.localToken
      );
    }

    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);

    (address[] memory localTokens, bytes[] memory remoteTokens) = s_rateLimiter.getAllRateLimitTokens(CHAIN_SELECTOR_1);

    assertEq(localTokens.length, adds.length);
    assertEq(localTokens.length, remoteTokens.length);

    for (uint256 i = 0; i < adds.length; ++i) {
      assertEq(adds[i].remoteToken, remoteTokens[i]);
      assertEq(adds[i].localTokenArgs.localToken, localTokens[i]);
    }
  }

  function test_UpdateRateLimitTokensMultipleChains() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](2);
    adds[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[0]
      }),
      remoteToken: abi.encode(s_sourceTokens[0])
    });
    adds[1] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_2,
        localToken: s_destTokens[1]
      }),
      remoteToken: abi.encode(s_sourceTokens[1])
    });

    for (uint256 i = 0; i < adds.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitAdded(
        adds[i].localTokenArgs.remoteChainSelector, adds[i].remoteToken, adds[i].localTokenArgs.localToken
      );
    }

    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);

    (address[] memory localTokensChain1, bytes[] memory remoteTokensChain1) =
      s_rateLimiter.getAllRateLimitTokens(CHAIN_SELECTOR_1);

    assertEq(localTokensChain1.length, 1);
    assertEq(localTokensChain1.length, remoteTokensChain1.length);
    assertEq(localTokensChain1[0], adds[0].localTokenArgs.localToken);
    assertEq(remoteTokensChain1[0], adds[0].remoteToken);

    (address[] memory localTokensChain2, bytes[] memory remoteTokensChain2) =
      s_rateLimiter.getAllRateLimitTokens(CHAIN_SELECTOR_2);

    assertEq(localTokensChain2.length, 1);
    assertEq(localTokensChain2.length, remoteTokensChain2.length);
    assertEq(localTokensChain2[0], adds[1].localTokenArgs.localToken);
    assertEq(remoteTokensChain2[0], adds[1].remoteToken);
  }

  function test_UpdateRateLimitTokens_AddsAndRemoves() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](2);
    adds[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[0]
      }),
      remoteToken: abi.encode(s_sourceTokens[0])
    });
    adds[1] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[1]
      }),
      remoteToken: abi.encode(s_sourceTokens[1])
    });

    MultiAggregateRateLimiter.LocalRateLimitToken[] memory removes =
      new MultiAggregateRateLimiter.LocalRateLimitToken[](1);
    removes[0] = adds[0].localTokenArgs;

    for (uint256 i = 0; i < adds.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitAdded(
        CHAIN_SELECTOR_1, adds[i].remoteToken, adds[i].localTokenArgs.localToken
      );
    }

    s_rateLimiter.updateRateLimitTokens(removes, adds);

    for (uint256 i = 0; i < removes.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitRemoved(CHAIN_SELECTOR_1, removes[i].localToken);
    }

    s_rateLimiter.updateRateLimitTokens(removes, new MultiAggregateRateLimiter.RateLimitTokenArgs[](0));

    (address[] memory localTokens, bytes[] memory remoteTokens) = s_rateLimiter.getAllRateLimitTokens(CHAIN_SELECTOR_1);

    assertEq(1, remoteTokens.length);
    assertEq(adds[1].remoteToken, remoteTokens[0]);

    assertEq(1, localTokens.length);
    assertEq(adds[1].localTokenArgs.localToken, localTokens[0]);
  }

  function test_UpdateRateLimitTokens_RemoveNonExistentToken() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](0);

    MultiAggregateRateLimiter.LocalRateLimitToken[] memory removes =
      new MultiAggregateRateLimiter.LocalRateLimitToken[](1);
    removes[0] = MultiAggregateRateLimiter.LocalRateLimitToken({
      remoteChainSelector: CHAIN_SELECTOR_1,
      localToken: s_destTokens[0]
    });

    vm.recordLogs();
    s_rateLimiter.updateRateLimitTokens(removes, adds);

    // No event since no remove occurred
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);

    (address[] memory localTokens, bytes[] memory remoteTokens) = s_rateLimiter.getAllRateLimitTokens(CHAIN_SELECTOR_1);

    assertEq(localTokens.length, 0);
    assertEq(localTokens.length, remoteTokens.length);
  }

  // Reverts

  function test_RevertWhen_ZeroSourceToken() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](1);
    adds[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: s_destTokens[0]
      }),
      remoteToken: new bytes(0)
    });

    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);
  }

  function test_RevertWhen_ZeroDestToken() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](1);
    adds[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: address(0)
      }),
      remoteToken: abi.encode(s_destTokens[0])
    });

    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);
  }

  function test_RevertWhen_ZeroDestToken_AbiEncoded() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](1);
    adds[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_1,
        localToken: address(0)
      }),
      remoteToken: abi.encode(address(0))
    });

    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);
  }

  function test_RevertWhen_NonOwner() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](4);

    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);
  }
}
