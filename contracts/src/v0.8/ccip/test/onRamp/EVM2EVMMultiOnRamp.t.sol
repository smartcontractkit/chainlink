// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {AggregateRateLimiter} from "../../AggregateRateLimiter.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {EVM2EVMMultiOnRamp} from "../../onRamp/EVM2EVMMultiOnRamp.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import "./EVM2EVMMultiOnRampSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract EVM2EVMMultiOnRamp_constructor is EVM2EVMMultiOnRampSetup {
  function test_Constructor_Success() public {
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = EVM2EVMMultiOnRamp.StaticConfig({
      linkToken: s_sourceTokens[0],
      chainSelector: SOURCE_CHAIN_SELECTOR,
      maxNopFeesJuels: MAX_NOP_FEES_JUELS,
      rmnProxy: address(s_mockRMN)
    });
    EVM2EVMMultiOnRamp.DynamicConfig memory dynamicConfig =
      generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry), address(s_tokenAdminRegistry));

    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
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

    s_onRamp = new EVM2EVMMultiOnRampHelper(
      staticConfig,
      dynamicConfig,
      destChainConfigArgs,
      getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getMultiOnRampNopsAndWeights()
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
    assertStaticConfigsEqual(staticConfig, gotStaticConfig);

    EVM2EVMMultiOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertDynamicConfigsEqual(dynamicConfig, gotDynamicConfig);

    EVM2EVMMultiOnRamp.DestChainConfig memory gotDestChainConfig = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR);
    assertDestChainConfigsEqual(expectedDestChainConfig, gotDestChainConfig);

    // Initial values
    assertEq("EVM2EVMMultiOnRamp 1.6.0-dev", s_onRamp.typeAndVersion());
    assertEq(OWNER, s_onRamp.owner());
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber(destChainConfigArg.destChainSelector));
  }

  function test_Constructor_InvalidConfigLinkTokenEqAddressZero_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory dynamicConfig =
      generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry), address(s_tokenAdminRegistry));
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = EVM2EVMMultiOnRamp.StaticConfig({
      linkToken: address(0),
      chainSelector: SOURCE_CHAIN_SELECTOR,
      maxNopFeesJuels: MAX_NOP_FEES_JUELS,
      rmnProxy: address(s_mockRMN)
    });

    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp = new EVM2EVMMultiOnRampHelper(
      staticConfig,
      dynamicConfig,
      destChainConfigArgs,
      getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getMultiOnRampNopsAndWeights()
    );
  }

  function test_Constructor_InvalidConfigLinkChainSelectorEqZero_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory dynamicConfig =
      generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry), address(s_tokenAdminRegistry));
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = EVM2EVMMultiOnRamp.StaticConfig({
      linkToken: s_sourceTokens[0],
      chainSelector: 0,
      maxNopFeesJuels: MAX_NOP_FEES_JUELS,
      rmnProxy: address(s_mockRMN)
    });

    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp = new EVM2EVMMultiOnRampHelper(
      staticConfig,
      dynamicConfig,
      destChainConfigArgs,
      getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getMultiOnRampNopsAndWeights()
    );
  }

  function test_Constructor_InvalidConfigRMNProxyEqAddressZero_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory dynamicConfig =
      generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry), address(s_tokenAdminRegistry));
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = EVM2EVMMultiOnRamp.StaticConfig({
      linkToken: s_sourceTokens[0],
      chainSelector: SOURCE_CHAIN_SELECTOR,
      maxNopFeesJuels: MAX_NOP_FEES_JUELS,
      rmnProxy: address(0)
    });

    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp = new EVM2EVMMultiOnRampHelper(
      staticConfig,
      dynamicConfig,
      destChainConfigArgs,
      getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getMultiOnRampNopsAndWeights()
    );
  }
}

contract EVM2EVMMultiOnRamp_applyDestChainConfigUpdates is EVM2EVMMultiOnRampSetup {
  function test_Fuzz_applyDestChainConfigUpdates_Success(
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArgs
  ) public {
    vm.assume(destChainConfigArgs.destChainSelector != 0);
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

    assertDestChainConfigsEqual(
      expectedDestChainConfig, s_onRamp.getDestChainConfig(destChainConfigArgs.destChainSelector)
    );
  }

  function test_applyDestChainConfigUpdates_Success() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](2);
    destChainConfigArgs[0] = generateDestChainConfigArgs()[0];
    destChainConfigArgs[0].dynamicConfig.isEnabled = false;
    destChainConfigArgs[1] = generateDestChainConfigArgs()[0];
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
    assertDestChainConfigsEqual(expectedDestChainConfig0, gotDestChainConfig0);
    assertDestChainConfigsEqual(expectedDestChainConfig1, gotDestChainConfig1);
  }

  function test_applyDestChainConfigUpdatesZeroIntput() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs =
      new EVM2EVMMultiOnRamp.DestChainConfigArgs[](0);

    vm.recordLogs();
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);

    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_InvalidDestChainConfigDestChainSelectorEqZero_Revert() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    // destChainSelector == address(0)
    destChainConfigArg.destChainSelector = 0;
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_onRamp.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestChainConfigNewPrevOnRampOnExistingChain_Revert() public {
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
    EVM2EVMMultiOnRamp.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    // Changing prevOnRamp on already set destination chain
    destChainConfigArg.prevOnRamp = address(1);
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOnRamp.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
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
    EVM2EVMMultiOnRampSetup.setUp();

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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(_messageToEvent(expectedMessage, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ForwardFromRouter_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = 1234567890;
    IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_onRamp), feeAmount);

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(_messageToEvent(message, 1, 1, feeAmount, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
  }

  function test_ShouldIncrementSeqNumAndNonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    for (uint64 i = 1; i < 4; ++i) {
      uint64 nonceBefore = s_onRamp.getSenderNonce(DEST_CHAIN_SELECTOR, OWNER);

      vm.expectEmit();
      emit EVM2EVMMultiOnRamp.CCIPSendRequested(_messageToEvent(message, i, i, 0, OWNER));

      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      uint64 nonceAfter = s_onRamp.getSenderNonce(DEST_CHAIN_SELECTOR, OWNER);
      assertEq(nonceAfter, nonceBefore + 1);
    }
  }

  function test_ShouldStoreLinkFees() public {
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
    uint256 feeTokenPrice = s_priceRegistry.getTokenPrice(s_sourceTokens[1]).value;
    uint256 linkTokenPrice = s_priceRegistry.getTokenPrice(s_sourceFeeToken).value;
    uint256 conversionRate = (feeTokenPrice * 1e18) / linkTokenPrice;
    uint256 expectedJuels = (feeAmount * conversionRate) / 1e18;

    assertEq(s_onRamp.getNopFeesJuels(), expectedJuels);
  }

  // Make sure any valid sender, receiver and feeAmount can be handled.
  // @TODO Temporarily setting lower fuzz run as 256 triggers snapshot gas off by 1 error.
  // https://github.com/foundry-rs/foundry/issues/5689
  /// forge-dynamicConfig: default.fuzz.runs = 32
  /// forge-dynamicConfig: ccip.fuzz.runs = 32
  function test_Fuzz_ForwardFromRouter_Success(address originalSender, address receiver, uint96 feeTokenAmount) public {
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
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(expectedEvent);

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
    Internal.PriceUpdates memory priceUpdates = getSingleTokenPriceUpdateStruct(s_sourceTokens[0], tokenPrice);
    s_priceRegistry.updatePrices(priceUpdates);
    vm.startPrank(address(s_sourceRouter));

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector,
        getOutboundRateLimiterConfig().capacity,
        (message.tokenAmounts[0].amount * tokenPrice) / 1e18
      )
    );
    // Expect to fail from ARL
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

    // Configure ARL off for token
    EVM2EVMMultiOnRamp.TokenTransferFeeConfig memory tokenTransferFeeConfig =
      s_onRamp.getTokenTransferFeeConfig(s_sourceTokens[0]);
    EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs({
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

  // Reverts

  function test_Paused_Revert() public {
    // We pause by disabling the whitelist
    vm.stopPrank();
    vm.startPrank(OWNER);
    address router = address(0);
    s_onRamp.setDynamicConfig(generateDynamicMultiOnRampConfig(router, address(2), address(s_tokenAdminRegistry)));
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
    s_mockRMN.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
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

  function test_MaxCapacityExceeded_Revert() public {
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

  function test_PriceNotFoundForToken_Revert() public {
    // Set token price to 0
    vm.stopPrank();
    vm.startPrank(OWNER);
    s_priceRegistry.updatePrices(getSingleTokenPriceUpdateStruct(CUSTOM_TOKEN, 0));

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

  function test_MaxFeeBalanceReached_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(EVM2EVMMultiOnRamp.MaxFeeBalanceReached.selector);

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, MAX_NOP_FEES_JUELS + 1, OWNER);
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
      new EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[](1);
    tokenTransferFeeConfigArgs[0] = EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs({
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
    EVM2EVMMultiOnRampSetup.setUp();

    // Add additional pool addresses for test tokens to mark them as supported
    s_tokenAdminRegistry.registerAdministratorPermissioned(s_sourceRouter.getWrappedNative(), OWNER);
    s_tokenAdminRegistry.registerAdministratorPermissioned(CUSTOM_TOKEN, OWNER);

    LockReleaseTokenPool wrappedNativePool = new LockReleaseTokenPool(
      IERC20(s_sourceRouter.getWrappedNative()), new address[](0), address(s_mockRMN), true, address(s_sourceRouter)
    );

    TokenPool.ChainUpdate[] memory wrappedNativeChainUpdate = new TokenPool.ChainUpdate[](1);
    wrappedNativeChainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(address(111111)),
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
      remotePoolAddress: abi.encode(makeAddr("random")),
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
    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = generateDestChainConfigArgs();
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

contract EVM2EVMMultiOnRamp_getFee is EVM2EVMMultiOnRamp_getFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      EVM2EVMMultiOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD =
        (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) * feeTokenConfig.premiumMultiplierWeiPerEth);
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
    EVM2EVMMultiOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

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

      EVM2EVMMultiOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);
      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + customDataSize * DEST_GAS_PER_PAYLOAD_BYTE;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD =
        (configUSDCentToWei(feeTokenConfig.networkFeeUSDCents) * feeTokenConfig.premiumMultiplierWeiPerEth);
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
      EVM2EVMMultiOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);
      uint32 tokenGasOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token).destGasOverhead;
      uint32 destBytesOverhead = s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[0].token).destBytesOverhead;
      uint32 tokenBytesOverhead =
        destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * feeTokenConfig.premiumMultiplierWeiPerEth);
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
      EVM2EVMMultiOnRamp.FeeTokenConfig memory feeTokenConfig = s_onRamp.getFeeTokenConfig(message.feeToken);

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

      uint256 gasUsed =
        customGasLimit + DEST_GAS_OVERHEAD + message.data.length * DEST_GAS_PER_PAYLOAD_BYTE + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * feeTokenConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_onRamp.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * feeTokenConfig.premiumMultiplierWeiPerEth);
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

contract EVM2EVMMultiOnRamp_setDynamicConfig is EVM2EVMMultiOnRampSetup {
  function test_SetDynamicConfig_Success() public {
    EVM2EVMMultiOnRamp.StaticConfig memory staticConfig = s_onRamp.getStaticConfig();
    EVM2EVMMultiOnRamp.DynamicConfig memory newConfig = EVM2EVMMultiOnRamp.DynamicConfig({
      router: address(2134),
      priceRegistry: address(23423),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.ConfigSet(staticConfig, newConfig);

    s_onRamp.setDynamicConfig(newConfig);

    EVM2EVMMultiOnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(newConfig.router, gotDynamicConfig.router);
    assertEq(newConfig.priceRegistry, gotDynamicConfig.priceRegistry);
    assertEq(newConfig.tokenAdminRegistry, gotDynamicConfig.tokenAdminRegistry);
  }

  // Reverts

  function test_SetConfigInvalidConfig_Revert() public {
    EVM2EVMMultiOnRamp.DynamicConfig memory newConfig = EVM2EVMMultiOnRamp.DynamicConfig({
      router: address(1),
      priceRegistry: address(23423),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });

    // Invalid price reg reverts.
    newConfig.priceRegistry = address(0);
    vm.expectRevert(EVM2EVMMultiOnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);

    // Succeeds if valid
    newConfig.priceRegistry = address(23423);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_SetConfigOnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(generateDynamicMultiOnRampConfig(address(1), address(2), address(3)));
    vm.startPrank(ADMIN);
    vm.expectRevert("Only callable by owner");
    s_onRamp.setDynamicConfig(generateDynamicMultiOnRampConfig(address(1), address(2), address(3)));
  }
}
