// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {AggregateRateLimiter} from "../../AggregateRateLimiter.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import "./EVM2EVMOnRampSetup.t.sol";

contract EVM2EVMOnRamp_constructor is EVM2EVMOnRampSetup {
  function test_Constructor_Success() public {
    EVM2EVMOnRamp.StaticConfig memory staticConfig = EVM2EVMOnRamp.StaticConfig({
      linkToken: s_sourceTokens[0],
      chainSelector: SOURCE_CHAIN_SELECTOR,
      destChainSelector: DEST_CHAIN_SELECTOR,
      defaultTxGasLimit: GAS_LIMIT,
      maxNopFeesJuels: MAX_NOP_FEES_JUELS,
      prevOnRamp: address(0),
      rmnProxy: address(s_mockRMN),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig =
      generateDynamicOnRampConfig(address(s_sourceRouter), address(s_feeQuoter));

    vm.expectEmit();
    emit EVM2EVMOnRamp.ConfigSet(staticConfig, dynamicConfig);

    s_onRamp = new EVM2EVMOnRampHelper(
      staticConfig,
      dynamicConfig,
      _getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getNopsAndWeights()
    );

    EVM2EVMOnRamp.StaticConfig memory gotStaticConfig = s_onRamp.getStaticConfig();
    assertEq(staticConfig.linkToken, gotStaticConfig.linkToken);
    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.destChainSelector, gotStaticConfig.destChainSelector);
    assertEq(staticConfig.defaultTxGasLimit, gotStaticConfig.defaultTxGasLimit);
    assertEq(staticConfig.maxNopFeesJuels, gotStaticConfig.maxNopFeesJuels);
    assertEq(staticConfig.prevOnRamp, gotStaticConfig.prevOnRamp);
    assertEq(staticConfig.rmnProxy, gotStaticConfig.rmnProxy);

    EVM2EVMOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(dynamicConfig.router, gotDynamicConfig.router);
    assertEq(dynamicConfig.maxNumberOfTokensPerMsg, gotDynamicConfig.maxNumberOfTokensPerMsg);
    assertEq(dynamicConfig.destGasOverhead, gotDynamicConfig.destGasOverhead);
    assertEq(dynamicConfig.destGasPerPayloadByte, gotDynamicConfig.destGasPerPayloadByte);
    assertEq(dynamicConfig.priceRegistry, gotDynamicConfig.priceRegistry);
    assertEq(dynamicConfig.maxDataBytes, gotDynamicConfig.maxDataBytes);
    assertEq(dynamicConfig.maxPerMsgGasLimit, gotDynamicConfig.maxPerMsgGasLimit);

    // Initial values
    assertEq("EVM2EVMOnRamp 1.5.0", s_onRamp.typeAndVersion());
    assertEq(OWNER, s_onRamp.owner());
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber());
  }
}

contract EVM2EVMOnRamp_payNops_fuzz is EVM2EVMOnRampSetup {
  function test_Fuzz_NopPayNops_Success(uint96 nopFeesJuels) public {
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    // To avoid NoFeesToPay
    vm.assume(nopFeesJuels > weightsTotal);
    vm.assume(nopFeesJuels < MAX_NOP_FEES_JUELS);

    // Set Nop fee juels
    deal(s_sourceFeeToken, address(s_onRamp), nopFeesJuels);
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);

    vm.startPrank(OWNER);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (totalJuels * nopsAndWeights[i].weight) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }
}

contract EVM2EVMNopsFeeSetup is EVM2EVMOnRampSetup {
  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    // Since we'll mostly be testing for valid calls from the router we'll
    // mock all calls to be originating from the router and re-mock in
    // tests that require failure.
    vm.startPrank(address(s_sourceRouter));

    uint256 feeAmount = 1234567890;
    uint256 numberOfMessages = 5;

    // Send a bunch of messages, increasing the juels in the contract
    for (uint256 i = 0; i < numberOfMessages; ++i) {
      IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), feeAmount, OWNER);
    }

    assertEq(s_onRamp.getNopFeesJuels(), feeAmount * numberOfMessages);
    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), feeAmount * numberOfMessages);
  }
}

contract EVM2EVMOnRamp_payNops is EVM2EVMNopsFeeSetup {
  function test_OwnerPayNops_Success() public {
    vm.startPrank(OWNER);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (nopsAndWeights[i].weight * totalJuels) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }

  function test_AdminPayNops_Success() public {
    vm.startPrank(ADMIN);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (nopsAndWeights[i].weight * totalJuels) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }

  function test_NopPayNops_Success() public {
    vm.startPrank(getNopsAndWeights()[0].nop);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (nopsAndWeights[i].weight * totalJuels) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }

  function test_PayNopsSuccessAfterSetNops() public {
    vm.startPrank(OWNER);

    // set 2 nops, 1 from previous, 1 new
    address prevNop = getNopsAndWeights()[0].nop;
    address newNop = STRANGER;
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](2);
    nopsAndWeights[0] = EVM2EVMOnRamp.NopAndWeight({nop: prevNop, weight: 1});
    nopsAndWeights[1] = EVM2EVMOnRamp.NopAndWeight({nop: newNop, weight: 1});
    s_onRamp.setNops(nopsAndWeights);

    // refill OnRamp nops fees
    vm.startPrank(address(s_sourceRouter));
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), feeAmount, OWNER);

    vm.startPrank(newNop);
    uint256 prevNopBalance = IERC20(s_sourceFeeToken).balanceOf(prevNop);
    uint256 totalJuels = s_onRamp.getNopFeesJuels();

    s_onRamp.payNops();

    assertEq(totalJuels / 2 + prevNopBalance, IERC20(s_sourceFeeToken).balanceOf(prevNop));
    assertEq(totalJuels / 2, IERC20(s_sourceFeeToken).balanceOf(newNop));
  }

  // Reverts

  function test_InsufficientBalance_Revert() public {
    vm.startPrank(address(s_onRamp));
    IERC20(s_sourceFeeToken).transfer(OWNER, IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)));
    vm.startPrank(OWNER);
    vm.expectRevert(EVM2EVMOnRamp.InsufficientBalance.selector);
    s_onRamp.payNops();
  }

  function test_WrongPermissions_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdminOrNop.selector);
    s_onRamp.payNops();
  }

  function test_NoFeesToPay_Revert() public {
    vm.startPrank(OWNER);
    s_onRamp.payNops();
    vm.expectRevert(EVM2EVMOnRamp.NoFeesToPay.selector);
    s_onRamp.payNops();
  }

  function test_NoNopsToPay_Revert() public {
    vm.startPrank(OWNER);
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](0);
    s_onRamp.setNops(nopsAndWeights);
    vm.expectRevert(EVM2EVMOnRamp.NoNopsToPay.selector);
    s_onRamp.payNops();
  }
}

contract EVM2EVMOnRamp_linkAvailableForPayment is EVM2EVMNopsFeeSetup {
  function test_LinkAvailableForPayment_Success() public {
    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    uint256 linkBalance = IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp));

    assertEq(int256(linkBalance - totalJuels), s_onRamp.linkAvailableForPayment());

    vm.startPrank(OWNER);
    s_onRamp.payNops();

    assertEq(int256(linkBalance - totalJuels), s_onRamp.linkAvailableForPayment());
  }

  function test_InsufficientLinkBalance_Success() public {
    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    uint256 linkBalance = IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp));

    vm.startPrank(address(s_onRamp));

    uint256 linkRemaining = 1;
    IERC20(s_sourceFeeToken).transfer(OWNER, linkBalance - linkRemaining);

    vm.startPrank(STRANGER);
    assertEq(int256(linkRemaining) - int256(totalJuels), s_onRamp.linkAvailableForPayment());
  }
}

contract EVM2EVMOnRamp_forwardFromRouter is EVM2EVMOnRampSetup {
  struct LegacyExtraArgs {
    uint256 gasLimit;
    bool strict;
  }

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];
    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));

    // Since we'll mostly be testing for valid calls from the router we'll
    // mock all calls to be originating from the router and re-mock in
    // tests that require failure.
    vm.startPrank(address(s_sourceRouter));
  }

  function test_ForwardFromRouterSuccessCustomExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterSuccessLegacyExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs =
      abi.encodeWithSelector(Client.EVM_EXTRA_ARGS_V1_TAG, LegacyExtraArgs({gasLimit: GAS_LIMIT * 2, strict: true}));
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    // We expect the message to be emitted with strict = false.
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterExtraArgsV2_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: true})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterExtraArgsV2AllowOutOfOrderTrue_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: true})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_Fuzz_EnforceOutOfOrder(bool enforce, bool allowOutOfOrderExecution) public {
    // Update dynamic config to enforce allowOutOfOrderExecution = defaultVal.
    vm.stopPrank();
    vm.startPrank(OWNER);
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    s_onRamp.setDynamicConfig(
      EVM2EVMOnRamp.DynamicConfig({
        router: dynamicConfig.router,
        maxNumberOfTokensPerMsg: dynamicConfig.maxNumberOfTokensPerMsg,
        destGasOverhead: dynamicConfig.destGasOverhead,
        destGasPerPayloadByte: dynamicConfig.destGasPerPayloadByte,
        destDataAvailabilityOverheadGas: dynamicConfig.destDataAvailabilityOverheadGas,
        destGasPerDataAvailabilityByte: dynamicConfig.destGasPerDataAvailabilityByte,
        destDataAvailabilityMultiplierBps: dynamicConfig.destDataAvailabilityMultiplierBps,
        priceRegistry: dynamicConfig.priceRegistry,
        maxDataBytes: dynamicConfig.maxDataBytes,
        maxPerMsgGasLimit: dynamicConfig.maxPerMsgGasLimit,
        defaultTokenFeeUSDCents: dynamicConfig.defaultTokenFeeUSDCents,
        defaultTokenDestGasOverhead: dynamicConfig.defaultTokenDestGasOverhead,
        enforceOutOfOrder: enforce
      })
    );
    vm.stopPrank();

    vm.startPrank(address(s_sourceRouter));
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG,
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: allowOutOfOrderExecution})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    if (enforce) {
      // If enforcement is on, only true should be allowed.
      if (allowOutOfOrderExecution) {
        vm.expectEmit();
        emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));
        s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
      } else {
        vm.expectRevert(EVM2EVMOnRamp.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
        s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
      }
    } else {
      // no enforcement should allow any value.
      vm.expectEmit();
      emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
    }
  }

  function test_ShouldIncrementSeqNumAndNonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_onRamp.getSenderNonce(OWNER);
      uint64 sequenceNumberBefore = s_onRamp.getSequenceNumber();

      vm.expectEmit();
      emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_onRamp.getSenderNonce(OWNER);
      uint64 sequenceNumberAfter = s_onRamp.getSequenceNumber();
      assertEq(nonceAfter, nonceBefore + 1);
      assertEq(sequenceNumberAfter, sequenceNumberBefore + 1);
    }
  }

  function test_ShouldIncrementNonceOnlyOnOrdered_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: true})
    );

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_onRamp.getSenderNonce(OWNER);
      uint64 sequenceNumberBefore = s_onRamp.getSequenceNumber();

      vm.expectEmit();
      emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_onRamp.getSenderNonce(OWNER);
      uint64 sequenceNumberAfter = s_onRamp.getSequenceNumber();
      assertEq(nonceAfter, nonceBefore);
      assertEq(sequenceNumberAfter, sequenceNumberBefore + 1);
    }
  }

  function test_forwardFromRouter_ShouldStoreLinkFees_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), feeAmount);
    assertEq(s_onRamp.getNopFeesJuels(), feeAmount);
  }

  function test_ShouldStoreNonLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceTokens[1]).transferFrom(OWNER, address(s_onRamp), feeAmount);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceTokens[1]).balanceOf(address(s_onRamp)), feeAmount);

    // Calculate conversion done by prices contract
    uint256 feeTokenPrice = s_feeQuoter.getTokenPrice(s_sourceTokens[1]).value;
    uint256 linkTokenPrice = s_feeQuoter.getTokenPrice(s_sourceFeeToken).value;
    uint256 conversionRate = (feeTokenPrice * 1e18) / linkTokenPrice;
    uint256 expectedJuels = (feeAmount * conversionRate) / 1e18;

    assertEq(s_onRamp.getNopFeesJuels(), expectedJuels);
  }

  // Make sure any valid sender, receiver and feeAmount can be handled.
  // @TODO Temporarily setting lower fuzz run as 256 triggers snapshot gas off by 1 error.
  // https://github.com/foundry-rs/foundry/issues/5689
  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function test_Fuzz_ForwardFromRouter_Success(address originalSender, address receiver, uint96 feeTokenAmount) public {
    // To avoid RouterMustSetOriginalSender
    vm.assume(originalSender != address(0));
    vm.assume(uint160(receiver) >= Internal.PRECOMPILE_SPACE);
    vm.assume(feeTokenAmount <= MAX_NOP_FEES_JUELS);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(receiver);

    // Make sure the tokens are in the contract
    deal(s_sourceFeeToken, address(s_onRamp), feeTokenAmount);

    Internal.EVM2EVMMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeTokenAmount, originalSender);

    vm.expectEmit(false, false, false, true);
    emit EVM2EVMOnRamp.CCIPSendRequested(expectedEvent);

    // Assert the message Id is correct
    assertEq(
      expectedEvent.messageId, s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeTokenAmount, originalSender)
    );
    // Assert the fee token amount is correctly assigned to the nop fee pool
    assertEq(feeTokenAmount, s_onRamp.getNopFeesJuels());
  }

  function test_OverValueWithARLOff_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 10;
    message.tokenAmounts[0].token = s_sourceTokens[0];

    IERC20(s_sourceTokens[0]).approve(address(s_onRamp), 10);

    vm.startPrank(OWNER);
    // Set a high price to trip the ARL
    uint224 tokenPrice = 3 ** 128;
    Internal.PriceUpdates memory priceUpdates = _getSingleTokenPriceUpdateStruct(s_sourceTokens[0], tokenPrice);
    s_feeQuoter.updatePrices(priceUpdates);
    vm.startPrank(address(s_sourceRouter));

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector,
        _getOutboundRateLimiterConfig().capacity,
        (message.tokenAmounts[0].amount * tokenPrice) / 1e18
      )
    );
    // Expect to fail from ARL
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Configure ARL off for token
    EVM2EVMOnRamp.TokenTransferFeeConfig memory tokenTransferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(s_sourceTokens[0]);
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: s_sourceTokens[0],
      minFeeUSDCents: tokenTransferFeeConfig.minFeeUSDCents,
      maxFeeUSDCents: tokenTransferFeeConfig.maxFeeUSDCents,
      deciBps: tokenTransferFeeConfig.deciBps,
      destGasOverhead: tokenTransferFeeConfig.destGasOverhead,
      destBytesOverhead: tokenTransferFeeConfig.destBytesOverhead,
      aggregateRateLimitEnabled: false
    });
    vm.startPrank(OWNER);
    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs, new address[](0));

    vm.startPrank(address(s_sourceRouter));
    // Expect the call now succeeds
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_forwardFromRouter_correctSourceTokenData_Success() public {
    Client.EVM2AnyMessage memory message = _generateTokenMessage();

    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      address token = message.tokenAmounts[i].token;
      deal(token, s_sourcePoolByToken[token], message.tokenAmounts[i].amount * 2);
    }

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount * 2);

    Internal.EVM2EVMMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeAmount, OWNER);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(expectedEvent);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    // Same message, but we change the onchain config which should be reflected in the event.
    // We get the event before changing the onRamp config, as the event generation code uses the current
    // onramp to generate the event. This test checks if it does so correctly.
    expectedEvent = _messageToEvent(message, 2, 2, feeAmount, OWNER);

    uint256 tokenIndexToChange = 1;
    address changedToken = message.tokenAmounts[tokenIndexToChange].token;

    // Set token config to change the destGasOverhead
    vm.startPrank(OWNER);
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: changedToken,
      minFeeUSDCents: 0,
      maxFeeUSDCents: 100,
      deciBps: 0,
      destGasOverhead: 1_000_111,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      aggregateRateLimitEnabled: false
    });
    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs, new address[](0));

    vm.startPrank(address(s_sourceRouter));

    expectedEvent.sourceTokenData[tokenIndexToChange] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[changedToken]),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[changedToken]),
        extraData: "",
        // The user will be billed either the default or the override, so we send the exact amount that we billed for
        // to the destination chain to be used for the token releaseOrMint and transfer.
        destGasAmount: tokenTransferFeeConfigArgs[0].destGasOverhead
      })
    );
    // Update the hash because we manually changed sourceTokenData
    expectedEvent.messageId = Internal._hash(expectedEvent, s_metadataHash);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(expectedEvent);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  // Reverts

  function test_Paused_Revert() public {
    // We pause by disabling the whitelist
    vm.stopPrank();
    vm.startPrank(OWNER);
    address router = address(0);
    s_onRamp.setDynamicConfig(generateDynamicOnRampConfig(router, address(2)));
    vm.expectRevert(EVM2EVMOnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_InvalidExtraArgsTag_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = bytes("bad args");

    vm.expectRevert(EVM2EVMOnRamp.InvalidExtraArgsTag.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(EVM2EVMOnRamp.CursedByRMN.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_Permissions_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.expectRevert(EVM2EVMOnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_OriginalSender_Revert() public {
    vm.expectRevert(EVM2EVMOnRamp.RouterMustSetOriginalSender.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, address(0));
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(EVM2EVMOnRamp.UnsupportedNumberOfTokens.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function test_CannotSendZeroTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 0;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    vm.expectRevert(EVM2EVMOnRamp.CannotSendZeroTokens.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function test_UnsupportedToken_Revert() public {
    address wrongToken = address(1);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].token = wrongToken;
    message.tokenAmounts[0].amount = 1;

    // We need to set the price of this new token to be able to reach
    // the proper revert point. This must be called by the owner.
    vm.stopPrank();
    vm.startPrank(OWNER);

    Internal.PriceUpdates memory priceUpdates = _getSingleTokenPriceUpdateStruct(wrongToken, 1);
    s_feeQuoter.updatePrices(priceUpdates);

    // Change back to the router
    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, wrongToken));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_MaxCapacityExceeded_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 2 ** 128;
    message.tokenAmounts[0].token = s_sourceTokens[0];

    IERC20(s_sourceTokens[0]).approve(address(s_onRamp), 2 ** 128);

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector,
        _getOutboundRateLimiterConfig().capacity,
        (message.tokenAmounts[0].amount * s_sourceTokenPrices[0]) / 1e18
      )
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_PriceNotFoundForToken_Revert() public {
    // Set token price to 0
    vm.stopPrank();
    vm.startPrank(OWNER);
    s_feeQuoter.updatePrices(_getSingleTokenPriceUpdateStruct(CUSTOM_TOKEN, 0));

    vm.startPrank(address(s_sourceRouter));

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].token = CUSTOM_TOKEN;
    message.tokenAmounts[0].amount = 1;

    vm.expectRevert(abi.encodeWithSelector(AggregateRateLimiter.PriceNotFoundForToken.selector, CUSTOM_TOKEN));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageGasLimitTooHigh.selector));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_InvalidAddressEncodePacked_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encodePacked(address(234));

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 1, OWNER);
  }

  function test_InvalidAddress_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(type(uint208).max);

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 1, OWNER);
  }

  // We disallow sending to addresses 0-9.
  function test_ZeroAddressReceiver_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint160 i = 0; i < 10; ++i) {
      message.receiver = abi.encode(address(i));

      vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 1, OWNER);
    }
  }

  function test_MaxFeeBalanceReached_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(EVM2EVMOnRamp.MaxFeeBalanceReached.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, MAX_NOP_FEES_JUELS + 1, OWNER);
  }

  function test_InvalidChainSelector_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 wrongChainSelector = DEST_CHAIN_SELECTOR + 1;
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidChainSelector.selector, wrongChainSelector));

    s_onRamp.forwardFromRouter(wrongChainSelector, message, 1, OWNER);
  }

  function test_SourceTokenDataTooLarge_Revert() public {
    address sourceETH = s_sourceTokens[1];
    vm.stopPrank();
    vm.startPrank(OWNER);

    MaybeRevertingBurnMintTokenPool newPool = new MaybeRevertingBurnMintTokenPool(
      BurnMintERC677(sourceETH), new address[](0), address(s_mockRMN), address(s_sourceRouter)
    );
    BurnMintERC677(sourceETH).grantMintAndBurnRoles(address(newPool));
    deal(address(sourceETH), address(newPool), type(uint256).max);

    // Add TokenPool to OnRamp
    s_tokenAdminRegistry.setPool(sourceETH, address(newPool));

    // Allow chain in TokenPool
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destTokenPool),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    newPool.applyChainUpdates(chainUpdates);

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(address(sourceETH), 1000);

    // No data set, should succeed
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set max data length, should succeed
    vm.startPrank(OWNER);
    newPool.setSourceTokenData(new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES));

    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set data to max length +1, should revert
    vm.startPrank(OWNER);
    newPool.setSourceTokenData(new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 1));

    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set token config to allow larger data
    vm.startPrank(OWNER);
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: sourceETH,
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      aggregateRateLimitEnabled: false
    });
    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs, new address[](0));

    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set the token data larger than the configured token data, should revert
    vm.startPrank(OWNER);
    newPool.setSourceTokenData(new bytes(uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32 + 1));

    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_forwardFromRouter_UnsupportedToken_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1;
    message.tokenAmounts[0].token = address(1);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, message.tokenAmounts[0].token));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_EnforceOutOfOrder_Revert() public {
    // Update dynamic config to enforce allowOutOfOrderExecution = true.
    vm.stopPrank();
    vm.startPrank(OWNER);
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    s_onRamp.setDynamicConfig(
      EVM2EVMOnRamp.DynamicConfig({
        router: dynamicConfig.router,
        maxNumberOfTokensPerMsg: dynamicConfig.maxNumberOfTokensPerMsg,
        destGasOverhead: dynamicConfig.destGasOverhead,
        destGasPerPayloadByte: dynamicConfig.destGasPerPayloadByte,
        destDataAvailabilityOverheadGas: dynamicConfig.destDataAvailabilityOverheadGas,
        destGasPerDataAvailabilityByte: dynamicConfig.destGasPerDataAvailabilityByte,
        destDataAvailabilityMultiplierBps: dynamicConfig.destDataAvailabilityMultiplierBps,
        priceRegistry: dynamicConfig.priceRegistry,
        maxDataBytes: dynamicConfig.maxDataBytes,
        maxPerMsgGasLimit: dynamicConfig.maxPerMsgGasLimit,
        defaultTokenFeeUSDCents: dynamicConfig.defaultTokenFeeUSDCents,
        defaultTokenDestGasOverhead: dynamicConfig.defaultTokenDestGasOverhead,
        enforceOutOfOrder: true
      })
    );
    vm.stopPrank();

    vm.startPrank(address(s_sourceRouter));
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Empty extraArgs to should revert since it enforceOutOfOrder is true.
    message.extraArgs = "";
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectRevert(EVM2EVMOnRamp.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }
}

contract EVM2EVMOnRamp_forwardFromRouter_upgrade is EVM2EVMOnRampSetup {
  uint256 internal constant FEE_AMOUNT = 1234567890;
  EVM2EVMOnRampHelper internal s_prevOnRamp;

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    s_prevOnRamp = s_onRamp;

    s_onRamp = new EVM2EVMOnRampHelper(
      EVM2EVMOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        destChainSelector: DEST_CHAIN_SELECTOR,
        defaultTxGasLimit: GAS_LIMIT,
        maxNopFeesJuels: MAX_NOP_FEES_JUELS,
        prevOnRamp: address(s_prevOnRamp),
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      generateDynamicOnRampConfig(address(s_sourceRouter), address(s_feeQuoter)),
      _getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getNopsAndWeights()
    );
    s_onRamp.setAdmin(ADMIN);

    s_metadataHash = keccak256(
      abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, address(s_onRamp))
    );

    vm.startPrank(address(s_sourceRouter));
  }

  function test_V2_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);
  }

  function test_V2SenderNoncesReadsPreviousRamp_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 startNonce = s_onRamp.getSenderNonce(OWNER);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      assertEq(startNonce + i, s_onRamp.getSenderNonce(OWNER));
    }
  }

  function test_V2NonceStartsAtV1Nonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 startNonce = s_onRamp.getSenderNonce(OWNER);

    // send 1 message from previous onramp
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 1, s_onRamp.getSenderNonce(OWNER));

    // new onramp nonce should start from 2, while sequence number start from 1
    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, startNonce + 2, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 2, s_onRamp.getSenderNonce(OWNER));

    // after another send, nonce should be 3, and sequence number be 2
    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 2, startNonce + 3, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 3, s_onRamp.getSenderNonce(OWNER));
  }

  function test_V2NonceNewSenderStartsAtZero_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    // send 1 message from previous onramp from OWNER
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    address newSender = address(1234567);
    // new onramp nonce should start from 1 for new sender
    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, FEE_AMOUNT, newSender));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, newSender);
  }
}

contract EVM2EVMOnRamp_getFeeSetup is EVM2EVMOnRampSetup {
  uint224 internal s_feeTokenPrice;
  uint224 internal s_wrappedTokenPrice;
  uint224 internal s_customTokenPrice;

  address internal s_selfServeTokenDefaultPricing = makeAddr("self-serve-token-default-pricing");

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    // Add additional pool addresses for test tokens to mark them as supported
    s_tokenAdminRegistry.proposeAdministrator(s_sourceRouter.getWrappedNative(), OWNER);
    s_tokenAdminRegistry.acceptAdminRole(s_sourceRouter.getWrappedNative());
    s_tokenAdminRegistry.proposeAdministrator(CUSTOM_TOKEN, OWNER);
    s_tokenAdminRegistry.acceptAdminRole(CUSTOM_TOKEN);

    LockReleaseTokenPool wrappedNativePool = new LockReleaseTokenPool(
      IERC20(s_sourceRouter.getWrappedNative()), new address[](0), address(s_mockRMN), true, address(s_sourceRouter)
    );

    TokenPool.ChainUpdate[] memory wrappedNativeChainUpdate = new TokenPool.ChainUpdate[](1);
    wrappedNativeChainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(address(111111)),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    wrappedNativePool.applyChainUpdates(wrappedNativeChainUpdate);
    s_tokenAdminRegistry.setPool(s_sourceRouter.getWrappedNative(), address(wrappedNativePool));

    LockReleaseTokenPool customPool = new LockReleaseTokenPool(
      IERC20(CUSTOM_TOKEN), new address[](0), address(s_mockRMN), true, address(s_sourceRouter)
    );
    TokenPool.ChainUpdate[] memory customChainUpdate = new TokenPool.ChainUpdate[](1);
    customChainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(makeAddr("random")),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    customPool.applyChainUpdates(customChainUpdate);
    s_tokenAdminRegistry.setPool(CUSTOM_TOKEN, address(customPool));

    s_feeTokenPrice = s_sourceTokenPrices[0];
    s_wrappedTokenPrice = s_sourceTokenPrices[2];
    s_customTokenPrice = CUSTOM_TOKEN_PRICE;

    // Ensure the self-serve token is set up on the admin registry
    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, s_selfServeTokenDefaultPricing),
      abi.encode(makeAddr("self-serve-pool"))
    );
  }

  function calcUSDValueFromTokenAmount(uint224 tokenPrice, uint256 tokenAmount) internal pure returns (uint256) {
    return (tokenPrice * tokenAmount) / 1e18;
  }

  function applyBpsRatio(uint256 tokenAmount, uint16 ratio) internal pure returns (uint256) {
    return (tokenAmount * ratio) / 1e5;
  }

  function configUSDCentToWei(uint256 usdCent) internal pure returns (uint256) {
    return usdCent * 1e16;
  }
}

contract EVM2EVMOnRamp_getDataAvailabilityCost is EVM2EVMOnRamp_getFeeSetup {
  function test_EmptyMessageCalculatesDataAvailabilityCost_Success() public view {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);

    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();

    uint256 dataAvailabilityGas = dynamicConfig.destDataAvailabilityOverheadGas
      + dynamicConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * dynamicConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCost_Success() public view {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();

    uint256 dataAvailabilityLengthBytes =
      Internal.MESSAGE_FIXED_BYTES + 100 + (5 * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + 50;
    uint256 dataAvailabilityGas = dynamicConfig.destDataAvailabilityOverheadGas
      + dynamicConfig.destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * dynamicConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function test_Fuzz_ZeroDataAvailabilityGasPriceAlwaysCalculatesZeroDataAvailabilityCost_Success(
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public view {
    uint256 dataAvailabilityCostUSD =
      s_onRamp.getDataAvailabilityCost(0, messageDataLength, numberOfTokens, tokenTransferBytesOverhead);

    assertEq(0, dataAvailabilityCostUSD);
  }

  function test_Fuzz_CalculateDataAvailabilityCost_Success(
    uint32 destDataAvailabilityOverheadGas,
    uint16 destGasPerDataAvailabilityByte,
    uint16 destDataAvailabilityMultiplierBps,
    uint112 dataAvailabilityGasPrice,
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public {
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    dynamicConfig.destDataAvailabilityOverheadGas = destDataAvailabilityOverheadGas;
    dynamicConfig.destGasPerDataAvailabilityByte = destGasPerDataAvailabilityByte;
    dynamicConfig.destDataAvailabilityMultiplierBps = destDataAvailabilityMultiplierBps;
    s_onRamp.setDynamicConfig(dynamicConfig);

    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(
      dataAvailabilityGasPrice, messageDataLength, numberOfTokens, tokenTransferBytesOverhead
    );

    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES + messageDataLength
      + (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + tokenTransferBytesOverhead;

    uint256 dataAvailabilityGas =
      destDataAvailabilityOverheadGas + destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      dataAvailabilityGasPrice * dataAvailabilityGas * destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }
}

contract EVM2EVMOnRamp_getSupportedTokens is EVM2EVMOnRampSetup {
  function test_GetSupportedTokens_Revert() public {
    vm.expectRevert(EVM2EVMOnRamp.GetSupportedTokensFunctionalityRemovedCheckAdminRegistry.selector);
    s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR);
  }
}

contract EVM2EVMOnRamp_getTokenTransferCost is EVM2EVMOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_NoTokenTransferChargesZeroFee_Success() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(0, feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  function test__getTokenTransferCost_selfServeUsesDefaults_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_selfServeTokenDefaultPricing, 1000);

    // Get config to assert it isn't set
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token);

    assertFalse(transferFeeConfig.isEnabled);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // Assert that the default values are used
    assertEq(uint256(DEFAULT_TOKEN_FEE_USD_CENTS) * 1e16, feeUSDWei);
    assertEq(DEFAULT_TOKEN_DEST_GAS_OVERHEAD, destGasOverhead);
    assertEq(DEFAULT_TOKEN_BYTES_OVERHEAD, destBytesOverhead);
  }

  function test_SmallTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1000);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_ZeroAmountTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 0);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_LargeTokenTransferChargesMaxFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.maxFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_FeeTokenBpsFee_Success() public view {
    uint256 tokenAmount = 10000e18;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_feeTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[0].deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_CustomTokenBpsFee_Success() public view {
    uint256 tokenAmount = 200000e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: tokenAmount});

    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_customTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[1].deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroFeeConfigChargesMinFee_Success() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: s_sourceFeeToken,
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      aggregateRateLimitEnabled: true
    });
    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs, new address[](0));

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // if token charges 0 bps, it should cost minFee to transfer
    assertEq(configUSDCentToWei(tokenTransferFeeConfigArgs[0].minFeeUSDCents), feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_Fuzz_TokenTransferFeeDuplicateTokens_Success(uint256 transfers, uint256 amount) public view {
    // It shouldn't be possible to pay materially lower fees by splitting up the transfers.
    // Note it is possible to pay higher fees since the minimum fees are added.
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    transfers = bound(transfers, 1, dynamicConfig.maxNumberOfTokensPerMsg);
    // Cap amount to avoid overflow
    amount = bound(amount, 0, 1e36);
    Client.EVMTokenAmount[] memory multiple = new Client.EVMTokenAmount[](transfers);
    for (uint256 i = 0; i < transfers; ++i) {
      multiple[i] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: amount});
    }
    Client.EVMTokenAmount[] memory single = new Client.EVMTokenAmount[](1);
    single[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: amount * transfers});

    address feeToken = s_sourceRouter.getWrappedNative();

    (uint256 feeSingleUSDWei, uint32 gasOverheadSingle, uint32 bytesOverheadSingle) =
      s_onRamp.getTokenTransferCost(feeToken, s_wrappedTokenPrice, single);
    (uint256 feeMultipleUSDWei, uint32 gasOverheadMultiple, uint32 bytesOverheadMultiple) =
      s_onRamp.getTokenTransferCost(feeToken, s_wrappedTokenPrice, multiple);

    // Note that there can be a rounding error once per split.
    assertTrue(feeMultipleUSDWei >= (feeSingleUSDWei - dynamicConfig.maxNumberOfTokensPerMsg));
    assertEq(gasOverheadMultiple, gasOverheadSingle * transfers);
    assertEq(bytesOverheadMultiple, bytesOverheadSingle * transfers);
  }

  function test_MixedTokenTransferFee_Success() public view {
    address[3] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative(), CUSTOM_TOKEN];
    uint224[3] memory tokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice, s_customTokenPrice];
    EVM2EVMOnRamp.TokenTransferFeeConfig[3] memory tokenTransferFeeConfigs = [
      s_onRamp.getTokenTransferFeeConfig(testTokens[0]),
      s_onRamp.getTokenTransferFeeConfig(testTokens[1]),
      s_onRamp.getTokenTransferFeeConfig(testTokens[2])
    ];

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](3),
      feeToken: s_sourceRouter.getWrappedNative(),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    uint256 expectedTotalGas = 0;
    uint256 expectedTotalBytes = 0;

    // Start with small token transfers, total bps fee is lower than min token transfer fee
    for (uint256 i = 0; i < testTokens.length; ++i) {
      message.tokenAmounts[i] = Client.EVMTokenAmount({token: testTokens[i], amount: 1e14});
      EVM2EVMOnRamp.TokenTransferFeeConfig memory tokenTransferFeeConfig =
        s_onRamp.getTokenTransferFeeConfig(testTokens[i]);
      expectedTotalGas += tokenTransferFeeConfig.destGasOverhead == 0
        ? DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        : tokenTransferFeeConfig.destGasOverhead;
      expectedTotalBytes += tokenTransferFeeConfig.destBytesOverhead == 0
        ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES)
        : tokenTransferFeeConfig.destBytesOverhead;
    }
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);

    uint256 expectedFeeUSDWei = 0;
    for (uint256 i = 0; i < testTokens.length; ++i) {
      expectedFeeUSDWei += configUSDCentToWei(
        tokenTransferFeeConfigs[i].minFeeUSDCents == 0
          ? DEFAULT_TOKEN_FEE_USD_CENTS
          : tokenTransferFeeConfigs[i].minFeeUSDCents
      );
    }

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 1");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 1");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 1");

    // Set 1st token transfer to a meaningful amount so its bps fee is now between min and max fee
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: testTokens[0], amount: 10000e18});

    uint256 token0USDWei = applyBpsRatio(
      calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount), tokenTransferFeeConfigs[0].deciBps
    );
    uint256 token1USDWei = configUSDCentToWei(DEFAULT_TOKEN_FEE_USD_CENTS);

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = token0USDWei + token1USDWei + configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 2");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 2");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 2");

    // Set 2nd token transfer to a large amount that is higher than maxFeeUSD
    message.tokenAmounts[2] = Client.EVMTokenAmount({token: testTokens[2], amount: 1e36});

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_onRamp.getTokenTransferCost(message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = token0USDWei + token1USDWei + configUSDCentToWei(tokenTransferFeeConfigs[2].maxFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 3");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 3");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 3");
  }

  // reverts

  function test_UnsupportedToken_Revert() public {
    address NOT_SUPPORTED_TOKEN = address(123);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(NOT_SUPPORTED_TOKEN, 200);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, NOT_SUPPORTED_TOKEN));

    s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);
  }
}

contract EVM2EVMOnRamp_getFee is EVM2EVMOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD =
        (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) * feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_GetFeeOfZeroForTokenMessage_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
    assertTrue(feeAmount > 0);

    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeTokenConfigArgs = new EVM2EVMOnRamp.FeeTokenConfigArgs[](1);
    feeTokenConfigArgs[0] = EVM2EVMOnRamp.FeeTokenConfigArgs({
      token: message.feeToken,
      networkFeeUSDCents: 0,
      gasMultiplierWeiPerEth: 0,
      premiumMultiplierWeiPerEth: 0,
      enabled: true
    });

    s_onRamp.setFeeTokenConfig(feeTokenConfigArgs);
    EVM2EVMOnRamp.DynamicConfig memory config =
      generateDynamicOnRampConfig(address(s_sourceRouter), address(s_feeQuoter));
    config.destDataAvailabilityMultiplierBps = 0;

    s_onRamp.setDynamicConfig(config);

    feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
    assertEq(0, feeAmount);
  }

  function test_ZeroDataAvailabilityMultiplier_Success() public {
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    dynamicConfig.destDataAvailabilityMultiplierBps = 0;
    s_onRamp.setDynamicConfig(dynamicConfig);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

    uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

    uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
    uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
    uint256 messageFeeUSD =
      (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) * feeTokenConfig.premiumMultiplierWeiPerEth);

    uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD) / s_feeTokenPrice;
    assertEq(totalPriceInFeeToken, feeAmount);
  }

  function test_HighGasMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = MAX_GAS_LIMIT;
    uint256 customDataSize = MAX_DATA_SIZE;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: new bytes(customDataSize),
        tokenAmounts: new Client.EVMTokenAmount[](0),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });

      EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);
      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + customDataSize * DEST_GAS_PER_PAYLOAD_BYTE;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD =
        (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) * feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_SingleTokenMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 tokenAmount = 10000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
      message.feeToken = testTokens[i];
      EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);
      uint32 tokenGasOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token).destGasOverhead;
      uint32 destBytesOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token).destBytesOverhead;
      uint32 tokenBytesOverhead =
        destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_onRamp.getTokenTransferCost(message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_MessageWithDataAndTokenTransfer_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = 1_000_000;
    uint256 feeTokenAmount = 10000e18;
    uint256 customTokenAmount = 200000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: "",
        tokenAmounts: new Client.EVMTokenAmount[](2),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });
      EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

      message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: feeTokenAmount});
      message.tokenAmounts[1] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: customTokenAmount});
      message.data = "random bits and bytes that should be factored into the cost of the message";

      uint32 tokenGasOverhead = 0;
      uint32 tokenBytesOverhead = 0;
      for (uint256 j = 0; j < message.tokenAmounts.length; ++j) {
        tokenGasOverhead += s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[j].token).destGasOverhead;
        uint32 destBytesOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[j].token).destBytesOverhead;
        tokenBytesOverhead += destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;
      }

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed =
        customGasLimit + DEST_GAS_OVERHEAD + message.data.length * DEST_GAS_PER_PAYLOAD_BYTE + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_onRamp.getTokenTransferCost(message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  // Reverts

  function test_NotAFeeToken_Revert() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.NotAFeeToken.selector, notAFeeToken));

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(EVM2EVMOnRamp.UnsupportedNumberOfTokens.selector);
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageGasLimitTooHigh.selector));
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract EVM2EVMOnRamp_setNops is EVM2EVMOnRampSetup {
  // Used because EnumerableMap doesn't guarantee order
  mapping(address nop => uint256 weight) internal s_nopsToWeights;

  function test_SetNops_Success() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[1].nop = USER_4;
    nopsAndWeights[1].weight = 20;
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      s_nopsToWeights[nopsAndWeights[i].nop] = nopsAndWeights[i].weight;
    }

    s_onRamp.setNops(nopsAndWeights);

    (EVM2EVMOnRamp.NopAndWeight[] memory actual,) = s_onRamp.getNops();
    for (uint256 i = 0; i < actual.length; ++i) {
      assertEq(actual[i].weight, s_nopsToWeights[actual[i].nop]);
    }
  }

  function test_AdminCanSetNops_Success() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    // Should not revert
    vm.startPrank(ADMIN);
    s_onRamp.setNops(nopsAndWeights);
  }

  function test_IncludesPayment_Success() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[1].nop = USER_4;
    nopsAndWeights[1].weight = 20;
    uint32 totalWeight;
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      totalWeight += nopsAndWeights[i].weight;
      s_nopsToWeights[nopsAndWeights[i].nop] = nopsAndWeights[i].weight;
    }

    // Make sure a payout happens regardless of what the weights are set to
    uint96 nopFeesJuels = totalWeight * 5;
    // Set Nop fee juels
    deal(s_sourceFeeToken, address(s_onRamp), nopFeesJuels);
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    vm.startPrank(OWNER);

    // We don't care about the fee calculation logic in this test
    // so we don't verify the amounts. We do verify the addresses to
    // make sure the existing nops get paid and not the new ones.
    EVM2EVMOnRamp.NopAndWeight[] memory existingNopsAndWeights = getNopsAndWeights();
    for (uint256 i = 0; i < existingNopsAndWeights.length; ++i) {
      vm.expectEmit(true, false, false, false);
      emit EVM2EVMOnRamp.NopPaid(existingNopsAndWeights[i].nop, 0);
    }

    s_onRamp.setNops(nopsAndWeights);

    (EVM2EVMOnRamp.NopAndWeight[] memory actual,) = s_onRamp.getNops();
    for (uint256 i = 0; i < actual.length; ++i) {
      assertEq(actual[i].weight, s_nopsToWeights[actual[i].nop]);
    }
  }

  function test_SetNopsRemovesOldNopsCompletely_Success() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](0);
    s_onRamp.setNops(nopsAndWeights);
    (EVM2EVMOnRamp.NopAndWeight[] memory actual, uint256 totalWeight) = s_onRamp.getNops();
    assertEq(actual.length, 0);
    assertEq(totalWeight, 0);

    address prevNop = getNopsAndWeights()[0].nop;
    vm.startPrank(prevNop);

    // prev nop should not have permission to call payNops
    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdminOrNop.selector);
    s_onRamp.payNops();
  }

  // Reverts

  function test_NotEnoughFundsForPayout_Revert() public {
    uint96 nopFeesJuels = MAX_NOP_FEES_JUELS;
    // Set Nop fee juels but don't transfer LINK. This can happen when users
    // pay in non-link tokens.
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    vm.startPrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.InsufficientBalance.selector);

    s_onRamp.setNops(getNopsAndWeights());
  }

  function test_NonOwnerOrAdmin_Revert() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    vm.startPrank(STRANGER);
    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);
    s_onRamp.setNops(nopsAndWeights);
  }

  function test_LinkTokenCannotBeNop_Revert() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[0].nop = address(s_sourceTokens[0]);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidNopAddress.selector, address(s_sourceTokens[0])));

    s_onRamp.setNops(nopsAndWeights);
  }

  function test_ZeroAddressCannotBeNop_Revert() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[0].nop = address(0);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidNopAddress.selector, address(0)));

    s_onRamp.setNops(nopsAndWeights);
  }

  function test_TooManyNops_Revert() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](257);

    vm.expectRevert(EVM2EVMOnRamp.TooManyNops.selector);

    s_onRamp.setNops(nopsAndWeights);
  }
}

contract EVM2EVMOnRamp_withdrawNonLinkFees is EVM2EVMOnRampSetup {
  IERC20 internal s_token;

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();
    // Send some non-link tokens to the onRamp
    s_token = IERC20(s_sourceTokens[1]);
    deal(s_sourceTokens[1], address(s_onRamp), 100);
  }

  function test_WithdrawNonLinkFees_Success() public {
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));

    assertEq(0, s_token.balanceOf(address(s_onRamp)));
    assertEq(100, s_token.balanceOf(address(this)));
  }

  function test_SettlingBalance_Success() public {
    // Set Nop fee juels
    uint96 nopFeesJuels = 10000000;
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    vm.startPrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.LinkBalanceNotSettled.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));

    // It doesnt matter how the link tokens get to the onRamp
    // In this case we simply deal them to the ramp to show
    // anyone can settle the balance
    deal(s_sourceTokens[0], address(s_onRamp), nopFeesJuels);

    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));
  }

  function test_Fuzz_FuzzWithdrawalOnlyLeftoverLink_Success(uint96 nopFeeJuels, uint64 extraJuels) public {
    nopFeeJuels = uint96(bound(nopFeeJuels, 1, MAX_NOP_FEES_JUELS));

    // Set Nop fee juels
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeeJuels, OWNER);
    vm.startPrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.LinkBalanceNotSettled.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));

    address linkToken = s_sourceTokens[0];
    // It doesnt matter how the link tokens get to the onRamp
    // In this case we simply deal them to the ramp to show
    // anyone can settle the balance
    deal(linkToken, address(s_onRamp), nopFeeJuels + uint96(extraJuels));

    // Now that we've sent nopFeesJuels + extraJuels, we should be able to withdraw extraJuels
    address linkRecipient = address(0x123456789);
    assertEq(0, IERC20(linkToken).balanceOf(linkRecipient));

    s_onRamp.withdrawNonLinkFees(linkToken, linkRecipient);

    assertEq(extraJuels, IERC20(linkToken).balanceOf(linkRecipient));
  }

  // Reverts

  function test_LinkBalanceNotSettled_Revert() public {
    // Set Nop fee juels
    uint96 nopFeesJuels = 10000000;
    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    vm.startPrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.LinkBalanceNotSettled.selector);

    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));
  }

  function test_NonOwnerOrAdmin_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));
  }

  function test_WithdrawToZeroAddress_Revert() public {
    vm.expectRevert(EVM2EVMOnRamp.InvalidWithdrawParams.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(0));
  }
}

contract EVM2EVMOnRamp_setFeeTokenConfig is EVM2EVMOnRampSetup {
  function test_SetFeeTokenConfig_Success() public {
    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeConfig;

    vm.expectEmit();
    emit EVM2EVMOnRamp.FeeConfigSet(feeConfig);

    s_onRamp.setFeeTokenConfig(feeConfig);
  }

  function test_SetFeeTokenConfigByAdmin_Success() public {
    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeConfig;

    vm.startPrank(ADMIN);

    vm.expectEmit();
    emit EVM2EVMOnRamp.FeeConfigSet(feeConfig);

    s_onRamp.setFeeTokenConfig(feeConfig);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeConfig;
    vm.startPrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);

    s_onRamp.setFeeTokenConfig(feeConfig);
  }
}

contract EVM2EVMOnRamp_setTokenTransferFeeConfig is EVM2EVMOnRampSetup {
  function test__setTokenTransferFeeConfig_Success() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeArgs =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](2);
    tokenTransferFeeArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: address(5),
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: 312,
      aggregateRateLimitEnabled: true
    });
    tokenTransferFeeArgs[1] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: address(11),
      minFeeUSDCents: 12,
      maxFeeUSDCents: 13,
      deciBps: 14,
      destGasOverhead: 15,
      destBytesOverhead: 394,
      aggregateRateLimitEnabled: false
    });

    vm.expectEmit();
    emit EVM2EVMOnRamp.TokenTransferFeeConfigSet(tokenTransferFeeArgs);

    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeArgs, new address[](0));

    EVM2EVMOnRamp.TokenTransferFeeConfig memory config0 =
      s_onRamp.getTokenTransferFeeConfig(tokenTransferFeeArgs[0].token);

    assertEq(tokenTransferFeeArgs[0].minFeeUSDCents, config0.minFeeUSDCents);
    assertEq(tokenTransferFeeArgs[0].maxFeeUSDCents, config0.maxFeeUSDCents);
    assertEq(tokenTransferFeeArgs[0].deciBps, config0.deciBps);
    assertEq(tokenTransferFeeArgs[0].destGasOverhead, config0.destGasOverhead);
    assertEq(tokenTransferFeeArgs[0].destBytesOverhead, config0.destBytesOverhead);
    assertEq(tokenTransferFeeArgs[0].aggregateRateLimitEnabled, config0.aggregateRateLimitEnabled);
    assertTrue(config0.isEnabled);

    EVM2EVMOnRamp.TokenTransferFeeConfig memory config1 =
      s_onRamp.getTokenTransferFeeConfig(tokenTransferFeeArgs[1].token);

    assertEq(tokenTransferFeeArgs[1].minFeeUSDCents, config1.minFeeUSDCents);
    assertEq(tokenTransferFeeArgs[1].maxFeeUSDCents, config1.maxFeeUSDCents);
    assertEq(tokenTransferFeeArgs[1].deciBps, config1.deciBps);
    assertEq(tokenTransferFeeArgs[1].destGasOverhead, config1.destGasOverhead);
    assertEq(tokenTransferFeeArgs[1].destBytesOverhead, config1.destBytesOverhead);
    assertEq(tokenTransferFeeArgs[1].aggregateRateLimitEnabled, config1.aggregateRateLimitEnabled);
    assertTrue(config0.isEnabled);

    // Remove only the first token and validate only the first token is removed
    address[] memory tokensToRemove = new address[](1);
    tokensToRemove[0] = tokenTransferFeeArgs[0].token;

    vm.expectEmit();
    emit EVM2EVMOnRamp.TokenTransferFeeConfigDeleted(tokensToRemove);

    s_onRamp.setTokenTransferFeeConfig(new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](0), tokensToRemove);

    config0 = s_onRamp.getTokenTransferFeeConfig(tokenTransferFeeArgs[0].token);

    assertEq(0, config0.minFeeUSDCents);
    assertEq(0, config0.maxFeeUSDCents);
    assertEq(0, config0.deciBps);
    assertEq(0, config0.destGasOverhead);
    assertEq(0, config0.destBytesOverhead);
    assertFalse(config0.aggregateRateLimitEnabled);
    assertFalse(config0.isEnabled);

    config1 = s_onRamp.getTokenTransferFeeConfig(tokenTransferFeeArgs[1].token);

    assertEq(tokenTransferFeeArgs[1].minFeeUSDCents, config1.minFeeUSDCents);
    assertEq(tokenTransferFeeArgs[1].maxFeeUSDCents, config1.maxFeeUSDCents);
    assertEq(tokenTransferFeeArgs[1].deciBps, config1.deciBps);
    assertEq(tokenTransferFeeArgs[1].destGasOverhead, config1.destGasOverhead);
    assertEq(tokenTransferFeeArgs[1].destBytesOverhead, config1.destBytesOverhead);
    assertEq(tokenTransferFeeArgs[1].aggregateRateLimitEnabled, config1.aggregateRateLimitEnabled);
    assertTrue(config1.isEnabled);
  }

  function test__setTokenTransferFeeConfig_byAdmin_Success() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory transferFeeConfig;
    vm.startPrank(ADMIN);

    vm.expectEmit();
    emit EVM2EVMOnRamp.TokenTransferFeeConfigSet(transferFeeConfig);

    s_onRamp.setTokenTransferFeeConfig(transferFeeConfig, new address[](0));
  }

  // Reverts

  function test__setTokenTransferFeeConfig_InvalidDestBytesOverhead_Revert() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory transferFeeConfig =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    transferFeeConfig[0].destBytesOverhead = uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) - 1;
    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOnRamp.InvalidDestBytesOverhead.selector,
        transferFeeConfig[0].token,
        transferFeeConfig[0].destBytesOverhead
      )
    );
    s_onRamp.setTokenTransferFeeConfig(transferFeeConfig, new address[](0));
  }

  function test__setTokenTransferFeeConfig_OnlyCallableByOwnerOrAdmin_Revert() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory transferFeeConfig;
    vm.startPrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);

    s_onRamp.setTokenTransferFeeConfig(transferFeeConfig, new address[](0));
  }
}

contract EVM2EVMOnRamp_getTokenPool is EVM2EVMOnRampSetup {
  function test_GetTokenPool_Success() public view {
    assertEq(
      s_sourcePoolByToken[s_sourceTokens[0]],
      address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_sourceTokens[0])))
    );
    assertEq(
      s_sourcePoolByToken[s_sourceTokens[1]],
      address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_sourceTokens[1])))
    );

    address wrongToken = address(123);
    address nonExistentPool = address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(wrongToken)));

    assertEq(address(0), nonExistentPool);
  }
}

contract EVM2EVMOnRamp_setDynamicConfig is EVM2EVMOnRampSetup {
  function test_SetDynamicConfig_Success() public {
    EVM2EVMOnRamp.StaticConfig memory staticConfig = s_onRamp.getStaticConfig();
    EVM2EVMOnRamp.DynamicConfig memory newConfig = EVM2EVMOnRamp.DynamicConfig({
      router: address(2134),
      maxNumberOfTokensPerMsg: 14,
      destGasOverhead: DEST_GAS_OVERHEAD / 2,
      destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE / 2,
      destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
      destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
      destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
      priceRegistry: address(23423),
      maxDataBytes: 400,
      maxPerMsgGasLimit: MAX_GAS_LIMIT / 2,
      defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
      defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
      enforceOutOfOrder: false
    });

    vm.expectEmit();
    emit EVM2EVMOnRamp.ConfigSet(staticConfig, newConfig);

    s_onRamp.setDynamicConfig(newConfig);

    EVM2EVMOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(newConfig.router, gotDynamicConfig.router);
    assertEq(newConfig.maxNumberOfTokensPerMsg, gotDynamicConfig.maxNumberOfTokensPerMsg);
    assertEq(newConfig.destGasOverhead, gotDynamicConfig.destGasOverhead);
    assertEq(newConfig.destGasPerPayloadByte, gotDynamicConfig.destGasPerPayloadByte);
    assertEq(newConfig.priceRegistry, gotDynamicConfig.priceRegistry);
    assertEq(newConfig.maxDataBytes, gotDynamicConfig.maxDataBytes);
    assertEq(newConfig.maxPerMsgGasLimit, gotDynamicConfig.maxPerMsgGasLimit);
  }

  // Reverts

  function test_SetConfigInvalidConfig_Revert() public {
    EVM2EVMOnRamp.DynamicConfig memory newConfig = EVM2EVMOnRamp.DynamicConfig({
      router: address(1),
      maxNumberOfTokensPerMsg: 14,
      destGasOverhead: DEST_GAS_OVERHEAD / 2,
      destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE / 2,
      destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
      destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
      destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
      priceRegistry: address(23423),
      maxDataBytes: 400,
      maxPerMsgGasLimit: MAX_GAS_LIMIT / 2,
      defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
      defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
      enforceOutOfOrder: false
    });

    // Invalid price reg reverts.
    newConfig.priceRegistry = address(0);
    vm.expectRevert(EVM2EVMOnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);

    // Succeeds if valid
    newConfig.priceRegistry = address(23423);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigOnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(generateDynamicOnRampConfig(address(1), address(2)));
    vm.startPrank(ADMIN);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(generateDynamicOnRampConfig(address(1), address(2)));
  }
}
