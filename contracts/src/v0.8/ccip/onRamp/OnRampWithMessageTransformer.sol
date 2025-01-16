// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../interfaces/IMessageTransformer.sol";
import {Internal} from "../libraries/Internal.sol";
import {OnRamp} from "./OnRamp.sol";

/// @notice OnRamp that uses a message transformer to transform messages from router
contract OnRampWithMessageTransformer is OnRamp {
  address internal s_messageTransformer;

  error ZeroAddressNotAllowed();

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    DestChainConfigArgs[] memory destChainConfigs,
    address messageTransformerAddr
  ) OnRamp(staticConfig, dynamicConfig, destChainConfigs) {
    if (address(messageTransformerAddr) == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    s_messageTransformer = messageTransformerAddr;
  }

  function getMessageTransformer() external view returns (address) {
    return s_messageTransformer;
  }

  function setMessageTransformer(
    address messageTransformerAddr
  ) external onlyOwner {
    if (address(messageTransformerAddr) == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    s_messageTransformer = messageTransformerAddr;
  }

  function _postProcessMessage(
    Internal.EVM2AnyRampMessage memory message
  ) internal override returns (Internal.EVM2AnyRampMessage memory) {
    try IMessageTransformer(s_messageTransformer).transformOutboundMessage(message) returns (
      Internal.EVM2AnyRampMessage memory transformedMessage
    ) {
      return transformedMessage;
    } catch (bytes memory err) {
      revert IMessageTransformer.MessageTransformError(err);
    }
  }
}
