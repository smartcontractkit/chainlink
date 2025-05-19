// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../interfaces/IMessageTransformer.sol";

import {Internal} from "../libraries/Internal.sol";
import {OffRamp} from "./OffRamp.sol";

/// @notice OffRamp that uses a message transformer to transform messages before execution
contract OffRampWithMessageTransformer is OffRamp {
  address internal s_messageTransformer;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    SourceChainConfigArgs[] memory sourceChainConfigs,
    address messageTransformerAddr
  ) OffRamp(staticConfig, dynamicConfig, sourceChainConfigs) {
    if (messageTransformerAddr == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    s_messageTransformer = messageTransformerAddr;
  }

  /// @notice Get the address of the message transformer
  /// @return messageTransformerAddr The address of the message transformer
  function getMessageTransformer() external view returns (address) {
    return s_messageTransformer;
  }

  /// @notice Set the address of the message transformer
  /// @param messageTransformerAddr The address of the message transformer
  function setMessageTransformer(
    address messageTransformerAddr
  ) external onlyOwner {
    if (messageTransformerAddr == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    s_messageTransformer = messageTransformerAddr;
  }

  /// @inheritdoc OffRamp
  function _beforeExecuteSingleMessage(
    Internal.Any2EVMRampMessage memory message
  ) internal override returns (Internal.Any2EVMRampMessage memory) {
    try IMessageTransformer(s_messageTransformer).transformInboundMessage(message) returns (
      Internal.Any2EVMRampMessage memory transformedMessage
    ) {
      return transformedMessage;
    } catch (bytes memory err) {
      revert IMessageTransformer.MessageTransformError(err);
    }
  }
}
