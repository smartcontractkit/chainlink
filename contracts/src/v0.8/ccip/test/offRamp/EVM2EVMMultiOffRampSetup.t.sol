// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {ICommitStore} from "../../interfaces/ICommitStore.sol";

import {Router} from "../../Router.sol";
import {IAny2EVMOffRamp} from "../../interfaces/IAny2EVMOffRamp.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMMultiOffRamp} from "../../offRamp/EVM2EVMMultiOffRamp.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {EVM2EVMMultiOffRampHelper} from "../helpers/EVM2EVMMultiOffRampHelper.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MockCommitStore} from "../mocks/MockCommitStore.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMMultiOffRampSetup is TokenSetup, PriceRegistrySetup, OCR2BaseSetup {
  uint64 internal constant SOURCE_CHAIN_SELECTOR_1 = SOURCE_CHAIN_SELECTOR;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_2 = 6433500567565415381;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_3 = 4051577828743386545;

  address internal constant ON_RAMP_ADDRESS_1 = ON_RAMP_ADDRESS;
  address internal constant ON_RAMP_ADDRESS_2 = 0xaA3f843Cf8E33B1F02dd28303b6bD87B1aBF8AE4;
  address internal constant ON_RAMP_ADDRESS_3 = 0x71830C37Cb193e820de488Da111cfbFcC680a1b9;

  MockCommitStore internal s_mockCommitStore;
  IAny2EVMMessageReceiver internal s_receiver;
  IAny2EVMMessageReceiver internal s_secondary_receiver;
  MaybeRevertMessageReceiver internal s_reverting_receiver;

  MaybeRevertingBurnMintTokenPool internal s_maybeRevertingPool;

  EVM2EVMMultiOffRampHelper internal s_offRamp;
  address internal s_sourceTokenPool = makeAddr("sourceTokenPool");

  event ExecutionStateChanged(
    uint64 indexed sourceChainSelector,
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    Internal.MessageExecutionState state,
    bytes returnData
  );
  event SkippedIncorrectNonce(uint64 sourceChainSelector, uint64 nonce, address indexed sender);
  event SkippedAlreadyExecutedMessage(uint64 indexed sequenceNumber);

  function setUp() public virtual override(TokenSetup, PriceRegistrySetup, OCR2BaseSetup) {
    TokenSetup.setUp();
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();

    s_mockCommitStore = new MockCommitStore();
    s_receiver = new MaybeRevertMessageReceiver(false);
    s_secondary_receiver = new MaybeRevertMessageReceiver(false);
    s_reverting_receiver = new MaybeRevertMessageReceiver(true);

    s_maybeRevertingPool = MaybeRevertingBurnMintTokenPool(s_destPoolByToken[s_destTokens[1]]);

    deployOffRamp(s_mockCommitStore, s_destRouter);
  }

  function deployOffRamp(ICommitStore commitStore, Router router) internal {
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](0);

    s_offRamp = new EVM2EVMMultiOffRampHelper(
      EVM2EVMMultiOffRamp.StaticConfig({
        commitStore: address(commitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        rmnProxy: address(s_mockRMN)
      }),
      sourceChainConfigs,
      getInboundRateLimiterConfig()
    );
    s_offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(generateDynamicMultiOffRampConfig(address(router), address(s_priceRegistry))),
      s_offchainConfigVersion,
      abi.encode("")
    );

    EVM2EVMMultiOffRamp.RateLimitToken[] memory tokensToAdd =
      new EVM2EVMMultiOffRamp.RateLimitToken[](s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = EVM2EVMMultiOffRamp.RateLimitToken({sourceToken: s_sourceTokens[i], destToken: s_destTokens[i]});
    }
    s_offRamp.updateRateLimitTokens(new EVM2EVMMultiOffRamp.RateLimitToken[](0), tokensToAdd);
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
        rmnProxy: address(s_mockRMN)
      }),
      getInboundRateLimiterConfig()
    );
    offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(generateDynamicOffRampConfig(address(router), address(s_priceRegistry))),
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
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      isEnabled: true,
      prevOffRamp: address(0),
      onRamp: ON_RAMP_ADDRESS_1
    });
    sourceChainConfigs[1] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_2,
      isEnabled: false,
      prevOffRamp: address(0),
      onRamp: ON_RAMP_ADDRESS_2
    });
    sourceChainConfigs[2] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_3,
      isEnabled: true,
      prevOffRamp: address(0),
      onRamp: ON_RAMP_ADDRESS_3
    });
    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);
  }

  function _setupMultipleOffRampsFromConfigs(EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs)
    internal
  {
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](sourceChainConfigs.length);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2 * onRampUpdates.length);

    for (uint256 i = 0; i < sourceChainConfigs.length; ++i) {
      uint64 sourceChainSelector = sourceChainConfigs[i].sourceChainSelector;

      onRampUpdates[i] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: sourceChainConfigs[i].onRamp});

      offRampUpdates[2 * i] = Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(s_offRamp)});
      offRampUpdates[2 * i + 1] =
        Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(sourceChainConfigs[i].prevOffRamp)});
    }

    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }

  function generateDynamicOffRampConfig(
    address router,
    address priceRegistry
  ) internal pure returns (EVM2EVMOffRamp.DynamicConfig memory) {
    return EVM2EVMOffRamp.DynamicConfig({
      permissionLessExecutionThresholdSeconds: PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS,
      router: router,
      priceRegistry: priceRegistry,
      maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
      maxDataBytes: MAX_DATA_SIZE,
      maxPoolReleaseOrMintGas: MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS
    });
  }

  function generateDynamicMultiOffRampConfig(
    address router,
    address priceRegistry
  ) internal pure returns (EVM2EVMMultiOffRamp.DynamicConfig memory) {
    return EVM2EVMMultiOffRamp.DynamicConfig({
      permissionLessExecutionThresholdSeconds: PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS,
      router: router,
      priceRegistry: priceRegistry,
      maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
      maxDataBytes: MAX_DATA_SIZE,
      maxPoolReleaseOrMintGas: MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS
    });
  }

  function _convertToGeneralMessage(Internal.EVM2EVMMessage memory original)
    internal
    view
    returns (Client.Any2EVMMessage memory message)
  {
    uint256 numberOfTokens = original.tokenAmounts.length;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](numberOfTokens);

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Internal.SourceTokenData memory sourceTokenData =
        abi.decode(original.sourceTokenData[i], (Internal.SourceTokenData));

      address destPoolAddress = abi.decode(sourceTokenData.destPoolAddress, (address));
      TokenPool pool = TokenPool(destPoolAddress);
      destTokenAmounts[i].token = address(pool.getToken());
      destTokenAmounts[i].amount = original.tokenAmounts[i].amount;
    }

    return Client.Any2EVMMessage({
      messageId: original.messageId,
      sourceChainSelector: original.sourceChainSelector,
      sender: abi.encode(original.sender),
      data: original.data,
      destTokenAmounts: destTokenAmounts
    });
  }

  function _generateAny2EVMMessageNoTokens(
    uint64 sourceChainSelector,
    address onRamp,
    uint64 sequenceNumber
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    return _generateAny2EVMMessage(sourceChainSelector, onRamp, sequenceNumber, new Client.EVMTokenAmount[](0));
  }

  function _generateAny2EVMMessageWithTokens(
    uint64 sourceChainSelector,
    address onRamp,
    uint64 sequenceNumber,
    uint256[] memory amounts
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].amount = amounts[i];
    }
    return _generateAny2EVMMessage(sourceChainSelector, onRamp, sequenceNumber, tokenAmounts);
  }

  function _generateAny2EVMMessage(
    uint64 sourceChainSelector,
    address onRamp,
    uint64 sequenceNumber,
    Client.EVMTokenAmount[] memory tokenAmounts
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    bytes memory data = abi.encode(0);
    Internal.EVM2EVMMessage memory message = Internal.EVM2EVMMessage({
      sequenceNumber: sequenceNumber,
      sender: OWNER,
      nonce: sequenceNumber,
      gasLimit: GAS_LIMIT,
      strict: false,
      sourceChainSelector: sourceChainSelector,
      receiver: address(s_receiver),
      data: data,
      tokenAmounts: tokenAmounts,
      sourceTokenData: new bytes[](tokenAmounts.length),
      feeToken: s_destFeeToken,
      feeTokenAmount: uint256(0),
      messageId: ""
    });

    // Correctly set the TokenDataPayload for each token. Tokens have to be set up in the TokenSetup.
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      message.sourceTokenData[i] = abi.encode(
        Internal.SourceTokenData({
          sourcePoolAddress: abi.encode(s_sourcePoolByToken[tokenAmounts[i].token]),
          destPoolAddress: abi.encode(s_destPoolBySourceToken[tokenAmounts[i].token]),
          extraData: ""
        })
      );
    }

    message.messageId = Internal._hash(
      message, keccak256(abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, sourceChainSelector, DEST_CHAIN_SELECTOR, onRamp))
    );

    return message;
  }

  function _generateBasicMessages(
    uint64 sourceChainSelector,
    address onRamp
  ) internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](1);
    messages[0] = _generateAny2EVMMessageNoTokens(sourceChainSelector, onRamp, 1);
    return messages;
  }

  function _generateMessagesWithTokens(
    uint64 sourceChainSelector,
    address onRamp
  ) internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](2);
    Client.EVMTokenAmount[] memory tokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = 1e18;
    tokenAmounts[1].amount = 5e18;
    messages[0] = _generateAny2EVMMessage(sourceChainSelector, onRamp, 1, tokenAmounts);
    messages[1] = _generateAny2EVMMessage(sourceChainSelector, onRamp, 2, tokenAmounts);

    return messages;
  }

  function _generateSingleRampReportFromMessages(
    uint64 sourceChainSelector,
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReport memory) {
    Internal.ExecutionReportSingleChain memory singleChainReport =
      _generateReportFromMessages(sourceChainSelector, messages);

    return Internal.ExecutionReport({
      proofs: singleChainReport.proofs,
      proofFlagBits: singleChainReport.proofFlagBits,
      messages: singleChainReport.messages,
      offchainTokenData: singleChainReport.offchainTokenData
    });
  }

  function _generateReportFromMessages(
    uint64 sourceChainSelector,
    Internal.EVM2EVMMessage[] memory messages
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
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReportSingleChain[] memory) {
    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](1);
    reports[0] = _generateReportFromMessages(sourceChainSelector, messages);
    return reports;
  }

  function _getGasLimitsFromMessages(Internal.EVM2EVMMessage[] memory messages)
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

  function _assertSameConfig(
    EVM2EVMMultiOffRamp.DynamicConfig memory a,
    EVM2EVMMultiOffRamp.DynamicConfig memory b
  ) public pure {
    assertEq(a.permissionLessExecutionThresholdSeconds, b.permissionLessExecutionThresholdSeconds);
    assertEq(a.router, b.router);
    assertEq(a.priceRegistry, b.priceRegistry);
    assertEq(a.maxNumberOfTokensPerMsg, b.maxNumberOfTokensPerMsg);
    assertEq(a.maxDataBytes, b.maxDataBytes);
    assertEq(a.maxPoolReleaseOrMintGas, b.maxPoolReleaseOrMintGas);
  }

  function _assertSourceChainConfigEquality(
    EVM2EVMMultiOffRamp.SourceChainConfig memory config1,
    EVM2EVMMultiOffRamp.SourceChainConfig memory config2
  ) internal pure {
    assertEq(config1.isEnabled, config2.isEnabled);
    assertEq(config1.prevOffRamp, config2.prevOffRamp);
    assertEq(config1.onRamp, config2.onRamp);
    assertEq(config1.metadataHash, config2.metadataHash);
  }

  function _getDefaultSourceTokenData(Client.EVMTokenAmount[] memory srcTokenAmounts)
    internal
    view
    returns (bytes[] memory)
  {
    bytes[] memory sourceTokenData = new bytes[](srcTokenAmounts.length);
    for (uint256 i = 0; i < srcTokenAmounts.length; ++i) {
      sourceTokenData[i] = abi.encode(
        Internal.SourceTokenData({
          sourcePoolAddress: abi.encode(s_sourcePoolByToken[srcTokenAmounts[i].token]),
          destPoolAddress: abi.encode(s_destPoolBySourceToken[srcTokenAmounts[i].token]),
          extraData: ""
        })
      );
    }
    return sourceTokenData;
  }
}
