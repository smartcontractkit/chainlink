// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";

import {OnRampWithMessageTransformer} from "../../../onRamp/OnRampWithMessageTransformer.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

contract OnRampWithMessageTransformer_setMessageTransformer is OnRampSetup {
  OnRampWithMessageTransformer internal s_onRampWithMessageTransformer;
  MessageTransformerHelper internal s_outboundMessageTransformer;

  function setUp() public virtual override {
    super.setUp();
    s_outboundMessageTransformer = new MessageTransformerHelper();
    s_onRampWithMessageTransformer = new OnRampWithMessageTransformer(
      s_onRamp.getStaticConfig(),
      s_onRamp.getDynamicConfig(),
      _generateDestChainConfigArgs(s_sourceRouter),
      address(s_outboundMessageTransformer)
    );
  }

  function test_setMessageTransformer() public {
    assertEq(s_onRampWithMessageTransformer.getMessageTransformer(), address(s_outboundMessageTransformer));
    IMessageTransformer newMessageTransformer = new MessageTransformerHelper();
    s_onRampWithMessageTransformer.setMessageTransformer(address(newMessageTransformer));
    assertEq(s_onRampWithMessageTransformer.getMessageTransformer(), address(newMessageTransformer));
    assertNotEq(address(s_outboundMessageTransformer), address(newMessageTransformer));
  }

  function test_setMessageTransformer_RevertWhen_ZeroAddress() public {
    vm.expectRevert(OnRampWithMessageTransformer.ZeroAddressNotAllowed.selector);
    s_onRampWithMessageTransformer.setMessageTransformer(address(0));
  }
}
