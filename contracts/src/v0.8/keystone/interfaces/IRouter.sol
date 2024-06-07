// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title IRouter - delivers keystone reports to receiver
interface IRouter {
  enum TransmissionState {
    NOT_ATTEMPTED,
    SUCCEEDED,
    FAILED
  }

  function route(
    bytes32 transmissionId,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata report
  ) external returns (bool);

  function getTransmitter(bytes32 transmissionId) external view returns (address);
  function getTransmissionState(bytes32 transmissionId) external view returns (TransmissionState);
}
