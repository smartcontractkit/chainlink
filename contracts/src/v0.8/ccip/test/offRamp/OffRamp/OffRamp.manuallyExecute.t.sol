// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {ConformingReceiver} from "../../helpers/receivers/ConformingReceiver.sol";
import {MaybeRevertMessageReceiver} from "../../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {ReentrancyAbuserMultiRamp} from "../../helpers/receivers/ReentrancyAbuserMultiRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {Vm} from "forge-std/Vm.sol";

contract OffRamp_manuallyExecute is OffRampSetup {
  uint32 internal constant MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS = 200_000;

  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();

    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_manuallyExecute() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
    );

    s_reverting_receiver.setRevert(false);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](messages.length);

    vm.recordLogs();
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_manuallyExecute_WithGasOverride() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
    );

    s_reverting_receiver.setRevert(false);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0][0].receiverExecutionGasLimit += 1;
    vm.recordLogs();
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_manuallyExecute_DoesNotRevertIfUntouched() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    assertEq(
      messages[0].header.nonce - 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender)
    );

    s_reverting_receiver.setRevert(true);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    vm.recordLogs();
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        OffRamp.ReceiverError.selector, abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, "")
      )
    );

    assertEq(
      messages[0].header.nonce, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender)
    );
  }

  function test_manuallyExecute_WithMultiReportGasOverride() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](3);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](2);

    for (uint64 i = 0; i < 3; ++i) {
      messages1[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, i + 1);
      messages1[i].receiver = address(s_reverting_receiver);
      messages1[i].header.messageId = _hashMessage(messages1[i], ON_RAMP_ADDRESS_1);
    }

    for (uint64 i = 0; i < 2; ++i) {
      messages2[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, i + 1);
      messages2[i].receiver = address(s_reverting_receiver);
      messages2[i].header.messageId = _hashMessage(messages2[i], ON_RAMP_ADDRESS_3);
    }

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](2));

    s_reverting_receiver.setRevert(false);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](2);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages1);
    gasLimitOverrides[1] = _getGasLimitsFromMessages(messages2);

    for (uint256 i = 0; i < 3; ++i) {
      gasLimitOverrides[0][i].receiverExecutionGasLimit += 1;
    }

    for (uint256 i = 0; i < 2; ++i) {
      gasLimitOverrides[1][i].receiverExecutionGasLimit += 1;
    }

    vm.recordLogs();
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    Vm.Log[] memory logs = vm.getRecordedLogs();

    for (uint256 j = 0; j < 3; ++j) {
      _assertExecutionStateChangedEventLogs(
        logs,
        SOURCE_CHAIN_SELECTOR_1,
        messages1[j].header.sequenceNumber,
        messages1[j].header.messageId,
        _hashMessage(messages1[j], ON_RAMP_ADDRESS_1),
        Internal.MessageExecutionState.SUCCESS,
        ""
      );
    }

    for (uint256 k = 0; k < 2; ++k) {
      _assertExecutionStateChangedEventLogs(
        logs,
        SOURCE_CHAIN_SELECTOR_3,
        messages2[k].header.sequenceNumber,
        messages2[k].header.messageId,
        _hashMessage(messages2[k], ON_RAMP_ADDRESS_3),
        Internal.MessageExecutionState.SUCCESS,
        ""
      );
    }
  }

  function test_manuallyExecute_WithPartialMessages() public {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](3);

    for (uint64 i = 0; i < 3; ++i) {
      messages[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, i + 1);
    }

    messages[1].receiver = address(s_reverting_receiver);
    messages[1].header.messageId = _hashMessage(messages[1], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
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
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        OffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, bytes(""))
      )
    );

    _assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[2].header.sequenceNumber,
      messages[2].header.messageId,
      _hashMessage(messages[2], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_reverting_receiver.setRevert(false);

    // Only the 2nd message reverted
    Internal.Any2EVMRampMessage[] memory newMessages = new Internal.Any2EVMRampMessage[](1);
    newMessages[0] = messages[1];

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(newMessages);
    gasLimitOverrides[0][0].receiverExecutionGasLimit += 1;

    vm.recordLogs();
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, newMessages), gasLimitOverrides);
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_manuallyExecute_LowGasLimit() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].gasLimit = 1;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
    );
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(OffRamp.ReceiverError.selector, "")
    );

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](1);
    gasLimitOverrides[0][0].receiverExecutionGasLimit = 100_000;

    vm.expectEmit();
    emit ConformingReceiver.MessageReceived();

    vm.recordLogs();
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  // Reverts

  function test_RevertWhen_manuallyExecute_ForkedChain() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    uint256 chain1 = uint200(block.chainid);
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ForkedChain.selector, chain1, chain2));

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_RevertWhen_ManualExecGasLimitMismatchSingleReport() public {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](2);
    messages[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);

    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    // No overrides for report
    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, new OffRamp.GasLimitOverride[][](0));

    // No messages
    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1 message missing
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](1);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1 message in excess
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](3);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_RevertWhen_manuallyExecute_GasLimitMismatchMultipleReports() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, 1);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, new OffRamp.GasLimitOverride[][](0));

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, new OffRamp.GasLimitOverride[][](1));

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](2);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 2nd report empty
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](2);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1st report empty
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](0);
    gasLimitOverrides[1] = new OffRamp.GasLimitOverride[](1);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1st report oversized
    gasLimitOverrides[0] = new OffRamp.GasLimitOverride[](3);

    vm.expectRevert(OffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_RevertWhen_manuallyExecute_InvalidReceiverExecutionGasLimit() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0][0].receiverExecutionGasLimit--;

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.InvalidManualExecutionGasLimit.selector,
        SOURCE_CHAIN_SELECTOR_1,
        messages[0].header.messageId,
        gasLimitOverrides[0][0].receiverExecutionGasLimit
      )
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_RevertWhen_manuallyExecute_DestinationGasAmountCountMismatch() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 1000;
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](1);
    messages[0] = _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    // empty tokenGasOverride array provided
    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.ManualExecutionGasAmountCountMismatch.selector, messages[0].header.messageId, 1)
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);

    //trying with excesss elements tokenGasOverride array provided
    gasLimitOverrides[0][0].tokenGasOverrides = new uint32[](3);
    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.ManualExecutionGasAmountCountMismatch.selector, messages[0].header.messageId, 1)
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_RevertWhen_manuallyExecute_InvalidTokenGasOverride() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 1000;
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](1);
    messages[0] = _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    uint32[] memory tokenGasOverrides = new uint32[](2);
    tokenGasOverrides[0] = DEFAULT_TOKEN_DEST_GAS_OVERHEAD;
    tokenGasOverrides[1] = DEFAULT_TOKEN_DEST_GAS_OVERHEAD - 1; //invalid token gas override value
    gasLimitOverrides[0][0].tokenGasOverrides = tokenGasOverrides;

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.InvalidManualExecutionTokenGasOverride.selector,
        messages[0].header.messageId,
        1,
        DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        tokenGasOverrides[1]
      )
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_RevertWhen_manuallyExecute_FailedTx() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
    );

    s_reverting_receiver.setRevert(true);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.ExecutionError.selector,
        messages[0].header.messageId,
        abi.encodeWithSelector(
          OffRamp.ReceiverError.selector,
          abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, bytes(""))
        )
      )
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_manuallyExecute_ReentrancyFails() public {
    uint256 tokenAmount = 1e9;
    IERC20 tokenToAbuse = IERC20(s_destFeeToken);

    // This needs to be deployed before the source chain message is sent
    // because we need the address for the receiver.
    ReentrancyAbuserMultiRamp receiver = new ReentrancyAbuserMultiRamp(address(s_destRouter), s_offRamp);
    uint256 balancePre = tokenToAbuse.balanceOf(address(receiver));

    // For this test any message will be flagged as correct by the
    // commitStore. In a real scenario the abuser would have to actually
    // send the message that they want to replay.
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].tokenAmounts = new Internal.Any2EVMTokenTransfer[](1);
    messages[0].tokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[s_sourceFeeToken]),
      destTokenAddress: s_destTokenBySourceToken[s_sourceFeeToken],
      extraData: "",
      amount: tokenAmount,
      destGasAmount: MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS
    });

    messages[0].receiver = address(receiver);

    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    // sets the report to be repeated on the ReentrancyAbuser to be able to replay
    receiver.setPayload(report);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0][0].tokenGasOverrides = new uint32[](messages[0].tokenAmounts.length);

    // The first entry should be fine and triggers the second entry which is skipped. Due to the reentrancy
    // the second completes first, so we expect the skip event before the success event.
    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(
      messages[0].header.sourceChainSelector, messages[0].header.sequenceNumber
    );

    vm.recordLogs();
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    // Since the tx failed we don't release the tokens
    assertEq(tokenToAbuse.balanceOf(address(receiver)), balancePre + tokenAmount);
  }

  function test_RevertWhen_manuallyExecute_MultipleReportsWithSingleCursedLane() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](3);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](2);

    for (uint64 i = 0; i < 3; ++i) {
      messages1[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, i + 1);
      messages1[i].receiver = address(s_reverting_receiver);
      messages1[i].header.messageId = _hashMessage(messages1[i], ON_RAMP_ADDRESS_1);
    }

    for (uint64 i = 0; i < 2; ++i) {
      messages2[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, i + 1);
      messages2[i].receiver = address(s_reverting_receiver);
      messages2[i].header.messageId = _hashMessage(messages2[i], ON_RAMP_ADDRESS_3);
    }

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](2);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages1);
    gasLimitOverrides[1] = _getGasLimitsFromMessages(messages2);

    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_3, true);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.CursedByRMN.selector, SOURCE_CHAIN_SELECTOR_3));

    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_RevertWhen_manuallyExecute_SourceChainSelectorMismatch() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](1);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);
    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](2);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages1);
    gasLimitOverrides[1] = _getGasLimitsFromMessages(messages2);

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.SourceChainSelectorMismatch.selector, SOURCE_CHAIN_SELECTOR_3, SOURCE_CHAIN_SELECTOR_1
      )
    );
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }
}
