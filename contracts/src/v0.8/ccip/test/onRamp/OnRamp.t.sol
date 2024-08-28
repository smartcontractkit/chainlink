// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IMessageInterceptor} from "../../interfaces/IMessageInterceptor.sol";
import {IRouter} from "../../interfaces/IRouter.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import "./OnRampSetup.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OnRamp_constructor is OnRampSetup {
  function test_Constructor_Success() public {
    OnRamp.StaticConfig memory staticConfig = OnRamp.StaticConfig({
      chainSelector: SOURCE_CHAIN_SELECTOR,
      rmnProxy: address(s_mockRMN),
      nonceManager: address(s_outboundNonceManager),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });
    OnRamp.DynamicConfig memory dynamicConfig = _generateDynamicOnRampConfig(address(s_feeQuoter));

    vm.expectEmit();
    emit OnRamp.ConfigSet(staticConfig, dynamicConfig);
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, OnRamp.DestChainConfig(0, s_sourceRouter));

    _deployOnRamp(SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry));

    OnRamp.StaticConfig memory gotStaticConfig = s_onRamp.getStaticConfig();
    _assertStaticConfigsEqual(staticConfig, gotStaticConfig);

    OnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    _assertDynamicConfigsEqual(dynamicConfig, gotDynamicConfig);

    // Initial values
    assertEq("OnRamp 1.6.0-dev", s_onRamp.typeAndVersion());
    assertEq(OWNER, s_onRamp.owner());
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR));
    assertEq(address(s_sourceRouter), address(s_onRamp.getRouter(DEST_CHAIN_SELECTOR)));
  }

  function test_Constructor_InvalidConfigChainSelectorEqZero_Revert() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: 0,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(s_outboundNonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }

  function test_Constructor_InvalidConfigRMNProxyEqAddressZero_Revert() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp = new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: SOURCE_CHAIN_SELECTOR,
        rmnProxy: address(0),
        nonceManager: address(s_outboundNonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }

  function test_Constructor_InvalidConfigNonceManagerEqAddressZero_Revert() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: SOURCE_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(0),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }

  function test_Constructor_InvalidConfigTokenAdminRegistryEqAddressZero_Revert() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: SOURCE_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(s_outboundNonceManager),
        tokenAdminRegistry: address(0)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }
}

contract OnRamp_forwardFromRouter is OnRampSetup {
  struct LegacyExtraArgs {
    uint256 gasLimit;
    bool strict;
  }

  function setUp() public virtual override {
    super.setUp();

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
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_Success_ConfigurableSourceRouter() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    // Change the source router for this lane
    IRouter newRouter = IRouter(makeAddr("NEW ROUTER"));
    vm.stopPrank();
    vm.prank(OWNER);
    s_onRamp.applyDestChainConfigUpdates(_generateDestChainConfigArgs(newRouter));

    // forward fails from wrong router
    vm.prank(address(s_sourceRouter));
    vm.expectRevert(OnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    // forward succeeds from correct router
    vm.prank(address(newRouter));
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
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterSuccessEmptyExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = "";
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    // We expect the message to be emitted with strict = false.
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterExtraArgsV2_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: false})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

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
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ShouldIncrementSeqNumAndNonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberBefore = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;

      vm.expectEmit();
      emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberAfter = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;
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
      uint64 nonceBefore = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberBefore = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;

      vm.expectEmit();
      emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberAfter = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;
      assertEq(nonceAfter, nonceBefore);
      assertEq(sequenceNumberAfter, sequenceNumberBefore + 1);
    }
  }

  function test_ShouldStoreLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.FeePaid(s_sourceFeeToken, feeAmount);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), feeAmount);
  }

  function test_ShouldStoreNonLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceTokens[1]).transferFrom(OWNER, address(s_onRamp), feeAmount);

    // Calculate conversion done by prices contract
    uint256 feeTokenPrice = s_feeQuoter.getTokenPrice(s_sourceTokens[1]).value;
    uint256 linkTokenPrice = s_feeQuoter.getTokenPrice(s_sourceFeeToken).value;
    uint256 conversionRate = (feeTokenPrice * 1e18) / linkTokenPrice;
    uint256 expectedJuels = (feeAmount * conversionRate) / 1e18;

    vm.expectEmit();
    emit OnRamp.FeePaid(s_sourceTokens[1], expectedJuels);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceTokens[1]).balanceOf(address(s_onRamp)), feeAmount);
  }

  // Make sure any valid sender, receiver and feeAmount can be handled.
  // @TODO Temporarily setting lower fuzz run as 256 triggers snapshot gas off by 1 error.
  // https://github.com/foundry-rs/foundry/issues/5689
  /// forge-dynamicConfig: default.fuzz.runs = 32
  /// forge-dynamicConfig: ccip.fuzz.runs = 32
  function test_Fuzz_ForwardFromRouter_Success(address originalSender, address receiver, uint96 feeTokenAmount) public {
    // To avoid RouterMustSetOriginalSender
    vm.assume(originalSender != address(0));
    vm.assume(uint160(receiver) >= Internal.PRECOMPILE_SPACE);
    feeTokenAmount = uint96(bound(feeTokenAmount, 0, MAX_MSG_FEES_JUELS));

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(receiver);

    // Make sure the tokens are in the contract
    deal(s_sourceFeeToken, address(s_onRamp), feeTokenAmount);

    Internal.EVM2AnyRampMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeTokenAmount, originalSender);

    vm.expectEmit();
    emit OnRamp.FeePaid(s_sourceFeeToken, feeTokenAmount);
    vm.expectEmit(false, false, false, true);
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, expectedEvent);

    // Assert the message Id is correct
    assertEq(
      expectedEvent.header.messageId,
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeTokenAmount, originalSender)
    );
  }

  function test_forwardFromRouter_WithValidation_Success() public {
    _enableOutboundMessageValidator();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1e18;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
    s_outboundMessageValidator.setMessageIdValidationState(keccak256(abi.encode(message)), false);

    vm.expectEmit();
    emit OnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  // Reverts

  function test_Paused_Revert() public {
    // We pause by disabling the whitelist
    vm.stopPrank();
    vm.startPrank(OWNER);
    s_onRamp.setDynamicConfig(_generateDynamicOnRampConfig(address(2)));
    vm.expectRevert(OnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_InvalidExtraArgsTag_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = bytes("bad args");

    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_Permissions_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.expectRevert(OnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_OriginalSender_Revert() public {
    vm.expectRevert(OnRamp.RouterMustSetOriginalSender.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, address(0));
  }

  function test_MessageValidationError_Revert() public {
    _enableOutboundMessageValidator();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1e18;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
    s_outboundMessageValidator.setMessageIdValidationState(keccak256(abi.encode(message)), true);

    vm.expectRevert(
      abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_CannotSendZeroTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 0;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    vm.expectRevert(OnRamp.CannotSendZeroTokens.selector);
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
    vm.expectRevert(abi.encodeWithSelector(OnRamp.UnsupportedToken.selector, wrongToken));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_forwardFromRouter_UnsupportedToken_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1;
    message.tokenAmounts[0].token = address(1);

    vm.expectRevert(abi.encodeWithSelector(OnRamp.UnsupportedToken.selector, message.tokenAmounts[0].token));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_MesssageFeeTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.MessageFeeTooHigh.selector, MAX_MSG_FEES_JUELS + 1, MAX_MSG_FEES_JUELS)
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, MAX_MSG_FEES_JUELS + 1, OWNER);
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
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set token config to allow larger data
    vm.startPrank(OWNER);
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = sourceETH;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      isEnabled: true
    });
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set the token data larger than the configured token data, should revert
    vm.startPrank(OWNER);
    newPool.setSourceTokenData(new bytes(uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32 + 1));

    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }
}

contract OnRamp_getSupportedTokens is OnRampSetup {
  function test_GetSupportedTokens_Revert() public {
    vm.expectRevert(OnRamp.GetSupportedTokensFunctionalityRemovedCheckAdminRegistry.selector);
    s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR);
  }
}

contract OnRamp_getFee is OnRampSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
      uint256 expectedFeeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      assertEq(expectedFeeAmount, feeAmount);
    }
  }

  function test_SingleTokenMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 tokenAmount = 10000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
      message.feeToken = testTokens[i];

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
      uint256 expectedFeeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      assertEq(expectedFeeAmount, feeAmount);
    }
  }

  // Reverts

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(abi.encodeWithSelector(OnRamp.CursedByRMN.selector, DEST_CHAIN_SELECTOR));
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, _generateEmptyMessage());
  }

  function test_EnforceOutOfOrder_Revert() public {
    // Update dynamic config to enforce allowOutOfOrderExecution = true.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = true;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Empty extraArgs to should revert since it enforceOutOfOrder is true.
    message.extraArgs = "";

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract OnRamp_setDynamicConfig is OnRampSetup {
  function test_SetDynamicConfig_Success() public {
    OnRamp.StaticConfig memory staticConfig = s_onRamp.getStaticConfig();
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(23423),
      messageValidator: makeAddr("messageValidator"),
      feeAggregator: FEE_AGGREGATOR
    });

    vm.expectEmit();
    emit OnRamp.ConfigSet(staticConfig, newConfig);

    s_onRamp.setDynamicConfig(newConfig);

    OnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(newConfig.feeQuoter, gotDynamicConfig.feeQuoter);
  }

  // Reverts

  function test_SetConfigInvalidConfigFeeQuoterEqAddressZero_Revert() public {
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(0),
      feeAggregator: FEE_AGGREGATOR,
      messageValidator: makeAddr("messageValidator")
    });

    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigInvalidConfig_Revert() public {
    OnRamp.DynamicConfig memory newConfig =
      OnRamp.DynamicConfig({feeQuoter: address(23423), messageValidator: address(0), feeAggregator: FEE_AGGREGATOR});

    // Invalid price reg reverts.
    newConfig.feeQuoter = address(0);
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigInvalidConfigFeeAggregatorEqAddressZero_Revert() public {
    OnRamp.DynamicConfig memory newConfig =
      OnRamp.DynamicConfig({feeQuoter: address(23423), messageValidator: address(0), feeAggregator: address(0)});
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigOnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(_generateDynamicOnRampConfig(address(2)));
    vm.startPrank(ADMIN);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(_generateDynamicOnRampConfig(address(2)));
  }
}

contract OnRamp_withdrawFeeTokens is OnRampSetup {
  mapping(address => uint256) internal s_nopFees;

  function setUp() public virtual override {
    super.setUp();

    // Since we'll mostly be testing for valid calls from the router we'll
    // mock all calls to be originating from the router and re-mock in
    // tests that require failure.
    vm.startPrank(address(s_sourceRouter));

    uint256 feeAmount = 1234567890;

    // Send a bunch of messages, increasing the juels in the contract
    for (uint256 i = 0; i < s_sourceFeeTokens.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = s_sourceFeeTokens[i % s_sourceFeeTokens.length];
      uint256 newFeeTokenBalance = IERC20(message.feeToken).balanceOf(address(s_onRamp)) + feeAmount;
      deal(message.feeToken, address(s_onRamp), newFeeTokenBalance);
      s_nopFees[message.feeToken] = newFeeTokenBalance;
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
    }
  }

  function test_Fuzz_WithdrawFeeTokens_Success(uint256[5] memory amounts) public {
    vm.startPrank(OWNER);
    address[] memory feeTokens = new address[](amounts.length);
    for (uint256 i = 0; i < amounts.length; ++i) {
      vm.assume(amounts[i] > 0);
      feeTokens[i] = _deploySourceToken("", amounts[i], 18);
      IERC20(feeTokens[i]).transfer(address(s_onRamp), amounts[i]);
    }

    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      vm.expectEmit();
      emit OnRamp.FeeTokenWithdrawn(FEE_AGGREGATOR, feeTokens[i], amounts[i]);
    }

    s_onRamp.withdrawFeeTokens();

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      assertEq(IERC20(feeTokens[i]).balanceOf(FEE_AGGREGATOR), amounts[i]);
      assertEq(IERC20(feeTokens[i]).balanceOf(address(s_onRamp)), 0);
    }
  }

  function test_WithdrawFeeTokens_Success() public {
    vm.expectEmit();
    emit OnRamp.FeeTokenWithdrawn(FEE_AGGREGATOR, s_sourceFeeToken, s_nopFees[s_sourceFeeToken]);

    s_onRamp.withdrawFeeTokens();

    assertEq(IERC20(s_sourceFeeToken).balanceOf(FEE_AGGREGATOR), s_nopFees[s_sourceFeeToken]);
    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), 0);
  }
}

contract OnRamp_getTokenPool is OnRampSetup {
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

contract OnRamp_applyDestChainConfigUpdates is OnRampSetup {
  function test_ApplyDestChainConfigUpdates_Success() external {
    vm.stopPrank();
    vm.startPrank(OWNER);
    OnRamp.DestChainConfigArgs[] memory configArgs = new OnRamp.DestChainConfigArgs[](1);
    configArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;

    // supports disabling a lane by setting a router to zero
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, OnRamp.DestChainConfig(0, IRouter(address(0))));
    s_onRamp.applyDestChainConfigUpdates(configArgs);
    assertEq(address(0), address(s_onRamp.getRouter(DEST_CHAIN_SELECTOR)));

    // supports updating and adding lanes simultaneously
    configArgs = new OnRamp.DestChainConfigArgs[](2);
    configArgs[0] = OnRamp.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, router: s_sourceRouter});
    configArgs[1] = OnRamp.DestChainConfigArgs({destChainSelector: 9999, router: IRouter(address(9999))});
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, OnRamp.DestChainConfig(0, s_sourceRouter));
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(9999, OnRamp.DestChainConfig(0, IRouter(address(9999))));
    s_onRamp.applyDestChainConfigUpdates(configArgs);
    assertEq(address(s_sourceRouter), address(s_onRamp.getRouter(DEST_CHAIN_SELECTOR)));
    assertEq(address(9999), address(s_onRamp.getRouter(9999)));

    // handles empty list
    uint256 numLogs = vm.getRecordedLogs().length;
    configArgs = new OnRamp.DestChainConfigArgs[](0);
    s_onRamp.applyDestChainConfigUpdates(configArgs);
    assertEq(numLogs, vm.getRecordedLogs().length); // indicates no changes made
  }

  function test_ApplyDestChainConfigUpdates_WithInalidChainSelector_Revert() external {
    vm.stopPrank();
    vm.startPrank(OWNER);
    OnRamp.DestChainConfigArgs[] memory configArgs = new OnRamp.DestChainConfigArgs[](1);
    configArgs[0].destChainSelector = 0; // invalid
    vm.expectRevert(abi.encodeWithSelector(OnRamp.InvalidDestChainConfig.selector, 0));
    s_onRamp.applyDestChainConfigUpdates(configArgs);
  }
}
