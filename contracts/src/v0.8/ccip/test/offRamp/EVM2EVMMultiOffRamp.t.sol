// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICommitStore} from "../../interfaces/ICommitStore.sol";
import {IMessageInterceptor} from "../../interfaces/IMessageInterceptor.sol";
import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";
import {IRMN} from "../../interfaces/IRMN.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {CallWithExactGas} from "../../../shared/call/CallWithExactGas.sol";
import {NonceManager} from "../../NonceManager.sol";
import {PriceRegistry} from "../../PriceRegistry.sol";
import {RMN} from "../../RMN.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {MerkleMultiProof} from "../../libraries/MerkleMultiProof.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";
import {EVM2EVMMultiOffRamp} from "../../offRamp/EVM2EVMMultiOffRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {EVM2EVMMultiOffRampHelper} from "../helpers/EVM2EVMMultiOffRampHelper.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MessageInterceptorHelper} from "../helpers/MessageInterceptorHelper.sol";
import {ConformingReceiver} from "../helpers/receivers/ConformingReceiver.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MaybeRevertMessageReceiverNo165} from "../helpers/receivers/MaybeRevertMessageReceiverNo165.sol";
import {ReentrancyAbuserMultiRamp} from "../helpers/receivers/ReentrancyAbuserMultiRamp.sol";
import {EVM2EVMMultiOffRampSetup} from "./EVM2EVMMultiOffRampSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMMultiOffRamp_constructor is EVM2EVMMultiOffRampSetup {
  function test_Constructor_Success() public {
    EVM2EVMMultiOffRamp.StaticConfig memory staticConfig = EVM2EVMMultiOffRamp.StaticConfig({
      chainSelector: DEST_CHAIN_SELECTOR,
      rmnProxy: address(s_mockRMN),
      tokenAdminRegistry: address(s_tokenAdminRegistry),
      nonceManager: address(s_inboundNonceManager)
    });
    EVM2EVMMultiOffRamp.DynamicConfig memory dynamicConfig =
      _generateDynamicMultiOffRampConfig(address(s_destRouter), address(s_priceRegistry));

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });
    sourceChainConfigs[1] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 1,
      onRamp: ON_RAMP_ADDRESS_2,
      isEnabled: true
    });

    EVM2EVMMultiOffRamp.SourceChainConfig memory expectedSourceChainConfig1 =
      EVM2EVMMultiOffRamp.SourceChainConfig({isEnabled: true, minSeqNr: 1, onRamp: sourceChainConfigs[0].onRamp});

    EVM2EVMMultiOffRamp.SourceChainConfig memory expectedSourceChainConfig2 =
      EVM2EVMMultiOffRamp.SourceChainConfig({isEnabled: true, minSeqNr: 1, onRamp: sourceChainConfigs[1].onRamp});

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.StaticConfigSet(staticConfig);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1 + 1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1 + 1, expectedSourceChainConfig2);

    s_offRamp = new EVM2EVMMultiOffRampHelper(staticConfig, sourceChainConfigs);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });

    s_offRamp.setDynamicConfig(dynamicConfig);
    s_offRamp.setOCR3Configs(ocrConfigs);

    // Static config
    EVM2EVMMultiOffRamp.StaticConfig memory gotStaticConfig = s_offRamp.getStaticConfig();
    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.rmnProxy, gotStaticConfig.rmnProxy);
    assertEq(staticConfig.tokenAdminRegistry, gotStaticConfig.tokenAdminRegistry);

    // Dynamic config
    EVM2EVMMultiOffRamp.DynamicConfig memory gotDynamicConfig = s_offRamp.getDynamicConfig();
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

    _assertSourceChainConfigEquality(
      s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig1
    );
    _assertSourceChainConfigEquality(
      s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1 + 1), expectedSourceChainConfig2
    );

    // OffRamp initial values
    assertEq("EVM2EVMMultiOffRamp 1.6.0-dev", s_offRamp.typeAndVersion());
    assertEq(OWNER, s_offRamp.owner());
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
  }

  // Revert
  function test_ZeroOnRampAddress_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: new bytes(0),
      isEnabled: true
    });

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      sourceChainConfigs
    );
  }

  function test_SourceChainSelector_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] =
      EVM2EVMMultiOffRamp.SourceChainConfigArgs({sourceChainSelector: 0, onRamp: ON_RAMP_ADDRESS_1, isEnabled: true});

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroChainSelectorNotAllowed.selector);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      sourceChainConfigs
    );
  }

  function test_ZeroRMNProxy_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: ZERO_ADDRESS,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      sourceChainConfigs
    );
  }

  function test_ZeroChainSelector_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroChainSelectorNotAllowed.selector);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        chainSelector: 0,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      sourceChainConfigs
    );
  }

  function test_ZeroTokenAdminRegistry_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: ZERO_ADDRESS,
        nonceManager: address(s_inboundNonceManager)
      }),
      sourceChainConfigs
    );
  }

  function test_ZeroNonceManager_Revert() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: ZERO_ADDRESS
      }),
      sourceChainConfigs
    );
  }
}

contract EVM2EVMMultiOffRamp_setDynamicConfig is EVM2EVMMultiOffRampSetup {
  function test_SetDynamicConfig_Success() public {
    EVM2EVMMultiOffRamp.DynamicConfig memory dynamicConfig =
      _generateDynamicMultiOffRampConfig(USER_3, address(s_priceRegistry));

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.DynamicConfigSet(dynamicConfig);

    s_offRamp.setDynamicConfig(dynamicConfig);

    EVM2EVMMultiOffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  function test_SetDynamicConfigWithValidator_Success() public {
    EVM2EVMMultiOffRamp.DynamicConfig memory dynamicConfig =
      _generateDynamicMultiOffRampConfig(USER_3, address(s_priceRegistry));
    dynamicConfig.messageValidator = address(s_inboundMessageValidator);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.DynamicConfigSet(dynamicConfig);

    s_offRamp.setDynamicConfig(dynamicConfig);

    EVM2EVMMultiOffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  // Reverts

  function test_NonOwner_Revert() public {
    vm.startPrank(STRANGER);
    EVM2EVMMultiOffRamp.DynamicConfig memory dynamicConfig =
      _generateDynamicMultiOffRampConfig(USER_3, address(s_priceRegistry));

    vm.expectRevert("Only callable by owner");

    s_offRamp.setDynamicConfig(dynamicConfig);
  }

  function test_RouterZeroAddress_Revert() public {
    EVM2EVMMultiOffRamp.DynamicConfig memory dynamicConfig =
      _generateDynamicMultiOffRampConfig(ZERO_ADDRESS, address(s_priceRegistry));

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp.setDynamicConfig(dynamicConfig);
  }

  function test_PriceRegistryZeroAddress_Revert() public {
    EVM2EVMMultiOffRamp.DynamicConfig memory dynamicConfig = _generateDynamicMultiOffRampConfig(USER_3, ZERO_ADDRESS);

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp.setDynamicConfig(dynamicConfig);
  }
}

contract EVM2EVMMultiOffRamp_ccipReceive is EVM2EVMMultiOffRampSetup {
  // Reverts

  function test_Reverts() public {
    Client.Any2EVMMessage memory message =
      _convertToGeneralMessage(_generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1));
    vm.expectRevert();
    s_offRamp.ccipReceive(message);
  }
}

contract EVM2EVMMultiOffRamp_executeSingleReport is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_SingleMessageNoTokens_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    messages[0].header.nonce++;
    messages[0].header.sequenceNumber++;
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender), nonceBefore);
  }

  function test_SingleMessageNoTokensUnordered_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].header.nonce = 0;
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    // Nonce never increments on unordered messages.
    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(
      s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender),
      nonceBefore,
      "nonce must remain unchanged on unordered messages"
    );

    messages[0].header.sequenceNumber++;
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    // Nonce never increments on unordered messages.
    nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
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
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messagesChain1), new uint256[](0)
    );

    uint64 nonceChain1 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messagesChain1[0].sender);
    assertGt(nonceChain1, 0);

    Internal.Any2EVMRampMessage[] memory messagesChain2 =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);
    assertEq(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, messagesChain2[0].sender), 0);

    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messagesChain2), new uint256[](0)
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
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1)
      )
    );
    // Nonce should increment on non-strict
    assertEq(uint64(0), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
    assertEq(uint64(1), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
  }

  function test_SkippedIncorrectNonce_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messages[0].header.nonce++;
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(
      messages[0].header.sourceChainSelector, messages[0].header.nonce, messages[0].sender
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function test_SkippedIncorrectNonceStillExecutes_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messages[1].header.nonce++;
    messages[1].header.messageId = Internal._hash(messages[1], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR_1, messages[1].header.nonce, messages[1].sender);

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function test__execute_SkippedAlreadyExecutedMessage_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function test__execute_SkippedAlreadyExecutedMessageUnordered_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].header.nonce = 0;
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  // Send a message to a contract that does not implement the CCIPReceiver interface
  // This should execute successfully.
  function test_SingleMessageToNonCCIPReceiver_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    messages[0].receiver = address(newReceiver);
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function test_SingleMessagesNoTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    Internal.ExecutionReportSingleChain memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.resumeGasMetering();
    s_offRamp.executeSingleReport(report, new uint256[](0));
  }

  function test_TwoMessagesWithTokensSuccess_gas() public {
    vm.pauseGasMetering();
    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].header.messageId = Internal._hash(messages[1], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[1].header.sequenceNumber,
      messages[1].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    Internal.ExecutionReportSingleChain memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.resumeGasMetering();
    s_offRamp.executeSingleReport(report, new uint256[](0));
  }

  function test_TwoMessagesWithTokensAndGE_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    // Set message 1 to use another receiver to simulate more fair gas costs
    messages[1].receiver = address(s_secondary_receiver);
    messages[1].header.messageId = Internal._hash(messages[1], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[1].header.sequenceNumber,
      messages[1].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    assertEq(uint64(0), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
    assertEq(uint64(2), s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER)));
  }

  function test_Fuzz_InterleavingOrderedAndUnorderedMessages_Success(bool[7] memory orderings) public {
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
      messages[i].header.messageId = Internal._hash(messages[i], ON_RAMP_ADDRESS_1);

      vm.expectEmit();
      emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
        SOURCE_CHAIN_SELECTOR_1,
        messages[i].header.sequenceNumber,
        messages[i].header.messageId,
        Internal.MessageExecutionState.SUCCESS,
        ""
      );
    }

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, abi.encode(OWNER));
    assertEq(uint64(0), nonceBefore, "nonce before exec should be 0");
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
    // all executions should succeed.
    for (uint256 i = 0; i < orderings.length; ++i) {
      assertEq(
        uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, messages[i].header.sequenceNumber)),
        uint256(Internal.MessageExecutionState.SUCCESS)
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
    messages[0].sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(fakePoolAddress),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[messages[0].tokenAmounts[0].token]),
        extraData: ""
      })
    );

    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);
    messages[1].header.messageId = Internal._hash(messages[1], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.TokenHandlingError.selector,
        abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, abi.encode(fakePoolAddress))
      )
    );

    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function test_WithCurseOnAnotherSourceChain_Success() public {
    s_mockRMN.setChainCursed(SOURCE_CHAIN_SELECTOR_2, true);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[](0)
    );
  }

  // Reverts

  function test_MismatchingDestChainSelector_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3);
    messages[0].header.destChainSelector = DEST_CHAIN_SELECTOR + 1;

    Internal.ExecutionReportSingleChain memory executionReport =
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.InvalidMessageDestChainSelector.selector, messages[0].header.destChainSelector
      )
    );
    s_offRamp.executeSingleReport(executionReport, new uint256[](0));
  }

  function test_MismatchingOnRampRoot_Revert() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 0);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    EVM2EVMMultiOffRamp.CommitReport memory commitReport = _constructCommitReport(
      // Root against mismatching on ramp
      Internal._hash(messages[0], ON_RAMP_ADDRESS_3)
    );
    _commit(commitReport, s_latestSequenceNumber);

    Internal.ExecutionReportSingleChain memory executionReport =
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.RootNotCommitted.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.executeSingleReport(executionReport, new uint256[](0));
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.CursedByRMN.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[](0)
    );
    // Uncurse should succeed
    s_mockRMN.setGlobalCursed(false);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[](0)
    );
  }

  function test_UnhealthySingleChainCurse_Revert() public {
    s_mockRMN.setChainCursed(SOURCE_CHAIN_SELECTOR_1, true);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.CursedByRMN.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[](0)
    );
    // Uncurse should succeed
    s_mockRMN.setChainCursed(SOURCE_CHAIN_SELECTOR_1, false);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[](0)
    );
  }

  function test_UnexpectedTokenData_Revert() public {
    Internal.ExecutionReportSingleChain memory report = _generateReportFromMessages(
      SOURCE_CHAIN_SELECTOR_1, _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
    );
    report.offchainTokenData = new bytes[][](report.messages.length + 1);

    vm.expectRevert(EVM2EVMMultiOffRamp.UnexpectedTokenData.selector);

    s_offRamp.executeSingleReport(report, new uint256[](0));
  }

  function test_EmptyReport_Revert() public {
    vm.expectRevert(EVM2EVMMultiOffRamp.EmptyReport.selector);
    s_offRamp.executeSingleReport(
      Internal.ExecutionReportSingleChain({
        sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
        proofs: new bytes32[](0),
        proofFlagBits: 0,
        messages: new Internal.Any2EVMRampMessage[](0),
        offchainTokenData: new bytes[][](0)
      }),
      new uint256[](0)
    );
  }

  function test_RootNotCommitted_Revert() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 0);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.RootNotCommitted.selector, SOURCE_CHAIN_SELECTOR_1));

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    s_offRamp.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), _getGasLimitsFromMessages(messages)
    );
  }

  function test_ManualExecutionNotYetEnabled_Revert() public {
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, BLOCK_TIME);

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOffRamp.ManualExecutionNotYetEnabled.selector, SOURCE_CHAIN_SELECTOR_1)
    );

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

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.SourceChainNotEnabled.selector, newSourceChainSelector));
    s_offRamp.executeSingleReport(_generateReportFromMessages(newSourceChainSelector, messages), new uint256[](0));
  }

  function test_DisabledSourceChain_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_2, ON_RAMP_ADDRESS_2);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.SourceChainNotEnabled.selector, SOURCE_CHAIN_SELECTOR_2));
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_2, messages), new uint256[](0));
  }

  function test_TokenDataMismatch_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReportSingleChain memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    report.offchainTokenData[0] = new bytes[](messages[0].tokenAmounts.length + 1);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.TokenDataMismatch.selector, SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber
      )
    );
    s_offRamp.executeSingleReport(report, new uint256[](0));
  }

  function test_RouterYULCall_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    // gas limit too high, Router's external call should revert
    messages[0].gasLimit = 1e36;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    Internal.ExecutionReportSingleChain memory executionReport =
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ExecutionError.selector,
        messages[0].header.messageId,
        abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector)
      )
    );
    s_offRamp.executeSingleReport(executionReport, new uint256[](0));
  }

  function test_RetryFailedMessageWithoutManualExecution_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1)
      )
    );
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.AlreadyAttempted.selector, SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber
      )
    );
    s_offRamp.executeSingleReport(_generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[](0));
  }

  function _constructCommitReport(bytes32 merkleRoot) internal view returns (EVM2EVMMultiOffRamp.CommitReport memory) {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: merkleRoot
    });

    return EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });
  }
}

contract EVM2EVMMultiOffRamp_executeSingleMessage is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    vm.startPrank(address(s_offRamp));
  }

  function test_executeSingleMessage_NoTokens_Success() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function test_executeSingleMessage_WithTokens_Success() public {
    Internal.Any2EVMRampMessage memory message =
      _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)[0];
    bytes[] memory offchainTokenData = new bytes[](message.tokenAmounts.length);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(message.sourceTokenData[0], (Internal.SourceTokenData));

    vm.expectCall(
      s_destPoolByToken[s_destTokens[0]],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: message.sender,
          receiver: message.receiver,
          amount: message.tokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[message.tokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: ""
        })
      )
    );

    s_offRamp.executeSingleMessage(message, offchainTokenData);
  }

  function test_executeSingleMessage_WithValidation_Success() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageValidator();
    vm.startPrank(address(s_offRamp));
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function test_NonContract_Success() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
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
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
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

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.TokenHandlingError.selector, errorMessage));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function test_ZeroGasDONExecution_Revert() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    message.gasLimit = 0;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.ReceiverError.selector, ""));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function test_MessageSender_Revert() public {
    vm.stopPrank();
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    vm.expectRevert(EVM2EVMMultiOffRamp.CanOnlySelfCall.selector);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function test_executeSingleMessage_WithFailingValidation_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageValidator();
    vm.startPrank(address(s_offRamp));
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_inboundMessageValidator.setMessageIdValidationState(message.header.messageId, true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }

  function test_executeSingleMessage_WithFailingValidationNoRouterCall_Revert() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageValidator();
    vm.startPrank(address(s_offRamp));

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);

    // Setup the receiver to a non-CCIP Receiver, which will skip the Router call (but should still perform the validation)
    MaybeRevertMessageReceiverNo165 newReceiver = new MaybeRevertMessageReceiverNo165(true);
    message.receiver = address(newReceiver);
    message.header.messageId = Internal._hash(message, ON_RAMP_ADDRESS_1);

    s_inboundMessageValidator.setMessageIdValidationState(message.header.messageId, true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length));
  }
}

contract EVM2EVMMultiOffRamp_batchExecute is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_3, 1);
  }

  function test_SingleReport_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender);
    s_offRamp.batchExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[][](1));

    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender), nonceBefore);
  }

  function test_MultipleReportsSameChain_Success() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint64 nonceBefore = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages1[0].sender);
    s_offRamp.batchExecute(reports, new uint256[][](2));
    assertGt(s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages1[0].sender), nonceBefore);
  }

  function test_MultipleReportsDifferentChains_Success() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, 1);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.batchExecute(reports, new uint256[][](2));

    uint64 nonceChain1 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages1[0].sender);
    uint64 nonceChain3 = s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_3, messages2[0].sender);

    assertTrue(nonceChain1 != nonceChain3);
    assertGt(nonceChain1, 0);
    assertGt(nonceChain3, 0);
  }

  function test_MultipleReportsSkipDuplicate_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SkippedAlreadyExecutedMessage(SOURCE_CHAIN_SELECTOR_1, messages[0].header.sequenceNumber);

    s_offRamp.batchExecute(reports, new uint256[][](2));
  }

  // Reverts
  function test_ZeroReports_Revert() public {
    vm.expectRevert(EVM2EVMMultiOffRamp.EmptyReport.selector);
    s_offRamp.batchExecute(new Internal.ExecutionReportSingleChain[](0), new uint256[][](1));
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.CursedByRMN.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[][](1)
    );
    // Uncurse should succeed
    s_mockRMN.setGlobalCursed(false);
    s_offRamp.batchExecute(
      _generateBatchReportFromMessages(
        SOURCE_CHAIN_SELECTOR_1, _generateMessagesWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1)
      ),
      new uint256[][](1)
    );
  }

  function test_OutOfBoundsGasLimitsAccess_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    vm.expectRevert();
    s_offRamp.batchExecute(reports, new uint256[][](1));
  }
}

contract EVM2EVMMultiOffRamp_manuallyExecute is EVM2EVMMultiOffRampSetup {
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
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);
    s_offRamp.batchExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[][](1));

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = new uint256[](messages.length);
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_manuallyExecute_WithGasOverride_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);
    s_offRamp.batchExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[][](1));

    s_reverting_receiver.setRevert(false);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0][0] += 1;

    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_manuallyExecute_DoesNotRevertIfUntouched_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    assertEq(
      messages[0].header.nonce - 1, s_inboundNonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR_1, messages[0].sender)
    );

    s_reverting_receiver.setRevert(true);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, "")
      )
    );

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);

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
      messages1[i].header.messageId = Internal._hash(messages1[i], ON_RAMP_ADDRESS_1);
    }

    for (uint64 i = 0; i < 2; ++i) {
      messages2[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, i + 1);
      messages2[i].receiver = address(s_reverting_receiver);
      messages2[i].header.messageId = Internal._hash(messages2[i], ON_RAMP_ADDRESS_3);
    }

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    s_offRamp.batchExecute(reports, new uint256[][](2));

    s_reverting_receiver.setRevert(false);

    uint256[][] memory gasLimitOverrides = new uint256[][](2);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages1);
    gasLimitOverrides[1] = _getGasLimitsFromMessages(messages2);

    for (uint256 i = 0; i < 3; ++i) {
      vm.expectEmit();
      emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
        SOURCE_CHAIN_SELECTOR_1,
        messages1[i].header.sequenceNumber,
        messages1[i].header.messageId,
        Internal.MessageExecutionState.SUCCESS,
        ""
      );

      gasLimitOverrides[0][i] += 1;
    }

    for (uint256 i = 0; i < 2; ++i) {
      vm.expectEmit();
      emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
        SOURCE_CHAIN_SELECTOR_3,
        messages2[i].header.sequenceNumber,
        messages2[i].header.messageId,
        Internal.MessageExecutionState.SUCCESS,
        ""
      );

      gasLimitOverrides[1][i] += 1;
    }

    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_manuallyExecute_WithPartialMessages_Success() public {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](3);

    for (uint64 i = 0; i < 3; ++i) {
      messages[i] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, i + 1);
    }
    messages[1].receiver = address(s_reverting_receiver);
    messages[1].header.messageId = Internal._hash(messages[1], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[1].header.sequenceNumber,
      messages[1].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, bytes(""))
      )
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[2].header.sequenceNumber,
      messages[2].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.batchExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[][](1));

    s_reverting_receiver.setRevert(false);

    // Only the 2nd message reverted
    Internal.Any2EVMRampMessage[] memory newMessages = new Internal.Any2EVMRampMessage[](1);
    newMessages[0] = messages[1];

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(newMessages);
    gasLimitOverrides[0][0] += 1;

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      newMessages[0].header.sequenceNumber,
      newMessages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, newMessages), gasLimitOverrides);
  }

  function test_manuallyExecute_LowGasLimit_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    messages[0].gasLimit = 1;
    messages[0].receiver = address(new ConformingReceiver(address(s_destRouter), s_destFeeToken));
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(EVM2EVMMultiOffRamp.ReceiverError.selector, "")
    );
    s_offRamp.batchExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[][](1));

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = new uint256[](1);
    gasLimitOverrides[0][0] = 100_000;

    vm.expectEmit();
    emit ConformingReceiver.MessageReceived();

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  // Reverts

  function test_manuallyExecute_ForkedChain_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ForkedChain.selector, chain1, chain2));

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_ManualExecGasLimitMismatchSingleReport_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](2);
    messages[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);

    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    // No overrides for report
    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, new uint256[][](0));

    // No messages
    uint256[][] memory gasLimitOverrides = new uint256[][](1);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1 message missing
    gasLimitOverrides[0] = new uint256[](1);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1 message in excess
    gasLimitOverrides[0] = new uint256[](3);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_manuallyExecute_GasLimitMismatchMultipleReports_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_3, ON_RAMP_ADDRESS_3, 1);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_3, messages2);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, new uint256[][](0));

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, new uint256[][](1));

    uint256[][] memory gasLimitOverrides = new uint256[][](2);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 2nd report empty
    gasLimitOverrides[0] = new uint256[](2);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1st report empty
    gasLimitOverrides[0] = new uint256[](0);
    gasLimitOverrides[1] = new uint256[](1);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);

    // 1st report oversized
    gasLimitOverrides[0] = new uint256[](3);

    vm.expectRevert(EVM2EVMMultiOffRamp.ManualExecutionGasLimitMismatch.selector);
    s_offRamp.manuallyExecute(reports, gasLimitOverrides);
  }

  function test_ManualExecInvalidGasLimit_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);
    gasLimitOverrides[0][0]--;

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.InvalidManualExecutionGasLimit.selector, SOURCE_CHAIN_SELECTOR_1, 0, gasLimitOverrides[0][0]
      )
    );
    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);
  }

  function test_manuallyExecute_FailedTx_Revert() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);

    messages[0].receiver = address(s_reverting_receiver);
    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    s_offRamp.batchExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new uint256[][](1));

    s_reverting_receiver.setRevert(true);

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ExecutionError.selector,
        messages[0].header.messageId,
        abi.encodeWithSelector(
          EVM2EVMMultiOffRamp.ReceiverError.selector,
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
    messages[0].tokenAmounts = new Client.EVMTokenAmount[](1);
    messages[0].tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: tokenAmount});
    messages[0].sourceTokenData = new bytes[](1);
    messages[0].sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[s_sourceFeeToken]),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[s_sourceFeeToken]),
        extraData: ""
      })
    );

    messages[0].receiver = address(receiver);

    messages[0].header.messageId = Internal._hash(messages[0], ON_RAMP_ADDRESS_1);

    Internal.ExecutionReportSingleChain memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    // sets the report to be repeated on the ReentrancyAbuser to be able to replay
    receiver.setPayload(report);

    uint256[][] memory gasLimitOverrides = new uint256[][](1);
    gasLimitOverrides[0] = _getGasLimitsFromMessages(messages);

    // The first entry should be fine and triggers the second entry. This one fails
    // but since it's an inner tx of the first one it is caught in the try-catch.
    // This means the first tx is marked `FAILURE` with the error message of the second tx.
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages[0].header.sourceChainSelector,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.ReceiverError.selector,
        abi.encodeWithSelector(
          EVM2EVMMultiOffRamp.AlreadyExecuted.selector,
          messages[0].header.sourceChainSelector,
          messages[0].header.sequenceNumber
        )
      )
    );

    s_offRamp.manuallyExecute(_generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), gasLimitOverrides);

    // Since the tx failed we don't release the tokens
    assertEq(tokenToAbuse.balanceOf(address(receiver)), balancePre);
  }
}

contract EVM2EVMMultiOffRamp_execute is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);
  }

  // Asserts that execute completes
  function test_SingleReport_Success() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR_1,
      messages[0].header.sequenceNumber,
      messages[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    _execute(reports);
  }

  function test_MultipleReports_Success() public {
    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    _execute(reports);
  }

  function test_LargeBatch_Success() public {
    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](10);
    for (uint64 i = 0; i < reports.length; ++i) {
      Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](3);
      messages[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1 + i * 3);
      messages[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2 + i * 3);
      messages[2] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3 + i * 3);

      reports[i] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);
    }

    for (uint64 i = 0; i < reports.length; ++i) {
      for (uint64 j = 0; j < reports[i].messages.length; ++j) {
        vm.expectEmit();
        emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
          reports[i].messages[j].header.sourceChainSelector,
          reports[i].messages[j].header.sequenceNumber,
          reports[i].messages[j].header.messageId,
          Internal.MessageExecutionState.SUCCESS,
          ""
        );
      }
    }

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    _execute(reports);
  }

  function test_MultipleReportsWithPartialValidationFailures_Success() public {
    _enableInboundMessageValidator();

    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    messages1[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    messages1[1] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 2);
    messages2[0] = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 3);

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages2);

    s_inboundMessageValidator.setMessageIdValidationState(messages1[0].header.messageId, true);
    s_inboundMessageValidator.setMessageIdValidationState(messages2[0].header.messageId, true);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[0].header.sourceChainSelector,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages1[1].header.sourceChainSelector,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      messages2[0].header.sourceChainSelector,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      Internal.MessageExecutionState.FAILURE,
      abi.encodeWithSelector(
        IMessageInterceptor.MessageValidationError.selector,
        abi.encodeWithSelector(IMessageInterceptor.MessageValidationError.selector, bytes("Invalid message"))
      )
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(
      uint8(Internal.OCRPluginType.Execution), s_configDigestExec, uint64(uint256(s_configDigestExec))
    );

    _execute(reports);
  }

  // Reverts

  function test_UnauthorizedTransmitter_Revert() public {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function test_NoConfig_Revert() public {
    _redeployOffRampWithNoOCRConfigs();
    s_offRamp.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

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
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

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
    Internal.ExecutionReportSingleChain[] memory reports =
      _generateBatchReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.expectRevert();
    _execute(reports);
  }

  function test_ZeroReports_Revert() public {
    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](0);

    vm.expectRevert(EVM2EVMMultiOffRamp.EmptyReport.selector);
    _execute(reports);
  }

  function test_IncorrectArrayType_Revert() public {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    uint256[] memory wrongData = new uint256[](1);
    wrongData[0] = 1;

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.execute(reportContext, abi.encode(wrongData));
  }

  function test_NonArray_Revert() public {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    Internal.ExecutionReportSingleChain memory report = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.execute(reportContext, abi.encode(report));
  }
}

contract EVM2EVMMultiOffRamp_getExecutionState is EVM2EVMMultiOffRampSetup {
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

contract EVM2EVMMultiOffRamp_trialExecute is EVM2EVMMultiOffRampSetup {
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
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length));
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
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMMultiOffRamp.TokenHandlingError.selector, errorMessage), err);

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
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMMultiOffRamp.TokenHandlingError.selector, errorMessage), err);
  }

  // TODO test actual pool exists but isn't compatible instead of just no pool
  function test_TokenPoolIsNotAContract_Success() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 10000;
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);

    // Happy path, pool is correct
    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length));

    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // address 0 has no contract
    assertEq(address(0).code.length, 0);
    message.sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(address(0)),
        destTokenAddress: abi.encode(address(0)),
        extraData: ""
      })
    );

    message.header.messageId = Internal._hash(message, ON_RAMP_ADDRESS_1);

    // Unhappy path, no revert but marked as failed.
    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, abi.encode(address(0))), err);

    address notAContract = makeAddr("not_a_contract");

    message.sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(address(0)),
        destTokenAddress: abi.encode(notAContract),
        extraData: ""
      })
    );

    message.header.messageId = Internal._hash(message, ON_RAMP_ADDRESS_1);

    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(EVM2EVMMultiOffRamp.NotACompatiblePool.selector, address(0)), err);
  }
}

contract EVM2EVMMultiOffRamp__releaseOrMintSingleToken is EVM2EVMMultiOffRampSetup {
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

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(s_destTokenBySourceToken[token]),
      extraData: ""
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
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData
        })
      )
    );

    s_offRamp.releaseOrMintSingleToken(
      amount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, offchainTokenData
    );

    assertEq(startingBalance + amount, dstToken1.balanceOf(OWNER));
  }

  function test__releaseOrMintSingleToken_NotACompatiblePool_Revert() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    vm.label(destToken, "destToken");
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(destToken),
      extraData: ""
    });

    // Address(0) should always revert
    address returnedPool = address(0);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintSingleToken(
      amount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, offchainTokenData
    );

    // A contract that doesn't support the interface should also revert
    returnedPool = address(s_offRamp);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintSingleToken(
      amount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, offchainTokenData
    );
  }

  function test__releaseOrMintSingleToken_TokenHandlingError_revert_Revert() public {
    address receiver = makeAddr("receiver");
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: abi.encode(destToken),
      extraData: ""
    });

    bytes memory revertData = "call reverted :o";

    vm.mockCallRevert(destToken, abi.encodeWithSelector(IERC20.transfer.selector, receiver, amount), revertData);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.TokenHandlingError.selector, revertData));
    s_offRamp.releaseOrMintSingleToken(
      amount, originalSender, receiver, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, offchainTokenData
    );
  }
}

contract EVM2EVMMultiOffRamp_releaseOrMintTokens is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test_releaseOrMintTokens_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destFeeToken);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(encodedSourceTokenData[0], (Internal.SourceTokenData));

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
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, encodedSourceTokenData, offchainTokenData
    );

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintTokens_destDenominatedDecimals_Success() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    address destToken = s_destFeeToken;
    uint256 amount = 100;
    uint256 destinationDenominationMultiplier = 1000;
    srcTokenAmounts[0].amount = amount;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(encodedSourceTokenData[0], (Internal.SourceTokenData));

    // Since the pool call is mocked, we manually release funds to the offRamp
    deal(destToken, address(s_offRamp), amount * destinationDenominationMultiplier);

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
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData[0]
        })
      ),
      abi.encode(amount * destinationDenominationMultiplier)
    );

    Client.EVMTokenAmount[] memory destTokenAmounts = s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, encodedSourceTokenData, offchainTokenData
    );

    assertEq(destTokenAmounts[0].amount, amount * destinationDenominationMultiplier);
    assertEq(destTokenAmounts[0].token, destToken);
  }

  // Revert

  function test_TokenHandlingError_Reverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes memory unknownError = bytes("unknown error");
    s_maybeRevertingPool.setShouldRevert(unknownError);

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.TokenHandlingError.selector, unknownError));

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts,
      abi.encode(OWNER),
      OWNER,
      SOURCE_CHAIN_SELECTOR_1,
      _getDefaultSourceTokenData(srcTokenAmounts),
      new bytes[](srcTokenAmounts.length)
    );
  }

  function test_releaseOrMintTokens_InvalidDataLengthReturnData_Revert() public {
    uint256 amount = 100;
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    srcTokenAmounts[0].amount = amount;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(encodedSourceTokenData[0], (Internal.SourceTokenData));

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
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData[0]
        })
      ),
      // Includes the amount twice, this will revert due to the return data being to long
      abi.encode(amount, amount)
    );

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOffRamp.InvalidDataLength.selector, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, 64)
    );

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, encodedSourceTokenData, offchainTokenData
    );
  }

  function test_releaseOrMintTokens_InvalidEVMAddress_Revert() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    bytes[] memory sourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    bytes memory wrongAddress = abi.encode(address(1000), address(10000), address(10000));

    sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[srcTokenAmounts[0].token]),
        destTokenAddress: wrongAddress,
        extraData: ""
      })
    );

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, wrongAddress));

    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, offchainTokenData
    );
  }

  function test__releaseOrMintTokens_PoolIsNotAPool_Reverts() public {
    // The offRamp is a contract, but not a pool
    address fakePoolAddress = address(s_offRamp);

    bytes[] memory sourceTokenData = new bytes[](1);
    sourceTokenData[0] = abi.encode(
      Internal.SourceTokenData({
        sourcePoolAddress: abi.encode(fakePoolAddress),
        destTokenAddress: abi.encode(s_offRamp),
        extraData: ""
      })
    );

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.NotACompatiblePool.selector, address(0)));
    s_offRamp.releaseOrMintTokens(
      new Client.EVMTokenAmount[](1), abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, new bytes[](1)
    );
  }

  function test_releaseOrMintTokens_PoolDoesNotSupportDest_Reverts() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    bytes[] memory encodedSourceTokenData = _getDefaultSourceTokenData(srcTokenAmounts);
    Internal.SourceTokenData memory sourceTokenData = abi.decode(encodedSourceTokenData[0], (Internal.SourceTokenData));

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
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );
    vm.expectRevert();
    s_offRamp.releaseOrMintTokens(
      srcTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_3, encodedSourceTokenData, offchainTokenData
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
        extraData: unusedVar
      })
    );

    try s_offRamp.releaseOrMintTokens(
      new Client.EVMTokenAmount[](1), abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, sourceTokenData, new bytes[](1)
    ) {} catch (bytes memory reason) {
      // Any revert should be a TokenHandlingError, InvalidEVMAddress, InvalidDataLength or NoContract as those are caught by the offramp
      assertTrue(
        bytes4(reason) == EVM2EVMMultiOffRamp.TokenHandlingError.selector
          || bytes4(reason) == Internal.InvalidEVMAddress.selector
          || bytes4(reason) == EVM2EVMMultiOffRamp.InvalidDataLength.selector
          || bytes4(reason) == CallWithExactGas.NoContract.selector
          || bytes4(reason) == EVM2EVMMultiOffRamp.NotACompatiblePool.selector,
        "Expected TokenHandlingError or InvalidEVMAddress"
      );

      if (destPool > type(uint160).max) {
        assertEq(reason, abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, abi.encode(destPool)));
      }
    }
  }
}

contract EVM2EVMMultiOffRamp_applySourceChainConfigUpdates is EVM2EVMMultiOffRampSetup {
  function test_ApplyZeroUpdates_Success() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](0);

    vm.recordLogs();
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    // No logs emitted
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);

    // assertEq(s_offRamp.getSourceChainSelectors().length, 0);
  }

  function test_AddNewChain_Success() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    EVM2EVMMultiOffRamp.SourceChainConfig memory expectedSourceChainConfig =
      EVM2EVMMultiOffRamp.SourceChainConfig({isEnabled: true, minSeqNr: 1, onRamp: ON_RAMP_ADDRESS_1});

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1);

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig);

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    _assertSourceChainConfigEquality(s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig);
  }

  function test_ReplaceExistingChain_Success() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].isEnabled = false;
    EVM2EVMMultiOffRamp.SourceChainConfig memory expectedSourceChainConfig =
      EVM2EVMMultiOffRamp.SourceChainConfig({isEnabled: false, minSeqNr: 1, onRamp: ON_RAMP_ADDRESS_1});

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig);

    vm.recordLogs();
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    // No log emitted for chain selector added (only for setting the config)
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    _assertSourceChainConfigEquality(s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig);

    // uint64[] memory resultSourceChainSelectors = s_offRamp.getSourceChainSelectors();
    // assertEq(resultSourceChainSelectors.length, 1);
    // assertEq(resultSourceChainSelectors[0], SOURCE_CHAIN_SELECTOR_1);
  }

  function test_AddMultipleChains_Success() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 0),
      isEnabled: true
    });
    sourceChainConfigs[1] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 1,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 1),
      isEnabled: false
    });
    sourceChainConfigs[2] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 2,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 2),
      isEnabled: true
    });

    EVM2EVMMultiOffRamp.SourceChainConfig[] memory expectedSourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfig[](3);
    for (uint256 i = 0; i < 3; ++i) {
      expectedSourceChainConfigs[i] = EVM2EVMMultiOffRamp.SourceChainConfig({
        isEnabled: sourceChainConfigs[i].isEnabled,
        minSeqNr: 1,
        onRamp: abi.encode(ON_RAMP_ADDRESS_1, i)
      });

      vm.expectEmit();
      emit EVM2EVMMultiOffRamp.SourceChainSelectorAdded(sourceChainConfigs[i].sourceChainSelector);

      vm.expectEmit();
      emit EVM2EVMMultiOffRamp.SourceChainConfigSet(
        sourceChainConfigs[i].sourceChainSelector, expectedSourceChainConfigs[i]
      );
    }

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    for (uint256 i = 0; i < 3; ++i) {
      _assertSourceChainConfigEquality(
        s_offRamp.getSourceChainConfig(sourceChainConfigs[i].sourceChainSelector), expectedSourceChainConfigs[i]
      );
    }
  }

  function test_Fuzz_applySourceChainConfigUpdate_Success(
    EVM2EVMMultiOffRamp.SourceChainConfigArgs memory sourceChainConfigArgs
  ) public {
    // Skip invalid inputs
    vm.assume(sourceChainConfigArgs.sourceChainSelector != 0);
    vm.assume(sourceChainConfigArgs.onRamp.length != 0);

    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
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

    EVM2EVMMultiOffRamp.SourceChainConfig memory expectedSourceChainConfig = EVM2EVMMultiOffRamp.SourceChainConfig({
      isEnabled: sourceChainConfigArgs.isEnabled,
      minSeqNr: 1,
      onRamp: sourceChainConfigArgs.onRamp
    });

    if (isNewChain) {
      vm.expectEmit();
      emit EVM2EVMMultiOffRamp.SourceChainSelectorAdded(sourceChainConfigArgs.sourceChainSelector);
    }

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.SourceChainConfigSet(sourceChainConfigArgs.sourceChainSelector, expectedSourceChainConfig);

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    _assertSourceChainConfigEquality(
      s_offRamp.getSourceChainConfig(sourceChainConfigArgs.sourceChainSelector), expectedSourceChainConfig
    );
  }

  // Reverts

  function test_ZeroOnRampAddress_Revert() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: new bytes(0),
      isEnabled: true
    });

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_ZeroSourceChainSelector_Revert() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] =
      EVM2EVMMultiOffRamp.SourceChainConfigArgs({sourceChainSelector: 0, onRamp: ON_RAMP_ADDRESS_1, isEnabled: true});

    vm.expectRevert(EVM2EVMMultiOffRamp.ZeroChainSelectorNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_ReplaceExistingChainOnRamp_Revert() public {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].onRamp = ON_RAMP_ADDRESS_2;

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.InvalidStaticConfig.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }
}

contract EVM2EVMMultiOffRamp_commit is EVM2EVMMultiOffRampSetup {
  uint64 internal s_maxInterval = 12;

  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();

    s_latestSequenceNumber = uint64(uint256(s_configDigestCommit));
  }

  function test_ReportAndPriceUpdate_Success() public {
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = _constructCommitReport();

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.CommitReportAccepted(commitReport);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_maxInterval + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_ReportOnlyRootSuccess_gas() public {
    uint64 max1 = 931;
    bytes32 root = "Only a single root";

    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, max1),
      merkleRoot: root
    });

    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.CommitReportAccepted(commitReport);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(max1 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(block.timestamp, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR_1, root));
  }

  function test_StaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenStartPrice =
      IPriceRegistry(s_offRamp.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value;

    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, maxSeq),
      merkleRoot: "stale report 1"
    });
    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.CommitReportAccepted(commitReport);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    commitReport.merkleRoots[0].interval = EVM2EVMMultiOffRamp.Interval(maxSeq + 1, maxSeq * 2);
    commitReport.merkleRoots[0].merkleRoot = "stale report 2";

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.CommitReportAccepted(commitReport);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq * 2 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(
      tokenStartPrice, IPriceRegistry(s_offRamp.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
  }

  function test_OnlyTokenPriceUpdates_Success() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](0);
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });

    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_OnlyGasPriceUpdates_Success() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](0);
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });

    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_PriceSequenceNumberCleared_Success() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](0);
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });

    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
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
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_ValidPriceUpdateThenStaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenPrice1 = 4e18;
    uint224 tokenPrice2 = 5e18;
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](0);
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice1),
      merkleRoots: roots
    });

    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, tokenPrice1, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, maxSeq),
      merkleRoot: "stale report"
    });
    commitReport.priceUpdates = getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice2);
    commitReport.merkleRoots = roots;

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.CommitReportAccepted(commitReport);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(
      tokenPrice1, IPriceRegistry(s_offRamp.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  // Reverts

  function test_UnauthorizedTransmitter_Revert() public {
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[3] memory reportContext =
      [s_configDigestCommit, bytes32(uint256(s_latestSequenceNumber)), s_configDigestCommit];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, s_F + 1);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_NoConfig_Revert() public {
    _redeployOffRampWithNoOCRConfigs();

    EVM2EVMMultiOffRamp.CommitReport memory commitReport = _constructCommitReport();

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

    EVM2EVMMultiOffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[3] memory reportContext = [bytes32(""), s_configDigestCommit, s_configDigestCommit];
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, s_F + 1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_WrongConfigWithoutSigners_Revert() public {
    _redeployOffRampWithNoOCRConfigs();

    EVM2EVMMultiOffRamp.CommitReport memory commitReport = _constructCommitReport();

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    vm.expectRevert();
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: "Only a single root"
    });

    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.CursedByRMN.selector, roots[0].sourceChainSelector));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_InvalidRootRevert() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, 4),
      merkleRoot: bytes32(0)
    });
    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectRevert(EVM2EVMMultiOffRamp.InvalidRoot.selector);
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_InvalidInterval_Revert() public {
    EVM2EVMMultiOffRamp.Interval memory interval = EVM2EVMMultiOffRamp.Interval(2, 2);
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: interval,
      merkleRoot: bytes32(0)
    });
    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOffRamp.InvalidInterval.selector, roots[0].sourceChainSelector, interval)
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_InvalidIntervalMinLargerThanMax_Revert() public {
    s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR);
    EVM2EVMMultiOffRamp.Interval memory interval = EVM2EVMMultiOffRamp.Interval(1, 0);
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: interval,
      merkleRoot: bytes32(0)
    });
    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectRevert(
      abi.encodeWithSelector(EVM2EVMMultiOffRamp.InvalidInterval.selector, roots[0].sourceChainSelector, interval)
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_ZeroEpochAndRound_Revert() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](0);
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });

    vm.expectRevert(EVM2EVMMultiOffRamp.StaleCommitReport.selector);
    _commit(commitReport, 0);
  }

  function test_OnlyPriceUpdateStaleReport_Revert() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](0);
    EVM2EVMMultiOffRamp.CommitReport memory commitReport = EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });

    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    _commit(commitReport, s_latestSequenceNumber);

    vm.expectRevert(EVM2EVMMultiOffRamp.StaleCommitReport.selector);
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_SourceChainNotEnabled_Revert() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: 0,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: "Only a single root"
    });

    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectRevert(abi.encodeWithSelector(EVM2EVMMultiOffRamp.SourceChainNotEnabled.selector, 0));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RootAlreadyCommitted_Revert() public {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: "Only a single root"
    });
    EVM2EVMMultiOffRamp.CommitReport memory commitReport =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    _commit(commitReport, s_latestSequenceNumber);
    commitReport.merkleRoots[0].interval = EVM2EVMMultiOffRamp.Interval(3, 3);

    vm.expectRevert(
      abi.encodeWithSelector(
        EVM2EVMMultiOffRamp.RootAlreadyCommitted.selector, roots[0].sourceChainSelector, roots[0].merkleRoot
      )
    );
    _commit(commitReport, ++s_latestSequenceNumber);
  }

  function _constructCommitReport() internal view returns (EVM2EVMMultiOffRamp.CommitReport memory) {
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      interval: EVM2EVMMultiOffRamp.Interval(1, s_maxInterval),
      merkleRoot: "test #2"
    });

    return EVM2EVMMultiOffRamp.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });
  }
}

contract EVM2EVMMultiOffRamp_resetUnblessedRoots is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupRealRMN();
    _deployOffRamp(s_destRouter, s_realRMN, s_inboundNonceManager);
    _setupMultipleOffRamps();
  }

  function test_ResetUnblessedRoots_Success() public {
    EVM2EVMMultiOffRamp.UnblessedRoot[] memory rootsToReset = new EVM2EVMMultiOffRamp.UnblessedRoot[](3);
    rootsToReset[0] = EVM2EVMMultiOffRamp.UnblessedRoot({sourceChainSelector: SOURCE_CHAIN_SELECTOR, merkleRoot: "1"});
    rootsToReset[1] = EVM2EVMMultiOffRamp.UnblessedRoot({sourceChainSelector: SOURCE_CHAIN_SELECTOR, merkleRoot: "2"});
    rootsToReset[2] = EVM2EVMMultiOffRamp.UnblessedRoot({sourceChainSelector: SOURCE_CHAIN_SELECTOR, merkleRoot: "3"});

    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](3);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: rootsToReset[0].merkleRoot
    });
    roots[1] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(3, 4),
      merkleRoot: rootsToReset[1].merkleRoot
    });
    roots[2] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(5, 5),
      merkleRoot: rootsToReset[2].merkleRoot
    });

    EVM2EVMMultiOffRamp.CommitReport memory report =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    _commit(report, ++s_latestSequenceNumber);

    IRMN.TaggedRoot[] memory blessedTaggedRoots = new IRMN.TaggedRoot[](1);
    blessedTaggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_offRamp), root: rootsToReset[1].merkleRoot});

    vm.startPrank(BLESS_VOTE_ADDR);
    s_realRMN.voteToBless(blessedTaggedRoots);

    vm.expectEmit(false, false, false, true);
    emit EVM2EVMMultiOffRamp.RootRemoved(rootsToReset[0].merkleRoot);

    vm.expectEmit(false, false, false, true);
    emit EVM2EVMMultiOffRamp.RootRemoved(rootsToReset[2].merkleRoot);

    vm.startPrank(OWNER);
    s_offRamp.resetUnblessedRoots(rootsToReset);

    assertEq(0, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR, rootsToReset[0].merkleRoot));
    assertEq(BLOCK_TIME, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR, rootsToReset[1].merkleRoot));
    assertEq(0, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR, rootsToReset[2].merkleRoot));
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    EVM2EVMMultiOffRamp.UnblessedRoot[] memory rootsToReset = new EVM2EVMMultiOffRamp.UnblessedRoot[](0);
    s_offRamp.resetUnblessedRoots(rootsToReset);
  }
}

contract EVM2EVMMultiOffRamp_verify is EVM2EVMMultiOffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupRealRMN();
    _deployOffRamp(s_destRouter, s_realRMN, s_inboundNonceManager);
    _setupMultipleOffRamps();
  }

  function test_NotBlessed_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";

    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: leaves[0]
    });
    EVM2EVMMultiOffRamp.CommitReport memory report =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    _commit(report, ++s_latestSequenceNumber);
    bytes32[] memory proofs = new bytes32[](0);
    // We have not blessed this root, should return 0.
    uint256 timestamp = s_offRamp.verify(SOURCE_CHAIN_SELECTOR, leaves, proofs, 0);
    assertEq(uint256(0), timestamp);
  }

  function test_Blessed_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: leaves[0]
    });
    EVM2EVMMultiOffRamp.CommitReport memory report =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    _commit(report, ++s_latestSequenceNumber);
    // Bless that root.
    IRMN.TaggedRoot[] memory taggedRoots = new IRMN.TaggedRoot[](1);
    taggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_offRamp), root: leaves[0]});
    vm.startPrank(BLESS_VOTE_ADDR);
    s_realRMN.voteToBless(taggedRoots);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_offRamp.verify(SOURCE_CHAIN_SELECTOR, leaves, proofs, 0);
    assertEq(BLOCK_TIME, timestamp);
  }

  function test_NotBlessedWrongChainSelector_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](1);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(1, 2),
      merkleRoot: leaves[0]
    });

    EVM2EVMMultiOffRamp.CommitReport memory report =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    _commit(report, ++s_latestSequenceNumber);

    // Bless that root.
    IRMN.TaggedRoot[] memory taggedRoots = new IRMN.TaggedRoot[](1);
    taggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_offRamp), root: leaves[0]});
    vm.startPrank(BLESS_VOTE_ADDR);
    s_realRMN.voteToBless(taggedRoots);

    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_offRamp.verify(SOURCE_CHAIN_SELECTOR + 1, leaves, proofs, 0);
    assertEq(uint256(0), timestamp);
  }

  // Reverts

  function test_TooManyLeaves_Revert() public {
    bytes32[] memory leaves = new bytes32[](258);
    bytes32[] memory proofs = new bytes32[](0);
    vm.expectRevert(MerkleMultiProof.InvalidProof.selector);
    s_offRamp.verify(SOURCE_CHAIN_SELECTOR, leaves, proofs, 0);
  }
}
