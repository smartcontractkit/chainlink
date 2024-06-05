// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouter} from "./interfaces/IRouter.sol";
import {IReceiver} from "./interfaces/IReceiver.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";

contract KeystoneRouter is IRouter, ConfirmedOwner, TypeAndVersionInterface {
  error Unauthorized();

  /// @notice This error is thrown whenever a message has already been processed.
  /// @param messageId The ID of the message that was already processed
  error AlreadyProcessed(bytes32 messageId);

  mapping(address forwarder => bool) internal s_forwarders;
  mapping(bytes32 reportId => DeliveryStatus status) internal s_reports;

  constructor() ConfirmedOwner(msg.sender) {}

  struct DeliveryStatus {
    address transmitter;
    bool success;
  }

  function addForwarder(address forwarder) external onlyOwner {
    // TODO: events
    s_forwarders[forwarder] = true;
  }

  function removeForwarder(address forwarder) external onlyOwner {
    // TODO: events
    s_forwarders[forwarder] = false;
  }

  function route(
    bytes32 id,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata report
  ) external returns (bool) {
    if (!s_forwarders[msg.sender]) { revert Unauthorized(); }

    if (s_reports[id].transmitter != address(0)) revert AlreadyProcessed(id);

    bool success;
    try IReceiver(receiver).onReport(metadata, report) {
      success = true;
    } catch {
      // Do nothing, success is already false
    }
    s_reports[id] = DeliveryStatus(transmitter, success);
    return success;
  }

  // get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(bytes32 id) external view returns (address) {
    return s_reports[id].transmitter;
  }

  // get delivery status of a given report
  function getDeliveryStatus(bytes32 id) external view returns (bool) {
    return s_reports[id].transmitter != address(0);
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "KeystoneRouter 1.0.0";
  }
}
