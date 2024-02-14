// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ICommitStore} from "../../interfaces/ICommitStore.sol";
import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {IPool} from "../../interfaces/pools/IPool.sol";

import {Internal} from "../../libraries/Internal.sol";
import {Client} from "../../libraries/Client.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";
import {MockCommitStore} from "../mocks/MockCommitStore.sol";
import {Router} from "../../Router.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {EVM2EVMOffRampHelper} from "../helpers/EVM2EVMOffRampHelper.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMOffRampSetup is TokenSetup, PriceRegistrySetup, OCR2BaseSetup {
  MockCommitStore internal s_mockCommitStore;
  IAny2EVMMessageReceiver internal s_receiver;
  IAny2EVMMessageReceiver internal s_secondary_receiver;
  MaybeRevertMessageReceiver internal s_reverting_receiver;

  EVM2EVMOffRampHelper internal s_offRamp;

  event ExecutionStateChanged(
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    Internal.MessageExecutionState state,
    bytes returnData
  );
  event SkippedIncorrectNonce(uint64 indexed nonce, address indexed sender);

  function setUp() public virtual override(TokenSetup, PriceRegistrySetup, OCR2BaseSetup) {
    TokenSetup.setUp();
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();

    s_mockCommitStore = new MockCommitStore();
    s_receiver = new MaybeRevertMessageReceiver(false);
    s_secondary_receiver = new MaybeRevertMessageReceiver(false);
    s_reverting_receiver = new MaybeRevertMessageReceiver(true);

    deployOffRamp(s_mockCommitStore, s_destRouter, address(0));

    TokenPool.ChainUpdate[] memory offRamps = new TokenPool.ChainUpdate[](1);
    offRamps[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });

    LockReleaseTokenPool(address(s_destPools[0])).applyChainUpdates(offRamps);
    LockReleaseTokenPool(address(s_destPools[1])).applyChainUpdates(offRamps);
  }

  function deployOffRamp(ICommitStore commitStore, Router router, address prevOffRamp) internal {
    s_offRamp = new EVM2EVMOffRampHelper(
      EVM2EVMOffRamp.StaticConfig({
        commitStore: address(commitStore),
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        prevOffRamp: prevOffRamp,
        armProxy: address(s_mockARM)
      }),
      getCastedSourceTokens(),
      getCastedDestinationPools(),
      getInboundRateLimiterConfig()
    );
    s_offRamp.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(generateDynamicOffRampConfig(address(router), address(s_priceRegistry))),
      s_offchainConfigVersion,
      abi.encode("")
    );

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: address(s_offRamp)});
    offRampUpdates[1] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: address(prevOffRamp)});
    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }

  function _convertToGeneralMessage(
    Internal.EVM2EVMMessage memory original
  ) internal view returns (Client.Any2EVMMessage memory message) {
    uint256 numberOfTokens = original.tokenAmounts.length;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](numberOfTokens);

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      IPool pool = s_offRamp.getPoolBySourceToken(IERC20(original.tokenAmounts[i].token));
      destTokenAmounts[i].token = address(pool.getToken());
      destTokenAmounts[i].amount = original.tokenAmounts[i].amount;
    }

    return
      Client.Any2EVMMessage({
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
    return _generateAny2EVMMessage(sequenceNumber, new Client.EVMTokenAmount[](0));
  }

  function _generateAny2EVMMessageWithTokens(
    uint64 sequenceNumber,
    uint256[] memory amounts
  ) internal view returns (Internal.EVM2EVMMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].amount = amounts[i];
    }
    return _generateAny2EVMMessage(sequenceNumber, tokenAmounts);
  }

  function _generateAny2EVMMessage(
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
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      receiver: address(s_receiver),
      data: data,
      tokenAmounts: tokenAmounts,
      sourceTokenData: new bytes[](tokenAmounts.length),
      feeToken: s_destFeeToken,
      feeTokenAmount: uint256(0),
      messageId: ""
    });
    message.messageId = Internal._hash(
      message,
      keccak256(
        abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, ON_RAMP_ADDRESS)
      )
    );

    return message;
  }

  function _generateBasicMessages() internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](1);
    messages[0] = _generateAny2EVMMessageNoTokens(1);
    return messages;
  }

  function _generateMessagesWithTokens() internal view returns (Internal.EVM2EVMMessage[] memory) {
    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](2);
    Client.EVMTokenAmount[] memory tokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = 1e18;
    tokenAmounts[1].amount = 5e18;
    messages[0] = _generateAny2EVMMessage(1, tokenAmounts);
    messages[1] = _generateAny2EVMMessage(2, tokenAmounts);
    return messages;
  }

  function _generateReportFromMessages(
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (Internal.ExecutionReport memory) {
    bytes[][] memory offchainTokenData = new bytes[][](messages.length);

    for (uint256 i = 0; i < messages.length; ++i) {
      offchainTokenData[i] = new bytes[](messages[i].tokenAmounts.length);
    }

    return
      Internal.ExecutionReport({
        proofs: new bytes32[](0),
        proofFlagBits: 2 ** 256 - 1,
        messages: messages,
        offchainTokenData: offchainTokenData
      });
  }

  function _getGasLimitsFromMessages(
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (uint256[] memory) {
    uint256[] memory gasLimits = new uint256[](messages.length);
    for (uint256 i = 0; i < messages.length; ++i) {
      gasLimits[i] = messages[i].gasLimit;
    }

    return gasLimits;
  }

  function _assertSameConfig(EVM2EVMOffRamp.DynamicConfig memory a, EVM2EVMOffRamp.DynamicConfig memory b) public {
    assertEq(a.permissionLessExecutionThresholdSeconds, b.permissionLessExecutionThresholdSeconds);
    assertEq(a.router, b.router);
    assertEq(a.priceRegistry, b.priceRegistry);
    assertEq(a.maxNumberOfTokensPerMsg, b.maxNumberOfTokensPerMsg);
    assertEq(a.maxDataBytes, b.maxDataBytes);
    assertEq(a.maxPoolReleaseOrMintGas, b.maxPoolReleaseOrMintGas);
  }
}
