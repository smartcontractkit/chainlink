// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title IRouter - delivers keystone reports to receiver
interface IRouter {
  error UnauthorizedForwarder();
  error AlreadyAttempted(bytes32 transmissionId);

  event ForwarderAdded(address indexed forwarder);
  event ForwarderRemoved(address indexed forwarder);

  enum TransmissionState {
    NOT_ATTEMPTED,
    SUCCEEDED,
    FAILED
  }

  struct TransmissionInfo {
    address transmitter;
    bool success;
  }

  function addForwarder(address forwarder) external;
  function removeForwarder(address forwarder) external;

  function route(
    bytes32 transmissionId,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata report
  ) external returns (bool);

  function getTransmitter(bytes32 transmissionId) external view returns (address);
  function getTransmissionState(bytes32 transmissionId) external view returns (TransmissionState);
  function isForwarder(address forwarder) external view returns (bool);
}
