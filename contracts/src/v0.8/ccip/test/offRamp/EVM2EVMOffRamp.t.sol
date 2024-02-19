// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../../interfaces/pools/IPool.sol";
import {ICommitStore} from "../../interfaces/ICommitStore.sol";

import {EVM2EVMOffRampSetup} from "./EVM2EVMOffRampSetup.t.sol";
import {OCR2Base} from "../ocr/OCR2Base.t.sol";
import {OCR2BaseNoChecks} from "../ocr/OCR2BaseNoChecks.t.sol";
import {Router} from "../../Router.sol";
import {ARM} from "../../ARM.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Client} from "../../libraries/Client.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";

import {MockCommitStore} from "../mocks/MockCommitStore.sol";
import {CallWithExactGas} from "../../../shared/call/CallWithExactGas.sol";
import {ConformingReceiver} from "../helpers/receivers/ConformingReceiver.sol";
import {MaybeRevertMessageReceiverNo165} from "../helpers/receivers/MaybeRevertMessageReceiverNo165.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {ReentrancyAbuser} from "../helpers/receivers/ReentrancyAbuser.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMOffRamp_constructor is EVM2EVMOffRampSetup {
  event ConfigSet(EVM2EVMOffRamp.StaticConfig staticConfig, EVM2EVMOffRamp.DynamicConfig dynamicConfig);
  event PoolAdded(address token, address pool);

  function testConstructorSuccess() public {
    EVM2EVMOffRamp.StaticConfig memory staticConfig = EVM2EVMOffRamp.StaticConfig({
      commitStore: address(s_mockCommitStore),
      chainSelector: DEST_CHAIN_SELECTOR,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      onRamp: ON_RAMP_ADDRESS,
      prevOffRamp: address(0),
      armProxy: address(s_mockARM)
    });
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(
      address(s_destRouter),
      address(s_priceRegistry)
    );
    IERC20[] memory sourceTokens = getCastedSourceTokens();
    IPool[] memory castedPools = getCastedDestinationPools();

    for (uint256 i = 0; i < sourceTokens.length; ++i) {
      vm.expectEmit();
      emit PoolAdded(address(sourceTokens[i]), address(castedPools[i]));
    }

    s_offRamp = new EVM2EVMOffRampHelper(staticConfig, sourceTokens, castedPools, getInboundRateLimiterConfig());

    s_offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );

    // Static config
    EVM2EVMOffRamp.StaticConfig memory gotStaticConfig = s_offRamp.getStaticConfig();
    assertEq(staticConfig.commitStore, gotStaticConfig.commitStore);
    assertEq(staticConfig.sourceChainSelector, gotStaticConfig.sourceChainSelector);
    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.onRamp, gotStaticConfig.onRamp);
    assertEq(staticConfig.prevOffRamp, gotStaticConfig.prevOffRamp);

    // Dynamic config
    EVM2EVMOffRamp.DynamicConfig memory gotDynamicConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, gotDynamicConfig);

    // Pools & tokens
    IERC20[] memory pools = s_offRamp.getSupportedTokens();
    assertEq(pools.length, s_sourceTokens.length);
    assertTrue(address(pools[0]) == address(s_sourceTokens[0]));
    assertTrue(address(pools[1]) == address(s_sourceTokens[1]));
    assertEq(address(s_offRamp.getPoolByDestToken(IERC20(s_destTokens[0]))), address(s_destPools[0]));

    (uint32 configCount, uint32 blockNumber, ) = s_offRamp.latestConfigDetails();
    assertEq(1, configCount);
    assertEq(block.number, blockNumber);

    // OffRamp initial values
    assertEq("EVM2EVMOffRamp 1.5.0-dev", s_offRamp.typeAndVersion());
    assertEq(OWNER, s_offRamp.owner());
  }

  // Revert
  function testTokenConfigMismatchReverts() public {
    vm.expectRevert(EVM2EVMOffRamp.InvalidTokenPoolConfig.selector);

    IPool[] memory pools = new IPool[](1);

    IERC20[] memory wrongTokens = new IERC20[](5);
    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(s_mockCommitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        prevOffRamp: address(0),
        armProxy: address(s_mockARM)
      }),
      wrongTokens,
      pools,
      getInboundRateLimiterConfig()
    );
  }

  function testZeroOnRampAddressReverts() public {
    IPool[] memory pools = new IPool[](2);
    pools[0] = IPool(s_sourcePools[0]);
    pools[1] = new LockReleaseTokenPool(
      IERC20(s_sourceTokens[1]),
      new address[](0),
      address(s_mockARM),
      true,
      address(s_destRouter)
    );

    vm.expectRevert(EVM2EVMOffRamp.ZeroAddressNotAllowed.selector);

    RateLimiter.Config memory rateLimitConfig = RateLimiter.Config({isEnabled: true, rate: 1e20, capacity: 1e20});

    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(s_mockCommitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ZERO_ADDRESS,
        prevOffRamp: address(0),
        armProxy: address(s_mockARM)
      }),
      getCastedSourceTokens(),
      pools,
      rateLimitConfig
    );
  }

  function testCommitStoreAlreadyInUseReverts() public {
    s_mockCommitStore.setExpectedNextSequenceNumber(2);

    vm.expectRevert(EVM2EVMOffRamp.CommitStoreAlreadyInUse.selector);

    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(s_mockCommitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        prevOffRamp: address(0),
        armProxy: address(s_mockARM)
      }),
      getCastedSourceTokens(),
      getCastedDestinationPools(),
      getInboundRateLimiterConfig()
    );
  }
}

contract EVM2EVMOffRamp_setDynamicConfig is EVM2EVMOffRampSetup {
  // OffRamp event
  event ConfigSet(EVM2EVMOffRamp.StaticConfig staticConfig, EVM2EVMOffRamp.DynamicConfig dynamicConfig);

  function testSetDynamicConfigSuccess() public {
    EVM2EVMOffRamp.StaticConfig memory staticConfig = s_offRamp.getStaticConfig();
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(USER_3, address(s_priceRegistry));
    bytes memory onchainConfig = abi.encode(dynamicConfig);

    vm.expectEmit();
    emit ConfigSet(staticConfig, dynamicConfig);

    vm.expectEmit();
    uint32 configCount = 1;
    emit ConfigSet(
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
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      onchainConfig,
      s_offchainConfigVersion,
      abi.encode("")
    );

    EVM2EVMOffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  function testNonOwnerReverts() public {
    changePrank(STRANGER);
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(USER_3, address(s_priceRegistry));

    vm.expectRevert("Only callable by owner");

    s_offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );
  }

  function testRouterZeroAddressReverts() public {
    EVM2EVMOffRamp.DynamicConfig memory dynamicConfig = generateDynamicOffRampConfig(ZERO_ADDRESS, ZERO_ADDRESS);

    vm.expectRevert(EVM2EVMOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );
  }
}

contract EVM2EVMOffRamp_metadataHash is EVM2EVMOffRampSetup {
  function testMetadataHashSuccess() public {
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

  function testReverts() public {
    Client.Any2EVMMessage memory message = _convertToGeneralMessage(_generateAny2EVMMessageNoTokens(1));
    vm.expectRevert();
    s_offRamp.ccipReceive(message);
  }
}

contract EVM2EVMOffRamp_execute is EVM2EVMOffRampSetup {
  error PausedError();

  function testSingleMessageNoTokensSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

    messages[0].nonce++;
    messages[0].sequenceNumber++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint64 nonceBefore = s_offRamp.getSenderNonce(messages[0].sender);
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertGt(s_offRamp.getSenderNonce(messages[0].sender), nonceBefore);
  }

  function testReceiverErrorSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
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
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(uint64(1), s_offRamp.getSenderNonce(address(OWNER)));
  }

  function testStrictUntouchedToSuccessSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    messages[0].strict = true;
    messages[0].receiver = address(s_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    // Nonce should increment on a strict untouched -> success.
    assertEq(uint64(0), s_offRamp.getSenderNonce(address(OWNER)));
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(uint64(1), s_offRamp.getSenderNonce(address(OWNER)));
  }

  function testSkippedIncorrectNonceSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    messages[0].nonce++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit SkippedIncorrectNonce(messages[0].nonce, messages[0].sender);

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
  }

  function testSkippedIncorrectNonceStillExecutesSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();

    messages[1].nonce++;
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit SkippedIncorrectNonce(messages[1].nonce, messages[1].sender);

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
  }

  // Send a message to a contract that does not implement the CCIPReceiver interface
  // This should execute successfully.
  function testSingleMessageToNonCCIPReceiverSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    messages[0].receiver = address(newReceiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
  }

  function testSingleMessagesNoTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.resumeGasMetering();
    s_offRamp.execute(report, new uint256[](0));
  }

  function testTwoMessagesWithTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[1].sequenceNumber,
      messages[1].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.resumeGasMetering();
    s_offRamp.execute(report, new uint256[](0));
  }

  function testTwoMessagesWithTokensAndGESuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[1].sequenceNumber,
      messages[1].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(uint64(0), s_offRamp.getSenderNonce(OWNER));
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    assertEq(uint64(2), s_offRamp.getSenderNonce(OWNER));
  }

  // Reverts

  function testInvalidMessageIdReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].nonce++;
    // MessageID no longer matches hash.
    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);
    vm.expectRevert(EVM2EVMOffRamp.InvalidMessageId.selector);
    s_offRamp.execute(executionReport, new uint256[](0));
  }

  function testPausedReverts() public {
    s_mockCommitStore.pause();
    vm.expectRevert(PausedError.selector);
    s_offRamp.execute(_generateReportFromMessages(_generateMessagesWithTokens()), new uint256[](0));
  }

  function testUnhealthyReverts() public {
    s_mockARM.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    vm.expectRevert(EVM2EVMOffRamp.BadARMSignal.selector);
    s_offRamp.execute(_generateReportFromMessages(_generateMessagesWithTokens()), new uint256[](0));
    // Uncurse should succeed
    ARM.UnvoteToCurseRecord[] memory records = new ARM.UnvoteToCurseRecord[](1);
    records[0] = ARM.UnvoteToCurseRecord({curseVoteAddr: OWNER, cursesHash: bytes32(uint256(0)), forceUnvote: true});
    s_mockARM.ownerUnvoteToCurse(records);
    s_offRamp.execute(_generateReportFromMessages(_generateMessagesWithTokens()), new uint256[](0));
  }

  function testUnexpectedTokenDataReverts() public {
    Internal.ExecutionReport memory report = _generateReportFromMessages(_generateBasicMessages());
    report.offchainTokenData = new bytes[][](report.messages.length + 1);

    vm.expectRevert(EVM2EVMOffRamp.UnexpectedTokenData.selector);

    s_offRamp.execute(report, new uint256[](0));
  }

  function testEmptyReportReverts() public {
    vm.expectRevert(EVM2EVMOffRamp.EmptyReport.selector);
    s_offRamp.execute(
      Internal.ExecutionReport({
        proofs: new bytes32[](0),
        proofFlagBits: 0,
        messages: new Internal.EVM2EVMMessage[](0),
        offchainTokenData: new bytes[][](0)
      }),
      new uint256[](0)
    );
  }

  function testRootNotCommittedReverts() public {
    vm.mockCall(address(s_mockCommitStore), abi.encodeWithSelector(ICommitStore.verify.selector), abi.encode(0));
    vm.expectRevert(EVM2EVMOffRamp.RootNotCommitted.selector);

    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    vm.clearMockedCalls();
  }

  function testManualExecutionNotYetEnabledReverts() public {
    vm.mockCall(
      address(s_mockCommitStore),
      abi.encodeWithSelector(ICommitStore.verify.selector),
      abi.encode(BLOCK_TIME)
    );
    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionNotYetEnabled.selector);

    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    s_offRamp.execute(_generateReportFromMessages(messages), _getGasLimitsFromMessages(messages));
    vm.clearMockedCalls();
  }

  function testAlreadyExecutedReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);
    s_offRamp.execute(executionReport, new uint256[](0));
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.AlreadyExecuted.selector, messages[0].sequenceNumber));
    s_offRamp.execute(executionReport, new uint256[](0));
  }

  function testInvalidSourceChainReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].sourceChainSelector = SOURCE_CHAIN_SELECTOR + 1;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.InvalidSourceChain.selector, SOURCE_CHAIN_SELECTOR + 1));
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
  }

  function testUnsupportedNumberOfTokensReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    Client.EVMTokenAmount[] memory newTokens = new Client.EVMTokenAmount[](MAX_TOKENS_LENGTH + 1);
    messages[0].tokenAmounts = newTokens;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.UnsupportedNumberOfTokens.selector, messages[0].sequenceNumber)
    );
    s_offRamp.execute(report, new uint256[](0));
  }

  function testTokenDataMismatchReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    report.offchainTokenData[0] = new bytes[](messages[0].tokenAmounts.length + 1);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenDataMismatch.selector, messages[0].sequenceNumber));
    s_offRamp.execute(report, new uint256[](0));
  }

  function testMessageTooLargeReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].data = new bytes(MAX_DATA_SIZE + 1);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);
    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMOffRamp.MessageTooLarge.selector, MAX_DATA_SIZE, messages[0].data.length)
    );
    s_offRamp.execute(executionReport, new uint256[](0));
  }

  function testUnsupportedTokenReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateMessagesWithTokens();
    messages[0].tokenAmounts[0] = getCastedDestinationEVMTokenAmountsWithZeroAmounts()[0];
    messages[0].feeToken = messages[0].tokenAmounts[0].token;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    messages[1].messageId = Internal._hash(messages[1], s_offRamp.metadataHash());
    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ExecutionError.selector,
        abi.encodeWithSelector(EVM2EVMOffRamp.UnsupportedToken.selector, s_destTokens[0])
      )
    );
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
  }

  function testRouterYULCallReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    // gas limit too high, Router's external call should revert
    messages[0].gasLimit = 1e36;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(messages);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ExecutionError.selector,
        abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector)
      )
    );
    s_offRamp.execute(executionReport, new uint256[](0));
  }
}

contract EVM2EVMOffRamp_execute_upgrade is EVM2EVMOffRampSetup {
  event SkippedSenderWithPreviousRampMessageInflight(uint64 indexed nonce, address indexed sender);

  EVM2EVMOffRampHelper internal s_prevOffRamp;

  function setUp() public virtual override {
    EVM2EVMOffRampSetup.setUp();

    s_prevOffRamp = s_offRamp;

    deployOffRamp(s_mockCommitStore, s_destRouter, address(s_prevOffRamp));
  }

  function testV2Success() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
  }

  function testV2SenderNoncesReadsPreviousRampSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    uint64 startNonce = s_offRamp.getSenderNonce(messages[0].sender);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOffRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

      messages[0].nonce++;
      messages[0].sequenceNumber++;
      messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

      assertEq(startNonce + i, s_offRamp.getSenderNonce(messages[0].sender));
    }
  }

  function testV2NonceStartsAtV1NonceSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint64 startNonce = s_offRamp.getSenderNonce(messages[0].sender);

    s_prevOffRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

    assertEq(startNonce + 1, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(startNonce + 2, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce++;
    messages[0].sequenceNumber++;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(startNonce + 3, s_offRamp.getSenderNonce(messages[0].sender));
  }

  function testV2NonceNewSenderStartsAtZeroSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_prevOffRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

    address newSender = address(1234567);
    messages[0].sender = newSender;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    // new sender nonce in new offramp should go from 0 -> 1
    assertEq(s_offRamp.getSenderNonce(newSender), 0);
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(s_offRamp.getSenderNonce(newSender), 1);
  }

  function testV2OffRampNonceSkipsIfMsgInFlightSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    address newSender = address(1234567);
    messages[0].sender = newSender;
    messages[0].nonce = 2;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    uint64 startNonce = s_offRamp.getSenderNonce(messages[0].sender);

    // new offramp sees msg nonce higher than senderNonce
    // it waits for previous offramp to execute
    vm.expectEmit();
    emit SkippedSenderWithPreviousRampMessageInflight(messages[0].nonce, newSender);
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(startNonce, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce = 1;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    // previous offramp executes msg and increases nonce
    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    s_prevOffRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(startNonce + 1, s_offRamp.getSenderNonce(messages[0].sender));

    messages[0].nonce = 2;
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    // new offramp is able to execute
    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));
    assertEq(startNonce + 2, s_offRamp.getSenderNonce(messages[0].sender));
  }
}

contract EVM2EVMOffRamp_executeSingleMessage is EVM2EVMOffRampSetup {
  event MessageReceived();
  event Released(address indexed sender, address indexed recipient, uint256 amount);
  event Minted(address indexed sender, address indexed recipient, uint256 amount);

  function setUp() public virtual override {
    EVM2EVMOffRampSetup.setUp();
    changePrank(address(s_offRamp));
  }

  function testNoTokensSuccess() public {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function testTokensSuccess() public {
    Internal.EVM2EVMMessage memory message = _generateMessagesWithTokens()[0];
    bytes[] memory offchainTokenData = new bytes[](message.tokenAmounts.length);
    vm.expectCall(
      s_destPools[0],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        abi.encode(message.sender),
        message.receiver,
        message.tokenAmounts[0].amount,
        SOURCE_CHAIN_SELECTOR,
        abi.encode(message.sourceTokenData[0], offchainTokenData[0])
      )
    );

    s_offRamp.executeSingleMessage(message, offchainTokenData);
  }

  function testNonContractSuccess() public {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function testNonContractWithTokensSuccess() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;
    vm.expectEmit();
    emit Released(address(s_offRamp), STRANGER, amounts[0]);
    vm.expectEmit();
    emit Minted(address(s_offRamp), STRANGER, amounts[1]);
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  // Reverts

  function testTokenHandlingErrorReverts() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    bytes memory errorMessage = "Random token pool issue";

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    MaybeRevertingBurnMintTokenPool(s_destPools[1]).setShouldRevert(errorMessage);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, errorMessage));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function testZeroGasDONExecutionReverts() public {
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    message.gasLimit = 0;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.ReceiverError.selector, ""));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function testMessageSenderReverts() public {
    vm.stopPrank();
    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageNoTokens(1);
    vm.expectRevert(EVM2EVMOffRamp.CanOnlySelfCall.selector);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }
}

contract EVM2EVMOffRamp__report is EVM2EVMOffRampSetup {
  // Asserts that execute completes
  function testReportSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    s_offRamp.report(abi.encode(report));
  }
}

contract EVM2EVMOffRamp_manuallyExecute is EVM2EVMOffRampSetup {
  event ReentrancySucceeded();

  function testManualExecSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), new uint256[](messages.length));
  }

  function testManualExecWithGasOverrideSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint256[] memory gasLimitOverrides = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0] += 1;

    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  event MessageReceived();

  function testLowGasLimitManualExecSuccess() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].gasLimit = 1;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(EVM2EVMOffRamp.ReceiverError.selector, "")
    );
    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

    uint256[] memory gasLimitOverrides = new uint256[](1);
    gasLimitOverrides[0] = 100_000;

    vm.expectEmit();
    emit MessageReceived();

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimitOverrides);
  }

  function testManualExecForkedChainReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);
    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(OCR2BaseNoChecks.ForkedChain.selector, chain1, chain2));

    s_offRamp.manuallyExecute(report, _getGasLimitsFromMessages(messages));
  }

  function testManualExecGasLimitMismatchReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), new uint256[](0));

    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), new uint256[](messages.length - 1));

    vm.expectRevert(EVM2EVMOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), new uint256[](messages.length + 1));
  }

  function testManualExecInvalidGasLimitReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    uint256[] memory gasLimits = _getGasLimitsFromMessages(messages);
    gasLimits[0]--;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.InvalidManualExecutionGasLimit.selector, 0, gasLimits[0]));
    s_offRamp.manuallyExecute(_generateReportFromMessages(messages), gasLimits);
  }

  function testManualExecFailedTxReverts() public {
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    s_offRamp.execute(_generateReportFromMessages(messages), new uint256[](0));

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

  function testReentrancyManualExecuteFAILS() public {
    uint256 tokenAmount = 1e9;
    IERC20 tokenToAbuse = IERC20(s_destFeeToken);

    // This needs to be deployed before the source chain message is sent
    // because we need the address for the receiver.
    ReentrancyAbuser receiver = new ReentrancyAbuser(address(s_destRouter), s_offRamp);
    uint256 balancePre = tokenToAbuse.balanceOf(address(receiver));

    // For this test any message will be flagged as correct by the
    // commitStore. In a real scenario the abuser would have to actually
    // send the message that they want to replay.
    Internal.EVM2EVMMessage[] memory messages = _generateBasicMessages();
    messages[0].tokenAmounts = new Client.EVMTokenAmount[](1);
    messages[0].tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: tokenAmount});
    messages[0].sourceTokenData = new bytes[](1);
    messages[0].receiver = address(receiver);
    messages[0].messageId = Internal._hash(messages[0], s_offRamp.metadataHash());

    Internal.ExecutionReport memory report = _generateReportFromMessages(messages);

    // sets the report to be repeated on the ReentrancyAbuser to be able to replay
    receiver.setPayload(report);

    // The first entry should be fine and triggers the second entry. This one fails
    // but since it's an inner tx of the first one it is caught in the try-catch.
    // Since this is manual exec, the entire tx fails on any failure.
    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMOffRamp.ExecutionError.selector,
        abi.encodeWithSelector(
          EVM2EVMOffRamp.ReceiverError.selector,
          abi.encodeWithSelector(EVM2EVMOffRamp.AlreadyExecuted.selector, messages[0].sequenceNumber)
        )
      )
    );

    s_offRamp.manuallyExecute(report, _getGasLimitsFromMessages(messages));

    // Since the tx failed we don't release the tokens
    assertEq(tokenToAbuse.balanceOf(address(receiver)), balancePre);
  }
}

contract EVM2EVMOffRamp_getExecutionState is EVM2EVMOffRampSetup {
  mapping(uint64 seqNum => Internal.MessageExecutionState state) internal s_differentialExecutionState;

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function testFuzz_DifferentialSuccess(uint16[500] memory seqNums, uint8[500] memory values) public {
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

  function test_GetExecutionStateSuccess() public {
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

  function testFillExecutionStateSuccess() public {
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
      assertEq(
        0x5555555555555555555555555555555555555555555555555555555555555555,
        s_offRamp.getExecutionStateBitMap(i)
      );
    }
  }
}

contract EVM2EVMOffRamp__trialExecute is EVM2EVMOffRampSetup {
  function test_trialExecuteSuccess() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    IERC20 dstToken0 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken0.balanceOf(message.receiver);

    (Internal.MessageExecutionState newState, bytes memory err) = s_offRamp.trialExecute(
      message,
      new bytes[](message.tokenAmounts.length)
    );
    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // Check that the tokens were transferred
    assertEq(startingBalance + amounts[0], dstToken0.balanceOf(message.receiver));
  }

  function testTokenHandlingErrorIsCaughtSuccess() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    IERC20 dstToken0 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken0.balanceOf(OWNER);

    bytes memory errorMessage = "Random token pool issue";

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    MaybeRevertingBurnMintTokenPool(s_destPools[1]).setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) = s_offRamp.trialExecute(
      message,
      new bytes[](message.tokenAmounts.length)
    );
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, errorMessage), err);

    // Expect the balance to remain the same
    assertEq(startingBalance, dstToken0.balanceOf(OWNER));
  }

  function testRateLimitErrorSuccess() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    bytes memory errorMessage = abi.encodeWithSelector(RateLimiter.BucketOverfilled.selector);

    Internal.EVM2EVMMessage memory message = _generateAny2EVMMessageWithTokens(1, amounts);
    MaybeRevertingBurnMintTokenPool(s_destPools[1]).setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) = s_offRamp.trialExecute(
      message,
      new bytes[](message.tokenAmounts.length)
    );
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, errorMessage), err);
  }
}

contract EVM2EVMOffRamp__releaseOrMintTokens is EVM2EVMOffRampSetup {
  function test_releaseOrMintTokensSuccess() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes memory originalSender = abi.encode(OWNER);

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    bytes[] memory sourceTokenData = new bytes[](srcTokenAmounts.length);
    sourceTokenData[0] = abi.encode(0x87654321);

    vm.expectCall(
      s_destPools[0],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        originalSender,
        OWNER,
        srcTokenAmounts[0].amount,
        SOURCE_CHAIN_SELECTOR,
        abi.encode(sourceTokenData[0], offchainTokenData[0])
      )
    );

    s_offRamp.releaseOrMintTokens(srcTokenAmounts, originalSender, OWNER, sourceTokenData, offchainTokenData);

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  // Revert

  function testTokenHandlingErrorReverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes memory unknownError = bytes("unknown error");
    MaybeRevertingBurnMintTokenPool(s_destPools[1]).setShouldRevert(unknownError);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, unknownError));

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts,
      abi.encode(OWNER),
      OWNER,
      new bytes[](srcTokenAmounts.length),
      new bytes[](srcTokenAmounts.length)
    );
  }

  function testRateLimitErrorsReverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes[] memory rateLimitErrors = new bytes[](5);
    rateLimitErrors[0] = abi.encodeWithSelector(RateLimiter.BucketOverfilled.selector);
    rateLimitErrors[1] = abi.encodeWithSelector(
      RateLimiter.AggregateValueMaxCapacityExceeded.selector,
      uint256(100),
      uint256(1000)
    );
    rateLimitErrors[2] = abi.encodeWithSelector(
      RateLimiter.AggregateValueRateLimitReached.selector,
      uint256(42),
      1,
      s_sourceTokens[0]
    );
    rateLimitErrors[3] = abi.encodeWithSelector(
      RateLimiter.TokenMaxCapacityExceeded.selector,
      uint256(100),
      uint256(1000),
      s_sourceTokens[0]
    );
    rateLimitErrors[4] = abi.encodeWithSelector(
      RateLimiter.TokenRateLimitReached.selector,
      uint256(42),
      1,
      s_sourceTokens[0]
    );

    for (uint256 i = 0; i < rateLimitErrors.length; ++i) {
      MaybeRevertingBurnMintTokenPool(s_destPools[1]).setShouldRevert(rateLimitErrors[i]);

      vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.TokenHandlingError.selector, rateLimitErrors[i]));

      s_offRamp.releaseOrMintTokens(
        srcTokenAmounts,
        abi.encode(OWNER),
        OWNER,
        new bytes[](srcTokenAmounts.length),
        new bytes[](srcTokenAmounts.length)
      );
    }
  }

  function testUnsupportedTokenReverts() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.UnsupportedToken.selector, address(0)));
    s_offRamp.releaseOrMintTokens(tokenAmounts, bytes(""), OWNER, new bytes[](0), new bytes[](0));
  }
}

contract EVM2EVMOffRamp_applyPoolUpdates is EVM2EVMOffRampSetup {
  event PoolAdded(address token, address pool);
  event PoolRemoved(address token, address pool);

  function testApplyPoolUpdatesSuccess() public {
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({
      token: address(1),
      pool: address(
        new LockReleaseTokenPool(IERC20(address(1)), new address[](0), address(s_mockARM), true, address(s_destRouter))
      )
    });

    vm.expectEmit();
    emit PoolAdded(adds[0].token, adds[0].pool);

    s_offRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    assertEq(adds[0].pool, address(s_offRamp.getPoolBySourceToken(IERC20(adds[0].token))));

    vm.expectEmit();
    emit PoolRemoved(adds[0].token, adds[0].pool);

    s_offRamp.applyPoolUpdates(adds, new Internal.PoolUpdate[](0));

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.UnsupportedToken.selector, adds[0].token));
    s_offRamp.getPoolBySourceToken(IERC20(adds[0].token));
  }

  // Reverts
  function testOnlyCallableByOwnerReverts() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");

    s_offRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), new Internal.PoolUpdate[](0));
  }

  function testPoolAlreadyExistsReverts() public {
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](2);
    adds[0] = Internal.PoolUpdate({
      token: address(1),
      pool: address(
        new LockReleaseTokenPool(IERC20(address(1)), new address[](0), address(s_mockARM), true, address(s_destRouter))
      )
    });
    adds[1] = Internal.PoolUpdate({
      token: address(1),
      pool: address(
        new LockReleaseTokenPool(IERC20(address(1)), new address[](0), address(s_mockARM), true, address(s_destRouter))
      )
    });

    vm.expectRevert(EVM2EVMOffRamp.PoolAlreadyAdded.selector);

    s_offRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);
  }

  function testInvalidTokenPoolConfigReverts() public {
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({token: address(0), pool: address(2)});

    vm.expectRevert(EVM2EVMOffRamp.InvalidTokenPoolConfig.selector);

    s_offRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    adds[0] = Internal.PoolUpdate({token: address(1), pool: address(0)});

    vm.expectRevert(EVM2EVMOffRamp.InvalidTokenPoolConfig.selector);

    s_offRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);
  }

  function testPoolDoesNotExistReverts() public {
    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](1);
    removes[0] = Internal.PoolUpdate({
      token: address(1),
      pool: address(
        new LockReleaseTokenPool(IERC20(address(1)), new address[](0), address(s_mockARM), true, address(s_destRouter))
      )
    });

    vm.expectRevert(EVM2EVMOffRamp.PoolDoesNotExist.selector);

    s_offRamp.applyPoolUpdates(removes, new Internal.PoolUpdate[](0));
  }

  function testTokenPoolMismatchReverts() public {
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0] = Internal.PoolUpdate({
      token: address(1),
      pool: address(
        new LockReleaseTokenPool(IERC20(address(1)), new address[](0), address(s_mockARM), true, address(s_destRouter))
      )
    });
    s_offRamp.applyPoolUpdates(new Internal.PoolUpdate[](0), adds);

    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](1);
    removes[0] = Internal.PoolUpdate({
      token: address(1),
      pool: address(
        new LockReleaseTokenPool(
          IERC20(address(1000)),
          new address[](0),
          address(s_mockARM),
          true,
          address(s_destRouter)
        )
      )
    });

    vm.expectRevert(EVM2EVMOffRamp.TokenPoolMismatch.selector);

    s_offRamp.applyPoolUpdates(removes, adds);
  }
}

contract EVM2EVMOffRamp_getDestinationToken is EVM2EVMOffRampSetup {
  function testGetDestinationTokenSuccess() public {
    address expectedToken = address(IPool(s_destPools[0]).getToken());
    address actualToken = address(s_offRamp.getDestinationToken(IERC20(s_sourceTokens[0])));

    assertEq(expectedToken, actualToken);

    expectedToken = address(IPool(s_destPools[1]).getToken());
    actualToken = address(s_offRamp.getDestinationToken(IERC20(s_sourceTokens[1])));

    assertEq(expectedToken, actualToken);
  }

  function testUnsupportedTokenReverts() public {
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMOffRamp.UnsupportedToken.selector, DUMMY_CONTRACT_ADDRESS));
    s_offRamp.getDestinationToken(IERC20(DUMMY_CONTRACT_ADDRESS));
  }
}

contract EVM2EVMOffRamp_getDestinationTokens is EVM2EVMOffRampSetup {
  function testGetDestinationTokensSuccess() public {
    IERC20[] memory actualTokens = s_offRamp.getDestinationTokens();

    for (uint256 i = 0; i < actualTokens.length; ++i) {
      assertEq(address(s_destTokens[i]), address(actualTokens[i]));
    }
  }
}
