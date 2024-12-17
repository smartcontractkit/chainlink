// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../libraries/Internal.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract OffRamp_batchExecute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_SingleReport() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);

    vm.recordLogs();
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
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

  function test_MultipleReportsSameChain() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages1[0].sender);
    vm.recordLogs();
    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](2));

    Vm.Log[] memory logs = vm.getRecordedLogs();
    _assertExecutionStateChangedEventLogs(
      logs,
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      logs,
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      logs,
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      _hashMessage(messages2[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages1[0].sender), nonceBefore);
  }

  function test_MultipleReportsDifferentChains() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, 1);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    vm.recordLogs();

    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](2));

    Vm.Log[] memory logs = vm.getRecordedLogs();

    _assertExecutionStateChangedEventLogs(
      logs,
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      logs,
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      logs,
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      _hashMessage(messages2[0], ON_RAMP_ADDRESS_3),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint64 nonceChain1 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages1[0].sender);
    uint64 nonceChain3 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, messages2[0].sender);

    assertTrue(nonceChain1 != nonceChain3);
    assertGt(nonceChain1, 0);
    assertGt(nonceChain3, 0);
  }

  function test_MultipleReportsDifferentChainsSkipCursedChain() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, true);

    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, 1);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    vm.recordLogs();

    vm.expectEmit();
    emit OffRamp.SkippedReportExecution(SOURCE_CHAIN_SELECTOR_1);

    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](2));

    Vm.Log[] memory logs = vm.getRecordedLogs();

    for (uint256 i = 0; i < logs.length; ++i) {
      if (logs[i].topics[0] == OffRamp.ExecutionStateChanged.selector) {
        uint64 logSourceChainSelector = uint64(uint256(logs[i].topics[1]));
        uint64 logSequenceNumber = uint64(uint256(logs[i].topics[2]));
        bytes32 logMessageId = bytes32(logs[i].topics[3]);
        (bytes32 logMessageHash, uint8 logState,,) = abi.decode(logs[i].data, (bytes32, uint8, bytes, uint256));
        assertEq(logMessageId, messages2[0].header.messageId);
        assertEq(logSourceChainSelector, messages2[0].header.sourceChainSelector);
        assertEq(logSequenceNumber, messages2[0].header.sequenceNumber);
        assertEq(logMessageId, messages2[0].header.messageId);
        assertEq(logMessageHash, _hashMessage(messages2[0], ON_RAMP_ADDRESS_3));
        assertEq(logState, uint8(Internal.MessageExecutionState.SUCCESS));
      }
    }
  }

  function test_MultipleReportsSkipDuplicate() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    vm.recordLogs();
    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](2));
    _assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_Unhealthy() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, true);
    vm.expectEmit();
    emit OffRamp.SkippedReportExecution(SOURCE_CHAIN_SELECTOR_1);
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[][](1)
    );

    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, false);

    vm.recordLogs();
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[][](1)
    );

    _assertNoEmit(OffRamp.SkippedReportExecution.selector);
  }

  // Reverts
  function test_RevertWhen_ZeroReports() public {
    vm.expectRevert(OffRamp.EmptyBatch.selector);
    s_offRamp.batchExecute(new Internal.ExecutionReport[](0), new OffRamp.GasLimitOverride[][](1));
  }

  function test_RevertWhen_OutOfBoundsGasLimitsAccess() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    vm.expectRevert();
    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](1));
  }
}
