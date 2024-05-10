// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {ICommitStore} from "../../interfaces/ICommitStore.sol";

import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMMultiOffRamp} from "../../offRamp/EVM2EVMMultiOffRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {EVM2EVMMultiOffRampHelper} from "../helpers/EVM2EVMMultiOffRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MockCommitStore} from "../mocks/MockCommitStore.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMMultiOffRampSetup is TokenSetup, PriceRegistrySetup, OCR2BaseSetup {
  uint64 internal constant SOURCE_CHAIN_SELECTOR_1 = 16015286601757825753;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_2 = 6433500567565415381;
  uint64 internal constant SOURCE_CHAIN_SELECTOR_3 = 4051577828743386545;

  address internal constant ON_RAMP_ADDRESS_1 = 0x11118e64e1FB0c487f25dD6D3601FF6aF8d32E4e;
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

  function _generateReportFromMessages(Internal.EVM2EVMMessage[] memory messages)
    internal
    pure
    returns (Internal.ExecutionReport memory)
  {
    bytes[][] memory offchainTokenData = new bytes[][](messages.length);

    for (uint256 i = 0; i < messages.length; ++i) {
      offchainTokenData[i] = new bytes[](messages[i].tokenAmounts.length);
    }

    return Internal.ExecutionReport({
      proofs: new bytes32[](0),
      proofFlagBits: 2 ** 256 - 1,
      messages: messages,
      offchainTokenData: offchainTokenData
    });
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
