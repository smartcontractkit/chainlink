// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CallWithExactGas} from "../../../../shared/call/CallWithExactGas.sol";
import {NonceManager} from "../../../NonceManager.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {ConformingReceiver} from "../../helpers/receivers/ConformingReceiver.sol";
import {MaybeRevertMessageReceiver} from "../../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MaybeRevertMessageReceiverNo165} from "../../helpers/receivers/MaybeRevertMessageReceiverNo165.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract OffRamp_executeSingleReport is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_SingleMessageNoTokens() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    messages[0].header.nonce++;
    messages[0].header.sequenceNumber++;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender), nonceBefore);
  }

  function test_SingleMessageNoTokensUnordered() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].header.nonce = 0;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    // Nonce never increments on unordered messages.
    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(
      s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender),
      nonceBefore,
      "nonce must remain unchanged on unordered messages"
    );

    messages[0].header.sequenceNumber++;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    // Nonce never increments on unordered messages.
    nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertEq(
      s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender),
      nonceBefore,
      "nonce must remain unchanged on unordered messages"
    );
  }

  function test_SingleMessageNoTokensOtherChain() public {
    Internal.Any2EVMRampMessage[] memory messagesChain1 =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messagesChain1), new OffRamp.GasLimitOverride[](0)
    );

    uint64 nonceChain1 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messagesChain1[0].sender);
    assertGt(nonceChain1, 0);

    Internal.Any2EVMRampMessage[] memory messagesChain2 =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, messagesChain2[0].sender), 0);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messagesChain2), new OffRamp.GasLimitOverride[](0)
    );
    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, messagesChain2[0].sender), 0);

    // Other chain's nonce is unaffected
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messagesChain1[0].sender), nonceChain1);
  }

  function test_ReceiverError() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    // Nonce should increment on non-strict
    assertEq(uint64(0), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        OffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1)
      )
    );
    assertEq(uint64(1), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
  }

  function test_SkippedIncorrectNonce() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messages[0].header.nonce++;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(
      messages[0].header.sourceChainSelector, messages[0].header.nonce, messages[0].sender
    );

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_SkippedIncorrectNonceStillExecutes() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messages[1].header.nonce++;
    messages[1].header.messageId = _hashMessage(messages[1], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR_1, messages[1].header.nonce, messages[1].sender);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test__execute_SkippedAlreadyExecutedMessage() public {
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

    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test__execute_SkippedAlreadyExecutedMessageUnordered() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].header.nonce = 0;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

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

    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  // Send a message to a contract that does not implement the CCIPReceiver interface
  // This should execute successfully.
  function test_SingleMessageToNonCCIPReceiver() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    messages[0].receiver = address(newReceiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_SingleMessagesNoTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.resumeGasMetering();
    vm.recordLogs();
    s_offRamp.executeSingleReport(report, new OffRamp.GasLimitOverride[](0));
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_TwoMessagesWithTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].header.messageId = _hashMessage(messages[1], ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.resumeGasMetering();
    vm.recordLogs();
    s_offRamp.executeSingleReport(report, new OffRamp.GasLimitOverride[](0));

    Vm.Log[] memory logs = vm.getRecordedLogs();
    _assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[1].header.sequenceNumber,
      messages[1].header.messageId,
      _hashMessage(messages[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_TwoMessagesWithTokensAndGE() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].header.messageId = _hashMessage(messages[1], ON_RAMP_ADDRESS_1);

    assertEq(uint64(0), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );

    Vm.Log[] memory logs = vm.getRecordedLogs();

    _assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    _assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[1].header.sequenceNumber,
      messages[1].header.messageId,
      _hashMessage(messages[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertEq(uint64(2), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
  }

  function testFuzz_InterleavingOrderedAndUnorderedMessages_Success(
    bool[7] memory orderings
  ) public {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](orderings.length);
    // number of tokens needs to be capped otherwise we hit UnsupportedNumberOfTokens.
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](3);
    for (uint256 i = 0; i < 3; ++i) {
      tokenAmounts[i].token = s_sourceTokens[i % s_sourceTokens.length];
      tokenAmounts[i].amount = 1e18;
    }
    uint64 expectedNonce = 0;

    for (uint256 i = 0; i < orderings.length; ++i) {
      messages[i] =
        _generateAny2EVMMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, uint64(i + 1), tokenAmounts, !orderings[i]);
      if (orderings[i]) {
        messages[i].header.nonce = ++expectedNonce;
      }
      messages[i].header.messageId = _hashMessage(messages[i], ON_RAMP_ADDRESS_1);
    }

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER));
    assertEq(uint64(0), nonceBefore, "nonce before exec should be 0");
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );

    Vm.Log[] memory logs = vm.getRecordedLogs();

    // all executions should succeed.
    for (uint256 i = 0; i < orderings.length; ++i) {
      assertEq(
        uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, messages[i].header.sequenceNumber)),
        uint256(Internal.MessageExecutionState.SUCCESS)
      );

      _assertExecutionStateChangedEventLogs(
        logs,
        SOURCE_CHAIN_SELECTOR_1,
        messages[i].header.sequenceNumber,
        messages[i].header.messageId,
        _hashMessage(messages[i], ON_RAMP_ADDRESS_1),
        Internal.MessageExecutionState.SUCCESS,
        ""
      );
    }
    assertEq(
      nonceBefore + expectedNonce, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER))
    );
  }

  function test_InvalidSourcePoolAddress() public {
    address fakePoolAddress = address(0x0000000000333333);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].tokenAmounts[0].sourcePoolAddress = abi.encode(fakePoolAddress);

    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);
    messages[1].header.messageId = _hashMessage(messages[1], ON_RAMP_ADDRESS_1);

    address destPool = s_destPoolByToken[messages[0].tokenAmounts[0].destTokenAddress];

    vm.recordLogs();

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        OffRamp.TokenHandlingError.selector,
        destPool,
        abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, abi.encode(fakePoolAddress))
      )
    );
  }

  function test_WithCurseOnAnotherSourceChain() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_2, true);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_Unhealthy() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, true);

    vm.expectEmit();
    emit OffRamp.SkippedReportExecution(SOURCE_CHAIN_SELECTOR_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[](0)
    );

    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, false);
    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[](0)
    );

    _assertNoEmit(OffRamp.SkippedReportExecution.selector);
  }

  // Reverts

  function test_RevertWhen_MismatchingDestChainSelector() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);
    messages[0].header.destChainSelector = DEST_CHAIN_SELECTOR + 1;

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.InvalidMessageDestChainSelector.selector, messages[0].header.destChainSelector)
    );
    s_offRamp.executeSingleReport(executionReport, new OffRamp.GasLimitOverride[](0));
  }

  function test_RevertWhen_UnhealthySingleChainCurse() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, true);
    vm.expectEmit();
    emit OffRamp.SkippedReportExecution(SOURCE_CHAIN_SELECTOR_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[](0)
    );
    vm.recordLogs();
    // Uncurse should succeed
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, false);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[](0)
    );
    _assertNoEmit(OffRamp.SkippedReportExecution.selector);
  }

  function test_RevertWhen_UnexpectedTokenData() public {
    Internal.ExecutionReport memory report = _generateReportFromMessages(
      SOURCE_CHAIN_SELECTOR_1, _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
    );
    report.offchainTokenData = new bytes[][](report.messages.length + 1);

    vm.expectRevert(OffRamp.UnexpectedTokenData.selector);

    s_offRamp.executeSingleReport(report, new OffRamp.GasLimitOverride[](0));
  }

  function test_RevertWhen_EmptyReport() public {
    vm.expectRevert(abi.encodeWithSelector(OffRamp.EmptyReport.selector, SOURCE_CHAIN_SELECTOR_1));

    s_offRamp.executeSingleReport(
      Internal.ExecutionReport({
        sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
        proofs: new bytes32[](0),
        proofFlagBits: 0,
        messages: new Internal.Any2EVMRampMessage[](0),
        offchainTokenData: new bytes[][](0)
      }),
      new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_RevertWhen_RootNotCommitted() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 0);
    vm.expectRevert(abi.encodeWithSelector(OffRamp.RootNotCommitted.selector, SOURCE_CHAIN_SELECTOR_1));

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
  }

  function test_RevertWhen_ManualExecutionNotYetEnabled() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, BLOCK_TIME);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.ManualExecutionNotYetEnabled.selector, SOURCE_CHAIN_SELECTOR_1));

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
  }

  function test_RevertWhen_NonExistingSourceChain() public {
    uint64 newSourceChainSelector = SOURCE_CHAIN_SELECTOR_1 + 1;
    bytes memory newOnRamp = abi.encode(ON_RAMP_ADDRESS, 1);

    Internal.Any2EVMRampMessage[] memory messages = _generateSingleBasicMessage(newSourceChainSelector, newOnRamp);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.SourceChainNotEnabled.selector, newSourceChainSelector));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(newSourceChainSelector, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_RevertWhen_DisabledSourceChain() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_2, ON_RAMP_ADDRESS_2);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.SourceChainNotEnabled.selector, SOURCE_CHAIN_SELECTOR_2));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_2, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_RevertWhen_TokenDataMismatch() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    report.offchainTokenData[0] = new bytes[](messages[0].tokenAmounts.length + 1);

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.TokenDataMismatch.selector, SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber
      )
    );
    s_offRamp.executeSingleReport(report, new OffRamp.GasLimitOverride[](0));
  }

  function test_RevertWhen_RouterYULCall() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    // gas limit too high, Router's external call should revert
    messages[0].gasLimit = 1e36;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.recordLogs();
    s_offRamp.executeSingleReport(executionReport, new OffRamp.GasLimitOverride[](0));
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector)
    );
  }

  function test_RevertWhen_RetryFailedMessageWithoutManualExecution() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        OffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1)
      )
    );

    // The second time should skip the msg
    vm.expectEmit();
    emit OffRamp.AlreadyAttempted(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function _constructCommitReport(
    bytes32 merkleRoot
  ) internal view returns (OffRamp.CommitReport memory) {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: abi.encode(ON_RAMP_ADDRESS_1),
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: merkleRoot
    });

    return OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });
  }
}
