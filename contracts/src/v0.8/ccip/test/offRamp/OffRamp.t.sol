// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IFeeQuoter} from "../../interfaces/IFeeQuoter.sol";
import {IMessageInterceptor} from "../../interfaces/IMessageInterceptor.sol";
import {IRMNRemote} from "../../interfaces/IRMNRemote.sol";
import {IRouter} from "../../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {CallWithExactGas} from "../../../shared/call/CallWithExactGas.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {NonceManager} from "../../NonceManager.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {OffRampHelper} from "../helpers/OffRampHelper.sol";
import {ConformingReceiver} from "../helpers/receivers/ConformingReceiver.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MaybeRevertMessageReceiverNo165} from "../helpers/receivers/MaybeRevertMessageReceiverNo165.sol";
import {ReentrancyAbuserMultiRamp} from "../helpers/receivers/ReentrancyAbuserMultiRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OffRamp_constructor is OffRampSetup {
  function test_Constructor_Success() public {
    OffRamp.StaticConfig memory staticConfig = OffRamp.StaticConfig({
      chainSelector: DEST_CHAIN_SELECTOR,
      rmnRemote: s_mockRMNRemote,
      tokenAdminRegistry: address(s_tokenAdminRegistry),
      nonceManager: address(s_inboundNonceManager)
    });
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 1,
      onRamp: ON_RAMP_ADDRESS_2,
      isEnabled: true
    });

    OffRamp.SourceChainConfig memory expectedSourceChainConfig1 = OffRamp.SourceChainConfig({
      router: s_destRouter,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: sourceChainConfigs[0].onRamp
    });

    OffRamp.SourceChainConfig memory expectedSourceChainConfig2 = OffRamp.SourceChainConfig({
      router: s_destRouter,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: sourceChainConfigs[1].onRamp
    });

    uint64[] memory expectedSourceChainSelectors = new uint64[](2);
    expectedSourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;
    expectedSourceChainSelectors[1] = SOURCE_CHAIN_SELECTOR_1 + 1;

    vm.expectEmit();
    emit OffRamp.StaticConfigSet(staticConfig);

    vm.expectEmit();
    emit OffRamp.DynamicConfigSet(dynamicConfig);

    vm.expectEmit();
    emit OffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1);

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig1);

    vm.expectEmit();
    emit OffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1 + 1);

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1 + 1, expectedSourceChainConfig2);

    s_offRamp = new OffRampHelper(staticConfig, dynamicConfig, sourceChainConfigs);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });

    s_offRamp.setOCR3Configs(ocrConfigs);

    // Static config
    OffRamp.StaticConfig memory gotStaticConfig = s_offRamp.getStaticConfig();
    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(address(staticConfig.rmnRemote), address(gotStaticConfig.rmnRemote));
    assertEq(staticConfig.tokenAdminRegistry, gotStaticConfig.tokenAdminRegistry);

    // Dynamic config
    OffRamp.DynamicConfig memory gotDynamicConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, gotDynamicConfig);

    // OCR Config
    MultiOCR3Base.OCRConfig memory expectedOCRConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: 0,
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    MultiOCR3Base.OCRConfig memory gotOCRConfig = s_offRamp.latestConfigDetails(uint8(Internal.OCRPluginType.Execution));
    _assertOCRConfigEquality(expectedOCRConfig, gotOCRConfig);

    (uint64[] memory actualSourceChainSelectors, OffRamp.SourceChainConfig[] memory actualSourceChainConfigs) =
      s_offRamp.getAllSourceChainConfigs();

    _assertSourceChainConfigEquality(actualSourceChainConfigs[0], expectedSourceChainConfig1);
    _assertSourceChainConfigEquality(actualSourceChainConfigs[1], expectedSourceChainConfig2);

    // OffRamp initial values
    assertEq("OffRamp 1.6.0-dev", s_offRamp.typeAndVersion());
    assertEq(OWNER, s_offRamp.owner());
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    // assertion for source chain selector
    for (uint256 i = 0; i < expectedSourceChainSelectors.length; i++) {
      assertEq(expectedSourceChainSelectors[i], actualSourceChainSelectors[i]);
    }
  }

  // Revert
  function test_ZeroOnRampAddress_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: new bytes(0),
      isEnabled: true
    });

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_SourceChainSelector_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: 0,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    vm.expectRevert(OffRamp.ZeroChainSelectorNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_ZeroRMNRemote_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: IRMNRemote(ZERO_ADDRESS),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_ZeroChainSelector_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroChainSelectorNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: 0,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_ZeroTokenAdminRegistry_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: ZERO_ADDRESS,
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_ZeroNonceManager_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: ZERO_ADDRESS
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }
}

contract OffRamp_setDynamicConfig is OffRampSetup {
  function test_SetDynamicConfig_Success() public {
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));

    vm.expectEmit();
    emit OffRamp.DynamicConfigSet(dynamicConfig);

    s_offRamp.setDynamicConfig(dynamicConfig);

    OffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  function test_SetDynamicConfigWithInterceptor_Success() public {
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));
    dynamicConfig.messageInterceptor = address(s_inboundMessageInterceptor);

    vm.expectEmit();
    emit OffRamp.DynamicConfigSet(dynamicConfig);

    s_offRamp.setDynamicConfig(dynamicConfig);

    OffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  // Reverts

  function test_NonOwner_Revert() public {
    vm.startPrank(STRANGER);
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));

    vm.expectRevert("Only callable by owner");

    s_offRamp.setDynamicConfig(dynamicConfig);
  }

  function test_FeeQuoterZeroAddress_Revert() public {
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(ZERO_ADDRESS);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp.setDynamicConfig(dynamicConfig);
  }
}

contract OffRamp_ccipReceive is OffRampSetup {
  // Reverts

  function test_Reverts() public {
    Client.Any2EVMMessage memory message =
      _convertToGeneralMessage(_generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1));
    vm.expectRevert();
    s_offRamp.ccipReceive(message);
  }
}

contract OffRamp_executeSingleReport is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_SingleMessageNoTokens_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
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
    assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender), nonceBefore);
  }

  function test_SingleMessageNoTokensUnordered_Success() public {
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
    assertExecutionStateChangedEventLogs(
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
    assertExecutionStateChangedEventLogs(
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

  function test_SingleMessageNoTokensOtherChain_Success() public {
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

  function test_ReceiverError_Success() public {
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
    assertExecutionStateChangedEventLogs(
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

  function test_SkippedIncorrectNonce_Success() public {
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

  function test_SkippedIncorrectNonceStillExecutes_Success() public {
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
    assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test__execute_SkippedAlreadyExecutedMessage_Success() public {
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

    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test__execute_SkippedAlreadyExecutedMessageUnordered_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].header.nonce = 0;
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

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

    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  // Send a message to a contract that does not implement the CCIPReceiver interface
  // This should execute successfully.
  function test_SingleMessageToNonCCIPReceiver_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    messages[0].receiver = address(newReceiver);
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
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
    assertExecutionStateChangedEventLogs(
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
    assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[1].header.sequenceNumber,
      messages[1].header.messageId,
      _hashMessage(messages[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_TwoMessagesWithTokensAndGE_Success() public {
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

    assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    assertExecutionStateChangedEventLogs(
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

  function test_Fuzz_InterleavingOrderedAndUnorderedMessages_Success(
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

      assertExecutionStateChangedEventLogs(
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

  function test_InvalidSourcePoolAddress_Success() public {
    address fakePoolAddress = address(0x0000000000333333);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].tokenAmounts[0].sourcePoolAddress = abi.encode(fakePoolAddress);

    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);
    messages[1].header.messageId = _hashMessage(messages[1], ON_RAMP_ADDRESS_1);

    vm.recordLogs();

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        OffRamp.TokenHandlingError.selector,
        abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, abi.encode(fakePoolAddress))
      )
    );
  }

  function test_WithCurseOnAnotherSourceChain_Success() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_2, true);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_Unhealthy_Success() public {
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

  function test_MismatchingDestChainSelector_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);
    messages[0].header.destChainSelector = DEST_CHAIN_SELECTOR + 1;

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.InvalidMessageDestChainSelector.selector, messages[0].header.destChainSelector)
    );
    s_offRamp.executeSingleReport(executionReport, new OffRamp.GasLimitOverride[](0));
  }

  function test_UnhealthySingleChainCurse_Revert() public {
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

  function test_UnexpectedTokenData_Revert() public {
    Internal.ExecutionReport memory report = _generateReportFromMessages(
      SOURCE_CHAIN_SELECTOR_1, _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
    );
    report.offchainTokenData = new bytes[][](report.messages.length + 1);

    vm.expectRevert(OffRamp.UnexpectedTokenData.selector);

    s_offRamp.executeSingleReport(report, new OffRamp.GasLimitOverride[](0));
  }

  function test_EmptyReport_Revert() public {
    vm.expectRevert(OffRamp.EmptyReport.selector);
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

  function test_RootNotCommitted_Revert() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 0);
    vm.expectRevert(abi.encodeWithSelector(OffRamp.RootNotCommitted.selector, SOURCE_CHAIN_SELECTOR_1));

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
  }

  function test_ManualExecutionNotYetEnabled_Revert() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, BLOCK_TIME);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.ManualExecutionNotYetEnabled.selector, SOURCE_CHAIN_SELECTOR_1));

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
  }

  function test_NonExistingSourceChain_Revert() public {
    uint64 newSourceChainSelector = SOURCE_CHAIN_SELECTOR_1 + 1;
    bytes memory newOnRamp = abi.encode(ON_RAMP_ADDRESS, 1);

    Internal.Any2EVMRampMessage[] memory messages = _generateSingleBasicMessage(newSourceChainSelector, newOnRamp);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.SourceChainNotEnabled.selector, newSourceChainSelector));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(newSourceChainSelector, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_DisabledSourceChain_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_2, ON_RAMP_ADDRESS_2);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.SourceChainNotEnabled.selector, SOURCE_CHAIN_SELECTOR_2));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_2, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_TokenDataMismatch_Revert() public {
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

  function test_RouterYULCall_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    // gas limit too high, Router's external call should revert
    messages[0].gasLimit = 1e36;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport memory executionReport = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.recordLogs();
    s_offRamp.executeSingleReport(executionReport, new OffRamp.GasLimitOverride[](0));
    assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector)
    );
  }

  function test_RetryFailedMessageWithoutManualExecution_Revert() public {
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
    assertExecutionStateChangedEventLogs(
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
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });
  }
}

contract OffRamp_executeSingleMessage is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    vm.startPrank(address(s_offRamp));
  }

  function test_executeSingleMessage_NoTokens_Success() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_WithTokens_Success() public {
    Internal.Any2EVMRampMessage memory message =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)[0];
    bytes[] memory offchainTokenData = new bytes[](message.tokenAmounts.length);

    vm.expectCall(
      s_destPoolByToken[s_destTokens[0]],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: message.sender,
          receiver: message.receiver,
          amount: message.tokenAmounts[0].amount,
          localToken: message.tokenAmounts[0].destTokenAddress,
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: message.tokenAmounts[0].sourcePoolAddress,
          sourcePoolData: message.tokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    s_offRamp.executeSingleMessage(message, offchainTokenData, new uint32[](0));
  }

  function test_executeSingleMessage_WithVInterception_Success() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageInterceptor();
    vm.startPrank(address(s_offRamp));
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_NonContract_Success() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
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
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  // Reverts

  function test_TokenHandlingError_Revert() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    bytes memory errorMessage = "Random token pool issue";

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, errorMessage));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_ZeroGasDONExecution_Revert() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    message.gasLimit = 0;

    vm.expectRevert(abi.encodeWithSelector(OffRamp.ReceiverError.selector, ""));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_MessageSender_Revert() public {
    vm.stopPrank();
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    vm.expectRevert(OffRamp.CanOnlySelfCall.selector);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_WithFailingValidation_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageInterceptor();
    vm.startPrank(address(s_offRamp));
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_inboundMessageInterceptor.setMessageIdValidationState(message.header.messageId, true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_WithFailingValidationNoRouterCall_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageInterceptor();
    vm.startPrank(address(s_offRamp));

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);

    // Setup the receiver to a non-CCIP Receiver, which will skip the Router call (but should still perform the validation)
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    message.receiver = address(newReceiver);
    message.header.messageId = _hashMessage(message, ON_RAMP_ADDRESS_1);

    s_inboundMessageInterceptor.setMessageIdValidationState(message.header.messageId, true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }
}

contract OffRamp_batchExecute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_SingleReport_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);

    vm.recordLogs();
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
    );
    assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender), nonceBefore);
  }

  function test_MultipleReportsSameChain_Success() public {
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
    assertExecutionStateChangedEventLogs(
      logs,
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
      logs,
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
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

  function test_MultipleReportsDifferentChains_Success() public {
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

    assertExecutionStateChangedEventLogs(
      logs,
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
      logs,
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
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

  function test_MultipleReportsDifferentChainsSkipCursedChain_Success() public {
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

  function test_MultipleReportsSkipDuplicate_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectEmit();
    emit OffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    vm.recordLogs();
    s_offRamp.batchExecute(reports, new OffRamp.GasLimitOverride[][](2));
    assertExecutionStateChangedEventLogs(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_Unhealthy_Success() public {
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
  function test_ZeroReports_Revert() public {
    vm.expectRevert(OffRamp.EmptyReport.selector);
    s_offRamp.batchExecute(new Internal.ExecutionReport[](0), new OffRamp.GasLimitOverride[][](1));
  }

  function test_OutOfBoundsGasLimitsAccess_Revert() public {
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

contract OffRamp_manuallyExecute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();

    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_manuallyExecute_Success() public {
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
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_manuallyExecute_WithGasOverride_Success() public {
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
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_manuallyExecute_DoesNotRevertIfUntouched_Success() public {
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
    assertExecutionStateChangedEventLogs(
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

  function test_manuallyExecute_WithMultiReportGasOverride_Success() public {
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
      assertExecutionStateChangedEventLogs(
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
      assertExecutionStateChangedEventLogs(
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

  function test_manuallyExecute_WithPartialMessages_Success() public {
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

    assertExecutionStateChangedEventLogs(
      logs,
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
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

    assertExecutionStateChangedEventLogs(
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
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_manuallyExecute_LowGasLimit_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].gasLimit = 1;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].header.messageId = _hashMessage(messages[0], ON_RAMP_ADDRESS_1);

    vm.recordLogs();
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[][](1)
    );
    assertExecutionStateChangedEventLogs(
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
    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  // Reverts

  function test_manuallyExecute_ForkedChain_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ForkedChain.selector, chain1, chain2));

    OffRamp.GasLimitOverride[][] memory gasLimitOverrides = new OffRamp.GasLimitOverride[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_ManualExecGasLimitMismatchSingleReport_Revert() public {
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

  function test_manuallyExecute_GasLimitMismatchMultipleReports_Revert() public {
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

  function test_manuallyExecute_InvalidReceiverExecutionGasLimit_Revert() public {
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

  function test_manuallyExecute_DestinationGasAmountCountMismatch_Revert() public {
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

  function test_manuallyExecute_InvalidTokenGasOverride_Revert() public {
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

  function test_manuallyExecute_FailedTx_Revert() public {
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

  function test_manuallyExecute_ReentrancyFails_Success() public {
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
    assertExecutionStateChangedEventLogs(
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

  function test_manuallyExecute_MultipleReportsWithSingleCursedLane_Revert() public {
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

  function test_manuallyExecute_SourceChainSelectorMismatch_Revert() public {
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

contract OffRamp_execute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
  }

  // Asserts that execute completes
  function test_SingleReport_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    vm.recordLogs();

    _execute(reports);

    assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      _hashMessage(messages[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_MultipleReports_Success() public {
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

    assertExecutionStateChangedEventLogs(
      logs,
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
      logs,
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
      logs,
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      _hashMessage(messages2[0], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function test_LargeBatch_Success() public {
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
        assertExecutionStateChangedEventLogs(
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

  function test_MultipleReportsWithPartialValidationFailures_Success() public {
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

    assertExecutionStateChangedEventLogs(
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

    assertExecutionStateChangedEventLogs(
      logs,
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], ON_RAMP_ADDRESS_1),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertExecutionStateChangedEventLogs(
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

  function test_UnauthorizedTransmitter_Revert() public {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_NoConfig_Revert() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    bytes32[3] memory reportContext = [bytes32(""), s_configDigestExec, s_configDigestExec];

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_NoConfigWithOtherConfigPresent_Revert() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: s_F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    bytes32[3] memory reportContext = [bytes32(""), s_configDigestExec, s_configDigestExec];

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_WrongConfigWithSigners_Revert() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    s_configDigestExec = _getBasicConfigDigest(1, s_validSigners, s_validTransmitters);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: s_F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport[] memory reports = _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert();
    _execute(reports);
  }

  function test_ZeroReports_Revert() public {
    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](0);

    vm.expectRevert(OffRamp.EmptyReport.selector);
    _execute(reports);
  }

  function test_IncorrectArrayType_Revert() public {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    uint256[] memory wrongData = new uint256[](2);
    wrongData[0] = 1;

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.execute(reportContext, abi.encode(wrongData));
  }

  function test_NonArray_Revert() public {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReport memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.execute(reportContext, abi.encode(report));
  }
}

contract OffRamp_getExecutionState is OffRampSetup {
  mapping(uint64 sourceChainSelector => mapping(uint64 seqNum => Internal.MessageExecutionState state)) internal
    s_differentialExecutionState;

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function test_Fuzz_Differential_Success(
    uint64 sourceChainSelector,
    uint16[500] memory seqNums,
    uint8[500] memory values
  ) public {
    for (uint256 i = 0; i < seqNums.length; ++i) {
      // Only use the first three slots. This makes sure existing slots get overwritten
      // as the tests uses 500 sequence numbers.
      uint16 seqNum = seqNums[i] % 386;
      Internal.MessageExecutionState state = Internal.MessageExecutionState(values[i] % 4);
      s_differentialExecutionState[sourceChainSelector][seqNum] = state;
      s_offRamp.setExecutionStateHelper(sourceChainSelector, seqNum, state);
      assertEq(uint256(state), uint256(s_offRamp.getExecutionState(sourceChainSelector, seqNum)));
    }

    for (uint256 i = 0; i < seqNums.length; ++i) {
      uint16 seqNum = seqNums[i] % 386;
      Internal.MessageExecutionState expectedState = s_differentialExecutionState[sourceChainSelector][seqNum];
      assertEq(uint256(expectedState), uint256(s_offRamp.getExecutionState(sourceChainSelector, seqNum)));
    }
  }

  function test_GetExecutionState_Success() public {
    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 0, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 1, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (3 << 2));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 1, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 2, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2) + (3 << 4));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 127, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2) + (3 << 4) + (1 << 254));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 128, Internal.MessageExecutionState.SUCCESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2) + (3 << 4) + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 1), 2);

    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 0))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.IN_PROGRESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 1))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 2))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.IN_PROGRESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 127))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.SUCCESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 128))
    );
  }

  function test_GetDifferentChainExecutionState_Success() public {
    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 0, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), 0);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 127, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), 0);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 128, Internal.MessageExecutionState.SUCCESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 1), 2);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), 0);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 1), 0);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1 + 1, 127, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 1), 2);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), (3 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 1), 0);

    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 0))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.IN_PROGRESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 127))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.SUCCESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 128))
    );

    assertEq(
      uint256(Internal.MessageExecutionState.UNTOUCHED),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1 + 1, 0))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1 + 1, 127))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.UNTOUCHED),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1 + 1, 128))
    );
  }

  function test_FillExecutionState_Success() public {
    for (uint64 i = 0; i < 384; ++i) {
      s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, i, Internal.MessageExecutionState.FAILURE);
    }

    for (uint64 i = 0; i < 384; ++i) {
      assertEq(
        uint256(Internal.MessageExecutionState.FAILURE),
        uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, i))
      );
    }

    for (uint64 i = 0; i < 3; ++i) {
      assertEq(type(uint256).max, s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, i));
    }

    for (uint64 i = 0; i < 384; ++i) {
      s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, i, Internal.MessageExecutionState.IN_PROGRESS);
    }

    for (uint64 i = 0; i < 384; ++i) {
      assertEq(
        uint256(Internal.MessageExecutionState.IN_PROGRESS),
        uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, i))
      );
    }

    for (uint64 i = 0; i < 3; ++i) {
      // 0x555... == 0b101010101010.....
      assertEq(
        0x5555555555555555555555555555555555555555555555555555555555555555,
        s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, i)
      );
    }
  }
}

contract OffRamp_trialExecute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test_trialExecute_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);
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

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, errorMessage), err);

    // Expect the balance to remain the same
    assertEq(startingBalance, dstToken0.balanceOf(OWNER));
  }

  function test_RateLimitError_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    bytes memory errorMessage = abi.encodeWithSelector(RateLimiter.BucketOverfilled.selector);

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, errorMessage), err);
  }

  // TODO test actual pool exists but isn't compatible instead of just no pool
  function test_TokenPoolIsNotAContract_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 10000;
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);

    // Happy path, pool is correct
    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // address 0 has no contract
    assertEq(address(0).code.length, 0);

    message.tokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(address(0)),
      destTokenAddress: address(0),
      extraData: "",
      amount: message.tokenAmounts[0].amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    message.header.messageId = _hashMessage(message, ON_RAMP_ADDRESS_1);

    // Unhappy path, no revert but marked as failed.
    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, address(0)), err);

    address notAContract = makeAddr("not_a_contract");

    message.tokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(address(0)),
      destTokenAddress: notAContract,
      extraData: "",
      amount: message.tokenAmounts[0].amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    message.header.messageId = _hashMessage(message, ON_RAMP_ADDRESS_1);

    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, address(0)), err);
  }
}

contract OffRamp_releaseOrMintSingleToken is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test__releaseOrMintSingleToken_Success() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    IERC20 dstToken1 = IERC20(s_destTokenBySourceToken[token]);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
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
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: tokenAmount.sourcePoolAddress,
          sourcePoolData: tokenAmount.extraData,
          offchainTokenData: offchainTokenData
        })
      )
    );

    s_offRamp.releaseOrMintSingleToken(tokenAmount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData);

    assertEq(startingBalance + amount, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintToken_InvalidDataLength_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // Mock the call so returns 2 slots of data
    vm.mockCall(
      s_destTokenBySourceToken[token], abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER), abi.encode(0, 0)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.InvalidDataLength.selector, Internal.MAX_BALANCE_OF_RET_BYTES, 64));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR, "");
  }

  function test_releaseOrMintToken_TokenHandlingError_BalanceOf_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    bytes memory revertData = "failed to balanceOf";

    // Mock the call so returns 2 slots of data
    vm.mockCallRevert(
      s_destTokenBySourceToken[token], abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER), revertData
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, revertData));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR, "");
  }

  function test_releaseOrMintToken_ReleaseOrMintBalanceMismatch_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    uint256 mockedStaticBalance = 50000;

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.mockCall(
      s_destTokenBySourceToken[token],
      abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER),
      abi.encode(mockedStaticBalance)
    );

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.ReleaseOrMintBalanceMismatch.selector, amount, mockedStaticBalance, mockedStaticBalance
      )
    );

    s_offRamp.releaseOrMintSingleToken(tokenAmount, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR, "");
  }

  function test_releaseOrMintToken_skip_ReleaseOrMintBalanceMismatch_if_pool_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    uint256 mockedStaticBalance = 50000;

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // This should make the call fail if it does not skip the check
    vm.mockCall(
      s_destTokenBySourceToken[token],
      abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER),
      abi.encode(mockedStaticBalance)
    );

    s_offRamp.releaseOrMintSingleToken(
      tokenAmount, abi.encode(OWNER), s_destPoolBySourceToken[token], SOURCE_CHAIN_SELECTOR, ""
    );
  }

  function test__releaseOrMintSingleToken_NotACompatiblePool_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    vm.label(destToken, "destToken");
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: destToken,
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // Address(0) should always revert
    address returnedPool = address(0);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData);

    // A contract that doesn't support the interface should also revert
    returnedPool = address(s_offRamp);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData);
  }

  function test__releaseOrMintSingleToken_TokenHandlingError_transfer_Revert() public {
    address receiver = makeAddr("receiver");
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: destToken,
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    bytes memory revertData = "call reverted :o";

    vm.mockCallRevert(destToken, abi.encodeWithSelector(IERC20.transfer.selector, receiver, amount), revertData);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, revertData));
    s_offRamp.releaseOrMintSingleToken(
      tokenAmount, originalSender, receiver, SOURCE_CHAIN_SELECTOR_1, offchainTokenData
    );
  }
}

contract OffRamp_releaseOrMintTokens is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test_releaseOrMintTokens_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destFeeToken);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, new uint32[](0)
    );

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintTokens_WithGasOverride_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destFeeToken);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    uint32[] memory gasOverrides = new uint32[](sourceTokenAmounts.length);
    for (uint256 i = 0; i < gasOverrides.length; i++) {
      gasOverrides[i] = DEFAULT_TOKEN_DEST_GAS_OVERHEAD + 1;
    }
    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, gasOverrides
    );

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintTokens_destDenominatedDecimals_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount = 100;
    uint256 destinationDenominationMultiplier = 1000;
    srcTokenAmounts[1].amount = amount;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    address pool = s_destPoolBySourceToken[srcTokenAmounts[1].token];
    address destToken = s_destTokenBySourceToken[srcTokenAmounts[1].token];

    MaybeRevertingBurnMintTokenPool(pool).setReleaseOrMintMultiplier(destinationDenominationMultiplier);

    Client.EVMTokenAmount[] memory destTokenAmounts = s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, new uint32[](0)
    );
    assertEq(destTokenAmounts[1].amount, amount * destinationDenominationMultiplier);
    assertEq(destTokenAmounts[1].token, destToken);
  }

  // Revert

  function test_TokenHandlingError_Reverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes memory unknownError = bytes("unknown error");
    s_maybeRevertingPool.setShouldRevert(unknownError);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, unknownError));

    s_offRamp.releaseOrMintTokens(
      _getDefaultSourceTokenData(srcTokenAmounts),
      abi.encode(OWNER),
      OWNER,
      SOURCE_CHAIN_SELECTOR_1,
      new bytes[](srcTokenAmounts.length),
      new uint32[](0)
    );
  }

  function test_releaseOrMintTokens_InvalidDataLengthReturnData_Revert() public {
    uint256 amount = 100;
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    srcTokenAmounts[0].amount = amount;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.mockCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      ),
      // Includes the amount twice, this will revert due to the return data being to long
      abi.encode(amount, amount)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.InvalidDataLength.selector, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, 64));

    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, new uint32[](0)
    );
  }

  function test__releaseOrMintTokens_PoolIsNotAPool_Reverts() public {
    // The offRamp is a contract, but not a pool
    address fakePoolAddress = address(s_offRamp);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = new Internal.Any2EVMTokenTransfer[](1);
    sourceTokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(fakePoolAddress),
      destTokenAddress: address(s_offRamp),
      extraData: "",
      amount: 1,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.expectRevert(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, address(0)));
    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, new bytes[](1), new uint32[](0)
    );
  }

  function test_releaseOrMintTokens_PoolDoesNotSupportDest_Reverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_3,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );
    vm.expectRevert();
    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_3, offchainTokenData, new uint32[](0)
    );
  }

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 1024
  // Uint256 gives a good range of values to test, both inside and outside of the eth address space.
  function test_Fuzz__releaseOrMintTokens_AnyRevertIsCaught_Success(
    address destPool
  ) public {
    // Input 447301751254033913445893214690834296930546521452, which is 0x4E59B44847B379578588920CA78FBF26C0B4956C
    // triggers some Create2Deployer and causes it to fail
    vm.assume(destPool != 0x4e59b44847b379578588920cA78FbF26c0B4956C);
    bytes memory unusedVar = abi.encode(makeAddr("unused"));
    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = new Internal.Any2EVMTokenTransfer[](1);
    sourceTokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: unusedVar,
      destTokenAddress: destPool,
      extraData: unusedVar,
      amount: 1,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    try s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, new bytes[](1), new uint32[](0)
    ) {} catch (bytes memory reason) {
      // Any revert should be a TokenHandlingError, InvalidEVMAddress, InvalidDataLength or NoContract as those are caught by the offramp
      assertTrue(
        bytes4(reason) == OffRamp.TokenHandlingError.selector || bytes4(reason) == Internal.InvalidEVMAddress.selector
          || bytes4(reason) == OffRamp.InvalidDataLength.selector
          || bytes4(reason) == CallWithExactGas.NoContract.selector
          || bytes4(reason) == OffRamp.NotACompatiblePool.selector,
        "Expected TokenHandlingError or InvalidEVMAddress"
      );

      if (uint160(destPool) > type(uint160).max) {
        assertEq(reason, abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, abi.encode(destPool)));
      }
    }
  }
}

contract OffRamp_applySourceChainConfigUpdates is OffRampSetup {
  function test_ApplyZeroUpdates_Success() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.recordLogs();
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    // No logs emitted
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);

    assertEq(s_offRamp.getSourceChainSelectors().length, 0);
  }

  function test_AddNewChain_Success() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    OffRamp.SourceChainConfig memory expectedSourceChainConfig =
      OffRamp.SourceChainConfig({router: s_destRouter, isEnabled: true, minSeqNr: 1, onRamp: ON_RAMP_ADDRESS_1});

    vm.expectEmit();
    emit OffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1);

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig);

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    _assertSourceChainConfigEquality(s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig);
  }

  function test_ReplaceExistingChain_Success() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].isEnabled = false;
    OffRamp.SourceChainConfig memory expectedSourceChainConfig =
      OffRamp.SourceChainConfig({router: s_destRouter, isEnabled: false, minSeqNr: 1, onRamp: ON_RAMP_ADDRESS_1});

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig);

    vm.recordLogs();
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    // No log emitted for chain selector added (only for setting the config)
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    _assertSourceChainConfigEquality(s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig);

    uint256[] memory resultSourceChainSelectors = s_offRamp.getSourceChainSelectors();
    assertEq(resultSourceChainSelectors.length, 1);
  }

  function test_AddMultipleChains_Success() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 0),
      isEnabled: true
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 1,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 1),
      isEnabled: false
    });
    sourceChainConfigs[2] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 2,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 2),
      isEnabled: true
    });

    OffRamp.SourceChainConfig[] memory expectedSourceChainConfigs = new OffRamp.SourceChainConfig[](3);
    for (uint256 i = 0; i < 3; ++i) {
      expectedSourceChainConfigs[i] = OffRamp.SourceChainConfig({
        router: s_destRouter,
        isEnabled: sourceChainConfigs[i].isEnabled,
        minSeqNr: 1,
        onRamp: abi.encode(ON_RAMP_ADDRESS_1, i)
      });

      vm.expectEmit();
      emit OffRamp.SourceChainSelectorAdded(sourceChainConfigs[i].sourceChainSelector);

      vm.expectEmit();
      emit OffRamp.SourceChainConfigSet(sourceChainConfigs[i].sourceChainSelector, expectedSourceChainConfigs[i]);
    }

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    for (uint256 i = 0; i < 3; ++i) {
      _assertSourceChainConfigEquality(
        s_offRamp.getSourceChainConfig(sourceChainConfigs[i].sourceChainSelector), expectedSourceChainConfigs[i]
      );
    }
  }

  function test_Fuzz_applySourceChainConfigUpdate_Success(
    OffRamp.SourceChainConfigArgs memory sourceChainConfigArgs
  ) public {
    // Skip invalid inputs
    vm.assume(sourceChainConfigArgs.sourceChainSelector != 0);
    vm.assume(sourceChainConfigArgs.onRamp.length != 0);
    vm.assume(address(sourceChainConfigArgs.router) != address(0));

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });
    sourceChainConfigs[1] = sourceChainConfigArgs;

    // Handle cases when an update occurs
    bool isNewChain = sourceChainConfigs[1].sourceChainSelector != SOURCE_CHAIN_SELECTOR_1;
    if (!isNewChain) {
      sourceChainConfigs[1].onRamp = sourceChainConfigs[0].onRamp;
    }

    OffRamp.SourceChainConfig memory expectedSourceChainConfig = OffRamp.SourceChainConfig({
      router: sourceChainConfigArgs.router,
      isEnabled: sourceChainConfigArgs.isEnabled,
      minSeqNr: 1,
      onRamp: sourceChainConfigArgs.onRamp
    });

    if (isNewChain) {
      vm.expectEmit();
      emit OffRamp.SourceChainSelectorAdded(sourceChainConfigArgs.sourceChainSelector);
    }

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(sourceChainConfigArgs.sourceChainSelector, expectedSourceChainConfig);

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    _assertSourceChainConfigEquality(
      s_offRamp.getSourceChainConfig(sourceChainConfigArgs.sourceChainSelector), expectedSourceChainConfig
    );
  }

  function test_ReplaceExistingChainOnRamp_Success() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].onRamp = ON_RAMP_ADDRESS_2;

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(
      SOURCE_CHAIN_SELECTOR_1,
      OffRamp.SourceChainConfig({router: s_destRouter, isEnabled: true, minSeqNr: 1, onRamp: ON_RAMP_ADDRESS_2})
    );
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  // Reverts

  function test_ZeroOnRampAddress_Revert() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: new bytes(0),
      isEnabled: true
    });

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].onRamp = abi.encode(address(0));
    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_RouterAddress_Revert() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: IRouter(address(0)),
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_ZeroSourceChainSelector_Revert() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: 0,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    vm.expectRevert(OffRamp.ZeroChainSelectorNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_InvalidOnRampUpdate_Revert() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "test #2"
    });

    _commit(
      OffRamp.CommitReport({
        priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
        merkleRoots: roots,
        rmnSignatures: s_rmnSignatures,
        rmnRawVs: 0
      }),
      s_latestSequenceNumber
    );

    vm.stopPrank();
    vm.startPrank(OWNER);

    sourceChainConfigs[0].onRamp = ON_RAMP_ADDRESS_2;

    vm.expectRevert(abi.encodeWithSelector(OffRamp.InvalidOnRampUpdate.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }
}

contract OffRamp_commit is OffRampSetup {
  uint64 internal s_maxInterval = 12;

  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();

    s_latestSequenceNumber = uint64(uint256(s_configDigestCommit));
  }

  function test_ReportAndPriceUpdate_Success() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(commitReport.merkleRoots, commitReport.priceUpdates);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_maxInterval + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_ReportOnlyRootSuccess_gas() public {
    uint64 max1 = 931;
    bytes32 root = "Only a single root";

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: max1,
      merkleRoot: root
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(commitReport.merkleRoots, commitReport.priceUpdates);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(max1 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(block.timestamp, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR_1, root));
  }

  function test_StaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenStartPrice = IFeeQuoter(s_offRamp.getDynamicConfig().feeQuoter).getTokenPrice(s_sourceFeeToken).value;

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: maxSeq,
      merkleRoot: "stale report 1"
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(commitReport.merkleRoots, commitReport.priceUpdates);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    commitReport.merkleRoots[0].minSeqNr = maxSeq + 1;
    commitReport.merkleRoots[0].maxSeqNr = maxSeq * 2;
    commitReport.merkleRoots[0].merkleRoot = "stale report 2";

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(commitReport.merkleRoots, commitReport.priceUpdates);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq * 2 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(tokenStartPrice, IFeeQuoter(s_offRamp.getDynamicConfig().feeQuoter).getTokenPrice(s_sourceFeeToken).value);
  }

  function test_OnlyTokenPriceUpdates_Success() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_OnlyGasPriceUpdates_Success() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_PriceSequenceNumberCleared_Success() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    vm.startPrank(OWNER);
    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    // Execution plugin OCR config should not clear latest epoch and round
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    // Commit plugin config should clear latest epoch & round
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: s_F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    // The same sequence number can be reported again
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_ValidPriceUpdateThenStaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenPrice1 = 4e18;
    uint224 tokenPrice2 = 5e18;
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice1),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, tokenPrice1, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: maxSeq,
      merkleRoot: "stale report"
    });
    commitReport.priceUpdates = _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice2);
    commitReport.merkleRoots = roots;

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(commitReport.merkleRoots, commitReport.priceUpdates);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(tokenPrice1, IFeeQuoter(s_offRamp.getDynamicConfig().feeQuoter).getTokenPrice(s_sourceFeeToken).value);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  // Reverts

  function test_UnauthorizedTransmitter_Revert() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[3] memory reportContext =
      [s_configDigestCommit, bytes32(uint256(s_latestSequenceNumber)), s_configDigestCommit];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, s_F + 1);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_NoConfig_Revert() public {
    _redeployOffRampWithNoOCRConfigs();

    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[3] memory reportContext = [bytes32(""), s_configDigestCommit, s_configDigestCommit];
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, s_F + 1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_NoConfigWithOtherConfigPresent_Revert() public {
    _redeployOffRampWithNoOCRConfigs();

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[3] memory reportContext = [bytes32(""), s_configDigestCommit, s_configDigestCommit];
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, s_F + 1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_FailedRMNVerification_Reverts() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    OffRamp.CommitReport memory commitReport = _constructCommitReport();
    vm.expectRevert();
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_Unhealthy_Revert() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, true);
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "Only a single root",
      onRampAddress: abi.encode(ON_RAMP_ADDRESS_1)
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectRevert(abi.encodeWithSelector(OffRamp.CursedByRMN.selector, roots[0].sourceChainSelector));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_InvalidRootRevert() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 4,
      merkleRoot: bytes32(0)
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectRevert(OffRamp.InvalidRoot.selector);
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_InvalidInterval_Revert() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 2,
      maxSeqNr: 2,
      merkleRoot: bytes32(0)
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.InvalidInterval.selector, roots[0].sourceChainSelector, roots[0].minSeqNr, roots[0].maxSeqNr
      )
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_InvalidIntervalMinLargerThanMax_Revert() public {
    s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR);
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 0,
      merkleRoot: bytes32(0)
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.InvalidInterval.selector, roots[0].sourceChainSelector, roots[0].minSeqNr, roots[0].maxSeqNr
      )
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_ZeroEpochAndRound_Revert() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectRevert(OffRamp.StaleCommitReport.selector);
    _commit(commitReport, 0);
  }

  function test_OnlyPriceUpdateStaleReport_Revert() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    _commit(commitReport, s_latestSequenceNumber);

    vm.expectRevert(OffRamp.StaleCommitReport.selector);
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_SourceChainNotEnabled_Revert() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: 0,
      onRampAddress: abi.encode(ON_RAMP_ADDRESS_1),
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "Only a single root"
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    vm.expectRevert(abi.encodeWithSelector(OffRamp.SourceChainNotEnabled.selector, 0));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RootAlreadyCommitted_Revert() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "Only a single root"
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });

    _commit(commitReport, s_latestSequenceNumber);
    commitReport.merkleRoots[0].minSeqNr = 3;
    commitReport.merkleRoots[0].maxSeqNr = 3;

    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.RootAlreadyCommitted.selector, roots[0].sourceChainSelector, roots[0].merkleRoot)
    );
    _commit(commitReport, ++s_latestSequenceNumber);
  }

  function test_CommitOnRampMismatch_Revert() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    commitReport.merkleRoots[0].onRampAddress = ON_RAMP_ADDRESS_2;

    vm.expectRevert(abi.encodeWithSelector(OffRamp.CommitOnRampMismatch.selector, ON_RAMP_ADDRESS_2, ON_RAMP_ADDRESS_1));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function _constructCommitReport() internal view returns (OffRamp.CommitReport memory) {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: s_maxInterval,
      merkleRoot: "test #2"
    });

    return OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots,
      rmnSignatures: s_rmnSignatures,
      rmnRawVs: 0
    });
  }
}

contract OffRamp_afterOC3ConfigSet is OffRampSetup {
  function test_afterOCR3ConfigSet_SignatureVerificationDisabled_Revert() public {
    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      new OffRamp.SourceChainConfigArgs[](0)
    );

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    vm.expectRevert(OffRamp.SignatureVerificationDisabled.selector);
    s_offRamp.setOCR3Configs(ocrConfigs);
  }
}
