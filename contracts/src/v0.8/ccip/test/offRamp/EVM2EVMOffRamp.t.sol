// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICommitStore} from "../../interfaces/ICommitStore.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {CallWithExactGas} from "../../../shared/call/CallWithExactGas.sol";

import {GenericReceiver} from "../../../shared/test/testhelpers/GenericReceiver.sol";
import {AggregateRateLimiter} from "../../AggregateRateLimiter.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {OCR2Abstract} from "../../ocr/OCR2Abstract.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {ConformingReceiver} from "../helpers/receivers/ConformingReceiver.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MaybeRevertMessageReceiverNo165} from "../helpers/receivers/MaybeRevertMessageReceiverNo165.sol";
import {ReentrancyAbuser} from "../helpers/receivers/ReentrancyAbuser.sol";
import {OCR2BaseNoChecks} from "../ocr/OCR2BaseNoChecks.t.sol";
import {EVM2EVMOffRampSetup} from "./EVM2EVMOffRampSetup.t.sol";
import {stdError} from "forge-std/Test.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMOffRamp_constructor is EVM2EVMOffRampSetup {
  function test_Constructor_Success() public {
    EVM2EVMOffRamp.StaticConfig memory staticConfig = EVM2EVMOffRamp.StaticConfig({
      commitStore: address(s_mockCommitStore),
      chainSelector: DEST_CHAIN_SELECTOR,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      onRamp: ON_RAMP_ADDRESS,
      prevOffRamp: address(0),
      rmnProxy: address(s_mockRMN),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig =
      generateDynamicOffRampConfig(address(s_destRouter), address(s_feeQuoter));

    s_offRamp = new EVM2EVMOffRampHelper(staticConfig, _getInboundRateLimiterConfig());

    s_offRamp.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );

    // Static config
    EVM2EVMOffRamp.StaticConfig memory gotStaticConfig = s_offRamp.getStaticConfig();
    assertEq(staticConfig.commitStore, gotStaticConfig.commitStore);
    assertEq(staticConfig.sourceChainSelector, gotStaticConfig.sourceChainSelector);
    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.onRamp, gotStaticConfig.onRamp);
    assertEq(staticConfig.prevOffRamp, gotStaticConfig.prevOffRamp);
    assertEq(staticConfig.tokenAdminRegistry, gotStaticConfig.tokenAdminRegistry);

    // Dynamic config
    EVM2EVMOffRamp.DynamicConfig memory gotDynamicConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, gotDynamicConfig);

    (uint32 configCount, uint32 blockNumber,) = s_offRamp.latestConfigDetails();
    assertEq(1, configCount);
    assertEq(block.number, blockNumber);

    // OffRamp initial values
    assertEq("EVM2EVMOffRamp 1.5.0", s_offRamp.typeAndVersion());
    assertEq(OWNER, s_offRamp.owner());
  }

  // Revert
  function test_ZeroOnRampAddress_Revert() public {
    vm.expectRevert(EVM2EVMOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(s_mockCommitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ZERO_ADDRESS,
        prevOffRamp: address(0),
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      RateLimiter.Config({isEnabled: true, rate: 1e20, capacity: 1e20})
    );
  }

  function test_CommitStoreAlreadyInUse_Revert() public {
    s_mockCommitStore.setExpectedNextSequenceNumber(2);

    vm.expectRevert(EVM2EVMOffRamp.CommitStoreAlreadyInUse.selector);

    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(s_mockCommitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        prevOffRamp: address(0),
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _getInboundRateLimiterConfig()
    );
  }
}

contract EVM2EVMOffRamp_setDynamicConfig is EVM2EVMOffRampSetup {
  function test_SetDynamicConfig_Success() public {
    EVM2EVMOffRamp.StaticConfig memory staticConfig = s_offRamp.getStaticConfig();
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(USER_3, address(s_feeQuoter));
    bytes memory onchainConfig = abi.encode(dynamicConfig);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ConfigSet(staticConfig, dynamicConfig);

    vm.expectEmit();
    uint32 configCount = 1;
    emit OCR2Abstract.ConfigSet(
      uint32(block.number),
      getBasicConfigDigest(address(s_offRamp), s_f, configCount, onchainConfig),
      configCount + 1,
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      onchainConfig,
      s_offchainConfigVersion,
      abi.encode("")
    );

    s_offRamp.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, onchainConfig, s_offchainConfigVersion, abi.encode("")
    );

    EVM2EVMOffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  function test_NonOwner_Revert() public {
    vm.startPrank(STRANGER);
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(USER_3, address(s_feeQuoter));

    vm.expectRevert("Only callable by owner");

    s_offRamp.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }

  function test_RouterZeroAddress_Revert() public {
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(ZERO_ADDRESS, ZERO_ADDRESS);

    vm.expectRevert(EVM2EVMOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }
}

contract EVM2EVMOffRamp_metadataHash is EVM2EVMOffRampSetup {
  function test_MetadataHash_Success() public view {
    bytes32 h = s_offRamp.metadataHash();
    assertEq(
      h,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );
  }
}

contract EVM2EVMOffRamp_ccipReceive is EVM2EVMOffRampSetup {
  // Reverts

  function test_Reverts() public {
    Client.Any2EVMMessage memory message = _convertToGeneralMessage(_generateAny2EVMMessageNoTokens(1));
    vm.expectRevert();
    s_offRamp.ccipReceive(message);
  }
}

contract EVM2EVMOffRamp_execute is EVM2EVMOffRampSetup {
  error PausedError();

  function _generateMsgWithoutTokens(
    uint256 gasLimit,
    bytes memory messageData
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    message.gasLimit = gasLimit;
    message.data = messageData;
    message.messageId = Internal._hash(
      message,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );
    return message;
  }

  function test_Fuzz_trialExecuteWithoutTokens_Success(bytes4 funcSelector, bytes memory messageData) public {
    vm.assume(
      funcSelector != GenericReceiver.setRevert.selector && funcSelector != GenericReceiver.setErr.selector
        && funcSelector != 0x5100fc21 && funcSelector != 0x00000000 // s_toRevert(), which is public and therefore has a function selector
    );

    // Convert bytes4 into bytes memory to use in the message
    Internal.EVM2EVMMessage memory message = _generateMsgWithoutTokens(GAS_LIMIT, messageData);

    // Convert an Internal.EVM2EVMMessage into a Client.Any2EVMMessage digestable by the client
    Client.Any2EVMMessage memory receivedMessage = _convertToGeneralMessage(message);
    bytes memory expectedCallData =
      abi.encodeWithSelector(MaybeRevertMessageReceiver.ccipReceive.selector, receivedMessage);

    vm.expectCall(address(s_receiver), expectedCallData);
    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);
  }

  function test_Fuzz_trialExecuteWithTokens_Success(uint16 tokenAmount, bytes calldata messageData) public {
    vm.assume(tokenAmount != 0);

    uint256[] memory amounts = new uint256[](2);
    amounts[0] = uint256(tokenAmount);
    amounts[1] = uint256(tokenAmount);

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    message.data = messageData;

    IERC20 dstToken0 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken0.balanceOf(message.receiver);

    vm.expectCall(address(dstToken0), abi.encodeWithSelector(IERC20.transfer.selector, address(s_receiver), amounts[0]));

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // Check that the tokens were transferred
    assertEq(startingBalance + amounts[0], dstToken0.balanceOf(message.receiver));
  }

  function test_Fuzz_getSenderNonce(uint8 trialExecutions) public {
    vm.assume(trialExecutions > 1);

    Internal.EVM2EVMMessage[] memory messages;

    if (trialExecutions == 1) {
      messages = new Internal.EVM2EVMMessage[](1);
      messages[0] = _generateAny2EVMMessageNoTokens(0);
    } else {
      messages = _generateSingleBasicMessage();
    }

    // Fuzz the number of calls from the sender to ensure that getSenderNonce works
    for (uint256 i = 1; i < trialExecutions; ++i) {
      s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

      messages[0].nonce++;
      messages[0].sequenceNumber++;
      messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    }

    messages[0].nonce = 0;
    messages[0].sequenceNumber = 0;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    uint64 nonceBefore = s_offRamp.getSenderNonce(messages[0].sender);
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(s_offRamp.getSenderNonce(messages[0].sender), nonceBefore, "sender nonce is not as expected");
  }

  function test_Fuzz_getSenderNonceWithPrevOffRamp_Success(uint8 trialExecutions) public {
    vm.assume(trialExecutions > 1);
    // Fuzz a random nonce for getSenderNonce
    test_Fuzz_getSenderNonce(trialExecutions);

    address prevOffRamp = address(s_offRamp);
    deployOffRamp(s_mockCommitStore, s_destRouter, prevOffRamp);

    // Make sure the off-ramp address has changed by querying the static config
    assertNotEq(address(s_offRamp), prevOffRamp);
    EVM2EVMOffRamp.StaticConfig memory staticConfig = s_offRamp.getStaticConfig();
    assertEq(staticConfig.prevOffRamp, prevOffRamp, "Previous offRamp does not match expected address");

    // Since i_prevOffRamp != address(0) and senderNonce == 0, there should be a call to the previous offRamp
    vm.expectCall(prevOffRamp, abi.encodeWithSelector(s_offRamp.getSenderNonce.selector, OWNER));
    uint256 currentSenderNonce = s_offRamp.getSenderNonce(OWNER);
    assertNotEq(currentSenderNonce, 0, "Sender nonce should not be zero");
    assertEq(currentSenderNonce, trialExecutions - 1, "Sender Nonce does not match expected trial executions");

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    currentSenderNonce = s_offRamp.getSenderNonce(OWNER);
    assertEq(currentSenderNonce, trialExecutions - 1, "Sender Nonce on new offramp does not match expected executions");
  }

  function test_SingleMessageNoTokens_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    messages[0].nonce++;
    messages[0].sequenceNumber++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    uint64 nonceBefore = s_offRamp.getSenderNonce(messages[0].sender);
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertGt(s_offRamp.getSenderNonce(messages[0].sender), nonceBefore);
  }

  function test_SingleMessageNoTokensUnordered_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].nonce = 0;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    // Nonce never increments on unordered messages.
    uint64 nonceBefore = s_offRamp.getSenderNonce(messages[0].sender);
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(
      s_offRamp.getSenderNonce(messages[0].sender), nonceBefore, "nonce must remain unchanged on unordered messages"
    );

    messages[0].sequenceNumber++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    // Nonce never increments on unordered messages.
    nonceBefore = s_offRamp.getSenderNonce(messages[0].sender);
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(
      s_offRamp.getSenderNonce(messages[0].sender), nonceBefore, "nonce must remain unchanged on unordered messages"
    );
  }

  function test_ReceiverError_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1)
      )
    );
    // Nonce should increment on non-strict
    assertEq(uint64(0), s_offRamp.getSenderNonce(address(OWNER)));
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(uint64(1), s_offRamp.getSenderNonce(address(OWNER)));
  }

  function test_StrictUntouchedToSuccess_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    messages[0].strict = true;
    messages[0].receiver = address(s_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    // Nonce should increment on a strict untouched -> success.
    assertEq(uint64(0), s_offRamp.getSenderNonce(address(OWNER)));
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(uint64(1), s_offRamp.getSenderNonce(address(OWNER)));
  }

  function test_SkippedIncorrectNonce_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    messages[0].nonce++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.SkippedIncorrectNonce(messages[0].nonce, messages[0].sender);

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_SkippedIncorrectNonceStillExecutes_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();

    messages[1].nonce++;
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    vm.expectEmit();
    emit EVM2EVMOffRamp.SkippedIncorrectNonce(messages[1].nonce, messages[1].sender);

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test__execute_SkippedAlreadyExecutedMessage_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    vm.expectEmit();
    emit EVM2EVMOffRamp.SkippedAlreadyExecutedMessage(messages[0].sequenceNumber);

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test__execute_SkippedAlreadyExecutedMessageUnordered_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].nonce = 0;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    vm.expectEmit();
    emit EVM2EVMOffRamp.SkippedAlreadyExecutedMessage(messages[0].sequenceNumber);

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  // Send a message to a contract that does not implement the CCIPReceiver interface
  // This should execute successfully.
  function test_SingleMessageToNonCCIPReceiver_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    messages[0].receiver = address(newReceiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_SingleMessagesNoTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.resumeGasMetering();
    s_offRamp.execute(report, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_TwoMessagesWithTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[1].sequenceNumber, messages[1].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.resumeGasMetering();
    s_offRamp.execute(report, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_TwoMessagesWithTokensAndGE_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[1].sequenceNumber, messages[1].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    assertEq(uint64(0), s_offRamp.getSenderNonce(OWNER));
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    assertEq(uint64(2), s_offRamp.getSenderNonce(OWNER));
  }

  function test_Fuzz_InterleavingOrderedAndUnorderedMessages_Success(bool[7] memory orderings) public {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](orderings.length);
    // number of tokens needs to be capped otherwise we hit UnsupportedNumberOfTokens.
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](3);
    for (uint256 i = 0; i < 3; ++i) {
      tokenAmounts[i].token = s_sourceTokens[i % s_sourceTokens.length];
      tokenAmounts[i].amount = 1e18;
    }
    uint64 expectedNonce = 0;
    for (uint256 i = 0; i < orderings.length; ++i) {
      messages[i] = _generateAny2EVMMessage(uint64(i + 1), tokenAmounts, !orderings[i]);
      if (orderings[i]) {
        messages[i].nonce = ++expectedNonce;
      }
      messages[i].messageId = Internal._hash(messages[i], s_offRamp.metadataHash());

      vm.expectEmit();
      emit EVM2EVMOffRamp.ExecutionStateChanged(
        messages[i].sequenceNumber, messages[i].messageId, Internal.MessageExecutionState.SUCCESS, ""
      );
    }

    uint64 nonceBefore = s_offRamp.getSenderNonce(OWNER);
    assertEq(uint64(0), nonceBefore, "nonce before exec should be 0");
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    // all executions should succeed.
    for (uint256 i = 0; i < orderings.length; ++i) {
      assertEq(
        uint256(s_offRamp.getExecutionState(messages[i].sequenceNumber)),
        uint256(Internal.MessageExecutionState.SUCCESS)
      );
    }
    assertEq(nonceBefore + expectedNonce, s_offRamp.getSenderNonce(OWNER));
  }

  function test_InvalidSourcePoolAddress_Success() public {
    address fakePoolAddress = address(0x0000000000333333);

    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    messages[0].sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(fakePoolAddress),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[messages[0].tokenAmounts[0].token]),
        extraData: "",
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMOffRamp.TokenHandlingError.selector,
        abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, abi.encode(fakePoolAddress))
      )
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_execute_RouterYULCall_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    // gas limit too high, Router's external call should revert
    messages[0].gasLimit = 1e36;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector)
    );

    s_offRamp.execute(executionReport, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_RetryFailedMessageWithoutManualExecution_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1)
      )
    );
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    // The second time should skip the msg
    vm.expectEmit();
    emit EVM2EVMOffRamp.AlreadyAttempted(messages[0].sequenceNumber);

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  // Reverts

  function test_InvalidMessageId_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].nonce++;
    // MessageID no longer matches hash.
    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);
    vm.expectRevert(EVM2EVMOffRamp.InvalidMessageId.selector);
    s_offRamp.execute(executionReport, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_Paused_Revert() public {
    s_mockCommitStore.pause();
    vm.expectRevert(PausedError.selector);
    s_offRamp.execute(
      _generateReportFromMessages(_generateMessagesWithTokens()), new EVM2EVMOffRamp.GasLimitOverride[](0)
    );
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(EVM2EVMOffRamp.CursedByRMN.selector);
    s_offRamp.execute(
      _generateReportFromMessages(_generateMessagesWithTokens()), new EVM2EVMOffRamp.GasLimitOverride[](0)
    );
    // Uncurse should succeed
    s_mockRMN.setGlobalCursed(false);
    s_offRamp.execute(
      _generateReportFromMessages(_generateMessagesWithTokens()), new EVM2EVMOffRamp.GasLimitOverride[](0)
    );
  }

  function test_UnexpectedTokenData_Revert() public {
    Internal.ExecutionReport memory report = _generateReportFromMessages(_generateSingleBasicMessage());
    report.offchainTokenData = new bytes[][](report.messages.length + 1);

    vm.expectRevert(EVM2EVMOffRamp.UnexpectedTokenData.selector);

    s_offRamp.execute(report, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_EmptyReport_Revert() public {
    vm.expectRevert(EVM2EVMOffRamp.EmptyReport.selector);
    s_offRamp.execute(
      Internal.ExecutionReport({
        proofs: new bytes32[](0),
        proofFlagBits: 0,
        messages: new Internal.EVM2EVMMessage[](0),
        offchainTokenData: new bytes[][](0)
      }),
      new EVM2EVMOffRamp.GasLimitOverride[](0)
    );
  }

  function test_RootNotCommitted_Revert() public {
    vm.mockCall(address(s_mockCommitStore), abi.encodeWithSelector(ICommitStore.verify.selector), abi.encode(0));
    vm.expectRevert(EVM2EVMOffRamp.RootNotCommitted.selector);

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    vm.clearMockedCalls();
  }

  function test_ManualExecutionNotYetEnabled_Revert() public {
    vm.mockCall(
      address(s_mockCommitStore), abi.encodeWithSelector(ICommitStore.verify.selector), abi.encode(BLOCK_TIME)
    );
    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionNotYetEnabled.selector);

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    vm.clearMockedCalls();
  }

  function test_InvalidSourceChain_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].sourceChainSelector = SOURCE_CHAIN_SELECTOR + 1;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.InvalidSourceChain.selector, SOURCE_CHAIN_SELECTOR + 1));
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_UnsupportedNumberOfTokens_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    Client.EVMTokenAmount[] memory newTokens = new Client.EVMTokenAmount[](MAX_TOKENS_LENGTH + 1);
    messages[0].tokenAmounts = newTokens;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.UnsupportedNumberOfTokens.selector, messages[0].sequenceNumber)
    );
    s_offRamp.execute(report, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_TokenDataMismatch_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    report.offchainTokenData[0] = new bytes[](messages[0].tokenAmounts.length + 1);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenDataMismatch.selector, messages[0].sequenceNumber));
    s_offRamp.execute(report, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_MessageTooLarge_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].data = new bytes(MAX_DATA_SIZE + 1);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.MessageTooLarge.selector, MAX_DATA_SIZE, messages[0].data.length)
    );
    s_offRamp.execute(executionReport, new EVM2EVMOffRamp.GasLimitOverride[](0));
  }
}

contract EVM2EVMOffRamp_execute_upgrade is EVM2EVMOffRampSetup {
  EVM2EVMOffRampHelper internal s_prevOffRamp;

  function setUp() public virtual override {
    super.setUp();

    s_prevOffRamp = s_offRamp;

    deployOffRamp(s_mockCommitStore, s_destRouter, address(s_prevOffRamp));
  }

  function test_V2_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
  }

  function test_V2SenderNoncesReadsPreviousRamp_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    uint64 startNonce = s_offRamp.getSenderNonce(messages[0].sender);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOffRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

      messages[0].nonce++;
      messages[0].sequenceNumber++;
      messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

      assertEq(startNonce + i, s_offRamp.getSenderNonce(messages[0].sender));
    }
  }

  function test_V2NonceStartsAtV1Nonce_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    uint64 startNonce = s_offRamp.getSenderNonce(messages[0].sender);

    s_prevOffRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    assertEq(startNonce + 1, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(startNonce + 2, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce++;
    messages[0].sequenceNumber++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(startNonce + 3, s_offRamp.getSenderNonce(messages[0].sender));
  }

  function test_V2NonceNewSenderStartsAtZero_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_prevOffRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    address newSender = address(1234567);
    messages[0].sender = newSender;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    // new sender nonce in new offramp should go from 0 -> 1
    assertEq(s_offRamp.getSenderNonce(newSender), 0);
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(s_offRamp.getSenderNonce(newSender), 1);
  }

  function test_V2OffRampNonceSkipsIfMsgInFlight_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    address newSender = address(1234567);
    messages[0].sender = newSender;
    messages[0].nonce = 2;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    uint64 startNonce = s_offRamp.getSenderNonce(messages[0].sender);

    // new offramp sees msg nonce higher than senderNonce
    // it waits for previous offramp to execute
    vm.expectEmit();
    emit EVM2EVMOffRamp.SkippedSenderWithPreviousRampMessageInflight(messages[0].nonce, newSender);
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(startNonce, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce = 1;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    // previous offramp executes msg and increases nonce
    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    s_prevOffRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(startNonce + 1, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce = 2;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    // new offramp is able to execute
    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));
    assertEq(startNonce + 2, s_offRamp.getSenderNonce(messages[0].sender));
  }
}

contract EVM2EVMOffRamp_executeSingleMessage is EVM2EVMOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    vm.startPrank(address(s_offRamp));
  }

  function test_executeSingleMessage_NoTokens_Success() public {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_WithTokens_Success() public {
    Internal.EVM2EVMMessage memory message = _generateMessagesWithTokens()[0];
    bytes[] memory offchainTokenData = new bytes[](message.tokenAmounts.length);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(message.sourceTokenData[0], (Internal.SourceTokenData));

    vm.expectCall(
      s_destPoolByToken[s_destTokens[0]],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(message.sender),
          receiver: message.receiver,
          amount: message.tokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[message.tokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR,
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: ""
        })
      )
    );

    s_offRamp.executeSingleMessage(message, offchainTokenData, new uint32[](0));
  }

  function test_executeSingleMessage_ZeroGasZeroData_Success() public {
    uint256 gasLimit = 0;
    Internal.EVM2EVMMessage memory message = _generateMsgWithoutTokens(gasLimit);
    Client.Any2EVMMessage memory receiverMsg = _convertToGeneralMessage(message);

    // expect 0 calls to be made as no gas is provided
    vm.expectCall(
      address(s_destRouter),
      abi.encodeCall(Router.routeMessage, (receiverMsg, Internal.GAS_FOR_CALL_EXACT_CHECK, gasLimit, message.receiver)),
      0
    );

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    // Ensure we encoded it properly, and didn't simply expect the wrong call
    gasLimit = 200_000;
    message = _generateMsgWithoutTokens(gasLimit);
    receiverMsg = _convertToGeneralMessage(message);

    vm.expectCall(
      address(s_destRouter),
      abi.encodeCall(Router.routeMessage, (receiverMsg, Internal.GAS_FOR_CALL_EXACT_CHECK, gasLimit, message.receiver)),
      1
    );

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function _generateMsgWithoutTokens(uint256 gasLimit) internal view returns (Internal.EVM2EVMMessage memory) {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    message.gasLimit = gasLimit;
    message.data = "";
    message.messageId = Internal._hash(
      message,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );
    return message;
  }

  function test_NonContract_Success() public {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_NonContractWithTokens_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;
    vm.expectEmit();
    emit TokenPool.Released(address(s_offRamp), STRANGER, amounts[0]);
    vm.expectEmit();
    emit TokenPool.Minted(address(s_offRamp), STRANGER, amounts[1]);
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  // Reverts

  function test_TokenHandlingError_Revert() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    bytes memory errorMessage = "Random token pool issue";

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, errorMessage));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_ZeroGasDONExecution_Revert() public {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    message.gasLimit = 0;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.ReceiverError.selector, ""));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_MessageSender_Revert() public {
    vm.stopPrank();
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    vm.expectRevert(EVM2EVMOffRamp.CanOnlySelfCall.selector);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }
}

contract EVM2EVMOffRamp__report is EVM2EVMOffRampSetup {
  // Asserts that execute completes
  function test_Report_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    s_offRamp.report(abi.encode(report));
  }
}

contract EVM2EVMOffRamp_manuallyExecute is EVM2EVMOffRampSetup {
  function test_ManualExec_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    s_offRamp.manuallyExecute(
      _generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length)
    );
  }

  function test_ManualExecWithSourceTokens_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessageWithTokens();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_ManualExecWithMultipleMessagesAndSourceTokens_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    messages[1].receiver = address(s_reverting_receiver);
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[1].sequenceNumber, messages[1].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_manuallyExecute_DoesNotRevertIfUntouched_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    assertEq(messages[0].nonce - 1, s_offRamp.getSenderNonce(messages[0].sender));

    s_reverting_receiver.setRevert(true);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, "")
      )
    );

    s_offRamp.manuallyExecute(
      _generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length)
    );

    assertEq(messages[0].nonce, s_offRamp.getSenderNonce(messages[0].sender));
  }

  function test_manuallyExecute_WithGasOverride_Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length));

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0].receiverExecutionGasLimit += 1;

    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_manuallyExecute_WithInvalidSourceTokenDataCount_Revert() public {
    uint256 messageIndex = 0;

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessageWithTokens();
    messages[messageIndex].receiver = address(s_reverting_receiver);
    messages[messageIndex].messageId = Internal._hash(messages[messageIndex], s_offRamp.metadataHash());

    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);

    messages[messageIndex].sourceTokenData = new bytes[](0);

    vm.expectRevert(stdError.indexOOBError);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_manuallyExecute_WithInvalidReceiverExecutionGasOverride_Revert() public {
    uint256 messageIndex = 0;

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[messageIndex].receiver = address(s_reverting_receiver);
    messages[messageIndex].messageId = Internal._hash(messages[messageIndex], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length));

    s_reverting_receiver.setRevert(false);
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[messageIndex].receiverExecutionGasLimit -= 1;

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.InvalidManualExecutionGasLimit.selector,
        messages[messageIndex].messageId,
        messages[messageIndex].gasLimit,
        gasLimitOverrides[messageIndex].receiverExecutionGasLimit
      )
    );

    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_manuallyExecute_LowGasLimitManualExec_Success() public {
    uint256 messageIndex = 0;

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[messageIndex].gasLimit = 1;
    messages[messageIndex].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[messageIndex].messageId = Internal._hash(messages[messageIndex], s_offRamp.metadataHash());

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[messageIndex].sequenceNumber,
      messages[messageIndex].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(EVM2EVMOffRamp.ReceiverError.selector, "")
    );
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = new EVM2EVMOffRamp.GasLimitOverride[](1);
    gasLimitOverrides[messageIndex].receiverExecutionGasLimit = 100_000;

    vm.expectEmit();
    emit MaybeRevertMessageReceiver.MessageReceived();

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[messageIndex].sequenceNumber,
      messages[messageIndex].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_ReentrancyManualExecuteFails_Success() public {
    uint256 tokenAmount = 1e9;
    IERC20 tokenToAbuse = IERC20(s_destFeeToken);

    // This needs to be deployed before the source chain message is sent
    // because we need the address for the receiver.
    ReentrancyAbuser receiver = new ReentrancyAbuser(address(s_destRouter), s_offRamp);
    uint256 balancePre = tokenToAbuse.balanceOf(address(receiver));

    // For this test any message will be flagged as correct by the
    // commitStore. In a real scenario the abuser would have to actually
    // send the message that they want to replay.
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    messages[0].tokenAmounts = new Client.EVMTokenAmount[](1);
    messages[0].tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: tokenAmount});
    messages[0].receiver = address(receiver);
    messages[0].sourceTokenData = new bytes[](1);
    messages[0].sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[s_sourceFeeToken]),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[s_sourceFeeToken]),
        extraData: "",
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    // sets the report to be repeated on the ReentrancyAbuser to be able to replay
    receiver.setPayload(report);

    // The first entry should be fine and triggers the second entry which is skipped. Due to the reentrancy
    // the second completes first, so we expect the skip event before the success event.
    vm.expectEmit();
    emit EVM2EVMOffRamp.SkippedAlreadyExecutedMessage(messages[0].sequenceNumber);

    vm.expectEmit();
    emit EVM2EVMOffRamp.ExecutionStateChanged(
      messages[0].sequenceNumber, messages[0].messageId, Internal.MessageExecutionState.SUCCESS, ""
    );

    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimits = _getGasLimitsFromMessages(messages);
    s_offRamp.manuallyExecute(report, gasLimits);

    // Assert that they only got the tokens once, not twice
    assertEq(tokenToAbuse.balanceOf(address(receiver)), balancePre + tokenAmount);
  }

  function test_manuallyExecute_InvalidTokenGasOverride_Revert() public {
    uint256 failingMessageIndex = 0;
    uint256 failingTokenIndex = 0;

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessageWithTokens();
    messages[failingMessageIndex].receiver = address(s_reverting_receiver);
    messages[failingMessageIndex].messageId = Internal._hash(messages[failingMessageIndex], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length));

    s_reverting_receiver.setRevert(false);
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[failingMessageIndex].tokenGasOverrides[failingTokenIndex] -= 2;

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.InvalidTokenGasOverride.selector,
        messages[failingMessageIndex].messageId,
        failingTokenIndex,
        DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        gasLimitOverrides[failingMessageIndex].tokenGasOverrides[failingTokenIndex]
      )
    );

    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_manuallyExecute_DestinationGasAmountCountMismatch_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessageWithTokens();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length));

    s_reverting_receiver.setRevert(false);
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0].tokenGasOverrides = new uint32[](0);

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.DestinationGasAmountCountMismatch.selector, messages[0].messageId, 1)
    );

    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function test_ManualExecForkedChain_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);
    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(OCR2BaseNoChecks.ForkedChain.selector, chain1, chain2));

    s_offRamp.manuallyExecute(report, _getGasLimitsFromMessages(messages));
  }

  function test_ManualExecGasLimitMismatch_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(
      _generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length - 1)
    );

    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(
      _generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](messages.length + 1)
    );
  }

  function test_ManualExecInvalidGasLimit_Revert() public {
    uint256 messageIndex = 0;

    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimits = _getGasLimitsFromMessages(messages);
    gasLimits[messageIndex].receiverExecutionGasLimit -= 1;

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.InvalidManualExecutionGasLimit.selector,
        messages[messageIndex].messageId,
        messages[messageIndex].gasLimit,
        gasLimits[messageIndex].receiverExecutionGasLimit
      )
    );
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimits);
  }

  function test_ManualExecFailedTx_Revert() public {
    Internal.EVM2EVMMessage[] memory messages = _generateSingleBasicMessage();

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    s_offRamp.execute(_generateReportFromMessages(messages), new EVM2EVMOffRamp.GasLimitOverride[](0));

    s_reverting_receiver.setRevert(true);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ExecutionError.selector,
        abi.encodeWithSelector(
          EVM2EVMOffRamp.ReceiverError.selector,
          abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, bytes(""))
        )
      )
    );
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
  }
}

contract EVM2EVMOffRamp_getExecutionState is EVM2EVMOffRampSetup {
  mapping(uint64 seqNum => Internal.MessageExecutionState state) internal s_differentialExecutionState;

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function test_Fuzz_Differential_Success(uint16[500] memory seqNums, uint8[500] memory values) public {
    for (uint256 i = 0; i < seqNums.length; ++i) {
      // Only use the first three slots. This makes sure existing slots get overwritten
      // as the tests uses 500 sequence numbers.
      uint16 seqNum = seqNums[i] % 386;
      Internal.MessageExecutionState state = Internal.MessageExecutionState(values[i] % 4);
      s_differentialExecutionState[seqNum] = state;
      s_offRamp.setExecutionStateHelper(seqNum, state);
      assertEq(uint256(state), uint256(s_offRamp.getExecutionState(seqNum)));
    }

    for (uint256 i = 0; i < seqNums.length; ++i) {
      uint16 seqNum = seqNums[i] % 386;
      Internal.MessageExecutionState expectedState = s_differentialExecutionState[seqNum];
      assertEq(uint256(expectedState), uint256(s_offRamp.getExecutionState(seqNum)));
    }
  }

  function test_GetExecutionState_Success() public {
    s_offRamp.setExecutionStateHelper(0, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(0), 3);

    s_offRamp.setExecutionStateHelper(1, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(0), 3 + (3 << 2));

    s_offRamp.setExecutionStateHelper(1, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(0), 3 + (1 << 2));

    s_offRamp.setExecutionStateHelper(2, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(0), 3 + (1 << 2) + (3 << 4));

    s_offRamp.setExecutionStateHelper(127, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(0), 3 + (1 << 2) + (3 << 4) + (1 << 254));

    s_offRamp.setExecutionStateHelper(128, Internal.MessageExecutionState.SUCCESS);
    assertEq(s_offRamp.getExecutionStateBitMap(0), 3 + (1 << 2) + (3 << 4) + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(1), 2);

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(0)));
    assertEq(uint256(Internal.MessageExecutionState.IN_PROGRESS), uint256(s_offRamp.getExecutionState(1)));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(2)));
    assertEq(uint256(Internal.MessageExecutionState.IN_PROGRESS), uint256(s_offRamp.getExecutionState(127)));
    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(s_offRamp.getExecutionState(128)));
  }

  function test_FillExecutionState_Success() public {
    for (uint64 i = 0; i < 384; ++i) {
      s_offRamp.setExecutionStateHelper(i, Internal.MessageExecutionState.FAILURE);
    }

    for (uint64 i = 0; i < 384; ++i) {
      assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(i)));
    }

    for (uint64 i = 0; i < 3; ++i) {
      assertEq(type(uint256).max, s_offRamp.getExecutionStateBitMap(i));
    }

    for (uint64 i = 0; i < 384; ++i) {
      s_offRamp.setExecutionStateHelper(i, Internal.MessageExecutionState.IN_PROGRESS);
    }

    for (uint64 i = 0; i < 384; ++i) {
      assertEq(uint256(Internal.MessageExecutionState.IN_PROGRESS), uint256(s_offRamp.getExecutionState(i)));
    }

    for (uint64 i = 0; i < 3; ++i) {
      // 0x555... == 0b101010101010.....
      assertEq(0x5555555555555555555555555555555555555555555555555555555555555555, s_offRamp.getExecutionStateBitMap(i));
    }
  }
}

contract EVM2EVMOffRamp__trialExecute is EVM2EVMOffRampSetup {
  function test_trialExecute_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    IERC20 dstToken0 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken0.balanceOf(message.receiver);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // Check that the tokens were transferred
    assertEq(startingBalance + amounts[0], dstToken0.balanceOf(message.receiver));
  }

  function test_TokenHandlingErrorIsCaught_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    IERC20 dstToken0 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken0.balanceOf(OWNER);

    bytes memory errorMessage = "Random token pool issue";

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint32(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, errorMessage), err);

    // Expect the balance to remain the same
    assertEq(startingBalance, dstToken0.balanceOf(OWNER));
  }

  function test_RateLimitError_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    bytes memory errorMessage = abi.encodeWithSelector(RateLimiter.BucketOverfilled.selector);

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, errorMessage), err);
  }

  function test_TokenPoolIsNotAContract_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 10000;
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);

    // Happy path, pool is correct
    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // address 0 has no contract
    assertEq(address(0).code.length, 0);
    message.sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(address(0)),
        destTokenAddress: abi.encode(address(0)),
        extraData: "",
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    message.messageId = Internal._hash(
      message,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );

    // Unhappy path, no revert but marked as failed.
    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, abi.encode(address(0))), err);

    address notAContract = makeAddr("not_a_contract");

    message.sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(address(0)),
        destTokenAddress: abi.encode(notAContract),
        extraData: "",
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    message.messageId = Internal._hash(
      message,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );

    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMOffRamp.NotACompatiblePool.selector, address(0)), err);
  }
}

contract EVM2EVMOffRamp__releaseOrMintToken is EVM2EVMOffRampSetup {
  function test__releaseOrMintToken_Success() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = "";

    IERC20 dstToken1 = IERC20(s_destTokenBySourceToken[token]);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(s_destTokenBySourceToken[token]),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.expectCall(
      s_destPoolBySourceToken[token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: originalSender,
          receiver: OWNER,
          amount: amount,
          localToken: s_destTokenBySourceToken[token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR,
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData
        })
      )
    );

    s_offRamp.releaseOrMintToken(amount, originalSender, OWNER, sourceTokenData, offchainTokenData);

    assertEq(startingBalance + amount, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintToken_InvalidDataLength_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(s_destTokenBySourceToken[token]),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // Mock the call so returns 2 slots of data
    vm.mockCall(
      s_destTokenBySourceToken[token], abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER), abi.encode(0, 0)
    );

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.InvalidDataLength.selector, Internal.MAX_BALANCE_OF_RET_BYTES, 64)
    );

    s_offRamp.releaseOrMintToken(amount, abi.encode(OWNER), OWNER, sourceTokenData, "");
  }

  function test_releaseOrMintToken_TokenHandlingError_BalanceOf_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(s_destTokenBySourceToken[token]),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    bytes memory revertData = "failed to balanceOf";

    // Mock the call so returns 2 slots of data
    vm.mockCallRevert(
      s_destTokenBySourceToken[token], abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER), revertData
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, revertData));

    s_offRamp.releaseOrMintToken(amount, abi.encode(OWNER), OWNER, sourceTokenData, "");
  }

  function test_releaseOrMintToken_ReleaseOrMintBalanceMismatch_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    uint256 mockedStaticBalance = 50000;

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(s_destTokenBySourceToken[token]),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.mockCall(
      s_destTokenBySourceToken[token],
      abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER),
      abi.encode(mockedStaticBalance)
    );

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ReleaseOrMintBalanceMismatch.selector, amount, mockedStaticBalance, mockedStaticBalance
      )
    );

    s_offRamp.releaseOrMintToken(amount, abi.encode(OWNER), OWNER, sourceTokenData, "");
  }

  function test_releaseOrMintToken_skip_ReleaseOrMintBalanceMismatch_if_pool_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    uint256 mockedStaticBalance = 50000;

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(s_destTokenBySourceToken[token]),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // This should make the call fail if it does not skip the check
    vm.mockCall(
      s_destTokenBySourceToken[token],
      abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER),
      abi.encode(mockedStaticBalance)
    );

    s_offRamp.releaseOrMintToken(amount, abi.encode(OWNER), s_destPoolBySourceToken[token], sourceTokenData, "");
  }

  function test__releaseOrMintToken_NotACompatiblePool_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    vm.label(destToken, "destToken");
    bytes memory originalSender = abi.encode(OWNER);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(destToken),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // Address(0) should always revert
    address returnedPool = address(0);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintToken(amount, originalSender, OWNER, sourceTokenData, "");

    // A contract that doesn't support the interface should also revert
    returnedPool = address(s_offRamp);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintToken(amount, originalSender, OWNER, sourceTokenData, "");
  }

  function test__releaseOrMintToken_TokenHandlingError_transfer_Revert() public {
    address receiver = makeAddr("receiver");
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(destToken),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    bytes memory revertData = "call reverted :o";

    vm.mockCallRevert(destToken, abi.encodeWithSelector(IERC20.transfer.selector, receiver, amount), revertData);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, revertData));
    s_offRamp.releaseOrMintToken(amount, originalSender, receiver, sourceTokenData, offchainTokenData);
  }
}

contract EVM2EVMOffRamp__releaseOrMintTokens is EVM2EVMOffRampSetup {
  function test_releaseOrMintTokens_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destFeeToken);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes memory originalSender = abi.encode(OWNER);

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(encodedSourceTokenData[0], (Internal.SourceTokenData));

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: originalSender,
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR,
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, encodedSourceTokenData, offchainTokenData, new uint32[](0)
    );

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintTokens_destDenominatedDecimals_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount = 100;
    uint256 destinationDenominationMultiplier = 1000;
    srcTokenAmounts[1].amount = amount;

    bytes memory originalSender = abi.encode(OWNER);
    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    address pool = s_destPoolBySourceToken[srcTokenAmounts[1].token];
    address destToken = s_destTokenBySourceToken[srcTokenAmounts[1].token];

    MaybeRevertingBurnMintTokenPool(pool).setReleaseOrMintMultiplier(destinationDenominationMultiplier);

    Client.EVMTokenAmount[] memory destTokenAmounts = s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, encodedSourceTokenData, offchainTokenData, new uint32[](0)
    );

    assertEq(destTokenAmounts[1].amount, amount * destinationDenominationMultiplier);
    assertEq(destTokenAmounts[1].token, destToken);
  }

  function test_OverValueWithARLOff_Success() public {
    // Set a high price to trip the ARL
    uint224 tokenPrice = 3 ** 128;
    Internal.PriceUpdates memory priceUpdates = _getSingleTokenPriceUpdateStruct(s_destFeeToken, tokenPrice);
    s_feeQuoter.updatePrices(priceUpdates);

    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes memory originalSender = abi.encode(OWNER);

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    bytes[] memory sourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector,
        _getInboundRateLimiterConfig().capacity,
        (amount1 * tokenPrice) / 1e18
      )
    );

    // // Expect to fail from ARL
    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, sourceTokenData, offchainTokenData, new uint32[](0)
    );

    // Configure ARL off for token
    EVM2EVMOffRamp.RateLimitToken[] memory removes = new EVM2EVMOffRamp.RateLimitToken[](1);
    removes[0] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceFeeToken, destToken: s_destFeeToken});
    s_offRamp.updateRateLimitTokens(removes, new EVM2EVMOffRamp.RateLimitToken[](0));

    // Expect the call now succeeds
    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, sourceTokenData, offchainTokenData, new uint32[](0)
    );
  }

  // Revert

  function test_TokenHandlingError_Reverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes memory unknownError = bytes("unknown error");
    s_maybeRevertingPool.setShouldRevert(unknownError);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, unknownError));

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts,
      abi.encode(OWNER),
      OWNER,
      _getDefaultSourceTokenData(srcTokenAmounts),
      new bytes[](srcTokenAmounts.length),
      new uint32[](0)
    );
  }

  function test_releaseOrMintTokens_InvalidDataLengthReturnData_Revert() public {
    uint256 amount = 100;
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    srcTokenAmounts[0].amount = amount;

    bytes memory originalSender = abi.encode(OWNER);
    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(encodedSourceTokenData[0], (Internal.SourceTokenData));

    vm.mockCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: originalSender,
          receiver: OWNER,
          amount: amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR,
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData[0]
        })
      ),
      // Includes the amount twice, this will revert due to the return data being to long
      abi.encode(amount, amount)
    );

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.InvalidDataLength.selector, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, 64)
    );

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, encodedSourceTokenData, offchainTokenData, new uint32[](0)
    );
  }

  function test_releaseOrMintTokens_InvalidEVMAddress_Revert() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes memory originalSender = abi.encode(OWNER);
    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    bytes[] memory sourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    bytes memory wrongAddress = abi.encode(address(1000), address(10000), address(10000));

    sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[srcTokenAmounts[0].token]),
        destTokenAddress: wrongAddress,
        extraData: "",
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, wrongAddress));

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, sourceTokenData, offchainTokenData, new uint32[](0)
    );
  }

  function test_RateLimitErrors_Reverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes[] memory rateLimitErrors = new bytes[](5);
    rateLimitErrors[0] = abi.encodeWithSelector(RateLimiter.BucketOverfilled.selector);
    rateLimitErrors[1] =
      abi.encodeWithSelector(RateLimiter.AggregateValueMaxCapacityExceeded.selector, uint256(100), uint256(1000));
    rateLimitErrors[2] =
      abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, uint256(42), 1, s_sourceTokens[0]);
    rateLimitErrors[3] = abi.encodeWithSelector(
      RateLimiter.TokenMaxCapacityExceeded.selector, uint256(100), uint256(1000), s_sourceTokens[0]
    );
    rateLimitErrors[4] =
      abi.encodeWithSelector(RateLimiter.TokenRateLimitReached.selector, uint256(42), 1, s_sourceTokens[0]);

    for (uint256 i = 0; i < rateLimitErrors.length; ++i) {
      s_maybeRevertingPool.setShouldRevert(rateLimitErrors[i]);

      vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, rateLimitErrors[i]));

      s_offRamp.releaseOrMintTokens(
        srcTokenAmounts,
        abi.encode(OWNER),
        OWNER,
        _getDefaultSourceTokenData(srcTokenAmounts),
        new bytes[](srcTokenAmounts.length),
        new uint32[](0)
      );
    }
  }

  function test__releaseOrMintTokens_NotACompatiblePool_Reverts() public {
    address fakePoolAddress = makeAddr("Doesn't exist");

    bytes[] memory sourceTokenData = new bytes[](1);
    sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(fakePoolAddress),
        destTokenAddress: abi.encode(fakePoolAddress),
        extraData: "",
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.NotACompatiblePool.selector, address(0)));
    s_offRamp.releaseOrMintTokens(
      new Client.EVMTokenAmount[](1),
      abi.encode(makeAddr("original_sender")),
      OWNER,
      sourceTokenData,
      new bytes[](1),
      new uint32[](0)
    );
  }

  function test_PriceNotFoundForToken_Reverts() public {
    // Set token price to 0
    s_feeQuoter.updatePrices(_getSingleTokenPriceUpdateStruct(s_destFeeToken, 0));

    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes memory originalSender = abi.encode(OWNER);

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    bytes[] memory sourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectRevert(abi.encodeWithSelector(AggregateRateLimiter.PriceNotFoundForToken.selector, s_destFeeToken));

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, originalSender, OWNER, sourceTokenData, offchainTokenData, new uint32[](0)
    );
  }

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 1024
  // Uint256 gives a good range of values to test, both inside and outside of the eth address space.
  function test_Fuzz__releaseOrMintTokens_AnyRevertIsCaught_Success(uint256 destPool) public {
    // Input 447301751254033913445893214690834296930546521452, which is 0x4E59B44847B379578588920CA78FBF26C0B4956C
    // triggers some Create2Deployer and causes it to fail
    vm.assume(destPool != 447301751254033913445893214690834296930546521452);
    bytes memory unusedVar = abi.encode(makeAddr("unused"));
    bytes[] memory sourceTokenData = new bytes[](1);
    sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: unusedVar,
        destTokenAddress: abi.encode(destPool),
        extraData: unusedVar,
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      })
    );

    try s_offRamp.releaseOrMintTokens(
      new Client.EVMTokenAmount[](1), unusedVar, OWNER, sourceTokenData, new bytes[](1), new uint32[](0)
    ) {} catch (bytes memory reason) {
      // Any revert should be a TokenHandlingError, InvalidEVMAddress, InvalidDataLength or NoContract as those are caught by the offramp
      assertTrue(
        bytes4(reason) == EVM2EVMOffRamp.TokenHandlingError.selector
          || bytes4(reason) == Internal.InvalidEVMAddress.selector
          || bytes4(reason) == EVM2EVMOffRamp.InvalidDataLength.selector
          || bytes4(reason) == CallWithExactGas.NoContract.selector
          || bytes4(reason) == EVM2EVMOffRamp.NotACompatiblePool.selector,
        "Expected TokenHandlingError or InvalidEVMAddress"
      );

      if (destPool > type(uint160).max) {
        assertEq(reason, abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, abi.encode(destPool)));
      }
    }
  }
}

contract EVM2EVMOffRamp_getAllRateLimitTokens is EVM2EVMOffRampSetup {
  function test_GetAllRateLimitTokens_Success() public view {
    (address[] memory sourceTokens, address[] memory destTokens) = s_offRamp.getAllRateLimitTokens();

    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      assertEq(s_sourceTokens[i], sourceTokens[i]);
      assertEq(s_destTokens[i], destTokens[i]);
    }
  }
}

contract EVM2EVMOffRamp_updateRateLimitTokens is EVM2EVMOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    // Clear rate limit tokens state
    EVM2EVMOffRamp.RateLimitToken[] memory remove = new EVM2EVMOffRamp.RateLimitToken[](s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      remove[i] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[i], destToken: s_destTokens[i]});
    }
    s_offRamp.updateRateLimitTokens(remove, new EVM2EVMOffRamp.RateLimitToken[](0));
  }

  function test_updateRateLimitTokens_Success() public {
    EVM2EVMOffRamp.RateLimitToken[] memory adds = new EVM2EVMOffRamp.RateLimitToken[](2);
    adds[0] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[0], destToken: s_destTokens[0]});
    adds[1] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[1], destToken: s_destTokens[1]});

    for (uint256 i = 0; i < adds.length; ++i) {
      vm.expectEmit();
      emit EVM2EVMOffRamp.TokenAggregateRateLimitAdded(adds[i].sourceToken, adds[i].destToken);
    }

    s_offRamp.updateRateLimitTokens(new EVM2EVMOffRamp.RateLimitToken[](0), adds);

    (address[] memory sourceTokens, address[] memory destTokens) = s_offRamp.getAllRateLimitTokens();

    for (uint256 i = 0; i < adds.length; ++i) {
      assertEq(adds[i].sourceToken, sourceTokens[i]);
      assertEq(adds[i].destToken, destTokens[i]);
    }
  }

  function test_updateRateLimitTokens_AddsAndRemoves_Success() public {
    EVM2EVMOffRamp.RateLimitToken[] memory adds = new EVM2EVMOffRamp.RateLimitToken[](3);
    adds[0] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[0], destToken: s_destTokens[0]});
    adds[1] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[1], destToken: s_destTokens[1]});
    // Add a duplicate, this should not revert the tx
    adds[2] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[1], destToken: s_destTokens[1]});

    EVM2EVMOffRamp.RateLimitToken[] memory removes = new EVM2EVMOffRamp.RateLimitToken[](1);
    removes[0] = adds[0];

    for (uint256 i = 0; i < adds.length - 1; ++i) {
      vm.expectEmit();
      emit EVM2EVMOffRamp.TokenAggregateRateLimitAdded(adds[i].sourceToken, adds[i].destToken);
    }

    s_offRamp.updateRateLimitTokens(removes, adds);

    for (uint256 i = 0; i < removes.length; ++i) {
      vm.expectEmit();
      emit EVM2EVMOffRamp.TokenAggregateRateLimitRemoved(removes[i].sourceToken, removes[i].destToken);
    }

    s_offRamp.updateRateLimitTokens(removes, new EVM2EVMOffRamp.RateLimitToken[](0));

    (address[] memory sourceTokens, address[] memory destTokens) = s_offRamp.getAllRateLimitTokens();

    assertEq(1, sourceTokens.length);
    assertEq(adds[1].sourceToken, sourceTokens[0]);

    assertEq(1, destTokens.length);
    assertEq(adds[1].destToken, destTokens[0]);
  }

  function test_Fuzz_UpdateRateLimitTokens(uint8 numTokens) public {
    // Needs to be more than 1 so that the division doesn't round down and the even makes the comparisons simpler
    vm.assume(numTokens > 1 && numTokens % 2 == 0);

    // Clear the Rate limit tokens array so the test can start from a baseline
    (address[] memory sourceTokens, address[] memory destTokens) = s_offRamp.getAllRateLimitTokens();
    EVM2EVMOffRamp.RateLimitToken[] memory removes = new EVM2EVMOffRamp.RateLimitToken[](sourceTokens.length);
    for (uint256 x = 0; x < removes.length; x++) {
      removes[x] = EVM2EVMOffRamp.RateLimitToken({sourceToken: sourceTokens[x], destToken: destTokens[x]});
    }
    s_offRamp.updateRateLimitTokens(removes, new EVM2EVMOffRamp.RateLimitToken[](0));

    // Sanity check that the rateLimitTokens were successfully cleared
    (sourceTokens, destTokens) = s_offRamp.getAllRateLimitTokens();
    assertEq(sourceTokens.length, 0, "sourceTokenLength should be zero");

    EVM2EVMOffRamp.RateLimitToken[] memory adds = new EVM2EVMOffRamp.RateLimitToken[](numTokens);

    for (uint256 x = 0; x < numTokens; x++) {
      address tokenAddr = vm.addr(x + 1);

      // Create an array of several fake tokens to add which are deployed on the same address on both chains for simplicity
      adds[x] = EVM2EVMOffRamp.RateLimitToken({sourceToken: tokenAddr, destToken: tokenAddr});
    }

    // Attempt to add the tokens to the RateLimitToken Array
    s_offRamp.updateRateLimitTokens(new EVM2EVMOffRamp.RateLimitToken[](0), adds);

    // Retrieve them from storage and make sure that they all match the expected adds
    (sourceTokens, destTokens) = s_offRamp.getAllRateLimitTokens();

    for (uint256 x = 0; x < sourceTokens.length; x++) {
      // Check that the tokens match the ones we generated earlier
      assertEq(sourceTokens[x], adds[x].sourceToken, "Source token doesn't match add");
      assertEq(destTokens[x], adds[x].sourceToken, "dest Token doesn't match add");
    }

    // Attempt to remove half of the numTokens by removing the second half of the list and copying it to a removes array
    removes = new EVM2EVMOffRamp.RateLimitToken[](adds.length / 2);

    for (uint256 x = 0; x < adds.length / 2; x++) {
      removes[x] = adds[x + (adds.length / 2)];
    }

    // Attempt to update again, this time adding nothing and removing the second half of the tokens
    s_offRamp.updateRateLimitTokens(removes, new EVM2EVMOffRamp.RateLimitToken[](0));

    (sourceTokens, destTokens) = s_offRamp.getAllRateLimitTokens();
    assertEq(sourceTokens.length, adds.length / 2, "Current Rate limit token length is not half of the original adds");
    for (uint256 x = 0; x < sourceTokens.length; x++) {
      // Check that the tokens match the ones we generated earlier and didn't remove in the previous step
      assertEq(sourceTokens[x], adds[x].sourceToken, "Source token doesn't match add after removes");
      assertEq(destTokens[x], adds[x].destToken, "dest Token doesn't match add after removes");
    }
  }

  // Reverts

  function test_updateRateLimitTokens_NonOwner_Revert() public {
    EVM2EVMOffRamp.RateLimitToken[] memory addsAndRemoves = new EVM2EVMOffRamp.RateLimitToken[](4);

    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");

    s_offRamp.updateRateLimitTokens(addsAndRemoves, addsAndRemoves);
  }
}
