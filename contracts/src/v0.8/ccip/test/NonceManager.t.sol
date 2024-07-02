// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {NonceManager} from "../NonceManager.sol";
import {ICommitStore} from "../interfaces/ICommitStore.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {EVM2EVMMultiOffRamp} from "../offRamp/EVM2EVMMultiOffRamp.sol";
import {EVM2EVMMultiOnRamp} from "../onRamp/EVM2EVMMultiOnRamp.sol";
import {EVM2EVMOnRamp} from "../onRamp/EVM2EVMOnRamp.sol";

import {BaseTest} from "./BaseTest.t.sol";
import {EVM2EVMMultiOnRampHelper} from "./helpers/EVM2EVMMultiOnRampHelper.sol";
import {EVM2EVMOffRampHelper} from "./helpers/EVM2EVMOffRampHelper.sol";
import {EVM2EVMOnRampHelper} from "./helpers/EVM2EVMOnRampHelper.sol";
import {MockCommitStore} from "./mocks/MockCommitStore.sol";
import {EVM2EVMMultiOffRampSetup} from "./offRamp/EVM2EVMMultiOffRampSetup.t.sol";
import {EVM2EVMMultiOnRampSetup} from "./onRamp/EVM2EVMMultiOnRampSetup.t.sol";

contract NonceManager_NonceIncrementation is BaseTest {
  NonceManager private s_nonceManager;

  function setUp() public override {
    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(this);
    s_nonceManager = new NonceManager(authorizedCallers);
  }

  function test_getIncrementedOutboundNonce_Success() public {
    address sender = address(this);

    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 0);

    uint64 outboundNonce = s_nonceManager.getIncrementedOutboundNonce(DEST_CHAIN_SELECTOR, sender);
    assertEq(outboundNonce, 1);
  }

  function test_incrementInboundNonce_Success() public {
    address sender = address(this);

    s_nonceManager.incrementInboundNonce(SOURCE_CHAIN_SELECTOR, 1, abi.encode(sender));

    assertEq(s_nonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR, abi.encode(sender)), 1);
  }

  function test_incrementInboundNonce_Skip() public {
    address sender = address(this);
    uint64 expectedNonce = 2;

    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR, expectedNonce, abi.encode(sender));

    s_nonceManager.incrementInboundNonce(SOURCE_CHAIN_SELECTOR, expectedNonce, abi.encode(sender));

    assertEq(s_nonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR, abi.encode(sender)), 0);
  }

  function test_incrementNoncesInboundAndOutbound_Success() public {
    address sender = address(this);

    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 0);
    uint64 outboundNonce = s_nonceManager.getIncrementedOutboundNonce(DEST_CHAIN_SELECTOR, sender);
    assertEq(outboundNonce, 1);

    // Inbound nonce unchanged
    assertEq(s_nonceManager.getInboundNonce(DEST_CHAIN_SELECTOR, abi.encode(sender)), 0);

    s_nonceManager.incrementInboundNonce(DEST_CHAIN_SELECTOR, 1, abi.encode(sender));
    assertEq(s_nonceManager.getInboundNonce(DEST_CHAIN_SELECTOR, abi.encode(sender)), 1);

    // Outbound nonce unchanged
    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 1);
  }
}

contract NonceManager_applyPreviousRampsUpdates is EVM2EVMMultiOnRampSetup {
  function test_SingleRampUpdate() public {
    address prevOnRamp = makeAddr("prevOnRamp");
    address prevOffRamp = makeAddr("prevOffRamp");
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp, prevOffRamp));

    vm.expectEmit();
    emit NonceManager.PreviousRampsUpdated(DEST_CHAIN_SELECTOR, previousRamps[0].prevRamps);

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    _assertPreviousRampsEqual(s_outboundNonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR), previousRamps[0].prevRamps);
  }

  function test_MultipleRampsUpdates() public {
    address prevOnRamp1 = makeAddr("prevOnRamp1");
    address prevOnRamp2 = makeAddr("prevOnRamp2");
    address prevOffRamp1 = makeAddr("prevOffRamp1");
    address prevOffRamp2 = makeAddr("prevOffRamp2");
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](2);
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp1, prevOffRamp1));
    previousRamps[1] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR + 1, NonceManager.PreviousRamps(prevOnRamp2, prevOffRamp2));

    vm.expectEmit();
    emit NonceManager.PreviousRampsUpdated(DEST_CHAIN_SELECTOR, previousRamps[0].prevRamps);
    vm.expectEmit();
    emit NonceManager.PreviousRampsUpdated(DEST_CHAIN_SELECTOR + 1, previousRamps[1].prevRamps);

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    _assertPreviousRampsEqual(s_outboundNonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR), previousRamps[0].prevRamps);
    _assertPreviousRampsEqual(
      s_outboundNonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR + 1), previousRamps[1].prevRamps
    );
  }

  function test_ZeroInput() public {
    vm.recordLogs();
    s_outboundNonceManager.applyPreviousRampsUpdates(new NonceManager.PreviousRampsArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_PreviousRampAlreadySetOnRamp_Revert() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOnRamp = makeAddr("prevOnRamp");
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp, address(0)));

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp, address(0)));

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function test_PreviousRampAlreadySetOffRamp_Revert() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOffRamp = makeAddr("prevOffRamp");
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(0), prevOffRamp));

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(0), prevOffRamp));

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function test_PreviousRampAlreadySetOnRampAndOffRamp_Revert() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOnRamp = makeAddr("prevOnRamp");
    address prevOffRamp = makeAddr("prevOffRamp");
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp, prevOffRamp));

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp, prevOffRamp));

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function _assertPreviousRampsEqual(
    NonceManager.PreviousRamps memory a,
    NonceManager.PreviousRamps memory b
  ) internal pure {
    assertEq(a.prevOnRamp, b.prevOnRamp);
    assertEq(a.prevOffRamp, b.prevOffRamp);
  }
}

contract NonceManager_OnRampUpgrade is EVM2EVMMultiOnRampSetup {
  uint256 internal constant FEE_AMOUNT = 1234567890;
  EVM2EVMOnRampHelper internal s_prevOnRamp;

  function setUp() public virtual override {
    super.setUp();

    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeTokenConfigArgs = new EVM2EVMOnRamp.FeeTokenConfigArgs[](1);
    feeTokenConfigArgs[0] = EVM2EVMOnRamp.FeeTokenConfigArgs({
      token: s_sourceFeeToken,
      networkFeeUSDCents: 1_00, // 1 USD
      gasMultiplierWeiPerEth: 1e18, // 1x
      premiumMultiplierWeiPerEth: 5e17, // 0.5x
      enabled: true
    });

    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfig =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);

    tokenTransferFeeConfig[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: s_sourceFeeToken,
      minFeeUSDCents: 1_00, // 1 USD
      maxFeeUSDCents: 1000_00, // 1,000 USD
      deciBps: 2_5, // 2.5 bps, or 0.025%
      destGasOverhead: 40_000,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      aggregateRateLimitEnabled: true
    });

    s_prevOnRamp = new EVM2EVMOnRampHelper(
      EVM2EVMOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        destChainSelector: DEST_CHAIN_SELECTOR,
        defaultTxGasLimit: GAS_LIMIT,
        maxNopFeesJuels: MAX_NOP_FEES_JUELS,
        prevOnRamp: address(0),
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      EVM2EVMOnRamp.DynamicConfig({
        router: address(s_sourceRouter),
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        destGasOverhead: DEST_GAS_OVERHEAD,
        destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
        destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
        destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
        destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
        priceRegistry: address(s_priceRegistry),
        maxDataBytes: MAX_DATA_SIZE,
        maxPerMsgGasLimit: MAX_GAS_LIMIT,
        defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
        defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        defaultTokenDestBytesOverhead: DEFAULT_TOKEN_BYTES_OVERHEAD,
        enforceOutOfOrder: false
      }),
      RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e15}),
      feeTokenConfigArgs,
      tokenTransferFeeConfig,
      new EVM2EVMOnRamp.NopAndWeight[](0)
    );

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(s_prevOnRamp), address(0)));
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    destChainConfigArgs[0].prevOnRamp = address(s_prevOnRamp);

    (s_onRamp, s_metadataHash) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, address(s_sourceRouter), address(s_outboundNonceManager), address(s_tokenAdminRegistry)
    );

    vm.startPrank(address(s_sourceRouter));
  }

  function test_Upgrade_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);
  }

  function test_UpgradeSenderNoncesReadsPreviousRamp_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 startNonce = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      assertEq(startNonce + i, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));
    }
  }

  function test_UpgradeNonceStartsAtV1Nonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 startNonce = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);

    // send 1 message from previous onramp
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 1, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // new onramp nonce should start from 2, while sequence number start from 1
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, startNonce + 2, FEE_AMOUNT, OWNER)
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 2, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // after another send, nonce should be 3, and sequence number be 2
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(message, 2, startNonce + 3, FEE_AMOUNT, OWNER)
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 3, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));
  }

  function test_UpgradeNonceNewSenderStartsAtZero_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    // send 1 message from previous onramp from OWNER
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    address newSender = address(1234567);
    // new onramp nonce should start from 1 for new sender
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, FEE_AMOUNT, newSender)
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, newSender);
  }
}

contract NonceManager_OffRampUpgrade is EVM2EVMMultiOffRampSetup {
  EVM2EVMOffRampHelper internal s_prevOffRamp;
  EVM2EVMOffRampHelper[] internal s_nestedPrevOffRamps;

  function setUp() public virtual override {
    super.setUp();

    ICommitStore mockPrevCommitStore = new MockCommitStore();
    s_prevOffRamp = _deploySingleLaneOffRamp(
      mockPrevCommitStore, s_destRouter, address(0), SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1
    );

    s_nestedPrevOffRamps = new EVM2EVMOffRampHelper[](2);
    s_nestedPrevOffRamps[0] = _deploySingleLaneOffRamp(
      mockPrevCommitStore, s_destRouter, address(0), SOURCE_CHAIN_SELECTOR_2, ON_RAMP_ADDRESS_2
    );
    s_nestedPrevOffRamps[1] = _deploySingleLaneOffRamp(
      mockPrevCommitStore, s_destRouter, address(s_nestedPrevOffRamps[0]), SOURCE_CHAIN_SELECTOR_2, ON_RAMP_ADDRESS_2
    );

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](3);
    previousRamps[0] = NonceManager.PreviousRampsArgs(
      SOURCE_CHAIN_SELECTOR_1, NonceManager.PreviousRamps(address(0), address(s_prevOffRamp))
    );
    previousRamps[1] = NonceManager.PreviousRampsArgs(
      SOURCE_CHAIN_SELECTOR_2, NonceManager.PreviousRamps(address(0), address(s_nestedPrevOffRamps[1]))
    );
    previousRamps[2] =
      NonceManager.PreviousRampsArgs(SOURCE_CHAIN_SELECTOR_3, NonceManager.PreviousRamps(ON_RAMP_ADDRESS_3, address(0)));
    s_inboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_1
    });
    sourceChainConfigs[1] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_2,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_2
    });
    sourceChainConfigs[2] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_3,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_3
    });

    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);

    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_Upgraded_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function test_NoPrevOffRampForChain_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    uint64 startNonceChain3 =
      s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, abi.encode(messages[0].sender));
    s_prevOffRamp.execute(_generateSingleRampReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    // Nonce unchanged for chain 3
    assertEq(
      startNonceChain3, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, abi.encode(messages[0].sender))
    );

    Internal.EVM2EVMMessage[] memory messagesChain3 =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_3,
      messagesChain3[0].sequenceNumber,
      messagesChain3[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messagesChain3), new uint256[](0)
    );
    assertEq(
      startNonceChain3 + 1,
      s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, abi.encode(messagesChain3[0].sender))
    );
  }

  function test_UpgradedSenderNoncesReadsPreviousRamp_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender));

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOffRamp.execute(_generateSingleRampReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

      // messages contains a single message - update for the next execution
      messages[0].nonce++;
      messages[0].sequenceNumber++;
      messages[0].messageId = Internal._hash(messages[0], s_prevOffRamp.metadataHash());

      assertEq(
        startNonce + i, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender))
      );
    }
  }

  function test_UpgradedSenderNoncesReadsPreviousRampTransitive_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_2, ON_RAMP_ADDRESS_2);
    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_2, abi.encode(messages[0].sender));

    for (uint64 i = 1; i < 4; ++i) {
      s_nestedPrevOffRamps[0].execute(
        _generateSingleRampReportFromMessages(SOURCE_CHAIN_SELECTOR_2, messages), new uint256[](0)
      );

      // messages contains a single message - update for the next execution
      messages[0].nonce++;
      messages[0].sequenceNumber++;
      messages[0].messageId = Internal._hash(messages[0], s_nestedPrevOffRamps[0].metadataHash());

      // Read through prev sender nonce through prevOffRamp -> prevPrevOffRamp
      assertEq(
        startNonce + i, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_2, abi.encode(messages[0].sender))
      );
    }
  }

  function test_UpgradedNonceStartsAtV1Nonce_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender));
    s_prevOffRamp.execute(_generateSingleRampReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    assertEq(
      startNonce + 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender))
    );

    messages[0].nonce++;
    messages[0].messageId =
      Internal._hash(messages[0], s_offRamp.metadataHash(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1));

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(
      startNonce + 2, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender))
    );

    messages[0].nonce++;
    messages[0].sequenceNumber++;
    messages[0].messageId =
      Internal._hash(messages[0], s_offRamp.metadataHash(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1));

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(
      startNonce + 3, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender))
    );
  }

  function test_UpgradedNonceNewSenderStartsAtZero_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    s_prevOffRamp.execute(_generateSingleRampReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    address newSender = address(1234567);
    messages[0].sender = newSender;
    messages[0].messageId =
      Internal._hash(messages[0], s_offRamp.metadataHash(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1));

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    // new sender nonce in new offramp should go from 0 -> 1
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(newSender)), 0);
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(newSender)), 1);
  }

  function test_UpgradedOffRampNonceSkipsIfMsgInFlight_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    address newSender = address(1234567);
    messages[0].sender = newSender;
    messages[0].nonce = 2;
    messages[0].messageId =
      Internal._hash(messages[0], s_offRamp.metadataHash(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1));

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender));

    // new offramp sees msg nonce higher than senderNonce
    // it waits for previous offramp to execute
    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].nonce, abi.encode(newSender));
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(startNonce, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender)));

    messages[0].nonce = 1;
    messages[0].messageId = Internal._hash(messages[0], s_prevOffRamp.metadataHash());

    // previous offramp executes msg and increases nonce
    s_prevOffRamp.execute(_generateSingleRampReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(
      startNonce + 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender))
    );

    messages[0].nonce = 2;
    messages[0].messageId =
      Internal._hash(messages[0], s_offRamp.metadataHash(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1));

    // new offramp is able to execute
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(
      startNonce + 2, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(messages[0].sender))
    );
  }
}
