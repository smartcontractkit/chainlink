// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.26;

import {IMessageTransformer} from "../../interfaces/IMessageTransformer.sol";

import {Internal} from "../../libraries/Internal.sol";

// @dev This helper is used to test the On/OffRamps
contract MessageTransformerHelper is IMessageTransformer {
  error UnknownChain();

  bool public s_shouldRevert;

  /// @dev Set whether the transformer should revert
  function setShouldRevert(
    bool _shouldRevert
  ) external {
    s_shouldRevert = _shouldRevert;
  }

  /// @inheritdoc IMessageTransformer
  function transformOutboundMessage(
    Internal.EVM2AnyRampMessage memory message
  ) public view returns (Internal.EVM2AnyRampMessage memory) {
    if (s_shouldRevert) {
      revert UnknownChain();
    }
    message.data = abi.encodePacked("transformedData", message.data);
    return message;
  }

  /// @inheritdoc IMessageTransformer
  function transformInboundMessage(
    Internal.Any2EVMRampMessage memory message
  ) public view returns (Internal.Any2EVMRampMessage memory) {
    if (s_shouldRevert) {
      revert UnknownChain();
    }
    message.data = abi.encodePacked("transformedData", message.data);
    return message;
  }
}
