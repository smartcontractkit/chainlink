// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";

import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampWithMessageTransformer} from "../../../offRamp/OffRampWithMessageTransformer.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRampWithMessageTransformer_setMessageTransformer is OffRampSetup {
  OffRampWithMessageTransformer internal s_offRampWithMessageTransformer;
  IMessageTransformer internal s_inboundMessageTransformer;

  function setUp() public virtual override {
    super.setUp();
    s_inboundMessageTransformer = new MessageTransformerHelper();
    s_offRampWithMessageTransformer = new OffRampWithMessageTransformer(
      s_offRamp.getStaticConfig(),
      s_offRamp.getDynamicConfig(),
      new OffRamp.SourceChainConfigArgs[](0),
      address(s_inboundMessageTransformer)
    );
  }

  function test_setMessageTransformer() public {
    assertEq(s_offRampWithMessageTransformer.getMessageTransformer(), address(s_inboundMessageTransformer));
    IMessageTransformer newMessageTransformer = new MessageTransformerHelper();
    s_offRampWithMessageTransformer.setMessageTransformer(address(newMessageTransformer));
    assertEq(s_offRampWithMessageTransformer.getMessageTransformer(), address(newMessageTransformer));
    assertNotEq(address(s_inboundMessageTransformer), address(newMessageTransformer));
  }

  function test_setMessageTransformer_RevertWhen_ZeroAddress() public {
    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRampWithMessageTransformer.setMessageTransformer(address(0));
  }
}
