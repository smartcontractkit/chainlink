// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Router} from "../../../Router.sol";
import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampWithMessageTransformer} from "../../../offRamp/OffRampWithMessageTransformer.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";

import {LogMessageDataReceiver} from "../../helpers/receivers/LogMessageDataReceiver.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRampWithMessageTransformer_executeSingleMessage is OffRampSetup {
  OffRampWithMessageTransformer internal s_offRampWithMessageTransformer;
  MessageTransformerHelper internal s_inboundMessageTransformer;

  function setUp() public virtual override {
    super.setUp();
    s_inboundMessageTransformer = new MessageTransformerHelper();
    s_offRampWithMessageTransformer = new OffRampWithMessageTransformer(
      s_offRamp.getStaticConfig(),
      s_offRamp.getDynamicConfig(),
      new OffRamp.SourceChainConfigArgs[](0),
      address(s_inboundMessageTransformer)
    );

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });
    s_offRampWithMessageTransformer.applySourceChainConfigUpdates(sourceChainConfigs);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2 * sourceChainConfigs.length);

    for (uint256 i = 0; i < sourceChainConfigs.length; ++i) {
      uint64 sourceChainSelector = sourceChainConfigs[i].sourceChainSelector;

      offRampUpdates[2 * i] =
        Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(s_offRampWithMessageTransformer)});
      offRampUpdates[2 * i + 1] = Router.OffRamp({
        sourceChainSelector: sourceChainSelector,
        offRamp: s_inboundNonceManager.getPreviousRamps(sourceChainSelector).prevOffRamp
      });
    }

    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
    vm.startPrank(address(s_offRampWithMessageTransformer));
  }

  function test_executeSingleMessage() public {
    s_receiver = new LogMessageDataReceiver();
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    bytes memory data = abi.encode(0);
    assertEq(message.data, data);
    vm.expectEmit();
    emit LogMessageDataReceiver.MessageReceived(abi.encodePacked("transformedData", data));
    s_offRampWithMessageTransformer.executeSingleMessage(
      message, new bytes[](message.tokenAmounts.length), new uint32[](0)
    );
  }

  function test_executeSingleMessage_RevertWhen_UnknownChain() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    // Fail with any error (UnknownChain in this case) to check if OffRamp wraps the error with MessageTransformError during the revert
    s_inboundMessageTransformer.setShouldRevert(true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageTransformer.MessageTransformError.selector,
        abi.encodeWithSelector(MessageTransformerHelper.UnknownChain.selector)
      )
    );
    s_offRampWithMessageTransformer.executeSingleMessage(
      message, new bytes[](message.tokenAmounts.length), new uint32[](0)
    );
  }
}
