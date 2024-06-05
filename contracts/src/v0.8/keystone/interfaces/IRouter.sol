// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title IRouter - delivers keystone reports to receiver
interface IRouter {
  function route(
    bytes32 id,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata report
  ) external returns (bool);

  function getTransmitter(bytes32 id) external view returns (address);
  function getDeliveryStatus(bytes32 id) external view returns (bool);
}
