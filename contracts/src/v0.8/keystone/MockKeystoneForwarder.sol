// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "./interfaces/IReceiver.sol";
import {IRouter} from "./interfaces/IRouter.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

/// @notice Simplified mock version of KeystoneForwarder for testing purposes.
/// The report function is permissionless and skips all validations.
contract MockKeystoneForwarder is OwnerIsCreator, ITypeAndVersion, IRouter {
  /// @notice This error is returned when the report is shorter than REPORT_METADATA_LENGTH,
  /// which is the minimum length of a report.
  error InvalidReport();

  struct Transmission {
    address transmitter;
    // This is true if the receiver is not a contract or does not implement the `IReceiver` interface.
    bool invalidReceiver;
    // Whether the transmission attempt was successful. If `false`, the transmission can be retried
    // with an increased gas limit.
    bool success;
    // The amount of gas allocated for the `IReceiver.onReport` call. uint80 allows storing gas for known EVM block
    // gas limits. Ensures that the minimum gas requested by the user is available during the transmission attempt.
    // If the transmission fails (indicated by a `false` success state), it can be retried with an increased gas limit.
    uint80 gasLimit;
  }

  /// @notice Emitted when a report is processed
  /// @param result The result of the attempted delivery. True if successful.
  event ReportProcessed(
    address indexed receiver,
    bytes32 indexed workflowExecutionId,
    bytes2 indexed reportId,
    bool result
  );

  string public constant override typeAndVersion = "MockKeystoneForwarder 1.0.0";

  constructor() OwnerIsCreator() {
    s_forwarders[address(this)] = true;
  }

  uint256 internal constant METADATA_LENGTH = 109;
  uint256 internal constant FORWARDER_METADATA_LENGTH = 45;

  /// @dev This is the gas required to store `success` after the report is processed.
  /// It is a warm storage write because of the packed struct. In practice it will cost less.
  uint256 internal constant INTERNAL_GAS_REQUIREMENTS_AFTER_REPORT = 5_000;
  /// @dev This is the gas required to store the transmission struct and perform other checks.
  uint256 internal constant INTERNAL_GAS_REQUIREMENTS = 25_000 + INTERNAL_GAS_REQUIREMENTS_AFTER_REPORT;
  /// @dev This is the minimum gas required to route a report. This includes internal gas requirements
  /// as well as the minimum gas that the user contract will receive. 30k * 3 gas is to account for
  /// cases where consumers need close to the 30k limit provided in the supportsInterface check.
  uint256 internal constant MINIMUM_GAS_LIMIT = INTERNAL_GAS_REQUIREMENTS + 30_000 * 3 + 10_000;

  // ================================================================
  // │                          Router                              │
  // ================================================================

  mapping(address forwarder => bool isForwarder) internal s_forwarders;
  mapping(bytes32 transmissionId => Transmission transmission) internal s_transmissions;

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
    bytes calldata validatedReport
  ) public returns (bool) {
    s_transmissions[transmissionId].transmitter = transmitter;
    s_transmissions[transmissionId].gasLimit = uint80(gasleft());

    // Always call onReport on the receiver
    bool success;
    bytes memory payload = abi.encodeCall(IReceiver.onReport, (metadata, validatedReport));

    assembly {
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gas(), receiver, 0, add(payload, 0x20), mload(payload), 0x0, 0x0)
    }

    s_transmissions[transmissionId].success = success;
    return success;
  }

  function getTransmissionId(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) public pure returns (bytes32) {
    // This is slightly cheaper compared to `keccak256(abi.encode(receiver, workflowExecutionId, reportId));`
    return keccak256(bytes.concat(bytes20(uint160(receiver)), workflowExecutionId, reportId));
  }

  function getTransmissionInfo(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (TransmissionInfo memory) {
    bytes32 transmissionId = getTransmissionId(receiver, workflowExecutionId, reportId);

    Transmission memory transmission = s_transmissions[transmissionId];

    TransmissionState state;

    if (transmission.transmitter == address(0)) {
      state = IRouter.TransmissionState.NOT_ATTEMPTED;
    } else if (transmission.invalidReceiver) {
      state = IRouter.TransmissionState.INVALID_RECEIVER;
    } else {
      state = transmission.success ? IRouter.TransmissionState.SUCCEEDED : IRouter.TransmissionState.FAILED;
    }

    return
      TransmissionInfo({
        gasLimit: transmission.gasLimit,
        invalidReceiver: transmission.invalidReceiver,
        state: state,
        success: transmission.success,
        transmissionId: transmissionId,
        transmitter: transmission.transmitter
      });
  }

  /// @notice Get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (address) {
    return s_transmissions[getTransmissionId(receiver, workflowExecutionId, reportId)].transmitter;
  }

  function isForwarder(address forwarder) external view returns (bool) {
    return s_forwarders[forwarder];
  }

  // ================================================================
  // │                          Forwarder                          │
  // ================================================================

  /// @notice Simplified permissionless report function that skips all validations
  /// and does not call onReport on consumer contracts
  function report(
    address receiver,
    bytes calldata rawReport,
    bytes calldata reportContext,
    bytes[] calldata signatures
  ) external {
    if (rawReport.length < METADATA_LENGTH) {
      revert InvalidReport();
    }

    bytes32 workflowExecutionId;
    bytes2 reportId;
    {
      uint64 configId;
      (workflowExecutionId, configId, reportId) = _getMetadata(rawReport);
    }

    // Skip all validations and signature checks
    // Skip onReport call to consumer contracts
    bool success = this.route(
      getTransmissionId(receiver, workflowExecutionId, reportId),
      msg.sender,
      receiver,
      rawReport[FORWARDER_METADATA_LENGTH:METADATA_LENGTH],
      rawReport[METADATA_LENGTH:]
    );

    emit ReportProcessed(receiver, workflowExecutionId, reportId, success);
  }

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getMetadata(
    bytes memory rawReport
  ) internal pure returns (bytes32 workflowExecutionId, uint64 configId, bytes2 reportId) {
    // (first 32 bytes of memory contain length of the report)
    // version                offset  32, size  1
    // workflow_execution_id  offset  33, size 32
    // timestamp              offset  65, size  4
    // don_id                 offset  69, size  4
    // don_config_version,    offset  73, size  4
    // workflow_cid           offset  77, size 32
    // workflow_name          offset 109, size 10
    // workflow_owner         offset 119, size 20
    // report_id              offset 139, size  2
    assembly {
      workflowExecutionId := mload(add(rawReport, 33))
      // shift right by 24 bytes to get the combined don_id and don_config_version
      configId := shr(mul(24, 8), mload(add(rawReport, 69)))
      reportId := mload(add(rawReport, 139))
    }
  }
}