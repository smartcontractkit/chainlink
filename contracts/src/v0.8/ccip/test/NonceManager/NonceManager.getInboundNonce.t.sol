// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {NonceManager} from "../../NonceManager.sol";
import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {OffRampSetup} from "../offRamp/OffRamp/OffRampSetup.t.sol";

contract NonceManager_getInboundNonce is OffRampSetup {
  EVM2EVMOffRampHelper internal s_prevOffRamp;

  address internal constant SINGLE_LANE_ON_RAMP_ADDRESS_1 = abi.decode(ON_RAMP_ADDRESS_1, (address));
  address internal constant SINGLE_LANE_ON_RAMP_ADDRESS_2 = abi.decode(ON_RAMP_ADDRESS_2, (address));
  address internal constant SINGLE_LANE_ON_RAMP_ADDRESS_3 = abi.decode(ON_RAMP_ADDRESS_3, (address));

  function setUp() public virtual override {
    super.setUp();

    s_prevOffRamp = new EVM2EVMOffRampHelper();

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
      prevRamps: NonceManager.PreviousRamps(address(0), address(s_prevOffRamp)),
      overrideExistingRamps: false
    });

    s_inboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_1,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_2,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_2,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[2] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_3,
      isEnabled: true,
      onRamp: ON_RAMP_ADDRESS_3,
      isRMNVerificationDisabled: false
    });

    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);

    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_getInboundNonce_Upgraded() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_getInboundNonce_NoPrevOffRampForChain() public {
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
    _assertExecutionStateChangedEventLogs(
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

  function test_getInboundNonce_UpgradedSenderNoncesReadsPreviousRamp() public {
    address[] memory senders = new address[](1);
    senders[0] = OWNER;

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0]));

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOffRamp.execute(senders);

      assertEq(startNonce + i, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0])));
    }
  }

  function test_getInboundNonce_UpgradedNonceStartsAtV1Nonce() public {
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

    _assertExecutionStateChangedEventLogs(
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
    _assertExecutionStateChangedEventLogs(
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

  function test_getInboundNonce_UpgradedNonceNewSenderStartsAtZero() public {
    address[] memory senders = new address[](1);
    senders[0] = OWNER;

    s_prevOffRamp.execute(senders);

    Internal.Any2EVMRampMessage[] memory messagesMultiRamp =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    bytes memory newSender = abi.encode(address(1234567));
    messagesMultiRamp[0].sender = newSender;
    messagesMultiRamp[0].header.messageId = _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1);

    // new sender nonce in new offRamp should go from 0 -> 1
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, newSender), 0);
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messagesMultiRamp), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messagesMultiRamp[0].header.sequenceNumber,
      messagesMultiRamp[0].header.messageId,
      _hashMessage(messagesMultiRamp[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, newSender), 1);
  }

  function test_getInboundNonce_UpgradedOffRampNonceSkipsIfMsgInFlight() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    address newSender = address(1234567);
    messages[0].sender = abi.encode(newSender);
    messages[0].header.nonce = 2;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    uint64 startNonce = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);

    // new offRamp sees msg nonce higher than senderNonce
    // it waits for previous offRamp to execute
    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].header.nonce, messages[0].sender);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    assertEq(startNonce, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender));

    address[] memory senders = new address[](1);
    senders[0] = newSender;

    // previous offRamp executes msg and increases nonce
    s_prevOffRamp.execute(senders);
    assertEq(startNonce + 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(senders[0])));

    messages[0].header.nonce = 2;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    // new offRamp is able to execute
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );

    _assertExecutionStateChangedEventLogs(
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
