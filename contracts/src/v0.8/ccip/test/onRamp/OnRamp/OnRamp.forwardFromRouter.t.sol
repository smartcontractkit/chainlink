// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageInterceptor} from "../../../interfaces/IMessageInterceptor.sol";

import {IPoolV1} from "../../../interfaces/IPool.sol";
import {IRouter} from "../../../interfaces/IRouter.sol";

import {BurnMintERC20} from "../../../../shared/token/ERC20/BurnMintERC20.sol";
import {FeeQuoter} from "../../../FeeQuoter.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {Pool} from "../../../libraries/Pool.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {MaybeRevertingBurnMintTokenPool} from "../../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MessageInterceptorHelper} from "../../helpers/MessageInterceptorHelper.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OnRamp_forwardFromRouter is OnRampSetup {
  struct LegacyExtraArgs {
    uint256 gasLimit;
    bool strict;
  }

  MessageInterceptorHelper internal s_outboundMessageInterceptor;
  FeeQuoter.DestChainConfig private s_svmDestChainConfig;

  address internal s_destTokenPool = makeAddr("destTokenPool");
  address internal s_destToken = makeAddr("destToken");

  function setUp() public virtual override {
    super.setUp();
    // setup for SVM chain
    s_svmDestChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
    s_svmDestChainConfig.enforceOutOfOrder = true; // Enforcing out of order execution for messages to SVM
    s_svmDestChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    s_outboundMessageInterceptor = new MessageInterceptorHelper();

    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];
    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));

    uint64[] memory destinationChainSelectors = new uint64[](1);
    destinationChainSelectors[0] = DEST_CHAIN_SELECTOR;
    address[] memory addAllowedList = new address[](1);
    addAllowedList[0] = OWNER;
    OnRamp.AllowlistConfigArgs memory allowlistConfigArgs = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: true,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: addAllowedList,
      removedAllowlistedSenders: new address[](0)
    });
    OnRamp.AllowlistConfigArgs[] memory applyAllowlistConfigArgsItems = new OnRamp.AllowlistConfigArgs[](1);
    applyAllowlistConfigArgsItems[0] = allowlistConfigArgs;
    s_onRamp.applyAllowlistUpdates(applyAllowlistConfigArgsItems);

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
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_ConfigurableSourceRouter() public {
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
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterSuccessEmptyExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = "";
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    // We expect the message to be emitted with strict = false.
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_EVM_WithTokenTransfer() public {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceTokens[0], 1000);
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    vm.expectCall(
      s_tokenAdminRegistry.getPool(s_sourceTokens[0]),
      abi.encodeCall(
        IPoolV1.lockOrBurn,
        (
          Pool.LockOrBurnInV1({
            receiver: message.receiver,
            remoteChainSelector: DEST_CHAIN_SELECTOR,
            originalSender: OWNER,
            amount: 1000,
            localToken: s_sourceTokens[0]
          })
        )
      )
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_SVM_WithTokenTransfer() public {
    changePrank(OWNER);
    // register SVM chain family excludeSelector
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: s_svmDestChainConfig});
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigs);

    changePrank(address(s_sourceRouter));
    bytes32 svmTokenReceiver = bytes32("TOKEN RECEIVER");
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceTokens[0], 1000);
    message.extraArgs = Client._svmArgsToBytes(
      Client.SVMExtraArgsV1({
        computeUnits: GAS_LIMIT,
        accountIsWritableBitmap: 0,
        allowOutOfOrderExecution: true,
        tokenReceiver: svmTokenReceiver,
        accounts: new bytes32[](0)
      })
    );
    uint256 feeAmount = 1234567890;
    Internal.EVM2AnyRampMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeAmount, OWNER);
    expectedEvent.extraArgs = message.extraArgs;
    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, expectedEvent);

    vm.expectCall(
      s_tokenAdminRegistry.getPool(s_sourceTokens[0]),
      abi.encodeCall(
        IPoolV1.lockOrBurn,
        (
          Pool.LockOrBurnInV1({
            receiver: abi.encode(svmTokenReceiver), // For SVM, the receiver is the SVMExtraArgsV1.tokenReceiver
            remoteChainSelector: DEST_CHAIN_SELECTOR,
            originalSender: OWNER,
            amount: 1000,
            localToken: s_sourceTokens[0]
          })
        )
      )
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterExtraArgsV2() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.GENERIC_EXTRA_ARGS_V2_TAG,
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: false})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterExtraArgsV2AllowOutOfOrderTrue() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.GENERIC_EXTRA_ARGS_V2_TAG,
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: true})
    );
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ShouldIncrementSeqNumAndNonce() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberBefore = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;

      vm.expectEmit();
      emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, i, _messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberAfter = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;
      assertEq(nonceAfter, nonceBefore + 1);
      assertEq(sequenceNumberAfter, sequenceNumberBefore + 1);
    }
  }

  function test_ShouldIncrementNonceOnlyOnOrdered() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.GENERIC_EXTRA_ARGS_V2_TAG,
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: true})
    );

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      uint64 sequenceNumberBefore = s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR) - 1;

      vm.expectEmit();
      emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, i, _messageToEvent(message, i, i, 0, OWNER));

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
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), feeAmount);
  }

  // Make sure any valid sender, receiver and feeAmount can be handled.
  function testFuzz_ForwardFromRouter_Success(address originalSender, address receiver, uint96 feeTokenAmount) public {
    // To avoid RouterMustSetOriginalSender
    vm.assume(originalSender != address(0));
    vm.assume(uint160(receiver) >= Internal.PRECOMPILE_SPACE);
    feeTokenAmount = uint96(bound(feeTokenAmount, 0, MAX_MSG_FEES_JUELS));
    vm.stopPrank();

    vm.startPrank(OWNER);
    uint64[] memory destinationChainSelectors = new uint64[](1);
    destinationChainSelectors[0] = uint64(DEST_CHAIN_SELECTOR);
    address[] memory addAllowedList = new address[](1);
    addAllowedList[0] = originalSender;
    OnRamp.AllowlistConfigArgs[] memory applyAllowlistConfigArgsItems = new OnRamp.AllowlistConfigArgs[](1);
    applyAllowlistConfigArgsItems[0] = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: true,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: addAllowedList,
      removedAllowlistedSenders: new address[](0)
    });
    s_onRamp.applyAllowlistUpdates(applyAllowlistConfigArgsItems);
    vm.stopPrank();

    vm.startPrank(address(s_sourceRouter));

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(receiver);

    // Make sure the tokens are in the contract
    deal(s_sourceFeeToken, address(s_onRamp), feeTokenAmount);

    Internal.EVM2AnyRampMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeTokenAmount, originalSender);

    // Assert the message Id is correct
    assertEq(
      expectedEvent.header.messageId,
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeTokenAmount, originalSender)
    );
  }

  function test_forwardFromRouter_WithInterception() public {
    _enableOutboundMessageInterceptor();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1e18;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
    s_outboundMessageInterceptor.setMessageIdValidationState(keccak256(abi.encode(message)), false);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  // Reverts

  function test_RevertWhen_Paused() public {
    // We pause by disabling the whitelist
    vm.stopPrank();
    vm.startPrank(OWNER);
    s_onRamp.setDynamicConfig(_generateDynamicOnRampConfig(address(2)));
    vm.expectRevert(OnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_RevertWhen_InvalidExtraArgsTag() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = bytes("bad args");

    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_RevertWhen_Permissions() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.expectRevert(OnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_RevertWhen_OriginalSender() public {
    vm.expectRevert(OnRamp.RouterMustSetOriginalSender.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, address(0));
  }

  function test_RevertWhen_UnAllowedOriginalSender() public {
    vm.stopPrank();
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(OnRamp.SenderNotAllowed.selector, STRANGER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, STRANGER);
  }

  function test_RevertWhen_MessageInterceptionError() public {
    _enableOutboundMessageInterceptor();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
    uint256 feeAmount = 1234567890;
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1e18;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);
    s_outboundMessageInterceptor.setMessageIdValidationState(keccak256(abi.encode(message)), true);

    vm.expectRevert(
      abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_RevertWhen_MultiCannotSendZeroTokens() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 0;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    vm.expectRevert(OnRamp.CannotSendZeroTokens.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_RevertWhen_UnsupportedToken() public {
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

  function test_RevertWhen_forwardFromRouter_UnsupportedToken() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1;
    message.tokenAmounts[0].token = address(1);

    vm.expectRevert(abi.encodeWithSelector(OnRamp.UnsupportedToken.selector, message.tokenAmounts[0].token));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_RevertWhen_MessageFeeTooHigh() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.MessageFeeTooHigh.selector, MAX_MSG_FEES_JUELS + 1, MAX_MSG_FEES_JUELS)
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, MAX_MSG_FEES_JUELS + 1, OWNER);
  }

  function test_RevertWhen_SourceTokenDataTooLarge() public {
    address sourceETH = s_sourceTokens[1];
    vm.stopPrank();
    vm.startPrank(OWNER);

    MaybeRevertingBurnMintTokenPool newPool = new MaybeRevertingBurnMintTokenPool(
      BurnMintERC20(sourceETH),
      DEFAULT_TOKEN_DECIMALS,
      new address[](0),
      address(s_mockRMNRemote),
      address(s_sourceRouter)
    );
    BurnMintERC20(sourceETH).grantMintAndBurnRoles(address(newPool));
    deal(address(sourceETH), address(newPool), type(uint256).max);

    // Add TokenPool to OnRamp
    s_tokenAdminRegistry.setPool(sourceETH, address(newPool));

    // Allow chain in TokenPool
    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = abi.encode(s_destTokenPool);

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: abi.encode(s_destToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    newPool.applyChainUpdates(new uint64[](0), chainUpdates);

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
      minFeeUSDCents: 0,
      maxFeeUSDCents: 1,
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

  function _enableOutboundMessageInterceptor() internal {
    (, address msgSender,) = vm.readCallers();

    bool resetPrank = false;

    if (msgSender != OWNER) {
      vm.stopPrank();
      vm.startPrank(OWNER);
      resetPrank = true;
    }

    OnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    dynamicConfig.messageInterceptor = address(s_outboundMessageInterceptor);
    s_onRamp.setDynamicConfig(dynamicConfig);

    if (resetPrank) {
      vm.stopPrank();
      vm.startPrank(msgSender);
    }
  }
}
