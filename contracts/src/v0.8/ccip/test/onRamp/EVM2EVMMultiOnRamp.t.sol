// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IMessageInterceptor} from "../../interfaces/IMessageInterceptor.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {MultiAggregateRateLimiter} from "../../MultiAggregateRateLimiter.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {EVM2EVMMultiOnRamp} from "../../onRamp/EVM2EVMMultiOnRamp.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {EVM2EVMOnRampHelper} from "../helpers/EVM2EVMOnRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MessageInterceptorHelper} from "../helpers/MessageInterceptorHelper.sol";
import "./EVM2EVMMultiOnRampSetup.t.sol";

contract EVM2EVMMultiOnRamp_constructor is EVM2EVMMultiOnRampSetup {
  function test_Constructor_Success() public {
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = EVM2EVMMultiOnRamp.StaticConfig({
      linkToken: s_sourceTokens[0],
      chainSelector: SOURCE_CHAIN_SELECTOR,
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      rmnProxy: address(s_mockRMN),
      nonceManager: address(s_nonceManager),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });
    EVM2EVMMultiOnRamp.DynamicConfig memory dynamicConfig =
      _generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry));

    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.ConfigSet(staticConfig, dynamicConfig);
    // We ignore the DestChainConfig values as metadataHash is reliant on contract address.
    vm.expectEmit(true, false, false, false);
    emit EVM2EVMMultiOnRamp.DestChainAdded(
      DEST_CHAIN_SELECTOR,
      EVM2EVMMultiOnRamp.DestChainConfig({
        dynamicConfig: destChainConfigArg.dynamicConfig,
        prevOnRamp: address(0),
        sequenceNumber: 0,
        metadataHash: ""
      })
    );
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthUpdated(
      s_premiumMultiplierWeiPerEthArgs[0].token, s_premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthUpdated(
      s_premiumMultiplierWeiPerEthArgs[1].token, s_premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth
    );

    _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, address(s_sourceRouter), address(s_nonceManager), address(s_tokenAdminRegistry)
    );

    EVM2EVMMultiOnRamp.DestChainConfig memory expectedDestChainConfig = EVM2EVMMultiOnRamp.DestChainConfig({
      dynamicConfig: destChainConfigArg.dynamicConfig,
      prevOnRamp: address(0),
      sequenceNumber: 0,
      metadataHash: keccak256(
        abi.encode(
          Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, destChainConfigArg.destChainSelector, address(s_onRamp)
        )
        )
    });

    EVM2EVMMultiOnRamp.StaticConfig memory gotStaticConfig = s_onRamp.getStaticConfig();
    _assertStaticConfigsEqual(staticConfig, gotStaticConfig);

    EVM2EVMMultiOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    _assertDynamicConfigsEqual(dynamicConfig, gotDynamicConfig);

    EVM2EVMMultiOnRamp.DestChainConfig memory gotDestChainConfig = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR);
    _assertDestChainConfigsEqual(expectedDestChainConfig, gotDestChainConfig);

    uint64 gotFeeTokenConfig0 = s_onRamp.getPremiumMultiplierWeiPerEth(s_premiumMultiplierWeiPerEthArgs[0].token);
    assertEq(s_premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth, gotFeeTokenConfig0);

    uint64 gotFeeTokenConfig1 = s_onRamp.getPremiumMultiplierWeiPerEth(s_premiumMultiplierWeiPerEthArgs[1].token);
    assertEq(s_premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth, gotFeeTokenConfig1);

    // Initial values
    assertEq("EVM2EVMMultiOnRamp 1.6.0-dev", s_onRamp.typeAndVersion());
    assertEq(OWNER, s_onRamp.owner());
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber(destChainConfigArg.destChainSelector));
  }

  function test_Constructor_InvalidConfigLinkTokenEqAddressZero_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    new EVM2EVMMultiOnRampHelper(
      EVM2EVMMultiOnRamp.StaticConfig({
        linkToken: address(0),
        chainSelector: SOURCE_CHAIN_SELECTOR,
        maxFeeJuelsPerMsg: MAX_NOP_FEES_JUELS,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(s_nonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry)),
      _generateDestChainConfigArgs(),
      s_premiumMultiplierWeiPerEthArgs,
      s_tokenTransferFeeConfigArgs
    );
  }

  function test_Constructor_InvalidConfigLinkChainSelectorEqZero_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    new EVM2EVMMultiOnRampHelper(
      EVM2EVMMultiOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: 0,
        maxFeeJuelsPerMsg: MAX_NOP_FEES_JUELS,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(s_nonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry)),
      _generateDestChainConfigArgs(),
      s_premiumMultiplierWeiPerEthArgs,
      s_tokenTransferFeeConfigArgs
    );
  }

  function test_Constructor_InvalidConfigRMNProxyEqAddressZero_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp = new EVM2EVMMultiOnRampHelper(
      EVM2EVMMultiOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        maxFeeJuelsPerMsg: MAX_NOP_FEES_JUELS,
        rmnProxy: address(0),
        nonceManager: address(s_nonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry)),
      _generateDestChainConfigArgs(),
      s_premiumMultiplierWeiPerEthArgs,
      s_tokenTransferFeeConfigArgs
    );
  }

  function test_Constructor_InvalidConfigNonceManagerEqAddressZero_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    new EVM2EVMMultiOnRampHelper(
      EVM2EVMMultiOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        maxFeeJuelsPerMsg: MAX_NOP_FEES_JUELS,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(0),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry)),
      _generateDestChainConfigArgs(),
      s_premiumMultiplierWeiPerEthArgs,
      s_tokenTransferFeeConfigArgs
    );
  }

  function test_Constructor_InvalidConfigTokenAdminRegistryEqAddressZero_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    new EVM2EVMMultiOnRampHelper(
      EVM2EVMMultiOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        maxFeeJuelsPerMsg: MAX_NOP_FEES_JUELS,
        rmnProxy: address(s_mockRMN),
        nonceManager: address(s_nonceManager),
        tokenAdminRegistry: address(0)
      }),
      _generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry)),
      _generateDestChainConfigArgs(),
      s_premiumMultiplierWeiPerEthArgs,
      s_tokenTransferFeeConfigArgs
    );
  }
}

contract EVM2EVMMultiOnRamp_applyDestChainConfigUpdates is EVM2EVMMultiOnRampSetup {
  function test_Fuzz_applyDestChainConfigUpdates_Success(
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArgs
  ) public {
    vm.assume(destChainConfigArgs.destChainSelector != 0);
    vm.assume(destChainConfigArgs.dynamicConfig.defaultTxGasLimit != 0);
    destChainConfigArgs.dynamicConfig.defaultTokenDestBytesOverhead = uint32(
      bound(
        destChainConfigArgs.dynamicConfig.defaultTokenDestBytesOverhead,
        Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES,
        type(uint32).max
      )
    );
    bool isNewChain = true;

    if (destChainConfigArgs.destChainSelector == DEST_CHAIN_SELECTOR) {
      destChainConfigArgs.prevOnRamp = address(0);
      isNewChain = false;
    }
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory newDestChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](1);
    newDestChainConfigArgs[0] = destChainConfigArgs;
    EVM2EVMMultiOnRamp.DestChainConfig memory expectedDestChainConfig = EVM2EVMMultiOnRamp.DestChainConfig({
      dynamicConfig: destChainConfigArgs.dynamicConfig,
      prevOnRamp: destChainConfigArgs.prevOnRamp,
      sequenceNumber: 0,
      metadataHash: keccak256(
        abi.encode(
          Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, destChainConfigArgs.destChainSelector, address(s_onRamp)
        )
        )
    });

    if (isNewChain) {
      vm.expectEmit();
      emit EVM2EVMMultiOnRamp.DestChainAdded(destChainConfigArgs.destChainSelector, expectedDestChainConfig);
    } else {
      vm.expectEmit();
      emit EVM2EVMMultiOnRamp.DestChainDynamicConfigUpdated(
        destChainConfigArgs.destChainSelector, expectedDestChainConfig.dynamicConfig
      );
    }

    s_onRamp.applyDestChainConfigUpdates(newDestChainConfigArgs);

    _assertDestChainConfigsEqual(
      expectedDestChainConfig, s_onRamp.getDestChainConfig(destChainConfigArgs.destChainSelector)
    );
  }

  function test_applyDestChainConfigUpdates_Success() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](2);
    destChainConfigArgs[0] = _generateDestChainConfigArgs()[0];
    destChainConfigArgs[0].dynamicConfig.isEnabled = false;
    destChainConfigArgs[1] = _generateDestChainConfigArgs()[0];
    destChainConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    EVM2EVMMultiOnRamp.DestChainConfig memory expectedDestChainConfig0 = EVM2EVMMultiOnRamp.DestChainConfig({
      dynamicConfig: destChainConfigArgs[0].dynamicConfig,
      prevOnRamp: address(0),
      sequenceNumber: 0,
      metadataHash: keccak256(
        abi.encode(
          Internal.EVM_2_EVM_MESSAGE_HASH,
          SOURCE_CHAIN_SELECTOR,
          destChainConfigArgs[0].destChainSelector,
          address(s_onRamp)
        )
        )
    });

    EVM2EVMMultiOnRamp.DestChainConfig memory expectedDestChainConfig1 = EVM2EVMMultiOnRamp.DestChainConfig({
      dynamicConfig: destChainConfigArgs[1].dynamicConfig,
      prevOnRamp: address(0),
      sequenceNumber: 0,
      metadataHash: keccak256(
        abi.encode(
          Internal.EVM_2_EVM_MESSAGE_HASH,
          SOURCE_CHAIN_SELECTOR,
          destChainConfigArgs[1].destChainSelector,
          address(s_onRamp)
        )
        )
    });

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.DestChainDynamicConfigUpdated(DEST_CHAIN_SELECTOR, expectedDestChainConfig0.dynamicConfig);
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.DestChainAdded(DEST_CHAIN_SELECTOR + 1, expectedDestChainConfig1);

    vm.recordLogs();
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);

    EVM2EVMMultiOnRamp.DestChainConfig memory gotDestChainConfig0 = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR);
    EVM2EVMMultiOnRamp.DestChainConfig memory gotDestChainConfig1 = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);

    assertEq(vm.getRecordedLogs().length, 2);
    _assertDestChainConfigsEqual(expectedDestChainConfig0, gotDestChainConfig0);
    _assertDestChainConfigsEqual(expectedDestChainConfig1, gotDestChainConfig1);
  }

  function test_applyDestChainConfigUpdatesZeroIntput() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](0);

    vm.recordLogs();
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);

    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_InvalidDestChainConfigDestChainSelectorEqZero_Revert() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainSelector = 0;
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitEqZero() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.dynamicConfig.defaultTxGasLimit = 0;
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestChainConfigNewPrevOnRampOnExistingChain_Revert() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.prevOnRamp = address(1);
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestBytesOverhead_Revert() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.dynamicConfig.defaultTokenDestBytesOverhead = uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES - 1);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOnRamp.InvalidDestBytesOverhead.selector,
        address(0),
        destChainConfigArg.dynamicConfig.defaultTokenDestBytesOverhead
      )
    );

    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);
  }
}

contract EVM2EVMMultiOnRamp_forwardFromRouter is EVM2EVMMultiOnRampSetup {
  struct LegacyExtraArgs {
    uint256 gasLimit;
    bool strict;
  }

  function setUp() public virtual override {
    super.setUp();

    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));

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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouterSuccessEmptyExtraArgs() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = "";
    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    Client.EVM2AnyMessage memory expectedMessage = _generateEmptyMessage();
    expectedMessage.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}));
    vm.expectEmit();
    // We expect the message to be emitted with strict = false.
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(expectedMessage, 1, 1, feeAmount, OWNER)
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ShouldIncrementSeqNumAndNonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);

      vm.expectEmit();
      emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
      assertEq(nonceAfter, nonceBefore + 1);
    }
  }

  function test_ShouldStoreLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.FeePaid(s_sourceFeeToken, feeAmount);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);

    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), feeAmount);
  }

  function test_ShouldStoreNonLinkFees() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceTokens[1]).transferFrom(OWNER, address(s_onRamp), feeAmount);

    // Calculate conversion done by prices contract
    uint256 feeTokenPrice = s_priceRegistry.getTokenPrice(s_sourceTokens[1]).value;
    uint256 linkTokenPrice = s_priceRegistry.getTokenPrice(s_sourceFeeToken).value;
    uint256 conversionRate = (feeTokenPrice * 1e18) / linkTokenPrice;
    uint256 expectedJuels = (feeAmount * conversionRate) / 1e18;

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.FeePaid(s_sourceTokens[1], expectedJuels);
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

    Internal.EVM2EVMMessage memory expectedEvent = _messageToEvent(message, 1, 1, feeTokenAmount, originalSender);

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.FeePaid(s_sourceFeeToken, feeTokenAmount);
    vm.expectEmit(false, false, false, true);
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, expectedEvent);

    // Assert the message Id is correct
    assertEq(
      expectedEvent.messageId, s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeTokenAmount, originalSender)
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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  // Reverts

  function test_Paused_Revert() public {
    // We pause by disabling the whitelist
    vm.stopPrank();
    vm.startPrank(OWNER);
    address router = address(0);
    s_onRamp.setDynamicConfig(_generateDynamicMultiOnRampConfig(router, address(2)));
    vm.expectRevert(EVM2EVMMultiOnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_InvalidExtraArgsTag_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = bytes("bad args");

    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidExtraArgsTag.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.voteToCurse(bytes16(type(uint128).max));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.CursedByRMN.selector, DEST_CHAIN_SELECTOR));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_Permissions_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.expectRevert(EVM2EVMMultiOnRamp.MustBeCalledByRouter.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, OWNER);
  }

  function test_OriginalSender_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.RouterMustSetOriginalSender.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, _generateEmptyMessage(), 0, address(0));
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length)
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(EVM2EVMMultiOnRamp.UnsupportedNumberOfTokens.selector);
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, STRANGER);
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
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_CannotSendZeroTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 0;
    message.tokenAmounts[0].token = s_sourceTokens[0];
    vm.expectRevert(EVM2EVMMultiOnRamp.CannotSendZeroTokens.selector);
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

    Internal.PriceUpdates memory priceUpdates = getSingleTokenPriceUpdateStruct(wrongToken, 1);
    s_priceRegistry.updatePrices(priceUpdates);

    // Change back to the router
    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.UnsupportedToken.selector, wrongToken));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.MessageGasLimitTooHigh.selector));
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

  function test_MesssageFeeTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.MessageFeeTooHigh.selector, MAX_MSG_FEES_JUELS + 1, MAX_MSG_FEES_JUELS)
    );

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, MAX_MSG_FEES_JUELS + 1, OWNER);
  }

  function test_InvalidChainSelector_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 wrongChainSelector = DEST_CHAIN_SELECTOR + 1;
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.DestinationChainNotEnabled.selector, wrongChainSelector));

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
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
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
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set token config to allow larger data
    vm.startPrank(OWNER);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = sourceETH;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = EVM2EVMMultiOnRamp
      .TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      isEnabled: true
    });
    s_onRamp.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0)
    );

    vm.startPrank(address(s_sourceRouter));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Set the token data larger than the configured token data, should revert
    vm.startPrank(OWNER);
    newPool.setSourceTokenData(new bytes(uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32 + 1));

    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.SourceTokenDataTooLarge.selector, sourceETH));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_forwardFromRouter_UnsupportedToken_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 1;
    message.tokenAmounts[0].token = address(1);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.UnsupportedToken.selector, message.tokenAmounts[0].token));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }
}

contract EVM2EVMMultiOnRamp_getFeeSetup is EVM2EVMMultiOnRampSetup {
  uint224 internal s_feeTokenPrice;
  uint224 internal s_wrappedTokenPrice;
  uint224 internal s_customTokenPrice;

  address internal s_selfServeTokenDefaultPricing = makeAddr("self-serve-token-default-pricing");

  function setUp() public virtual override {
    super.setUp();

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
      remotePoolAddress: abi.encode(s_destTokenPool),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    wrappedNativePool.applyChainUpdates(wrappedNativeChainUpdate);
    s_tokenAdminRegistry.setPool(s_sourceRouter.getWrappedNative(), address(wrappedNativePool));

    LockReleaseTokenPool customPool = new LockReleaseTokenPool(
      IERC20(CUSTOM_TOKEN), new address[](0), address(s_mockRMN), true, address(s_sourceRouter)
    );
    TokenPool.ChainUpdate[] memory customChainUpdate = new TokenPool.ChainUpdate[](1);
    customChainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destTokenPool),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
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

contract EVM2EVMMultiOnRamp_getDataAvailabilityCost is EVM2EVMMultiOnRamp_getFeeSetup {
  function test_EmptyMessageCalculatesDataAvailabilityCost_Success() public {
    uint256 dataAvailabilityCostUSD =
      s_onRamp.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);

    EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
      s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;

    uint256 dataAvailabilityGas = destChainDynamicConfig.destDataAvailabilityOverheadGas
      + destChainDynamicConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    uint256 expectedDataAvailabilityCostUSD = USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas
      * destChainDynamicConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);

    // Test that the cost is destnation chain specific
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    destChainConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR + 1;
    destChainConfigArgs[0].dynamicConfig.destDataAvailabilityOverheadGas =
      destChainDynamicConfig.destDataAvailabilityOverheadGas * 2;
    destChainConfigArgs[0].dynamicConfig.destGasPerDataAvailabilityByte =
      destChainDynamicConfig.destGasPerDataAvailabilityByte * 2;
    destChainConfigArgs[0].dynamicConfig.destDataAvailabilityMultiplierBps =
      destChainDynamicConfig.destDataAvailabilityMultiplierBps * 2;
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);

    destChainDynamicConfig = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR + 1).dynamicConfig;
    uint256 dataAvailabilityCostUSD2 =
      s_onRamp.getDataAvailabilityCost(DEST_CHAIN_SELECTOR + 1, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);
    dataAvailabilityGas = destChainDynamicConfig.destDataAvailabilityOverheadGas
      + destChainDynamicConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    expectedDataAvailabilityCostUSD = USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas
      * destChainDynamicConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD2);
    assertFalse(dataAvailabilityCostUSD == dataAvailabilityCostUSD2);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCost_Success() public view {
    uint256 dataAvailabilityCostUSD =
      s_onRamp.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
      s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;

    uint256 dataAvailabilityLengthBytes =
      Internal.MESSAGE_FIXED_BYTES + 100 + (5 * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + 50;
    uint256 dataAvailabilityGas = destChainDynamicConfig.destDataAvailabilityOverheadGas
      + destChainDynamicConfig.destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD = USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas
      * destChainDynamicConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCostUnsupportedDestChainSelector_Success() public view {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(0, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    assertEq(dataAvailabilityCostUSD, 0);
  }

  function test_Fuzz_ZeroDataAvailabilityGasPriceAlwaysCalculatesZeroDataAvailabilityCost_Success(
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public view {
    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(
      DEST_CHAIN_SELECTOR, 0, messageDataLength, numberOfTokens, tokenTransferBytesOverhead
    );

    assertEq(0, dataAvailabilityCostUSD);
  }

  function test_Fuzz_CalculateDataAvailabilityCost_Success(
    uint64 destChainSelector,
    uint32 destDataAvailabilityOverheadGas,
    uint16 destGasPerDataAvailabilityByte,
    uint16 destDataAvailabilityMultiplierBps,
    uint112 dataAvailabilityGasPrice,
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public {
    vm.assume(destChainSelector != 0);
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](1);
    EVM2EVMMultiOnRamp.DestChainConfig memory destChainConfig = s_onRamp.getDestChainConfig(destChainSelector);
    destChainConfigArgs[0] = EVM2EVMMultiOnRamp.DestChainConfigArgs({
      destChainSelector: destChainSelector,
      dynamicConfig: destChainConfig.dynamicConfig,
      prevOnRamp: destChainConfig.prevOnRamp
    });
    destChainConfigArgs[0].dynamicConfig.destDataAvailabilityOverheadGas = destDataAvailabilityOverheadGas;
    destChainConfigArgs[0].dynamicConfig.destGasPerDataAvailabilityByte = destGasPerDataAvailabilityByte;
    destChainConfigArgs[0].dynamicConfig.destDataAvailabilityMultiplierBps = destDataAvailabilityMultiplierBps;
    destChainConfigArgs[0].dynamicConfig.defaultTxGasLimit = GAS_LIMIT;

    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 dataAvailabilityCostUSD = s_onRamp.getDataAvailabilityCost(
      destChainConfigArgs[0].destChainSelector,
      dataAvailabilityGasPrice,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
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

contract EVM2EVMMultiOnRamp_getSupportedTokens is EVM2EVMMultiOnRampSetup {
  function test_GetSupportedTokens_Revert() public {
    vm.expectRevert(EVM2EVMMultiOnRamp.GetSupportedTokensFunctionalityRemovedCheckAdminRegistry.selector);
    s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR);
  }
}

contract EVM2EVMMultiOnRamp_getFee is EVM2EVMMultiOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      uint64 premiumMultiplierWeiPerEth = s_onRamp.getPremiumMultiplierWeiPerEth(message.feeToken);
      EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
        s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * destChainDynamicConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD =
        (configUSDCentToWei(destChainDynamicConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_ZeroDataAvailabilityMultiplier_Success() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](1);
    EVM2EVMMultiOnRamp.DestChainConfig memory destChainConfig = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR);
    destChainConfigArgs[0] = EVM2EVMMultiOnRamp.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      dynamicConfig: destChainConfig.dynamicConfig,
      prevOnRamp: destChainConfig.prevOnRamp
    });
    destChainConfigArgs[0].dynamicConfig.destDataAvailabilityMultiplierBps = 0;
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 premiumMultiplierWeiPerEth = s_onRamp.getPremiumMultiplierWeiPerEth(message.feeToken);
    EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
      s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;

    uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

    uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
    uint256 gasFeeUSD = (gasUsed * destChainDynamicConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
    uint256 messageFeeUSD = (configUSDCentToWei(destChainDynamicConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);

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

      uint64 premiumMultiplierWeiPerEth = s_onRamp.getPremiumMultiplierWeiPerEth(message.feeToken);
      EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
        s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + customDataSize * DEST_GAS_PER_PAYLOAD_BYTE;
      uint256 gasFeeUSD = (gasUsed * destChainDynamicConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD =
        (configUSDCentToWei(destChainDynamicConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
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
      EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
        s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;
      uint32 destBytesOverhead =
        s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destBytesOverhead;
      uint32 tokenBytesOverhead =
        destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD
        + s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainDynamicConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * s_onRamp.getPremiumMultiplierWeiPerEth(message.feeToken));
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_MessageWithDataAndTokenTransfer_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = 1_000_000;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: "",
        tokenAmounts: new Client.EVMTokenAmount[](2),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });
      uint64 premiumMultiplierWeiPerEth = s_onRamp.getPremiumMultiplierWeiPerEth(message.feeToken);
      EVM2EVMMultiOnRamp.DestChainDynamicConfig memory destChainDynamicConfig =
        s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;

      message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: 10000e18}); // feeTokenAmount
      message.tokenAmounts[1] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: 200000e18}); // customTokenAmount
      message.data = "random bits and bytes that should be factored into the cost of the message";

      uint32 tokenGasOverhead = 0;
      uint32 tokenBytesOverhead = 0;
      for (uint256 j = 0; j < message.tokenAmounts.length; ++j) {
        tokenGasOverhead +=
          s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destGasOverhead;
        uint32 destBytesOverhead =
          s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destBytesOverhead;
        tokenBytesOverhead += destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;
      }

      uint256 gasUsed =
        customGasLimit + DEST_GAS_OVERHEAD + message.data.length * DEST_GAS_PER_PAYLOAD_BYTE + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainDynamicConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_onRamp.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, s_onRamp.getFee(DEST_CHAIN_SELECTOR, message));
    }
  }

  // Reverts

  function test_DestinationChainNotEnabled_Revert() public {
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.DestinationChainNotEnabled.selector, DEST_CHAIN_SELECTOR + 1)
    );
    s_onRamp.getFee(DEST_CHAIN_SELECTOR + 1, _generateEmptyMessage());
  }

  function test_NotAFeeToken_Revert() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.NotAFeeToken.selector, notAFeeToken));

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length)
    );

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(EVM2EVMMultiOnRamp.UnsupportedNumberOfTokens.selector);
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.MessageGasLimitTooHigh.selector));
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract EVM2EVMMultiOnRamp_getTokenTransferCost is EVM2EVMMultiOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_NoTokenTransferChargesZeroFee_Success() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(0, feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  function test_getTokenTransferCost_selfServeUsesDefaults_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_selfServeTokenDefaultPricing, 1000);

    // Get config to assert it isn't set
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    assertFalse(transferFeeConfig.isEnabled);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // Assert that the default values are used
    assertEq(uint256(DEFAULT_TOKEN_FEE_USD_CENTS) * 1e16, feeUSDWei);
    assertEq(DEFAULT_TOKEN_DEST_GAS_OVERHEAD, destGasOverhead);
    assertEq(DEFAULT_TOKEN_BYTES_OVERHEAD, destBytesOverhead);
  }

  function test_SmallTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1000);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroAmountTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 0);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_LargeTokenTransferChargesMaxFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.maxFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_FeeTokenBpsFee_Success() public view {
    uint256 tokenAmount = 10000e18;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_feeTokenPrice, tokenAmount);
    uint256 bpsUSDWei =
      applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_WETHTokenBpsFee_Success() public view {
    uint256 tokenAmount = 100e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceRouter.getWrappedNative(),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceRouter.getWrappedNative(), amount: tokenAmount});

    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_wrappedTokenPrice, tokenAmount);
    uint256 bpsUSDWei =
      applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig.deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
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

    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory transferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_customTokenPrice, tokenAmount);
    uint256 bpsUSDWei =
      applyBpsRatio(usdWei, s_tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[2].tokenTransferFeeConfig.deciBps);

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroFeeConfigChargesMinFee_Success() public {
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = s_sourceFeeToken;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = EVM2EVMMultiOnRamp
      .TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      isEnabled: true
    });
    s_onRamp.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0)
    );

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // if token charges 0 bps, it should cost minFee to transfer
    assertEq(
      configUSDCentToWei(tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.minFeeUSDCents),
      feeUSDWei
    );
    assertEq(0, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_Fuzz_TokenTransferFeeDuplicateTokens_Success(uint256 transfers, uint256 amount) public view {
    // It shouldn't be possible to pay materially lower fees by splitting up the transfers.
    // Note it is possible to pay higher fees since the minimum fees are added.
    EVM2EVMMultiOnRamp.DestChainDynamicConfig memory dynamicConfig =
      s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR).dynamicConfig;
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
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, single);
    (uint256 feeMultipleUSDWei, uint32 gasOverheadMultiple, uint32 bytesOverheadMultiple) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, multiple);

    // Note that there can be a rounding error once per split.
    assertGe(feeMultipleUSDWei, (feeSingleUSDWei - dynamicConfig.maxNumberOfTokensPerMsg));
    assertEq(gasOverheadMultiple, gasOverheadSingle * transfers);
    assertEq(bytesOverheadMultiple, bytesOverheadSingle * transfers);
  }

  function test_MixedTokenTransferFee_Success() public view {
    address[3] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative(), CUSTOM_TOKEN];
    uint224[3] memory tokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice, s_customTokenPrice];
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig[3] memory tokenTransferFeeConfigs = [
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[0]),
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[1]),
      s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[2])
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
      expectedTotalGas += s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[i]).destGasOverhead;
      expectedTotalBytes += s_onRamp.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[i]).destBytesOverhead;
    }
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);

    uint256 expectedFeeUSDWei = 0;
    for (uint256 i = 0; i < testTokens.length; ++i) {
      expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[i].minFeeUSDCents);
    }

    assertEq(expectedFeeUSDWei, feeUSDWei);
    assertEq(expectedTotalGas, destGasOverhead);
    assertEq(expectedTotalBytes, destBytesOverhead);

    // Set 1st token transfer to a meaningful amount so its bps fee is now between min and max fee
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: testTokens[0], amount: 10000e18});

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = applyBpsRatio(
      calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount), tokenTransferFeeConfigs[0].deciBps
    );
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[1].minFeeUSDCents);
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei);
    assertEq(expectedTotalGas, destGasOverhead);
    assertEq(expectedTotalBytes, destBytesOverhead);

    // Set 2nd token transfer to a large amount that is higher than maxFeeUSD
    message.tokenAmounts[1] = Client.EVMTokenAmount({token: testTokens[1], amount: 1e36});

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = applyBpsRatio(
      calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount), tokenTransferFeeConfigs[0].deciBps
    );
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[1].maxFeeUSDCents);
    expectedFeeUSDWei += configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei);
    assertEq(expectedTotalGas, destGasOverhead);
    assertEq(expectedTotalBytes, destBytesOverhead);
  }

  // reverts

  function test_UnsupportedToken_Revert() public {
    address NOT_SUPPORTED_TOKEN = address(123);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(NOT_SUPPORTED_TOKEN, 200);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOnRamp.UnsupportedToken.selector, NOT_SUPPORTED_TOKEN));

    s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);
  }

  function test_ValidatedPriceStaleness_Revert() public {
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

    s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);
  }
}

contract EVM2EVMMultiOnRamp_setDynamicConfig is EVM2EVMMultiOnRampSetup {
  function test_SetDynamicConfig_Success() public {
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = s_onRamp.getStaticConfig();
    EVM2EVMMultiOnRamp.DynamicConfig memory newConfig = EVM2EVMMultiOnRamp.DynamicConfig({
      router: address(2134),
      priceRegistry: address(23423),
      messageValidator: makeAddr("messageValidator"),
      feeAggregator: FEE_AGGREGATOR
    });

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.ConfigSet(staticConfig, newConfig);

    s_onRamp.setDynamicConfig(newConfig);

    EVM2EVMMultiOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(newConfig.router, gotDynamicConfig.router);
    assertEq(newConfig.priceRegistry, gotDynamicConfig.priceRegistry);
  }

  // Reverts

  function test_SetConfigInvalidConfigPriceRegistryEqAddressZero_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory newConfig = EVM2EVMMultiOnRamp.DynamicConfig({
      router: address(2134),
      priceRegistry: address(0),
      feeAggregator: FEE_AGGREGATOR,
      messageValidator: makeAddr("messageValidator")
    });

    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigInvalidConfig_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory newConfig = EVM2EVMMultiOnRamp.DynamicConfig({
      router: address(1),
      priceRegistry: address(23423),
      messageValidator: address(0),
      feeAggregator: FEE_AGGREGATOR
    });

    // Invalid price reg reverts.
    newConfig.priceRegistry = address(0);
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigInvalidConfigFeeAggregatorEqAddressZero_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory newConfig = EVM2EVMMultiOnRamp.DynamicConfig({
      router: address(2134),
      priceRegistry: address(23423),
      messageValidator: address(0),
      feeAggregator: address(0)
    });
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigOnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(_generateDynamicMultiOnRampConfig(address(1), address(2)));
    vm.startPrank(ADMIN);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(_generateDynamicMultiOnRampConfig(address(1), address(2)));
  }
}

contract EVM2EVMMultiOnRamp_withdrawFeeTokens is EVM2EVMMultiOnRampSetup {
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

    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      vm.expectEmit();
      emit EVM2EVMMultiOnRamp.FeeTokenWithdrawn(FEE_AGGREGATOR, feeTokens[i], amounts[i]);
    }

    s_onRamp.withdrawFeeTokens();

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      assertEq(IERC20(feeTokens[i]).balanceOf(FEE_AGGREGATOR), amounts[i]);
      assertEq(IERC20(feeTokens[i]).balanceOf(address(s_onRamp)), 0);
    }
  }

  function test_WithdrawFeeTokens_Success() public {
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.FeeTokenWithdrawn(FEE_AGGREGATOR, s_sourceFeeToken, s_nopFees[s_sourceFeeToken]);

    s_onRamp.withdrawFeeTokens();

    assertEq(IERC20(s_sourceFeeToken).balanceOf(FEE_AGGREGATOR), s_nopFees[s_sourceFeeToken]);
    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), 0);
  }
}

contract EVM2EVMOnRamp_applyPremiumMultiplierWeiPerEthUpdates is EVM2EVMMultiOnRampSetup {
  function test_Fuzz_applyPremiumMultiplierWeiPerEthUpdates_Success(
    EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs memory premiumMultiplierWeiPerEthArg
  ) public {
    EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = premiumMultiplierWeiPerEthArg;

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthUpdated(
      premiumMultiplierWeiPerEthArg.token, premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth
    );

    s_onRamp.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth,
      s_onRamp.getPremiumMultiplierWeiPerEth(premiumMultiplierWeiPerEthArg.token)
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesSingleToken_Success() public {
    EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = s_premiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );

    s_onRamp.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      s_premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth, s_onRamp.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesMultipleTokens_Success() public {
    EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[](2);
    premiumMultiplierWeiPerEthArgs[0] = s_premiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);
    premiumMultiplierWeiPerEthArgs[1].token = vm.addr(2);

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(2), premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth
    );

    s_onRamp.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth, s_onRamp.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
    assertEq(
      premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth, s_onRamp.getPremiumMultiplierWeiPerEth(vm.addr(2))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesZeroInput() public {
    vm.recordLogs();
    s_onRamp.applyPremiumMultiplierWeiPerEthUpdates(new EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    EVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs;
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");

    s_onRamp.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);
  }
}

contract EVM2EVMMultiOnRamp_applyTokenTransferFeeConfigUpdates is EVM2EVMMultiOnRampSetup {
  function test_Fuzz_ApplyTokenTransferFeeConfig_Success(
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig[2] memory tokenTransferFeeConfigs
  ) public {
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(2, 2);
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
        emit EVM2EVMMultiOnRamp.TokenTransferFeeConfigUpdated(
          tokenTransferFeeConfigArgs[i].destChainSelector, feeToken, tokenTransferFeeConfigs[j]
        );
      }
    }

    s_onRamp.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0)
    );

    for (uint256 i = 0; i < tokenTransferFeeConfigs.length; ++i) {
      _assertTokenTransferFeeConfigEqual(
        tokenTransferFeeConfigs[i],
        s_onRamp.getTokenTransferFeeConfig(
          tokenTransferFeeConfigArgs[0].destChainSelector,
          tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[i].token
        )
      );
    }
  }

  function test_ApplyTokenTransferFeeConfig_Success() public {
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = EVM2EVMMultiOnRamp
      .TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: 312,
      isEnabled: true
    });
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token = address(11);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig = EVM2EVMMultiOnRamp
      .TokenTransferFeeConfig({
      minFeeUSDCents: 12,
      maxFeeUSDCents: 13,
      deciBps: 14,
      destGasOverhead: 15,
      destBytesOverhead: 394,
      isEnabled: true
    });

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig
    );
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig
    );

    EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[] memory tokensToRemove =
      new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0);
    s_onRamp.applyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs, tokensToRemove);

    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory config0 = s_onRamp.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory config1 = s_onRamp.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig, config0
    );
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );

    // Remove only the first token and validate only the first token is removed
    tokensToRemove = new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](1);
    tokensToRemove[0] = EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs({
      destChainSelector: tokenTransferFeeConfigArgs[0].destChainSelector,
      token: tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    });

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.TokenTransferFeeConfigDeleted(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );

    s_onRamp.applyTokenTransferFeeConfigUpdates(new EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[](0), tokensToRemove);

    config0 = s_onRamp.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    config1 = s_onRamp.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory emptyConfig;

    _assertTokenTransferFeeConfigEqual(emptyConfig, config0);
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );
  }

  function test_ApplyTokenTransferFeeZeroInput() public {
    vm.recordLogs();
    s_onRamp.applyTokenTransferFeeConfigUpdates(
      new EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[](0),
      new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0)
    );

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    vm.startPrank(STRANGER);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs;

    vm.expectRevert("Only callable by owner");

    s_onRamp.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }

  function test_InvalidDestBytesOverhead_Revert() public {
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = EVM2EVMMultiOnRamp
      .TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES - 1),
      isEnabled: true
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOnRamp.InvalidDestBytesOverhead.selector,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.destBytesOverhead
      )
    );

    s_onRamp.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new EVM2EVMMultiOnRamp.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }
}

contract EVM2EVMMultiOnRamp_getTokenPool is EVM2EVMMultiOnRampSetup {
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
