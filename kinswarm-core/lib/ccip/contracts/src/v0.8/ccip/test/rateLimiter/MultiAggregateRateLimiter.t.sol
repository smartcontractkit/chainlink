// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {MultiAggregateRateLimiter} from "../../MultiAggregateRateLimiter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {BaseTest} from "../BaseTest.t.sol";

import {FeeQuoterSetup} from "../feeQuoter/FeeQuoterSetup.t.sol";
import {MultiAggregateRateLimiterHelper} from "../helpers/MultiAggregateRateLimiterHelper.sol";
import {stdError} from "forge-std/Test.sol";
import {Vm} from "forge-std/Vm.sol";

contract MultiAggregateRateLimiterSetup is BaseTest, FeeQuoterSetup {
  MultiAggregateRateLimiterHelper internal s_rateLimiter;

  address internal immutable TOKEN = 0x21118E64E1fB0c487F25Dd6d3601FF6af8D32E4e;
  uint224 internal constant TOKEN_PRICE = 4e18;

  uint64 internal constant CHAIN_SELECTOR_1 = 5009297550715157269;
  uint64 internal constant CHAIN_SELECTOR_2 = 4949039107694359620;

  RateLimiter.Config internal RATE_LIMITER_CONFIG_1 = RateLimiter.Config({isEnabled: true, rate: 5, capacity: 100});
  RateLimiter.Config internal RATE_LIMITER_CONFIG_2 = RateLimiter.Config({isEnabled: true, rate: 10, capacity: 200});

  address internal immutable MOCK_OFFRAMP = address(1111);
  address internal immutable MOCK_ONRAMP = address(1112);

  address[] internal s_authorizedCallers;

  function setUp() public virtual override(BaseTest, FeeQuoterSetup) {
    BaseTest.setUp();
    FeeQuoterSetup.setUp();

    Internal.PriceUpdates memory priceUpdates = _getSingleTokenPriceUpdateStruct(TOKEN, TOKEN_PRICE);
    s_feeQuoter.updatePrices(priceUpdates);

    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](4);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    configUpdates[1] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_2,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });
    configUpdates[2] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    configUpdates[3] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_2,
      isOutboundLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
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

contract MultiAggregateRateLimiter_constructor is MultiAggregateRateLimiterSetup {
  function test_ConstructorNoAuthorizedCallers_Success() public {
    address[] memory authorizedCallers = new address[](0);

    vm.recordLogs();
    s_rateLimiter = new MultiAggregateRateLimiterHelper(address(s_feeQuoter), authorizedCallers);

    // FeeQuoterSet
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_feeQuoter), s_rateLimiter.getFeeQuoter());
  }

  function test_Constructor_Success() public {
    address[] memory authorizedCallers = new address[](2);
    authorizedCallers[0] = MOCK_OFFRAMP;
    authorizedCallers[1] = MOCK_ONRAMP;

    vm.expectEmit();
    emit MultiAggregateRateLimiter.FeeQuoterSet(address(s_feeQuoter));

    s_rateLimiter = new MultiAggregateRateLimiterHelper(address(s_feeQuoter), authorizedCallers);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_feeQuoter), s_rateLimiter.getFeeQuoter());
    assertEq(s_rateLimiter.typeAndVersion(), "MultiAggregateRateLimiter 1.6.0-dev");
  }
}

contract MultiAggregateRateLimiter_setFeeQuoter is MultiAggregateRateLimiterSetup {
  function test_Owner_Success() public {
    address newAddress = address(42);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.FeeQuoterSet(newAddress);

    s_rateLimiter.setFeeQuoter(newAddress);
    assertEq(newAddress, s_rateLimiter.getFeeQuoter());
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(bytes("Only callable by owner"));

    s_rateLimiter.setFeeQuoter(STRANGER);
  }

  function test_ZeroAddress_Revert() public {
    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);
    s_rateLimiter.setFeeQuoter(address(0));
  }
}

contract MultiAggregateRateLimiter_getTokenBucket is MultiAggregateRateLimiterSetup {
  function test_GetTokenBucket_Success() public view {
    RateLimiter.TokenBucket memory bucketInbound = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_1, bucketInbound);
    assertEq(BLOCK_TIME, bucketInbound.lastUpdated);

    RateLimiter.TokenBucket memory bucketOutbound = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_1, bucketOutbound);
    assertEq(BLOCK_TIME, bucketOutbound.lastUpdated);
  }

  function test_Refill_Success() public {
    RATE_LIMITER_CONFIG_1.capacity = RATE_LIMITER_CONFIG_1.capacity * 2;

    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    assertEq(RATE_LIMITER_CONFIG_1.rate, bucket.rate);
    assertEq(RATE_LIMITER_CONFIG_1.capacity, bucket.capacity);
    assertEq(RATE_LIMITER_CONFIG_1.capacity / 2, bucket.tokens);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    uint256 warpTime = 4;
    vm.warp(BLOCK_TIME + warpTime);

    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    assertEq(RATE_LIMITER_CONFIG_1.rate, bucket.rate);
    assertEq(RATE_LIMITER_CONFIG_1.capacity, bucket.capacity);
    assertEq(RATE_LIMITER_CONFIG_1.capacity / 2 + warpTime * RATE_LIMITER_CONFIG_1.rate, bucket.tokens);
    assertEq(BLOCK_TIME + warpTime, bucket.lastUpdated);

    vm.warp(BLOCK_TIME + warpTime * 100);

    // Bucket overflow
    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(RATE_LIMITER_CONFIG_1.capacity, bucket.tokens);
  }

  // Reverts

  function test_TimeUnderflow_Revert() public {
    vm.warp(BLOCK_TIME - 1);

    vm.expectRevert(stdError.arithmeticError);
    s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
  }
}

contract MultiAggregateRateLimiter_applyRateLimiterConfigUpdates is MultiAggregateRateLimiterSetup {
  function test_ZeroConfigs_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](0);

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_SingleConfig_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, false, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    RateLimiter.TokenBucket memory bucket1 =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);
    _assertConfigWithTokenBucketEquality(configUpdates[0].rateLimiterConfig, bucket1);
    assertEq(BLOCK_TIME, bucket1.lastUpdated);
  }

  function test_SingleConfigOutbound_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, true, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    RateLimiter.TokenBucket memory bucket1 =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, true);
    _assertConfigWithTokenBucketEquality(configUpdates[0].rateLimiterConfig, bucket1);
    assertEq(BLOCK_TIME, bucket1.lastUpdated);
  }

  function test_MultipleConfigs_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](5);

    for (uint64 i; i < configUpdates.length; ++i) {
      configUpdates[i] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
        remoteChainSelector: CHAIN_SELECTOR_1 + i + 1,
        isOutboundLane: i % 2 == 0 ? false : true,
        rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 5 + i, capacity: 100 + i})
      });

      vm.expectEmit();
      emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
        configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane, configUpdates[i].rateLimiterConfig
      );
    }

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, configUpdates.length);

    for (uint256 i; i < configUpdates.length; ++i) {
      RateLimiter.TokenBucket memory bucket =
        s_rateLimiter.currentRateLimiterState(configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane);
      _assertConfigWithTokenBucketEquality(configUpdates[i].rateLimiterConfig, bucket);
      assertEq(BLOCK_TIME, bucket.lastUpdated);
    }
  }

  function test_MultipleConfigsBothLanes_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](2);

    for (uint64 i; i < configUpdates.length; ++i) {
      configUpdates[i] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
        remoteChainSelector: CHAIN_SELECTOR_1 + 1,
        isOutboundLane: i % 2 == 0 ? false : true,
        rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 5 + i, capacity: 100 + i})
      });

      vm.expectEmit();
      emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
        configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane, configUpdates[i].rateLimiterConfig
      );
    }

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, configUpdates.length);

    for (uint256 i; i < configUpdates.length; ++i) {
      RateLimiter.TokenBucket memory bucket =
        s_rateLimiter.currentRateLimiterState(configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane);
      _assertConfigWithTokenBucketEquality(configUpdates[i].rateLimiterConfig, bucket);
      assertEq(BLOCK_TIME, bucket.lastUpdated);
    }
  }

  function test_UpdateExistingConfig_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });

    RateLimiter.TokenBucket memory bucket1 =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);

    // Capacity equals tokens
    assertEq(bucket1.capacity, bucket1.tokens);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, false, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    vm.warp(BLOCK_TIME + 1);
    bucket1 = s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);
    assertEq(BLOCK_TIME + 1, bucket1.lastUpdated);

    // Tokens < capacity since capacity doubled
    assertTrue(bucket1.capacity != bucket1.tokens);

    // Outbound lane config remains unchanged
    _assertConfigWithTokenBucketEquality(
      RATE_LIMITER_CONFIG_1, s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true)
    );
  }

  function test_UpdateExistingConfigWithNoDifference_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    RateLimiter.TokenBucket memory bucketPreUpdate =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, false, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    vm.warp(BLOCK_TIME + 1);
    RateLimiter.TokenBucket memory bucketPostUpdate =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);
    _assertTokenBucketEquality(bucketPreUpdate, bucketPostUpdate);
    assertEq(BLOCK_TIME + 1, bucketPostUpdate.lastUpdated);
  }

  // Reverts
  function test_ZeroChainSelector_Revert() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: 0,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    vm.expectRevert(MultiAggregateRateLimiter.ZeroChainSelectorNotAllowed.selector);
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_OnlyCallableByOwner_Revert() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    vm.startPrank(STRANGER);

    vm.expectRevert(bytes("Only callable by owner"));
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }
}

contract MultiAggregateRateLimiter_getTokenValue is MultiAggregateRateLimiterSetup {
  function test_GetTokenValue_Success() public view {
    uint256 numberOfTokens = 10;
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: TOKEN, amount: 10});
    uint256 value = s_rateLimiter.getTokenValue(tokenAmount);
    assertEq(value, (numberOfTokens * TOKEN_PRICE) / 1e18);
  }

  // Reverts
  function test_NoTokenPrice_Reverts() public {
    address tokenWithNoPrice = makeAddr("Token with no price");
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: tokenWithNoPrice, amount: 10});

    vm.expectRevert(abi.encodeWithSelector(MultiAggregateRateLimiter.PriceNotFoundForToken.selector, tokenWithNoPrice));
    s_rateLimiter.getTokenValue(tokenAmount);
  }
}

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

  function test_UpdateRateLimitTokensSingleChain_Success() public {
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

  function test_UpdateRateLimitTokensMultipleChains_Success() public {
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

  function test_UpdateRateLimitTokens_AddsAndRemoves_Success() public {
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

  function test_UpdateRateLimitTokens_RemoveNonExistentToken_Success() public {
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

  function test_ZeroSourceToken_Revert() public {
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

  function test_ZeroDestToken_Revert() public {
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

  function test_NonOwner_Revert() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory adds = new MultiAggregateRateLimiter.RateLimitTokenArgs[](4);

    vm.startPrank(STRANGER);

    vm.expectRevert(bytes("Only callable by owner"));
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), adds);
  }
}

contract MultiAggregateRateLimiter_onInboundMessage is MultiAggregateRateLimiterSetup {
  address internal immutable MOCK_RECEIVER = address(1113);

  function setUp() public virtual override {
    super.setUp();

    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitTokenArgs[](s_sourceTokens.length);
    for (uint224 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = MultiAggregateRateLimiter.RateLimitTokenArgs({
        localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
          remoteChainSelector: CHAIN_SELECTOR_1,
          localToken: s_destTokens[i]
        }),
        remoteToken: abi.encode(s_sourceTokens[i])
      });

      Internal.PriceUpdates memory priceUpdates =
        _getSingleTokenPriceUpdateStruct(s_destTokens[i], TOKEN_PRICE * (i + 1));
      s_feeQuoter.updatePrices(priceUpdates);
    }
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), tokensToAdd);
  }

  function test_ValidateMessageWithNoTokens_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    vm.recordLogs();
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessageNoTokens(CHAIN_SELECTOR_1));

    // No consumed rate limit events
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_ValidateMessageWithTokens_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 3});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    // 3 tokens * TOKEN_PRICE + 1 token * (2 * TOKEN_PRICE)
    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  function test_ValidateMessageWithDisabledRateLimitToken_Success() public {
    MultiAggregateRateLimiter.LocalRateLimitToken[] memory removes =
      new MultiAggregateRateLimiter.LocalRateLimitToken[](1);
    removes[0] = MultiAggregateRateLimiter.LocalRateLimitToken({
      remoteChainSelector: CHAIN_SELECTOR_1,
      localToken: s_destTokens[1]
    });
    s_rateLimiter.updateRateLimitTokens(removes, new MultiAggregateRateLimiter.RateLimitTokenArgs[](0));

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 5});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    vm.startPrank(MOCK_OFFRAMP);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  function test_ValidateMessageWithRateLimitDisabled_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    configUpdates[0].rateLimiterConfig.isEnabled = false;

    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 1000});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 50});

    vm.startPrank(MOCK_OFFRAMP);
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // No consumed rate limit events
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_ValidateMessageWithTokensOnDifferentChains_Success() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitTokenArgs[](s_sourceTokens.length);
    for (uint224 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = MultiAggregateRateLimiter.RateLimitTokenArgs({
        localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
          remoteChainSelector: CHAIN_SELECTOR_2,
          localToken: s_destTokens[i]
        }),
        // Create a remote token address that is different from CHAIN_SELECTOR_1
        remoteToken: abi.encode(uint256(uint160(s_sourceTokens[i])) + type(uint160).max + 1)
      });
    }
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), tokensToAdd);

    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 2});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    // 2 tokens * (TOKEN_PRICE) + 1 token * (2 * TOKEN_PRICE)
    uint256 totalValue = (4 * TOKEN_PRICE) / 1e18;

    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Chain 1 changed
    RateLimiter.TokenBucket memory bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 unchanged
    RateLimiter.TokenBucket memory bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    assertEq(bucketChain2.capacity, bucketChain2.tokens);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(totalValue);

    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_2, tokenAmounts));

    // Chain 1 unchanged
    bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 changed
    bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    assertEq(bucketChain2.capacity - totalValue, bucketChain2.tokens);
  }

  function test_ValidateMessageWithDifferentTokensOnDifferentChains_Success() public {
    MultiAggregateRateLimiter.RateLimitTokenArgs[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitTokenArgs[](1);

    // Only 1 rate limited token on different chain
    tokensToAdd[0] = MultiAggregateRateLimiter.RateLimitTokenArgs({
      localTokenArgs: MultiAggregateRateLimiter.LocalRateLimitToken({
        remoteChainSelector: CHAIN_SELECTOR_2,
        localToken: s_destTokens[0]
      }),
      // Create a remote token address that is different from CHAIN_SELECTOR_1
      remoteToken: abi.encode(uint256(uint160(s_sourceTokens[0])) + type(uint160).max + 1)
    });
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.LocalRateLimitToken[](0), tokensToAdd);

    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 3});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    // 3 tokens * (TOKEN_PRICE) + 1 token * (2 * TOKEN_PRICE)
    uint256 totalValue = (5 * TOKEN_PRICE) / 1e18;

    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Chain 1 changed
    RateLimiter.TokenBucket memory bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 unchanged
    RateLimiter.TokenBucket memory bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    assertEq(bucketChain2.capacity, bucketChain2.tokens);

    // 3 tokens * (TOKEN_PRICE)
    uint256 totalValue2 = (3 * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(totalValue2);

    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_2, tokenAmounts));

    // Chain 1 unchanged
    bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 changed
    bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    assertEq(bucketChain2.capacity - totalValue2, bucketChain2.tokens);
  }

  function test_ValidateMessageWithRateLimitReset_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 20});

    // Remaining capacity: 100 -> 20
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Cannot fit 80 rate limit value (need to wait at least 12 blocks, current capacity is 20)
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 12, 20));
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Remaining capacity: 20 -> 35 (need to wait 9 more blocks)
    vm.warp(BLOCK_TIME + 3);
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 9, 35));
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Remaining capacity: 35 -> 80 (can fit exactly 80)
    vm.warp(BLOCK_TIME + 12);
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  // Reverts

  function test_ValidateMessageWithRateLimitExceeded_Revert() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 80});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 30});

    uint256 totalValue = (80 * TOKEN_PRICE + 2 * (30 * TOKEN_PRICE)) / 1e18;
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueMaxCapacityExceeded.selector, 100, totalValue));
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  function test_ValidateMessageFromUnauthorizedCaller_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_rateLimiter.onInboundMessage(_generateAny2EVMMessageNoTokens(CHAIN_SELECTOR_1));
  }
}

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

  function test_ValidateMessageWithNoTokens_Success() public {
    vm.startPrank(MOCK_ONRAMP);

    vm.recordLogs();
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessageNoTokens());

    // No consumed rate limit events
    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_onOutboundMessage_ValidateMessageWithTokens_Success() public {
    vm.startPrank(MOCK_ONRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 3});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 1});

    // 3 tokens * TOKEN_PRICE + 1 token * (2 * TOKEN_PRICE)
    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
  }

  function test_onOutboundMessage_ValidateMessageWithDisabledRateLimitToken_Success() public {
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

  function test_onOutboundMessage_ValidateMessageWithRateLimitDisabled_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
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

  function test_onOutboundMessage_ValidateMessageWithTokensOnDifferentChains_Success() public {
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

  function test_onOutboundMessage_ValidateMessageWithDifferentTokensOnDifferentChains_Success() public {
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

  function test_onOutboundMessage_ValidateMessageWithRateLimitReset_Success() public {
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

  function test_RateLimitValueDifferentLanes_Success() public {
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

  function test_onOutboundMessage_ValidateMessageWithRateLimitExceeded_Revert() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: 80});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_sourceTokens[1], amount: 30});

    uint256 totalValue = (80 * TOKEN_PRICE + 2 * (30 * TOKEN_PRICE)) / 1e18;
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueMaxCapacityExceeded.selector, 100, totalValue));
    s_rateLimiter.onOutboundMessage(CHAIN_SELECTOR_1, _generateEVM2AnyMessage(tokenAmounts));
  }

  function test_onOutboundMessage_ValidateMessageFromUnauthorizedCaller_Revert() public {
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
