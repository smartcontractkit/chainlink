// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICommitStore} from "../interfaces/ICommitStore.sol";
import {IEVM2AnyOnRamp} from "../interfaces/IEVM2AnyOnRamp.sol";

import {NonceManager} from "../NonceManager.sol";
import {Router} from "../Router.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {OffRamp} from "../offRamp/OffRamp.sol";
import {OnRamp} from "../onRamp/OnRamp.sol";
import {BaseTest} from "./BaseTest.t.sol";
import {EVM2EVMOffRampHelper} from "./helpers/EVM2EVMOffRampHelper.sol";
import {OnRampHelper} from "./helpers/OnRampHelper.sol";
import {OffRampSetup} from "./offRamp/OffRampSetup.t.sol";
import {OnRampSetup} from "./onRamp/OnRampSetup.t.sol";

import {Test} from "forge-std/Test.sol";

contract NonceManager_typeAndVersion is Test {
  NonceManager private s_nonceManager;

  function setUp() public {
    s_nonceManager = new NonceManager(new address[](0));
  }

  function test_typeAndVersion() public view {
    assertEq(s_nonceManager.typeAndVersion(), "NonceManager 1.6.0-dev");
  }
}

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

contract NonceManager_applyPreviousRampsUpdates is OnRampSetup {
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

contract NonceManager_OnRampUpgrade is OnRampSetup {
  uint256 internal constant FEE_AMOUNT = 1234567890;
  OnRampHelper internal s_prevOnRamp;

  function setUp() public virtual override {
    super.setUp();

    (s_prevOnRamp,) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry)
    );

    // Since the previous onRamp is not a 1.5 ramp it doesn't have the getSenderNonce function. We mock it to return 0
    vm.mockCall(address(s_prevOnRamp), abi.encodeWithSelector(IEVM2AnyOnRamp.getSenderNonce.selector), abi.encode(0));

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(s_prevOnRamp), address(0)));
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    (s_onRamp, s_metadataHash) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry)
    );

    vm.startPrank(address(s_sourceRouter));
  }

  function test_Upgrade_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, FEE_AMOUNT, OWNER));

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
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, startNonce + 2, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 2, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // after another send, nonce should be 3, and sequence number be 2
    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 2, _messageToEvent(message, 2, startNonce + 3, FEE_AMOUNT, OWNER));
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
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, FEE_AMOUNT, newSender));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, newSender);
  }
}

contract NonceManager_OffRampUpgrade is OffRampSetup {
  EVM2EVMOffRampHelper internal s_prevOffRamp;

  address internal constant SINGLE_LANE_ON_RAMP_ADDRESS_1 = abi.decode(ON_RAMP_ADDRESS_1, (address));
  address internal constant SINGLE_LANE_ON_RAMP_ADDRESS_2 = abi.decode(ON_RAMP_ADDRESS_2, (address));
  address internal constant SINGLE_LANE_ON_RAMP_ADDRESS_3 = abi.decode(ON_RAMP_ADDRESS_3, (address));

  function setUp() public virtual override {
    super.setUp();

    s_prevOffRamp = new EVM2EVMOffRampHelper();

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] = NonceManager.PreviousRampsArgs(
      SOURCE_CHAIN_SELECTOR_1, NonceManager.PreviousRamps(address(0), address(s_prevOffRamp))
    );

    s_inboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_1
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_2,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_2
    });
    sourceChainConfigs[2] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_3,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_3
    });

    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);

    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_Upgraded_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_NoPrevOffRampForChain_Success() public {
    address[] memory senders = new address[](1);
    senders[0] = OWNER;

    uint64 startNonceChain3 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, abi.encode(senders[0]));
    s_prevOffRamp.execute(senders);

    // Nonce unchanged for chain 3
    assertEq(startNonceChain3, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, abi.encode(senders[0])));

    Internal.Any2EVMRampMessage[] memory messagesChain3 =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);

    vm.recordLogs();

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messagesChain3), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_3,
      messagesChain3[0].header.sequenceNumber,
      messagesChain3[0].header.messageId,
      _hashMessage(messagesChain3[0], ON_RAMP_ADDRESS_3),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(
      startNonceChain3 + 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, messagesChain3[0].sender)
    );
  }

  function test_UpgradedSenderNoncesReadsPreviousRamp_Success() public {
    address[] memory senders = new address[](1);
    senders[0] = OWNER;

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0]));

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOffRamp.execute(senders);

      assertEq(startNonce + i, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0])));
    }
  }

  function test_UpgradedNonceStartsAtV1Nonce_Success() public {
    address[] memory senders = new address[](1);
    senders[0] = OWNER;

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0]));
    s_prevOffRamp.execute(senders);

    assertEq(startNonce + 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0])));

    Internal.Any2EVMRampMessage[] memory messagesMultiRamp =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messagesMultiRamp[0].header.nonce++;
    messagesMultiRamp[0].header.messageId = _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messagesMultiRamp), new OffRamp.GasLimitOverride[](0)
    );

    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messagesMultiRamp[0].header.sequenceNumber,
      messagesMultiRamp[0].header.messageId,
      _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(
      startNonce + 2, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messagesMultiRamp[0].sender)
    );

    messagesMultiRamp[0].header.nonce++;
    messagesMultiRamp[0].header.sequenceNumber++;
    messagesMultiRamp[0].header.messageId = _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messagesMultiRamp), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messagesMultiRamp[0].header.sequenceNumber,
      messagesMultiRamp[0].header.messageId,
      _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(
      startNonce + 3, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messagesMultiRamp[0].sender)
    );
  }

  function test_UpgradedNonceNewSenderStartsAtZero_Success() public {
    address[] memory senders = new address[](1);
    senders[0] = OWNER;

    s_prevOffRamp.execute(senders);

    Internal.Any2EVMRampMessage[] memory messagesMultiRamp =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    bytes memory newSender = abi.encode(address(1234567));
    messagesMultiRamp[0].sender = newSender;
    messagesMultiRamp[0].header.messageId = _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1);

    // new sender nonce in new offramp should go from 0 -> 1
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, newSender), 0);
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messagesMultiRamp), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messagesMultiRamp[0].header.sequenceNumber,
      messagesMultiRamp[0].header.messageId,
      _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, newSender), 1);
  }

  function test_UpgradedOffRampNonceSkipsIfMsgInFlight_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    address newSender = address(1234567);
    messages[0].sender = abi.encode(newSender);
    messages[0].header.nonce = 2;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);

    // new offramp sees msg nonce higher than senderNonce
    // it waits for previous offramp to execute
    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].header.nonce, messages[0].sender);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    assertEq(startNonce, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender));

    address[] memory senders = new address[](1);
    senders[0] = newSender;

    // previous offramp executes msg and increases nonce
    s_prevOffRamp.execute(senders);
    assertEq(startNonce + 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0])));

    messages[0].header.nonce = 2;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    // new offramp is able to execute
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );

    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(startNonce + 2, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender));
  }
}
