// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {ICommitStore} from "../../interfaces/ICommitStore.sol";
import {IRMN} from "../../interfaces/IRMN.sol";

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {NonceManager} from "../../NonceManager.sol";
import {RMN} from "../../RMN.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MessageInterceptorHelper} from "../helpers/MessageInterceptorHelper.sol";
import {OffRampHelper} from "../helpers/OffRampHelper.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MultiOCR3BaseSetup} from "../ocr/MultiOCR3BaseSetup.t.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistrySetup.t.sol";
import {Vm} from "forge-std/Test.sol";

contract OffRampSetup is PriceRegistrySetup, MultiOCR3BaseSetup {
  uint64 internal constant SOURCE_CHAIN_SELECTOR_1 = SOURCE_CHAIN_SELECTOR;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_2 = 6433500567565415381;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_3 = 4051577828743386545;
  bytes32 internal constant EXECUTION_STATE_CHANGE_TOPIC_HASH =
    keccak256("ExecutionStateChanged(uint64,uint64,bytes32,uint8,bytes,uint256)");

  bytes internal constant ON_RAMP_ADDRESS_1 = abi.encode(ON_RAMP_ADDRESS);
  bytes internal constant ON_RAMP_ADDRESS_2 = abi.encode(0xaA3f843Cf8E33B1F02dd28303b6bD87B1aBF8AE4);
  bytes internal constant ON_RAMP_ADDRESS_3 = abi.encode(0x71830C37Cb193e820de488Da111cfbFcC680a1b9);

  address internal constant BLESS_VOTE_ADDR = address(8888);

  IAny2EVMMessageReceiver internal s_receiver;
  IAny2EVMMessageReceiver internal s_secondary_receiver;
  MaybeRevertMessageReceiver internal s_reverting_receiver;

  MaybeRevertingBurnMintTokenPool internal s_maybeRevertingPool;

  OffRampHelper internal s_offRamp;
  MessageInterceptorHelper internal s_inboundMessageValidator;
  NonceManager internal s_inboundNonceManager;
  RMN internal s_realRMN;
  address internal s_sourceTokenPool = makeAddr("sourceTokenPool");

  bytes32 internal s_configDigestExec;
  bytes32 internal s_configDigestCommit;
  uint64 internal constant s_offchainConfigVersion = 3;
  uint8 internal constant s_F = 1;

  uint64 internal s_latestSequenceNumber;

  function setUp() public virtual override(PriceRegistrySetup, MultiOCR3BaseSetup) {
    PriceRegistrySetup.setUp();
    MultiOCR3BaseSetup.setUp();

    s_inboundMessageValidator = new MessageInterceptorHelper();
    s_receiver = new MaybeRevertMessageReceiver(false);
    s_secondary_receiver = new MaybeRevertMessageReceiver(false);
    s_reverting_receiver = new MaybeRevertMessageReceiver(true);

    s_maybeRevertingPool = MaybeRevertingBurnMintTokenPool(s_destPoolByToken[s_destTokens[1]]);
    s_inboundNonceManager = new NonceManager(new address[](0));

    _deployOffRamp(s_mockRMN, s_inboundNonceManager);
  }

  function _deployOffRamp(IRMN rmnProxy, NonceManager nonceManager) internal {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(rmnProxy),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(nonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_priceRegistry)),
      sourceChainConfigs
    );

    s_configDigestExec = _getBasicConfigDigest(s_F, s_emptySigners, s_validTransmitters);
    s_configDigestCommit = _getBasicConfigDigest(s_F, s_validSigners, s_validTransmitters);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](2);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: s_F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[1] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: s_F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    s_offRamp.setDynamicConfig(_generateDynamicOffRampConfig(address(s_priceRegistry)));
    s_offRamp.setOCR3Configs(ocrConfigs);

    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(s_offRamp);
    NonceManager(nonceManager).applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = address(s_offRamp);
    s_priceRegistry.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );
  }

  // TODO: function can be made common across OffRampSetup and MultiOffRampSetup
  function _deploySingleLaneOffRamp(
    ICommitStore commitStore,
    Router router,
    address prevOffRamp,
    uint64 sourceChainSelector,
    address onRampAddress
  ) internal returns (EVM2EVMOffRampHelper) {
    EVM2EVMOffRampHelper offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(commitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: sourceChainSelector,
        onRamp: onRampAddress,
        prevOffRamp: prevOffRamp,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _getInboundRateLimiterConfig()
    );
    offRamp.setOCR2Config(
      s_validSigners,
      s_validTransmitters,
      s_F,
      abi.encode(_generateDynamicOffRampConfig(address(router), address(s_priceRegistry))),
      s_offchainConfigVersion,
      abi.encode("")
    );

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(s_offRamp)});
    offRampUpdates[1] = Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(prevOffRamp)});
    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
    EVM2EVMOffRamp.RateLimitToken[] memory tokensToAdd = new EVM2EVMOffRamp.RateLimitToken[](s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[i], destToken: s_destTokens[i]});
    }
    offRamp.updateRateLimitTokens(new EVM2EVMOffRamp.RateLimitToken[](0), tokensToAdd);

    return offRamp;
  }

  function _setupMultipleOffRamps() internal {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_2,
      onRamp: ON_RAMP_ADDRESS_2,
      isEnabled: false
    });
    sourceChainConfigs[2] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_3,
      onRamp: ON_RAMP_ADDRESS_3,
      isEnabled: true
    });
    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);
  }

  function _setupMultipleOffRampsFromConfigs(OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs) internal {
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
    address router,
    address priceRegistry
  ) internal pure returns (EVM2EVMOffRamp.DynamicConfig memory) {
    return EVM2EVMOffRamp.DynamicConfig({
      permissionLessExecutionThresholdSeconds: PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS,
      router: router,
      priceRegistry: priceRegistry,
      maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
      maxDataBytes: MAX_DATA_SIZE
    });
  }

  uint32 internal constant MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS = 200_000;
  uint32 internal constant MAX_TOKEN_POOL_TRANSFER_GAS = 50_000;

  function _generateDynamicOffRampConfig(address priceRegistry) internal pure returns (OffRamp.DynamicConfig memory) {
    return OffRamp.DynamicConfig({
      permissionLessExecutionThresholdSeconds: PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS,
      priceRegistry: priceRegistry,
      messageValidator: address(0),
      maxPoolReleaseOrMintGas: MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS,
      maxTokenTransferGas: MAX_TOKEN_POOL_TRANSFER_GAS
    });
  }

  function _convertToGeneralMessage(Internal.Any2EVMRampMessage memory original)
    internal
    view
    returns (Client.Any2EVMMessage memory message)
  {
    uint256 numberOfTokens = original.tokenAmounts.length;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](numberOfTokens);

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Internal.RampTokenAmount memory tokenAmount = original.tokenAmounts[i];

      address destPoolAddress = abi.decode(tokenAmount.destTokenAddress, (address));
      TokenPool pool = TokenPool(destPoolAddress);
      destTokenAmounts[i].token = address(pool.getToken());
      destTokenAmounts[i].amount = tokenAmount.amount;
    }

    return Client.Any2EVMMessage({
      messageId: original.header.messageId,
      sourceChainSelector: original.header.sourceChainSelector,
      sender: abi.encode(original.sender),
      data: original.data,
      destTokenAmounts: destTokenAmounts
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

  function _generateAny2EVMMessage(
    uint64 sourceChainSelector,
    bytes memory onRamp,
    uint64 sequenceNumber,
    Client.EVMTokenAmount[] memory tokenAmounts,
    bool allowOutOfOrderExecution
  ) internal view returns (Internal.Any2EVMRampMessage memory) {
    bytes memory data = abi.encode(0);

    Internal.RampTokenAmount[] memory rampTokenAmounts = new Internal.RampTokenAmount[](tokenAmounts.length);

    // Correctly set the TokenDataPayload for each token. Tokens have to be set up in the TokenSetup.
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      rampTokenAmounts[i] = Internal.RampTokenAmount({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[tokenAmounts[i].token]),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[tokenAmounts[i].token]),
        extraData: "",
        amount: tokenAmounts[i].amount
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
      tokenAmounts: rampTokenAmounts,
      gasLimit: GAS_LIMIT
    });

    message.header.messageId = Internal._hash(message, onRamp);

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

  function _generateReportFromMessages(
    uint64 sourceChainSelector,
    Internal.Any2EVMRampMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReportSingleChain memory) {
    bytes[][] memory offchainTokenData = new bytes[][](messages.length);

    for (uint256 i = 0; i < messages.length; ++i) {
      offchainTokenData[i] = new bytes[](messages[i].tokenAmounts.length);
    }

    return Internal.ExecutionReportSingleChain({
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
  ) internal pure returns (Internal.ExecutionReportSingleChain[] memory) {
    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](1);
    reports[0] = _generateReportFromMessages(sourceChainSelector, messages);
    return reports;
  }

  function _getGasLimitsFromMessages(Internal.Any2EVMRampMessage[] memory messages)
    internal
    pure
    returns (uint256[] memory)
  {
    uint256[] memory gasLimits = new uint256[](messages.length);
    for (uint256 i = 0; i < messages.length; ++i) {
      gasLimits[i] = messages[i].gasLimit;
    }

    return gasLimits;
  }

  function _assertSameConfig(OffRamp.DynamicConfig memory a, OffRamp.DynamicConfig memory b) public pure {
    assertEq(a.permissionLessExecutionThresholdSeconds, b.permissionLessExecutionThresholdSeconds);
    assertEq(a.maxPoolReleaseOrMintGas, b.maxPoolReleaseOrMintGas);
    assertEq(a.maxTokenTransferGas, b.maxTokenTransferGas);
    assertEq(a.messageValidator, b.messageValidator);
    assertEq(a.priceRegistry, b.priceRegistry);
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

  function _getDefaultSourceTokenData(Client.EVMTokenAmount[] memory srcTokenAmounts)
    internal
    view
    returns (Internal.RampTokenAmount[] memory)
  {
    Internal.RampTokenAmount[] memory sourceTokenData = new Internal.RampTokenAmount[](srcTokenAmounts.length);
    for (uint256 i = 0; i < srcTokenAmounts.length; ++i) {
      sourceTokenData[i] = Internal.RampTokenAmount({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[srcTokenAmounts[i].token]),
        destTokenAddress: abi.encode(s_destTokenBySourceToken[srcTokenAmounts[i].token]),
        extraData: "",
        amount: srcTokenAmounts[i].amount
      });
    }
    return sourceTokenData;
  }

  function _enableInboundMessageValidator() internal {
    OffRamp.DynamicConfig memory dynamicConfig = s_offRamp.getDynamicConfig();
    dynamicConfig.messageValidator = address(s_inboundMessageValidator);
    s_offRamp.setDynamicConfig(dynamicConfig);
  }

  function _redeployOffRampWithNoOCRConfigs() internal {
    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_priceRegistry)),
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
    s_priceRegistry.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );
  }

  function _setupRealRMN() internal {
    RMN.Voter[] memory voters = new RMN.Voter[](1);
    voters[0] =
      RMN.Voter({blessVoteAddr: BLESS_VOTE_ADDR, curseVoteAddr: address(9999), blessWeight: 1, curseWeight: 1});
    // Overwrite base mock rmn with real.
    s_realRMN = new RMN(RMN.Config({voters: voters, blessWeightThreshold: 1, curseWeightThreshold: 1}));
  }

  function _commit(OffRamp.CommitReport memory commitReport, uint64 sequenceNumber) internal {
    bytes32[3] memory reportContext = [s_configDigestCommit, bytes32(uint256(sequenceNumber)), s_configDigestCommit];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, s_F + 1);

    vm.startPrank(s_validTransmitters[0]);
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function _execute(Internal.ExecutionReportSingleChain[] memory reports) internal {
    bytes32[3] memory reportContext = [s_configDigestExec, s_configDigestExec, s_configDigestExec];

    vm.startPrank(s_validTransmitters[0]);
    s_offRamp.execute(reportContext, abi.encode(reports));
  }

  function assertExecutionStateChangedEventLogs(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    bytes32 messageId,
    Internal.MessageExecutionState state,
    bytes memory returnData
  ) public {
    Vm.Log[] memory logs = vm.getRecordedLogs();
    for (uint256 i = 0; i < logs.length; ++i) {
      if (logs[i].topics[0] == EXECUTION_STATE_CHANGE_TOPIC_HASH) {
        uint64 logSourceChainSelector = uint64(uint256(logs[i].topics[1]));
        uint64 logSequenceNumber = uint64(uint256(logs[i].topics[2]));
        bytes32 logMessageId = bytes32(logs[i].topics[3]);
        (uint8 logState, bytes memory logReturnData,) = abi.decode(logs[i].data, (uint8, bytes, uint256));
        if (logMessageId == messageId) {
          assertEq(logSourceChainSelector, sourceChainSelector);
          assertEq(logSequenceNumber, sequenceNumber);
          assertEq(logMessageId, messageId);
          assertEq(logState, uint8(state));
          assertEq(logReturnData, returnData);
        }
      }
    }
  }
}
