// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {Router} from "../../../Router.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";
import {OffRampWithMessageTransformerHelper} from "../../helpers/OffRampWithMessageTransformerHelper.sol";
import {LogMessageDataReceiver} from "../../helpers/receivers/LogMessageDataReceiver.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRampWithMessageTransformer_executeSingleReport is OffRampSetup {
  OffRampWithMessageTransformerHelper internal s_offRampWithMessageTransformer;
  MessageTransformerHelper internal s_inboundMessageTransformer;

  function setUp() public virtual override {
    super.setUp();
    s_inboundMessageTransformer = new MessageTransformerHelper();
    s_offRampWithMessageTransformer = new OffRampWithMessageTransformerHelper(
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
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    // set up off ramp with message transformer from configs
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

    // override for exec testing
    s_offRampWithMessageTransformer.setVerifyOverrideResult(SOURCE_CHAIN_SELECTOR_1, 1);

    // set nonce manager authorized callers
    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(s_offRampWithMessageTransformer);
    s_inboundNonceManager.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );
    vm.startPrank(address(s_offRampWithMessageTransformer));
  }

  function test_executeSingleReport() public {
    s_receiver = new LogMessageDataReceiver();
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    bytes memory data = abi.encode(0);
    assertEq(messages[0].data, data);
    vm.expectEmit();
    emit LogMessageDataReceiver.MessageReceived(abi.encodePacked("transformedData", data));
    s_offRampWithMessageTransformer.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }

  function test_executeSingleReport_RevertWhen_UnknownChain() public {
    Internal.Any2EVMRampMessage[] memory messages =
      _generateSingleBasicMessage(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1);
    // Fail with any error (UnknownChain in this case) to check if OffRamp wraps the error with MessageTransformError during the revert
    s_inboundMessageTransformer.setShouldRevert(true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageTransformer.MessageTransformError.selector,
        abi.encodeWithSelector(MessageTransformerHelper.UnknownChain.selector)
      )
    );
    s_offRampWithMessageTransformer.executeSingleReport(
      _generateReportFromMessages(SOURCE_CHAIN_SELECTOR_1, messages), new OffRamp.GasLimitOverride[](0)
    );
  }
}
