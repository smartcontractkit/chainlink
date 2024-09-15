// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IMessageInterceptor} from "../../interfaces/IMessageInterceptor.sol";
import {Client} from "../../libraries/Client.sol";

contract MessageInterceptorHelper is IMessageInterceptor {
  mapping(bytes32 messageId => bool isInvalid) internal s_invalidMessageIds;

  constructor() {}

  function setMessageIdValidationState(bytes32 messageId, bool isInvalid) external {
    s_invalidMessageIds[messageId] = isInvalid;
  }

  /// @inheritdoc IMessageInterceptor
  function onInboundMessage(Client.Any2EVMMessage memory message) external view {
    if (s_invalidMessageIds[message.messageId]) {
      revert MessageValidationError(bytes("Invalid message"));
    }
  }

  /// @inheritdoc IMessageInterceptor
  function onOutboundMessage(uint64, Client.EVM2AnyMessage calldata message) external view {
    if (s_invalidMessageIds[keccak256(abi.encode(message))]) {
      revert MessageValidationError(bytes("Invalid message"));
    }
    return;
  }
}
