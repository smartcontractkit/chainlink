// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {FeeQuoterSetup} from "../../feeQuoter/FeeQuoterSetup.t.sol";
import {MultiAggregateRateLimiterHelper} from "../../helpers/MultiAggregateRateLimiterHelper.sol";

contract MultiAggregateRateLimiterSetup is FeeQuoterSetup {
  MultiAggregateRateLimiterHelper internal s_rateLimiter;

  address internal constant TOKEN = 0x21118E64E1fB0c487F25Dd6d3601FF6af8D32E4e;
  uint224 internal constant TOKEN_PRICE = 4e18;

  uint64 internal constant CHAIN_SELECTOR_1 = 5009297550715157269;
  uint64 internal constant CHAIN_SELECTOR_2 = 4949039107694359620;

  RateLimiter.Config internal s_rateLimiterConfig1 = RateLimiter.Config({isEnabled: true, rate: 5, capacity: 100});
  RateLimiter.Config internal s_rateLimiterConfig2 = RateLimiter.Config({isEnabled: true, rate: 10, capacity: 200});

  address internal constant MOCK_OFFRAMP = address(1111);
  address internal constant MOCK_ONRAMP = address(1112);

  address[] internal s_authorizedCallers;

  function setUp() public virtual override {
    FeeQuoterSetup.setUp();

    Internal.PriceUpdates memory priceUpdates = _getSingleTokenPriceUpdateStruct(TOKEN, TOKEN_PRICE);
    s_feeQuoter.updatePrices(priceUpdates);

    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](4);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig1
    });
    configUpdates[1] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_2,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig2
    });
    configUpdates[2] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: true,
      rateLimiterConfig: s_rateLimiterConfig1
    });
    configUpdates[3] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_2,
      isOutboundLane: true,
      rateLimiterConfig: s_rateLimiterConfig2
    });

    s_authorizedCallers = new address[](2);
    s_authorizedCallers[0] = MOCK_OFFRAMP;
    s_authorizedCallers[1] = MOCK_ONRAMP;

    s_rateLimiter = new MultiAggregateRateLimiterHelper(address(s_feeQuoter), s_authorizedCallers);
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function _assertConfigWithTokenBucketEquality(
    RateLimiter.Config memory config,
    RateLimiter.TokenBucket memory tokenBucket
  ) internal pure {
    assertEq(config.rate, tokenBucket.rate);
    assertEq(config.capacity, tokenBucket.capacity);
    assertEq(config.capacity, tokenBucket.tokens);
    assertEq(config.isEnabled, tokenBucket.isEnabled);
  }

  function _assertTokenBucketEquality(
    RateLimiter.TokenBucket memory tokenBucketA,
    RateLimiter.TokenBucket memory tokenBucketB
  ) internal pure {
    assertEq(tokenBucketA.rate, tokenBucketB.rate);
    assertEq(tokenBucketA.capacity, tokenBucketB.capacity);
    assertEq(tokenBucketA.tokens, tokenBucketB.tokens);
    assertEq(tokenBucketA.isEnabled, tokenBucketB.isEnabled);
  }

  function _generateAny2EVMMessage(
    uint64 sourceChainSelector,
    Client.EVMTokenAmount[] memory tokenAmounts
  ) internal pure returns (Client.Any2EVMMessage memory) {
    return Client.Any2EVMMessage({
      messageId: keccak256(bytes("messageId")),
      sourceChainSelector: sourceChainSelector,
      sender: abi.encode(OWNER),
      data: abi.encode(0),
      destTokenAmounts: tokenAmounts
    });
  }

  function _generateAny2EVMMessageNoTokens(
    uint64 sourceChainSelector
  ) internal pure returns (Client.Any2EVMMessage memory) {
    return _generateAny2EVMMessage(sourceChainSelector, new Client.EVMTokenAmount[](0));
  }
}
