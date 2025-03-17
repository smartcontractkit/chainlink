// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";
import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {NonceManager} from "../../../NonceManager.sol";
import {Router} from "../../../Router.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {FeeQuoterSetup} from "../../feeQuoter/FeeQuoterSetup.t.sol";
import {MaybeRevertingBurnMintTokenPool} from "../../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MessageInterceptorHelper} from "../../helpers/MessageInterceptorHelper.sol";
import {OffRampHelper} from "../../helpers/OffRampHelper.sol";
import {MaybeRevertMessageReceiver} from "../../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MultiOCR3BaseSetup} from "../../ocr/MultiOCR3Base/MultiOCR3BaseSetup.t.sol";
import {Vm} from "forge-std/Test.sol";

contract OffRampSetup is FeeQuoterSetup, MultiOCR3BaseSetup {
  uint64 internal constant SOURCE_CHAIN_SELECTOR_1 = SOURCE_CHAIN_SELECTOR;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_2 = 6433500567565415381;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_3 = 4051577828743386545;

  address internal constant ON_RAMP_ADDRESS = 0x11118e64e1FB0c487f25dD6D3601FF6aF8d32E4e;

  bytes internal constant ON_RAMP_ADDRESS_1 = abi.encode(ON_RAMP_ADDRESS);
  bytes internal constant ON_RAMP_ADDRESS_2 = abi.encode(0xaA3f843Cf8E33B1F02dd28303b6bD87B1aBF8AE4);
  bytes internal constant ON_RAMP_ADDRESS_3 = abi.encode(0x71830C37Cb193e820de488Da111cfbFcC680a1b9);

  IAny2EVMMessageReceiver internal s_receiver;
  IAny2EVMMessageReceiver internal s_secondary_receiver;
  MaybeRevertMessageReceiver internal s_reverting_receiver;

  MaybeRevertingBurnMintTokenPool internal s_maybeRevertingPool;

  OffRampHelper internal s_offRamp;
  MessageInterceptorHelper internal s_inboundMessageInterceptor;
  NonceManager internal s_inboundNonceManager;

  bytes32 internal s_configDigestExec;
  bytes32 internal s_configDigestCommit;
  uint8 internal constant F = 1;

  uint64 internal s_latestSequenceNumber;

  IRMNRemote.Signature[] internal s_rmnSignatures;

  function setUp() public virtual override(FeeQuoterSetup, MultiOCR3BaseSetup) {
    FeeQuoterSetup.setUp();
    MultiOCR3BaseSetup.setUp();

    s_inboundMessageInterceptor = new MessageInterceptorHelper();
    s_receiver = new MaybeRevertMessageReceiver(false);
    s_secondary_receiver = new MaybeRevertMessageReceiver(false);
    s_reverting_receiver = new MaybeRevertMessageReceiver(true);

    s_maybeRevertingPool = MaybeRevertingBurnMintTokenPool(s_destPoolByToken[s_destTokens[1]]);
    s_inboundNonceManager = new NonceManager(new address[](0));

    _deployOffRamp(s_mockRMNRemote, s_inboundNonceManager);
  }

  function _deployOffRamp(IRMNRemote rmnRemote, NonceManager nonceManager) internal {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnRemote: rmnRemote,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(nonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );

    s_configDigestExec = _getBasicConfigDigest(F, s_emptySigners, s_validTransmitters);
    s_configDigestCommit = _getBasicConfigDigest(F, s_validSigners, s_validTransmitters);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](2);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[1] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    s_offRamp.setDynamicConfig(_generateDynamicOffRampConfig(address(s_feeQuoter)));
    s_offRamp.setOCR3Configs(ocrConfigs);

    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(s_offRamp);
    NonceManager(nonceManager).applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = address(s_offRamp);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: address(s_offRamp)});
    s_destRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);
  }

  function _setupMultipleOffRamps() internal {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_2,
      onRamp: ON_RAMP_ADDRESS_2,
      isEnabled: false,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[2] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_3,
      onRamp: ON_RAMP_ADDRESS_3,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });
    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);
  }

  function _setupMultipleOffRampsFromConfigs(
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs
  ) internal {
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2 * sourceChainConfigs.length);

    for (uint256 i = 0; i < sourceChainConfigs.length; ++i) {
      uint64 sourceChainSelector = sourceChainConfigs[i].sourceChainSelector;

      offRampUpdates[2 * i] = Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(s_offRamp)});
      offRampUpdates[2 * i + 1] = Router.OffRamp({
        sourceChainSelector: sourceChainSelector,
        offRamp: s_inboundNonceManager.getPreviousRamps(sourceChainSelector).prevOffRamp
      });
    }

    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }

  function _generateDynamicOffRampConfig(
    address feeQuoter
  ) internal pure returns (OffRamp.DynamicConfig memory) {
    return OffRamp.DynamicConfig({
      feeQuoter: feeQuoter,
      permissionLessExecutionThresholdSeconds: 60 * 60,
      messageInterceptor: address(0)
    });
  }

  function _generateAny2EVMMessageNoTokens(
    uint64 sourceChainSelector,
    bytes memory onRamp,
    uint64 sequenceNumber
  ) internal view returns (Internal.Any2EVMRampMessage memory) {
    return _generateAny2EVMMessage(sourceChainSelector, onRamp, sequenceNumber, new Client.EVMTokenAmount[](0), false);
  }

  function _generateAny2EVMMessageWithTokens(
    uint64 sourceChainSelector,
    bytes memory onRamp,
    uint64 sequenceNumber,
    uint256[] memory amounts
  ) internal view returns (Internal.Any2EVMRampMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].amount = amounts[i];
    }
    return _generateAny2EVMMessage(sourceChainSelector, onRamp, sequenceNumber, tokenAmounts, false);
  }

  function _generateAny2EVMMessageWithMaybeRevertingSingleToken(
    uint64 sequenceNumber,
    uint256 amount
  ) internal view returns (Internal.Any2EVMRampMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0].token = s_sourceTokens[1];
    tokenAmounts[0].amount = amount;

    return _generateAny2EVMMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, sequenceNumber, tokenAmounts, false);
  }

  function _generateAny2EVMMessage(
    uint64 sourceChainSelector,
    bytes memory onRamp,
    uint64 sequenceNumber,
    Client.EVMTokenAmount[] memory sourceTokenAmounts,
    bool allowOutOfOrderExecution
  ) internal view returns (Internal.Any2EVMRampMessage memory) {
    bytes memory data = abi.encode(0);

    Internal.Any2EVMTokenTransfer[] memory any2EVMTokenTransfer =
      new Internal.Any2EVMTokenTransfer[](sourceTokenAmounts.length);

    // Correctly set the TokenDataPayload for each token. Tokens have to be set up in the TokenSetup.
    for (uint256 i = 0; i < sourceTokenAmounts.length; ++i) {
      any2EVMTokenTransfer[i] = Internal.Any2EVMTokenTransfer({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[sourceTokenAmounts[i].token]),
        destTokenAddress: s_destTokenBySourceToken[sourceTokenAmounts[i].token],
        extraData: "",
        amount: sourceTokenAmounts[i].amount,
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      });
    }

    Internal.Any2EVMRampMessage memory message = Internal.Any2EVMRampMessage({
      header: Internal.RampMessageHeader({
        messageId: "",
        sourceChainSelector: sourceChainSelector,
        destChainSelector: DEST_CHAIN_SELECTOR,
        sequenceNumber: sequenceNumber,
        nonce: allowOutOfOrderExecution ? 0 : sequenceNumber
      }),
      sender: abi.encode(OWNER),
      data: data,
      receiver: address(s_receiver),
      tokenAmounts: any2EVMTokenTransfer,
      gasLimit: GAS_LIMIT
    });

    message.header.messageId = _hashMessage(message, onRamp);

    return message;
  }

  function _generateSingleBasicMessage(
    uint64 sourceChainSelector,
    bytes memory onRamp
  ) internal view returns (Internal.Any2EVMRampMessage[] memory) {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](1);
    messages[0] = _generateAny2EVMMessageNoTokens(sourceChainSelector, onRamp, 1);
    return messages;
  }

  function _generateMessagesWithTokens(
    uint64 sourceChainSelector,
    bytes memory onRamp
  ) internal view returns (Internal.Any2EVMRampMessage[] memory) {
    Internal.Any2EVMRampMessage[] memory messages = new Internal.Any2EVMRampMessage[](2);
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = 1e18;
    tokenAmounts[1].amount = 5e18;
    messages[0] = _generateAny2EVMMessage(sourceChainSelector, onRamp, 1, tokenAmounts, false);
    messages[1] = _generateAny2EVMMessage(sourceChainSelector, onRamp, 2, tokenAmounts, false);

    return messages;
  }

  function _getCastedSourceEVMTokenAmountsWithZeroAmounts()
    internal
    view
    returns (Client.EVMTokenAmount[] memory tokenAmounts)
  {
    tokenAmounts = new Client.EVMTokenAmount[](s_sourceTokens.length);
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].token = s_sourceTokens[i];
    }
    return tokenAmounts;
  }

  function _generateReportFromMessages(
    uint64 sourceChainSelector,
    Internal.Any2EVMRampMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReport memory) {
    bytes[][] memory offchainTokenData = new bytes[][](messages.length);

    for (uint256 i = 0; i < messages.length; ++i) {
      offchainTokenData[i] = new bytes[](messages[i].tokenAmounts.length);
    }

    return Internal.ExecutionReport({
      sourceChainSelector: sourceChainSelector,
      proofs: new bytes32[](0),
      proofFlagBits: 2 ** 256 - 1,
      messages: messages,
      offchainTokenData: offchainTokenData
    });
  }

  function _generateBatchReportFromMessages(
    uint64 sourceChainSelector,
    Internal.Any2EVMRampMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReport[] memory) {
    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](1);
    reports[0] = _generateReportFromMessages(sourceChainSelector, messages);
    return reports;
  }

  function _getGasLimitsFromMessages(
    Internal.Any2EVMRampMessage[] memory messages
  ) internal pure returns (OffRamp.GasLimitOverride[] memory) {
    OffRamp.GasLimitOverride[] memory gasLimits = new OffRamp.GasLimitOverride[](messages.length);
    for (uint256 i = 0; i < messages.length; ++i) {
      gasLimits[i].receiverExecutionGasLimit = messages[i].gasLimit;
    }

    return gasLimits;
  }

  function _assertSameConfig(OffRamp.DynamicConfig memory a, OffRamp.DynamicConfig memory b) public pure {
    assertEq(a.permissionLessExecutionThresholdSeconds, b.permissionLessExecutionThresholdSeconds);
    assertEq(a.messageInterceptor, b.messageInterceptor);
    assertEq(a.feeQuoter, b.feeQuoter);
  }

  function _assertSourceChainConfigEquality(
    OffRamp.SourceChainConfig memory config1,
    OffRamp.SourceChainConfig memory config2
  ) internal pure {
    assertEq(config1.isEnabled, config2.isEnabled);
    assertEq(config1.minSeqNr, config2.minSeqNr);
    assertEq(config1.onRamp, config2.onRamp);
    assertEq(address(config1.router), address(config2.router));
  }

  function _enableInboundMessageInterceptor() internal {
    OffRamp.DynamicConfig memory dynamicConfig = s_offRamp.getDynamicConfig();
    dynamicConfig.messageInterceptor = address(s_inboundMessageInterceptor);
    s_offRamp.setDynamicConfig(dynamicConfig);
  }

  function _redeployOffRampWithNoOCRConfigs() internal {
    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      new OffRamp.SourceChainConfigArgs[](0)
    );

    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(s_offRamp);
    s_inboundNonceManager.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );
    _setupMultipleOffRamps();

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = address(s_offRamp);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );
  }

  function _commit(OffRamp.CommitReport memory commitReport, uint64 sequenceNumber) internal {
    bytes32[2] memory reportContext = [s_configDigestCommit, bytes32(uint256(sequenceNumber))];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, F + 1);

    vm.startPrank(s_validTransmitters[0]);
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function _execute(
    Internal.ExecutionReport[] memory reports
  ) internal {
    bytes32[2] memory reportContext = [s_configDigestExec, s_configDigestExec];

    vm.startPrank(s_validTransmitters[0]);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function _assertExecutionStateChangedEventLogs(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    bytes32 messageId,
    bytes32 messageHash,
    Internal.MessageExecutionState state,
    bytes memory returnData
  ) internal {
    Vm.Log[] memory logs = vm.getRecordedLogs();
    for (uint256 i = 0; i < logs.length; ++i) {
      if (logs[i].topics[0] == OffRamp.ExecutionStateChanged.selector) {
        uint64 logSourceChainSelector = uint64(uint256(logs[i].topics[1]));
        uint64 logSequenceNumber = uint64(uint256(logs[i].topics[2]));
        bytes32 logMessageId = bytes32(logs[i].topics[3]);
        (bytes32 logMessageHash, uint8 logState, bytes memory logReturnData,) =
          abi.decode(logs[i].data, (bytes32, uint8, bytes, uint256));
        if (logMessageId == messageId) {
          assertEq(logSourceChainSelector, sourceChainSelector);
          assertEq(logSequenceNumber, sequenceNumber);
          assertEq(logMessageId, messageId);
          assertEq(logMessageHash, messageHash);
          assertEq(logState, uint8(state));
          assertEq(logReturnData, returnData);
        }
      }
    }
  }

  function _assertExecutionStateChangedEventLogs(
    Vm.Log[] memory logs,
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    bytes32 messageId,
    bytes32 messageHash,
    Internal.MessageExecutionState state,
    bytes memory returnData
  ) internal pure {
    for (uint256 i = 0; i < logs.length; ++i) {
      if (logs[i].topics[0] == OffRamp.ExecutionStateChanged.selector) {
        uint64 logSourceChainSelector = uint64(uint256(logs[i].topics[1]));
        uint64 logSequenceNumber = uint64(uint256(logs[i].topics[2]));
        bytes32 logMessageId = bytes32(logs[i].topics[3]);
        (bytes32 logMessageHash, uint8 logState, bytes memory logReturnData,) =
          abi.decode(logs[i].data, (bytes32, uint8, bytes, uint256));
        if (logMessageId == messageId) {
          assertEq(logSourceChainSelector, sourceChainSelector);
          assertEq(logSequenceNumber, sequenceNumber);
          assertEq(logMessageId, messageId);
          assertEq(logMessageHash, messageHash);
          assertEq(logState, uint8(state));
          assertEq(logReturnData, returnData);
        }
      }
    }
  }

  function _assertNoEmit(
    bytes32 eventSelector
  ) internal {
    Vm.Log[] memory logs = vm.getRecordedLogs();

    for (uint256 i = 0; i < logs.length; ++i) {
      assertTrue(logs[i].topics[0] != eventSelector);
    }
  }

  function _hashMessage(
    Internal.Any2EVMRampMessage memory message,
    bytes memory onRamp
  ) internal pure returns (bytes32) {
    return Internal._hash(
      message,
      keccak256(
        abi.encode(
          Internal.ANY_2_EVM_MESSAGE_HASH,
          message.header.sourceChainSelector,
          message.header.destChainSelector,
          keccak256(onRamp)
        )
      )
    );
  }

  function _getEmptyPriceUpdates() internal pure returns (Internal.PriceUpdates memory priceUpdates) {
    return Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
  }
}
