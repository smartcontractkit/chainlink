// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {NonceManager} from "../../../NonceManager.sol";
import {Router} from "../../../Router.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampWithMessageTransformer} from "../../../onRamp/OnRampWithMessageTransformer.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

contract OnRampWithMessageTransformer_executeSingleMessage is OnRampSetup {
  OnRampWithMessageTransformer internal s_onRampWithMessageTransformer;
  MessageTransformerHelper internal s_messageTransformer;

  function setUp() public virtual override {
    super.setUp();
    s_messageTransformer = new MessageTransformerHelper();
    s_onRampWithMessageTransformer = new OnRampWithMessageTransformer(
      s_onRamp.getStaticConfig(),
      s_onRamp.getDynamicConfig(),
      _generateDestChainConfigArgs(s_sourceRouter),
      address(s_messageTransformer)
    );
    s_metadataHash = keccak256(
      abi.encode(
        Internal.EVM_2_ANY_MESSAGE_HASH,
        SOURCE_CHAIN_SELECTOR,
        DEST_CHAIN_SELECTOR,
        address(s_onRampWithMessageTransformer)
      )
    );
    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(s_onRampWithMessageTransformer);

    NonceManager(s_outboundNonceManager).applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] =
      Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(s_onRampWithMessageTransformer)});

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: makeAddr("offRamp0")});
    offRampUpdates[1] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: makeAddr("offRamp1")});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
    vm.startPrank(address(s_sourceRouter));
  }

  function test_forwardFromRouter() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // transformedMessage is used to assert that message was transformed as expected by IMessageTransformer
    Client.EVM2AnyMessage memory transformedMessage = _generateEmptyMessage();
    transformedMessage.data = abi.encodePacked("transformedData", transformedMessage.data);
    Internal.EVM2AnyRampMessage memory messageSentLog = _messageToEvent(transformedMessage, 1, 1, 0, OWNER);
    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, messageSentLog);
    s_onRampWithMessageTransformer.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_forwardFromRouter_RevertWhen_UnknownChain() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Fail with any error (UnknownChain in this case) to check if OnRamp wraps the error with MessageTransformError during the revert
    s_messageTransformer.setShouldRevert(true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageTransformer.MessageTransformError.selector,
        abi.encodeWithSelector(MessageTransformerHelper.UnknownChain.selector)
      )
    );
    s_onRampWithMessageTransformer.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }
}
