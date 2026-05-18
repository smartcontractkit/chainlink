// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

/// @title IRouter - delivers keystone reports to receiver
interface IRouter {
  error UnauthorizedForwarder();
  /// @dev Thrown when the gas limit is insufficient for handling state after
  /// calling the receiver function.
  error InsufficientGasForRouting(bytes32 transmissionId);
  error AlreadyAttempted(bytes32 transmissionId);

  event ForwarderAdded(address indexed forwarder);
  event ForwarderRemoved(address indexed forwarder);

  enum TransmissionState {
    NOT_ATTEMPTED,
    SUCCEEDED,
    INVALID_RECEIVER,
    FAILED
  }

  struct TransmissionInfo {
    bytes32 transmissionId;
    TransmissionState state;
    address transmitter;
    // This is true if the receiver is not a contract or does not implement the
    // `IReceiver` interface.
    bool invalidReceiver;
    // Whether the transmission attempt was successful. If `false`, the
    // transmission can be retried with an increased gas limit.
    bool success;
    // The amount of gas allocated for the `IReceiver.onReport` call. uint80
    // allows storing gas for known EVM block gas limits.
    // Ensures that the minimum gas requested by the user is available during
    // the transmission attempt. If the transmission fails (indicated by a
    // `false` success state), it can be retried with an increased gas limit.
    uint80 gasLimit;
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

  function getTransmissionId(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external pure returns (bytes32);
  function getTransmissionInfo(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (TransmissionInfo memory);
  function getTransmitter(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (address);
}
