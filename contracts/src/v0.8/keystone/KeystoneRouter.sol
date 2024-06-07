// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouter} from "./interfaces/IRouter.sol";
import {IReceiver} from "./interfaces/IReceiver.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

contract KeystoneRouter is IRouter, OwnerIsCreator, ITypeAndVersion {
  error Unauthorized();

  /// @notice This error is thrown whenever a message has already been processed.
  /// @param messageId The ID of the message that was already processed
  error AlreadyProcessed(bytes32 messageId);

  event ForwarderAdded(address indexed forwarder);
  event ForwarderRemoved(address indexed forwarder);

  mapping(address forwarder => bool) internal s_forwarders;
  mapping(bytes32 reportId => DeliveryStatus status) internal s_reports;

  string public constant override typeAndVersion = "KeystoneRouter 1.0.0";

  struct DeliveryStatus {
    address transmitter;
    bool state;
  }

  function addForwarder(address forwarder) external onlyOwner {
    s_forwarders[forwarder] = true;
    emit ForwarderAdded(forwarder);
  }

  function removeForwarder(address forwarder) external onlyOwner {
    s_forwarders[forwarder] = false;
    emit ForwarderRemoved(forwarder);
  }

  function route(
    bytes32 id,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata report
  ) external returns (bool) {
    if (!s_forwarders[msg.sender]) {
      revert Unauthorized();
    }

    if (s_reports[id].transmitter != address(0)) revert AlreadyProcessed(id);
    s_reports[id].transmitter = transmitter;

    bool success;
    try IReceiver(receiver).onReport(metadata, report) {
      success = true;
      s_reports[id].state = true;
    } catch {
      // Do nothing, success is already false
    }
    return success;
  }

  // @notice Get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(bytes32 id) external view returns (address) {
    return s_reports[id].transmitter;
  }

  // @notice Get delivery status of a given report
  function getDeliveryStatus(bytes32 id) external view returns (bool) {
    return s_reports[id].transmitter != address(0);
  }
}
