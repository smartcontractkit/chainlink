// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";

import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";

contract MultiAggregateRateLimiter_onOutboundMessage is MultiAggregateRateLimiterSetup {
  function setUp() public virtual override {
    super.setUp();

    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitTokenArgs[](s_sourceTokens.length);
    for (uint224 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = MultiAggregateRateLimiter.RateLimitTokenArgs({
        localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
          remoteChainSelector: CHAIN_SELECTOR_1,
          localToken: s_sourceTokens[i]
        }),
        remoteToken: abi.encode(bytes20(s_destTokenBySourceToken[s_sourceTokens[i]]))
      });

      Internal.PriceUpdates memory priceUpdates =
        _getSingleTokenPriceUpdateStruct(s_sourceTokens[i], TOKEN_PRICE * (i + 1));
      s_feeQuoter.updatePrices(priceUpdates);
    }
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), tokensToAdd);
  }

  function test_ValidateMessageWithNoTokens() public {
    vm.startPrank(MOCK_ONRAMP);

    vm.recordLogs();
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessageNoTokens());

    // No consumed rate limit events
    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_onOutboundMessage_ValidateMessageWithTokens() public {
    vm.startPrank(MOCK_ONRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 3});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 1});

    // 3 tokens * TOKEN_PRICE + 1 token * (2 * TOKEN_PRICE)
    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
  }

  function test_onOutboundMessage_ValidateMessageWithDisabledRateLimitToken() public {
    MultiAggregateRateLimiter.LocalRateLimitToken[] memory removes =
      new MultiAggregateRateLimiter.LocalRateLimitToken[](1);
    removes[0] = MultiAggregateRateLimiter.LocalRateLimitToken({
      remoteChainSelector: CHAIN_SELECTOR_1,
      localToken: s_sourceTokens[1]
    });
    s_rateLimiter.updateRateLimitTokens(removes, new MultiAggregateRateLimiter.RateLimitTokenArgs[](0));

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 5});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 1});

    vm.startPrank(MOCK_ONRAMP);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
  }

  function test_onOutboundMessage_ValidateMessageWithRateLimitDisabled() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: true,
      rateLimiterConfig: s_rateLimiterConfig1
    });
    configUpdates[0].rateLimiterConfig.isEnabled = false;

    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 1000});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 50});

    vm.startPrank(MOCK_ONRAMP);
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));

    // No consumed rate limit events
    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_onOutboundMessage_ValidateMessageWithTokensOnDifferentChains() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitTokenArgs[](s_sourceTokens.length);
    for (uint224 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = MultiAggregateRateLimiter.RateLimitTokenArgs({
        localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
          remoteChainSelector: CHAIN_SELECTOR_2,
          localToken: s_sourceTokens[i]
        }),
        // Create a remote token address that is different from CHAIN_SELECTOR_1
        remoteToken: abi.encode(uint256(uint160(s_destTokenBySourceToken[s_sourceTokens[i]])) + type(uint160).max + 1)
      });
    }
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), tokensToAdd);

    vm.startPrank(MOCK_ONRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 2});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 1});

    // 2 tokens * (TOKEN_PRICE) + 1 token * (2 * TOKEN_PRICE)
    uint256 totalValue = (4 * TOKEN_PRICE) / 1e18;

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));

    // Chain 1 changed
    RateLimiter.TokenBucket memory bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 unchanged
    RateLimiter.TokenBucket memory bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, true);
    assertEq(bucketChain2.capacity, bucketChain2.tokens);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(totalValue);

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_2, _generateEVM2AnyMessage(tokenAmounts));

    // Chain 1 unchanged
    bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 changed
    bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, true);
    assertEq(bucketChain2.capacity - totalValue, bucketChain2.tokens);
  }

  function test_onOutboundMessage_ValidateMessageWithDifferentTokensOnDifferentChains() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitTokenArgs[](1);

    // Only 1 rate limited token on different chain
    tokensToAdd[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_2,
        localToken: s_sourceTokens[0]
      }),
      // Create a remote token address that is different from CHAIN_SELECTOR_1
      remoteToken: abi.encode(uint256(uint160(s_destTokenBySourceToken[s_sourceTokens[0]])) + type(uint160).max + 1)
    });
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), tokensToAdd);

    vm.startPrank(MOCK_ONRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 3});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 1});

    // 3 tokens * (TOKEN_PRICE) + 1 token * (2 * TOKEN_PRICE)
    uint256 totalValue = (5 * TOKEN_PRICE) / 1e18;

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));

    // Chain 1 changed
    RateLimiter.TokenBucket memory bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 unchanged
    RateLimiter.TokenBucket memory bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, true);
    assertEq(bucketChain2.capacity, bucketChain2.tokens);

    // 3 tokens * (TOKEN_PRICE)
    uint256 totalValue2 = (3 * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(totalValue2);

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_2, _generateEVM2AnyMessage(tokenAmounts));

    // Chain 1 unchanged
    bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 changed
    bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, true);
    assertEq(bucketChain2.capacity - totalValue2, bucketChain2.tokens);
  }

  function test_onOutboundMessage_ValidateMessageWithRateLimitReset() public {
    vm.startPrank(MOCK_ONRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 20});

    // Remaining capacity: 100 -> 20
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));

    // Cannot fit 80 rate limit value (need to wait at least 12 blocks, current capacity is 20)
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 12, 20));
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));

    // Remaining capacity: 20 -> 35 (need to wait 9 more blocks)
    vm.warp(BLOCK_TIME + 3);
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 9, 35));
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));

    // Remaining capacity: 35 -> 80 (can fit exactly 80)
    vm.warp(BLOCK_TIME + 12);
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
  }

  function test_RateLimitValueDifferentLanes() public {
    vm.pauseGasMetering();
    // start from blocktime that does not equal rate limiter init timestamp
    vm.warp(BLOCK_TIME + 1);

    // 10 (tokens) * 4 (price) * 2 (number of times) = 80 < 100 (capacity)
    uint256 numberOfTokens = 10;
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: numberOfTokens});
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    vm.startPrank(MOCK_ONRAMP);
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
    vm.pauseGasMetering();

    // Get the updated bucket status
    RateLimiter.TokenBucket memory bucket1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    RateLimiter.TokenBucket memory bucket2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    // Assert the proper value has been taken out of the bucket
    assertEq(bucket1.capacity - value, bucket1.tokens);
    // Inbound lane should remain unchanged
    assertEq(bucket2.capacity, bucket2.tokens);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
    vm.pauseGasMetering();

    bucket1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    bucket2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    // Inbound lane should remain unchanged
    assertEq(bucket1.capacity - value, bucket1.tokens);
    assertEq(bucket2.capacity - value, bucket2.tokens);
  }

  // Reverts

  function test_RevertWhen_onOutboundMessage_ValidateMessageWithRateLimitExceeded() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 80});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 30});

    uint256 totalValue = (80 * TOKEN_PRICE + 2 * (30 * TOKEN_PRICE)) / 1e18;
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueMaxCapacityExceeded.selector, 100, totalValue));
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
  }

  function test_RevertWhen_onOutboundMessage_ValidateMessageFromUnauthorizedCaller() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessageNoTokens());
  }

  function _generateEVM2AnyMessage(
    Client.EVMTokenAmount[] memory tokenAmounts
  ) public view returns (Client.EVM2AnyMessage memory) {
    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }

  function _generateEVM2AnyMessageNoTokens() internal view returns (Client.EVM2AnyMessage memory) {
    return _generateEVM2AnyMessage(new Client.EVMTokenAmount[](0));
  }
}
