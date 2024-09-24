// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {ICommitStore} from "../../interfaces/ICommitStore.sol";

import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {TokenSetup} from "../TokenSetup.t.sol";

import {FeeQuoterSetup} from "../feeQuoter/FeeQuoterSetup.t.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {MockCommitStore} from "../mocks/MockCommitStore.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";

contract EVM2EVMOffRampSetup is TokenSetup, FeeQuoterSetup, OCR2BaseSetup {
  MockCommitStore internal s_mockCommitStore;
  IAny2EVMMessageReceiver internal s_receiver;
  IAny2EVMMessageReceiver internal s_secondary_receiver;
  MaybeRevertMessageReceiver internal s_reverting_receiver;

  MaybeRevertingBurnMintTokenPool internal s_maybeRevertingPool;

  EVM2EVMOffRampHelper internal s_offRamp;
  address internal s_sourceTokenPool = makeAddr("sourceTokenPool");

  function setUp() public virtual override(TokenSetup, FeeQuoterSetup, OCR2BaseSetup) {
    TokenSetup.setUp();
    FeeQuoterSetup.setUp();
    OCR2BaseSetup.setUp();

    s_mockCommitStore = new MockCommitStore();
    s_receiver = new MaybeRevertMessageReceiver(false);
    s_secondary_receiver = new MaybeRevertMessageReceiver(false);
    s_reverting_receiver = new MaybeRevertMessageReceiver(true);

    s_maybeRevertingPool = MaybeRevertingBurnMintTokenPool(s_destPoolByToken[s_destTokens[1]]);

    deployOffRamp(s_mockCommitStore, s_destRouter, address(0));
  }

  function deployOffRamp(ICommitStore commitStore, Router router, address prevOffRamp) internal {
    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(commitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        prevOffRamp: prevOffRamp,
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _getInboundRateLimiterConfig()
    );
    s_offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(generateDynamicOffRampConfig(address(router), address(s_feeQuoter))),
      s_offchainConfigVersion,
      abi.encode("")
    );

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: address(s_offRamp)});
    offRampUpdates[1] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: address(prevOffRamp)});
    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
    EVM2EVMOffRamp.RateLimitToken[] memory tokensToAdd = new EVM2EVMOffRamp.RateLimitToken[](s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] = EVM2EVMOffRamp.RateLimitToken({sourceToken: s_sourceTokens[i], destToken: s_destTokens[i]});
    }
    s_offRamp.updateRateLimitTokens(new EVM2EVMOffRamp.RateLimitToken[](0), tokensToAdd);
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
      maxDataBytes: MAX_DATA_SIZE
    });
  }

  function _convertToGeneralMessage(
    Internal.EVM2EVMMessage memory original
  ) internal view returns (Client.Any2EVMMessage memory message) {
    uint256 numberOfTokens = original.tokenAmounts.length;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](numberOfTokens);

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Internal.SourceTokenData memory sourceTokenData =
        abi.decode(original.sourceTokenData[i], (Internal.SourceTokenData));

      address destPoolAddress = abi.decode(sourceTokenData.destTokenAddress, (address));
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
    uint64 sequenceNumber
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    return _generateAny2EVMMessage(sequenceNumber, new Client.EVMTokenAmount[](0), false);
  }

  function _generateAny2EVMMessageWithTokens(
    uint64 sequenceNumber,
    uint256[] memory amounts
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].amount = amounts[i];
    }
    return _generateAny2EVMMessage(sequenceNumber, tokenAmounts, false);
  }

  function _generateAny2EVMMessage(
    uint64 sequenceNumber,
    Client.EVMTokenAmount[] memory tokenAmounts,
    bool allowOutOfOrderExecution
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    bytes memory data = abi.encode(0);
    Internal.EVM2EVMMessage memory message = Internal.EVM2EVMMessage({
      sequenceNumber: sequenceNumber,
      sender: OWNER,
      nonce: allowOutOfOrderExecution ? 0 : sequenceNumber,
      gasLimit: GAS_LIMIT,
      strict: false,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
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
          destTokenAddress: abi.encode(s_destTokenBySourceToken[tokenAmounts[i].token]),
          extraData: "",
          destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        })
      );
    }

    message.messageId = Internal._hash(
      message,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );

    return message;
  }

  function _generateSingleBasicMessage() internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](1);
    messages[0] = _generateAny2EVMMessageNoTokens(1);
    return messages;
  }

  function _generateSingleBasicMessageWithTokens() internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](1);
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = 1e18;
    messages[0] = _generateAny2EVMMessage(1, tokenAmounts, false);
    return messages;
  }

  function _generateMessagesWithTokens() internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](2);
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = 1e18;
    tokenAmounts[1].amount = 5e18;
    messages[0] = _generateAny2EVMMessage(1, tokenAmounts, false);
    messages[1] = _generateAny2EVMMessage(2, tokenAmounts, false);

    return messages;
  }

  function _generateReportFromMessages(
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReport memory) {
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

  function _getGasLimitsFromMessages(
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (EVM2EVMOffRamp.GasLimitOverride[] memory) {
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = new EVM2EVMOffRamp.GasLimitOverride[](messages.length);
    for (uint256 i = 0; i < messages.length; ++i) {
      gasLimitOverrides[i].receiverExecutionGasLimit = messages[i].gasLimit;
      gasLimitOverrides[i].tokenGasOverrides = new uint32[](messages[i].tokenAmounts.length);

      for (uint256 j = 0; j < messages[i].tokenAmounts.length; ++j) {
        gasLimitOverrides[i].tokenGasOverrides[j] = DEFAULT_TOKEN_DEST_GAS_OVERHEAD + 1;
      }
    }

    return gasLimitOverrides;
  }

  function _assertSameConfig(EVM2EVMOffRamp.DynamicConfig memory a, EVM2EVMOffRamp.DynamicConfig memory b) public pure {
    assertEq(a.permissionLessExecutionThresholdSeconds, b.permissionLessExecutionThresholdSeconds);
    assertEq(a.router, b.router);
    assertEq(a.priceRegistry, b.priceRegistry);
    assertEq(a.maxNumberOfTokensPerMsg, b.maxNumberOfTokensPerMsg);
    assertEq(a.maxDataBytes, b.maxDataBytes);
  }

  function _getDefaultSourceTokenData(
    Client.EVMTokenAmount[] memory srcTokenAmounts
  ) internal view returns (bytes[] memory) {
    bytes[] memory sourceTokenData = new bytes[](srcTokenAmounts.length);
    for (uint256 i = 0; i < srcTokenAmounts.length; ++i) {
      sourceTokenData[i] = abi.encode(
        Internal.SourceTokenData({
          sourcePoolAddress: abi.encode(s_sourcePoolByToken[srcTokenAmounts[i].token]),
          destTokenAddress: abi.encode(s_destTokenBySourceToken[srcTokenAmounts[i].token]),
          extraData: "",
          destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        })
      );
    }
    return sourceTokenData;
  }
}
