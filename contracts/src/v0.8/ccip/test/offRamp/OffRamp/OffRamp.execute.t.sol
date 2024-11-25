// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageInterceptor} from "../../../interfaces/IMessageInterceptor.sol";

import {Internal} from "../../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract OffRamp_execute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
  }

  // Asserts that execute completes
  function test_SingleReport() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    vm.recordLogs();

    _execute(reports);

    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_MultipleReports() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    vm.recordLogs();
    _execute(reports);

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
  }

  function test_LargeBatch() public {
    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](10);
    for (uint64 i = 0; i < reports.length; ++i) {
      Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](3);
      messages[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1 + i * 3);
      messages[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2 + i * 3);
      messages[2] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3 + i * 3);

      reports[i] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    }

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    vm.recordLogs();
    _execute(reports);

    Vm.Log[] memory logs = vm.getRecordedLogs();

    for (uint64 i = 0; i < reports.length; ++i) {
      for (uint64 j = 0; j < reports[i].messages.length; ++j) {
        _assertExecutionStateChangedEventLogs(
          logs,
          reports[i].messages[j].header.sourceChainSelector,
          reports[i].messages[j].header.sequenceNumber,
          reports[i].messages[j].header.messageId,
          _hashMessage(reports[i].messages[j], ON_RAMP_ADDRESS_1),
          Internal.MessageExecutionState.SUCCESS,
          ""
        );
      }
    }
  }

  function test_MultipleReportsWithPartialValidationFailures() public {
    _enableInboundMessageInterceptor();

    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    s_inboundMessageInterceptor.setMessageIdValidationState(messages1[0].header.messageId, true);
    s_inboundMessageInterceptor.setMessageIdValidationState(messages2[0].header.messageId, true);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    vm.recordLogs();
    _execute(reports);

    Vm.Log[] memory logs = vm.getRecordedLogs();

    _assertExecutionStateChangedEventLogs(
      logs,
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
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
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );
  }

  // Reverts

  function test_RevertWhen_UnauthorizedTransmitter() public {
    bytes32[2] memory reportContext = [s_configDigestExec, s_configDigestExec];

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_RevertWhen_NoConfig() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    bytes32[2] memory reportContext = [bytes32(""), s_configDigestExec];

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_RevertWhen_NoConfigWithOtherConfigPresent() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    bytes32[2] memory reportContext = [bytes32(""), s_configDigestExec];

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_RevertWhen_WrongConfigWithSigners() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    s_configDigestExec = _getBasicConfigDigest(1, s_validSigners, s_validTransmitters);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    vm.expectRevert(OffRamp.SignatureVerificationNotAllowedInExecutionPlugin.selector);
    s_offRamp.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_ZeroReports() public {
    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](0);

    vm.expectRevert(OffRamp.EmptyBatch.selector);
    _execute(reports);
  }

  function test_RevertWhen_IncorrectArrayType() public {
    bytes32[2] memory reportContext = [s_configDigestExec, s_configDigestExec];

    uint256[] memory wrongData = new uint256[](2);
    wrongData[0] = 1;

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.execute(reportContext, abi.encode(wrongData));
  }

  function test_RevertWhen_NonArray() public {
    bytes32[2] memory reportContext = [s_configDigestExec, s_configDigestExec];

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.execute(reportContext, abi.encode(report));
  }
}
