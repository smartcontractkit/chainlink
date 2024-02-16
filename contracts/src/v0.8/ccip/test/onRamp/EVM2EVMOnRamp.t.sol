// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "./EVM2EVMOnRampSetup.t.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {AggregateRateLimiter} from "../../AggregateRateLimiter.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {MockTokenPool} from "../mocks/MockTokenPool.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";

/// @notice #constructor
contract EVM2EVMOnRamp_constructor is EVM2EVMOnRampSetup {
  event ConfigSet(EVM2EVMOnRamp.StaticConfig staticConfig, EVM2EVMOnRamp.DynamicConfig dynamicConfig);
  event PoolAdded(address token, address pool);

  function testConstructorSuccess() public {
    EVM2EVMOnRamp.StaticConfig memory staticConfig = EVM2EVMOnRamp.StaticConfig({
      linkToken: s_sourceTokens[0],
      chainSelector: SOURCE_CHAIN_SELECTOR,
      destChainSelector: DEST_CHAIN_SELECTOR,
      defaultTxGasLimit: GAS_LIMIT,
      maxNopFeesJuels: MAX_NOP_FEES_JUELS,
      prevOnRamp: address(0),
      armProxy: address(s_mockARM)
    });
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = generateDynamicOnRampConfig(
      address(s_sourceRouter),
      address(s_priceRegistry)
    );
    Internal.PoolUpdate[] memory tokensAndPools = getTokensAndPools(s_sourceTokens, getCastedSourcePools());

    vm.expectEmit();
    emit ConfigSet(staticConfig, dynamicConfig);
    vm.expectEmit();
    emit PoolAdded(tokensAndPools[0].token, tokensAndPools[0].pool);

    s_onRamp = new EVM2EVMOnRampHelper(
      staticConfig,
      dynamicConfig,
      tokensAndPools,
      getOutboundRateLimiterConfig(),
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
    assertEq(staticConfig.armProxy, gotStaticConfig.armProxy);

    EVM2EVMOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(dynamicConfig.router, gotDynamicConfig.router);
    assertEq(dynamicConfig.maxNumberOfTokensPerMsg, gotDynamicConfig.maxNumberOfTokensPerMsg);
    assertEq(dynamicConfig.destGasOverhead, gotDynamicConfig.destGasOverhead);
    assertEq(dynamicConfig.destGasPerPayloadByte, gotDynamicConfig.destGasPerPayloadByte);
    assertEq(dynamicConfig.priceRegistry, gotDynamicConfig.priceRegistry);
    assertEq(dynamicConfig.maxDataBytes, gotDynamicConfig.maxDataBytes);
    assertEq(dynamicConfig.maxPerMsgGasLimit, gotDynamicConfig.maxPerMsgGasLimit);

    // Tokens
    assertEq(s_sourceTokens, s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR));

    // Initial values
    assertEq("EVM2EVMOnRamp 1.5.0-dev", s_onRamp.typeAndVersion());
    assertEq(OWNER, s_onRamp.owner());
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber());
  }
}

contract EVM2EVMOnRamp_payNops_fuzz is EVM2EVMOnRampSetup {
  function testFuzz_NopPayNopsSuccess(uint96 nopFeesJuels) public {
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    // To avoid NoFeesToPay
    vm.assume(nopFeesJuels > weightsTotal);
    vm.assume(nopFeesJuels < MAX_NOP_FEES_JUELS);

    // Set Nop fee juels
    deal(s_sourceFeeToken, address(s_onRamp), nopFeesJuels);
    changePrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);

    changePrank(OWNER);

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
    changePrank(address(s_sourceRouter));

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
  function testOwnerPayNopsSuccess() public {
    changePrank(OWNER);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (nopsAndWeights[i].weight * totalJuels) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }

  function testAdminPayNopsSuccess() public {
    changePrank(ADMIN);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (nopsAndWeights[i].weight * totalJuels) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }

  function testNopPayNopsSuccess() public {
    changePrank(getNopsAndWeights()[0].nop);

    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    s_onRamp.payNops();
    (EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) = s_onRamp.getNops();
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      uint256 expectedPayout = (nopsAndWeights[i].weight * totalJuels) / weightsTotal;
      assertEq(IERC20(s_sourceFeeToken).balanceOf(nopsAndWeights[i].nop), expectedPayout);
    }
  }

  function testPayNopsSuccessAfterSetNops() public {
    changePrank(OWNER);

    // set 2 nops, 1 from previous, 1 new
    address prevNop = getNopsAndWeights()[0].nop;
    address newNop = STRANGER;
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](2);
    nopsAndWeights[0] = EVM2EVMOnRamp.NopAndWeight({nop: prevNop, weight: 1});
    nopsAndWeights[1] = EVM2EVMOnRamp.NopAndWeight({nop: newNop, weight: 1});
    s_onRamp.setNops(nopsAndWeights);

    // refill OnRamp nops fees
    changePrank(address(s_sourceRouter));
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), feeAmount, OWNER);

    changePrank(newNop);
    uint256 prevNopBalance = IERC20(s_sourceFeeToken).balanceOf(prevNop);
    uint256 totalJuels = s_onRamp.getNopFeesJuels();

    s_onRamp.payNops();

    assertEq(totalJuels / 2 + prevNopBalance, IERC20(s_sourceFeeToken).balanceOf(prevNop));
    assertEq(totalJuels / 2, IERC20(s_sourceFeeToken).balanceOf(newNop));
  }

  // Reverts

  function testInsufficientBalanceReverts() public {
    changePrank(address(s_onRamp));
    IERC20(s_sourceFeeToken).transfer(OWNER, IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)));
    changePrank(OWNER);
    vm.expectRevert(EVM2EVMOnRamp.InsufficientBalance.selector);
    s_onRamp.payNops();
  }

  function testWrongPermissionsReverts() public {
    changePrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdminOrNop.selector);
    s_onRamp.payNops();
  }

  function testNoFeesToPayReverts() public {
    changePrank(OWNER);
    s_onRamp.payNops();
    vm.expectRevert(EVM2EVMOnRamp.NoFeesToPay.selector);
    s_onRamp.payNops();
  }

  function testNoNopsToPayReverts() public {
    changePrank(OWNER);
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](0);
    s_onRamp.setNops(nopsAndWeights);
    vm.expectRevert(EVM2EVMOnRamp.NoNopsToPay.selector);
    s_onRamp.payNops();
  }
}

/// @notice #linkAvailableForPayment
contract EVM2EVMOnRamp_linkAvailableForPayment is EVM2EVMNopsFeeSetup {
  function testLinkAvailableForPaymentSuccess() public {
    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    uint256 linkBalance = IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp));

    assertEq(int256(linkBalance - totalJuels), s_onRamp.linkAvailableForPayment());

    changePrank(OWNER);
    s_onRamp.payNops();

    assertEq(int256(linkBalance - totalJuels), s_onRamp.linkAvailableForPayment());
  }

  function testInsufficientLinkBalanceSuccess() public {
    uint256 totalJuels = s_onRamp.getNopFeesJuels();
    uint256 linkBalance = IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp));

    changePrank(address(s_onRamp));

    uint256 linkRemaining = 1;
    IERC20(s_sourceFeeToken).transfer(OWNER, linkBalance - linkRemaining);

    changePrank(STRANGER);
    assertEq(int256(linkRemaining) - int256(totalJuels), s_onRamp.linkAvailableForPayment());
  }
}

/// @notice #forwardFromRouter
contract EVM2EVMOnRamp_forwardFromRouter is EVM2EVMOnRampSetup {
  struct LegacyExtraArgs {
    uint256 gasLimit;
    bool strict;
  }

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));

    // Since we'll mostly be testing for valid calls from the router we'll
    // mock all calls to be originating from the router and re-mock in
    // tests that require failure.
    changePrank(address(s_sourceRouter));
  }

  function testForwardFromRouterSuccessCustomExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function testForwardFromRouterSuccessLegacyExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V1_TAG,
      LegacyExtraArgs({gasLimit: GAS_LIMIT * 2, strict: true})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    // We expect the message to be emitted with strict = false.
    emit CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function testForwardFromRouterSuccess() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function testShouldIncrementSeqNumAndNonceSuccess() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_onRamp.getSenderNonce(OWNER);

      vm.expectEmit();
      emit CCIPSendRequested(_messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_onRamp.getSenderNonce(OWNER);
      assertEq(nonceAfter, nonceBefore + 1);
    }
  }

  event Transfer(address indexed from, address indexed to, uint256 value);

  function testShouldStoreLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), feeAmount);
    assertEq(s_onRamp.getNopFeesJuels(), feeAmount);
  }

  function testShouldStoreNonLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceTokens[1]).transferFrom(OWNER, address(s_onRamp), feeAmount);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceTokens[1]).balanceOf(address(s_onRamp)), feeAmount);

    // Calculate conversion done by prices contract
    uint256 feeTokenPrice = s_priceRegistry.getTokenPrice(s_sourceTokens[1]).value;
    uint256 linkTokenPrice = s_priceRegistry.getTokenPrice(s_sourceFeeToken).value;
    uint256 conversionRate = (feeTokenPrice * 1e18) / linkTokenPrice;
    uint256 expectedJuels = (feeAmount * conversionRate) / 1e18;

    assertEq(s_onRamp.getNopFeesJuels(), expectedJuels);
  }

  // Make sure any valid sender, receiver and feeAmount can be handled.
  // @TODO Temporarily setting lower fuzz run as 256 triggers snapshot gas off by 1 error.
  // https://github.com/foundry-rs/foundry/issues/5689
  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function testFuzz_ForwardFromRouterSuccess(address originalSender, address receiver, uint96 feeTokenAmount) public {
    // To avoid RouterMustSetOriginalSender
    vm.assume(originalSender != address(0));
    vm.assume(uint160(receiver) >= 10);
    vm.assume(feeTokenAmount <= MAX_NOP_FEES_JUELS);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(receiver);

    // Make sure the tokens are in the contract
    deal(s_sourceFeeToken, address(s_onRamp), feeTokenAmount);

    Internal.EVM2EVMMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeTokenAmount, originalSender);

    vm.expectEmit(false, false, false, true);
    emit CCIPSendRequested(expectedEvent);

    // Assert the message Id is correct
    assertEq(
      expectedEvent.messageId,
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeTokenAmount, originalSender)
    );
    // Assert the fee token amount is correctly assigned to the nop fee pool
    assertEq(feeTokenAmount, s_onRamp.getNopFeesJuels());
  }

  // Reverts

  function testPausedReverts() public {
    // We pause by disabling the whitelist
    changePrank(OWNER);
    address router = address(0);
    s_onRamp.setDynamicConfig(generateDynamicOnRampConfig(router, address(2)));
    vm.expectRevert(EVM2EVMOnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function testInvalidExtraArgsTagReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = bytes("bad args");

    vm.expectRevert(EVM2EVMOnRamp.InvalidExtraArgsTag.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function testUnhealthyReverts() public {
    s_mockARM.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    vm.expectRevert(EVM2EVMOnRamp.BadARMSignal.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function testPermissionsReverts() public {
    changePrank(OWNER);
    vm.expectRevert(EVM2EVMOnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function testOriginalSenderReverts() public {
    vm.expectRevert(EVM2EVMOnRamp.RouterMustSetOriginalSender.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, address(0));
  }

  function testMessageTooLargeReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function testTooManyTokensReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(EVM2EVMOnRamp.UnsupportedNumberOfTokens.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function testCannotSendZeroTokensReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 0;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    vm.expectRevert(EVM2EVMOnRamp.CannotSendZeroTokens.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function testUnsupportedTokenReverts() public {
    address wrongToken = address(1);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].token = wrongToken;
    message.tokenAmounts[0].amount = 1;

    // We need to set the price of this new token to be able to reach
    // the proper revert point. This must be called by the owner.
    changePrank(OWNER);

    Internal.PriceUpdates memory priceUpdates = getSingleTokenPriceUpdateStruct(wrongToken, 1);
    s_priceRegistry.updatePrices(priceUpdates);

    // Change back to the router
    changePrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, wrongToken));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function testMaxCapacityExceededReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 2 ** 128;
    message.tokenAmounts[0].token = s_sourceTokens[0];

    IERC20(s_sourceTokens[0]).approve(address(s_onRamp), 2 ** 128);

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector,
        getOutboundRateLimiterConfig().capacity,
        (message.tokenAmounts[0].amount * s_sourceTokenPrices[0]) / 1e18
      )
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function testPriceNotFoundForTokenReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    address fakeToken = address(1);
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].token = fakeToken;
    message.tokenAmounts[0].amount = 1;

    vm.expectRevert(abi.encodeWithSelector(AggregateRateLimiter.PriceNotFoundForToken.selector, fakeToken));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function testMessageGasLimitTooHighReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageGasLimitTooHigh.selector));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function testInvalidAddressEncodePackedReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encodePacked(address(234));

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidAddress.selector, message.receiver));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 1, OWNER);
  }

  function testInvalidAddressReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(type(uint208).max);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidAddress.selector, message.receiver));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 1, OWNER);
  }

  // We disallow sending to addresses 0-9.
  function testZeroAddressReceiverReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint160 i = 0; i < 10; ++i) {
      message.receiver = abi.encode(address(i));

      vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidAddress.selector, message.receiver));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 1, OWNER);
    }
  }

  function testMaxFeeBalanceReachedReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(EVM2EVMOnRamp.MaxFeeBalanceReached.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, MAX_NOP_FEES_JUELS + 1, OWNER);
  }

  function testInvalidChainSelectorReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 wrongChainSelector = DEST_CHAIN_SELECTOR + 1;
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidChainSelector.selector, wrongChainSelector));

    s_onRamp.forwardFromRouter(wrongChainSelector, message, 1, OWNER);
  }

  function testSourceTokenDataTooLargeReverts() public {
    address sourceETH = s_sourceTokens[1];
    changePrank(OWNER);

    MaybeRevertingBurnMintTokenPool newPool = new MaybeRevertingBurnMintTokenPool(
      BurnMintERC677(sourceETH),
      new address[](0),
      address(s_mockARM),
      address(s_sourceRouter)
    );
    // Allow Pool to burn/mint Eth
    BurnMintERC677(sourceETH).grantMintAndBurnRoles(address(newPool));
    // Pool will be burning its own balance
    deal(address(sourceETH), address(newPool), type(uint256).max);

    // Set destBytesOverhead to 0, and let tokenPool return 1 byte
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[]
      memory tokenTransferFeeConfigArgs = new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: sourceETH,
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: 0
    });
    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs);
    newPool.setSourceTokenData(new bytes(1));

    // Add TokenPool to OnRamp
    Internal.PoolUpdate[] memory removePool = new Internal.PoolUpdate[](1);
    removePool[0] = Internal.PoolUpdate({token: address(sourceETH), pool: s_sourcePools[1]});
    Internal.PoolUpdate[] memory addPool = new Internal.PoolUpdate[](1);
    addPool[0] = Internal.PoolUpdate({token: address(sourceETH), pool: address(newPool)});
    s_onRamp.applyPoolUpdates(removePool, addPool);

    // Whitelist OnRamp in TokenPool
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    newPool.applyChainUpdates(chainUpdates);

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(address(sourceETH), 1000);

    // only call OnRamp from Router
    changePrank(address(s_sourceRouter));

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }
}

/// @notice #forwardFromRouter with ramp upgrade
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
        armProxy: address(s_mockARM)
      }),
      generateDynamicOnRampConfig(address(s_sourceRouter), address(s_priceRegistry)),
      getTokensAndPools(s_sourceTokens, getCastedSourcePools()),
      getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getNopsAndWeights()
    );
    s_onRamp.setAdmin(ADMIN);

    s_metadataHash = keccak256(
      abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, address(s_onRamp))
    );

    changePrank(address(s_sourceRouter));
  }

  function testV2Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectEmit();
    emit CCIPSendRequested(_messageToEvent(message, 1, 1, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);
  }

  function testV2SenderNoncesReadsPreviousRampSuccess() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 startNonce = s_onRamp.getSenderNonce(OWNER);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      assertEq(startNonce + i, s_onRamp.getSenderNonce(OWNER));
    }
  }

  function testV2NonceStartsAtV1NonceSuccess() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 startNonce = s_onRamp.getSenderNonce(OWNER);

    // send 1 message from previous onramp
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 1, s_onRamp.getSenderNonce(OWNER));

    // new onramp nonce should start from 2, while sequence number start from 1
    vm.expectEmit();
    emit CCIPSendRequested(_messageToEvent(message, 1, startNonce + 2, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 2, s_onRamp.getSenderNonce(OWNER));

    // after another send, nonce should be 3, and sequence number be 2
    vm.expectEmit();
    emit CCIPSendRequested(_messageToEvent(message, 2, startNonce + 3, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 3, s_onRamp.getSenderNonce(OWNER));
  }

  function testV2NonceNewSenderStartsAtZeroSuccess() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    // send 1 message from previous onramp from OWNER
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    address newSender = address(1234567);
    // new onramp nonce should start from 1 for new sender
    vm.expectEmit();
    emit CCIPSendRequested(_messageToEvent(message, 1, 1, FEE_AMOUNT, newSender));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, newSender);
  }
}

contract EVM2EVMOnRamp_getFeeSetup is EVM2EVMOnRampSetup {
  uint224 internal s_feeTokenPrice;
  uint224 internal s_wrappedTokenPrice;
  uint224 internal s_customTokenPrice;

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    // Add additional pool addresses for test tokens to mark them as supported
    Internal.PoolUpdate[] memory newRamps = new Internal.PoolUpdate[](2);
    address wrappedNativePool = address(
      new LockReleaseTokenPool(
        IERC20(s_sourceRouter.getWrappedNative()),
        new address[](0),
        address(s_mockARM),
        true,
        address(s_sourceRouter)
      )
    );
    newRamps[0] = Internal.PoolUpdate({token: s_sourceRouter.getWrappedNative(), pool: wrappedNativePool});

    address customPool = address(
      new LockReleaseTokenPool(
        IERC20(CUSTOM_TOKEN),
        new address[](0),
        address(s_mockARM),
        true,
        address(s_sourceRouter)
      )
    );
    newRamps[1] = Internal.PoolUpdate({token: CUSTOM_TOKEN, pool: customPool});
    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), newRamps);

    s_feeTokenPrice = s_sourceTokenPrices[0];
    s_wrappedTokenPrice = s_sourceTokenPrices[2];
    s_customTokenPrice = CUSTOM_TOKEN_PRICE;
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

/// @notice #getDataAvailabilityCost
contract EVM2EVMOnRamp_getDataAvailabilityCost is EVM2EVMOnRamp_getFeeSetup {
  function testEmptyMessageCalculatesDataAvailabilityCostSuccess() public {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);

    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();

    uint256 dataAvailabilityGas = dynamicConfig.destDataAvailabilityOverheadGas +
      dynamicConfig.destGasPerDataAvailabilityByte *
      Internal.MESSAGE_FIXED_BYTES;
    uint256 expectedDataAvailabilityCostUSD = USD_PER_DATA_AVAILABILITY_GAS *
      dataAvailabilityGas *
      dynamicConfig.destDataAvailabilityMultiplierBps *
      1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function testSimpleMessageCalculatesDataAvailabilityCostSuccess() public {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();

    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES +
      100 +
      (5 * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) +
      50;
    uint256 dataAvailabilityGas = dynamicConfig.destDataAvailabilityOverheadGas +
      dynamicConfig.destGasPerDataAvailabilityByte *
      dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD = USD_PER_DATA_AVAILABILITY_GAS *
      dataAvailabilityGas *
      dynamicConfig.destDataAvailabilityMultiplierBps *
      1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function testFuzz_ZeroDataAvailabilityGasPriceAlwaysCalculatesZeroDataAvailabilityCostSuccess(
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(
      0,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
    );

    assertEq(0, dataAvailabilityCostUSD);
  }

  function testFuzz_CalculateDataAvailabilityCostSuccess(
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
      dataAvailabilityGasPrice,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
    );

    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES +
      messageDataLength +
      (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) +
      tokenTransferBytesOverhead;

    uint256 dataAvailabilityGas = destDataAvailabilityOverheadGas +
      destGasPerDataAvailabilityByte *
      dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD = dataAvailabilityGasPrice *
      dataAvailabilityGas *
      destDataAvailabilityMultiplierBps *
      1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }
}

/// @notice #getTokenTransferFee
contract EVM2EVMOnRamp_getTokenTransferCost is EVM2EVMOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function testNoTokenTransferChargesZeroFeeSuccess() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    assertEq(0, feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  function testSmallTokenTransferChargesMinFeeAndGasSuccess() public {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1000);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig = s_onRamp.getTokenTransferFeeConfig(
      message.tokenAmounts[0].token
    );

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function testZeroAmountTokenTransferChargesMinFeeAndAgasSuccess() public {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 0);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig = s_onRamp.getTokenTransferFeeConfig(
      message.tokenAmounts[0].token
    );

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function testLargeTokenTransferChargesMaxFeeAndGasSuccess() public {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig = s_onRamp.getTokenTransferFeeConfig(
      message.tokenAmounts[0].token
    );

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    assertEq(configUSDCentToWei(transferFeeConfig.maxFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function testFeeTokenBpsFeeSuccess() public {
    uint256 tokenAmount = 10000e18;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig = s_onRamp.getTokenTransferFeeConfig(
      message.tokenAmounts[0].token
    );

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    uint256 usdWei = calcUSDValueFromTokenAmount(s_feeTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[0].deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function testWETHTokenBpsFeeSuccess() public {
    uint256 tokenAmount = 100e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceRouter.getWrappedNative(),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceRouter.getWrappedNative(), amount: tokenAmount});

    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig = s_onRamp.getTokenTransferFeeConfig(
      message.tokenAmounts[0].token
    );

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_wrappedTokenPrice,
      message.tokenAmounts
    );

    uint256 usdWei = calcUSDValueFromTokenAmount(s_wrappedTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[1].deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function testCustomTokenBpsFeeSuccess() public {
    uint256 tokenAmount = 200000e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: tokenAmount});

    EVM2EVMOnRamp.TokenTransferFeeConfig memory transferFeeConfig = s_onRamp.getTokenTransferFeeConfig(
      message.tokenAmounts[0].token
    );

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    uint256 usdWei = calcUSDValueFromTokenAmount(s_customTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[2].deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function testZeroFeeConfigChargesMinFeeSuccess() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[]
      memory tokenTransferFeeConfigArgs = new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: s_sourceFeeToken,
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: 0
    });
    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs);

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_feeTokenPrice,
      message.tokenAmounts
    );

    // if token charges 0 bps, it should cost minFee to transfer
    assertEq(configUSDCentToWei(tokenTransferFeeConfigArgs[0].minFeeUSDCents), feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  // Temporarily setting lower fuzz run as 256 triggers snapshot gas off by 1 error.
  /// forge-config: default.fuzz.runs = 16
  /// forge-config: ccip.fuzz.runs = 16
  function testFuzz_TokenTransferFeeDuplicateTokensSuccess(uint256 transfers, uint256 amount) public {
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

    (uint256 feeSingleUSDWei, uint32 gasOverheadSingle, uint32 bytesOverheadSingle) = s_onRamp.getTokenTransferCost(
      feeToken,
      s_wrappedTokenPrice,
      single
    );
    (uint256 feeMultipleUSDWei, uint32 gasOverheadMultiple, uint32 bytesOverheadMultiple) = s_onRamp
      .getTokenTransferCost(feeToken, s_wrappedTokenPrice, multiple);

    // Note that there can be a rounding error once per split.
    assertTrue(feeMultipleUSDWei >= (feeSingleUSDWei - dynamicConfig.maxNumberOfTokensPerMsg));
    assertEq(gasOverheadMultiple, gasOverheadSingle * transfers);
    assertEq(bytesOverheadMultiple, bytesOverheadSingle * transfers);
  }

  function testMixedTokenTransferFeeSuccess() public {
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
      expectedTotalGas += s_onRamp.getTokenTransferFeeConfig(testTokens[i]).destGasOverhead;
      expectedTotalBytes += s_onRamp.getTokenTransferFeeConfig(testTokens[i]).destBytesOverhead;
    }
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_wrappedTokenPrice,
      message.tokenAmounts
    );

    uint256 expectedFeeUSDWei = 0;
    for (uint256 i = 0; i < testTokens.length; ++i) {
      expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[i].minFeeUSDCents);
    }

    assertEq(expectedFeeUSDWei, feeUSDWei);
    assertEq(expectedTotalGas, destGasOverhead);
    assertEq(expectedTotalBytes, destBytesOverhead);

    // Set 1st token transfer to a meaningful amount so its bps fee is now between min and max fee
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: testTokens[0], amount: 10000e18});

    (feeUSDWei, destGasOverhead, destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_wrappedTokenPrice,
      message.tokenAmounts
    );
    expectedFeeUSDWei = applyBpsRatio(
      calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount),
      tokenTransferFeeConfigs[0].deciBps
    );
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[1].minFeeUSDCents);
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei);
    assertEq(expectedTotalGas, destGasOverhead);
    assertEq(expectedTotalBytes, destBytesOverhead);

    // Set 2nd token transfer to a large amount that is higher than maxFeeUSD
    message.tokenAmounts[1] = Client.EVMTokenAmount({token: testTokens[1], amount: 1e36});

    (feeUSDWei, destGasOverhead, destBytesOverhead) = s_onRamp.getTokenTransferCost(
      message.feeToken,
      s_wrappedTokenPrice,
      message.tokenAmounts
    );
    expectedFeeUSDWei = applyBpsRatio(
      calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount),
      tokenTransferFeeConfigs[0].deciBps
    );
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[1].maxFeeUSDCents);
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei);
    assertEq(expectedTotalGas, destGasOverhead);
    assertEq(expectedTotalBytes, destBytesOverhead);
  }

  // reverts

  function testUnsupportedTokenReverts() public {
    address NOT_SUPPORTED_TOKEN = address(123);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(NOT_SUPPORTED_TOKEN, 200000e18);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, NOT_SUPPORTED_TOKEN));

    s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);
  }

  function testValidatedPriceStalenessReverts() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    message.tokenAmounts[0].token = s_sourceRouter.getWrappedNative();

    vm.expectRevert(
      abi.encodeWithSelector(
        PriceRegistry.StaleTokenPrice.selector,
        s_sourceRouter.getWrappedNative(),
        uint128(TWELVE_HOURS),
        uint128(TWELVE_HOURS + 1)
      )
    );

    s_onRamp.getTokenTransferCost(message.feeToken, s_feeTokenPrice, message.tokenAmounts);
  }
}

/// @notice #getFee
contract EVM2EVMOnRamp_getFee is EVM2EVMOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function testEmptyMessageSuccess() public {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) *
        feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function testZeroDataAvailabilityMultiplierSuccess() public {
    EVM2EVMOnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    dynamicConfig.destDataAvailabilityMultiplierBps = 0;
    s_onRamp.setDynamicConfig(dynamicConfig);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

    uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

    uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
    uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
    uint256 messageFeeUSD = (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) *
      feeTokenConfig.premiumMultiplierWeiPerEth);

    uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD) / s_feeTokenPrice;
    assertEq(totalPriceInFeeToken, feeAmount);
  }

  function testHighGasMessageSuccess() public {
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
      uint256 messageFeeUSD = (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) *
        feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function testSingleTokenMessageSuccess() public {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 tokenAmount = 10000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
      message.feeToken = testTokens[i];
      EVM2EVMOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);
      uint32 tokenGasOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token).destGasOverhead;
      uint32 tokenBytesOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token).destBytesOverhead;

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD, , ) = s_onRamp.getTokenTransferCost(
        message.feeToken,
        feeTokenPrices[i],
        message.tokenAmounts
      );
      uint256 messageFeeUSD = (transferFeeUSD * feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function testMessageWithDataAndTokenTransferSuccess() public {
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
        tokenBytesOverhead += s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[j].token).destBytesOverhead;
      }

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = customGasLimit +
        DEST_GAS_OVERHEAD +
        message.data.length *
        DEST_GAS_PER_PAYLOAD_BYTE +
        tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD, , ) = s_onRamp.getTokenTransferCost(
        message.feeToken,
        feeTokenPrices[i],
        message.tokenAmounts
      );
      uint256 messageFeeUSD = (transferFeeUSD * feeTokenConfig.premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  // Reverts

  function testNotAFeeTokenReverts() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.NotAFeeToken.selector, notAFeeToken));

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function testMessageTooLargeReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function testTooManyTokensReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(EVM2EVMOnRamp.UnsupportedNumberOfTokens.selector);
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function testMessageGasLimitTooHighReverts() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.MessageGasLimitTooHigh.selector));
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract EVM2EVMOnRamp_setNops is EVM2EVMOnRampSetup {
  event NopPaid(address indexed nop, uint256 amount);

  // Used because EnumerableMap doesn't guarantee order
  mapping(address nop => uint256 weight) internal s_nopsToWeights;

  function testSetNopsSuccess() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[1].nop = USER_4;
    nopsAndWeights[1].weight = 20;
    for (uint256 i = 0; i < nopsAndWeights.length; ++i) {
      s_nopsToWeights[nopsAndWeights[i].nop] = nopsAndWeights[i].weight;
    }

    s_onRamp.setNops(nopsAndWeights);

    (EVM2EVMOnRamp.NopAndWeight[] memory actual, ) = s_onRamp.getNops();
    for (uint256 i = 0; i < actual.length; ++i) {
      assertEq(actual[i].weight, s_nopsToWeights[actual[i].nop]);
    }
  }

  function testAdminCanSetNopsSuccess() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    // Should not revert
    changePrank(ADMIN);
    s_onRamp.setNops(nopsAndWeights);
  }

  function testIncludesPaymentSuccess() public {
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
    changePrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    changePrank(OWNER);

    // We don't care about the fee calculation logic in this test
    // so we don't verify the amounts. We do verify the addresses to
    // make sure the existing nops get paid and not the new ones.
    EVM2EVMOnRamp.NopAndWeight[] memory existingNopsAndWeights = getNopsAndWeights();
    for (uint256 i = 0; i < existingNopsAndWeights.length; ++i) {
      vm.expectEmit(true, false, false, false);
      emit NopPaid(existingNopsAndWeights[i].nop, 0);
    }

    s_onRamp.setNops(nopsAndWeights);

    (EVM2EVMOnRamp.NopAndWeight[] memory actual, ) = s_onRamp.getNops();
    for (uint256 i = 0; i < actual.length; ++i) {
      assertEq(actual[i].weight, s_nopsToWeights[actual[i].nop]);
    }
  }

  function testSetNopsRemovesOldNopsCompletelySuccess() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](0);
    s_onRamp.setNops(nopsAndWeights);
    (EVM2EVMOnRamp.NopAndWeight[] memory actual, uint256 totalWeight) = s_onRamp.getNops();
    assertEq(actual.length, 0);
    assertEq(totalWeight, 0);

    address prevNop = getNopsAndWeights()[0].nop;
    changePrank(prevNop);

    // prev nop should not have permission to call payNops
    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdminOrNop.selector);
    s_onRamp.payNops();
  }

  // Reverts

  function testNotEnoughFundsForPayoutReverts() public {
    uint96 nopFeesJuels = MAX_NOP_FEES_JUELS;
    // Set Nop fee juels but don't transfer LINK. This can happen when users
    // pay in non-link tokens.
    changePrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    changePrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.InsufficientBalance.selector);

    s_onRamp.setNops(getNopsAndWeights());
  }

  function testNonOwnerOrAdminReverts() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    changePrank(STRANGER);
    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);
    s_onRamp.setNops(nopsAndWeights);
  }

  function testLinkTokenCannotBeNopReverts() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[0].nop = address(s_sourceTokens[0]);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidNopAddress.selector, address(s_sourceTokens[0])));

    s_onRamp.setNops(nopsAndWeights);
  }

  function testZeroAddressCannotBeNopReverts() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = getNopsAndWeights();
    nopsAndWeights[0].nop = address(0);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.InvalidNopAddress.selector, address(0)));

    s_onRamp.setNops(nopsAndWeights);
  }

  function testTooManyNopsReverts() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](257);

    vm.expectRevert(EVM2EVMOnRamp.TooManyNops.selector);

    s_onRamp.setNops(nopsAndWeights);
  }
}

/// @notice #withdrawNonLinkFees
contract EVM2EVMOnRamp_withdrawNonLinkFees is EVM2EVMOnRampSetup {
  IERC20 internal s_token;

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();
    // Send some non-link tokens to the onRamp
    s_token = IERC20(s_sourceTokens[1]);
    deal(s_sourceTokens[1], address(s_onRamp), 100);
  }

  function testWithdrawNonLinkFeesSuccess() public {
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));

    assertEq(0, s_token.balanceOf(address(s_onRamp)));
    assertEq(100, s_token.balanceOf(address(this)));
  }

  function testSettlingBalanceSuccess() public {
    // Set Nop fee juels
    uint96 nopFeesJuels = 10000000;
    changePrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    changePrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.LinkBalanceNotSettled.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));

    // It doesnt matter how the link tokens get to the onRamp
    // In this case we simply deal them to the ramp to show
    // anyone can settle the balance
    deal(s_sourceTokens[0], address(s_onRamp), nopFeesJuels);

    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));
  }

  function testFuzz_FuzzWithdrawalOnlyLeftoverLinkSuccess(uint96 nopFeeJuels, uint64 extraJuels) public {
    nopFeeJuels = uint96(bound(nopFeeJuels, 1, MAX_NOP_FEES_JUELS));

    // Set Nop fee juels
    changePrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeeJuels, OWNER);
    changePrank(OWNER);

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

  function testLinkBalanceNotSettledReverts() public {
    // Set Nop fee juels
    uint96 nopFeesJuels = 10000000;
    changePrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), nopFeesJuels, OWNER);
    changePrank(OWNER);

    vm.expectRevert(EVM2EVMOnRamp.LinkBalanceNotSettled.selector);

    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));
  }

  function testNonOwnerOrAdminReverts() public {
    changePrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(this));
  }

  function testWithdrawToZeroAddressReverts() public {
    vm.expectRevert(EVM2EVMOnRamp.InvalidWithdrawParams.selector);
    s_onRamp.withdrawNonLinkFees(address(s_token), address(0));
  }
}

/// @notice #setFeeTokenConfig
contract EVM2EVMOnRamp_setFeeTokenConfig is EVM2EVMOnRampSetup {
  event FeeConfigSet(EVM2EVMOnRamp.FeeTokenConfigArgs[] feeConfig);

  function testSetFeeTokenConfigSuccess() public {
    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeConfig;

    vm.expectEmit();
    emit FeeConfigSet(feeConfig);

    s_onRamp.setFeeTokenConfig(feeConfig);
  }

  function testSetFeeTokenConfigByAdminSuccess() public {
    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeConfig;

    changePrank(ADMIN);

    vm.expectEmit();
    emit FeeConfigSet(feeConfig);

    s_onRamp.setFeeTokenConfig(feeConfig);
  }

  // Reverts

  function testOnlyCallableByOwnerOrAdminReverts() public {
    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeConfig;
    changePrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);

    s_onRamp.setFeeTokenConfig(feeConfig);
  }
}

/// @notice #setTokenTransferFeeConfig
contract EVM2EVMOnRamp_setTokenTransferFeeConfig is EVM2EVMOnRampSetup {
  event TokenTransferFeeConfigSet(EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] transferFeeConfig);

  function testSetTokenTransferFeeSuccess() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[]
      memory tokenTransferFeeConfigArgs = new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](2);
    tokenTransferFeeConfigArgs[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: address(0),
      minFeeUSDCents: 0,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: 0
    });
    tokenTransferFeeConfigArgs[1] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: address(1),
      minFeeUSDCents: 1,
      maxFeeUSDCents: 1,
      deciBps: 1,
      destGasOverhead: 1,
      destBytesOverhead: 1
    });

    vm.expectEmit();
    emit TokenTransferFeeConfigSet(tokenTransferFeeConfigArgs);

    s_onRamp.setTokenTransferFeeConfig(tokenTransferFeeConfigArgs);

    EVM2EVMOnRamp.TokenTransferFeeConfig memory tokenTransferFeeConfig0 = s_onRamp.getTokenTransferFeeConfig(
      address(0)
    );
    assertEq(0, tokenTransferFeeConfig0.minFeeUSDCents);
    assertEq(0, tokenTransferFeeConfig0.maxFeeUSDCents);
    assertEq(0, tokenTransferFeeConfig0.deciBps);
    assertEq(0, tokenTransferFeeConfig0.destGasOverhead);
    assertEq(0, tokenTransferFeeConfig0.destBytesOverhead);

    EVM2EVMOnRamp.TokenTransferFeeConfig memory tokenTransferFeeConfig1 = s_onRamp.getTokenTransferFeeConfig(
      address(1)
    );
    assertEq(1, tokenTransferFeeConfig1.minFeeUSDCents);
    assertEq(1, tokenTransferFeeConfig1.maxFeeUSDCents);
    assertEq(1, tokenTransferFeeConfig1.deciBps);
    assertEq(1, tokenTransferFeeConfig1.destGasOverhead);
    assertEq(1, tokenTransferFeeConfig1.destBytesOverhead);
  }

  function testSetFeeTokenConfigByAdminSuccess() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory transferFeeConfig;
    changePrank(ADMIN);

    vm.expectEmit();
    emit TokenTransferFeeConfigSet(transferFeeConfig);

    s_onRamp.setTokenTransferFeeConfig(transferFeeConfig);
  }

  // Reverts

  function testOnlyCallableByOwnerOrAdminReverts() public {
    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory transferFeeConfig;
    changePrank(STRANGER);

    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdmin.selector);

    s_onRamp.setTokenTransferFeeConfig(transferFeeConfig);
  }
}

// #getTokenPool
contract EVM2EVMOnRamp_getTokenPool is EVM2EVMOnRampSetup {
  function testGetTokenPoolSuccess() public {
    assertEq(s_sourcePools[0], address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_sourceTokens[0]))));
    assertEq(s_sourcePools[1], address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_sourceTokens[1]))));

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, IERC20(s_destTokens[0])));
    s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_destTokens[0]));
  }
}

contract EVM2EVMOnRamp_applyPoolUpdates is EVM2EVMOnRampSetup {
  event PoolAdded(address token, address pool);
  event PoolRemoved(address token, address pool);

  function testApplyPoolUpdatesSuccess() public {
    address token = address(1);
    MockTokenPool mockPool = new MockTokenPool(token);

    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({token: token, pool: address(mockPool)});

    vm.expectEmit();
    emit PoolAdded(adds[0].token, adds[0].pool);

    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    assertEq(adds[0].pool, address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(adds[0].token))));

    vm.expectEmit();
    emit PoolRemoved(adds[0].token, adds[0].pool);

    s_onRamp.applyPoolUpdates(adds, new Internal.PoolUpdate[](0));

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.UnsupportedToken.selector, adds[0].token));
    s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(adds[0].token));
  }

  function testAtomicPoolReplacementSuccess() public {
    address token = address(1);
    MockTokenPool mockPool = new MockTokenPool(token);

    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({token: token, pool: address(mockPool)});

    vm.expectEmit();
    emit PoolAdded(token, adds[0].pool);

    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    assertEq(adds[0].pool, address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(token))));

    MockTokenPool newMockPool = new MockTokenPool(token);

    Internal.PoolUpdate[] memory updates = new Internal.PoolUpdate[](1);
    updates[0] = Internal.PoolUpdate({token: token, pool: address(newMockPool)});

    vm.expectEmit();
    emit PoolRemoved(token, adds[0].pool);
    vm.expectEmit();
    emit PoolAdded(token, updates[0].pool);

    s_onRamp.applyPoolUpdates(adds, updates);

    assertEq(updates[0].pool, address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(token))));
  }

  // Reverts
  function testOnlyCallableByOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), new Internal.PoolUpdate[](0));
    changePrank(ADMIN);
    vm.expectRevert("Only callable by owner");
    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), new Internal.PoolUpdate[](0));
  }

  function testPoolAlreadyExistsReverts() public {
    address token = address(1);
    MockTokenPool mockPool = new MockTokenPool(token);

    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](2);
    adds[0] = Internal.PoolUpdate({token: token, pool: address(mockPool)});
    adds[1] = Internal.PoolUpdate({token: token, pool: address(mockPool)});

    vm.expectRevert(EVM2EVMOnRamp.PoolAlreadyAdded.selector);

    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);
  }

  function testInvalidTokenPoolConfigReverts() public {
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({token: address(0), pool: address(2)});

    vm.expectRevert(EVM2EVMOnRamp.InvalidTokenPoolConfig.selector);

    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    adds[0] = Internal.PoolUpdate({token: address(1), pool: address(0)});

    vm.expectRevert(EVM2EVMOnRamp.InvalidTokenPoolConfig.selector);

    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);
  }

  function testPoolDoesNotExistReverts() public {
    address token = address(1);
    MockTokenPool mockPool = new MockTokenPool(token);

    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](1);
    removes[0] = Internal.PoolUpdate({token: token, pool: address(mockPool)});

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOnRamp.PoolDoesNotExist.selector, removes[0].token));

    s_onRamp.applyPoolUpdates(removes, new Internal.PoolUpdate[](0));
  }

  function testRemoveTokenPoolMismatchReverts() public {
    address token = address(1);
    MockTokenPool[] memory mockPools = new MockTokenPool[](2);
    mockPools[0] = new MockTokenPool(token);
    mockPools[1] = new MockTokenPool(token);

    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({token: token, pool: address(mockPools[0])});
    s_onRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](1);
    removes[0] = Internal.PoolUpdate({token: token, pool: address(mockPools[1])});

    vm.expectRevert(EVM2EVMOnRamp.TokenPoolMismatch.selector);

    s_onRamp.applyPoolUpdates(removes, adds);
  }

  function testAddTokenPoolMismatchReverts() public {
    address token = address(1);
    MockTokenPool mockPool = new MockTokenPool(address(2));

    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](0);
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({token: token, pool: address(mockPool)});

    vm.expectRevert(EVM2EVMOnRamp.TokenPoolMismatch.selector);

    s_onRamp.applyPoolUpdates(removes, adds);
  }
}

// #getSupportedTokens
contract EVM2EVMOnRamp_getSupportedTokens is EVM2EVMOnRampSetup {
  function testGetSupportedTokensSuccess() public {
    address[] memory supportedTokens = s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR);

    assertEq(s_sourceTokens, supportedTokens);

    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](1);
    removes[0] = Internal.PoolUpdate({token: s_sourceTokens[0], pool: s_sourcePools[0]});

    s_onRamp.applyPoolUpdates(removes, new Internal.PoolUpdate[](0));

    supportedTokens = s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR);

    assertEq(address(s_sourceTokens[1]), supportedTokens[0]);
    assertEq(s_sourceTokens.length - 1, supportedTokens.length);
  }
}

// #getExpectedNextSequenceNumber
contract EVM2EVMOnRamp_getExpectedNextSequenceNumber is EVM2EVMOnRampSetup {
  function testGetExpectedNextSequenceNumberSuccess() public {
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber());
  }
}

// #setDynamicConfig
contract EVM2EVMOnRamp_setDynamicConfig is EVM2EVMOnRampSetup {
  event ConfigSet(EVM2EVMOnRamp.StaticConfig staticConfig, EVM2EVMOnRamp.DynamicConfig dynamicConfig);

  function testSetDynamicConfigSuccess() public {
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
      maxPerMsgGasLimit: MAX_GAS_LIMIT / 2
    });

    vm.expectEmit();
    emit ConfigSet(staticConfig, newConfig);

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

  function testSetConfigInvalidConfigReverts() public {
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
      maxPerMsgGasLimit: MAX_GAS_LIMIT / 2
    });

    // Invalid price reg reverts.
    newConfig.priceRegistry = address(0);
    vm.expectRevert(EVM2EVMOnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);

    // Succeeds if valid
    newConfig.priceRegistry = address(23423);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function testSetConfigOnlyOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(generateDynamicOnRampConfig(address(1), address(2)));
    changePrank(ADMIN);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(generateDynamicOnRampConfig(address(1), address(2)));
  }
}
