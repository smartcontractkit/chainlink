// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IRouter} from "./interfaces/IRouter.sol";
import {IReceiver} from "./interfaces/IReceiver.sol";

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

contract KeystoneRouter is IRouter, OwnerIsCreator, ITypeAndVersion {
  error Unauthorized();
  error AlreadyAttempted(bytes32 transmissionId);

  event ForwarderAdded(address indexed forwarder);
  event ForwarderRemoved(address indexed forwarder);

  mapping(address forwarder => bool) internal s_forwarders;
  mapping(bytes32 transmissionId => TransmissionInfo) internal s_transmissions;

  string public constant override typeAndVersion = "KeystoneRouter 1.0.0";

  struct TransmissionInfo {
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
    bytes32 transmissionId,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata report
  ) external returns (bool) {
    if (!s_forwarders[msg.sender]) {
      revert Unauthorized();
    }

    if (s_transmissions[transmissionId].transmitter != address(0)) revert AlreadyAttempted(transmissionId);
    s_transmissions[transmissionId].transmitter = transmitter;

    if (receiver.code.length == 0) return false;

    try IReceiver(receiver).onReport(metadata, report) {
      s_transmissions[transmissionId].state = true;
      return true;
    } catch {
      return false;
    }
  }

  /// @notice Get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(bytes32 transmissionId) external view returns (address) {
    return s_transmissions[transmissionId].transmitter;
  }

  /// @notice Get delivery status of a given report
  function getTransmissionState(bytes32 transmissionId) external view returns (IRouter.TransmissionState) {
    if (s_transmissions[transmissionId].transmitter == address(0)) return IRouter.TransmissionState.NOT_ATTEMPTED;
    return
      s_transmissions[transmissionId].state ? IRouter.TransmissionState.SUCCEEDED : IRouter.TransmissionState.FAILED;
  }

  function isForwarder(address forwarder) external view returns (bool) {
    return s_forwarders[forwarder];
  }
}
