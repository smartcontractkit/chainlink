// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageInterceptor} from "../../../interfaces/IMessageInterceptor.sol";
import {IRouter} from "../../../interfaces/IRouter.sol";

import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {Pool} from "../../../libraries/Pool.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {MaybeRevertMessageReceiverNo165} from "../../helpers/receivers/MaybeRevertMessageReceiverNo165.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_executeSingleMessage is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
    vm.startPrank(address(s_offRamp));
  }

  function test_executeSingleMessage_NoTokens() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);

    Client.Any2EVMMessage memory expectedAny2EvmMessage = Client.Any2EVMMessage({
      messageId: message.header.messageId,
      sourceChainSelector: message.header.sourceChainSelector,
      sender: message.sender,
      data: message.data,
      destTokenAmounts: new Client.EVMTokenAmount[](0)
    });
    vm.expectCall(
      address(s_destRouter),
      abi.encodeWithSelector(
        IRouter.routeMessage.selector,
        expectedAny2EvmMessage,
        GAS_FOR_CALL_EXACT_CHECK,
        message.gasLimit,
        message.receiver
      )
    );
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_WithTokens() public {
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

  function test_executeSingleMessage_WithMessageInterceptor() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    _enableInboundMessageInterceptor();
    vm.startPrank(address(s_offRamp));
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_NonContract() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    message.receiver = STRANGER;
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_NonContractWithTokens() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);

    message.receiver = STRANGER;

    vm.expectEmit();
    emit TokenPool.Released(address(s_offRamp), STRANGER, amounts[0]);
    vm.expectEmit();
    emit TokenPool.Minted(address(s_offRamp), STRANGER, amounts[1]);

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  // Reverts

  function test_RevertWhen_executeSingleMessageWhen_TokenHandlingError() public {
    Internal.Any2EVMRampMessage memory message = _generateAny2EVMMessageWithMaybeRevertingSingleToken(1, 50);
    address destPool = s_destPoolByToken[message.tokenAmounts[0].destTokenAddress];

    bytes memory errorMessage = "Random token pool issue";

    s_maybeRevertingPool.setShouldRevert(errorMessage);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, destPool, errorMessage));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_RevertWhen_executeSingleMessageWhen_ZeroGasDONExecution() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    message.gasLimit = 0;

    vm.expectRevert(abi.encodeWithSelector(OffRamp.ReceiverError.selector, ""));

    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_RevertWhen_executeSingleMessageWhen_MessageSender() public {
    vm.stopPrank();
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    vm.expectRevert(OffRamp.CanOnlySelfCall.selector);
    s_offRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_RevertWhen_executeSingleMessageWhen_MessageValidationError() public {
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

  function test_RevertWhen_executeSingleMessageWhen_WithFailingValidationNoRouterCall() public {
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
