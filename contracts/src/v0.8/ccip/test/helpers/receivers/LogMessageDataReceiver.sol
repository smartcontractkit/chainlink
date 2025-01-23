// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../../../shared/interfaces/ITypeAndVersion.sol";
import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";

import {Client} from "../../../libraries/Client.sol";

import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

/// @dev A contract that logs the data of a CCIP message received
contract LogMessageDataReceiver is IAny2EVMMessageReceiver, ITypeAndVersion, IERC165 {
  event MessageReceived(bytes data);

  string public constant override typeAndVersion = "LogMessageDataReceiver 1.0.0";

  /// @notice IERC165 supports an interfaceId
  /// @param interfaceId The interfaceId to check
  /// @return true if the interfaceId is supported
  function supportsInterface(
    bytes4 interfaceId
  ) public pure override returns (bool) {
    return interfaceId == type(IAny2EVMMessageReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }

  /// @dev Logs the data of the message received
  /// @param message The message received
  function ccipReceive(
    Client.Any2EVMMessage calldata message
  ) external override {
    emit MessageReceived(message.data);
  }
}
